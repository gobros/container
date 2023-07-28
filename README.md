# GoBros Container
A flexible Inversion of Control container for GoLang that uses lazy loading
and generics to make it as easy to use as possible.

<img src="mascot.png" width="200" alt="Gobros mascot, a picture of a gopher sitting in a cardboard box">

# Installing
Run the following command in your project directory to install.

`go get github.com/gobros/container@latest`

# Quick Start
The following attempts to give a quick high-level overview of the container.

## Global Container Functions
These act upon the global container created by this module.

* [Bind\[T any\](resolver any) error](#Bind)
* [ResolveAll\[T any\]() ([]T, error)](#ResolveAll)
* [Resolve\[T any\]() (T, error)](#Resolve)

## Instance Container Functions
These act upon the container provided as a paramter. Can be used if you need
multiple containers don't want to use the global container provided.

* [BindInstance\[T any\](container *Container, resolver any) error](#BindInstance)
* [ResolveAllInstance\[T any\](container *Container) ([]T, error)](#ResolveAllInstance)
* [ResolveInstance\[T any\](container *Container) (T, error)](#ResolveInstance)


## Must Container Functions
For convenience, there is an alternate version of all Global and Instance Container
Functions that wrap them and panic instead of returning an error.

```golang
func MustBind[T any](resolver any)
func MustResolveAll[T any]() []T
func MustResolve[T any]() T
func MustBindInstance[T any](container *Container, resolver any)
func MustResolveAllInstance[T any](container *Container) []T
func MustResolveInstance[T any](container *Container) T
```


# Usage Guide
TODO

## Functions
TODO

### Bind
TODO

### ResolveAll
TODO

### Resolve
TODO

### BindInstance
TODO

### ResolveAllInstance
TODO

### ResolveInstance
TODO
