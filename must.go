package container

// Binds a resolver to a bound type. Can later be resolved for use. Uses the
// global container instance.
func MustBind[T any](resolver any) {
	if err := Bind[T](resolver); err != nil {
		panic(err.Error())
	}
}

// Binds a resolver to a bound type. Can later be resolved for use. Uses the
// provided container instance.
func MustBindInstance[T any](container *Container, resolver any) {
	if err := BindInstance[T](container, resolver); err != nil {
		panic(err.Error())
	}
}

// Attempts to resolve and return all concretes bound to the provided type as
// a slice. Uses the global container instance.
func MustResolveAll[T any]() []T {
	if retVal, err := ResolveAll[T](); err != nil {
		panic(err.Error())
	} else {
		return retVal
	}
}

// Attempts to resolve and return all concretes bound to the provided type as
// a slice. Uses the provided container instance.
func MustResolveAllInstance[T any](container *Container) []T {
	if retVal, err := ResolveAllInstance[T](container); err != nil {
		panic(err.Error())
	} else {
		return retVal
	}
}

// Resolves a single concrete bound to the provided type. If multiple resolvers
// were bound, the concrete from the most recent one is returned. Uses the
// global container instance.
func MustResolve[T any]() T {
	if retVal, err := Resolve[T](); err != nil {
		panic(err.Error())
	} else {
		return retVal
	}
}

// Resolves a single concrete bound to the provided type. If multiple resolvers
// were bound, the concrete from the most recent one is returned. Uses the
// provided container instance.
func MustResolveInstance[T any](container *Container) T {
	if retVal, err := ResolveInstance[T](container); err != nil {
		panic(err.Error())
	} else {
		return retVal
	}
}
