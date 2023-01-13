// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	Query  = queryBinding{}
	Uri    = uriBinding{}
	Header = headerBinding{}
)
