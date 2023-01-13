// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type appkey struct {
	Appkey string `json:"appkey" form:"appkey"`
}

type QueryTest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	appkey
}

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required,max=32"`
}

type FooBarStruct struct {
	FooStruct
	Bar string `msgpack:"bar" json:"bar" form:"bar" xml:"bar" binding:"required"`
}

type FooStructUseNumber struct {
	Foo any `json:"foo" binding:"required"`
}

type FooStructDisallowUnknownFields struct {
	Foo any `json:"foo" binding:"required"`
}

type FooStructForMapType struct {
	MapFoo map[string]any `form:"map_foo"`
}

type FooStructForBoolType struct {
	BoolFoo bool `form:"bool_foo"`
}

func TestBindingQuery(t *testing.T) {
	testQueryBinding(t, "POST",
		"/?foo=bar&bar=foo", "/",
		"foo=unused", "bar2=foo")
}

func TestBindingQuery2(t *testing.T) {
	testQueryBinding(t, "GET",
		"/?foo=bar&bar=foo", "/?bar2=foo",
		"foo=unused", "")
}

func TestBindingQueryFail(t *testing.T) {
	testQueryBindingFail(t, "POST",
		"/?map_foo=", "/",
		"map_foo=unused", "bar2=foo")
}

func TestBindingQueryFail2(t *testing.T) {
	testQueryBindingFail(t, "GET",
		"/?map_foo=", "/?bar2=foo",
		"map_foo=unused", "")
}

func TestBindingQueryBoolFail(t *testing.T) {
	testQueryBindingBoolFail(t, "GET",
		"/?bool_foo=fasl", "/?bar2=foo",
		"bool_foo=unused", "")
}

func TestBindingQueryStringMap(t *testing.T) {
	obj := make(map[string]string)
	req := requestWithBody("GET", "/?foo=bar&hello=world", "")
	err := BindQuery(req, &obj)
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	assert.Equal(t, "bar", obj["foo"])
	assert.Equal(t, "world", obj["hello"])

	obj = make(map[string]string)
	req = requestWithBody("GET", "/?foo=bar&foo=2&hello=world", "") // should pick last
	err = BindQuery(req, &obj)
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Len(t, obj, 2)
	assert.Equal(t, "2", obj["foo"])
	assert.Equal(t, "world", obj["hello"])
}

func TestUriBinding(t *testing.T) {

	type Tag struct {
		Name string `uri:"name"`
	}
	var tag Tag
	m := make(map[string][]string)
	m["name"] = []string{"thinkerou"}
	assert.NoError(t, BindURI(m, &tag))
	assert.Equal(t, "thinkerou", tag.Name)

	type NotSupportStruct struct {
		Name map[string]any `uri:"name"`
	}
	var not NotSupportStruct
	assert.Error(t, BindURI(m, &not))
	assert.Equal(t, map[string]any(nil), not.Name)
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
	assert.NoError(t, BindURI(m, &tag))
	assert.Equal(t, tag.Name, expectedName)
	assert.Equal(t, tag.S.Age, expectedAge)
}

const MIMEPOSTForm = "application/x-www-form-urlencoded"

func testQueryBinding(t *testing.T, method, path, badPath, body, badBody string) {
	obj := FooBarStruct{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := BindQuery(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)
	assert.Equal(t, "foo", obj.Bar)
}

func testQueryBindingFail(t *testing.T, method, path, badPath, body, badBody string) {
	obj := FooStructForMapType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := BindQuery(req, &obj)
	assert.Error(t, err)
}

func testQueryBindingBoolFail(t *testing.T, method, path, badPath, body, badBody string) {
	obj := FooStructForBoolType{}
	req := requestWithBody(method, path, body)
	if method == "POST" {
		req.Header.Add("Content-Type", MIMEPOSTForm)
	}
	err := BindQuery(req, &obj)
	assert.Error(t, err)
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
