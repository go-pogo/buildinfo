// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"encoding/json"
	"github.com/go-pogo/errors"
	"io"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

//goland:noinspection GoUnusedConst
const (
	// ShortFlag is the default flag to print the current build information.
	ShortFlag = "v"
	// LongFlag is an alternative long version that may be used together with ShortFlag.
	LongFlag = "version"

	// MetricName is the default name for the metric (without namespace).
	MetricName = "buildinfo"
	// MetricHelp is the default help text to describe the metric.
	MetricHelp = "Metric with build information labels and a constant value of '1'."

	// PathPattern is the default path for a http handler.
	PathPattern = "/version"

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
// version, revision and time.
type BuildInfo struct {
	info *debug.BuildInfo

	// AltName is an alternative name for the release.
	AltName string
	// AltVersion is an alternative version of the release.
	AltVersion string
	// Extra additional information to show.
	Extra map[string]string
}

const ErrNoBuildInfo = "no build information available"

// New creates a new BuildInfo with the given altVersion string.
func New(altVersion string) (*BuildInfo, error) {
	bld := BuildInfo{AltVersion: altVersion}
	if !bld.init() {
		return nil, errors.New(ErrNoBuildInfo)
	}
	return &bld, nil
}

func (bld *BuildInfo) init() bool {
	if bld.info != nil {
		return true
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		bld.info = info
		return true
	}
	return false
}

func (bld *BuildInfo) Internal() *debug.BuildInfo { return bld.info }

func (bld *BuildInfo) Setting(key string) string {
	if !bld.init() {
		return ""
	}
	for _, set := range bld.info.Settings {
		if set.Key == key {
			return set.Value
		}
	}
	return ""
}

// GoVersion returns the Go runtime version used to make the current build.
func (bld *BuildInfo) GoVersion() string {
	if !bld.init() || bld.info.GoVersion == "" {
		return runtime.Version()
	}
	return bld.info.GoVersion
}

func (bld *BuildInfo) Name() string {
	if bld.AltName != "" {
		return bld.AltName
	}
	if !bld.init() {
		return ""
	}
	return bld.info.Path[:strings.LastIndex(bld.info.Path, "/")+1]
}

func (bld *BuildInfo) Version() string {
	if bld.AltVersion != "" {
		return bld.AltVersion
	}
	if !bld.init() || bld.info.Main.Version == "" || bld.info.Main.Version == "(devel)" {
		return EmptyVersion
	}
	return bld.info.Main.Version
}

const (
	settingRevision = "vcs.revision"
	settingTime     = "vcs.time"
)

// Revision is the (short) commit hash the release is build from.
func (bld *BuildInfo) Revision() string { return bld.Setting(settingRevision) }

// Time of the commit the release was build.
func (bld *BuildInfo) Time() time.Time {
	t, _ := time.Parse(time.RFC3339, bld.Setting(settingTime))
	return t
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

// Map returns the build information as a map. Field names are lowercase.
// Empty fields are omitted.
func (bld *BuildInfo) Map() map[string]string {
	m := make(map[string]string, 5+len(bld.Extra))
	m[keyVersion] = bld.Version()
	m[keyGoversion] = bld.GoVersion()

	if rev := bld.Revision(); rev != "" {
		m[keyRevision] = rev
	}
	if tim := bld.Time(); !tim.IsZero() {
		m[keyTime] = tim.Format(time.RFC3339)
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
	rev := bld.Revision()
	tim := bld.Time()
	if rev == "" && tim.IsZero() {
		return bld.Version()
	}

	var buf strings.Builder
	_, _ = buf.WriteString(bld.Version())

	if rev != "" {
		_, _ = buf.WriteRune(' ')
		_, _ = buf.WriteString(rev)
	}
	if !tim.IsZero() {
		_, _ = buf.WriteString(" (")
		_, _ = buf.WriteString(tim.Format(time.RFC3339))
		_, _ = buf.WriteString(")")
	}
	return buf.String()
}

var _ json.Marshaler = (*BuildInfo)(nil)

// MarshalJSON returns valid JSON output.
// Empty fields within buildInfo are omitted.
func (bld *BuildInfo) MarshalJSON() ([]byte, error) {
	// WriteString on strings.Builder never returns an error
	var buf strings.Builder
	bld.writeJson(&buf)
	return []byte(buf.String()), nil
}

func (bld *BuildInfo) writeJson(w io.StringWriter) {
	_, _ = w.WriteString(`{"version":"`)
	_, _ = w.WriteString(bld.Version())

	if rev := bld.Revision(); rev != "" {
		_, _ = w.WriteString(`","revision":"`)
		_, _ = w.WriteString(rev)
	}
	if tim := bld.Time(); !tim.IsZero() {
		_, _ = w.WriteString(`","time":"`)
		_, _ = w.WriteString(tim.Format(time.RFC3339))
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
