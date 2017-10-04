package promadapter

import (
	"fmt"
	"log"
	"time"

	"github.com/solarwinds/p2l/config"
	"github.com/solarwinds/p2l/librato"
)

// BatchMeasurements reads slices of librato.Measurement types off a channel populated by the web handler
// and packages them into batches conforming to the limitations imposed by the API.
func BatchMeasurements(prepChan <-chan []*librato.Measurement, pushChan chan<- []*librato.Measurement, stopChan <-chan bool) {
	var currentBatch []*librato.Measurement
	for {
		select {
		case mslice := <-prepChan:
			currentBatch = append(currentBatch, mslice...)
			if len(currentBatch) >= librato.MeasurementPostMaxBatchSize {
				pushBatch := currentBatch[:librato.MeasurementPostMaxBatchSize]
				pushChan <- pushBatch
				currentBatch = currentBatch[librato.MeasurementPostMaxBatchSize:]
			}
		case <-stopChan:
			break
		}
	}
}

// PersistBatches reads maximal slices of librato.Measurement types off a channel and persists them to the remote Librato
// API. Errors are placed on the error channel.
func PersistBatches(lc librato.ServiceAccessor, pushChan <-chan []*librato.Measurement, stopChan <-chan bool, errorChan chan<- error) {
	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-ticker.C:
			batch := <-pushChan
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

// persistBatch sends to the remote Librato endpoint unless config.SendStats() returns false, when it prints to stdout
func persistBatch(lc librato.ServiceAccessor, batch []*librato.Measurement) error {
	if config.SendStats() {
		_, err := lc.MeasurementsService().Create(batch)
		log.Printf("persisting %d Measurements to Librato\n", len(batch))
		return err
	} else {
		printMeasurements(batch)
		return nil
	}
}

// printMeasurements pretty-prints the supplied measurements to stdout
func printMeasurements(data []*librato.Measurement) {
	for _, measurement := range data {
		fmt.Printf("\nMetric name: '%s' \n", measurement.Name)
		fmt.Printf("\t\tTags: ")
		for k, v := range measurement.Tags {
			fmt.Printf("\n\t\t\t%s: %s", k, v)
		}
	}
}
