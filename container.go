package container

import (
	"fmt"
	"reflect"
)

var Global = Container{
	bindingToResolver:  make(map[reflect.Type][]reflect.Value),
	resolverToConcrete: make(map[reflect.Value]any),
}

type Container struct {
	// Binds a pointer/interface to a resolver function
	bindingToResolver  map[reflect.Type][]reflect.Value
	resolverToConcrete map[reflect.Value]any
}

func Empty() {
	EmptyInstance(&Global)
}

func EmptyInstance(container *Container) {
	container.bindingToResolver = make(map[reflect.Type][]reflect.Value)
	container.resolverToConcrete = make(map[reflect.Value]any)
}

// TODO DJJ Here! Have to decide on global/instance naming scheme
func Bind[T any](resolver any) {
	BindInstance[T](&Global, resolver)
}

func BindInstance[T any](container *Container, resolver any) {
	genericType := getBindingType[T]()
	resolverType := reflect.ValueOf(resolver)

	// Ensure the resolver is valid and has a chance of functioning
	validateResolver(resolverType, genericType)

	var foundIdx = -1
	for idx, existingResolverType := range container.bindingToResolver[genericType] {
		if existingResolverType == resolverType {
			foundIdx = idx
			break
		}
	}

	// If the concrete type is already bound, drop it so we can re-add it to the
	// end, making it take precedence in a single Resolve() call.
	if foundIdx != -1 {
		container.bindingToResolver[genericType] =
			append(container.bindingToResolver[genericType][:foundIdx],
				container.bindingToResolver[genericType][foundIdx+1:]...)
	}

	container.bindingToResolver[genericType] = append(container.bindingToResolver[genericType], resolverType)

	if _, ok := container.resolverToConcrete[resolverType]; !ok {
		container.resolverToConcrete[resolverType] = nil
	}
}

func ResolveAll[T any]() []T {
	return ResolveAllInstance[T](&Global)
}

func ResolveAllInstance[T any](container *Container) []T {
	bindingType := getBindingType[T]()
	return resolveAllInstanceInternal(bindingType, container).([]T)
}

func resolveAllInstanceInternal(bindingType reflect.Type, container *Container) any {
	resolved := reflect.MakeSlice(reflect.SliceOf(bindingType), 0, 0)
	resolverTypes := container.bindingToResolver[bindingType]

	for _, resolverValue := range resolverTypes {
		if container.resolverToConcrete[resolverValue] == nil {
			resolverType := resolverValue.Type()
			argCount := resolverType.NumIn()
			args := make([]reflect.Value, argCount)
			for i := 0; i < argCount; i++ {
				argType := resolverType.In(i)
				if argType.Kind() == reflect.Slice {
					sliceType := argType.Elem()
					argVal := reflect.ValueOf(resolveAllInstanceInternal(sliceType, container))
					if argVal.Len() > 0 {
						args[i] = argVal
					} else {
						panic(fmt.Sprintf("unable to resolve dependency (%v) for interface (%v)", argType, bindingType.Name()))
					}
				} else {
					argVal := reflect.ValueOf(resolveAllInstanceInternal(argType, container))
					if argVal.Len() > 0 {
						args[i] = argVal.Index(argVal.Len() - 1)
					} else {
						panic(fmt.Sprintf("unable to resolve dependency (%v) for interface (%v)", argType, bindingType.Name()))
					}
				}
			}
			// TODO wirecat handle panics to Call if possible
			values := resolverValue.Call(args)
			container.resolverToConcrete[resolverValue] = values[0].Interface()
		}

		resolved = reflect.Append(resolved, reflect.ValueOf(container.resolverToConcrete[resolverValue]))
	}
	return resolved.Interface()
}

func Resolve[T any]() T {
	return ResolveInstance[T](&Global)
}

func ResolveInstance[T any](container *Container) T {
	resolvedConcretes := ResolveAllInstance[T](container)
	if len(resolvedConcretes) > 0 {
		return resolvedConcretes[len(resolvedConcretes)-1]
	} else {
		var def T
		return def
	}
}

func getBindingType[T any]() reflect.Type {
	return reflect.TypeOf(new(T)).Elem()
}

func validateResolver(resolverType reflect.Value, genericType reflect.Type) {
	if resolverType.Kind() != reflect.Func {
		// TODO wirecat real errors
		panic("resolver must be a function")
	}
	if resolverType.Type().NumOut() > 0 {
		firtReturn := resolverType.Type().Out(0)
		if genericType.Kind() != reflect.Ptr && genericType.Kind() != reflect.Interface {
			// TODO wirecat real errors
			panic("generic must be a pointer or interface")
		}
		if genericType.Kind() == reflect.Interface && !firtReturn.Implements(genericType) {
			// TODO wirecat real errors
			panic("resolver must return a type that implements the provided interface T")
		}
		if genericType.Kind() != reflect.Interface && !firtReturn.AssignableTo(genericType) {
			// TODO wirecat real errors
			panic("resolver must return a type assignable to T")
		}
	} else {
		// TODO wirecat real errors
		panic("resolver must return a concrete as it's first return")
	}
	if resolverType.Type().NumIn() > 0 {
		for i := 0; i < resolverType.Type().NumIn(); i++ {
			paramType := resolverType.Type().In(i)

			if paramType.Kind() != reflect.Ptr &&
				paramType.Kind() != reflect.Interface &&
				paramType.Kind() != reflect.Slice {
				// TODO wirecat real errors
				panic("resolver input parameters must all be of type pointer, interface, or slice")
			}
		}
	}
}
