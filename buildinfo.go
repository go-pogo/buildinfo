package buildinfo

import (
	"bytes"
	"fmt"
	"runtime"
)

const (
	// ShortFlag is the default flag to print the current build information of the app.
	ShortFlag = "v"
	// LongFlag is an alternative long version that may be used together with ShortFlag.
	LongFlag = "version"

	MetricName = "build_info"
	MetricHelp = "Metric with build information labels and a constant value of '1'."
)

type BuildInfo struct {
	Version string
	Date    string
	Branch  string
	Commit  string
}

// GoVersion returns the version of the used Go runtime. See `runtime.Version()`
// for additional details.
func (bld BuildInfo) GoVersion() string { return runtime.Version() }

// String returns the string representation of the build information.
func (bld BuildInfo) String() string {
	return fmt.Sprintf("%s, #%s @ %s", bld.Version, bld.Commit, bld.Date)
}

// ToMap returns the build information as a strings map.
// Empty fields within BuildInfo are omitted.
func (bld BuildInfo) ToMap() map[string]string {
	m := make(map[string]string, 5)
	m["version"] = bld.Version

	if bld.Date != "" {
		m["date"] = bld.Date
	}
	if bld.Branch != "" {
		m["branch"] = bld.Branch
	}
	if bld.Commit != "" {
		m["commit"] = bld.Commit
	}

	m["goversion"] = bld.GoVersion()
	return m
}

func (bld BuildInfo) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString(`{"version":"`)
	buf.WriteString(bld.Version)
	buf.WriteString(`","date":"`)
	buf.WriteString(bld.Date)
	buf.WriteString(`","branch":"`)
	buf.WriteString(bld.Branch)
	buf.WriteString(`","commit":"`)
	buf.WriteString(bld.Commit)
	buf.WriteString(`","goversion":"`)
	buf.WriteString(bld.GoVersion())
	buf.WriteString(`"}`)

	return buf.Bytes(), nil
}

const (
	DummyVersion = "0.0.0"
	DummyDate    = "1997-08-29 13:37:00"
	DummyBranch  = "HEAD"
	DummyCommit  = "abcdef"
)

// IsDummy returns `true` when all fields' values within a `BuildInfo` are dummy
// values. This may indicate the build information variables are not properly
// overwritten when a new build is made.
func IsDummy(bld BuildInfo) bool {
	return bld.Version == DummyVersion &&
		bld.Date == DummyDate &&
		bld.Branch == DummyBranch &&
		bld.Commit == DummyCommit
}
