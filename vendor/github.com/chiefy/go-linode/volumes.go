package golinode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

// Volume represents a linode volume object
type Volume struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID       int
	Label    string
	Status   string
	Region   string
	Size     int
	LinodeID int        `json:"linode_id"`
	Created  *time.Time `json:"-"`
	Updated  *time.Time `json:"-"`
}

type VolumeAttachOptions struct {
	LinodeID int
	ConfigID int
}

// LinodeVolumesPagedResponse represents a linode API response for listing of volumes
type VolumesPagedResponse struct {
	*PageOptions
	Data []*Volume
}

// Endpoint gets the endpoint URL for Volume
func (VolumesPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.Volumes.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends Volumes when processing paginated Volume responses
func (resp *VolumesPagedResponse) AppendData(r *VolumesPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of Volume
func (VolumesPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(VolumesPagedResponse{})
}

// ListVolumes lists Volumes
func (c *Client) ListVolumes(opts *ListOptions) ([]*Volume, error) {
	response := VolumesPagedResponse{}
	err := c.ListHelper(response, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *Volume) fixDates() *Volume {
	v.Created, _ = parseDates(v.CreatedStr)
	v.Updated, _ = parseDates(v.UpdatedStr)
	return v
}

// GetVolume gets the template with the provided ID
func (c *Client) GetVolume(id int) (*Volume, error) {
	e, err := c.Volumes.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)
	r, err := c.R().SetResult(&Volume{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*Volume).fixDates(), nil
}

// AttachVolume attaches volume to linode instance
func (c *Client) AttachVolume(id int, options *VolumeAttachOptions) (bool, error) {
	body := ""
	if bodyData, err := json.Marshal(options); err == nil {
		body = string(bodyData)
	} else {
		return false, err
	}

	e, err := c.Volumes.Endpoint()
	if err != nil {
		return false, err
	}

	e = fmt.Sprintf("%s/%d/attach", e, id)
	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	return settleBoolResponseOrError(resp, err)
}

// CloneVolume clones a Linode instance
func (c *Client) CloneVolume(id int, label string) (*Volume, error) {
	body := fmt.Sprintf("{\"label\":\"%s\"}", label)

	e, err := c.Volumes.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/clone", e, id)

	req := c.R().SetResult(&Volume{})

	resp, err := req.
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*Volume).fixDates(), nil
}

// DetachVolume detaches a Linode instance
func (c *Client) DetachVolume(id int) (bool, error) {
	body := ""

	e, err := c.Volumes.Endpoint()
	if err != nil {
		return false, err
	}

	e = fmt.Sprintf("%s/%d/detach", e, id)

	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	return settleBoolResponseOrError(resp, err)
}

// ResizeVolume resizes an instance to new Linode type
func (c *Client) ResizeVolume(id int, size int) (bool, error) {
	body := fmt.Sprintf("{\"size\":\"%d\"}", size)

	e, err := c.Volumes.Endpoint()
	if err != nil {
		return false, err
	}
	e = fmt.Sprintf("%s/%d/resize", e, id)

	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	return settleBoolResponseOrError(resp, err)
}
