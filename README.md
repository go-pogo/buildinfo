go-buildinfo
============

[![Latest release][latest-release-img]][latest-release-url]
[![Tests status][tests-status-img]][tests-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/roeldev/go-buildinfo.svg?label=latest
[latest-release-url]: https://github.com/roeldev/go-buildinfo/releases
[tests-status-img]: https://github.com/roeldev/go-buildinfo/workflows/Run%20tests/badge.svg
[tests-status-url]: https://github.com/roeldev/go-buildinfo/actions?query=workflow%3A%22Run+tests%22
[report-img]: https://goreportcard.com/badge/github.com/roeldev/go-buildinfo
[report-url]: https://goreportcard.com/report/github.com/roeldev/go-buildinfo
[doc-img]: https://godoc.org/github.com/roeldev/go-buildinfo?status.svg
[doc-url]: https://pkg.go.dev/github.com/roeldev/go-buildinfo


This package provides an easy way to add build/version information to your application.
By providing   
Usage of the revision, branch and date fields is not mandatory but is recommended.

```sh
go get -u github.com/roeldev/go-buildinfo
```
```go
import "github.com/roeldev/go-buildinfo"
```

## Basic usage

```go
package main

import "github.com/roeldev/go-buildinfo"

// these values are changed via ldflags when building a new release
var (
	version = buildinfo.DummyVersion
	revision = buildinfo.DummyRevision
	branch = buildinfo.DummyBranch
	date = buildinfo.DummyDate
)

func main() {
	bld := buildinfo.BuildInfo{
		Version:  version,
		Revision: revision,
		Branch:   branch,
		Date:     date,
	}
}
```

Build your Go project and include the following _ldflags_:
```sh
go build -ldflags=" \
  -X main.version=`git describe --tags` \
  -X main.revision=`git rev-parse --short HEAD` \
  -X main.branch=`git branch --show-current` \
  -X main.date=`date +%FT%T%z`" \
  main.go
```

The resulting build info string is something like `v0.1.0 master (rev 328396d, date 1997-08-29 13:37:00)`.
Any empty fields, but version, are ommitted. 


## Prometheus metric collector
It is often a good idea to make the build information of your app available for Prometheus to collect. Below example shows just how easy it is to create and register a collector with the build information as constant labels.
```go
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
Additional detailed documentation is available at [go.dev][doc-url]


### Created with
<a href="https://www.jetbrains.com/?from=roeldev" target="_blank"><img src="https://pbs.twimg.com/profile_images/1206615658638856192/eiS7UWLo_400x400.jpg" width="35" /></a>


## License
[GPL-3.0+](LICENSE) © 2020 [Roel Schut](https://roelschut.nl)
