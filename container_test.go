package container_test

import (
	"errors"
	"testing"

	"github.com/gobros/container"
	"github.com/stretchr/testify/assert"
)

func setup() {
	container.EmptyInstance(container.Global)

	// Reset the struct IDs now that the container is empty
	Str1InstanceNumber = 0
	Str2InstanceNumber = 0
}

func cleanup() {
	// NTD for now
}

func TestSimpleBind(t *testing.T) {
	// Given
	setup()

	container.Bind[PrimaryIDGiver](NewTestStruct1)

	// When
	str1Prim, err := container.Resolve[PrimaryIDGiver]()

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, str1Prim)
	assert.Equal(t, 1, Str1InstanceNumber)
	assert.Equal(t, 0, Str2InstanceNumber)

	primID := str1Prim.GivePrimaryID()
	assert.Equal(t, TestStruct1Name, primID.Name)
	assert.Equal(t, Str1InstanceNumber, primID.Number)

	cleanup()
}

func TestSimpleBindOverride(t *testing.T) {
	// Given
	setup()

	container.Bind[PrimaryIDGiver](NewTestStruct1)
	container.Bind[PrimaryIDGiver](NewTestStruct1)

	// When
	str1Prim, err := container.Resolve[PrimaryIDGiver]()

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, str1Prim)
	assert.Equal(t, 1, Str1InstanceNumber)
	assert.Equal(t, 0, Str2InstanceNumber)

	primID := str1Prim.GivePrimaryID()
	assert.Equal(t, TestStruct1Name, primID.Name)
	assert.Equal(t, Str1InstanceNumber, primID.Number)

	cleanup()
}

func TestSimpleBindPointer(t *testing.T) {
	//Given
	setup()

	type ptrStruct struct {
		Name       string
		FaveNumber int
	}
	testVal := ptrStruct{
		Name:       "wirecat",
		FaveNumber: 1337,
	}
	container.Bind[*ptrStruct](func() *ptrStruct {
		return &testVal
	})

	// When
	resolvedVal, err := container.Resolve[*ptrStruct]()

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, resolvedVal)
	assert.Equal(t, "wirecat", resolvedVal.Name)
	assert.Equal(t, 1337, resolvedVal.FaveNumber)

	cleanup()
}

func TestNothingBoundResolve(t *testing.T) {
	// Given
	setup()

	// When
	str1Prim, err := container.Resolve[PrimaryIDGiver]()
	assert.Equal(t, 0, Str1InstanceNumber)
	assert.Equal(t, 0, Str2InstanceNumber)

	// Then
	assert.Nil(t, str1Prim)
	assert.Error(t, err)

	cleanup()
}

func TestNothingBoundResolveAll(t *testing.T) {
	// Given
	setup()

	// When
	str1PrimSlice, err := container.ResolveAll[PrimaryIDGiver]()

	// Then
	assert.Empty(t, str1PrimSlice)
	assert.Error(t, err)
	assert.Equal(t, 0, Str1InstanceNumber)
	assert.Equal(t, 0, Str2InstanceNumber)

	cleanup()
}

func TestResolverError(t *testing.T) {
	// Given
	setup()
	container.Bind[PrimaryIDGiver](func() (*TestStruct1, error) {
		return nil, errors.New("resolver did a bad!")
	})

	// When
	val, err := container.Resolve[PrimaryIDGiver]()

	// Then
	assert.Error(t, err)
	assert.Nil(t, val)

	cleanup()
}

func TestResolverErrorPanic(t *testing.T) {
	// Given
	setup()
	container.Bind[PrimaryIDGiver](func() *TestStruct1 {
		panic("resolver did a bad!")
	})

	// When
	val, err := container.Resolve[PrimaryIDGiver]()

	// Then
	assert.Error(t, err)
	assert.Nil(t, val)

	cleanup()
}

func TestResolverErrorNotAFunction(t *testing.T) {
	// Given
	setup()

	// When
	err := container.Bind[PrimaryIDGiver](5)

	// Then
	assert.Error(t, err)

	cleanup()
}

func TestResolverErrorInterfaceBad(t *testing.T) {
	// Given
	setup()

	// When
	err := container.Bind[int](func() int {
		return 5
	})

	// Then
	assert.Error(t, err)

	cleanup()
}

func TestResolverErrorReturnDoesntImplementInterface(t *testing.T) {
	// Given
	setup()

	// When
	err := container.Bind[PrimaryIDGiver](func() int {
		return 5
	})

	// Then
	assert.Error(t, err)

	cleanup()
}

func TestResolverErrorReturnIsntAssignableToInterface(t *testing.T) {
	// Given
	setup()

	// When
	err := container.Bind[*int](func() float32 {
		return float32(5)
	})

	// Then
	assert.Error(t, err)

	cleanup()
}

func TestResolverErrorNoReturn(t *testing.T) {
	// Given
	setup()

	// When
	err := container.Bind[*int](func() {
	})

	// Then
	assert.Error(t, err)

	cleanup()
}

func TestResolverErrorBadArg(t *testing.T) {
	// Given
	setup()

	// When
	err := container.Bind[PrimaryIDGiver](func(a int) *TestStruct1 {
		return NewTestStruct1()
	})

	// Then
	assert.Error(t, err)

	cleanup()
}

func TestMultipleInterfacesToOneConcrete(t *testing.T) {
	// Given
	setup()

	container.Bind[PrimaryIDGiver](NewTestStruct1)
	container.Bind[SecondaryIDGiver](NewTestStruct1)

	// When
	str1Prim, str1PrimErr := container.Resolve[PrimaryIDGiver]()
	str1Sec, str1SecErr := container.Resolve[SecondaryIDGiver]()

	// Then
	assert.NotNil(t, str1Prim)
	assert.NotNil(t, str1Sec)
	assert.NoError(t, str1PrimErr)
	assert.NoError(t, str1SecErr)
	assert.Equal(t, 1, Str1InstanceNumber)
	assert.Equal(t, 0, Str2InstanceNumber)

	primID := str1Prim.GivePrimaryID()
	assert.Equal(t, TestStruct1Name, primID.Name)
	assert.Equal(t, Str1InstanceNumber, primID.Number)

	secID := str1Sec.GiveSecondaryID()
	assert.Equal(t, TestStruct1Name, secID.Name)
	assert.Equal(t, Str1InstanceNumber, secID.Number)

	cleanup()
}

func TestOneInterfaceToMultipleConcretes(t *testing.T) {
	// Given
	setup()

	container.Bind[PrimaryIDGiver](NewTestStruct1)
	container.Bind[PrimaryIDGiver](NewTestStruct2)

	// When
	primSlice, err := container.ResolveAll[PrimaryIDGiver]()

	// Then
	assert.NotNil(t, primSlice)
	assert.NoError(t, err)
	assert.Len(t, primSlice, 2)
	assert.Equal(t, 1, Str1InstanceNumber)
	assert.Equal(t, 1, Str2InstanceNumber)

	str1Prim := primSlice[0]
	primID := str1Prim.GivePrimaryID()
	assert.Equal(t, TestStruct1Name, primID.Name)
	assert.Equal(t, Str1InstanceNumber, primID.Number)

	str2Prim := primSlice[1]
	secID := str2Prim.GivePrimaryID()
	assert.Equal(t, TestStruct2Name, secID.Name)
	assert.Equal(t, Str2InstanceNumber, secID.Number)

	cleanup()
}

func TestMultipleInterfaceToMultipleConcretes(t *testing.T) {
	// Given
	setup()

	container.Bind[PrimaryIDGiver](NewTestStruct1)
	container.Bind[SecondaryIDGiver](NewTestStruct1)
	container.Bind[PrimaryIDGiver](NewTestStruct2)
	container.Bind[SecondaryIDGiver](NewTestStruct2)

	// When
	PrimSlice, PrimSliceErr := container.ResolveAll[PrimaryIDGiver]()
	SecSlice, SecSliceErr := container.ResolveAll[SecondaryIDGiver]()

	// Then
	assert.NotNil(t, PrimSlice)
	assert.NotNil(t, SecSlice)
	assert.NoError(t, PrimSliceErr)
	assert.NoError(t, SecSliceErr)
	assert.Len(t, PrimSlice, 2)
	assert.Equal(t, 1, Str1InstanceNumber)
	assert.Equal(t, 1, Str2InstanceNumber)

	str1Prim := PrimSlice[0]
	str1PrimID := str1Prim.GivePrimaryID()
	assert.Equal(t, TestStruct1Name, str1PrimID.Name)
	assert.Equal(t, Str1InstanceNumber, str1PrimID.Number)

	str2Prim := PrimSlice[1]
	str2PrimID := str2Prim.GivePrimaryID()
	assert.Equal(t, TestStruct2Name, str2PrimID.Name)
	assert.Equal(t, Str2InstanceNumber, str2PrimID.Number)

	str1Sec := SecSlice[0]
	str1SecID := str1Sec.GiveSecondaryID()
	assert.Equal(t, TestStruct1Name, str1SecID.Name)
	assert.Equal(t, Str1InstanceNumber, str1SecID.Number)

	str2Sec := SecSlice[1]
	str2SecID := str2Sec.GiveSecondaryID()
	assert.Equal(t, TestStruct2Name, str2SecID.Name)
	assert.Equal(t, Str2InstanceNumber, str2SecID.Number)

	cleanup()
}

func TestResolverWithArgs(t *testing.T) {
	// Given
	setup()

	container.Bind[PrimaryIDGiver](NewTestStruct1)
	container.Bind[SecondaryIDGiver](NewTestStruct1)
	container.Bind[PrimaryIDGiver](NewTestStruct2)
	container.Bind[SecondaryIDGiver](NewTestStruct2)
	container.Bind[IDAggregator](NewTestIDAggregatorStruct)

	// When
	agg, err := container.Resolve[IDAggregator]()

	// Then
	assert.NotNil(t, agg)
	assert.NoError(t, err)
	primIDs := agg.GivePrimaryIDs()
	assert.Len(t, primIDs, 2)
	assert.Equal(t, TestStruct1Name, primIDs[0].Name)
	assert.Equal(t, Str1InstanceNumber, primIDs[0].Number)
	assert.Equal(t, TestStruct2Name, primIDs[1].Name)
	assert.Equal(t, Str2InstanceNumber, primIDs[1].Number)

	secID := agg.GiveSecondaryID()
	assert.Equal(t, TestStruct2Name, secID.Name)
	assert.Equal(t, Str2InstanceNumber, secID.Number)

	cleanup()
}

func TestResolverWithArgsMissingDependency(t *testing.T) {
	// Given
	setup()

	container.Bind[PrimaryIDGiver](NewTestStruct1)
	container.Bind[PrimaryIDGiver](NewTestStruct2)
	container.Bind[IDAggregator](NewTestIDAggregatorStruct)

	// When
	val, err := container.Resolve[IDAggregator]()

	// Then
	assert.Nil(t, val)
	assert.Error(t, err)

	cleanup()
}

func TestResolverWithArgsMissingDependencies(t *testing.T) {
	// Given
	setup()

	container.Bind[SecondaryIDGiver](NewTestStruct1)
	container.Bind[SecondaryIDGiver](NewTestStruct2)
	container.Bind[IDAggregator](NewTestIDAggregatorStruct)

	// When
	val, err := container.Resolve[IDAggregator]()

	// Then
	assert.Nil(t, val)
	assert.Error(t, err)

	cleanup()
}

// Test types
type ID struct {
	Name   string
	Number int
}

// Test interfaces
type PrimaryIDGiver interface {
	GivePrimaryID() ID
}

type SecondaryIDGiver interface {
	GiveSecondaryID() ID
}

// TestStruct1 struct - Implements PrimaryIDGiver and SecondaryIDGiver
var Str1InstanceNumber = 0

const TestStruct1Name = "TestStruct1"

type TestStruct1 struct {
	InstanceId int
}

var _ PrimaryIDGiver = &TestStruct1{}
var _ SecondaryIDGiver = &TestStruct1{}

func NewTestStruct1() *TestStruct1 {
	Str1InstanceNumber += 1
	return &TestStruct1{
		InstanceId: Str1InstanceNumber,
	}
}

func (c *TestStruct1) GivePrimaryID() ID {
	return ID{Name: TestStruct1Name, Number: c.InstanceId}
}

func (c *TestStruct1) GiveSecondaryID() ID {
	return ID{Name: TestStruct1Name, Number: c.InstanceId}
}

// TestStruct2 struct - Implements PrimaryIDGiver and SecondaryIDGiver
var Str2InstanceNumber = 0

const TestStruct2Name = "TestStruct2"

type TestStruct2 struct {
	InstanceId int
}

var _ PrimaryIDGiver = &TestStruct2{}
var _ SecondaryIDGiver = &TestStruct2{}

func NewTestStruct2() *TestStruct2 {
	Str2InstanceNumber += 1
	return &TestStruct2{
		InstanceId: Str2InstanceNumber,
	}
}

func (c *TestStruct2) GivePrimaryID() ID {
	return ID{Name: TestStruct2Name, Number: c.InstanceId}
}

func (c *TestStruct2) GiveSecondaryID() ID {
	return ID{Name: TestStruct2Name, Number: c.InstanceId}
}

// TestCompositeStruct struct - Has dependencies that must be fulfilled
type IDAggregator interface {
	GivePrimaryIDs() []ID
	GiveSecondaryID() ID
}

var TestCompositeStructID = 0

type TestIDAggregatorStruct struct {
	primeIDGivers    []PrimaryIDGiver
	secondaryIDGiver SecondaryIDGiver
}

var _ IDAggregator = &TestIDAggregatorStruct{}

func NewTestIDAggregatorStruct(primIDGivers []PrimaryIDGiver, secondaryIDGiver SecondaryIDGiver) *TestIDAggregatorStruct {
	return &TestIDAggregatorStruct{
		primeIDGivers:    primIDGivers,
		secondaryIDGiver: secondaryIDGiver,
	}
}

func (c *TestIDAggregatorStruct) GivePrimaryIDs() []ID {
	retVal := make([]ID, len(c.primeIDGivers))
	for idx, val := range c.primeIDGivers {
		retVal[idx] = val.GivePrimaryID()
	}
	return retVal
}

func (c *TestIDAggregatorStruct) GiveSecondaryID() ID {
	return c.secondaryIDGiver.GiveSecondaryID()
}
