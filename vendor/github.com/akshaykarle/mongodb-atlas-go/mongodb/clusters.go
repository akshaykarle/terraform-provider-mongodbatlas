package mongodb

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// ClusterService provides methods for accessing MongoDB Atlas Clusters API endpoints.
type ClusterService struct {
	sling *sling.Sling
}

// newClusterService returns a new ClusterService.
func newClusterService(sling *sling.Sling) *ClusterService {
	return &ClusterService{
		sling: sling.Path("groups/"),
	}
}

// ProviderSettings is the configuration for the provisioned servers on which MongoDB runs.
// The available options are specific to the cloud service provider.
type ProviderSettings struct {
	ProviderName        string `json:"providerName,omitempty"`
	BackingProviderName string `json:"backingProviderName,omitempty"`
	RegionName          string `json:"regionName,omitempty"`
	InstanceSizeName    string `json:"instanceSizeName,omitempty"`
}

// Cluster represents a Cluster configuration in MongoDB.
type Cluster struct {
	ID                  string           `json:"id,omitempty"`
	GroupID             string           `json:"groupId,omitempty"`
	Name                string           `json:"name,omitempty"`
	MongoDBVersion      string           `json:"mongoDBVersion,omitempty"`
	DiskSizeGB          float64          `json:"diskSizeGB,omitempty"`
	MongoDBMajorVersion string           `json:"mongoDBMajorVersion,omitempty"`
	BackupEnabled       bool             `json:"backupEnabled"`
	StateName           string           `json:"stateName,omitempty"`
	ReplicationFactor   int              `json:"replicationFactor,omitempty"`
	ProviderSettings    ProviderSettings `json:"providerSettings,omitempty"`
}

// clusterListResponse is the response from the ClusterService.List.
type clusterListResponse struct {
	Results    []Cluster `json:"results"`
	TotalCount int       `json:"totalCount"`
}

// List all clusters for the specified group.
// https://docs.atlas.mongodb.com/reference/api/clusters-get-all/
func (c *ClusterService) List(gid string) ([]Cluster, *http.Response, error) {
	response := new(clusterListResponse)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/clusters", gid)
	resp, err := c.sling.New().Get(path).Receive(response, apiError)
	return response.Results, resp, relevantError(err, *apiError)
}

// Get a cluster in the specified group.
// https://docs.atlas.mongodb.com/reference/api/clusters-get-one/
func (c *ClusterService) Get(gid string, name string) (*Cluster, *http.Response, error) {
	cluster := new(Cluster)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/clusters/%s", gid, name)
	resp, err := c.sling.New().Get(path).Receive(cluster, apiError)
	return cluster, resp, relevantError(err, *apiError)
}

// Create a cluster in the specified group.
// https://docs.atlas.mongodb.com/reference/api/clusters-create-one/
func (c *ClusterService) Create(gid string, cluster *Cluster) (*Cluster, *http.Response, error) {
	apiError := new(APIError)
	path := fmt.Sprintf("%s/clusters", gid)
	resp, err := c.sling.New().Post(path).BodyJSON(cluster).Receive(cluster, apiError)
	return cluster, resp, relevantError(err, *apiError)
}

// Update a cluster in the specified group.
// https://docs.atlas.mongodb.com/reference/api/clusters-modify-one/
func (c *ClusterService) Update(gid string, name string, cluster *Cluster) (*Cluster, *http.Response, error) {
	apiError := new(APIError)
	path := fmt.Sprintf("%s/clusters/%s", gid, name)
	resp, err := c.sling.New().Patch(path).BodyJSON(cluster).Receive(cluster, apiError)
	return cluster, resp, relevantError(err, *apiError)
}

// Delete a cluster in the specified group.
// https://docs.atlas.mongodb.com/reference/api/clusters-delete-one/
func (c *ClusterService) Delete(gid string, name string) (*http.Response, error) {
	cluster := new(Cluster)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/clusters/%s", gid, name)
	resp, err := c.sling.New().Delete(path).Receive(cluster, apiError)
	return resp, relevantError(err, *apiError)
}
