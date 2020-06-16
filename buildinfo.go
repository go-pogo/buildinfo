package buildinfo

import (
	"bytes"
	"fmt"
	"runtime"
)

const (
	ShortFlag = "v"
	LongFlag  = "version"

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

// ToMap
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
