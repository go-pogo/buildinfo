package main

import (
	"fmt"

	"github.com/roeldev/go-buildinfo"
)

// these values are changed via ldflags when building a new release
var (
	version   = buildinfo.DummyVersion
	buildDate = buildinfo.DummyDate
	gitBranch = buildinfo.DummyBranch
	gitCommit = buildinfo.DummyCommit
)

func main() {
	fmt.Println(buildinfo.BuildInfo{
		Version: version,
		Date:    buildDate,
		Branch:  gitBranch,
		Commit:  gitCommit,
	})
}
