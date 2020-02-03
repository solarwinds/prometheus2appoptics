package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
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

func App(aoClient *appoptics.Client, bp *appoptics.BatchPersister) http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy!")) //nolint:errcheck
	})

	r.Post("/receive", receiveHandler(bp.MeasurementsSink()))
	r.Post("/test", testMetricHandler(aoClient))

	return r
}

// receiveHandler implements the code path for handling incoming Prometheus metrics
func receiveHandler(prepChan chan<- []appoptics.Measurement) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		compressed, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err := processRequestData(compressed)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		convertedData := adapter.PromDataToAppOpticsMeasurements(&data)
		msg := fmt.Sprintf("measurements received - %d", len(convertedData))
		log.Debug(msg)

		prepChan <- convertedData
		w.WriteHeader(http.StatusAccepted)
	}
}

// testMetricHandler sends a single fixture test Metric to AppOptics and is used in debugging
func testMetricHandler(lc appoptics.ServiceAccessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			log.Error("*http.Response was nil")
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)

		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}
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
