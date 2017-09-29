package librato

import (
	"net/http"
)

// https://www.librato.com/docs/api/?shell#measurements
// Measurements are the individual time series samples sent to Librato. They are
// associated by name with a Metric.

type MeasurementTag struct {
	// Key is the name of the tag
	Key string
	// Value is the value of the tag
	Value string
}

// Measurement corresponds to the Librato API type of the same name
// TODO: support the full set of Measurement fields
type Measurement struct {
	// Name is the name of the Metric this Measurement is associated with
	Name string `json:"name"`
	// Tags add dimensionality to data, similar to Labels in Prometheus
	Tags []*MeasurementTag `json:"tags,omitempty"`
	// Time is the UNIX epoch timestamp of the Measurement
	Time int64 `json:"time"`
	// Value is the value of the
	Value float64 `json:"value"`
}

// MeasurementsCommunicator defines an interface for communicating with the Measurements portion of the Librato API
type MeasurementsCommunicator interface {
	Create([]*Measurement) (*http.Response, error)
}

// MeasurementsService implements MeasurementsCommunicator
type MeasurementsService struct {
	client *Client
}

// Create persists the given MeasurementCollection to Librato
func (ms *MeasurementsService) Create(mc []*Measurement) (*http.Response, error) {
	req, _ := ms.client.NewRequest("POST", "metrics", mc)
	return ms.client.Do(req)
}
