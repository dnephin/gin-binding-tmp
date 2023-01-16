// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"net/http"
	"strconv"
	"testing"

	"gotest.tools/v3/assert"
)

type FooStruct struct {
	Foo string `json:"foo" form:"foo"`
}

type FooBarStruct struct {
	FooStruct
	Bar string `json:"bar" form:"bar"`
}

type FooStructForMapType struct {
	MapFoo map[string]any `form:"map_foo"`
}

type FooStructForBoolType struct {
	BoolFoo bool `form:"bool_foo"`
}

func TestBindingQuery(t *testing.T) {
	testQueryBinding(t, "POST", "/?foo=bar&bar=foo", "foo=unused")
}

func TestBindingQuery2(t *testing.T) {
	testQueryBinding(t, "GET", "/?foo=bar&bar=foo", "foo=unused")
}

func TestBindingQueryFail(t *testing.T) {
	testQueryBindingFail(t, "POST", "/?map_foo=", "map_foo=unused")
}

func TestBindingQueryFail2(t *testing.T) {
	testQueryBindingFail(t, "GET", "/?map_foo=", "map_foo=unused")
}

func TestBindingQueryBoolFail(t *testing.T) {
	testQueryBindingBoolFail(t, "GET", "/?bool_foo=fasl", "bool_foo=unused")
}

func TestUriBinding(t *testing.T) {

	type Tag struct {
		Name string `uri:"name"`
	}
	var tag Tag
	m := make(map[string][]string)
	m["name"] = []string{"thinkerou"}
	assert.NilError(t, BindURI(m, &tag))
	assert.Equal(t, "thinkerou", tag.Name)

	type NotSupportStruct struct {
		Name map[string]any `uri:"name"`
	}
	var not NotSupportStruct
	assert.ErrorContains(t, BindURI(m, &not), "invalid character")
	assert.DeepEqual(t, map[string]any(nil), not.Name)
}

func TestUriInnerBinding(t *testing.T) {
	type Tag struct {
		Name string `uri:"name"`
		S    struct {
			Age int `uri:"age"`
		}
	}

	expectedName := "mike"
	expectedAge := 25

	m := map[string][]string{
		"name": {expectedName},
		"age":  {strconv.Itoa(expectedAge)},
	}

	var tag Tag
	assert.NilError(t, BindURI(m, &tag))
	assert.Equal(t, tag.Name, expectedName)
	assert.Equal(t, tag.S.Age, expectedAge)
}

const MIMEPOSTForm = "application/x-www-form-urlencoded"

func testQueryBinding(t *testing.T, method, path, body string) {
	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := BindQuery(req, &obj)
	assert.NilError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func testQueryBindingFail(t *testing.T, method, path, body string) {
	obj := FooStructForMapType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := BindQuery(req, &obj)
	assert.Error(t, err, "unexpected end of JSON input")
}

func testQueryBindingBoolFail(t *testing.T, method, path, body string) {
	obj := FooStructForBoolType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := BindQuery(req, &obj)
	assert.ErrorContains(t, err, "invalid syntax")
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
