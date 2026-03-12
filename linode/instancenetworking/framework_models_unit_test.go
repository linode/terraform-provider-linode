//go:build unit

package instancenetworking

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

	var diags diag.Diagnostics
	dataSourceModel.parseInstanceIPAddressResponse(instanceIPResponse, &diags)

	assert.False(t, diags.HasError())

	assert.Contains(t, dataSourceModel.IPV4.String(), "1.2.3.4")
	assert.Contains(t, dataSourceModel.IPV6.String(), "fe80::1")
}

func TestParseInstanceIPAddressResponse_VPCDualStack(t *testing.T) {
	ipv4Addr := "10.0.0.5"
	ipv6Range := "fd00:abcd:1234:5678::/64"
	isPublic := false
	slaacAddr := "fd00:abcd:1234:5678::1"

	instanceIPResponse := &linodego.InstanceIPAddressResponse{
		IPv4: &linodego.InstanceIPv4Response{
			Public: []*linodego.InstanceIP{
				{
					Address: "1.2.3.4",
					Type:    "ipv4",
					Public:  true,
				},
			},
			VPC: []*linodego.VPCIP{
				{
					Address:     &ipv4Addr,
					Active:      true,
					VPCID:       100,
					SubnetID:    200,
					ConfigID:    300,
					InterfaceID: 400,
					Gateway:     "10.0.0.1",
					Prefix:      24,
					Region:      "us-east",
					SubnetMask:  "255.255.255.0",
					LinodeID:    500,
				},
			},
		},
		IPv6: &linodego.InstanceIPv6Response{
			LinkLocal: &linodego.InstanceIP{
				Address: "fe80::1",
				Type:    "ipv6",
			},
			SLAAC: &linodego.InstanceIP{
				Address: "2600:3c00::1",
				Type:    "ipv6",
			},
			VPC: []linodego.VPCIP{
				{
					Active:       true,
					VPCID:        100,
					SubnetID:     200,
					ConfigID:     300,
					InterfaceID:  400,
					Region:       "us-east",
					LinodeID:     500,
					IPv6Range:    &ipv6Range,
					IPv6IsPublic: &isPublic,
					IPv6Addresses: []linodego.VPCIPIPv6Address{
						{SLAACAddress: slaacAddr},
					},
				},
			},
		},
	}

	dataSourceModel := &DataSourceModel{}

	var diags diag.Diagnostics
	dataSourceModel.parseInstanceIPAddressResponse(instanceIPResponse, &diags)

	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	// Verify IPv4 VPC IP is present
	assert.Contains(t, dataSourceModel.IPV4.String(), "10.0.0.5")

	// Verify IPv6 VPC IP fields are present
	ipv6Str := dataSourceModel.IPV6.String()
	assert.Contains(t, ipv6Str, ipv6Range)
	assert.Contains(t, ipv6Str, slaacAddr)
}
