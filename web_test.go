package main

import (
	"testing"

	"net/http/httptest"

	"net/http"

	"bytes"

	"github.com/solarwinds/p2l/librato"
)

func TestReceiveHandler(t *testing.T) {
	// simple hack to ensure we don't block forever
	prepChan := make(chan []*librato.Measurement)
	go func(prepChan <-chan []*librato.Measurement) {
		_ = <-prepChan
	}(prepChan)

	server := httptest.NewServer(receiveHandler(prepChan))
	defer server.Close()

	t.Run("data is well-formed", func(t *testing.T) {
		payload := FixtureSamplePayload()
		resp, err := postToReceive(server, payload)

		if err != nil {
			t.Errorf("Expected no error but received %s", err.Error())
		}

		if resp.StatusCode != http.StatusAccepted {
			t.Errorf("Expected status 202 but received %d", resp.StatusCode)
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
