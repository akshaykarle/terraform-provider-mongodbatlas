package mongodbatlas

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// ProjectService provides methods for accessing MongoDB Atlas Projects API endpoints.
type ProjectService struct {
	sling *sling.Sling
}

// newProjectService returns a new ProjectService.
func newProjectService(sling *sling.Sling) *ProjectService {
	return &ProjectService{
		sling: sling.Path("groups/"),
	}
}

// Project represents a projecting connection information in MongoDB.
type Project struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	OrgID        string `json:"orgId,omitempty"`
	Created      string `json:"created,omitempty"`
	ClusterCount int    `json:"clusterCount,omitempty"`
}

// projectListResponse is the response from the ProjectService.List.
type projectListResponse struct {
	Results    []Project `json:"results"`
	TotalCount int       `json:"totalCount"`
}

// List all projects the authenticated user has access to.
// https://docs.atlas.mongodb.com/reference/api/project-get-all/
func (c *ProjectService) List() ([]Project, *http.Response, error) {
	response := new(projectListResponse)
	apiError := new(APIError)
	resp, err := c.sling.New().Get("").Receive(response, apiError)
	return response.Results, resp, relevantError(err, *apiError)
}

// Get information about the project associated to group id
// https://docs.atlas.mongodb.com/reference/api/project-get-one/
func (c *ProjectService) Get(id string) (*Project, *http.Response, error) {
	project := new(Project)
	apiError := new(APIError)
	path := fmt.Sprintf("%s", id)
	resp, err := c.sling.New().Get(path).Receive(project, apiError)
	return project, resp, relevantError(err, *apiError)
}

// GetByName information about the project associated to group name
// https://docs.atlas.mongodb.com/reference/api/project-get-one-by-name/
func (c *ProjectService) GetByName(name string) (*Project, *http.Response, error) {
	project := new(Project)
	apiError := new(APIError)
	path := fmt.Sprintf("byName/%s", name)
	resp, err := c.sling.New().Get(path).Receive(project, apiError)
	return project, resp, relevantError(err, *apiError)
}

// Create a project.
// https://docs.atlas.mongodb.com/reference/api/project-create-one/
func (c *ProjectService) Create(projectParams *Project) (*Project, *http.Response, error) {
	project := new(Project)
	apiError := new(APIError)
	resp, err := c.sling.New().Post("").BodyJSON(projectParams).Receive(project, apiError)
	return project, resp, relevantError(err, *apiError)
}

// Delete a project.
// https://docs.atlas.mongodb.com/reference/api/project-delete-one/
func (c *ProjectService) Delete(id string) (*http.Response, error) {
	project := new(Project)
	apiError := new(APIError)
	path := fmt.Sprintf("%s", id)
	resp, err := c.sling.New().Delete(path).Receive(project, apiError)
	return resp, relevantError(err, *apiError)
}
