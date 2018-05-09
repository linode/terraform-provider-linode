package golinode

import (
	"fmt"

	"github.com/go-resty/resty"
)

// IPv6RangesPagedResponse represents a paginated IPv6Range API response
type IPv6RangesPagedResponse struct {
	*PageOptions
	Data []*IPv6Range
}

// Endpoint gets the endpoint URL for IPv6Range
func (IPv6RangesPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.IPv6Ranges.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends IPv6Ranges when processing paginated IPv6Range responses
func (resp *IPv6RangesPagedResponse) AppendData(r *IPv6RangesPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of IPv6Range
func (IPv6RangesPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(IPv6RangesPagedResponse{})
}

// ListIPv6Ranges lists IPv6Ranges
func (c *Client) ListIPv6Ranges(opts *ListOptions) ([]*IPv6Range, error) {
	response := IPv6RangesPagedResponse{}
	err := c.ListHelper(response, opts)
	return response.Data, err
}

// GetIPv6Range gets the template with the provided ID
func (c *Client) GetIPv6Range(id string) (*IPv6Range, error) {
	e, err := c.IPv6Ranges.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R().SetResult(&IPv6Range{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*IPv6Range), nil
}
