package mongodbatlas

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// OrganizationService provides methods for accessing MongoDB Atlas Organizations API endpoints.
type OrganizationService struct {
	sling *sling.Sling
}

// newOrganizationService returns a new OrganizationService.
func newOrganizationService(sling *sling.Sling) *OrganizationService {
	return &OrganizationService{
		sling: sling.Path("orgs/"),
	}
}

// Organization represents an organization's connection information in MongoDB.
type Organization struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// organizationListResponse is the response from the OrganizationService.List.
type organizationListResponse struct {
	Results    []Organization `json:"results"`
	TotalCount int            `json:"totalCount"`
}

// List all organizations the authenticated user has access to.
// https://docs.atlas.mongodb.com/reference/api/organization-get-all/
func (c *OrganizationService) List() ([]Organization, *http.Response, error) {
	response := new(organizationListResponse)
	apiError := new(APIError)
	resp, err := c.sling.New().Get("").Receive(response, apiError)
	return response.Results, resp, relevantError(err, *apiError)
}

// Get information about the organization associated to org ID
// https://docs.atlas.mongodb.com/reference/api/organization-get-one/
func (c *OrganizationService) Get(id string) (*Organization, *http.Response, error) {
	organization := new(Organization)
	apiError := new(APIError)
	path := fmt.Sprintf("%s", id)
	resp, err := c.sling.New().Get(path).Receive(organization, apiError)
	return organization, resp, relevantError(err, *apiError)
}

// Create an organization.
// https://docs.atlas.mongodb.com/reference/api/organization-create-one/
func (c *OrganizationService) Create(organizationParams *Organization) (*Organization, *http.Response, error) {
	organization := new(Organization)
	apiError := new(APIError)
	resp, err := c.sling.New().Post("").BodyJSON(organizationParams).Receive(organization, apiError)
	return organization, resp, relevantError(err, *apiError)
}

// Update name of an organization.
// https://docs.atlas.mongodb.com/reference/api/organization-rename/
func (c *OrganizationService) Update(id string, organizationParams *Organization) (*Organization, *http.Response, error) {
	organization := new(Organization)
	apiError := new(APIError)
	path := fmt.Sprintf("%s", id)
	resp, err := c.sling.New().Patch(path).BodyJSON(organizationParams).Receive(organization, apiError)
	return organization, resp, relevantError(err, *apiError)
}

// Delete an organization
// https://docs.atlas.mongodb.com/reference/api/organization-delete-one/
func (c *OrganizationService) Delete(id string) (*http.Response, error) {
	apiError := new(APIError)
	path := fmt.Sprintf("%s", id)
	resp, err := c.sling.New().Delete(path).Receive(nil, apiError)
	return resp, relevantError(err, *apiError)
}
