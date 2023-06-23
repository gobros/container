package container_test

import (
	"testing"

	"github.com/gobros/container"
	"github.com/stretchr/testify/assert"
)

func setup() {
	container.Empty()

	// Reset the struct IDs now that the container is empty
	Str1InstanceID = 0
	Str2InstanceID = 0
}

func cleanup() {
	// NTD for now
}

func TestSimpleBind(t *testing.T) {
	// Given
	setup()
	container.Bind[TestInterfaceOne](NewTestStruct1)

	// When
	str1Int1 := container.Resolve[TestInterfaceOne]()

	// Then
	assert.NotNil(t, str1Int1)
	assert.Equal(t, 1, Str1InstanceID)
	assert.Equal(t, 0, Str2InstanceID)

	str1Int1Name, str1Int1ID := str1Int1.GiveOneID()
	assert.Equal(t, TestStruct1Name, str1Int1Name)
	assert.Equal(t, Str1InstanceID, str1Int1ID)

	cleanup()
}

func TestNothingBoundResolve(t *testing.T) {
	// Given
	setup()
	// Nothing

	// When
	str1 := container.Resolve[TestInterfaceOne]()
	assert.Equal(t, 0, Str1InstanceID)
	assert.Equal(t, 0, Str2InstanceID)

	// Then
	assert.Nil(t, str1)
}

func TestNothingBoundResolveAll(t *testing.T) {
	// Given
	setup()
	// Nothing

	// When
	str1Slice := container.ResolveAll[TestInterfaceOne]()

	// Then
	assert.Empty(t, str1Slice)
	assert.Equal(t, 0, Str1InstanceID)
	assert.Equal(t, 0, Str2InstanceID)

	cleanup()
}

func TestMultipleInterfacesToOneConcrete(t *testing.T) {
	// Given
	setup()
	container.Bind[TestInterfaceOne](NewTestStruct1)
	container.Bind[TestInterfaceTwo](NewTestStruct1)

	// When
	str1Int1 := container.Resolve[TestInterfaceOne]()
	str1Int2 := container.Resolve[TestInterfaceTwo]()

	// Then
	assert.NotNil(t, str1Int1)
	assert.Equal(t, 1, Str1InstanceID)
	assert.Equal(t, 0, Str2InstanceID)

	str1Int1Name, str1Int1ID := str1Int1.GiveOneID()
	assert.Equal(t, TestStruct1Name, str1Int1Name)
	assert.Equal(t, Str1InstanceID, str1Int1ID)

	str1Int2Name, str1Int2ID := str1Int2.GiveTwoID()
	assert.Equal(t, TestStruct1Name, str1Int2Name)
	assert.Equal(t, Str1InstanceID, str1Int2ID)

	cleanup()
}

func TestOneInterfaceToMultipleConcretes(t *testing.T) {
	// Given
	setup()
	container.Bind[TestInterfaceOne](NewTestStruct1)
	container.Bind[TestInterfaceOne](NewTestStruct2)

	// When
	int1Slice := container.ResolveAll[TestInterfaceOne]()

	// Then
	assert.NotNil(t, int1Slice)
	assert.Len(t, int1Slice, 2)
	assert.Equal(t, 1, Str1InstanceID)
	assert.Equal(t, 1, Str2InstanceID)

	str1Int1 := int1Slice[0]
	str1Int1Name, str1Int1ID := str1Int1.GiveOneID()
	assert.Equal(t, TestStruct1Name, str1Int1Name)
	assert.Equal(t, Str1InstanceID, str1Int1ID)

	str2Int1 := int1Slice[1]
	str2Int1Name, str2Int1ID := str2Int1.GiveOneID()
	assert.Equal(t, TestStruct2Name, str2Int1Name)
	assert.Equal(t, Str2InstanceID, str2Int1ID)

	cleanup()
}

func TestMultipleInterfaceToMultipleConcretes(t *testing.T) {
	// Given
	setup()
	container.Bind[TestInterfaceOne](NewTestStruct1)
	container.Bind[TestInterfaceTwo](NewTestStruct1)
	container.Bind[TestInterfaceOne](NewTestStruct2)
	container.Bind[TestInterfaceTwo](NewTestStruct2)

	// When
	int1Slice := container.ResolveAll[TestInterfaceOne]()
	int2Slice := container.ResolveAll[TestInterfaceTwo]()

	// Then
	assert.NotNil(t, int1Slice)
	assert.Len(t, int1Slice, 2)
	assert.Equal(t, 1, Str1InstanceID)
	assert.Equal(t, 1, Str2InstanceID)

	str1Int1 := int1Slice[0]
	str1Int1Name, str1Int1ID := str1Int1.GiveOneID()
	assert.Equal(t, TestStruct1Name, str1Int1Name)
	assert.Equal(t, Str1InstanceID, str1Int1ID)

	str2Int1 := int1Slice[1]
	str2Int1Name, str2Int1ID := str2Int1.GiveOneID()
	assert.Equal(t, TestStruct2Name, str2Int1Name)
	assert.Equal(t, Str2InstanceID, str2Int1ID)

	str1Int2 := int2Slice[0]
	str1Int2Name, str1Int2ID := str1Int2.GiveTwoID()
	assert.Equal(t, TestStruct1Name, str1Int2Name)
	assert.Equal(t, Str1InstanceID, str1Int2ID)

	str2Int2 := int2Slice[1]
	str2Int2Name, str2Int2ID := str2Int2.GiveTwoID()
	assert.Equal(t, TestStruct2Name, str2Int2Name)
	assert.Equal(t, Str2InstanceID, str2Int2ID)

	cleanup()
}

// Test structs
type TestInterfaceOne interface {
	GiveOneID() (string, int)
}

type TestInterfaceTwo interface {
	GiveTwoID() (string, int)
}

// TestStruct1 struct - Implements TestInterfaceOne and TestInterfaceTwo
var Str1InstanceID = 0

const TestStruct1Name = "TestStruct1"

type TestStruct1 struct {
	InstanceId int
}

var _ TestInterfaceOne = &TestStruct1{}
var _ TestInterfaceTwo = &TestStruct1{}

func NewTestStruct1() *TestStruct1 {
	Str1InstanceID += 1
	return &TestStruct1{
		InstanceId: Str1InstanceID,
	}
}

func (c *TestStruct1) GiveOneID() (string, int) {
	return TestStruct1Name, c.InstanceId
}

func (c *TestStruct1) GiveTwoID() (string, int) {
	return TestStruct1Name, c.InstanceId
}

// TestStruct2 struct - Implements TestInterfaceOne and TestInterfaceTwo
var Str2InstanceID = 0

const TestStruct2Name = "TestStruct2"

type TestStruct2 struct {
	InstanceId int
}

var _ TestInterfaceOne = &TestStruct2{}
var _ TestInterfaceTwo = &TestStruct2{}

func NewTestStruct2() *TestStruct2 {
	Str2InstanceID += 1
	return &TestStruct2{
		InstanceId: Str2InstanceID,
	}
}

func (c *TestStruct2) GiveOneID() (string, int) {
	return TestStruct2Name, c.InstanceId
}

func (c *TestStruct2) GiveTwoID() (string, int) {
	return TestStruct2Name, c.InstanceId
}
