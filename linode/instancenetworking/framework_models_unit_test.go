//go:build unit

package instancenetworking

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseInstanceIPAddressResponse(t *testing.T) {
	instanceIPResponse := &linodego.InstanceIPAddressResponse{
		IPv4: &linodego.InstanceIPv4Response{
			Public: []*linodego.InstanceIP{
				{
					Address: "1.2.3.4",
					Type:    "ipv4",
					Public:  true,
				},
			},
			Private: []*linodego.InstanceIP{
				{
					Address: "10.0.0.1",
					Type:    "ipv4",
					Public:  false,
				},
			},
		},
		IPv6: &linodego.InstanceIPv6Response{
			LinkLocal: &linodego.InstanceIP{
				Address: "fe80::1",
				Type:    "ipv6",
			},
			SLAAC: &linodego.InstanceIP{
				Address: "fe80::1",
				Type:    "ipv6",
			},
		},
	}

	dataSourceModel := &DataSourceModel{}

	diags := dataSourceModel.parseInstanceIPAddressResponse(context.Background(), instanceIPResponse)

	assert.False(t, diags.HasError())

	assert.Contains(t, dataSourceModel.IPV4.String(), "1.2.3.4")
	assert.Contains(t, dataSourceModel.IPV6.String(), "fe80::1")
}
