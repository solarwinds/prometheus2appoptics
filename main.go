package main

import (
	"fmt"
	"net/http"
	"os/signal"
	"time"

	"github.com/solarwinds/prometheus2appoptics/internal/app/api"

	"os"

	"github.com/solarwinds/prometheus2appoptics/config"

	"github.com/appoptics/appoptics-api-go"
)

var startTime = time.Now().UTC()
var osSignalChan = make(chan os.Signal, 1)
var stopChan chan<- struct{}

func main() {
	signal.Notify(osSignalChan, os.Interrupt)
	go handleShutdown()

	portString := fmt.Sprintf(":%d", config.BindPort())
	fmt.Println("[-] Starting on ", portString)

	userAgentFragment := fmt.Sprintf("%s", config.AppName)

	aoClient := appoptics.NewClient(config.AccessToken(), appoptics.UserAgentClientOption(userAgentFragment))

	bp := appoptics.NewBatchPersister(aoClient.MeasurementsService(), config.SendStats())
	bp.BatchAndPersistMeasurementsForever()

	stopChan = bp.MeasurementsStopBatchingChannel()

	http.ListenAndServe(portString, api.App(aoClient, bp)) //nolint:errcheck
}

// handleShutdown defines the behavior of the application when it receives SIGINT
func handleShutdown() {
	<-osSignalChan
	runDuration := time.Since(startTime) / time.Second
	fmt.Println("\n[-] Sending stop signal and shutting down")
	fmt.Printf("[-] Process ran for %d seconds\n", runDuration)
	stopChan <- struct{}{}
	os.Exit(0)
}
