buildinfo
=========

[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/go-pogo/buildinfo.svg?label=latest

[latest-release-url]: https://github.com/go-pogo/buildinfo/releases

[build-status-img]: https://github.com/go-pogo/buildinfo/actions/workflows/test.yml/badge.svg

[build-status-url]: https://github.com/go-pogo/buildinfo/actions/workflows/test.yml

[report-img]: https://goreportcard.com/badge/github.com/go-pogo/buildinfo

[report-url]: https://goreportcard.com/report/github.com/go-pogo/buildinfo

[doc-img]: https://godoc.org/github.com/go-pogo/buildinfo?status.svg

[doc-url]: https://pkg.go.dev/github.com/go-pogo/buildinfo

Package `buildinfo` provides basic building blocks and instructions to easily add
build and release information to your app.

```sh
go get github.com/go-pogo/buildinfo
```

```
import "github.com/go-pogo/buildinfo"
```

## Using ldflags

Declare build info variables in your main package:

```
package main

// this value is changed via ldflags when building a new release
var version string

func main() {
    bld, err := buildinfo.New(version)
}
```

Build your Go project and include the following `ldflags`:

```sh
go build -ldflags=" \
  -X main.version=`$(git describe --tags)` \
  main.go
```

## Observability usage

When using a metrics scraper like Prometheus or OpenTelemetry, it is often a
good idea to make the build information of your app available. Below example
shows just how easy it is to create and register a collector with the build
information as constant labels.

### Prometheus metric collector

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

### OTEL resource

```
resource.Merge(
    resource.Default(),
    resource.NewSchemaless(
        semconv.ServiceName("myapp"),
        semconv.ServiceVersion(bld.Version()),
        attribute.String("vcs.revision", bld.Revision()),
        attribute.String("vcs.time", bld.Time().Format(time.RFC3339)),
    ),
)
```

## Documentation

Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Created with

<a href="https://www.jetbrains.com/?from=go-pogo" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" width="35" /></a>

## License

Copyright © 2020-2024 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
