package librato

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	promremote "github.com/prometheus/prometheus/storage/remote"
)

// ServiceAccessor defines an interface for talking to Librato
type ServiceAccessor interface {
	SendPrometheusMetrics([]*promremote.TimeSeries) (*http.Response, error)
}

const (
	defaultBaseURL   = "https://metrics-api.librato.com/v1"
	defaultMediaType = "application/json"
)

// Client implements ServiceAccessor
type Client struct {
	// baseURL is the base endpoint of the remote Librato service
	baseURL *url.URL
	// client is the http.Client singleton used for wire interaction
	client *http.Client
	// email is the public part of the API credential pair
	email string
	// token is the private part of the API credential pair
	token string
	// measurementsService embeds the client and implements access to the Measurements API
	MeasurementsService MeasurementsCommunicator
}

func NewClient(email, token string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:  new(http.Client),
		email:   email,
		token:   token,
		baseURL: baseURL,
	}
	c.MeasurementsService = &MeasurementsService{c}
	return c
}

// NewRequest standardizes the request being sent
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	requestURL := c.baseURL.ResolveReference(rel)

	var buffer io.ReadWriter

	if body != nil {
		buffer = &bytes.Buffer{}
		encodeErr := json.NewEncoder(buffer).Encode(body)
		if encodeErr != nil {
			return nil, encodeErr
		}

	}
	req, err := http.NewRequest(method, requestURL.String(), buffer)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Accept", defaultMediaType)
	req.Header.Set("Content-Type", defaultMediaType)

	return req, nil
}

// TODO: use this as a way to standardize error responses
// Do performs the HTTP request on the wire
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
