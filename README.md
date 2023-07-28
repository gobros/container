# GoBros Container
<img align="left" src="mascot.png" width="200" alt="Gobros mascot, a picture of a gopher sitting in a cardboard box"/>

A flexible Inversion of Control container for GoLang that uses lazy loading
and generics to make it as easy to use as possible.

## Key Features:
* Lazy loading (bound objects aren't resolved until requested)
* Multiple interfaces to one resolver. This allows an object that satisfies 
  multiple interfaces to be bound to the container and resolved by any of it's
  bound interfaces.
* Multiple resolvers to one interface. Allows a slice of objects to be resolved
  for a given interface.
* Runtime checks when binding to help ensure the container is used correctly.

<br>

# Installing
Run the following command in your project directory to install.

`go get github.com/gobros/container@latest`

# Why Use Inversion of Control?
With Inversion of Control, dependencies can be satisfied automatically and
provided to a constructor function only when they're needed.. Without it, most
systems end up with a complex tree of dependencies that must be passed from
parent to child.

# Binding & Resolving Overview
The primary operations this container performs are binding and resolving. Below
is a short example of it in action.

```golang
// Bind our implementation of IRandomIntGenerator
container.MustBind[IRandomIntGenerator](NewFairRandomIntGenerator)

// Resolve the bound implementation of IRandomIntGenerator
generator := container.MustResolve[IRandomIntGenerator]()

// Use the generator
fmt.Printf("Random Number: %v\n", generator.Generate())
```

## Requirements To Bind
* Must bind against an Interface or Pointer type
* Resolvers must be a function
* Resolver must return a type that either implements or is assignable to the
  bound type
* Resolver may have arguments, but they must be of type Interface, Pointer, or
  Slice so the container can attemp to resolve them.

## Requirements To Resolve
* Provide a valid Interface or Pointer type to resolve

# Global Container Functions
These act upon the global container created by this module.

## Bind
Binds a resolver to an Interface or Pointer type. Can later be resolved for use.

### Definition
`Bind[T any](resolver any) error`

### Example
```golang
```

---
## ResolveAll
### Definition
`ResolveAll[T any]() ([]T, error)`

### Example
```golang
```

---
## Resolve
### Definition
`Resolve[T any]() (T, error)`

### Example
```golang
```

# Instance Container Functions
These act upon the container provided as a paramter. Can be used if you need
multiple containers don't want to use the global container provided.

* [BindInstance\[T any\](container *Container, resolver any) error](#BindInstance)
* [ResolveAllInstance\[T any\](container *Container) ([]T, error)](#ResolveAllInstance)
* [ResolveInstance\[T any\](container *Container) (T, error)](#ResolveInstance)


# Must Container Functions
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
