package linodego

import (
	"fmt"

	"github.com/go-resty/resty"
)

// LinodeType represents a linode type object
type LinodeType struct {
	ID         string
	Disk       int
	Class      string // enum: nanode, standard, highmem
	Price      *LinodePrice
	Label      string
	Addons     *LinodeAddons
	NetworkOut int `json:"network_out"`
	Memory     int
	Transfer   int
	VCPUs      int
}

// LinodePrice represents a linode type price object
type LinodePrice struct {
	Hourly  float32
	Monthly float32
}

// LinodeBackupsAddon represents a linode backups addon object
type LinodeBackupsAddon struct {
	Price *LinodePrice
}

// LinodeAddons represent the linode addons object
type LinodeAddons struct {
	Backups *LinodeBackupsAddon
}

// LinodeTypesPagedResponse represents a linode types API response for listing
type LinodeTypesPagedResponse struct {
	*PageOptions
	Data []*LinodeType
}

func (LinodeTypesPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Types.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *LinodeTypesPagedResponse) appendData(r *LinodeTypesPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

func (LinodeTypesPagedResponse) setResult(r *resty.Request) {
	r.SetResult(LinodeTypesPagedResponse{})
}

// ListTypes lists linode types
func (c *Client) ListTypes(opts *ListOptions) ([]*LinodeType, error) {
	response := LinodeTypesPagedResponse{}
	err := c.listHelper(&response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetType gets the type with the provided ID
func (c *Client) GetType(typeID string) (*LinodeType, error) {
	e, err := c.Types.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, typeID)

	r, err := coupleAPIErrors(c.Types.R().Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeType), nil
}
