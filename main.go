package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/model"

	promremote "github.com/prometheus/prometheus/storage/remote"
)

// Flag vars
var bindPort int

func init() {
	flag.IntVar(&bindPort, "bind-port", 4567, "the port the HTTP server binds to")

	flag.Parse()
}

func main() {
	portString := fmt.Sprintf(":%d", bindPort)
	fmt.Println("[-] Starting on ", portString)
	http.Handle("/receive", receiveHandler())
	http.ListenAndServe(portString, nil)
}

// receiveHandler implements the code path for handling incoming Prometheus metrics
func receiveHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		compressed, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err := timeseriesFromRequestData(compressed)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, ts := range data {
			m := make(model.Metric, len(ts.Labels))
			for _, l := range ts.Labels {
				m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
			}
			fmt.Println(m)

			for _, s := range ts.Samples {
				fmt.Printf("  %f %d\n", s.Value, s.TimestampMs)
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
	return ts, nil
}
