// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpHandler(t *testing.T) {
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			HTTPHandler(&tc.wantStruct).ServeHTTP(rec, nil)
			assert.Exactly(t, []byte(tc.wantJson), rec.Body.Bytes())
		})
	}
}
