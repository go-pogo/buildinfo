// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/go-pogo/buildinfo/internal"
	"github.com/go-pogo/errors"
)

func run(ctx context.Context, args []string, latest bool) (err error) {
	var tag string
	if latest {
		tag, _, err = internal.LatestTag(ctx)
		if err != nil {
			return err
		}
	} else {
		tag, _, err = internal.CurrentTag(ctx)
		if err != nil {
			return err
		}
	}

	n := len(args)
	if n == 0 {
		_, _ = fmt.Fprint(os.Stdout, tag)
		return nil
	}

	ver, err := semver.NewVersion(tag)
	if err != nil {
		return err
	}

	vars := state{
		ctx:     ctx,
		version: version{v: *ver},
	}

	if n == 1 && strings.Contains(args[0], "{{") && strings.Contains(args[0], "}}") {
		tmpl, err := template.New("").Parse(args[0])
		if err != nil {
			return err
		}
		if err = tmpl.Execute(os.Stdout, vars); err != nil {
			return err
		}
	}

	var sb strings.Builder
	for _, arg := range args {
		if arg == "" {
			continue
		}

		if sb.Len() != 0 {
			sb.WriteRune(' ')
		}

		switch strings.ToLower(arg) {
		case "revision", "rev":
			rev, err := vars.Revision()
			if err != nil {
				return err
			}
			sb.WriteString(rev)

		case "time":
			tim, err := vars.Time()
			if err != nil {
				return err
			}
			sb.WriteString(tim.Format(time.RFC3339))

		case "original":
			sb.WriteString(vars.Original())
		case "full":
			sb.WriteString(vars.FullVersion())
		case "version", "major.minor.patch":
			sb.WriteString(vars.Version())
		case "major.minor":
			sb.WriteString(vars.MajorMinor())
		case "major":
			sb.WriteString(strconv.FormatUint(vars.Major(), 10))
		case "minor":
			sb.WriteString(strconv.FormatUint(vars.Minor(), 10))
		case "patch":
			sb.WriteString(strconv.FormatUint(vars.Patch(), 10))
		// case "prerelease":
		// 	sb.WriteString(vars.Prerelease())
		// case "metadata":
		// 	sb.WriteString(vars.Metadata())
		case "+patch":
			sb.WriteString(vars.version.IncPatch().Version())
		case "+minor":
			sb.WriteString(vars.version.IncMinor().Version())
		case "+major":
			sb.WriteString(vars.version.IncMajor().Version())

		default:
			return errors.Newf("Invalid argument `%s`\n", arg)
		}
	}

	_, _ = fmt.Fprint(os.Stdout, sb.String())
	return nil
}

type state struct {
	version
	ctx      context.Context
	revision string
	time     time.Time
}

func (v *state) details() (err error) {
	if v.revision != "" && !v.time.IsZero() {
		return
	}
	v.revision, v.time, err = internal.TagDetails(v.ctx, v.version.Original())
	return err
}

func (v *state) Revision() (string, error) {
	if err := v.details(); err != nil {
		return "", err
	}
	return v.revision, nil
}

func (v *state) Time() (time.Time, error) {
	if err := v.details(); err != nil {
		return v.time, err
	}
	return v.time, nil
}

type version struct {
	v semver.Version
}

func (v version) Original() string { return v.v.Original() }

func (v version) FullVersion() string { return v.v.String() }

func (v version) Version() string {
	return fmt.Sprintf("%d.%d.%d", v.v.Major(), v.v.Minor(), v.v.Patch())
}

func (v version) MajorMinor() string {
	return fmt.Sprintf("%d.%d", v.v.Major(), v.v.Minor())
}

func (v version) Major() uint64 { return v.v.Major() }

func (v version) Minor() uint64 { return v.v.Minor() }

func (v version) Patch() uint64 { return v.v.Patch() }

func (v version) Prerelease() string { return v.v.Prerelease() }

func (v version) Metadata() string { return v.v.Metadata() }

func (v version) IncPatch() version { return version{v.v.IncPatch()} }

func (v version) IncMinor() version { return version{v.v.IncMinor()} }

func (v version) IncMajor() version { return version{v.v.IncMajor()} }
