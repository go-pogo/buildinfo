// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
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

func isReservedKey(key string) bool {
	return key == keyVersion || key == keyGoversion || key == keyRevision || key == keyTime
}

// EmptyVersion is the default version string when no version is set.
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

// New creates a new *BuildInfo with the given version string.
func New(ver string) *BuildInfo {
	bld := &BuildInfo{Version: ver}
	bld.init()
	return bld
}

// Read reads from io.Reader r and json unmarshalls it's content into a new
// BuildInfo.
func Read(r io.Reader) (*BuildInfo, error) {
	var bld BuildInfo
	bld.init()

	if err := json.NewDecoder(r).Decode(&bld); err != nil {
		return nil, errors.WithStack(err)
	}

	return &bld, nil
}

// Open the file, then read and decode its contents using Read.
func Open(file string) (*BuildInfo, error) {
	return OpenFS(os.DirFS(""), file)
}

// OpenFS opens file from fsys. It then reads and decodes its contents using
// Read.
func OpenFS(fsys fs.FS, file string) (*BuildInfo, error) {
	f, err := fsys.Open(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer f.Close()
	return Read(f)
}

const panicReservedKey = "buildinfo: cannot add reserved key "

// WithExtra adds an extra key value pair.
func (bld *BuildInfo) WithExtra(key, value string) *BuildInfo {
	if isReservedKey(key) {
		panic(panicReservedKey + key)
	}

	bld.withExtra(key, value)
	return bld
}

func (bld *BuildInfo) withExtra(key, value string) {
	if bld.Extra == nil {
		bld.Extra = make(map[string]string, 2)
	}

	bld.Extra[key] = value
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
		if !isReservedKey(key) {
			m[key] = val
		}
	}
	return m
}

// String returns the string representation of the build information.
// It always includes the release version. Other fields are omitted when empty.
// Examples:
//   - version only: `8.5.0`
//   - version and revision `8.5.0 (#fedcba)`
//   - version and date: `8.5.0 (2020-06-16 19:53)`
//   - all: `8.5.0 (#fedcba @ 2020-06-16 19:53)`
func (bld *BuildInfo) String() string {
	if bld.Revision == "" && bld.Time.IsZero() {
		return bld.version()
	}

	var buf strings.Builder
	_, _ = buf.WriteString(bld.version())

	if bld.Revision != "" {
		_, _ = buf.WriteRune(' ')
		_, _ = buf.WriteString(bld.Revision)
	}
	if !bld.Time.IsZero() {
		_, _ = buf.WriteString(" (")
		_, _ = buf.WriteString(bld.Time.Format(time.RFC3339))
		_, _ = buf.WriteString(")")
	}
	return buf.String()
}

var (
	_ json.Marshaler   = new(BuildInfo)
	_ json.Unmarshaler = new(BuildInfo)
)

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

func (bld *BuildInfo) UnmarshalJSON(bytes []byte) error {
	fields := make(map[string]string, 0)
	if err := json.Unmarshal(bytes, &fields); err != nil {
		return errors.WithStack(err)
	}

	for k, v := range fields {
		switch k {
		case keyGoversion:
			continue
		case keyVersion:
			bld.Version = v
		case keyRevision:
			bld.Revision = v
		case keyTime:
			var err error
			bld.Time, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return errors.WithStack(err)
			}
		default:
			bld.withExtra(k, v)
		}
	}
	return nil
}

// HttpHandler is the http.Handler that writes BuildInfo bld as a json response
// to the received request.
func HttpHandler(bld *BuildInfo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		bld.writeJson(writing.ToStringWriter(w))
	})
}
