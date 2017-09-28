package main

import (
	"testing"
	"time"

	"net/http/httptest"

	"net/http"

	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	promremote "github.com/prometheus/prometheus/storage/remote"
)

func TestReceiveHandler(t *testing.T) {
	server := httptest.NewServer(receiveHandler())
	defer server.Close()

	t.Run("data is well-formed", func(t *testing.T) {
		payload := fixtureSamplePayload()
		resp, err := postToReceive(server, payload)

		if err != nil {
			t.Errorf("Expected no error but received %s", err.Error())
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 but received %d", resp.StatusCode)
		}
	})
}

// postToReceive sends the payload bytes to the endpoint via HTTP POST
func postToReceive(server *httptest.Server, payload []byte) (*http.Response, error) {
	client := new(http.Client)
	reader := bytes.NewReader(payload)
	req, _ := http.NewRequest("POST", server.URL+"/receive", reader)
	// header values pulled from Prometheus remote storage client implementation
	req.Header.Add("Content-Encoding", "snappy")
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	return client.Do(req)
}

// fixtureSamplePayload returns a Snappy-compressed TimeSeries
func fixtureSamplePayload() []byte {
	stubLabelPair := &promremote.LabelPair{Name: "environment", Value: "production"}
	stubSample := &promremote.Sample{Value: 123.45, TimestampMs: time.Now().UTC().Unix()}
	stubTimeSeries := promremote.TimeSeries{
		Labels:  []*promremote.LabelPair{stubLabelPair},
		Samples: []*promremote.Sample{stubSample},
	}

	writeRequest := promremote.WriteRequest{Timeseries: []*promremote.TimeSeries{&stubTimeSeries}}

	protoBytes, _ := proto.Marshal(&writeRequest)
	compressedBytes := snappy.Encode(nil, protoBytes)
	return compressedBytes
}
