package mongodbatlas

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// DatabaseUserService provides methods for accessing MongoDB Atlas DatabaseUsers API endpoints.
type DatabaseUserService struct {
	sling *sling.Sling
}

// newDatabaseUserService returns a new DatabaseUserService.
func newDatabaseUserService(sling *sling.Sling) *DatabaseUserService {
	return &DatabaseUserService{
		sling: sling.Path("groups/"),
	}
}

// Role allows the user to perform particular actions on the specified database.
// A role on the admin database can include privileges that apply to the other databases as well.
type Role struct {
	DatabaseName   string `json:"databaseName,omitempty"`
	CollectionName string `json:"collectionName,omitempty"`
	RoleName       string `json:"roleName,omitempty"`
}

// DatabaseUser represents MongoDB users in your cluster.
type DatabaseUser struct {
	GroupID         string `json:"groupId,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	DatabaseName    string `json:"databaseName,omitempty"`
	DeleteAfterDate string `json:"deleteAfterDate,omitempty"`
	Roles           []Role `json:"roles,omitempty"`
}

// databaseUserListResponse is the response from the DatabaseUserService.List.
type databaseUserListResponse struct {
	Results    []DatabaseUser `json:"results"`
	TotalCount int            `json:"totalCount"`
}

// List all databaseUsers for the specified group.
// https://docs.atlas.mongodb.com/reference/api/database-users-get-all-users/
func (c *DatabaseUserService) List(gid string) ([]DatabaseUser, *http.Response, error) {
	response := new(databaseUserListResponse)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/databaseUsers", gid)
	resp, err := c.sling.New().Get(path).Receive(response, apiError)
	return response.Results, resp, relevantError(err, *apiError)
}

// Get a databaseUser in the specified group.
// https://docs.atlas.mongodb.com/reference/api/database-users-get-single-user/
func (c *DatabaseUserService) Get(gid string, username string) (*DatabaseUser, *http.Response, error) {
	databaseUser := new(DatabaseUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/databaseUsers/admin/%s", gid, username)
	resp, err := c.sling.New().Get(path).Receive(databaseUser, apiError)
	return databaseUser, resp, relevantError(err, *apiError)
}

// Create a databaseUser in the specified group.
// https://docs.atlas.mongodb.com/reference/api/databaseUsers-create-one/
func (c *DatabaseUserService) Create(gid string, databaseUserParams *DatabaseUser) (*DatabaseUser, *http.Response, error) {
	databaseUser := new(DatabaseUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/databaseUsers", gid)
	resp, err := c.sling.New().Post(path).BodyJSON(databaseUserParams).Receive(databaseUser, apiError)
	return databaseUser, resp, relevantError(err, *apiError)
}

// Update a databaseUser in the specified group.
// https://docs.atlas.mongodb.com/reference/api/databaseUsers-modify-one/
func (c *DatabaseUserService) Update(gid string, username string, databaseUserParams *DatabaseUser) (*DatabaseUser, *http.Response, error) {
	databaseUser := new(DatabaseUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/databaseUsers/admin/%s", gid, username)
	resp, err := c.sling.New().Patch(path).BodyJSON(databaseUserParams).Receive(databaseUser, apiError)
	return databaseUser, resp, relevantError(err, *apiError)
}

// Delete a databaseUser in the specified group.
// https://docs.atlas.mongodb.com/reference/api/databaseUsers-delete-one/
func (c *DatabaseUserService) Delete(gid string, username string) (*http.Response, error) {
	databaseUser := new(DatabaseUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/databaseUsers/admin/%s", gid, username)
	resp, err := c.sling.New().Delete(path).Receive(databaseUser, apiError)
	return resp, relevantError(err, *apiError)
}
