// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"net/http"

	"github.com/go-pogo/writing"
)

// HTTPHandler is the http.Handler that writes BuildInfo bld as a JSON response
// to the http response.
func HTTPHandler(bld *BuildInfo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		h := w.Header()
		h.Set("Content-Type", "application/json")
		if !bld.Time.IsZero() {
			h.Set("Last-Modified", bld.Time.Format(http.TimeFormat))
		}
		bld.writeJson(writing.ToStringWriter(w))
	})
}
