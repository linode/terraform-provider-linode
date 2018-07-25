package linodego

import (
	"context"
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

// endpoint gets the endpoint URL for InstanceSnapshot
func (InstanceSnapshotsPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.InstanceSnapshots.endpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends InstanceSnapshots when processing paginated InstanceSnapshot responses
func (resp *InstanceSnapshotsPagedResponse) appendData(r *InstanceSnapshotsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// setResult sets the Resty response type of InstanceSnapshot
func (InstanceSnapshotsPagedResponse) setResult(r *resty.Request) {
	r.SetResult(InstanceSnapshotsPagedResponse{})
}

// ListInstanceSnapshots lists InstanceSnapshots
func (c *Client) ListInstanceSnapshots(ctx context.Context, linodeID int, opts *ListOptions) ([]*InstanceSnapshot, error) {
	response := InstanceSnapshotsPagedResponse{}
	err := c.listHelperWithID(ctx, &response, linodeID, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetInstanceSnapshot gets the snapshot with the provided ID
func (c *Client) GetInstanceSnapshot(ctx context.Context, linodeID int, snapshotID int) (*InstanceSnapshot, error) {
	e, err := c.InstanceSnapshots.endpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, snapshotID)
	r, err := coupleAPIErrors(c.R(ctx).SetResult(&InstanceSnapshot{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceSnapshot).fixDates(), nil
}
