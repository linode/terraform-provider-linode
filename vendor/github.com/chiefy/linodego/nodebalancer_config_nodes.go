package linodego

import (
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

// NodeBalancerNodesPagedResponse represents a paginated NodeBalancerNode API response
type NodeBalancerNodesPagedResponse struct {
	*PageOptions
	Data []*NodeBalancerNode
}

// endpoint gets the endpoint URL for NodeBalancerNode
func (NodeBalancerNodesPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.NodeBalancerNodes.endpointWithID(id)
	if err != nil {
		panic(err)
	}
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
	r, err := c.R().SetResult(&NodeBalancerNode{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancerNode).fixDates(), nil
}
