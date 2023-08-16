//go:build integration

package firewall

import (
	"context"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Unit tests for private functions in framework_models
// Functions under test: parseComputedAttributes, parseNonComputedAttributes, parseFirewallRules, parseFirewallLinodes

func TestParseComputedAttributes(t *testing.T) {
	firewall := &linodego.Firewall{
		ID:     123,
		Status: linodego.FirewallEnabled,
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

func TestParseFirewallLinodes(t *testing.T) {
	deviceEntity1 := linodego.FirewallDeviceEntity{
		ID:    1234,
		Type:  linodego.FirewallDeviceLinode,
		Label: "device_entity_1",
		URL:   "test-firewall.example.com",
	}

	deviceEntity2 := linodego.FirewallDeviceEntity{
		ID:    1235,
		Type:  linodego.FirewallDeviceLinode,
		Label: "device_entity_2",
		URL:   "test-firewall.example-2.com",
	}

	created := time.Date(2022, time.January, 15, 12, 30, 0, 0, time.UTC)
	updated := time.Date(2023, time.March, 7, 9, 45, 0, 0, time.UTC)

	device1 := linodego.FirewallDevice{
		ID:      123,
		Entity:  deviceEntity1,
		Created: &created,
		Updated: &updated,
	}

	device2 := linodego.FirewallDevice{
		ID:      123,
		Entity:  deviceEntity2,
		Created: &created,
		Updated: &updated,
	}

	devices := []linodego.FirewallDevice{
		device1, device2,
	}

	cases := []struct {
		devices  []linodego.FirewallDevice
		expected []int
	}{
		{
			devices:  devices,
			expected: []int{1234, 1235},
		},
	}

	for _, c := range cases {
		out := parseFirewallLinodes(c.devices)

		assert.ElementsMatch(t, c.expected, out)
	}
}
