package mongodbatlas

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// PrivateIPModeService provides many needfuls
type PrivateIPModeService struct {
	sling *sling.Sling
}

// newPrivateIPModeService returns a new instance of PrivateIPModeService
func newPrivateIPModeService(sling *sling.Sling) *PrivateIPModeService {
	return &PrivateIPModeService{
		sling: sling.Path("groups/"),
	}
}

// PrivateIPMode struct is the response from both the Enable and Disable functions
type PrivateIPMode struct {
	Enabled bool `json:"enabled,omitempty"`
}

// Enable – Enables PrivateIPMode on the Container
// https://docs.atlas.mongodb.com/reference/api/set-private-ip-mode-for-project/
func (p *PrivateIPModeService) Enable(gid string) (*http.Response, error) {
	privateIPMode := new(PrivateIPMode)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/privateIpMode", gid)
	params := PrivateIPMode{
		Enabled: true,
	}
	resp, err := p.sling.New().Patch(path).BodyJSON(params).Receive(privateIPMode, apiError)
	return resp, relevantError(err, *apiError)
}

// Disable – Disables the PrivateIPMode on the Container
// https://docs.atlas.mongodb.com/reference/api/set-private-ip-mode-for-project/
func (p *PrivateIPModeService) Disable(gid string) (*http.Response, error) {
	privateIPMode := new(PrivateIPMode)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/privateIpMode", gid)
	params := PrivateIPMode{
		Enabled: false,
	}
	resp, err := p.sling.New().Patch(path).BodyJSON(params).Receive(privateIPMode, apiError)
	return resp, relevantError(err, *apiError)
}
