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
	devices := []linodego.FirewallDevice{
		{
			ID:     111,
			Entity: deviceEntity1,
		},
	}

	data := &FirewallModel{}
	diags := data.parseComputedAttributes(context.Background(), firewall, nil, devices)
	assert.Nil(t, diags)

	assert.Equal(t, int64(123), data.ID.ValueInt64())
	assert.Contains(t, data.Status.String(), string(linodego.FirewallEnabled))

	expectedLinodeID := "1234"
	assert.Contains(t, data.Linodes.String(), expectedLinodeID)

	expectedDevicesID := "111"
	assert.Contains(t, data.Devices.String(), expectedDevicesID)
}

func TestParseNonComputedAttributes(t *testing.T) {
	ctx := context.Background()

	firewall := &linodego.Firewall{
		ID:     123,
		Status: linodego.FirewallEnabled,
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

	data := &FirewallModel{}

	diags := data.parseNonComputedAttributes(ctx, firewall, firewallRules, nil)
	assert.Nil(t, diags)

	assert.Contains(t, data.OutboundPolicy.String(), firewallRules.OutboundPolicy)
	assert.Contains(t, data.InboundPolicy.String(), firewallRules.InboundPolicy)

	assert.Contains(t, data.Inbound.String(), inboundRules[0].Action)
	assert.Contains(t, data.Inbound.String(), inboundRules[0].Protocol)
	assert.Contains(t, data.Inbound.String(), inboundRules[0].Ports)
	assert.Contains(t, data.Inbound.String(), (*inboundRules[0].Addresses.IPv4)[0])

	assert.Contains(t, data.Outbound.String(), outboundRules[0].Action)
	assert.Contains(t, data.Outbound.String(), outboundRules[0].Protocol)
	assert.Contains(t, data.Outbound.String(), outboundRules[0].Ports)
	assert.Contains(t, data.Outbound.String(), (*outboundRules[0].Addresses.IPv4)[0])
}
