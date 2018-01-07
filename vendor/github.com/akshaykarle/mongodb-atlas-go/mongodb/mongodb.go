package mongodb

import (
	"net/http"

	"github.com/dghubble/sling"
)

const apiURL = "https://cloud.mongodb.com/api/atlas/v1.0/"

// Client is a MongoDB Atlas client for making MongoDB API requests.
type Client struct {
	sling    *sling.Sling
	Clusters *ClusterService
	VPC      *VPCService
}

// NewClient returns a new Client.
func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(apiURL)
	return &Client{
		sling:    base,
		Clusters: newClusterService(base.New()),
		VPC:      newVPCService(base.New()),
	}
}
