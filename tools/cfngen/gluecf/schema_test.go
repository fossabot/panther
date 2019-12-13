package gluecf

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Field1 string
	Field2 int32

	TagSuppressed int `json:"-"` // should be skipped cuz of tag

	// should not be emitted because they are private
	privateField       int      // nolint
	setOfPrivateFields struct { // nolint
		subField1 int
		subField2 int
	}
}

func (ts *TestStruct) Foo() { // admits to TestInterface
}

type NestedStruct struct {
	A TestStruct
	B TestStruct
	C *TestStruct
}

type TestInterface interface {
	Foo()
}

func TestInferJsonColumns(t *testing.T) {
	// used to test pointers
	var s string = "S"
	var i int32 = 1
	var f float32 = 1

	obj := struct { // nolint
		BoolField bool

		StringField    string  `json:"stringField"`              // test we use json tags
		StringPtrField *string `json:"stringPtrField,omitempty"` // test we use json tags

		IntField    int
		Int8Field   int8
		Int16Field  int16
		Int32Field  int32
		Int64Field  int64
		IntPtrField *int32

		Float32Field    float32
		Float64Field    float64
		Float32PtrField *float32

		StringSlice []string

		IntSlice   []int
		Int32Slice []int32
		Int64Slice []int64

		Float32Slice []float32
		Float64Slice []float64

		StructSlice []TestStruct

		MapSlice []map[string]string

		TimeField time.Time

		MapStringToInterface map[string]interface{}
		MapStringToString    map[string]string
		MapStringToStruct    map[string]TestStruct

		StructField       TestStruct
		NestedStructField NestedStruct
	}{
		BoolField: true,

		StringField:    s,
		StringPtrField: &s,

		IntField:    1,
		Int8Field:   1,
		Int16Field:  1,
		Int32Field:  1,
		Int64Field:  1,
		IntPtrField: &i,

		Float32Field:    1,
		Float64Field:    1,
		Float32PtrField: &f,

		StringSlice: []string{"S1", "S2"},

		IntSlice:   []int{1, 2, 3},
		Int32Slice: []int32{1, 2, 3},
		Int64Slice: []int64{1, 2, 3},

		Float32Slice: []float32{1, 2, 3},
		Float64Slice: []float64{1, 2, 3},

		StructSlice: []TestStruct{},

		MapSlice: []map[string]string{
			make(map[string]string),
		},

		TimeField: time.Date(2019, 12, 1, 1, 1, 1, 1, time.UTC),

		MapStringToInterface: make(map[string]interface{}),
		MapStringToString:    make(map[string]string),
		MapStringToStruct:    make(map[string]TestStruct),

		StructField: TestStruct{},
		NestedStructField: NestedStruct{
			C: &TestStruct{}, // test with ptrs
		},
	}

	// adjust for native int expected results
	nativeIntMapping := func() string {
		switch strconv.IntSize {
		case 32:
			return "int"
		case 64:
			return "bigint"
		default:
			panic(fmt.Sprintf("Size of native int unexpected: %d", strconv.IntSize))
		}
	}

	excpectedCols := []Column{
		{Name: "BoolField", Type: "boolean"},
		{Name: "stringField", Type: "string"},
		{Name: "stringPtrField", Type: "string"},
		{Name: "IntField", Type: nativeIntMapping()},
		{Name: "Int8Field", Type: "tinyint"},
		{Name: "Int16Field", Type: "smallint"},
		{Name: "Int32Field", Type: "int"},
		{Name: "Int64Field", Type: "bigint"},
		{Name: "IntPtrField", Type: "int"},
		{Name: "Float32Field", Type: "float"},
		{Name: "Float64Field", Type: "double"},
		{Name: "Float32PtrField", Type: "float"},
		{Name: "StringSlice", Type: "array<string>"},
		{Name: "IntSlice", Type: "array<" + nativeIntMapping() + ">"},
		{Name: "Int32Slice", Type: "array<int>"},
		{Name: "Int64Slice", Type: "array<bigint>"},
		{Name: "Float32Slice", Type: "array<float>"},
		{Name: "Float64Slice", Type: "array<double>"},
		{Name: "StructSlice", Type: "array<struct<Field1:string,Field2:int>>"},
		{Name: "MapSlice", Type: "array<map<string,string>>"},
		{Name: "TimeField", Type: "timestamp"},
		{Name: "MapStringToInterface", Type: "map<string,string>"}, // special case
		{Name: "MapStringToString", Type: "map<string,string>"},
		{Name: "MapStringToStruct", Type: "map<string,struct<Field1:string,Field2:int>>"},
		{Name: "StructField", Type: "struct<Field1:string,Field2:int>"},
		{Name: "NestedStructField", Type: "struct<A:struct<Field1:string,Field2:int>,B:struct<Field1:string,Field2:int>,C:struct<Field1:string,Field2:int>>"}, // nolint
	}

	cols := InferJSONColumns(obj)

	// uncomment to see results
	/*
		for _, col := range cols {
			fmt.Printf("{Name: \"%s\", Type: \"%s\"},\n", col.Name, col.Type)
		}
	*/
	assert.Equal(t, excpectedCols, cols, "Expected columns not found")

	// Test using interface
	var testInterface TestInterface = &TestStruct{}
	cols = InferJSONColumns(testInterface)
	assert.Equal(t, []Column{{Name: "Field1", Type: "string"}, {Name: "Field2", Type: "int"}}, cols, "Interface test failed")
}
