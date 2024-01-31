package firewalldevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_firewall_device",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create linode_firewall_device")

	var plan FirewallDeviceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entityID := helper.FrameworkSafeInt64ToInt(
		plan.EntityID.ValueInt64(),
		&resp.Diagnostics,
	)
	firewallID := helper.FrameworkSafeInt64ToInt(
		plan.FirewallID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.FirewallDeviceCreateOptions{
		ID:   entityID,
		Type: linodego.FirewallDeviceType(plan.EntityType.ValueString()),
	}

	tflog.Debug(ctx, "client.CreateFirewallDevice(...)", map[string]any{
		"firewall_id": firewallID,
		"options":     createOpts,
	})
	device, err := client.CreateFirewallDevice(ctx, firewallID, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating a Linode Firewall Device",
			err.Error(),
		)
		return
	}

	plan.FlattenFirewallDevice(device, true)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_firewall_device")

	var state FirewallDeviceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"firewall_id": state.FirewallID.ValueInt64(),
		"device_id":   state.ID.ValueInt64(),
	})

	client := r.Meta.Client

	id := helper.FrameworkSafeInt64ToInt(
		state.ID.ValueInt64(),
		&resp.Diagnostics,
	)
	firewallID := helper.FrameworkSafeInt64ToInt(
		state.FirewallID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetFirewallDevice(...)")

	device, err := client.GetFirewallDevice(ctx, firewallID, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Firewall Device No Longer Exists",
				fmt.Sprintf(
					"Removing firewall device %d from state because it no longer exists",
					state.ID.ValueInt64(),
				),
			)
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				"Error finding the specified Linode Firewall Device",
				err.Error(),
			)
		}
		return
	}

	state.FlattenFirewallDevice(device, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_firewall_device")
	resp.Diagnostics.AddWarning(
		"Unintended Calling to Update Function",
		"The Update function of 'linode_firewall_device' should never be "+
			"invoked by design. This function has been redundantly implemented "+
			"for improved reliability. Please consider reporting this as a bug "+
			"to the provider developers.",
	)

	var state, plan FirewallDeviceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.CopyFrom(state, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_firewall_device")
	var state FirewallDeviceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	client := r.Meta.Client

	id := helper.FrameworkSafeInt64ToInt(
		state.ID.ValueInt64(),
		&resp.Diagnostics,
	)
	firewallID := helper.FrameworkSafeInt64ToInt(
		state.FirewallID.ValueInt64(),
		&resp.Diagnostics,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteFirewallDevice(...)", map[string]any{
		"firewall_id": firewallID,
		"device_id":   id,
	})

	err := client.DeleteFirewallDevice(ctx, firewallID, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf(
					"Attempted to Delete Firewall Device %d But Resource Not Found",
					id,
				),
				err.Error(),
			)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error deleting Linode Firewall Device %d", id),
				err.Error(),
			)
		}
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import linode_firewall_device")

	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: firewall_id,device_id. Got: %q", req.ID),
		)
		return
	}

	firewallID := helper.StringToInt64(idParts[0], &resp.Diagnostics)
	deviceID := helper.StringToInt64(idParts[1], &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), deviceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("firewall_id"), firewallID)...)
}
