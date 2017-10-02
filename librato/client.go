package librato

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// ServiceAccessor defines an interface for talking to Librato via domain-specific service constructs
type ServiceAccessor interface {
	// MeasurementsService implements an interface for dealing with Librato Measurements
	MeasurementsService() MeasurementsCommunicator
	// SpacesService implements an interface for dealing with Librato Spaces
	SpacesService() SpacesCommunicator
}

const (
	defaultBaseURL   = "https://metrics-api.librato.com/v1/"
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
	measurementsService MeasurementsCommunicator
	// spacesService embeds the client and implements access to the Spaces API
	spacesService SpacesCommunicator
}

func NewClient(email, token string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:  new(http.Client),
		email:   email,
		token:   token,
		baseURL: baseURL,
	}
	c.measurementsService = &MeasurementsService{c}
	c.spacesService = &SpacesService{c}

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

// MeasurementsService represents the subset of the API that deals with Librato Measurements
func (c *Client) MeasurementsService() MeasurementsCommunicator {
	return c.measurementsService
}

// SpacesService represents the subset of the API that deals with Librato Measurements
func (c *Client) SpacesService() SpacesCommunicator {
	return c.spacesService
}

// TODO: use this as a way to standardize error responses
// Do performs the HTTP request on the wire
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
