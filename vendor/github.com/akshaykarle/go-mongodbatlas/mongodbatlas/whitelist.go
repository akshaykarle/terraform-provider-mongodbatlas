package mongodbatlas

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
)

// WhitelistService provides methods for accessing MongoDB Atlas's Group IP Whitelist API endpoints.
type WhitelistService struct {
	sling *sling.Sling
}

// newWhitelistService returns a new WhitelistService.
func newWhitelistService(sling *sling.Sling) *WhitelistService {
	return &WhitelistService{
		sling: sling.Path("groups/"),
	}
}

// Whitelist represents a IP whitelist, which controls client access to your group’s MongoDB clusters.
// Clients can connect to clusters only from IP addresses on the whitelist.
type Whitelist struct {
	CidrBlock string `json:"cidrBlock,omitempty"`
	Comment   string `json:"comment,omitempty"`
	GroupID   string `json:"groupId,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
}

// whitelistListResponse is the response from the WhitelistService.List.
type whitelistListResponse struct {
	Results    []Whitelist `json:"results"`
	TotalCount int         `json:"totalCount"`
}

// List a Group’s IP Whitelist.
// https://docs.atlas.mongodb.com/reference/api/whitelist/#get-a-group-s-ip-whitelist
func (c *WhitelistService) List(gid string) ([]Whitelist, *http.Response, error) {
	response := new(whitelistListResponse)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/whitelist", gid)
	resp, err := c.sling.New().Get(path).Receive(response, apiError)
	return response.Results, resp, relevantError(err, *apiError)
}

// Get the Entry for a Specific Address in a Group’s IP Whitelist
// https://docs.atlas.mongodb.com/reference/api/whitelist/#get-the-entry-for-a-specific-address-in-a-group-s-ip-whitelist
func (c *WhitelistService) Get(gid string, ip string) (*Whitelist, *http.Response, error) {
	whitelist := new(Whitelist)
	apiError := new(APIError)
	escapedIP := url.PathEscape(ip)
	path := fmt.Sprintf("%s/whitelist/%s", gid, escapedIP)
	resp, err := c.sling.New().Get(path).Receive(whitelist, apiError)
	return whitelist, resp, relevantError(err, *apiError)
}

// Create entries in a Group's whitelist.
// https://docs.atlas.mongodb.com/reference/api/whitelist/#add-entries-to-a-group-s-ip-whitelist
func (c *WhitelistService) Create(gid string, whitelistParams []Whitelist) ([]Whitelist, *http.Response, error) {
	response := new(whitelistListResponse)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/whitelist", gid)
	resp, err := c.sling.New().Post(path).BodyJSON(whitelistParams).Receive(response, apiError)
	return response.Results, resp, relevantError(err, *apiError)
}

// Delete an Entry from Group's IP Whitelist.
// https://docs.atlas.mongodb.com/reference/api/whitelist/#delete-an-entry-from-a-group-s-ip-whitelist
func (c *WhitelistService) Delete(gid string, ip string) (*http.Response, error) {
	whitelist := new(Whitelist)
	apiError := new(APIError)
	escapedIP := url.PathEscape(ip)
	path := fmt.Sprintf("%s/whitelist/%s", gid, escapedIP)
	resp, err := c.sling.New().Delete(path).Receive(whitelist, apiError)
	return resp, relevantError(err, *apiError)
}
