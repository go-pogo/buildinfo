// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var goVersion = runtime.Version()

func TestNew(t *testing.T) {
	have, err := New("v1.2.3")
	assert.Nil(t, err)
	assert.Exactly(t, "v1.2.3", have.AltVersion)
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
			input: BuildInfo{AltVersion: "v0.12.1"},
			want:  "v0.12.1",
		},
		"version and revision": {
			input: BuildInfo{
				info: &debug.BuildInfo{
					Settings: []debug.BuildSetting{
						{Key: keyRevision, Value: "fedcba"},
					},
				},
				AltVersion: "v1.0.66",
			},
			want: "v1.0.66 fedcba",
		},
		"version and time": {
			input: BuildInfo{
				info: &debug.BuildInfo{
					Settings: []debug.BuildSetting{
						{Key: keyTime, Value: time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC).Format(time.RFC3339)},
					},
				},
				AltVersion: "0.0.2-rc1",
			},
			want: "0.0.2-rc1 (2020-06-16T19:53:00Z)",
		},
		"all": {
			input: BuildInfo{
				info: &debug.BuildInfo{
					Settings: []debug.BuildSetting{
						{Key: keyRevision, Value: "fedcba"},
						{Key: keyTime, Value: time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC).Format(time.RFC3339)},
					},
				},
				AltVersion: "v1.0.66",
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
			info: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: keyTime, Value: time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC).Format(time.RFC3339)},
				},
			},
			AltVersion: "v0.66",
		},
		wantMap: map[string]string{
			keyVersion:   "v0.66",
			keyGoversion: goVersion,
			keyTime:      "2020-06-16T19:53:00Z",
		},
		wantJson: `{"version":"v0.66","time":"2020-06-16T19:53:00Z","goversion":"` + goVersion + `"}`,
	},
	"full": {
		wantStruct: BuildInfo{
			info: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: keyRevision, Value: "abcdefghi"},
					{Key: keyTime, Value: time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC).Format(time.RFC3339)},
				},
			},
			AltVersion: "v0.66",
		},
		wantMap: map[string]string{
			keyVersion:   "v0.66",
			keyGoversion: goVersion,
			keyRevision:  "abcdefghi",
			keyTime:      "2020-06-16T19:53:00Z",
		},
		wantJson: `{"version":"v0.66","revision":"abcdefghi","time":"2020-06-16T19:53:00Z","goversion":"` + goVersion + `"}`,
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

			assert.Exactly(t, tc.wantJson, string(haveBytes))
			assert.Nil(t, haveErr)
		})
	}
}
