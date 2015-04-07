// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/json"

	"net/http"
)

type jsonBinding struct{}

func (_ jsonBinding) Name() string {
	return "json"
}

func (_ jsonBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	if err := _validator.ValidateStruct(obj); err != nil {
		return error(err)
	}
	return nil
}
