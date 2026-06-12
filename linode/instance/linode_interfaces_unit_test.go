//go:build unit

package instance

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpandLinodeInstanceInterfaces_RDMAVPC(t *testing.T) {
	input := []any{
		map[string]any{
			"firewall_id":   0,
			"default_route": []any{},
			"public":        []any{},
			"vlan":          []any{},
			"vpc":           []any{},
			"rdma_vpc": []any{
				map[string]any{
					"subnet_id": 1234,
					"ipv4": []any{
						map[string]any{
							"addresses": []any{
								map[string]any{
									"address": "10.0.0.5",
									"primary": true,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := expandLinodeInstanceInterfaces(input)
	require.NoError(t, err)
	require.Len(t, result, 1)

	require.Nil(t, result[0].Public)
	require.Nil(t, result[0].VLAN)
	require.Nil(t, result[0].VPC)
	require.NotNil(t, result[0].RDMAVPC)

	require.Equal(t, 1234, result[0].RDMAVPC.SubnetID)
	require.Len(t, result[0].RDMAVPC.IPv4.Addresses, 1)
	require.Equal(t, "10.0.0.5", result[0].RDMAVPC.IPv4.Addresses[0].Address)
	require.True(t, result[0].RDMAVPC.IPv4.Addresses[0].Primary)
}

func TestExpandLinodeInstanceInterfaces_Public(t *testing.T) {
	input := []any{
		map[string]any{
			"firewall_id": 0,
			"default_route": []any{
				map[string]any{
					"ipv4": true,
					"ipv6": false,
				},
			},
			"public": []any{
				map[string]any{
					"ipv4": []any{
						map[string]any{
							"addresses": []any{
								map[string]any{
									"address": "auto",
									"primary": true,
								},
							},
						},
					},
					"ipv6": []any{},
				},
			},
			"vlan":     []any{},
			"vpc":      []any{},
			"rdma_vpc": []any{},
		},
	}

	result, err := expandLinodeInstanceInterfaces(input)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotNil(t, result[0].Public)
	require.Nil(t, result[0].VLAN)
	require.Nil(t, result[0].VPC)
	require.Nil(t, result[0].RDMAVPC)
	require.NotNil(t, result[0].DefaultRoute)
	require.Equal(t, true, *result[0].DefaultRoute.IPv4)
}

func TestExpandLinodeInstanceInterfaces_VPC(t *testing.T) {
	input := []any{
		map[string]any{
			"firewall_id":   0,
			"default_route": []any{},
			"public":        []any{},
			"vlan":          []any{},
			"rdma_vpc":      []any{},
			"vpc": []any{
				map[string]any{
					"subnet_id": 42,
					"ipv4": []any{
						map[string]any{
							"addresses": []any{
								map[string]any{
									"address":         "10.0.0.100",
									"primary":         true,
									"nat_1_1_address": "203.0.113.5",
								},
							},
							"ranges": []any{},
						},
					},
				},
			},
		},
	}

	result, err := expandLinodeInstanceInterfaces(input)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotNil(t, result[0].VPC)
	require.Equal(t, 42, result[0].VPC.SubnetID)
	require.NotNil(t, result[0].VPC.IPv4)
	require.NotNil(t, result[0].VPC.IPv4.Addresses)
	addrs := *result[0].VPC.IPv4.Addresses
	require.Len(t, addrs, 1)
	require.Equal(t, "10.0.0.100", *addrs[0].Address)
	require.Equal(t, true, *addrs[0].Primary)
	require.Equal(t, "203.0.113.5", *addrs[0].NAT1To1Address)
}

func TestExpandLinodeInstanceInterfaces_VLAN(t *testing.T) {
	input := []any{
		map[string]any{
			"firewall_id":   0,
			"default_route": []any{},
			"public":        []any{},
			"rdma_vpc":      []any{},
			"vpc":           []any{},
			"vlan": []any{
				map[string]any{
					"vlan_label":   "my-vlan",
					"ipam_address": "10.0.0.1/24",
				},
			},
		},
	}

	result, err := expandLinodeInstanceInterfaces(input)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotNil(t, result[0].VLAN)
	require.Equal(t, "my-vlan", result[0].VLAN.VLANLabel)
	require.Equal(t, "10.0.0.1/24", *result[0].VLAN.IPAMAddress)
}

func TestExpandLinodeInstanceInterfaces_ConflictError(t *testing.T) {
	// Setting both public and vpc should error
	input := []any{
		map[string]any{
			"firewall_id":   0,
			"default_route": []any{},
			"rdma_vpc":      []any{},
			"vlan":          []any{},
			"public": []any{
				map[string]any{
					"ipv4": []any{},
					"ipv6": []any{},
				},
			},
			"vpc": []any{
				map[string]any{
					"subnet_id": 10,
					"ipv4":      []any{},
				},
			},
		},
	}

	_, err := expandLinodeInstanceInterfaces(input)
	require.Error(t, err)
	require.Contains(t, err.Error(), "at most one of")
}
