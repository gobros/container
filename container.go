package container

import (
	"fmt"
	"reflect"
)

var Global = &Container{
	bindingToResolver:          make(map[reflect.Type][]reflect.Value),
	resolverToConcreteInstance: make(map[reflect.Value]any),
}

type Container struct {
	// Binds a pointer/interface to a resolver function
	bindingToResolver map[reflect.Type][]reflect.Value
	// Binds a resolver function to an instantiated concrete instance
	resolverToConcreteInstance map[reflect.Value]any
}

func EmptyContainer(container *Container) {
	container.bindingToResolver = make(map[reflect.Type][]reflect.Value)
	container.resolverToConcreteInstance = make(map[reflect.Value]any)
}

// Binds a resolver to a bound type. Can later be resolved for use. Uses the
// global container instance.
func Bind[T any](resolver any) error {
	return BindInstance[T](Global, resolver)
}

// Binds a resolver to a bound type. Can later be resolved for use. Uses the
// provided container instance.
func BindInstance[T any](container *Container, resolver any) error {
	resolveReturnType := getBindingType[T]()
	resolverType := reflect.ValueOf(resolver)

	// Ensure the resolver is valid and has a chance of functioning
	err := validateResolver(resolverType, resolveReturnType)
	if err != nil {
		return fmt.Errorf("resolver validation failed: %w", err)
	}

	// If the concrete type is already bound, drop it so we can re-add it to the
	// end, making it take precedence in a Resolve() call.
	hasResolver, resolverIdx := findBoundResolver(container, resolverType, resolveReturnType)
	if hasResolver {
		container.bindingToResolver[resolveReturnType] =
			append(container.bindingToResolver[resolveReturnType][:resolverIdx],
				container.bindingToResolver[resolveReturnType][resolverIdx+1:]...)
	}

	container.bindingToResolver[resolveReturnType] = append(container.bindingToResolver[resolveReturnType], resolverType)

	// Create an entry in the resolver to concrete instance map so we can reference it safely later
	if _, ok := container.resolverToConcreteInstance[resolverType]; !ok {
		container.resolverToConcreteInstance[resolverType] = nil
	}

	return nil
}

// Attempts to resolve and return all concretes bound to the provided type as
// a slice. Uses the global container instance.
func ResolveAll[T any]() ([]T, error) {
	return ResolveAllInstance[T](Global)
}

// Attempts to resolve and return all concretes bound to the provided type as
// a slice. Uses the provided container instance.
func ResolveAllInstance[T any](container *Container) ([]T, error) {
	resolverReturnType := getBindingType[T]()

	resolvedInstance, err := resolveAllInstanceInternal(resolverReturnType, container)
	if err != nil {
		return nil, err
	}
	if arrRetVal, ok := resolvedInstance.([]T); !ok || len(arrRetVal) == 0 {
		return nil, fmt.Errorf("failed to resolve for interface (%v), nothing bound", resolverReturnType.Name())
	}

	return resolvedInstance.([]T), nil
}

// Resolves a single concrete bound to the provided type. If multiple resolvers
// were bound, the concrete from the most recent one is returned. Uses the
// global container instance.
func Resolve[T any]() (T, error) {
	return ResolveInstance[T](Global)
}

// Resolves a single concrete bound to the provided type. If multiple resolvers
// were bound, the concrete from the most recent one is returned. Uses the
// provided container instance.
func ResolveInstance[T any](container *Container) (T, error) {
	resolvedInstances, err := ResolveAllInstance[T](container)
	if err != nil {
		return *new(T), err
	}

	return resolvedInstances[len(resolvedInstances)-1], nil
}

// Shared logic for resolving all concrete instances for the given bound type
func resolveAllInstanceInternal(bindingType reflect.Type, container *Container) (resolvedRet any, errRet error) {
	resolvedInstances := reflect.MakeSlice(reflect.SliceOf(bindingType), 0, 0)
	resolvers := container.bindingToResolver[bindingType]
	var resolverReturnType reflect.Type

	// Rare case where it's much better to handle the panic and give a descriptive error
	defer func() {
		if r := recover(); r != nil {
			returnType := resolverReturnType
			resolvedRet = nil
			errRet = fmt.Errorf("failed to resolve for interface (%v) for resolver with return type (%v), encountered panic (%v) ", bindingType.Name(), returnType, r)
		}
	}()

	// Call resolvers to get a concrete instance if we don't already have one
	for _, resolver := range resolvers {
		resolverReturnType = resolver.Type().Out(0)

		if container.resolverToConcreteInstance[resolver] != nil {
			break
		}

		args, err := resolveArguments(container, resolver, bindingType)
		if err != nil {
			return nil, err
		}

		values := resolver.Call(args)

		// If we have 2 or more returns, the second return may be in an error state
		if len(values) >= 2 && values[1].Interface() != nil {
			return nil, fmt.Errorf("failed to resolve for interface (%v), resolver returned error: %w", bindingType.Name(), err)
		}

		container.resolverToConcreteInstance[resolver] = values[0].Interface()
	}

	// Add all the resolved instances to the slice for return
	for _, resolver := range resolvers {
		resolvedInstances = reflect.Append(resolvedInstances, reflect.ValueOf(container.resolverToConcreteInstance[resolver]))
	}

	return resolvedInstances.Interface(), nil
}

// Attempts to resolve all concrete instances for a resolver function's
// arguments so the resolver can be called
func resolveArguments(container *Container, resolverValue reflect.Value, bindingType reflect.Type) ([]reflect.Value, error) {
	resolverType := resolverValue.Type()
	argCount := resolverType.NumIn()
	resolvedArgs := make([]reflect.Value, argCount)

	for i := 0; i < argCount; i++ {
		argType := resolverType.In(i)
		if argType.Kind() == reflect.Slice {
			sliceType := argType.Elem()
			arg, err := resolveAllInstanceInternal(sliceType, container)
			if err != nil {
				return resolvedArgs, fmt.Errorf("resolver dependency error, failed to resolve dependency (%v) for interface (%v): %w", argType, bindingType.Name(), err)
			}
			argVal := reflect.ValueOf(arg)
			resolvedArgs[i] = argVal
		} else {
			arg, err := resolveAllInstanceInternal(argType, container)
			if err != nil {
				return nil, fmt.Errorf("resolver dependency error, failed to resolve dependency (%v) for interface (%v): %w", argType, bindingType.Name(), err)
			}

			argVal := reflect.ValueOf(arg)

			if argVal.Len() == 0 {
				return nil, fmt.Errorf("resolver dependency error, failed to resolve dependency (%v) for interface (%v)", argType, bindingType.Name())
			}

			resolvedArgs[i] = argVal.Index(argVal.Len() - 1)
		}
	}

	return resolvedArgs, nil
}

// Returns the type of the generic interface T
func getBindingType[T any]() reflect.Type {
	return reflect.TypeOf(new(T)).Elem()
}

// Validates that a resolver function is valid and can be used to resolve a type
func validateResolver(resolverType reflect.Value, genericType reflect.Type) error {
	if resolverType.Kind() != reflect.Func {
		return fmt.Errorf("resolver error, resolver must be a function")
	}

	resolverReturnCount := resolverType.Type().NumOut()
	if resolverReturnCount > 0 {
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

		// If we have 2 or more returns, the second return must be an error
		if resolverReturnCount >= 2 && !resolverType.Type().Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return fmt.Errorf("resolver error, resolvers with two or more parameters must return an error as the second parameter")
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

// Searches for a bound resolver that returns the specified type. If found,
// returns true and the index. Otherwise, returns false and -1.
func findBoundResolver(container *Container, resolverType reflect.Value, resolverReturnType reflect.Type) (found bool, idx int) {
	var foundIdx = -1
	for idx, existingResolverType := range container.bindingToResolver[resolverReturnType] {
		if existingResolverType == resolverType {
			foundIdx = idx
			return true, foundIdx
		}
	}

	return false, foundIdx
}
