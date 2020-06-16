package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	buildInfo := buildinfo.BuildInfo{
		Version: version,
		Date:    buildDate,
		Branch:  gitBranch,
		Commit:  gitCommit,
	}

	var displayBuildInfo bool

	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flags.SetOutput(os.Stdout)
	flags.BoolVar(&displayBuildInfo, buildinfo.ShortFlag, false, "Display build version information")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if displayBuildInfo {
		fmt.Println(buildInfo)
		return
	}

	flags.Usage()
}
