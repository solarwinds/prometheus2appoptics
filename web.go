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

// receiveHandler implements the code path for handling incoming Prometheus metrics
func receiveHandler() http.Handler {
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

		mc := promadapter.PromDataToLibratoMeasurements(&data)

		if config.SendStats() {
			log.Println("DEF SENDING STATS")
			log.Println(config.AccessEmail())
			log.Println(config.AccessToken())
		} else {
			log.Println("NOT SENDING STATS")
			printMeasurements(mc)
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
