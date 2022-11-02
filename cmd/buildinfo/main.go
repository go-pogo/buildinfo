// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-pogo/errors"
)

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	latest := flags.Bool("l", false, "Use latest tag")
	help := flags.Bool("h", false, "Display this help message")

	if err := flags.Parse(os.Args[1:]); *help || errors.Is(err, flag.ErrHelp) {
		_, _ = fmt.Fprintf(flags.Output(), "Arguments:\n")
		flags.PrintDefaults()
		_, _ = fmt.Fprintf(flags.Output(), `
Output commands:
	original            original version (default)
	version             version without metadata
	full                version including metadata
	major.minor.patch   alias of verion
	major.minor         major and minor version parts
	major               major version part only
	minor               minor version part only
	patch               patch part only
	+major              version with increased major part
	+minor              version with increased minor part
	+patch              version with increased patch part
	revision            commit revision
	rev                 alias of rev
	time                time of commit
`)
		os.Exit(0)
	}

	if err := run(context.Background(), flags.Args(), *latest); err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "fatal error: %s", err.Error())
		os.Exit(errors.GetExitCodeOr(err, 1))
	}
}
