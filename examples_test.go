// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"fmt"
)

// this value can be changed via ldflags when building a new release
var version = "1.2.3"

func ExampleNew() {
	bld := New(version)
	fmt.Println(bld.String())
	// Output: 1.2.3
}
