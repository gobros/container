Here lies the initial design of the gobros/container.

# Project Goals
* Lazy loaded binds. Concretes are not resolved until requested.
* Binding multiple interfaces to a single concrete.
* Binding multiple concretes to a single interface.
* Resolving a slice of concretes bound to an interface.
* Resolving the most recent concrete bound to an interface.
* Use of generics to make consumption easier.

# Proposed Methods
Below are the proposed public methods that the module will provide.

* `Bind[T any](resolver any)`

# Example
A complete example of the module in use with example types to help
illustrate usage.

```go
	type NameGiver interface {
		GiveName() string
	}
	
	type FaveNumGiver interface {
		GiveFaveNum() int
	}
	
	// Dano struct
	type Dano struct {
	}
	
	var _ NameGiver = &Dano{}
	
	func (c *Dano) GiveName() string {
		return "Dano"
	}
	
	func (c *Dano) GiveFaveNum() int {
		return 1337
	}
	
	func NewDano() *Dano {
		return &Dano{}
	}
	
	// Composition struct
	type Composition struct {
		subName NameGiver
	}
	
	var _ NameGiver = &Composition{}
	
	func (c *Composition) GiveName() string {
		return c.subName.GiveName()
	}
	
	func NewComposition(subName NameGiver) *Composition {
		return &Composition{subName: subName}
	}

	func main() {
		// Multiple interfaces to a single concrete
		container.Bind[NameGiver](NewDano)
		container.Bind[FaveNumGiver](NewDano)

		// Multiple concretes to a single interface
		container.Bind[NameGiver](NewComposition)
	}
  ```
