package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/model"
	promremote "github.com/prometheus/prometheus/storage/remote"
	"github.com/solarwinds/prometheus2appoptics/promadapter"

	"github.com/appoptics/appoptics-api-go"
)

var adapter promadapter.PrometheusAdapter

func init() {
	adapter = promadapter.NewPromAdapter()
}

// receiveHandler implements the code path for handling incoming Prometheus metrics
func receiveHandler(prepChan chan<- []appoptics.Measurement) http.Handler {
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

		// TODO: make this conditional upon log level
		convertedData := adapter.PromDataToAppOpticsMeasurements(&data)
		log.Println("measurements received - ", len(convertedData))

		prepChan <- convertedData
		w.WriteHeader(http.StatusAccepted)
	})
}

// listSpacesHandler returns the AppOptics Spaces on the associated account and can be used as a test for credentials
func listSpacesHandler(lc appoptics.ServiceAccessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spaces, resp, err := lc.SpacesService().List()

		if err != nil {
			if resp == nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Println(err)
			w.WriteHeader(resp.StatusCode)
			w.Write([]byte(err.Error()))
			return
		}

		for _, space := range spaces {
			fmt.Printf("%+v\n", space)
		}
	})
}

// testMetricHandler sends a single fixture test Metric to AppOptics and is used in debugging
func testMetricHandler(lc appoptics.ServiceAccessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := processRequestData(FixtureSamplePayload())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		mc := adapter.PromDataToAppOpticsMeasurements(&data)
		batch := &appoptics.MeasurementsBatch{
			Measurements: mc,
		}
		resp, err := lc.MeasurementsService().Create(batch)

		if resp == nil {
			log.Println("*http.Response was nil")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)

		if err != nil {
			w.Write([]byte(err.Error()))
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

// FixtureSamplePayload returns a Snappy-compressed TimeSeries
func FixtureSamplePayload() []byte {
	nameLabelPair := &promremote.LabelPair{Name: model.MetricNameLabel, Value: "mah-test-metric"}
	stubLabelPair := &promremote.LabelPair{Name: "environment", Value: "production"}
	stubSample := &promremote.Sample{Value: 123.45, TimestampMs: time.Now().UTC().Unix()}
	stubTimeSeries := promremote.TimeSeries{
		Labels:  []*promremote.LabelPair{stubLabelPair, nameLabelPair},
		Samples: []*promremote.Sample{stubSample},
	}

	writeRequest := promremote.WriteRequest{Timeseries: []*promremote.TimeSeries{&stubTimeSeries}}

	protoBytes, _ := proto.Marshal(&writeRequest)
	compressedBytes := snappy.Encode(nil, protoBytes)
	return compressedBytes
}
