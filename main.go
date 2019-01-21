package main

import (
	"fmt"
	"net/http"
	"os/signal"
	"time"

	"os"

	"github.com/solarwinds/prometheus2appoptics/config"

	"github.com/appoptics/appoptics-api-go"
)

// startTime helps us collect information on how long this process runs
var startTime = time.Now().UTC()

// osSignalChan is used to handle SIGINT
var osSignalChan = make(chan os.Signal, 1)

var stopChan chan<- bool

func main() {

	signal.Notify(osSignalChan, os.Interrupt)
	go handleShutdown()

	portString := fmt.Sprintf(":%d", config.BindPort())
	fmt.Println("[-] Starting on ", portString)

	userAgentFragment := fmt.Sprintf("%s", config.AppName)

	lc := appoptics.NewClient(config.AccessToken(), appoptics.UserAgentClientOption(userAgentFragment))

	bp := appoptics.NewBatchPersister(lc.MeasurementsService(), config.SendStats())
	bp.BatchAndPersistMeasurementsForever()

	stopChan = bp.MeasurementsStopBatchingChannel()

	http.Handle("/receive", receiveHandler(bp.MeasurementsSink()))
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
