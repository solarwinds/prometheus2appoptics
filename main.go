package main

import (
	"fmt"
	"net/http"
	"os/signal"
	"time"

	"os"

	"github.com/solarwinds/p2l/config"
	"github.com/solarwinds/p2l/appoptics"
	"github.com/solarwinds/p2l/promadapter"
)

// startTime helps us collect information on how long this process runs
var startTime = time.Now().UTC()

// osSignalChan is used to handle SIGINT
var osSignalChan = make(chan os.Signal, 1)

// stopChan is used for controlling the batching/sending flow & graceful shutdown
var stopChan = make(chan bool)

func main() {
	signal.Notify(osSignalChan, os.Interrupt)
	go handleShutdown()

	portString := fmt.Sprintf(":%d", config.BindPort())
	fmt.Println("[-] Starting on ", portString)

	lc := appoptics.NewClient(config.AccessEmail(), config.AccessToken())

	// prepChan holds groups of Measurements to be batched
	prepChan := make(chan []*appoptics.Measurement)

	// pushChan holds groups of Measurements conforming to the size constraint described
	// by librato.MeasurementPostMaxBatchSize
	pushChan := make(chan []*appoptics.Measurement)

	// errorChan is used to track persistence errors and shutdown when too many are seen
	errorChan := make(chan error)

	go promadapter.BatchMeasurements(prepChan, pushChan, stopChan)
	go promadapter.PersistBatches(lc, pushChan, stopChan, errorChan)
	go promadapter.ManagePersistenceErrors(errorChan, stopChan)

	http.Handle("/receive", receiveHandler(prepChan))
	http.Handle("/spaces", listSpacesHandler(lc))
	http.Handle("/test", testMetricHandler(lc))

	http.ListenAndServe(portString, nil)
}

// handleShutdown defines the behavior of the application when it receives SIGINT
func handleShutdown() {
	<-osSignalChan
	runDuration := time.Since(startTime) / time.Second
	fmt.Println("\n[-] Sending stop signal and shutting down")
	fmt.Printf("[-] Process ran for %d seconds\n", runDuration)
	stopChan <- true
	os.Exit(0)
}
