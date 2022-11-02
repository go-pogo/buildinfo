// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"context"
	"io"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/go-pogo/buildinfo/internal"
	"github.com/go-pogo/errors"
)

// ReaderFunc reads and returns a version string from any source.
type ReaderFunc func(ctx context.Context) (string, error)

// GitTag ReaderFunc returns the repository's current tag.
func GitTag(ctx context.Context) (string, error) {
	tag, _, err := internal.CurrentTag(ctx)
	return tag, err
}

const (
	CallerVar    = "Caller"
	CallerDirVar = "CallerDir"
	PackageVar   = "Package"
	FuncNameVar  = "FuncName"
	VersionVar   = "String"

	DefaultPackage  = "main"
	DefaultFuncName = "getBuildinfo"
	DefaultVersion  = "0.0.0"
)

// Generator
type Generator struct {
	reader         ReaderFunc
	Vars           map[string]interface{}
	DefaultVersion string
}

// New creates a new Generator.
func New(r ReaderFunc, skipCaller uint) *Generator {
	_, caller, _, _ := runtime.Caller(int(skipCaller) + 1)
	return &Generator{
		reader:         r,
		DefaultVersion: DefaultVersion,
		Vars: map[string]interface{}{
			CallerVar:    filepath.Base(caller),
			CallerDirVar: filepath.Dir(caller),
			PackageVar:   DefaultPackage,
			FuncNameVar:  DefaultFuncName,
		},
	}
}

// Read calls ReadContext with context.Background as context.
func (g *Generator) Read() error {
	return g.ReadContext(context.Background())
}

// ReadContext reads from the ReaderFunc and stores any valid version value to
// Vars using key VersionVar.
func (g *Generator) ReadContext(ctx context.Context) error {
	ver, err := g.reader(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	if ver == "" {
		g.Vars[VersionVar] = g.DefaultVersion
	} else {
		g.Vars[VersionVar] = ver
	}
	return nil
}

func (g *Generator) Execute(tmpl *template.Template, w io.Writer) error {
	return g.ExecuteContext(context.Background(), tmpl, w)
}

func (g *Generator) ExecuteContext(ctx context.Context, tmpl *template.Template, w io.Writer) error {
	if g.Vars[VersionVar] == nil || g.Vars[VersionVar] == "" || g.Vars[VersionVar] == g.DefaultVersion {
		if err := g.ReadContext(ctx); err != nil {
			return err
		}
	}
	return tmpl.Execute(w, g.Vars)
}
