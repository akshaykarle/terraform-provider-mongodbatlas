package mongodbatlas

import (
	"net/http"

	"github.com/dghubble/sling"
)

// RootService checks connectivity to MongoDB Atlas API.
type RootService struct {
	sling *sling.Sling
}

// newRootService returns a new RootService.
func newRootService(sling *sling.Sling) *RootService {
	return &RootService{
		sling: sling,
	}
}

// Root is the response from the RootService.List.
type Root struct {
	AppName string `json:"appName"`
	Build   string `json:"build"`
}

// Get the root resource which is the starting point for the Atlas API.
// https://docs.atlas.mongodb.com/reference/api/root/
func (c *RootService) Get() (*Root, *http.Response, error) {
	response := new(Root)
	apiError := new(APIError)
	resp, err := c.sling.Get("").Receive(response, apiError)
	return response, resp, relevantError(err, *apiError)
}
