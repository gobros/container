package container

import (
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

// TODO DJJ Here! Have to decide on global/instance naming scheme
func Bind[T any](resolver any) {
	BindInstance[T](&Global, resolver)
}

func BindInstance[T any](container *Container, resolver any) {
	genericType := getBindingType[T]()
	resolverType := reflect.ValueOf(resolver)

	// Ensure the resolver is valid and has a chance of functioning
	validateResolver(resolverType, genericType)

	// It's okay to attempt and bind a resolver multiple times, but it's a no-op if you do
	resolverAlreadyBound := false
	for _, existingResolverType := range container.bindingToResolver[genericType] {
		if existingResolverType == resolverType {
			resolverAlreadyBound = true
			break
		}
	}

	if !resolverAlreadyBound {
		container.bindingToResolver[genericType] = append(container.bindingToResolver[genericType], resolverType)
	}

	if _, ok := container.resolverToConcrete[resolverType]; !ok {
		container.resolverToConcrete[resolverType] = nil
	}
}

// TODO DJJ Needs all changed to support multiple concretes
func ResolveAll[T any]() []T {
	return ResolveAllInstance[T](&Global)
}

func ResolveAllInstance[T any](container *Container) []T {
	genericType := getBindingType[T]()
	// TODO Check for error

	resolverTypes := container.bindingToResolver[genericType]
	resolved := make([]T, 0)

	for _, resolverType := range resolverTypes {
		if container.resolverToConcrete[resolverType] == nil {
			// TODO DJJ Break out to helper and reflect parameters
			// to resolver. Attempt to resolve parameters.
			// TODO wirecat check for errors returns here
			values := resolverType.Call(nil)
			container.resolverToConcrete[resolverType] = values[0].Interface()
		}

		resolved = append(resolved, container.resolverToConcrete[resolverType].(T))
	}

	return resolved
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

			if paramType.Kind() != reflect.Ptr && paramType.Kind() != reflect.Interface {
				// TODO wirecat real errors
				panic("resolver input parameters must all be of type pointer or interface")
			}
		}
	}
}
