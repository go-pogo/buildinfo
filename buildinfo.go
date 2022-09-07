// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package buildinfo provides basic building blocks and instructions to easily add
build and release information to your app. This is done by replacing variables
in main during build with ldflags.

# Usage

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

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

//goland:noinspection GoUnusedConst
const (
	// ShortFlag is the default flag to print the current build information of the app.
	ShortFlag = "v"
	// LongFlag is an alternative long version that may be used together with ShortFlag.
	LongFlag = "version"

	// MetricName is the default name for the metric (without namespace).
	MetricName = "build_info"
	// MetricHelp is the default help text that describes the metric.
	MetricHelp = "Metric with build information labels and a constant value of '1'."

	// Route is the default path for a http handler.
	Route = "/version"
)

// BuildInfo contains the relevant information of the current release's build
// version, revision and build date.
type BuildInfo struct {
	goVersion string
	// Version is the current version of the release.
	Version string
	// Revision is the (short) commit hash the release is build from.
	Revision string
	// Time of when the release was build.
	Time time.Time
}

func (bld *BuildInfo) GoVersion() string {
	if bld.goVersion == "" {
		bld.goVersion = runtime.Version()
	}
	return bld.goVersion
}

// Map returns the build information as a map. Field names are lowercase.
// Empty fields within BuildInfo are omitted.
func (bld *BuildInfo) Map() map[string]string {
	m := make(map[string]string, 5)
	m["version"] = bld.Version

	if bld.Revision != "" {
		m["revision"] = bld.Revision
	}
	if !bld.Time.IsZero() {
		m["time"] = bld.Time.Format(time.RFC3339)
	}

	m["goversion"] = bld.GoVersion()
	return m
}

// MarshalJSON returns valid JSON output.
// Empty fields within BuildInfo are omitted.
func (bld *BuildInfo) MarshalJSON() ([]byte, error) {
	// WriteString on strings.Builder never returns an error
	var buf strings.Builder
	bld.writeJson(&buf)
	return []byte(buf.String()), nil
}

func (bld *BuildInfo) writeJson(w io.StringWriter) {
	_, _ = w.WriteString(`{"version":"`)
	_, _ = w.WriteString(bld.Version)

	if bld.Revision != "" {
		_, _ = w.WriteString(`","revision":"`)
		_, _ = w.WriteString(bld.Revision)
	}
	if !bld.Time.IsZero() {
		_, _ = w.WriteString(`","time":"`)
		_, _ = w.WriteString(bld.Time.Format(time.RFC3339))
	}

	_, _ = w.WriteString(`","goversion":"`)
	_, _ = w.WriteString(bld.GoVersion())
	_, _ = w.WriteString(`"}`)
}

// String returns the string representation of the build information.
// It always includes the release version. Other fields are omitted when empty.
// Examples:
//   - version only: `v8.0.0`
//   - version and revision `v8.5.0 (fedcba)`
//   - version and date: `v8.5.0 (2020-06-16 19:53)`
//   - all: `v8.5.0 (fedcba @ 2020-06-16 19:53)`
func (bld *BuildInfo) String() string {
	if bld.Revision == "" && bld.Time.IsZero() {
		return bld.Version
	}

	var buf strings.Builder
	_, _ = bld.WriteTo(&buf)
	return buf.String()
}

func (bld *BuildInfo) Format(s fmt.State, v rune) {
	switch v {
	case 'V':
		if s.Flag('#') {
			_, _ = bld.WriteTo(s)
		} else {
			_, _ = s.Write([]byte(bld.Version))
		}

	case 'R':
		if bld.Revision != "" {
			_, _ = s.Write([]byte(bld.Revision))
		}

	case 'D':
		if !bld.Time.IsZero() {
			_, _ = s.Write([]byte(bld.Time.Format(time.RFC3339)))
		}

	case 'G':
		_, _ = s.Write([]byte(bld.GoVersion()))
	}
}

// WriteTo writes a string representation of the build information, similar to
// String, to w. The return value is the number of bytes written. Any error
// encountered during writing is ignored.
func (bld *BuildInfo) WriteTo(w io.Writer) (int64, error) {
	cw := countingWriter{target: stringWriter(w)}
	cw.write(bld.Version)

	if bld.Revision != "" {
		cw.write(" (")
		cw.write(bld.Revision)
		if !bld.Time.IsZero() {
			cw.write(" @ ")
			cw.write(bld.Time.Format(time.RFC3339))
		}
		cw.write(")")
	} else if !bld.Time.IsZero() {
		cw.write(" (")
		cw.write(bld.Time.Format(time.RFC3339))
		cw.write(")")
	}

	return int64(cw.size), nil
}

func HttpHandler(bld *BuildInfo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		bld.writeJson(stringWriter(w))
	})
}

func stringWriter(w io.Writer) io.StringWriter {
	if sw, ok := w.(io.StringWriter); ok {
		return sw
	} else {
		return &wrappedWriter{w}
	}
}

type wrappedWriter struct {
	io.Writer
}

func (w *wrappedWriter) WriteString(s string) (int, error) {
	return w.Writer.Write([]byte(s))
}

type countingWriter struct {
	target io.StringWriter
	size   int
	errs   []error
}

func (cw *countingWriter) write(s string) {
	n, err := cw.target.WriteString(s)
	if err != nil {
		cw.errs = append(cw.errs, err)
	}

	cw.size += n
}
