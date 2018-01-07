package mongodb

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// VPCService provides methods for accessing MongoDB Atlas VPC API endpoints.
type VPCService struct {
	sling *sling.Sling
}

// newVPCService returns a new VPCService.
func newVPCService(sling *sling.Sling) *VPCService {
	return &VPCService{
		sling: sling.Path("groups/"),
	}
}

// Peer represents a peering connection information in MongoDB.
type Peer struct {
	ID                  string `json:"id,omitempty"`
	ProviderName        string `json:"providerName,omitempty"`
	RouteTableCidrBlock string `json:"routeTableCidrBlock,omitempty"`
	VpcID               string `json:"vpcId,omitempty"`
	AwsAccountID        string `json:"awsAccountId,omitempty"`
	ConnectionID        string `json:"connectionI,omitempty"`
	StatusName          string `json:"statusName,omitempty"`
	ErrorStateName      string `json:"errorStateName,omitempty"`
	ContainerID         string `json:"containerId,omitempty"`
}

// peerListResponse is the response from the VPCService.List.
type peerListResponse struct {
	Results    []Peer `json:"results"`
	TotalCount int    `json:"totalCount"`
}

// List all peers for the specified group.
// https://docs.atlas.mongodb.com/reference/api/vpc-get-connections-list/
func (c *VPCService) List(gid string) ([]Peer, *http.Response, error) {
	response := new(peerListResponse)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/peers", gid)
	resp, err := c.sling.New().Get(path).Receive(response, apiError)
	return response.Results, resp, relevantError(err, *apiError)
}

// Get a peer in the specified group.
// https://docs.atlas.mongodb.com/reference/api/vpc-get-connection/
func (c *VPCService) Get(gid string, name string) (*Peer, *http.Response, error) {
	peer := new(Peer)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/peers/%s", gid, name)
	resp, err := c.sling.New().Get(path).Receive(peer, apiError)
	return peer, resp, relevantError(err, *apiError)
}

// Create a peer in the specified group.
// https://docs.atlas.mongodb.com/reference/api/vpc-create-peering-connection/
func (c *VPCService) Create(gid string, peer *Peer) (*Peer, *http.Response, error) {
	apiError := new(APIError)
	path := fmt.Sprintf("%s/peers", gid)
	resp, err := c.sling.New().Post(path).BodyJSON(peer).Receive(peer, apiError)
	return peer, resp, relevantError(err, *apiError)
}

// Update a peer in the specified group.
// https://docs.atlas.mongodb.com/reference/api/vpc-update-peering-connection/
func (c *VPCService) Update(gid string, name string, peer *Peer) (*Peer, *http.Response, error) {
	apiError := new(APIError)
	path := fmt.Sprintf("%s/peers/%s", gid, name)
	resp, err := c.sling.New().Patch(path).BodyJSON(peer).Receive(peer, apiError)
	return peer, resp, relevantError(err, *apiError)
}

// Delete a peer in the specified group.
// https://docs.atlas.mongodb.com/reference/api/vpc-delete-peering-connection/
func (c *VPCService) Delete(gid string, name string) (*http.Response, error) {
	peer := new(Peer)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/peers/%s", gid, name)
	resp, err := c.sling.New().Delete(path).Receive(peer, apiError)
	return resp, relevantError(err, *apiError)
}
