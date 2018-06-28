package linodego

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty"
)

type NodeBalancerNode struct {
	ID             int
	Address        string
	Label          string
	Status         string
	Weight         int
	Mode           string
	ConfigID       int `json:"config_id"`
	NodeBalancerID int `json:"nodebalancer_id"`
}

type NodeBalancerNodeCreateOptions struct {
	Address        string `json:"address"`
	Label          string `json:"label"`
	Weight         int    `json:"weight"`
	Mode           string `json:"mode"`
	ConfigID       int    `json:"config_id"`
	NodeBalancerID int    `json:"nodebalancer_id"`
}

type NodeBalancerNodeUpdateOptions struct {
	Address string `json:"address"`
	Label   string `json:"label"`
	Weight  int    `json:"weight"`
	Mode    string `json:"mode"`
}

// NodeBalancerNodesPagedResponse represents a paginated NodeBalancerNode API response
type NodeBalancerNodesPagedResponse struct {
	*PageOptions
	Data []*NodeBalancerNode
}

// endpoint gets the endpoint URL for NodeBalancerNode
func (NodeBalancerNodesPagedResponse) endpointWithTwoIDs(c *Client, nodebalancerID int, configID int) string {
	endpoint, err := c.NodeBalancerConfigs.endpointWithID(nodebalancerID)
	if err != nil {
		panic(err)
	}
	endpoint = fmt.Sprintf("%s/%d/nodes/", endpoint, configID)
	return endpoint
}

// appendData appends NodeBalancerNodes when processing paginated NodeBalancerNode responses
func (resp *NodeBalancerNodesPagedResponse) appendData(r *NodeBalancerNodesPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// setResult sets the Resty response type of NodeBalancerNode
func (NodeBalancerNodesPagedResponse) setResult(r *resty.Request) {
	r.SetResult(NodeBalancerNodesPagedResponse{})
}

// ListNodeBalancerNodes lists NodeBalancerNodes
func (c *Client) ListNodeBalancerNodes(nodebalancerID int, configID int, opts *ListOptions) ([]*NodeBalancerNode, error) {
	response := NodeBalancerNodesPagedResponse{}
	err := c.listHelperWithID(&response, nodebalancerID, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *NodeBalancerNode) fixDates() *NodeBalancerNode {
	return v
}

// GetNodeBalancerNode gets the template with the provided ID
func (c *Client) GetNodeBalancerNode(nodebalancerID int, configID int, nodeID int) (*NodeBalancerNode, error) {
	e, err := c.NodeBalancerConfigs.endpointWithID(nodebalancerID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/nodes/%d", e, configID, nodeID)
	r, err := coupleAPIErrors(c.R().SetResult(&NodeBalancerNode{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancerNode).fixDates(), nil
}

// CreateNodeBalancerNode creates a NodeBalancerNode
func (c *Client) CreateNodeBalancerNode(nodebalancerID int, configID int, createOpts *NodeBalancerNodeCreateOptions) (*NodeBalancerNode, error) {
	var body string
	e, err := c.NodeBalancerNodes.endpointWithID(nodebalancerID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/nodes/", e, configID)

	req := c.R().SetResult(&NodeBalancerNode{})

	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Post(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancerNode).fixDates(), nil
}

// UpdateNodeBalancerNode updates the NodeBalancerNode with the specified id
func (c *Client) UpdateNodeBalancerNode(nodebalancerID int, configID int, nodeID int, updateOpts NodeBalancerNodeUpdateOptions) (*NodeBalancerNode, error) {
	var body string
	e, err := c.NodeBalancers.endpointWithID(nodebalancerID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/configs/%d/nodes/%d", e, configID, nodeID)

	req := c.R().SetResult(&NodeBalancerNode{})

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
	return r.Result().(*NodeBalancerNode).fixDates(), nil
}

// DeleteNodeBalancerNode deletes the NodeBalancerNode with the specified id
func (c *Client) DeleteNodeBalancerNode(nodebalancerID int, configID int, nodeID int) error {
	e, err := c.NodeBalancers.endpointWithID(nodebalancerID)
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/configs/%d/nodes/%d", e, configID, nodeID)

	if _, err := coupleAPIErrors(c.R().Delete(e)); err != nil {
		return err
	}

	return nil
}
