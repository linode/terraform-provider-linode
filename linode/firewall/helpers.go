package firewall

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/linode/linodego"
	"golang.org/x/net/context"
)

// firewallDeviceAssignment is a helper struct intended to be used in conjunction
// with updateFirewallDevices.
type firewallDeviceAssignment struct {
	ID   int
	Type linodego.FirewallDeviceType
}

func expandFirewallStatus(disabled bool) linodego.FirewallStatus {
	return map[bool]linodego.FirewallStatus{
		true:  linodego.FirewallDisabled,
		false: linodego.FirewallEnabled,
	}[disabled]
}

func fwUpdateFirewallDevices(
	ctx context.Context,
	client linodego.Client,
	id int,
	configuredDevices []firewallDeviceAssignment,
) error {
	currentDevices, err := client.ListFirewallDevices(ctx, id, nil)
	if err != nil {
		return err
	}

	// Populate a map to track existing devices by assignment
	deviceMap := make(map[firewallDeviceAssignment]linodego.FirewallDevice)
	for _, device := range currentDevices {
		deviceMap[firewallDeviceAssignment{ID: device.Entity.ID, Type: device.Entity.Type}] = device
	}

	for _, device := range configuredDevices {
		if _, ok := deviceMap[device]; ok {
			// Device exists, drop it from the map so it won't be removed
			delete(deviceMap, device)
			continue
		}

		// Device doesn't exist, create a new one
		createOpts := linodego.FirewallDeviceCreateOptions{
			ID:   device.ID,
			Type: device.Type,
		}

		tflog.Debug(ctx, "client.CreateFirewallDevice(...)", map[string]any{
			"options": createOpts,
		})

		_, err := client.CreateFirewallDevice(ctx, id, createOpts)
		if err != nil {
			return err
		}
	}

	// Clean up remaining devices
	for _, device := range deviceMap {
		tflog.Debug(ctx, "client.DeleteFirewallDevice(...)", map[string]any{
			"device_id": device.ID,
		})

		if err := client.DeleteFirewallDevice(ctx, id, device.ID); err != nil {
			return err
		}
	}

	return nil
}

func refreshDevices(
	ctx context.Context,
	client *linodego.Client,
	firewallID int,
	data *FirewallResourceModel,
	diags *diag.Diagnostics,
	preserveKnown bool,
) {
	devices, err := client.ListFirewallDevices(ctx, firewallID, nil)
	if err != nil {
		diags.AddError(fmt.Sprintf("Failed to Get Devices for Firewall %d", firewallID), err.Error())
		return
	}

	data.flattenDevices(ctx, devices, preserveKnown, diags)
}

func refreshRules(
	ctx context.Context,
	client *linodego.Client,
	firewallID int,
	data *FirewallResourceModel,
	diags *diag.Diagnostics,
	preserveKnown bool,
) {
	rules, err := client.GetFirewallRules(ctx, firewallID)
	if err != nil {
		diags.AddError(fmt.Sprintf("Failed to Get Rules for Firewall %d", firewallID), err.Error())
		return
	}

	data.flattenRules(ctx, rules, preserveKnown, diags)
}

func disableFirewall(
	ctx context.Context,
	firewallID int,
	client *linodego.Client,
	diags *diag.Diagnostics,
) *linodego.Firewall {
	updateOpts := linodego.FirewallUpdateOptions{
		Status: linodego.FirewallDisabled,
	}

	tflog.Debug(ctx, "client.UpdateFirewall(...)", map[string]any{
		"options": updateOpts,
	})

	firewall, err := client.UpdateFirewall(ctx, firewallID, updateOpts)
	if err != nil {
		diags.AddError("Failed to Disable Fireall", err.Error())
	}
	return firewall
}
