package container

import "reflect"

var Global = Container{
	bindings: make(map[reflect.Type][]Binding),
}

type Binding struct {
	resolver any
	concrete any
}

type Container struct {
	bindings map[reflect.Type][]Binding
}

resolverToConcrete map[any]

// TODO DJJ Here! Have to decide on global/instance naming scheme
func Bind[T any](resolver any) {
	BindInstance[T](&Global, resolver)
}

func BindInstance[T any](container *Container, resolver any) {
	// TODO DJJ Ensure resolver is a function that returns a type
	// assignable to T.
	genericType := reflect.TypeOf(new(T))
	typeBindings := container.bindings[genericType]
	newBinding := Binding{
		resolver: resolver,
	}
	container.bindings[genericType] = append(typeBindings, newBinding)
}

// TODO DJJ Needs all changed to support multiple concretes
func Resolve[T any]() []T {
	return ResolveInstance[T](&Global)
}

func ResolveInstance[T any](container *Container) []T {
	genericType := reflect.TypeOf(new(T))
	bindings := container.bindings[genericType]
	// TODO Check for error

	resolved := make([]T, len(bindings))

	for i, binding := range bindings {
		if binding.concrete == nil {
			// TODO DJJ Break out to helper and reflect parameters
			// to resolver. Attempt to resolve parameters.
			values := reflect.ValueOf(binding.resolver).Call(nil)
			binding.concrete = values[0].Interface()
		}

		resolved[i] = binding.concrete.(T)
	}

	return resolved
}
