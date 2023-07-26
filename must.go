package container

func MustBind[T any](resolver any) {
	if err := Bind[T](resolver); err != nil {
		panic(err.Error())
	}
}

func MustBindInstance[T any](container *Container, resolver any) {
	if err := BindInstance[T](container, resolver); err != nil {
		panic(err.Error())
	}
}

func MustResolveAll[T any]() []T {
	if retVal, err := ResolveAll[T](); err != nil {
		panic(err.Error())
	} else {
		return retVal
	}
}

func MustResolveAllInstance[T any](container *Container) []T {
	if retVal, err := ResolveAllInstance[T](container); err != nil {
		panic(err.Error())
	} else {
		return retVal
	}
}

func MustResolve[T any]() T {
	if retVal, err := Resolve[T](); err != nil {
		panic(err.Error())
	} else {
		return retVal
	}
}

func MustResolveInstance[T any](container *Container) T {
	if retVal, err := ResolveInstance[T](container); err != nil {
		panic(err.Error())
	} else {
		return retVal
	}
}
