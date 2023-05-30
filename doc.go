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

# Using an embedded file

	package main

	import _ "embed"

	//go:embed buildinfo.json
	var buildInfo []byte

	func main() {
		var bld buildinfo.BuildInfo
		if err := json.Unmarshal(buildInfo, &bld); err != nil {
			panic(err)
		}
	}

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
