package librato

import (
	"encoding/json"
)

// SpacesResponse represents the returned data payload from Spaces API's List command (/spaces)
type SpacesResponse struct {
	Query  map[string]int `json:"query"`
	Spaces []*Space       `json:"spaces"`
}

// Space represents a single Librato Space
type Space struct {
	// ID is the unique identifier of the Space
	ID int `json:"id"`
	// Name is the name of the Space
	Name string `json:"name"`
}

type SpacesCommunicator interface {
	List() ([]*Space, error)
}

type SpacesService struct {
	client *Client
}

// List implements the Librato Spaces API's List command
func (s *SpacesService) List() ([]*Space, error) {
	var spaces []*Space
	req, err := s.client.NewRequest("GET", "spaces", nil)

	if err != nil {
		return spaces, err
	}

	resp, err := s.client.Do(req)

	if err != nil {
		return spaces, err
	}

	var spacesResponse SpacesResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&spacesResponse); err != nil {
		return spaces, err
	}
	spaces = spacesResponse.Spaces
	return spaces, nil
}
