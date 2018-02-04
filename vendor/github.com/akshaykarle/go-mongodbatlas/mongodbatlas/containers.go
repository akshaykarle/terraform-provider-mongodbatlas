package mongodbatlas

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// ContainerService provides methods for accessing MongoDB Atlas Containers API endpoints.
type ContainerService struct {
	sling *sling.Sling
}

// newContainerService returns a new ContainerService.
func newContainerService(sling *sling.Sling) *ContainerService {
	return &ContainerService{
		sling: sling.Path("groups/"),
	}
}

// Container represents a Cloud Services Containers in MongoDB.
type Container struct {
	ID             string `json:"id,omitempty"`
	ProviderName   string `json:"providerName,omitempty"`
	AtlasCidrBlock string `json:"atlasCidrBlock,omitempty"`
	RegionName     string `json:"regionName,omitempty"`
	VpcID          string `json:"vpcId,omitempty"`
	Provisioned    bool   `json:"provisioned,omitempty"`
}

// containerListResponse is the response from the ContainerService.List.
type containerListResponse struct {
	Results    []Container `json:"results"`
	TotalCount int         `json:"totalCount"`
}

// List all containers for the specified group.
// https://docs.atlas.mongodb.com/reference/api/vpc-get-containers-list/
func (c *ContainerService) List(gid string) ([]Container, *http.Response, error) {
	response := new(containerListResponse)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/containers", gid)
	resp, err := c.sling.New().Get(path).Receive(response, apiError)
	return response.Results, resp, relevantError(err, *apiError)
}

// Get a container in the specified group.
// https://docs.atlas.mongodb.com/reference/api/vpc-get-container/
func (c *ContainerService) Get(gid string, id string) (*Container, *http.Response, error) {
	container := new(Container)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/containers/%s", gid, id)
	resp, err := c.sling.New().Get(path).Receive(container, apiError)
	return container, resp, relevantError(err, *apiError)
}

// Create a container in the specified group.
// https://docs.atlas.mongodb.com/reference/api/vpc-create-container/
func (c *ContainerService) Create(gid string, containerParams *Container) (*Container, *http.Response, error) {
	container := new(Container)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/containers", gid)
	resp, err := c.sling.New().Post(path).BodyJSON(containerParams).Receive(container, apiError)
	return container, resp, relevantError(err, *apiError)
}

// Update a container in the specified group.
// https://docs.atlas.mongodb.com/reference/api/vpc-update-container/
func (c *ContainerService) Update(gid string, id string, containerParams *Container) (*Container, *http.Response, error) {
	container := new(Container)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/containers/%s", gid, id)
	resp, err := c.sling.New().Patch(path).BodyJSON(containerParams).Receive(container, apiError)
	return container, resp, relevantError(err, *apiError)
}

// Delete a container in the specified group.
func (c *ContainerService) Delete(gid string, id string) (*http.Response, error) {
	container := new(Container)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/containers/%s", gid, id)
	resp, err := c.sling.New().Delete(path).Receive(container, apiError)
	return resp, relevantError(err, *apiError)
}
