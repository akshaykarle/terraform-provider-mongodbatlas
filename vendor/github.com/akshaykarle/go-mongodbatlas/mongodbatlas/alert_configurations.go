package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// AlertConfigurationService provides methods for accessing MongoDB Atlas Alert Configurations API endpoints.
type AlertConfigurationService struct {
	sling *sling.Sling
}

// newAlertConfigurationService returns a new AlertConfigurationService.
func newAlertConfigurationService(sling *sling.Sling) *AlertConfigurationService {
	return &AlertConfigurationService{
		sling: sling.Path("groups/"),
	}
}

// Notification is a way to get notified when a metric crosses the threshold
type Notification struct {
	TypeName            string `json:"typeName,omitempty"`
	IntervalMin         int    `json:"intervalMin,omitempty"`
	DelayMin            int    `json:"delayMin,omitempty"`
	EmailEnabled        bool   `json:"emailEnabled,omitempty"`
	SMSEnabled          bool   `json:"smsEnabled,omitempty"`
	Username            string `json:"username,omitempty"`
	TeamID              string `json:"teamId,omitempty"`
	EmailAddress        string `json:"emailAddress,omitempty"`
	MobileNumber        string `json:"mobileNumber,omitempty"`
	NotificationToken   string `json:"notificationToken,omitempty"`
	RoomName            string `json:"roomName,omitempty"`
	ChannelName         string `json:"channelName,omitempty"`
	APIToken            string `json:"apiToken,omitempty"`
	OrgName             string `json:"orgName,omitempty"`
	FlowName            string `json:"flowName,omitempty"`
	FlowdockAPIToken    string `json:"flowdockApiToken,omitempty"`
	ServiceKey          string `json:"serviceKey,omitempty"`
	VictorOpsAPIKey     string `json:"victorOpsApiKey,omitempty"`
	VictorOpsRoutingKey string `json:"victorOpsRoutingKey,omitempty"`
	OpsGenieAPIKey      string `json:"opsGenieApiKey,omitempty"`
}

// MetricThreshold describes how to know when to trigger this alert
type MetricThreshold struct {
	MetricName string  `json:"metricName,omitempty"`
	Operator   string  `json:"operator,omitempty"`
	Threshold  float64 `json:"threshold,omitempty"`
	Units      string  `json:"units,omitempty"`
	Mode       string  `json:"mode,omitempty"`
}

// Matcher contains the metric(s) we'd like to alert on
type Matcher struct {
	FieldName string `json:"fieldName,omitempty"`
	Operator  string `json:"operator,omitempty"`
	Value     string `json:"value,omitempty"`
}

// AlertConfiguration represents an AlertConfiguration in MongoDB.
type AlertConfiguration struct {
	ID              string          `json:"id,omitempty"`
	GroupID         string          `json:"groupId,omitempty"`
	EventTypeName   string          `json:"eventTypeName,omitempty"`
	Enabled         bool            `json:"enabled,omitempty"`
	Notifications   []Notification  `json:"notifications,omitempty"`
	MetricThreshold MetricThreshold `json:"metricThreshold,omitempty"`
	Matchers        []Matcher       `json:"matchers,omitempty"`
}

// MarshalJSON is custom defined here because the API pukes if you specify an empty metricThreshold ("metricThreshold":{}) when it doesn't want one
func (r AlertConfiguration) MarshalJSON() ([]byte, error) {
	encoded := struct {
		ID              string           `json:"id,omitempty"`
		GroupID         string           `json:"groupId,omitempty"`
		EventTypeName   string           `json:"eventTypeName,omitempty"`
		Enabled         bool             `json:"enabled,omitempty"`
		Notifications   []Notification   `json:"notifications,omitempty"`
		MetricThreshold *MetricThreshold `json:"metricThreshold,omitempty"`
		Matchers        []Matcher        `json:"matchers,omitempty"`
	}{
		ID:            r.ID,
		GroupID:       r.GroupID,
		EventTypeName: r.EventTypeName,
		Enabled:       r.Enabled,
		Notifications: r.Notifications,
		Matchers:      r.Matchers,
	}
	// only add the metric threshold if it's not empty
	if (MetricThreshold{}) != r.MetricThreshold {
		encoded.MetricThreshold = &r.MetricThreshold
	}

	return json.Marshal(encoded)
}

// alertConfigsListResponse is the response from the AlertConfigurationService.List.
type alertConfigsListResponse struct {
	Results    []AlertConfiguration `json:"results"`
	TotalCount int                  `json:"totalCount"`
}

// List all alert configurations for the specified group.
// https://docs.atlas.mongodb.com/reference/api/alert-configurations-get-all-configs/
func (c *AlertConfigurationService) List(gid string) ([]AlertConfiguration, *http.Response, error) {
	response := new(alertConfigsListResponse)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/alertConfigs", gid)
	resp, err := c.sling.New().Get(path).Receive(response, apiError)
	return response.Results, resp, relevantError(err, *apiError)
}

// Get an alert configuration in the specified group.
// https://docs.atlas.mongodb.com/reference/api/alert-configurations-get-config/
func (c *AlertConfigurationService) Get(gid string, id string) (*AlertConfiguration, *http.Response, error) {
	alert := new(AlertConfiguration)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/alertConfigs/%s", gid, id)
	resp, err := c.sling.New().Get(path).Receive(alert, apiError)
	return alert, resp, relevantError(err, *apiError)
}

// Create an alert configuration in the specified group.
// https://docs.atlas.mongodb.com/reference/api/alert-configurations-create-config/
func (c *AlertConfigurationService) Create(gid string, alertConfigurationParams *AlertConfiguration) (*AlertConfiguration, *http.Response, error) {
	alert := new(AlertConfiguration)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/alertConfigs", gid)
	resp, err := c.sling.New().Post(path).BodyJSON(alertConfigurationParams).Receive(alert, apiError)
	return alert, resp, relevantError(err, *apiError)
}

// Update an alert configuration in the specified group.
// https://docs.atlas.mongodb.com/reference/api/alert-configurations-update-config/
func (c *AlertConfigurationService) Update(gid string, id string, alertConfigurationParams *AlertConfiguration) (*AlertConfiguration, *http.Response, error) {
	alert := new(AlertConfiguration)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/alertConfigs/%s", gid, id)
	resp, err := c.sling.New().Put(path).BodyJSON(alertConfigurationParams).Receive(alert, apiError)
	return alert, resp, relevantError(err, *apiError)
}

// Delete an alert configuration in the specified group.
// https://docs.atlas.mongodb.com/reference/api/alert-configurations-update-config/
func (c *AlertConfigurationService) Delete(gid string, id string) (*http.Response, error) {
	alert := new(AlertConfiguration)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/alertConfigs/%s", gid, id)
	resp, err := c.sling.New().Delete(path).Receive(alert, apiError)
	return resp, relevantError(err, *apiError)
}
