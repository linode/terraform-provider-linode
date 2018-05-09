package golinode

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

// InstanceSnapshot represents a linode backup snapshot
type InstanceSnapshot struct {
	CreatedStr  string `json:"created"`
	UpdatedStr  string `json:"updated"`
	FinishedStr string `json:"finished"`

	ID       int
	Label    string
	Status   string
	Type     string
	Created  *time.Time `json:"-"`
	Updated  *time.Time `json:"-"`
	Finished *time.Time `json:"-"`
	Configs  []string
	Disks    []*InstanceSnapshotDisk
}

type InstanceSnapshotDisk struct {
	Label      string
	Size       int
	Filesystem string
}

func (l *InstanceSnapshot) fixDates() *InstanceSnapshot {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	l.Finished, _ = parseDates(l.FinishedStr)
	return l
}

// InstanceSnapshotsPagedResponse represents a paginated InstanceSnapshot API response
type InstanceSnapshotsPagedResponse struct {
	*PageOptions
	Data []*InstanceSnapshot
}

// Endpoint gets the endpoint URL for InstanceSnapshot
func (InstanceSnapshotsPagedResponse) EndpointWithID(c *Client, id int) string {
	endpoint, err := c.InstanceSnapshots.EndpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends InstanceSnapshots when processing paginated InstanceSnapshot responses
func (resp *InstanceSnapshotsPagedResponse) AppendData(r *InstanceSnapshotsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of InstanceSnapshot
func (InstanceSnapshotsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(InstanceSnapshotsPagedResponse{})
}

// ListInstanceSnapshots lists InstanceSnapshots
func (c *Client) ListInstanceSnapshots(linodeID int, opts *ListOptions) ([]*InstanceSnapshot, error) {
	response := InstanceSnapshotsPagedResponse{}
	err := c.ListHelperWithID(response, linodeID, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// GetInstanceSnapshot gets the snapshot with the provided ID
func (c *Client) GetInstanceSnapshot(linodeID int, snapshotID int) (*InstanceSnapshot, error) {
	e, err := c.InstanceSnapshots.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, snapshotID)
	r, err := c.R().SetResult(&InstanceSnapshot{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceSnapshot).fixDates(), nil
}
