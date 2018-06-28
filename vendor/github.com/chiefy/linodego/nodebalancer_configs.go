package linodego

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty"
)

type NodeBalancerConfig struct {
	ID             int
	Port           int
	Protocol       ConfigProtocol
	Algorithm      ConfigAlgorithm
	Stickiness     ConfigStickiness
	Check          ConfigCheck
	CheckInterval  int                     `json:"check_interval"`
	CheckAttempts  int                     `json:"check_attempts"`
	CheckPath      string                  `json:"check_path"`
	CheckBody      string                  `json:"check_body"`
	CheckPassive   bool                    `json:"check_passive"`
	CipherSuite    ConfigCipher            `json:"cipher_suite"`
	NodeBalancerID int                     `json:"nodebalancer_id"`
	SSLCommonName  string                  `json:"ssl_commonname"`
	SSLFingerprint string                  `json:"ssl_fingerprint"`
	SSLCert        string                  `json:"ssl_cert"`
	SSLKey         string                  `json:"ssl_key"`
	NodesStatus    *NodeBalancerNodeStatus `json:"nodes_status"`
}

type ConfigAlgorithm string

var (
	AlgorithmRoundRobin ConfigAlgorithm = "roundrobin"
	AlgorithmLeastConn  ConfigAlgorithm = "leastconn"
	AlgorithmSource     ConfigAlgorithm = "source"
)

type ConfigStickiness string

var (
	StickinessNone       ConfigStickiness = "none"
	StickinessTable      ConfigStickiness = "table"
	StickinessHTTPCookie ConfigStickiness = "http_cookie"
)

type ConfigCheck string

var (
	CheckNone       ConfigCheck = "none"
	CheckConnection ConfigCheck = "connection"
	CheckHTTP       ConfigCheck = "http"
	CheckHTTPBody   ConfigCheck = "http_body"
)

type ConfigProtocol string

var (
	ProtocolHTTP  ConfigProtocol = "http"
	ProtocolHTTPS ConfigProtocol = "https"
	ProtocolTCP   ConfigProtocol = "tcp"
)

type ConfigCipher string

var (
	CipherRecommended ConfigCipher = "recommended"
	CipherLegacy      ConfigCipher = "legacy"
)

type NodeBalancerNodeStatus struct {
	Up   int
	Down int
}

// NodeBalancerConfigUpdateOptions are permitted by CreateNodeBalancerConfig
type NodeBalancerConfigCreateOptions struct {
	NodeBalancerID int              `json:"nodebalancer_id"`
	Port           int              `json:"port"`
	Protocol       ConfigProtocol   `json:"protocol"`
	Algorithm      ConfigAlgorithm  `json:"algorithm"`
	Stickiness     ConfigStickiness `json:"stickiness"`
	Check          ConfigCheck      `json:"check"`
	CheckInterval  int              `json:"check_interval"`
	CheckAttempts  int              `json:"check_attempts"`
	CheckPath      string           `json:"check_path"`
	CheckBody      string           `json:"check_body"`
	CheckPassive   bool             `json:"check_passive"`
	CipherSuite    ConfigCipher     `json:"cipher_suite"`
	SSLCommonName  string           `json:"ssl_commonname"`
	SSLFingerprint string           `json:"ssl_fingerprint"`
	SSLCert        string           `json:"ssl_cert"`
	SSLKey         string           `json:"ssl_key"`
}

// NodeBalancerConfigUpdateOptions are permitted by UpdateNodeBalancerConfig
type NodeBalancerConfigUpdateOptions NodeBalancerConfigCreateOptions

// NodeBalancerConfigsPagedResponse represents a paginated NodeBalancerConfig API response
type NodeBalancerConfigsPagedResponse struct {
	*PageOptions
	Data []*NodeBalancerConfig
}

// endpoint gets the endpoint URL for NodeBalancerConfig
func (NodeBalancerConfigsPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.NodeBalancerConfigs.endpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends NodeBalancerConfigs when processing paginated NodeBalancerConfig responses
func (resp *NodeBalancerConfigsPagedResponse) appendData(r *NodeBalancerConfigsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// setResult sets the Resty response type of NodeBalancerConfig
func (NodeBalancerConfigsPagedResponse) setResult(r *resty.Request) {
	r.SetResult(NodeBalancerConfigsPagedResponse{})
}

// ListNodeBalancerConfigs lists NodeBalancerConfigs
func (c *Client) ListNodeBalancerConfigs(nodebalancerID int, opts *ListOptions) ([]*NodeBalancerConfig, error) {
	response := NodeBalancerConfigsPagedResponse{}
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
func (v *NodeBalancerConfig) fixDates() *NodeBalancerConfig {
	return v
}

// GetNodeBalancerConfig gets the template with the provided ID
func (c *Client) GetNodeBalancerConfig(nodebalancerID int, configID int) (*NodeBalancerConfig, error) {
	e, err := c.NodeBalancerConfigs.endpointWithID(nodebalancerID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, configID)
	r, err := coupleAPIErrors(c.R().SetResult(&NodeBalancerConfig{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*NodeBalancerConfig).fixDates(), nil
}

// CreateNodeBalancerConfig creates a NodeBalancerConfig
func (c *Client) CreateNodeBalancerConfig(nodebalancerConfig *NodeBalancerConfigCreateOptions) (*NodeBalancerConfig, error) {
	var body string
	e, err := c.NodeBalancerConfigs.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R().SetResult(&NodeBalancerConfig{})

	if bodyData, err := json.Marshal(nodebalancerConfig); err == nil {
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
	return r.Result().(*NodeBalancerConfig).fixDates(), nil
}

// UpdateNodeBalancerConfig updates the NodeBalancerConfig with the specified id
func (c *Client) UpdateNodeBalancerConfig(nodebalancerID int, configID int, updateOpts NodeBalancerConfigUpdateOptions) (*NodeBalancerConfig, error) {
	var body string
	e, err := c.NodeBalancers.endpointWithID(nodebalancerID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/configs/%d", e, configID)

	req := c.R().SetResult(&NodeBalancerConfig{})

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
	return r.Result().(*NodeBalancerConfig).fixDates(), nil
}

// DeleteNodeBalancerConfig deletes the NodeBalancerConfig with the specified id
func (c *Client) DeleteNodeBalancerConfig(nodebalancerID int, configID int) error {
	e, err := c.NodeBalancers.endpointWithID(nodebalancerID)
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/configs/%d", e, configID)

	if _, err := coupleAPIErrors(c.R().Delete(e)); err != nil {
		return err
	}

	return nil
}
