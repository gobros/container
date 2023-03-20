package main

import (
	"fmt"
	"strings"

	"github.com/gobros/container"
)

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

// Joe struct
type Joe struct {
}

var _ NameGiver = &Joe{}

func (c *Joe) GiveName() string {
	return "Joe"
}

func NewJoe() *Joe {
	return &Joe{}
}

// Composition struct
type Composition struct {
	subNameGivers []NameGiver
}

var _ NameGiver = &Composition{}

func (c *Composition) GiveName() string {
	subNames := make([]string, len(c.subNameGivers))
	for i, subNameGiver := range c.subNameGivers {
		subNames[i] = subNameGiver.GiveName()
	}
	return strings.Join(subNames, ", ")
}

func NewComposition(subNameGivers []NameGiver) *Composition {
	return &Composition{subNameGivers: subNameGivers}
}

func main() {
	// Multiple interfaces to a single concrete
	container.Bind[NameGiver](NewDano)
	container.Bind[FaveNumGiver](NewDano)

	// Multiple concretes to a single interface
	container.Bind[NameGiver](NewJoe)
	//container.Bind[NameGiver](NewComposition)

	// Resolve
	for _, concrete := range container.Resolve[NameGiver]() {
		fmt.Println(concrete.GiveName())
	}
}
