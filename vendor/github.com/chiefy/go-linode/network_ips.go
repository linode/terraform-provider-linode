package golinode

import (
	"fmt"

	"github.com/go-resty/resty"
)

// IPAddressesPagedResponse represents a paginated IPAddress API response
type IPAddressesPagedResponse struct {
	*PageOptions
	Data []*InstanceIP
}

// Endpoint gets the endpoint URL for IPAddress
func (IPAddressesPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.IPAddresses.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends IPAddresses when processing paginated InstanceIPAddress responses
func (resp *IPAddressesPagedResponse) AppendData(r *IPAddressesPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of IPAddress
func (IPAddressesPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(IPAddressesPagedResponse{})
}

// ListIPAddresses lists IPAddresses
func (c *Client) ListIPAddresses(opts *ListOptions) ([]*InstanceIP, error) {
	response := IPAddressesPagedResponse{}
	err := c.ListHelper(response, opts)
	return response.Data, err
}

// GetIPAddress gets the template with the provided ID
func (c *Client) GetIPAddress(id string) (*InstanceIP, error) {
	e, err := c.IPAddresses.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R().SetResult(&InstanceIP{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceIP), nil
}
