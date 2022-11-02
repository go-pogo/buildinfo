// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package buildinfo provides basic building blocks and instructions to easily add
build and release information to your app.

# Using ldflags
Declare build info variables in your main package:

	package main

	// this value is changed via ldflags when building a new release
	var version

	func main() {
		bld := buildinfo.New(version)
	}

Build your Go project and include the following ldflags:

	go build -ldflags=" \
	  -X main.version=`$(git describe --tags)` \
	  main.go

# Using gen
Create a file that is called by go generate:

	//go:build ignore
	package main

	import "github.com/go-pogo/buildinfo/gen"

	func main() {
		gen.GenerateFile("buildinfo.go")
	}

Use it to generate a buildinfo.go file containing the latest tag of your
project's repository. The file should be renewed any time a new tag is created.
This is typically done during build.

# CLI tool

	go install github.com/go-pogo/buildinfo/cmd/buildinfo@latest

# Prometheus metric collector
When using a metrics scraper like Prometheus, it is often a good idea to make
the build information of your app available. Below example shows just how easy
it is to create and register a collector with the build information as
constant labels.

	prometheus.MustRegister(prometheus.NewGaugeFunc(
	    prometheus.GaugeOpts{
	        Namespace:   "myapp",
	        Name:        buildinfo.MetricName,
	        Help:        buildinfo.MetricHelp,
	        ConstLabels: bld.Map(),
	    },
	    func() float64 { return 1 },
	))
*/
package buildinfo
