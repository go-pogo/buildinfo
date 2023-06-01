// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.18
// +build go1.18

package buildinfo

import (
	"runtime/debug"
	"time"
)

func (bld *BuildInfo) init() {
	bi, _ := debug.ReadBuildInfo()
	bld.goVersion = bi.GoVersion

	for _, set := range bi.Settings {
		switch set.Key {
		case "vcs.revision":
			bld.Revision = set.Value
		case "vcs.time":
			bld.Created, _ = time.Parse(time.RFC3339, set.Value)
		}
	}
}
