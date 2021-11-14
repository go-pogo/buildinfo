// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/go-pogo/buildinfo"
)

// these values should be changed via ldflags when building a new release
var (
	version  = buildinfo.DummyVersion
	revision = buildinfo.DummyRevision
	date     = buildinfo.DummyDate
)

func main() {
	fmt.Println(buildinfo.BuildInfo{
		Version:  version,
		Revision: revision,
		Date:     date,
	})
}
