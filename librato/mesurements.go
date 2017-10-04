package librato

import (
	"fmt"
	"math"
	"net/http"
)

// https://www.librato.com/docs/api/?shell#measurements
// Measurements are the individual time series samples sent to Librato. They are
// associated by name with a Metric.

type MeasurementTags map[string]string

// Measurement corresponds to the Librato API type of the same name
// TODO: support the full set of Measurement fields
type Measurement struct {
	// Name is the name of the Metric this Measurement is associated with
	Name string `json:"name"`
	// Tags add dimensionality to data, similar to Labels in Prometheus
	Tags MeasurementTags `json:"tags,omitempty"`
	// Time is the UNIX epoch timestamp of the Measurement
	Time int64 `json:"time"`
	// Value is the value of the
	Value float64 `json:"value"`
}

// MeasurementPayload is the construct we POST to the API
type MeasurementPayload struct {
	Measurements []*Measurement `json:"measurements"`
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
	payload := MeasurementPayload{mc}
	req, err := ms.client.NewRequest("POST", "measurements", payload)

	if err != nil {
		fmt.Println("error creating request:", err)
		return nil, err
	}
	return ms.client.Do(req, nil)
}

func dumpMeasurements(measurements interface{}) {
	ms := measurements.(MeasurementPayload)
	for i, measurement := range ms.Measurements {
		if math.IsNaN(measurement.Value) {
			fmt.Println("Found at index ", i)
			fmt.Printf("found in '%s'", measurement.Name)
		}
	}
}
