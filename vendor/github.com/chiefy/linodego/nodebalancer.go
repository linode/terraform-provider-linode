package linodego

import (
	"encoding/json"
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

// NodeBalancerCreateOptions are the options permitted for CreateNodeBalancer
type NodeBalancerCreateOptions struct {
	Label              string `json:"label"`
	Region             string `json:"region"`
	ClientConnThrottle int    `json:"client_conn_throttle"`
}

// NodeBalancerUpdateOptions are the options permitted for UpdateNodeBalancer
type NodeBalancerUpdateOptions NodeBalancerCreateOptions

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

func (n *NodeBalancer) fixDates() *NodeBalancer {
	n.Created, _ = parseDates(n.CreatedStr)
	n.Updated, _ = parseDates(n.UpdatedStr)
	return n
}

// GetNodeBalancer gets the NodeBalancer with the provided ID
func (c *Client) GetNodeBalancer(id int) (*NodeBalancer, error) {
	e, err := c.NodeBalancers.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)
	r, err := c.R().
		SetResult(&NodeBalancer{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancer), nil
}

// CreateNodeBalancer creates a NodeBalancer
func (c *Client) CreateNodeBalancer(nodebalancer *NodeBalancerCreateOptions) (*NodeBalancer, error) {
	var body string
	e, err := c.NodeBalancers.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R().SetResult(&NodeBalancer{})

	if bodyData, err := json.Marshal(nodebalancer); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancer).fixDates(), nil
}

// UpdateNodeBalancer updates the NodeBalancer with the specified id
func (c *Client) UpdateNodeBalancer(id int, updateOpts NodeBalancerUpdateOptions) (*NodeBalancer, error) {
	var body string
	e, err := c.NodeBalancers.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	req := c.R().SetResult(&NodeBalancer{})

	if bodyData, err := json.Marshal(updateOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Put(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancer).fixDates(), nil
}

// DeleteNodeBalancer deletes the NodeBalancer with the specified id
func (c *Client) DeleteNodeBalancer(id int) error {
	e, err := c.NodeBalancers.Endpoint()
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	if _, err := coupleAPIErrors(c.R().Delete(e)); err != nil {
		return err
	}

	return nil
}
