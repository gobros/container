package container

import (
	"fmt"
	"reflect"
)

var Global = &Container{
	bindingToResolver:  make(map[reflect.Type][]reflect.Value),
	resolverToConcrete: make(map[reflect.Value]any),
}

type Container struct {
	// Binds a pointer/interface to a resolver function
	bindingToResolver  map[reflect.Type][]reflect.Value
	resolverToConcrete map[reflect.Value]any
}

func EmptyInstance(container *Container) {
	container.bindingToResolver = make(map[reflect.Type][]reflect.Value)
	container.resolverToConcrete = make(map[reflect.Value]any)
}

func Bind[T any](resolver any) error {
	return BindInstance[T](Global, resolver)
}

func BindInstance[T any](container *Container, resolver any) error {
	genericType := getBindingType[T]()
	resolverType := reflect.ValueOf(resolver)

	// Ensure the resolver is valid and has a chance of functioning
	err := validateResolver(resolverType, genericType)
	if err != nil {
		return fmt.Errorf("resolver validation failed: %w", err)
	}

	var foundIdx = -1
	for idx, existingResolverType := range container.bindingToResolver[genericType] {
		if existingResolverType == resolverType {
			foundIdx = idx
			break
		}
	}

	// If the concrete type is already bound, drop it so we can re-add it to the
	// end, making it take precedence in a Resolve() call.
	if foundIdx != -1 {
		container.bindingToResolver[genericType] =
			append(container.bindingToResolver[genericType][:foundIdx],
				container.bindingToResolver[genericType][foundIdx+1:]...)
	}

	container.bindingToResolver[genericType] = append(container.bindingToResolver[genericType], resolverType)

	if _, ok := container.resolverToConcrete[resolverType]; !ok {
		container.resolverToConcrete[resolverType] = nil
	}

	return nil
}

func ResolveAll[T any]() ([]T, error) {
	return ResolveAllInstance[T](Global)
}

func ResolveAllInstance[T any](container *Container) ([]T, error) {
	bindingType := getBindingType[T]()
	retVal, err := resolveAllInstanceInternal(bindingType, container)
	if err != nil {
		return nil, err
	}
	return retVal.([]T), nil
}

func Resolve[T any]() (T, error) {
	return ResolveInstance[T](Global)
}

func ResolveInstance[T any](container *Container) (T, error) {
	resolvedConcretes, err := ResolveAllInstance[T](container)
	if err != nil {
		return *new(T), err
	}

	return resolvedConcretes[len(resolvedConcretes)-1], nil
}

func resolveAllInstanceInternal(bindingType reflect.Type, container *Container) (resolvedRet any, errRet error) {
	resolved := reflect.MakeSlice(reflect.SliceOf(bindingType), 0, 0)
	resolverTypes := container.bindingToResolver[bindingType]

	if len(resolverTypes) == 0 {
		return nil, fmt.Errorf("failed to resolve for interface (%v), nothing bound", bindingType.Name())
	}

	for _, resolverValue := range resolverTypes {
		if container.resolverToConcrete[resolverValue] == nil {
			args, err := resolveArguments(container, resolverValue, bindingType)
			if err != nil {
				return nil, err
			}

			// Rare case where it's much better to handle the panic and give a descriptive error
			defer func() {
				if r := recover(); r != nil {
					returnType := resolverValue.Type().Out(0)
					resolvedRet = nil
					errRet = fmt.Errorf("failed to resolve for interface (%v) for resolver with return type (%v), encountered panic (%v) ", bindingType.Name(), returnType, r)
				}
			}()

			values := resolverValue.Call(args)

			if len(values) >= 2 && values[1].CanInterface() {
				if err, ok := values[1].Interface().(error); ok {
					return nil, fmt.Errorf("failed to resolve for interface (%v), resolver returned error: %w", bindingType.Name(), err)
				}
			}

			container.resolverToConcrete[resolverValue] = values[0].Interface()
		}

		resolved = reflect.Append(resolved, reflect.ValueOf(container.resolverToConcrete[resolverValue]))
	}

	return resolved.Interface(), nil
}

func resolveArguments(container *Container, resolverValue reflect.Value, bindingType reflect.Type) ([]reflect.Value, error) {
	resolverType := resolverValue.Type()
	argCount := resolverType.NumIn()
	args := make([]reflect.Value, argCount)

	for i := 0; i < argCount; i++ {
		argType := resolverType.In(i)
		if argType.Kind() == reflect.Slice {
			sliceType := argType.Elem()
			arg, err := resolveAllInstanceInternal(sliceType, container)
			if err != nil {
				// It's better to return an empty slice to let the resolver handle the empty slice than error
				emptySlice := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, 0)
				arg = emptySlice.Interface()
			}
			argVal := reflect.ValueOf(arg)
			args[i] = argVal
		} else {
			arg, err := resolveAllInstanceInternal(argType, container)
			if err != nil {
				return nil, fmt.Errorf("resolver dependency error, failed to resolve dependency (%v) for interface (%v): %w", argType, bindingType.Name(), err)
			}
			argVal := reflect.ValueOf(arg)
			args[i] = argVal.Index(argVal.Len() - 1)
		}
	}

	return args, nil
}

func getBindingType[T any]() reflect.Type {
	return reflect.TypeOf(new(T)).Elem()
}

func validateResolver(resolverType reflect.Value, genericType reflect.Type) error {
	if resolverType.Kind() != reflect.Func {
		return fmt.Errorf("resolver error, resolver must be a function")
	}

	if resolverType.Type().NumOut() > 0 {
		firtReturn := resolverType.Type().Out(0)
		if genericType.Kind() != reflect.Ptr && genericType.Kind() != reflect.Interface {
			return fmt.Errorf("resolver error, interface T must be a pointer or interface")
		}
		if genericType.Kind() == reflect.Interface && !firtReturn.Implements(genericType) {
			return fmt.Errorf("resolver error, resolver must return a type that implements the provided interface T")
		}
		if genericType.Kind() != reflect.Interface && !firtReturn.AssignableTo(genericType) {
			return fmt.Errorf("resolver error, resolver must return a type assignable to interface T")
		}
	} else {
		return fmt.Errorf("resolver error, resolver must return a concrete as it's first return")
	}

	for i := 0; i < resolverType.Type().NumIn(); i++ {
		paramType := resolverType.Type().In(i)

		if paramType.Kind() != reflect.Ptr &&
			paramType.Kind() != reflect.Interface &&
			paramType.Kind() != reflect.Slice {
			return fmt.Errorf("resolver error, resolver input parameters must all be of type pointer, interface, or slice")
		}
	}

	return nil
}
