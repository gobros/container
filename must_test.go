package container_test

import (
	"testing"

	"github.com/gobros/container"
	"github.com/stretchr/testify/assert"
)

func TestMustSimpleBind(t *testing.T) {
	// Given
	setup()

	container.MustBind[PrimaryIDGiver](NewTestStruct1)

	// When
	str1Prim := container.MustResolve[PrimaryIDGiver]()

	// Then
	assert.NotNil(t, str1Prim)
	assert.Equal(t, 1, Str1InstanceNumber)
	assert.Equal(t, 0, Str2InstanceNumber)

	primID := str1Prim.GivePrimaryID()
	assert.Equal(t, TestStruct1Name, primID.Name)
	assert.Equal(t, Str1InstanceNumber, primID.Number)

	cleanup()
}

func TestMustSimpleBindOverride(t *testing.T) {
	// Given
	setup()

	container.MustBind[PrimaryIDGiver](NewTestStruct1)
	container.MustBind[PrimaryIDGiver](NewTestStruct1)

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

func TestMustSimpleBindPointer(t *testing.T) {
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
	container.MustBind[*ptrStruct](func() *ptrStruct {
		return &testVal
	})

	// When
	resolvedVal := container.MustResolve[*ptrStruct]()

	// Then
	assert.NotNil(t, resolvedVal)
	assert.Equal(t, "wirecat", resolvedVal.Name)
	assert.Equal(t, 1337, resolvedVal.FaveNumber)

	cleanup()
}

func TestMustNothingBoundResolve(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() { container.MustResolve[PrimaryIDGiver]() })

	cleanup()
}

func TestMustNothingBoundResolveAll(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() { container.MustResolveAll[PrimaryIDGiver]() })

	cleanup()
}

func TestMustResolverErrorNotAFunction(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() { container.MustBind[PrimaryIDGiver](5) })

	cleanup()
}

func TestMustResolverErrorInterfaceBad(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() {
		container.MustBind[int](func() int {
			return 5
		})
	})

	cleanup()
}

func TestMustResolverErrorReturnDoesntImplementInterface(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() {
		container.MustBind[PrimaryIDGiver](func() int {
			return 5
		})
	})

	cleanup()
}

func TestMustResolverErrorReturnIsntAssignableToInterface(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() {
		container.MustBind[*int](func() float32 {
			return float32(5)
		})
	})

	cleanup()
}

func TestMustResolverErrorNoReturn(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() {
		container.MustBind[*int](func() {
		})
	})

	cleanup()
}

func TestMustResolverErrorBadArg(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() {
		container.MustBind[PrimaryIDGiver](func(a int) *TestStruct1 {
			return NewTestStruct1()
		})
	})

	cleanup()
}

func TestMustMultipleInterfacesToOneConcrete(t *testing.T) {
	// Given
	setup()

	container.MustBind[PrimaryIDGiver](NewTestStruct1)
	container.MustBind[SecondaryIDGiver](NewTestStruct1)

	// When
	str1Prim := container.MustResolve[PrimaryIDGiver]()
	str1Sec := container.MustResolve[SecondaryIDGiver]()

	// Then
	assert.NotNil(t, str1Prim)
	assert.NotNil(t, str1Sec)
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

func TestMustOneInterfaceToMultipleConcretes(t *testing.T) {
	// Given
	setup()

	container.MustBind[PrimaryIDGiver](NewTestStruct1)
	container.MustBind[PrimaryIDGiver](NewTestStruct2)

	// When
	primSlice := container.MustResolveAll[PrimaryIDGiver]()

	// Then
	assert.NotNil(t, primSlice)
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

func TestMustMultipleInterfaceToMultipleConcretes(t *testing.T) {
	// Given
	setup()

	container.MustBind[PrimaryIDGiver](NewTestStruct1)
	container.MustBind[SecondaryIDGiver](NewTestStruct1)
	container.MustBind[PrimaryIDGiver](NewTestStruct2)
	container.MustBind[SecondaryIDGiver](NewTestStruct2)

	// When
	PrimSlice := container.MustResolveAll[PrimaryIDGiver]()
	SecSlice := container.MustResolveAll[SecondaryIDGiver]()

	// Then
	assert.NotNil(t, PrimSlice)
	assert.NotNil(t, SecSlice)
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

func TestMustResolverWithArgs(t *testing.T) {
	// Given
	setup()

	container.MustBind[PrimaryIDGiver](NewTestStruct1)
	container.MustBind[SecondaryIDGiver](NewTestStruct1)
	container.MustBind[PrimaryIDGiver](NewTestStruct2)
	container.MustBind[SecondaryIDGiver](NewTestStruct2)
	container.MustBind[IDAggregator](NewTestIDAggregatorStruct)

	// When
	agg := container.MustResolve[IDAggregator]()

	// Then
	assert.NotNil(t, agg)
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

func TestMustResolverWithArgsMissingDependency(t *testing.T) {
	// Given
	setup()

	container.MustBind[PrimaryIDGiver](NewTestStruct1)
	container.MustBind[PrimaryIDGiver](NewTestStruct2)
	container.MustBind[IDAggregator](NewTestIDAggregatorStruct)

	// When & Then
	assert.Panics(t, func() { container.MustResolve[IDAggregator]() })

	cleanup()
}

func TestMustResolverWithArgsMissingDependencies(t *testing.T) {
	// Given
	setup()

	container.MustBind[SecondaryIDGiver](NewTestStruct1)
	container.MustBind[SecondaryIDGiver](NewTestStruct2)
	container.MustBind[IDAggregator](NewTestIDAggregatorStruct)

	// When & Then
	val := container.MustResolve[IDAggregator]()
	assert.NotNil(t, val)
	assert.Empty(t, val.GivePrimaryIDs())

	cleanup()
}

func TestMustBindInstanceHappy(t *testing.T) {
	// Given
	setup()

	container.MustBindInstance[PrimaryIDGiver](container.Global, NewTestStruct1)

	// When
	str1Prim := container.MustResolve[PrimaryIDGiver]()

	// Then
	assert.NotNil(t, str1Prim)
	assert.Equal(t, 1, Str1InstanceNumber)
	assert.Equal(t, 0, Str2InstanceNumber)

	primID := str1Prim.GivePrimaryID()
	assert.Equal(t, TestStruct1Name, primID.Name)
	assert.Equal(t, Str1InstanceNumber, primID.Number)

	cleanup()
}

func TestMustBindInstancePanic(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() { container.MustBindInstance[PrimaryIDGiver](container.Global, 5) })

	cleanup()
}

func TestMustResolveAllInstanceHappy(t *testing.T) {
	// Given
	setup()

	container.MustBindInstance[PrimaryIDGiver](container.Global, NewTestStruct1)

	// When
	val := container.MustResolveAllInstance[PrimaryIDGiver](container.Global)

	// Then
	assert.NotNil(t, val)
	assert.NotEmpty(t, val)

	cleanup()
}

func TestMustResolveAllInstancePanic(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() { container.MustResolveAllInstance[PrimaryIDGiver](container.Global) })

	cleanup()
}

func TestMustResolveInstanceHappy(t *testing.T) {
	// Given
	setup()

	container.MustBindInstance[PrimaryIDGiver](container.Global, NewTestStruct1)

	// When
	val := container.MustResolveInstance[PrimaryIDGiver](container.Global)

	// Then
	assert.NotNil(t, val)

	cleanup()
}

func TestMustResolveInstancePanic(t *testing.T) {
	// Given
	setup()

	// When & Then
	assert.Panics(t, func() { container.MustResolveInstance[PrimaryIDGiver](container.Global) })

	cleanup()
}
