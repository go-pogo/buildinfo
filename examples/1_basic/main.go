package main

import (
	"fmt"

	"github.com/roeldev/go-buildinfo"
)

// these values are changed via ldflags when building a new release
var (
	v = buildinfo.DummyVersion
	r = buildinfo.DummyRevision
	b = buildinfo.DummyBranch
	d = buildinfo.DummyDate
)

func main() {
	fmt.Println(buildinfo.BuildInfo{
		Version:  v,
		Revision: r,
		Branch:   b,
		Date:     d,
	})
}
