package linodego

import (
	"fmt"

	"github.com/go-resty/resty"
)

// LinodeRegion represents a linode region object
type Region struct {
	ID      string
	Country string
}

// LinodeRegionsPagedResponse represents a linode API response for listing
type RegionsPagedResponse struct {
	*PageOptions
	Data []*Region
}

// Endpoint gets the endpoint URL for Region
func (RegionsPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.Regions.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends Regions when processing paginated Region responses
func (resp *RegionsPagedResponse) AppendData(r *RegionsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of Region
func (RegionsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(RegionsPagedResponse{})
}

// ListRegions lists Regions
func (c *Client) ListRegions(opts *ListOptions) ([]*Region, error) {
	response := RegionsPagedResponse{}
	err := c.ListHelper(&response, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *Region) fixDates() *Region {
	return v
}

// GetRegion gets the template with the provided ID
func (c *Client) GetRegion(id string) (*Region, error) {
	e, err := c.Regions.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R().SetResult(&Region{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*Region).fixDates(), nil
}
