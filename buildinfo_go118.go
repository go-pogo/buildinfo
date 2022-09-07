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

func New(ver string) *BuildInfo {
	bi, _ := debug.ReadBuildInfo()
	bld := BuildInfo{
		goVersion: bi.GoVersion,
		Version:   ver,
	}

	for _, set := range bi.Settings {
		switch set.Key {
		case "vcs.revision":
			bld.Revision = set.Value
		case "vcs.time":
			bld.Time, _ = time.Parse(time.RFC3339, set.Value)
		}
	}
	return &bld
}
