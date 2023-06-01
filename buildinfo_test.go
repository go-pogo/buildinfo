// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"net/http/httptest"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var goVersion = runtime.Version()

func TestNew(t *testing.T) {
	have := New("v1.2.3")
	want := &BuildInfo{Version: "v1.2.3", goVersion: goVersion}
	assert.Exactly(t, want, have)
}

func TestBuildInfo_WithExtra(t *testing.T) {
	t.Run("add value", func(t *testing.T) {
		var bi BuildInfo
		bi.WithExtra("foo", "bar")
		assert.Exactly(t, map[string]string{"foo": "bar"}, bi.Extra)
	})

	reserved := []string{
		keyVersion,
		keyGoversion,
		keyRevision,
		keyCreated,
	}
	for _, key := range reserved {
		t.Run("panic on reserved key "+key, func(t *testing.T) {
			assert.Panics(t, func() {
				var bi BuildInfo
				bi.WithExtra(key, "some value")
			})
		})
	}
}

func TestBuildInfo_GoVersion(t *testing.T) {
	assert.Exactly(t, goVersion, new(BuildInfo).GoVersion())
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
			want: "v1.0.66 fedcba",
		},
		"version and time": {
			input: BuildInfo{
				Version: "0.0.2-rc1",
				Created: time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
			},
			want: "0.0.2-rc1 (2020-06-16T19:53:00Z)",
		},
		"all": {
			input: BuildInfo{
				Version:  "v1.0.66",
				Revision: "fedcba",
				Created:  time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
			},
			want: "v1.0.66 fedcba (2020-06-16T19:53:00Z)",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, tc.input.String())
		})
	}
}

var tests = map[string]struct {
	wantStruct BuildInfo
	wantMap    map[string]string
	wantJson   string
}{
	"empty": {
		wantMap:  map[string]string{"version": EmptyVersion, "goversion": goVersion},
		wantJson: `{"version":"` + EmptyVersion + `","goversion":"` + goVersion + `"}`,
	},
	"partial": {
		wantStruct: BuildInfo{
			Version: "v0.66",
			Created: time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
		},
		wantMap: map[string]string{
			keyVersion:   "v0.66",
			keyGoversion: goVersion,
			keyCreated:   "2020-06-16T19:53:00Z",
		},
		wantJson: `{"version":"v0.66","created":"2020-06-16T19:53:00Z","goversion":"` + goVersion + `"}`,
	},
	"full": {
		wantStruct: BuildInfo{
			Version:  "v0.66",
			Revision: "abcdefghi",
			Created:  time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
		},
		wantMap: map[string]string{
			keyVersion:   "v0.66",
			keyGoversion: goVersion,
			keyRevision:  "abcdefghi",
			keyCreated:   "2020-06-16T19:53:00Z",
		},
		wantJson: `{"version":"v0.66","revision":"abcdefghi","created":"2020-06-16T19:53:00Z","goversion":"` + goVersion + `"}`,
	},
	"extras": {
		wantStruct: BuildInfo{
			Version:  "v0.66",
			Revision: "abcdefghi",
			Created:  time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
			Extra: map[string]string{
				"foo": "bar",
			},
		},
		wantMap: map[string]string{
			keyVersion:   "v0.66",
			keyGoversion: goVersion,
			keyRevision:  "abcdefghi",
			keyCreated:   "2020-06-16T19:53:00Z",
			"foo":        "bar",
		},
		wantJson: `{"version":"v0.66","revision":"abcdefghi","created":"2020-06-16T19:53:00Z","goversion":"` + goVersion + `","foo":"bar"}`,
	},
}

func TestBuildInfo_Map(t *testing.T) {
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.wantMap, tc.wantStruct.Map())
		})
	}
}

func TestBuildInfo_MarshalJSON(t *testing.T) {
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			haveBytes, haveErr := tc.wantStruct.MarshalJSON()

			assert.Exactly(t, []byte(tc.wantJson), haveBytes)
			assert.Nil(t, haveErr)
		})
	}
}

func TestBuildInfo_UnmarshalJSON(t *testing.T) {
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var haveStruct BuildInfo
			haveErr := haveStruct.UnmarshalJSON([]byte(tc.wantJson))

			assert.Exactly(t, tc.wantStruct, haveStruct)
			assert.Nil(t, haveErr)
		})
		t.Run(name, func(t *testing.T) {
			haveStruct := New("")
			haveErr := haveStruct.UnmarshalJSON([]byte(tc.wantJson))

			wantStruct := tc.wantStruct
			wantStruct.goVersion = goVersion

			assert.Exactly(t, wantStruct, *haveStruct)
			assert.Nil(t, haveErr)
		})
	}

	t.Run("empty json values", func(t *testing.T) {
		var have BuildInfo
		haveErr := have.UnmarshalJSON([]byte(`{"version":"","foo":""}`))

		assert.Exactly(t, BuildInfo{}, have)
		assert.Nil(t, haveErr)
	})
}

func TestHttpHandler(t *testing.T) {
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			HttpHandler(&tc.wantStruct).ServeHTTP(rec, nil)
			assert.Exactly(t, []byte(tc.wantJson), rec.Body.Bytes())
		})
	}
}
