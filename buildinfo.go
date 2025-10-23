// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"encoding/json"
	"io"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-pogo/errors"
)

//goland:noinspection GoUnusedConst
const (
	TimeFormat = time.RFC3339

	// ShortFlag is the default flag to print the current build information.
	ShortFlag = "v"
	// LongFlag is an alternative long version that may be used together with ShortFlag.
	LongFlag = "version"

	// MetricName is the default name for the metric (without namespace).
	MetricName = "buildinfo"
	// MetricHelp is the default help text to describe the metric.
	MetricHelp = "Metric with build information labels and a constant value of '1'."

	// PathPattern is the default path for an http handler route.
	PathPattern = "/version"

	// reserved keys
	keyVersion   = "version"
	keyGoVersion = "goversion"
	keyRevision  = "vcs.revision"
	keyTime      = "vcs.time"
	keyModified  = "vcs.modified"
)

func Version(info *debug.BuildInfo) string {
	if info == nil {
		return ""
	}
	return info.Main.Version
}

// Revision returns the (short) commit hash the release is build from according
// to the [debug.BuildSetting] with "vcs.revision" key.
func Revision(info *debug.BuildInfo) string {
	return Setting(info, keyRevision)
}

func Time(info *debug.BuildInfo) time.Time {
	if set := Setting(info, keyTime); set != "" {
		t, _ := time.Parse(TimeFormat, set)
		return t
	}
	return time.Time{}
}

// BuildInfo contains the relevant information of the current release's build
// version, revision and time.
type BuildInfo struct {
	info *debug.BuildInfo

	// GoVersion is the version of the Go toolchain that built the binary
	// (for example, "go1.19.2").
	GoVersion string
	// Version of the release.
	Version string
	// Revision is the (short) commit hash the release is build from.
	Revision string
	// Time of the commit the release is build from.
	Time time.Time
}

// New creates a new BuildInfo with a set version.
func New(version string) *BuildInfo {
	return &BuildInfo{
		GoVersion: runtime.Version(),
		Version:   version,
	}
}

const ErrNoBuildInfo = "no build information available"

// Read creates a new [BuildInfo] by reading the debug buildinfo using
// [debug.ReadBuildInfo] and Setting its keys using the read values.
// It returns an [ErrNoBuildInfo] error when reading is not successful.
func Read() (*BuildInfo, error) {
	var bld BuildInfo
	if !bld.init() {
		return nil, errors.New(ErrNoBuildInfo)
	}

	bld.GoVersion = bld.info.GoVersion
	bld.Version = Version(bld.info)
	bld.Revision = Revision(bld.info)
	bld.Time = Time(bld.info)

	return &bld, nil
}

func (bld *BuildInfo) init() bool {
	if bld.info != nil {
		return bld.info.GoVersion != ""
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		bld.info = info
		return true
	}

	bld.info = new(debug.BuildInfo)
	return false
}

// Internal returns the used [debug.BuildInfo], or tries to read it using
// [debug.ReadBuildInfo] when not already done.
func (bld *BuildInfo) Internal() *debug.BuildInfo {
	if !bld.init() {
		return nil
	}
	return bld.info
}

// AppName returns the name of the application based on the package path of
// the main package for the binary. For example, a package path of
// "golang.org/x/tools/cmd/stringer" returns "stringer".
func (bld *BuildInfo) AppName() string {
	if !bld.init() {
		return ""
	}
	return AppName(bld.info)
}

// AppName returns the name of the application based on the package path of
// the main package for the binary. For example, a package path of
// "golang.org/x/tools/cmd/stringer" returns "stringer".
func AppName(info *debug.BuildInfo) string {
	if info == nil || info.Path == "" || info.Path == "/" {
		return ""
	}
	if i := strings.LastIndex(info.Path, "/"); i >= 0 {
		return info.Path[i+1:]
	}
	return info.Path
}

// Module returns the matching [debug.Module] description of a dependency
// module, both direct and indirect, that contributed packages to the build of
// this binary.
func (bld *BuildInfo) Module(name string) *debug.Module {
	if !bld.init() {
		return nil
	}
	return Module(bld.info, name)
}

// Module returns the matching [debug.Module] description of a dependency
// module, both direct and indirect, that contributed packages to the build of
// this binary.
func Module(info *debug.BuildInfo, name string) *debug.Module {
	if info == nil {
		return nil
	}
	if name == "main" {
		return &info.Main
	}
	for _, mod := range info.Deps {
		if mod.Path == name {
			return mod
		}
	}
	return nil
}

// Setting returns the value of the matching [debug.BuildSetting] used to build
// the binary.
func (bld *BuildInfo) Setting(key string) string {
	if !bld.init() {
		return ""
	}
	return Setting(bld.info, key)
}

// Setting returns the value of the matching [debug.BuildSetting] used to build
// the binary.
func Setting(info *debug.BuildInfo, key string) string {
	if info == nil {
		return ""
	}
	for _, set := range info.Settings {
		if set.Key == key {
			return set.Value
		}
	}
	return ""
}

// Map returns the build information as a map. Field names are lowercase.
// Empty fields are omitted.
func (bld *BuildInfo) Map() map[string]string {
	m := make(map[string]string, 5)
	m[keyVersion] = bld.Version
	if bld.GoVersion != "" {
		m[keyGoVersion] = bld.GoVersion
	}
	if bld.Revision != "" {
		m[keyRevision] = bld.Revision
	}
	if !bld.Time.IsZero() {
		m[keyTime] = bld.Time.Format(TimeFormat)
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
		return bld.Version
	}

	var buf strings.Builder
	_, _ = buf.WriteString(bld.Version)

	if bld.Revision != "" {
		_, _ = buf.WriteRune(' ')
		_, _ = buf.WriteString(bld.Revision)
	}
	if !bld.Time.IsZero() {
		_, _ = buf.WriteString(" (")
		_, _ = buf.WriteString(bld.Time.Format(TimeFormat))
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
	_, _ = w.WriteString(bld.Version)

	if bld.Revision != "" {
		_, _ = w.WriteString(`","revision":"`)
		_, _ = w.WriteString(bld.Revision)
	}
	if !bld.Time.IsZero() {
		_, _ = w.WriteString(`","time":"`)
		_, _ = w.WriteString(bld.Time.Format(TimeFormat))
	}
	if bld.GoVersion != "" {
		_, _ = w.WriteString(`","goversion":"`)
		_, _ = w.WriteString(bld.GoVersion)
	}
	_, _ = w.WriteString(`"}`)
}
