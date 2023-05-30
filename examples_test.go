// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"bytes"
	_ "embed"
	"fmt"
)

func ExampleNew() {
	fmt.Println(New("1.2.3").String())
	// Output: 1.2.3
}

//go:embed example.json
var someEmbeddedJsonData []byte

func ExampleBuildInfo_UnmarshalJSON() {
	data := someEmbeddedJsonData
	var bld BuildInfo
	if err := bld.UnmarshalJSON(data); err != nil {
		panic(err)
	}

	fmt.Printf("version=%s, something=%s\n", bld.Version, bld.Extra["something"])
	// Output: version=1.2.3, something=else
}

func ExampleRead() {
	buf := bytes.NewBufferString(`{"version":"1.2.3","something":"else"}`)

	bld, err := Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("version=%s, something=%s\n", bld.Version, bld.Extra["something"])
	// Output: version=1.2.3, something=else
}

func ExampleBuildInfo_WithExtra() {
	var bld BuildInfo
	bld.WithExtra("extra", "value")

	fmt.Printf("version=%s, extra=%s\n", bld.Version, bld.Extra["extra"])
	// Output: version=, extra=value
}
