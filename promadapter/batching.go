package promadapter

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"bytes"

	"github.com/solarwinds/prometheus2appoptics/config"
	"github.com/librato/appoptics-api-go"

)

// BatchMeasurements reads slices of AppOptics.Measurement types off a channel populated by the web handler
// and packages them into batches conforming to the limitations imposed by the API.
func BatchMeasurements(prepChan <-chan []appoptics.Measurement, batchChan chan<- *appoptics.MeasurementsBatch, stopChan <-chan bool) {
	var ms = []appoptics.Measurement{}
	for {
		select {
		case mslice := <-prepChan:
			ms = append(ms, mslice...)
			if len(ms) >= appoptics.MeasurementPostMaxBatchSize {
				pushBatch := &appoptics.MeasurementsBatch{
					Measurements: ms[:appoptics.MeasurementPostMaxBatchSize],
				}
				batchChan <- pushBatch
				ms = ms[appoptics.MeasurementPostMaxBatchSize:]
			}
		case <-stopChan:
			break
		}
	}
}

// PersistBatches reads maximal slices of AppOptics.Measurement types off a channel and persists them to the remote AppOptics
// API. Errors are placed on the error channel.
func PersistBatches(lc appoptics.ServiceAccessor, batchChan <-chan *appoptics.MeasurementsBatch, stopChan <-chan bool, errorChan chan<- error) {
	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-ticker.C:
			batch := <-batchChan
			err := persistBatch(lc, batch)
			if err != nil {
				errorChan <- err
			}
		case <-stopChan:
			ticker.Stop()
			break
		}
	}
}

// ManagePersistenceErrors tracks errors on the provided channel and sends a stop signal if the ErrorLimit is reached
func ManagePersistenceErrors(errorChan <-chan error, stopChan chan<- bool) {
	var errors []error
	for {
		select {
		case err := <-errorChan:
			errors = append(errors, err)
			if len(errors) > config.PushErrorLimit() {
				stopChan <- true
				break
			}
		}

	}
}

// persistBatch sends to the remote AppOptics endpoint unless config.SendStats() returns false, when it prints to stdout
func persistBatch(lc appoptics.ServiceAccessor, batch *appoptics.MeasurementsBatch) error {
	if config.SendStats() {
		log.Printf("persisting %d Measurements to AppOptics\n", len(batch.Measurements))
		resp, err := lc.MeasurementsService().Create(batch)
		if resp == nil {
			fmt.Println("response is nil")
			return err
		}
		dumpResponse(resp)
	} else {
		printMeasurements(batch.Measurements)
	}
	return nil
}

// printMeasurements pretty-prints the supplied measurements to stdout
func printMeasurements(data []appoptics.Measurement) {
	for _, measurement := range data {
		fmt.Printf("\nMetric name: '%s' \n", measurement.Name)
		fmt.Printf("\t value: %d \n", measurement.Value)
		fmt.Printf("\t\tTags: ")
		for k, v := range measurement.Tags {
			fmt.Printf("\n\t\t\t%s: %s", k, v)
		}
	}
}

func dumpResponse(resp *http.Response) {
	buf := new(bytes.Buffer)
	fmt.Printf("response status: %s\n", resp.Status)
	if resp.Body != nil {
		buf.ReadFrom(resp.Body)
		fmt.Printf("response body: %s\n\n", string(buf.Bytes()))
	}
}
