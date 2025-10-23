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

func TestVersion(t *testing.T) {
	t.Run("nil info", func(t *testing.T) {
		assert.Empty(t, Version(nil))
	})
	t.Run("from main module", func(t *testing.T) {
		const want = "v1.0.345"
		assert.Exactly(t, want, Version(&debug.BuildInfo{Main: debug.Module{Version: want}}))
	})
}

func TestRevision(t *testing.T) {
	t.Run("nil info", func(t *testing.T) {
		assert.Empty(t, Revision(nil))
	})
	t.Run("nil settings", func(t *testing.T) {
		assert.Empty(t, Revision(new(debug.BuildInfo)))
	})
	t.Run("from settings", func(t *testing.T) {
		const want = "abcdef"
		assert.Exactly(t, want, Revision(&debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{
					Key:   keyRevision,
					Value: want,
				},
			},
		}))
	})
}

func TestTime(t *testing.T) {
	t.Run("nil info", func(t *testing.T) {
		assert.Exactly(t, time.Time{}, Time(nil))
	})
	t.Run("from settings", func(t *testing.T) {
		assert.Exactly(t,
			time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
			Time(&debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{
						Key:   keyTime,
						Value: "2020-06-16T19:53:00Z",
					},
				},
			}),
		)
	})
}

func TestNew(t *testing.T) {
	have := New("v1.2.3")
	assert.Exactly(t, "v1.2.3", have.Version)
	assert.Exactly(t, goVersion, have.GoVersion)
}

func TestRead(t *testing.T) {
	have, haveErr := Read()
	assert.NoError(t, haveErr)

	info := have.Internal()
	assert.NotNil(t, info)

	assert.Exactly(t, goVersion, have.GoVersion)
	assert.Exactly(t, Version(info), have.Version)
	assert.Exactly(t, Revision(info), have.Revision)
	assert.Exactly(t, Time(info), have.Time)
}

func TestBuildInfo_Internal(t *testing.T) {
	t.Run("on read err", func(t *testing.T) {
		bld := BuildInfo{info: new(debug.BuildInfo)}
		assert.Nil(t, bld.Internal())
	})
}

func TestAppName(t *testing.T) {
	t.Run("nil info", func(t *testing.T) {
		assert.Empty(t, AppName(nil))
	})

	tests := map[string]string{
		"":                            "",
		"/":                           "",
		"/foobar":                     "foobar",
		"barbaz":                      "barbaz",
		"/my-project/cmd/name-of-app": "name-of-app",
	}

	for path, want := range tests {
		t.Run(path, func(t *testing.T) {
			assert.Exactly(t, want, AppName(&debug.BuildInfo{Path: path}))
		})
	}
}

func TestBuildInfo_AppName(t *testing.T) {
	t.Run("after read err", func(t *testing.T) {
		bld := BuildInfo{info: new(debug.BuildInfo)}
		assert.Empty(t, bld.AppName())
	})
	t.Run("after init", func(t *testing.T) {
		have := new(BuildInfo).AppName()
		if info, ok := debug.ReadBuildInfo(); ok {
			assert.Exactly(t, AppName(info), have)
		} else {
			assert.Empty(t, have)
		}
	})
}

func TestModule(t *testing.T) {
	t.Run("nil info", func(t *testing.T) {
		assert.Nil(t, Module(nil, "something"))
	})
	t.Run("main module", func(t *testing.T) {
		info := debug.BuildInfo{
			Main: debug.Module{Path: "some-module", Version: "1.2.3"},
		}
		assert.Exactly(t, &info.Main, Module(&info, "main"))
	})
	t.Run("from deps", func(t *testing.T) {
		want := debug.Module{Path: "some-module", Version: "1.2.3"}
		info := debug.BuildInfo{
			Deps: []*debug.Module{&want},
		}
		assert.Exactly(t, &want, Module(&info, want.Path))
	})
	t.Run("non-existing module", func(t *testing.T) {
		assert.Nil(t, Module(new(debug.BuildInfo), "foobar"))
	})
}

func TestBuildInfo_Module(t *testing.T) {
	t.Run("after read err", func(t *testing.T) {
		bld := BuildInfo{info: new(debug.BuildInfo)}
		assert.Nil(t, bld.Module("main"))
	})
	t.Run("after init", func(t *testing.T) {
		have := new(BuildInfo).Module("main")
		if want, ok := debug.ReadBuildInfo(); ok {
			assert.Exactly(t, want.Main, *have)
		} else {
			assert.Empty(t, *have)
		}
	})
}

func TestSetting(t *testing.T) {
	t.Run("nil info", func(t *testing.T) {
		assert.Empty(t, Setting(nil, "something"))
	})
	t.Run("exists", func(t *testing.T) {
		info := debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{Key: "qux", Value: "xoo"},
			},
		}
		assert.Exactly(t, "xoo", Setting(&info, "qux"))
	})
	t.Run("not exists", func(t *testing.T) {
		assert.Empty(t, Setting(new(debug.BuildInfo), "foobar"))
	})
}

func TestBuildInfo_Setting(t *testing.T) {
	t.Run("after read err", func(t *testing.T) {
		bld := BuildInfo{info: new(debug.BuildInfo)}
		assert.Empty(t, bld.Setting("foobar"))
	})
	t.Run("after init", func(t *testing.T) {
		have := new(BuildInfo).Setting("GOOS")
		if info, ok := debug.ReadBuildInfo(); ok {
			assert.Exactly(t, Setting(info, "GOOS"), have)
		} else {
			assert.NotEmpty(t, have)
		}
	})
}

var outputTests = map[string]struct {
	input      BuildInfo
	wantMap    map[string]string
	wantString string
	wantJson   string
}{
	"empty": {
		wantMap:  map[string]string{"version": ""},
		wantJson: `{"version":""}`,
	},
	"version": {
		input:      BuildInfo{Version: "0.1.0"},
		wantMap:    map[string]string{"version": "0.1.0"},
		wantString: "0.1.0",
		wantJson:   `{"version":"0.1.0"}`,
	},
	"versions": {
		input: BuildInfo{
			GoVersion: goVersion,
			Version:   "0.1.0",
		},
		wantMap:    map[string]string{"version": "0.1.0", "goversion": goVersion},
		wantString: "0.1.0",
		wantJson:   `{"version":"0.1.0","goversion":"` + goVersion + `"}`,
	},
	"version and revision": {
		input: BuildInfo{
			GoVersion: goVersion,
			Version:   "v1.0.66",
			Revision:  "fedcba",
		},
		wantMap: map[string]string{
			keyVersion:   "v1.0.66",
			keyGoVersion: goVersion,
			keyRevision:  "fedcba",
		},
		wantString: "v1.0.66 fedcba",
		wantJson:   `{"version":"v1.0.66","revision":"fedcba","goversion":"` + goVersion + `"}`,
	},
	"version and time": {
		input: BuildInfo{
			GoVersion: goVersion,
			Version:   "0.0.2-rc1",
			Time:      time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
		},
		wantMap: map[string]string{
			keyVersion:   "0.0.2-rc1",
			keyGoVersion: goVersion,
			keyTime:      "2020-06-16T19:53:00Z",
		},
		wantString: "0.0.2-rc1 (2020-06-16T19:53:00Z)",
		wantJson:   `{"version":"0.0.2-rc1","time":"2020-06-16T19:53:00Z","goversion":"` + goVersion + `"}`,
	},
	"full": {
		input: BuildInfo{
			GoVersion: goVersion,
			Version:   "v0.66",
			Revision:  "abcdef",
			Time:      time.Date(2020, 6, 16, 19, 53, 0, 0, time.UTC),
		},
		wantMap: map[string]string{
			keyVersion:   "v0.66",
			keyGoVersion: goVersion,
			keyRevision:  "abcdef",
			keyTime:      "2020-06-16T19:53:00Z",
		},
		wantString: "v0.66 abcdef (2020-06-16T19:53:00Z)",
		wantJson:   `{"version":"v0.66","revision":"abcdef","time":"2020-06-16T19:53:00Z","goversion":"` + goVersion + `"}`,
	},
}

func TestBuildInfo_Map(t *testing.T) {
	for name, tc := range outputTests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.wantMap, tc.input.Map())
		})
	}
}

func TestBuildInfo_String(t *testing.T) {
	for name, tc := range outputTests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.wantString, tc.input.String())
		})
	}
}

func TestBuildInfo_MarshalJSON(t *testing.T) {
	for name, tc := range outputTests {
		t.Run(name, func(t *testing.T) {
			haveBytes, haveErr := tc.input.MarshalJSON()

			assert.Exactly(t, tc.wantJson, string(haveBytes))
			assert.Nil(t, haveErr)
		})
	}
}
