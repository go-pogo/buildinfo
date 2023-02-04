// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/writing"
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

	// reserved keys
	keyVersion   = "version"
	keyGoversion = "goversion"
	keyRevision  = "revision"
	keyTime      = "time"
)

var EmptyVersion = "0.0.0"

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
	// Extra additional information to show.
	Extra map[string]string
}

const panicReservedKey = "buildinfo: cannot add reserved key "

// WithExtra adds an extra key value pair.
func (bld *BuildInfo) WithExtra(key, value string) *BuildInfo {
	if key == keyVersion || key == keyGoversion || key == keyRevision || key == keyTime {
		panic(panicReservedKey + key)
	}
	if bld.Extra == nil {
		bld.Extra = make(map[string]string, 2)
	}

	bld.Extra[key] = value
	return bld
}

func (bld *BuildInfo) version() string {
	if bld.Version == "" {
		return EmptyVersion
	}
	return bld.Version
}

// GoVersion returns the Go runtime version used to make the current build.
func (bld *BuildInfo) GoVersion() string {
	if bld.goVersion == "" {
		bld.goVersion = runtime.Version()
	}
	return bld.goVersion
}

// Map returns the build information as a map. Field names are lowercase.
// Empty fields within BuildInfo are omitted.
func (bld *BuildInfo) Map() map[string]string {
	m := make(map[string]string, 5+len(bld.Extra))
	m[keyVersion] = bld.version()
	m[keyGoversion] = bld.GoVersion()

	if bld.Revision != "" {
		m[keyRevision] = bld.Revision
	}
	if !bld.Time.IsZero() {
		m[keyTime] = bld.Time.Format(time.RFC3339)
	}

	for key, val := range bld.Extra {
		m[key] = val
	}
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
	_, _ = w.WriteString(bld.version())

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

	for key, val := range bld.Extra {
		_, _ = w.WriteString(`","`)
		_, _ = w.WriteString(key)
		_, _ = w.WriteString(`":"`)
		_, _ = w.WriteString(val)
	}

	_, _ = w.WriteString(`"}`)
}

// String returns the string representation of the build information.
// It always includes the release version. Other fields are omitted when empty.
// Examples:
//   - version only: `8.5.0`
//   - version and revision `8.5.0 (fedcba)`
//   - version and date: `8.5.0 (2020-06-16 19:53)`
//   - all: `8.5.0 (fedcba @ 2020-06-16 19:53)`
func (bld *BuildInfo) String() string {
	if bld.Revision == "" && bld.Time.IsZero() {
		return bld.version()
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
			_, _ = io.WriteString(s, bld.version())
		}

	case 'R':
		if bld.Revision != "" {
			_, _ = io.WriteString(s, bld.Revision)
		}

	case 'D':
		if !bld.Time.IsZero() {
			_, _ = io.WriteString(s, bld.Time.Format(time.RFC3339))
		}

	case 'G':
		_, _ = io.WriteString(s, bld.GoVersion())
	}
}

// WriteTo writes a string representation of the build information, similar to
// String, to w. The return value is the number of bytes written. Any error
// encountered during writing is ignored.
func (bld *BuildInfo) WriteTo(w io.Writer) (int64, error) {
	cw := writing.ToCountingStringWriter(w)
	_, _ = cw.WriteString(bld.version())

	if bld.Revision != "" {
		_, _ = cw.WriteString(" (")
		_, _ = cw.WriteString(bld.Revision)
		if !bld.Time.IsZero() {
			_, _ = cw.WriteString(" @ ")
			_, _ = cw.WriteString(bld.Time.Format(time.RFC3339))
		}
		_, _ = cw.WriteString(")")
	} else if !bld.Time.IsZero() {
		_, _ = cw.WriteString(" (")
		_, _ = cw.WriteString(bld.Time.Format(time.RFC3339))
		_, _ = cw.WriteString(")")
	}

	return int64(cw.Count()), errors.Join(cw.Errors()...)
}

// HttpHandler is the http.Handler that writes BuildInfo bld as a json response
// to the received request.
func HttpHandler(bld *BuildInfo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		bld.writeJson(writing.ToStringWriter(w))
	})
}
