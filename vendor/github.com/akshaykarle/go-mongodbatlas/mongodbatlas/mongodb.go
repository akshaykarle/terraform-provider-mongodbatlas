package mongodbatlas

import (
	"net/http"

	"github.com/dghubble/sling"
)

const apiURL = "https://cloud.mongodb.com/api/atlas/v1.0/"

// Client is a MongoDB Atlas client for making MongoDB API requests.
type Client struct {
	sling         *sling.Sling
	Root          *RootService
	Whitelist     *WhitelistService
	Projects      *ProjectService
	Clusters      *ClusterService
	Containers    *ContainerService
	Peers         *PeerService
	DatabaseUsers *DatabaseUserService
	Organizations *OrganizationService
}

// NewClient returns a new Client.
func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(apiURL)
	return &Client{
		sling:         base,
		Root:          newRootService(base.New()),
		Whitelist:     newWhitelistService(base.New()),
		Projects:      newProjectService(base.New()),
		Clusters:      newClusterService(base.New()),
		Containers:    newContainerService(base.New()),
		Peers:         newPeerService(base.New()),
		DatabaseUsers: newDatabaseUserService(base.New()),
		Organizations: newOrganizationService(base.New()),
	}
}
