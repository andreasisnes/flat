package goflat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyStruct(t *testing.T) {
	type nestedNestedStruct struct {
		Field bool
	}

	type nestedStruct struct {
		Field              string
		NestedNestedStruct nestedNestedStruct
	}

	type testStruct struct {
		Field        string
		NestedStruct nestedStruct
		AnotherField int
	}

	value := &testStruct{
		Field: "FieldOuter",
		NestedStruct: nestedStruct{
			Field: "FieldInner",
			NestedNestedStruct: nestedNestedStruct{
				Field: true,
			},
		},
		AnotherField: 1,
	}

	v := &testStruct{}
	Copy(value, &v)

	assert.Equal(t, v.Field, value.Field)
	assert.Equal(t, v.AnotherField, value.AnotherField)
	assert.Equal(t, v.NestedStruct.Field, value.NestedStruct.Field)
	assert.Equal(t, v.NestedStruct.NestedNestedStruct.Field, value.NestedStruct.NestedNestedStruct.Field)
}

func TestCopyPointer(t *testing.T) {
	ptrstr1 := "ptr1"
	ptrstr2 := &ptrstr1
	type testStruct struct {
		Ptr1 *string
		Ptr2 **string
	}

	value := &testStruct{
		Ptr1: &ptrstr1,
		Ptr2: &ptrstr2,
	}

	v := &testStruct{}
	Copy(value, &v)

	assert.Equal(t, ptrstr1, value.Ptr1)
	assert.Equal(t, *ptrstr2, value.Ptr1)
}
