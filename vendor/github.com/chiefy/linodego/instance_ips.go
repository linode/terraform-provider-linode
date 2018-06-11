package linodego

import (
	"encoding/json"
	"fmt"
)

type InstanceIPAddressResponse struct {
	IPv4 *InstanceIPv4Response
	IPv6 *InstanceIPv6Response
}

type InstanceIPv4Response struct {
	Public  []*InstanceIP
	Private []*InstanceIP
	Shared  []*InstanceIP
}

type InstanceIP struct {
	Address    string
	Gateway    string
	SubnetMask string
	Prefix     int
	Type       string
	Public     bool
	RDNS       string
	LinodeID   int `json:"linode_id"`
	Region     string
}

type InstanceIPv6Response struct {
	LinkLocal *InstanceIP `json:"link_local"`
	SLAAC     *InstanceIP
	Global    []*IPv6Range
}

type IPv6Range struct {
	Range  string
	Region string
}

// GetInstanceIPAddresses gets the IPAddresses for a Linode instance
func (c *Client) GetInstanceIPAddresses(linodeID int) (*InstanceIPAddressResponse, error) {
	e, err := c.InstanceIPs.endpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	r, err := coupleAPIErrors(c.R().SetResult(&InstanceIPAddressResponse{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceIPAddressResponse), nil
}

// GetInstanceIPAddress gets the IPAddress for a Linode instance matching a supplied IP address
func (c *Client) GetInstanceIPAddress(linodeID int, ipaddress string) (*InstanceIP, error) {
	e, err := c.InstanceIPs.endpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, ipaddress)
	r, err := coupleAPIErrors(c.R().SetResult(&InstanceIP{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceIP), nil
}

// AddInstanceIPAddress adds a public or private IP to a Linode instance
func (c *Client) AddInstanceIPAddress(linodeID int, public bool) (*InstanceIP, error) {
	var body string
	e, err := c.InstanceIPs.endpointWithID(linodeID)
	if err != nil {
		return nil, err
	}

	req := c.R().SetResult(&InstanceIP{})

	instanceipRequest := struct {
		Type   string `json:"type"`
		Public bool   `json:"public"`
	}{"ipv4", true}

	if bodyData, err := json.Marshal(instanceipRequest); err == nil {
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

	return r.Result().(*InstanceIP), nil
}
