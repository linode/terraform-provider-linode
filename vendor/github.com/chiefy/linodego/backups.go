package linodego

import (
	"context"
	"fmt"
)

// InstanceBackupsResponse response struct for backup snapshot
type InstanceBackupsResponse struct {
	Automatic []*InstanceSnapshot
	Snapshot  *InstanceBackupSnapshotResponse
}

type InstanceBackupSnapshotResponse struct {
	Current    *InstanceSnapshot
	InProgress *InstanceSnapshot `json:"in_progress"`
}

// GetInstanceBackups gets the Instance's available Backups
func (c *Client) GetInstanceBackups(ctx context.Context, linodeID int) (*InstanceBackupsResponse, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/backups", e, linodeID)
	r, err := coupleAPIErrors(c.R(ctx).
		SetResult(&InstanceBackupsResponse{}).
		Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceBackupsResponse).fixDates(), nil
}

func (l *InstanceBackupSnapshotResponse) fixDates() *InstanceBackupSnapshotResponse {
	if l.Current != nil {
		l.Current.fixDates()
	}
	if l.InProgress != nil {
		l.InProgress.fixDates()
	}
	return l
}

func (l *InstanceBackupsResponse) fixDates() *InstanceBackupsResponse {
	for _, el := range l.Automatic {
		el.fixDates()
	}
	if l.Snapshot != nil {
		l.Snapshot.fixDates()
	}
	return l
}
