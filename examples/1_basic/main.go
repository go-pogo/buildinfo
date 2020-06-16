package main

import (
	"fmt"

	"github.com/roeldev/go-buildinfo"
)

// these values are changed via ldflags when building a new release
var (
	version   = ""
	buildDate = ""
	gitBranch = ""
	gitCommit = ""
)

func main() {
	fmt.Println(buildinfo.BuildInfo{
		Version: version,
		Date:    buildDate,
		Branch:  gitBranch,
		Commit:  gitCommit,
	})
}
