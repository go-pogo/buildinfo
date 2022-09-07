// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.18
// +build !go1.18

package buildinfo

import (
	"runtime"
)

func New(ver string) *BuildInfo {
	return &BuildInfo{
		goVersion: runtime.Version(),
		Version:   ver,
	}
}
