// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package buildinfo provides basic building blocks and instructions to easily add
build and release information to your app.

	var (
		version = buildinfo.DummyVersion
		revision = buildinfo.DummyRevision
		date = buildinfo.DummyDate
	)

	func main() {
		bld := buildinfo.BuildInfo{
			Version:  version,
			Revision: revision,
			Date:     date,
		}
	}

Build your Go project and include the following ldflags:

	go build -ldflags=" \
	  -X main.version=`$(git describe --tags)` \
	  -X main.revision=`$(git rev-parse --short HEAD)` \
	  -X main.date=`$(date +%FT%T%z`)" \
	  main.go

Prometheus metric collector

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

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

//goland:noinspection GoUnusedConst
const (
	// ShortFlag is the default flag to print the current build information of the app.
	ShortFlag = "v"
	// LongFlag is an alternative long version that may be used together with ShortFlag.
	LongFlag = "version"

	// MetricName is a default name for the metric (without namespace).
	MetricName = "build_info"
	// MetricHelp is a default help text that describes the metric.
	MetricHelp = "Metric with build information labels and a constant value of '1'."
)

// BuildInfo contains the relevant information of the current release's build
// version, revision and build date.
type BuildInfo struct {
	// Version is the current version of the release.
	Version string
	// Revision is the (short) commit hash the release is build from.
	Revision string
	// Date of when the release was build.
	Date string
}

// GoVersion returns the version of the used Go runtime. See runtime.Version
// for additional details.
func (bld BuildInfo) GoVersion() string { return runtime.Version() }

// Map returns the build information as a map. Field names are lowercase.
// Empty fields within BuildInfo are omitted.
func (bld BuildInfo) Map() map[string]string {
	m := make(map[string]string, 5)
	m["version"] = bld.Version

	if bld.Revision != "" {
		m["revision"] = bld.Revision
	}
	if bld.Date != "" {
		m["date"] = bld.Date
	}

	m["goversion"] = bld.GoVersion()
	return m
}

// MarshalJSON returns valid JSON output.
// Empty fields within BuildInfo are omitted.
func (bld BuildInfo) MarshalJSON() ([]byte, error) {
	var buf strings.Builder
	buf.WriteString(`{"version":"`)
	buf.WriteString(bld.Version)

	if bld.Revision != "" {
		buf.WriteString(`","revision":"`)
		buf.WriteString(bld.Revision)
	}
	if bld.Date != "" {
		buf.WriteString(`","date":"`)
		buf.WriteString(bld.Date)
	}

	buf.WriteString(`","goversion":"`)
	buf.WriteString(bld.GoVersion())
	buf.WriteString(`"}`)

	return []byte(buf.String()), nil
}

// String returns the string representation of the build information.
// It always includes the release version. Other fields are omitted when empty.
// Examples:
//  - version only: `v8.0.0`
//  - version and revision `v8.5.0 (fedcba)`
//  - version and date: `v8.5.0 (2020-06-16 19:53)`
//  - all: `v8.5.0 (fedcba @ 2020-06-16 19:53)`
func (bld BuildInfo) String() string {
	if bld.Revision == "" && bld.Date == "" {
		return bld.Version
	}

	var buf strings.Builder
	bld.WriteTo(&buf)
	return buf.String()
}

func (bld BuildInfo) Format(s fmt.State, v rune) {
	switch v {
	case 'V':
		if s.Flag('#') {
			bld.WriteTo(s)
		} else {
			s.Write([]byte(bld.Version))
		}

	case 'R':
		if bld.Revision != "" {
			s.Write([]byte(bld.Revision))
		}

	case 'D':
		if bld.Date != "" {
			s.Write([]byte(bld.Date))
		}

	case 'G':
		s.Write([]byte(bld.GoVersion()))
	}
}

// WriteTo writes a string representation of the build information, similar to
// String, to w. The return value is the number of bytes written. Any error
// encountered during writing is ignored.
func (bld BuildInfo) WriteTo(w io.Writer) (int64, error) {
	c := countingWriter{target: w}
	c.Write([]byte(bld.Version))

	if bld.Revision != "" {
		c.Write([]byte(" (" + bld.Revision))
		if bld.Date != "" {
			c.Write([]byte(" @ " + bld.Date))
		}
		c.Write([]byte(")"))
	} else if bld.Date != "" {
		c.Write([]byte(" (" + bld.Date + ")"))
	}

	return c.size, nil
}

type countingWriter struct {
	target io.Writer
	size   int64
}

func (cw *countingWriter) Write(data []byte) (int, error) {
	n, _ := cw.target.Write(data)
	cw.size += int64(n)
	return n, nil
}

const (
	DummyVersion  = "0.0.0"
	DummyRevision = "abcdef"
	DummyDate     = "1997-08-29 13:37:00"
)

// IsDummy returns true when all fields' values within a BuildInfo are dummy
// values. This may indicate the build information variables are not properly
// set when a new build is made.
func IsDummy(bld BuildInfo) bool {
	return bld.Version == DummyVersion &&
		bld.Revision == DummyRevision &&
		bld.Date == DummyDate
}
