// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
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
			want: "v1.0.66 (fedcba)",
		},
		"version and time": {
			input: BuildInfo{
				Version: "0.0.2-rc1",
				Time:    time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
			},
			want: "0.0.2-rc1 (2020-06-16T19:53:00Z)",
		},
		"all": {
			input: BuildInfo{
				Version:  "v1.0.66",
				Revision: "fedcba",
				Time:     time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
			},
			want: "v1.0.66 (fedcba @ 2020-06-16T19:53:00Z)",
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
		Time:     time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
	}

	want := map[string]string{
		"version":   "v1.0.66",
		"revision":  "fedcba",
		"time":      "2020-06-16T19:53:00Z",
		"goversion": goVersion,
	}
	assert.Exactly(t, want, have.Map())
}

func TestBuildInfo_MarshalJSON(t *testing.T) {
	tests := map[string]struct {
		input BuildInfo
		want  string
	}{
		"empty": {
			want: `{"version":"","goversion":"` + goVersion + `"}`,
		},
		"partial": {
			input: BuildInfo{
				Version: "v0.66",
				Time:    time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
			},
			want: `{"version":"v0.66","time":"2020-06-16T19:53:00Z","goversion":"` + goVersion + `"}`,
		},
		"full": {
			input: BuildInfo{
				Version:  "v0.66",
				Revision: "abcdefghi",
				Time:     time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
			},
			want: `{"version":"v0.66","revision":"abcdefghi","time":"2020-06-16T19:53:00Z","goversion":"` + goVersion + `"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			haveBytes, haveErr := tc.input.MarshalJSON()

			assert.Exactly(t, []byte(tc.want), haveBytes)
			assert.Nil(t, haveErr)
		})
	}
}
