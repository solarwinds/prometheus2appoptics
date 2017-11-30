package appoptics

import "net/http"

// SpacesResponse represents the returned data payload from Spaces API's List command (/spaces)
type SpacesResponse struct {
	Query  map[string]int `json:"query"`
	Spaces []*Space       `json:"spaces"`
}

// Space represents a single AppOptics Space
type Space struct {
	// ID is the unique identifier of the Space
	ID int `json:"id"`
	// Name is the name of the Space
	Name string `json:"name"`
}

type SpacesCommunicator interface {
	List() ([]*Space, *http.Response, error)
}

type SpacesService struct {
	client *Client
}

// List implements the  Spaces API's List command
func (s *SpacesService) List() ([]*Space, *http.Response, error) {
	var spaces []*Space
	req, err := s.client.NewRequest("GET", "spaces", nil)

	if err != nil {
		return spaces, nil, err
	}

	var spacesResponse SpacesResponse
	resp, err := s.client.Do(req, &spacesResponse)

	if err != nil {
		return spaces, resp, err
	}

	spaces = spacesResponse.Spaces
	return spaces, resp, nil
}
