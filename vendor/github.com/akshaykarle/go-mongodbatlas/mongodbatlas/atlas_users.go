package mongodbatlas

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// AtlasUserService provides methods for accessing MongoDB Atlas AtlasUsers API endpoints.
// https://docs.atlas.mongodb.com/reference/api/user/
type AtlasUserService struct {
	sling *sling.Sling
}

// newAtlasUserService returns a new AtlasUserService.
func newAtlasUserService(sling *sling.Sling) *AtlasUserService {
	return &AtlasUserService{
		sling: sling.Path("users/"),
	}
}

// AtlasRole represents the permission on either the organization or group level
type AtlasRole struct {
	OrgID    string `json:"orgId,omitempty"`
	GroupID  string `json:"groupId,omitempty"`
	RoleName string `json:"roleName,omitempty"`
}

// AtlasUser represents users in your MongoDB Atlas UI.
type AtlasUser struct {
	EmailAddress string      `json:"emailAddress,omitempty"`
	ID           string      `json:"id,omitempty"`
	Username     string      `json:"username,omitempty"`
	FirstName    string      `json:"firstName,omitempty"`
	LastName     string      `json:"lastName,omitempty"`
	Password     string      `json:"password,omitempty"`
	MobileNumber string      `json:"mobileNumber,omitempty"`
	Country      string      `json:"country,omitempty"`
	Roles        []AtlasRole `json:"roles,omitempty"`
	TeamIDs      []string    `json:"teamIds,omitempty"`
}

// Get an atlasUser by ID.
// https://docs.atlas.mongodb.com/reference/api/user-get-by-id/
func (c *AtlasUserService) Get(id string) (*AtlasUser, *http.Response, error) {
	atlasUser := new(AtlasUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s", id)
	resp, err := c.sling.New().Get(path).Receive(atlasUser, apiError)
	return atlasUser, resp, relevantError(err, *apiError)
}

// GetByName gets an atlasUser by Name.
// https://docs.atlas.mongodb.com/reference/api/user-get-one-by-name/
func (c *AtlasUserService) GetByName(name string) (*AtlasUser, *http.Response, error) {
	atlasUser := new(AtlasUser)
	apiError := new(APIError)
	path := fmt.Sprintf("byName/%s", name)
	resp, err := c.sling.New().Get(path).Receive(atlasUser, apiError)
	return atlasUser, resp, relevantError(err, *apiError)
}

// Create an atlasUser
// https://docs.atlas.mongodb.com/reference/api/user-create/
func (c *AtlasUserService) Create(atlasUserParams *AtlasUser) (*AtlasUser, *http.Response, error) {
	atlasUser := new(AtlasUser)
	apiError := new(APIError)
	resp, err := c.sling.New().Post("").BodyJSON(atlasUserParams).Receive(atlasUser, apiError)
	return atlasUser, resp, relevantError(err, *apiError)
}

// Update an atlasUser
// https://docs.atlas.mongodb.com/reference/api/user-update/
func (c *AtlasUserService) Update(id string, atlasUserParams *AtlasUser) (*AtlasUser, *http.Response, error) {
	atlasUser := new(AtlasUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s", id)
	resp, err := c.sling.New().Patch(path).BodyJSON(atlasUserParams).Receive(atlasUser, apiError)
	return atlasUser, resp, relevantError(err, *apiError)
}
