package networkingip

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
				Name:   "linode_networking_ip",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// reconcileIPReassignment marks changes to linode_id as RequiresReplace
	// if the operation cannot be completed using a resource update.
	reconcileIPReassignment := func() {
		var planReserved types.Bool
		var planLinodeID, stateLinodeID types.Int64

		resp.Diagnostics.Append(
			req.Plan.GetAttribute(ctx, path.Root("reserved"), &planReserved)...,
		)

		resp.Diagnostics.Append(
			req.Plan.GetAttribute(ctx, path.Root("linode_id"), &planLinodeID)...,
		)

		resp.Diagnostics.Append(
			req.State.GetAttribute(ctx, path.Root("linode_id"), &stateLinodeID)...,
		)

		if resp.Diagnostics.HasError() {
			return
		}

		if planReserved.ValueBool() || planLinodeID.IsUnknown() || planLinodeID.Equal(stateLinodeID) {
			// Nothing to do here
			return
		}

		resp.RequiresReplace.Append(path.Root("linode_id"))
		resp.Diagnostics.AddAttributeWarning(
			path.Root("linode_id"),
			"Resource must be recreated to update assigned linode_id",
			"Ephemeral IPs cannot be reassigned through a resource update.",
		)
	}

	reconcileIPReassignment()
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create linode_networking_ip")
	var plan ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	createOpts := linodego.AllocateReserveIPOptions{
		Type:   plan.Type.ValueString(),
		Public: plan.Public.ValueBool(),
	}

	if !plan.LinodeID.IsNull() {
		createOpts.LinodeID = int(plan.LinodeID.ValueInt64())
	}
	if !plan.Reserved.IsNull() {
		createOpts.Reserved = plan.Reserved.ValueBool()
	}
	if !plan.Region.IsNull() {
		createOpts.Region = plan.Region.ValueString()
	}

	ip, err := client.AllocateReserveIP(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating IP Address",
			fmt.Sprintf("Could not create IP address: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(plan.FlattenIPAddress(ip, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(ip.Address)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read linode_networking_ip")
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	// Use ListIPAddresses as a workaround to retrieve the specific private IP address
	// since GetIPAddress doesnt retrieve private IP addresses
	ips, err := client.ListIPAddresses(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading IP Addresses",
			fmt.Sprintf("Could not list IP addresses: %s", err),
		)
		return
	}

	var foundIP *linodego.InstanceIP
	for _, ip := range ips {
		if ip.Address == state.ID.ValueString() {
			foundIP = &ip
			break
		}
	}

	if foundIP == nil {
		// IP address not found; remove the resource from the state
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(state.FlattenIPAddress(foundIP, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update linode_networking_ip")
	var plan, state ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	// Updates to the reserved field need to happen before any reassignments occur
	if !plan.Reserved.Equal(state.Reserved) {
		updateOpts := linodego.IPAddressUpdateOptionsV2{
			Reserved: plan.Reserved.ValueBoolPointer(),
		}

		tflog.Info(ctx, "Updating IP address to reconcile reserved status")
		tflog.Debug(ctx, "client.UpdateIPAddressV2(...)", map[string]any{
			"options": updateOpts,
		})

		if _, err := client.UpdateIPAddressV2(ctx, state.Address.ValueString(), updateOpts); err != nil {
			resp.Diagnostics.AddError(
				"Failed to Update IP Address",
				fmt.Sprintf("Could not update reserved status of IP address: %s", err),
			)
			return
		}
	}

	// Reconcile the assignment of the IP
	resp.Diagnostics.Append(reconcileIPAssignments(ctx, client, &plan, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.CopyFrom(state, true)

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete linode_networking_ip")
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	if !state.Reserved.ValueBool() {
		// This is a regular ephemeral IP address
		linodeID := helper.FrameworkSafeInt64ToInt(state.LinodeID.ValueInt64(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		// Proceed with deleting the IP if it's not the only one
		tflog.Debug(ctx, "client.DeleteInstanceIPAddress(...)", map[string]any{
			"linode_id": linodeID,
			"address":   state.Address.ValueString(),
		})
		err := client.DeleteInstanceIPAddress(ctx, linodeID, state.Address.ValueString())
		if err != nil {
			if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
				resp.Diagnostics.AddError(
					"Failed to Delete Ephemeral IP",
					err.Error(),
				)
				return
			}
		}
	} else {
		// Reserved IP (unassigned) that needs to be deleted
		// If it's a reserved IP is not assigned to a Linode, proceed with deletion
		tflog.Debug(ctx, "client.DeleteReservedIPAddress(...)", map[string]any{
			"address": state.Address.ValueString(),
		})
		err := client.DeleteReservedIPAddress(ctx, state.Address.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to delete Reserved IP",
				err.Error(),
			)
			return
		}
	}
}

// reconcileIPAssignments reconciles the assignment of the IP managed by the given plan and state.
// This is intended to be called during update operations and is necessary due to the complexity
// of the reserved assignment validation.
func reconcileIPAssignments(
	ctx context.Context,
	client *linodego.Client,
	plan, state *ResourceModel,
) (d diag.Diagnostics) {
	if plan.LinodeID.IsUnknown() || plan.LinodeID.Equal(state.LinodeID) {
		// Nothing to do here
		return
	}

	if !state.LinodeID.IsNull() {
		// First, we should unassign the IP address from its current Linode if necessary.
		// NOTE: Without swapping, IP addresses must be unassigned before
		// being reassigned to a new Linode.
		stateLinodeID := helper.FrameworkSafeInt64ToInt(state.LinodeID.ValueInt64(), &d)
		if d.HasError() {
			return
		}

		if err := client.DeleteInstanceIPAddress(ctx, stateLinodeID, state.Address.ValueString()); err != nil {
			d.AddError(
				"Failed to unassign reserved IP from Linode",
				err.Error(),
			)
			return
		}

		state.LinodeID = types.Int64Null()
	}

	if !plan.LinodeID.IsNull() {
		// Asign the IP to a new Linode if necessary
		planLinodeID := helper.FrameworkSafeInt64ToInt(plan.LinodeID.ValueInt64(), &d)
		if d.HasError() {
			return
		}

		ip, err := client.AssignInstanceReservedIP(
			ctx,
			planLinodeID,
			linodego.InstanceReserveIPOptions{
				Type:    state.Type.ValueString(),
				Public:  state.Public.ValueBool(),
				Address: state.Address.ValueString(),
			},
		)
		if err != nil {
			d.AddError(
				"Failed to assign reserved IP to Linode",
				err.Error(),
			)
			return
		}

		state.LinodeID = types.Int64Value(int64(ip.LinodeID))
	}

	return
}
