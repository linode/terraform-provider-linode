//go:build unit

package firewall

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

// Unit tests for private functions in framework_models
// Functions under test: parseComputedAttributes, parseNonComputedAttributes, parseFirewallRules

func TestParseComputedAttributes(t *testing.T) {
	createdTime := time.Date(2023, time.August, 17, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, time.August, 18, 12, 0, 0, 0, time.UTC)

	firewall := &linodego.Firewall{
		ID:      123,
		Status:  linodego.FirewallEnabled,
		Created: &createdTime,
		Updated: &updatedTime,
	}

	deviceEntity1 := linodego.FirewallDeviceEntity{
		ID:    1234,
		Type:  linodego.FirewallDeviceLinode,
		Label: "device_entity_1",
		URL:   "test-firewall.example.com",
	}
	deviceEntity2 := linodego.FirewallDeviceEntity{
		ID:    4321,
		Type:  linodego.FirewallDeviceNodeBalancer,
		Label: "device_entity_2",
		URL:   "test-firewall.example.com",
	}
	devices := []linodego.FirewallDevice{
		{
			ID:     111,
			Entity: deviceEntity1,
		},
		{
			ID:     112,
			Entity: deviceEntity2,
		},
	}

	inboundRules := []linodego.FirewallRule{
		{
			Action:      "allow",
			Label:       "Rule 1",
			Description: "Allow SSH connections",
			Ports:       "22",
			Protocol:    "SSH",
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"192.168.1.1/24"},
				IPv6: &[]string{},
			},
		},
	}

	outboundRules := []linodego.FirewallRule{
		{
			Action:      "deny",
			Label:       "Rule 3",
			Description: "Allow SSH connections",
			Ports:       "22",
			Protocol:    "SSH",
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"192.168.1.3/24"},
				IPv6: &[]string{},
			},
		},
	}

	firewallRules := &linodego.FirewallRuleSet{
		InboundPolicy:  "ACCEPT",
		Inbound:        inboundRules,
		OutboundPolicy: "DROP",
		Outbound:       outboundRules,
	}

	data := &FirewallDataSourceModel{}
	diags := data.flattenFirewallForDataSource(context.Background(), firewall, devices, firewallRules)
	assert.Nil(t, diags)

	assert.Equal(t, int64(123), data.ID.ValueInt64())
	assert.Contains(t, data.Status.String(), string(linodego.FirewallEnabled))

	assert.Contains(t, data.Linodes.String(), "1234")
	assert.Contains(t, data.NodeBalancers.String(), "4321")

	assert.Nil(t, diags)

	assert.Contains(t, data.OutboundPolicy.String(), firewallRules.OutboundPolicy)
	assert.Contains(t, data.InboundPolicy.String(), firewallRules.InboundPolicy)

	assert.Equal(t, data.Inbound[0].Action.ValueString(), inboundRules[0].Action)
	assert.Equal(t, data.Inbound[0].Protocol.ValueString(), string(inboundRules[0].Protocol))
	assert.Equal(t, data.Inbound[0].Ports.ValueString(), inboundRules[0].Ports)
	assert.Contains(t, data.Inbound[0].IPv4.String(), (*inboundRules[0].Addresses.IPv4)[0])

	assert.Equal(t, data.Outbound[0].Action.ValueString(), outboundRules[0].Action)
	assert.Equal(t, data.Outbound[0].Protocol.ValueString(), string(outboundRules[0].Protocol))
	assert.Equal(t, data.Outbound[0].Ports.ValueString(), outboundRules[0].Ports)
	assert.Contains(t, data.Outbound[0].IPv4.String(), (*outboundRules[0].Addresses.IPv4)[0])

	assert.Equal(t, data.Devices[0].ID.ValueInt64(), int64(111))
	assert.Equal(t, data.Devices[1].ID.ValueInt64(), int64(112))
}
