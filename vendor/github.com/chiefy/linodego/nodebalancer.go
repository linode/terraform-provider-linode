package linodego

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

// NodeBalancer represents a NodeBalancer object
type NodeBalancer struct {
	CreatedStr         string `json:"created"`
	UpdatedStr         string `json:"updated"`
	ID                 int
	Label              string
	Region             string
	Hostname           string
	IPv4               string
	IPv6               string
	ClientConnThrottle int        `json:"client_conn_throttle"`
	Created            *time.Time `json:"-"`
	Updated            *time.Time `json:"-"`
}

// NodeBalancersPagedResponse represents a paginated NodeBalancer API response
type NodeBalancersPagedResponse struct {
	*PageOptions
	Data []*NodeBalancer
}

func (NodeBalancersPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.NodeBalancers.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *NodeBalancersPagedResponse) appendData(r *NodeBalancersPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

func (NodeBalancersPagedResponse) setResult(r *resty.Request) {
	r.SetResult(NodeBalancersPagedResponse{})
}

// ListNodeBalancers lists NodeBalancers
func (c *Client) ListNodeBalancers(opts *ListOptions) ([]*NodeBalancer, error) {
	response := NodeBalancersPagedResponse{}
	err := c.listHelper(&response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetNodeBalancer gets the NodeBalancer with the provided ID
func (c *Client) GetNodeBalancer(id string) (*NodeBalancer, error) {
	e, err := c.NodeBalancers.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R().
		SetResult(&NodeBalancer{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancer), nil
}
