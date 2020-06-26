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
	t.Run("only version", func(t *testing.T) {
		want := "v0.12.1"
		assert.Exactly(t, want, BuildInfo{Version: want}.String())
	})
	t.Run("version and branch", func(t *testing.T) {
		assert.Exactly(t, "3.0.2 stable", BuildInfo{Version: "3.0.2", Branch: "stable"}.String())
	})
	t.Run("no branch", func(t *testing.T) {
		assert.Exactly(t, "v1.0.66 (rev fedcba, date 2020-06-16 19:53)", BuildInfo{
			Version:  "v1.0.66",
			Revision: "fedcba",
			Date:     "2020-06-16 19:53",
		}.String())
	})
	t.Run("all", func(t *testing.T) {
		assert.Exactly(t, "v1.0.66 develop (rev fedcba, date 2020-06-16 19:53)", BuildInfo{
			Version:  "v1.0.66",
			Revision: "fedcba",
			Branch:   "develop",
			Date:     "2020-06-16 19:53",
		}.String())
	})
}

func TestBuildInfo_Map(t *testing.T) {
	have := BuildInfo{
		Version:  "v1.0.66",
		Revision: "fedcba",
		Branch:   "foobar",
		Date:     DummyDate,
	}

	want := map[string]string{
		"version":   "v1.0.66",
		"rev":       "fedcba",
		"branch":    "foobar",
		"date":      DummyDate,
		"goversion": runtime.Version(),
	}
	assert.Exactly(t, want, have.Map())
}

func TestBuildInfo_MarshalJSON(t *testing.T) {
	tests := map[string]BuildInfo{
		"empty":   {},
		"partial": {Version: DummyVersion, Date: DummyDate},
		"full": {
			Version:  DummyVersion,
			Revision: DummyRevision,
			Branch:   DummyBranch,
			Date:     DummyDate,
		},
	}

	for name, x := range tests {
		t.Run(name, func(t *testing.T) {
			haveBytes, haveErr := x.MarshalJSON()
			wantBytes, wantErr := json.Marshal(struct {
				Version   string `json:"version"`
				Revision  string `json:"rev,omitempty"`
				Branch    string `json:"branch,omitempty"`
				Date      string `json:"date,omitempty"`
				GoVersion string `json:"goversion"`
			}{x.Version, x.Revision, x.Branch, x.Date, x.GoVersion()})

			assert.Exactly(t, wantBytes, haveBytes)
			assert.Exactly(t, wantErr, haveErr)
		})
	}
}

func TestIsDummy(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		assert.True(t, IsDummy(BuildInfo{
			Version:  DummyVersion,
			Revision: DummyRevision,
			Branch:   DummyBranch,
			Date:     DummyDate,
		}))
	})
	t.Run("false", func(t *testing.T) {
		assert.False(t, IsDummy(BuildInfo{}))
	})
}
