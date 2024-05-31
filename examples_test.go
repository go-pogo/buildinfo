// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"fmt"
)

func ExampleNew() {
	bld, _ := New("1.2.3")
	fmt.Println(bld.String())
	// Output: 1.2.3
}
