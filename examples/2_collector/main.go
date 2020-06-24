package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/roeldev/go-buildinfo"
)

// Default build-time variables. These values are changed via ldflags when
// building a new release.
var (
	version   = buildinfo.DummyVersion
	revision  = buildinfo.DummyRevision
	gitBranch = buildinfo.DummyBranch
	buildDate = buildinfo.DummyDate
)

func main() {
	buildInfo := buildinfo.BuildInfo{
		Version:  version,
		Revision: revision,
		Branch:   gitBranch,
		Date:     buildDate,
	}
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   "example2",
			Name:        buildinfo.MetricName,
			Help:        buildinfo.MetricHelp,
			ConstLabels: buildInfo.Map(),
		},
		func() float64 { return 1 },
	))

	// allow the metrics server to start on a custom port
	var port int
	flag.IntVar(&port, "port", 8090, "Metrics server port")
	flag.Parse()

	// run the metrics server in a seperate go routine so it does not block
	// our main program
	go func() {
		http.Handle("/", promhttp.Handler())
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Printf("\nThe web server is running on `http://localhost:%d`.\n", port)
	fmt.Println("Visit it using your browser to see the build info metric in action.")
	fmt.Println("\n  ", buildInfo.Map())
	fmt.Println()

	// listen for SIGINT or SIGKILL signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, os.Kill)

	// wait here until we actually receive a signal
	// this blocks the main program
	sig := <-signalCh
	fmt.Println("Received signal:", sig)
	fmt.Println("Stopped!")
}
