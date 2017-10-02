package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	promremote "github.com/prometheus/prometheus/storage/remote"
	"github.com/solarwinds/p2l/config"
	"github.com/solarwinds/p2l/librato"
	"github.com/solarwinds/p2l/promadapter"
)

// incremented monotonically when a push to Librato results in an error
var globalPushErrorCounter int

// TODO: more robust/intelligent error handling
// receiveHandler implements the code path for handling incoming Prometheus metrics
func receiveHandler(lc librato.ServiceAccessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		compressed, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err := processRequestData(compressed)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if globalPushErrorCounter < config.PushErrorLimit() {
			mc := promadapter.PromDataToLibratoMeasurements(&data)

			// Either persist to Librato or print to stdout depending on how the app was started
			if config.SendStats() {
				resp, err := lc.MeasurementsService().Create(mc)
				if err != nil || resp.StatusCode > 399 {
					log.Println(err)
					globalPushErrorCounter++
					return
				}
				// TODO: get the real one and just look for it
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					log.Printf("Sent %d metrics to Librato\n", len(mc))
				}
			} else {
				printMeasurements(mc)
			}
		} else {
			log.Fatalln("too many errors - exiting")
		}
	})
}

// listSpacesHandler returns the Librato Spaces on the associated account and can be used as a test for credentials
func listSpacesHandler(lc librato.ServiceAccessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaces, err := lc.SpacesService().List()

		fmt.Println("IN THE HANDLER")

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, space := range spaces {
			fmt.Printf("%+v\n", space)
		}
	})
}

// processRequestData returns a Prometheus remote storage WriteRequest from the raw HTTP body data
func processRequestData(reqBytes []byte) (promremote.WriteRequest, error) {
	var req promremote.WriteRequest
	reqBuf, err := snappy.Decode(nil, reqBytes)
	if err != nil {
		return req, err
	}

	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return req, err
	}
	return req, nil
}

func printMeasurements(data []*librato.Measurement) {
	for _, measurement := range data {
		fmt.Printf("\nMetric name: '%s' \n", measurement.Name)
		fmt.Printf("\t\tTags: ")
		for _, tag := range measurement.Tags {
			fmt.Printf("\n\t\t\t%s: %s", tag.Key, tag.Value)
		}
	}
}
