buildinfo
=========

[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/go-pogo/buildinfo.svg?label=latest
[latest-release-url]: https://github.com/go-pogo/buildinfo/releases
[build-status-img]: https://github.com/go-pogo/buildinfo/workflows/Go/badge.svg
[build-status-url]: https://github.com/go-pogo/buildinfo/actions?query=workflow%3AGo
[report-img]: https://goreportcard.com/badge/github.com/go-pogo/buildinfo
[report-url]: https://goreportcard.com/report/github.com/go-pogo/buildinfo
[doc-img]: https://godoc.org/github.com/go-pogo/buildinfo?status.svg
[doc-url]: https://pkg.go.dev/github.com/go-pogo/buildinfo


```sh
go get github.com/go-pogo/buildinfo
```
```go
import "github.com/go-pogo/buildinfo"
```

## Basic usage

```go
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

Build your Go project and include the following _ldflags_:
```sh
go build -ldflags=" \
  -X main.version=`git describe --tags` \
  -X main.revision=`git rev-parse --short HEAD` \
  -X main.date=`date +%FT%T%z`" \
  main.go
```


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
<a href="https://www.jetbrains.com/?from=go-pogo" target="_blank"><img src="https://pbs.twimg.com/profile_images/1206615658638856192/eiS7UWLo_400x400.jpg" width="35" /></a>


## License
[GPL-3.0+](LICENSE) © 2020 [Roel Schut](https://roelschut.nl)
