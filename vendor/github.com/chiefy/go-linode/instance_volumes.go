package golinode

import (
	"fmt"

	"github.com/go-resty/resty"
)

// InstanceVolumesPagedResponse represents a paginated InstanceVolume API response
type InstanceVolumesPagedResponse struct {
	*PageOptions
	Data []*Volume
}

// Endpoint gets the endpoint URL for InstanceVolume
func (InstanceVolumesPagedResponse) EndpointWithID(c *Client, id int) string {
	endpoint, err := c.InstanceVolumes.EndpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends InstanceVolumes when processing paginated InstanceVolume responses
func (resp *InstanceVolumesPagedResponse) AppendData(r *InstanceVolumesPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of InstanceVolume
func (InstanceVolumesPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(InstanceVolumesPagedResponse{})
}

// ListInstanceVolumes lists InstanceVolumes
func (c *Client) ListInstanceVolumes(linodeID int, opts *ListOptions) ([]*Volume, error) {
	response := InstanceVolumesPagedResponse{}
	err := c.ListHelperWithID(response, linodeID, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// GetInstanceVolume gets the snapshot with the provided ID
func (c *Client) GetInstanceVolume(linodeID int, snapshotID int) (*Volume, error) {
	e, err := c.InstanceVolumes.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, snapshotID)
	r, err := c.R().SetResult(&Volume{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*Volume).fixDates(), nil
}
