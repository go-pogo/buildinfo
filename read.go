// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"encoding/json"
	"github.com/go-pogo/errors"
	"io"
	"io/fs"
	"os"
)

// Read reads from io.Reader r and json unmarshalls it's content into a new
// BuildInfo.
func Read(r io.Reader) (*BuildInfo, error) {
	var bld BuildInfo
	bld.init()

	if err := json.NewDecoder(r).Decode(&bld); err != nil {
		return nil, errors.WithStack(err)
	}

	return &bld, nil
}

// Open the file, then read and decode its contents using Read.
func Open(file string) (*BuildInfo, error) {
	return OpenFS(os.DirFS(""), file)
}

// OpenFS opens file from fsys. It then reads and decodes the file's contents
// using Read.
func OpenFS(fsys fs.FS, file string) (bld *BuildInfo, err error) {
	f, err := fsys.Open(file)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	defer errors.AppendFunc(&err, f.Close)
	bld, err = Read(f)
	return
}
