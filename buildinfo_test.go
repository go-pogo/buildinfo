package buildinfo

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInfo_GoVersion(t *testing.T) {
	assert.Exactly(t, runtime.Version(), new(BuildInfo).GoVersion())
}

func TestBuildInfo_String(t *testing.T) {
	tests := map[string]struct {
		input BuildInfo
		want  string
	}{
		"version only": {
			input: BuildInfo{Version: "v0.12.1"},
			want:  "v0.12.1",
		},
		"version and revision": {
			input: BuildInfo{
				Version:  "v1.0.66",
				Revision: "fedcba",
			},
			want: "v1.0.66 (fedcba)",
		},
		"version and date": {
			input: BuildInfo{
				Version: "0.0.2-rc1",
				Date:    "2020-06-16 19:53",
			},
			want: "0.0.2-rc1 (2020-06-16 19:53)",
		},
		"all": {
			input: BuildInfo{
				Version:  "v1.0.66",
				Revision: "fedcba",
				Date:     "2020-06-16 19:53",
			},
			want: "v1.0.66 (fedcba @ 2020-06-16 19:53)",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, tc.input.String())
		})
	}
}

func TestBuildInfo_Map(t *testing.T) {
	have := BuildInfo{
		Version:  "v1.0.66",
		Revision: "fedcba",
		Date:     DummyDate,
	}

	want := map[string]string{
		"version":   "v1.0.66",
		"revision":  "fedcba",
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
			Date:     DummyDate,
		},
	}

	for name, bld := range tests {
		t.Run(name, func(t *testing.T) {
			haveBytes, haveErr := bld.MarshalJSON()
			wantBytes, wantErr := json.Marshal(struct {
				Version   string `json:"version"`
				Revision  string `json:"revision,omitempty"`
				Date      string `json:"date,omitempty"`
				GoVersion string `json:"goversion"`
			}{bld.Version, bld.Revision, bld.Date, bld.GoVersion()})

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
			Date:     DummyDate,
		}))
	})
	t.Run("false", func(t *testing.T) {
		assert.False(t, IsDummy(BuildInfo{}))
	})
}
