package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/common/model"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	promremote "github.com/prometheus/prometheus/storage/remote"
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

		data, err := timeseriesFromRequestData(compressed)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, ts := range data {
			m := make(model.Metric, len(ts.Labels))
			for _, l := range ts.Labels {
				m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
			}
			log.Println(m)

			for _, s := range ts.Samples {
				log.Printf("  %f %d\n", s.Value, s.TimestampMs)
			}
		}
	})
}

// timeseriesFromRequestData produces a slice of Prometheus remote storage TimeSeries from the request
func timeseriesFromRequestData(reqBytes []byte) ([]*promremote.TimeSeries, error) {
	var ts []*promremote.TimeSeries
	reqBuf, err := snappy.Decode(nil, reqBytes)
	if err != nil {
		return ts, err
	}

	var req promremote.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return ts, err
	}
	return req.Timeseries, nil
}
