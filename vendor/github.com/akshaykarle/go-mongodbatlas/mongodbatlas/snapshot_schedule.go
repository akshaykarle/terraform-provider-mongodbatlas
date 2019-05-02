package mongodbatlas

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// SnapshotScheduleService provides methods for accessing MongoDB Atlas Snapshot Schedule API endpoints.
type SnapshotScheduleService struct {
	sling *sling.Sling
}

// newSnapshotScheduleService returns a new SnapshotScheduleService.
func newSnapshotScheduleService(sling *sling.Sling) *SnapshotScheduleService {
	return &SnapshotScheduleService{
		sling: sling.Path("groups/"),
	}
}

// SnapshotSchedule represents a snapshot schedule's connection information in MongoDB.
type SnapshotSchedule struct {
	GroupID                        string  `json:"groupId,omitempty"`
	ClusterID                      string  `json:"clusterId,omitempty"`
	SnapshotIntervalHours          float64 `json:"snapshotIntervalHours,omitempty"`
	SnapshotRetentionDays          float64 `json:"snapshotRetentionDays,omitempty"`
	DailySnapshotRetentionDays     float64 `json:"dailySnapshotRetentionDays,omitempty"`
	PointInTimeWindowHours         float64 `json:"pointInTimeWindowHours,omitempty"`
	WeeklySnapshotRetentionWeeks   float64 `json:"weeklySnapshotRetentionWeeks,omitempty"`
	MonthlySnapshotRetentionMonths float64 `json:"monthlySnapshotRetentionMonths,omitempty"`
	ClusterCheckpintIntervalMin    float64 `json:"clusterCheckpintIntervalMin,omitempty"`
}

// Get the snapshot schedule for the specified cluster
// https://docs.atlas.mongodb.com/reference/api/snapshot-schedule-get/
func (c *SnapshotScheduleService) Get(gid string, clusterName string) (*SnapshotSchedule, *http.Response, error) {
	snapshotSchedule := new(SnapshotSchedule)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/clusters/%s/snapshotSchedule", gid, clusterName)
	resp, err := c.sling.New().Get(path).Receive(snapshotSchedule, apiError)
	return snapshotSchedule, resp, relevantError(err, *apiError)
}

// Update the snapshot schedule for the specified cluster.
// https://docs.atlas.mongodb.com/reference/api/snapshot-schedule/
func (c *SnapshotScheduleService) Update(gid string, clusterName string, snapshotScheduleParams *SnapshotSchedule) (*SnapshotSchedule, *http.Response, error) {
	snapshotSchedule := new(SnapshotSchedule)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/clusters/%s/snapshotSchedule", gid, clusterName)
	resp, err := c.sling.New().Patch(path).BodyJSON(snapshotScheduleParams).Receive(snapshotSchedule, apiError)
	return snapshotSchedule, resp, relevantError(err, *apiError)
}
