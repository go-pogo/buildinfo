buildinfo
=========

[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/go-pogo/buildinfo.svg?label=latest
[latest-release-url]: https://github.com/go-pogo/buildinfo/releases
[build-status-img]: https://github.com/go-pogo/buildinfo/workflows/Test/badge.svg
[build-status-url]: https://github.com/go-pogo/buildinfo/actions?query=workflow%3Test
[report-img]: https://goreportcard.com/badge/github.com/go-pogo/buildinfo
[report-url]: https://goreportcard.com/report/github.com/go-pogo/buildinfo
[doc-img]: https://godoc.org/github.com/go-pogo/buildinfo?status.svg
[doc-url]: https://pkg.go.dev/github.com/go-pogo/buildinfo

Package `buildinfo` provides basic building blocks and instructions to easily add
build and release information to your app. This is done by replacing variables
in main during build with `ldflags`.

```sh
go get github.com/go-pogo/buildinfo
```
```go
import "github.com/go-pogo/buildinfo"
```

## Usage

Declare build info variables in your main package:
```
package main

// these values are changed via ldflags when building a new release
var (
    version = buildinfo.DummyVersion
    revision = buildinfo.DummyRevision
    date = buildinfo.DummyDate
)

func main() {
    bld := buildinfo.BuildInfo{
        Version:  version,
        Revision: revision,
        Date:     date,
    }
}
```
Build your Go project and include the following `ldflags`:
```
go build -ldflags=" \
  -X main.version=`$(git describe --tags)` \
  -X main.revision=`$(git rev-parse --short HEAD)` \
  -X main.date=`$(date +%FT%T%z`)" \
  main.go
```

## Prometheus metric collector

When using a metrics scraper like Prometheus, it is often a good idea to make
the build information of your app available. Below example shows just how easy
it is to create and register a collector with the build information as
constant labels.
```
prometheus.MustRegister(prometheus.NewGaugeFunc(
    prometheus.GaugeOpts{
        Namespace:   "myapp",
        Name:        buildinfo.MetricName,
        Help:        buildinfo.MetricHelp,
        ConstLabels: bld.Map(),
    },
    func() float64 { return 1 },
))
```

## Documentation
Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Created with
<a href="https://www.jetbrains.com/?from=go-pogo" target="_blank"><img src="https://pbs.twimg.com/profile_images/1206615658638856192/eiS7UWLo_400x400.jpg" width="35" /></a>

## License
Copyright © 2020 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
