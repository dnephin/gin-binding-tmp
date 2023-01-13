// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMappingBaseTypes(t *testing.T) {
	intPtr := func(i int) *int {
		return &i
	}
	for _, tt := range []struct {
		name   string
		value  any
		form   string
		expect any
	}{
		{"base type", struct{ F int }{}, "9", int(9)},
		{"base type", struct{ F int8 }{}, "9", int8(9)},
		{"base type", struct{ F int16 }{}, "9", int16(9)},
		{"base type", struct{ F int32 }{}, "9", int32(9)},
		{"base type", struct{ F int64 }{}, "9", int64(9)},
		{"base type", struct{ F uint }{}, "9", uint(9)},
		{"base type", struct{ F uint8 }{}, "9", uint8(9)},
		{"base type", struct{ F uint16 }{}, "9", uint16(9)},
		{"base type", struct{ F uint32 }{}, "9", uint32(9)},
		{"base type", struct{ F uint64 }{}, "9", uint64(9)},
		{"base type", struct{ F bool }{}, "True", true},
		{"base type", struct{ F float32 }{}, "9.1", float32(9.1)},
		{"base type", struct{ F float64 }{}, "9.1", float64(9.1)},
		{"base type", struct{ F string }{}, "test", string("test")},
		{"base type", struct{ F *int }{}, "9", intPtr(9)},

		// zero values
		{"zero value", struct{ F int }{}, "", int(0)},
		{"zero value", struct{ F uint }{}, "", uint(0)},
		{"zero value", struct{ F bool }{}, "", false},
		{"zero value", struct{ F float32 }{}, "", float32(0)},
	} {
		tp := reflect.TypeOf(tt.value)
		testName := tt.name + ":" + tp.Field(0).Type.String()

		val := reflect.New(reflect.TypeOf(tt.value))
		val.Elem().Set(reflect.ValueOf(tt.value))

		field := val.Elem().Type().Field(0)

		_, err := mapping(val, emptyField, formSource{field.Name: {tt.form}}, "form")
		assert.NoError(t, err, testName)

		actual := val.Elem().Field(0).Interface()
		assert.Equal(t, tt.expect, actual, testName)
	}
}

func TestMappingDefault(t *testing.T) {
	var s struct {
		Int   int    `form:",default=9"`
		Slice []int  `form:",default=9"`
		Array [1]int `form:",default=9"`
	}
	err := decode(&s, formSource{}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 9, s.Int)
	assert.Equal(t, []int{9}, s.Slice)
	assert.Equal(t, [1]int{9}, s.Array)
}

func TestMappingSkipField(t *testing.T) {
	var s struct {
		A int
	}
	err := decode(&s, formSource{}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 0, s.A)
}

func TestMappingIgnoreField(t *testing.T) {
	var s struct {
		A int `form:"A"`
		B int `form:"-"`
	}
	err := decode(&s, formSource{"A": {"9"}, "B": {"9"}}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 9, s.A)
	assert.Equal(t, 0, s.B)
}

func TestMappingUnexportedField(t *testing.T) {
	var s struct {
		A int `form:"a"`
		b int `form:"b"`
	}
	err := decode(&s, formSource{"a": {"9"}, "b": {"9"}}, "form")
	assert.NoError(t, err)

	assert.Equal(t, 9, s.A)
	assert.Equal(t, 0, s.b)
}

func TestMappingPrivateField(t *testing.T) {
	var s struct {
		f int `form:"field"`
	}
	err := decode(&s, formSource{"field": {"6"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, 0, s.f)
}

func TestMappingUnknownFieldType(t *testing.T) {
	var s struct {
		U uintptr
	}

	err := decode(&s, formSource{"U": {"unknown"}}, "form")
	assert.Error(t, err)
	assert.Equal(t, errUnknownType, err)
}

func TestBindURI(t *testing.T) {
	var s struct {
		F int `uri:"field"`
	}
	err := BindURI(map[string][]string{"field": {"6"}}, &s)
	assert.NoError(t, err)
	assert.Equal(t, 6, s.F)
}

func TestMappingForm(t *testing.T) {
	var s struct {
		F int `form:"field"`
	}
	err := decode(&s, map[string][]string{"field": {"6"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, 6, s.F)
}

func TestMappingSlice(t *testing.T) {
	var s struct {
		Slice []int `form:"slice,default=9"`
	}

	// default value
	err := decode(&s, formSource{}, "form")
	assert.NoError(t, err)
	assert.Equal(t, []int{9}, s.Slice)

	// ok
	err = decode(&s, formSource{"slice": {"3", "4"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4}, s.Slice)

	// error
	err = decode(&s, formSource{"slice": {"wrong"}}, "form")
	assert.Error(t, err)
}

func TestMappingArray(t *testing.T) {
	var s struct {
		Array [2]int `form:"array,default=9"`
	}

	// wrong default
	err := decode(&s, formSource{}, "form")
	assert.Error(t, err)

	// ok
	err = decode(&s, formSource{"array": {"3", "4"}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, [2]int{3, 4}, s.Array)

	// error - not enough vals
	err = decode(&s, formSource{"array": {"3"}}, "form")
	assert.Error(t, err)

	// error - wrong value
	err = decode(&s, formSource{"array": {"wrong"}}, "form")
	assert.Error(t, err)
}

func TestMappingStructField(t *testing.T) {
	var s struct {
		J struct {
			I int
		}
	}

	err := decode(&s, formSource{"J": {`{"I": 9}`}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, 9, s.J.I)
}

func TestMappingMapField(t *testing.T) {
	var s struct {
		M map[string]int
	}

	err := decode(&s, formSource{"M": {`{"one": 1}`}}, "form")
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"one": 1}, s.M)
}

func TestMappingIgnoredCircularRef(t *testing.T) {
	type S struct {
		S *S `form:"-"`
	}
	var s S

	err := decode(&s, formSource{}, "form")
	assert.NoError(t, err)
}
