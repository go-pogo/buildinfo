package buildinfo

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInfo_GoVersion(t *testing.T) {
	assert.Exactly(t, runtime.Version(), BuildInfo{}.GoVersion())
}

func TestBuildInfo_String(t *testing.T) {
	have := BuildInfo{"v1.0.66", "2020-06-16 19:53", "develop", "fedcba"}
	assert.Exactly(t, "v1.0.66, #fedcba @ 2020-06-16 19:53", have.String())
}

func TestBuildInfo_ToMap(t *testing.T) {
	have := BuildInfo{"v1.0.66", DummyDate, "foobar", "fedcba"}
	want := map[string]string{
		"version":   "v1.0.66",
		"date":      DummyDate,
		"branch":    "foobar",
		"commit":    "fedcba",
		"goversion": runtime.Version(),
	}
	assert.Exactly(t, want, have.ToMap())
}

func TestBuildInfo_MarshalJSON(t *testing.T) {
	var x BuildInfo
	haveBytes, haveErr := x.MarshalJSON()
	wantBytes, wantErr := json.Marshal(struct {
		Version   string `json:"version"`
		Date      string `json:"date"`
		Branch    string `json:"branch"`
		Commit    string `json:"commit"`
		GoVersion string `json:"goversion"`
	}{x.Version, x.Date, x.Branch, x.Commit, x.GoVersion()})

	assert.Exactly(t, wantBytes, haveBytes)
	assert.Exactly(t, wantErr, haveErr)
}

func TestIsDummy(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		assert.True(t, IsDummy(BuildInfo{
			Version: DummyVersion,
			Date:    DummyDate,
			Branch:  DummyBranch,
			Commit:  DummyCommit,
		}))
	})
	t.Run("false", func(t *testing.T) {
		assert.False(t, IsDummy(BuildInfo{}))
	})
}
