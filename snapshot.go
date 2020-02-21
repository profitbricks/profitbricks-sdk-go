package profitbricks

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

//Snapshot object
type Snapshot struct {
	ID         string             `json:"id,omitempty"`
	PBType     string             `json:"type,omitempty"`
	Href       string             `json:"href,omitempty"`
	Metadata   Metadata           `json:"metadata,omitempty"`
	Properties SnapshotProperties `json:"properties,omitempty"`
	Response   string             `json:"Response,omitempty"`
	Headers    *http.Header       `json:"headers,omitempty"`
	StatusCode int                `json:"statuscode,omitempty"`
}

// SnapshotProperties properties
type SnapshotProperties struct {
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	Location            string `json:"location,omitempty"`
	Size                int    `json:"size,omitempty"`
	CPUHotPlug          bool   `json:"cpuHotPlug,omitempty"`
	CPUHotUnplug        bool   `json:"cpuHotUnplug,omitempty"`
	RAMHotPlug          bool   `json:"ramHotPlug,omitempty"`
	RAMHotUnplug        bool   `json:"ramHotUnplug,omitempty"`
	NicHotPlug          bool   `json:"nicHotPlug,omitempty"`
	NicHotUnplug        bool   `json:"nicHotUnplug,omitempty"`
	DiscVirtioHotPlug   bool   `json:"discVirtioHotPlug,omitempty"`
	DiscVirtioHotUnplug bool   `json:"discVirtioHotUnplug,omitempty"`
	DiscScsiHotPlug     bool   `json:"discScsiHotPlug,omitempty"`
	DiscScsiHotUnplug   bool   `json:"discScsiHotUnplug,omitempty"`
	LicenceType         string `json:"licenceType,omitempty"`
}

//Snapshots object
type Snapshots struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []Snapshot   `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

//ListSnapshots lists all snapshots
func (c *Client) ListSnapshots() (*Snapshots, error) {
	url := snapshotsPath()
	ret := &Snapshots{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

//GetSnapshot gets a specific snapshot
func (c *Client) GetSnapshot(snapshotID string) (*Snapshot, error) {
	url := snapshotPath(snapshotID)
	ret := &Snapshot{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

// DeleteSnapshot deletes a specified snapshot
func (c *Client) DeleteSnapshot(snapshotID string) (*http.Header, error) {
	url := snapshotPath(snapshotID)
	ret := &http.Header{}
	err := c.Delete(url, ret, http.StatusAccepted)
	return ret, err
}

// UpdateSnapshot updates a snapshot
func (c *Client) UpdateSnapshot(snapshotID string, request SnapshotProperties) (*Snapshot, error) {
	url := snapshotPath(snapshotID)
	ret := &Snapshot{}
	err := c.Patch(url, request, ret, http.StatusAccepted)
	return ret, err
}

// DeleteSnapshotAndWait deletes a specified snapshot and waits for the request
// to complete. The default timeout is 10 minutes.
func (c *Client) DeleteSnapshotAndWait(snapshotID string, timeout time.Duration) error {
	ret, err := c.DeleteSnapshot(snapshotID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), DurationOrDefault(timeout, 10*time.Minute))
	defer cancel()
	return c.WaitTillProvisionedOrCanceled(ctx, ret.Get("location"))
}

// ListSnapshotsWithSelector retrieves all snapshots and performs client-side
// filtering according to the list of selectors. Each selector is concatenated
// with logical AND.
func (c *Client) ListSnapshotsWithSelector(selectors ...SnapshotSelector) ([]Snapshot, error) {
	url := snapshotsPath()
	ret := &Snapshots{}
	err := c.Get(url, ret, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var result []Snapshot
outerLoop:
	for _, snapshot := range ret.Items {
		for _, selector := range selectors {
			if !selector(&snapshot) {
				continue outerLoop
			}
		}
		result = append(result, snapshot)
	}
	return result, nil
}

// SnapshotSelector is used to do client-side filtering of a list of Snapshots
type SnapshotSelector func(*Snapshot) bool

// SnapshotByState selects snapshots with the given state
func SnapshotByState(state string) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		return snapshot.Metadata.State == state
	}
}

// SnapshotByName selects snapshots with the given name
func SnapshotByName(name string) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		return snapshot.Properties.Name == name
	}
}

// SnapshotByDescription selects snapshots with the given description
func SnapshotByDescription(description string) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		return snapshot.Properties.Description == description
	}
}

// SelectExactSnapshot concatenates the provided selectors with logical AND.
func SelectExactSnapshot(selectors ...SnapshotSelector) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		for _, selector := range selectors {
			if !selector(snapshot) {
				return false
			}
		}
		return true
	}
}

// SelectAnySnapshot concatenates the provided selectors with logical OR.
func SelectAnySnapshot(selectors ...SnapshotSelector) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		for _, selector := range selectors {
			if selector(snapshot) {
				return true
			}
		}
		return false
	}
}

// HasDeleteSnapshotInProgress tries to find an active delete snapshot request (QUEUED or RUNNING).
func (c *Client) HasDeleteSnapshotInProgress(snapshotID string) (bool, error) {
	f := NewRequestListFilter().WithUrl(snapshotPath(snapshotID)).WithMethod(http.MethodDelete)
	result, err := c.ListRequestsWithFilter(f.Clone().WithRequestStatus(RequestStatusQueued))
	if err != nil {
		return false, err
	}
	if len(result.Items) > 0 {
		return true, nil
	}
	result, err = c.ListRequestsWithFilter(f.Clone().WithRequestStatus(RequestStatusRunning))
	if err != nil {
		return false, err
	}
	return len(result.Items) > 0, nil
}

// IsSnapshotInUse checks whether a snapshot is being used by either an existing volume or a volume
// that is in the process being created (has an active request in the queue).
func (c *Client) IsSnapshotInUse(dcId, volumeId, snapshotId string) (bool, error) {
	volumes, err := c.ListVolumes(dcId)
	if err != nil {
		return false, nil
	}
	for _, volume := range volumes.Items {
		if volume.Properties.Image == snapshotId {
			return true, nil
		}
	}

	var requests []Request
	f := NewRequestListFilter().WithUrl(volumesPath(dcId)).WithMethod(http.MethodPost)
	queued, err := c.ListRequestsWithFilter(f.Clone().WithRequestStatus(RequestStatusQueued))
	if err != nil {
		return false, err
	}
	requests = append(requests, queued.Items...)
	running, err := c.ListRequestsWithFilter(f.Clone().WithRequestStatus(RequestStatusRunning))
	if err != nil {
		return false, err
	}
	requests = append(requests, running.Items...)

	for _, request := range requests {
		// Assume that the request properties contain a valid volume
		var volume Volume
		_ = json.Unmarshal([]byte(request.Properties.Body), &volume)
		if volume.Properties.Image == snapshotId {
			return true, nil
		}
	}

	return false, nil
}
