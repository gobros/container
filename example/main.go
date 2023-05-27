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

var danoInstanceNum = 0

func (c *Dano) GiveName() string {
	return "Dano" + fmt.Sprint(danoInstanceNum)
}

func (c *Dano) GiveFaveNum() int {
	return danoInstanceNum
}

func NewDano() *Dano {
	danoInstanceNum += 1
	return &Dano{}
}

// Joe struct
type Joe struct {
}

var _ NameGiver = &Joe{}

var joeInstanceNum = 0

func (c *Joe) GiveName() string {
	return "Joe" + fmt.Sprint(joeInstanceNum)
}

func NewJoe() *Joe {
	joeInstanceNum += 1
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
	container.Bind[*Dano](NewDano)
	container.Bind[NameGiver](NewDano)
	container.Bind[FaveNumGiver](NewDano)

	// Multiple concretes to a single interface
	container.Bind[NameGiver](NewJoe)
	//container.Bind[NameGiver](NewComposition)

	grabbedDano := container.Resolve[NameGiver]()
	if grabbedDano == nil {
		fmt.Println("Dano is nil")
	} else {
		fmt.Printf("Dano is not nil. Name is: %v\n", grabbedDano.GiveName())
	}

	// Resolve
	for _, concrete := range container.ResolveAll[NameGiver]() {
		fmt.Printf("Name: %v\n", concrete.GiveName())
	}

	for _, concrete := range container.ResolveAll[FaveNumGiver]() {
		fmt.Printf("FaveNum: %v\n", concrete.GiveFaveNum())
	}
}
