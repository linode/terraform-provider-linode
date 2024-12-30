package networkingip

import (
	"context"
	"fmt"

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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create linode_networking_ip")
	var plan NetworkingIPModel

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

	plan.FlattenIPAddress(ctx, ip, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read linode_networking_ip")
	var state NetworkingIPModel

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

	state.FlattenIPAddress(ctx, foundIP, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update linode_networking_ip")
	var plan, state NetworkingIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	var reservedValue *bool
	if plan.Reserved != state.Reserved {
		value := plan.Reserved.ValueBoolPointer()
		reservedValue = value
	}

	updateOpts := linodego.IPAddressUpdateOptionsV2{
		Reserved: reservedValue,
	}

	ip, err := client.UpdateIPAddressV2(ctx, state.Address.ValueString(), updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update IP Address",
			fmt.Sprintf("Could not update reserved status of IP address: %s", err),
		)
		return
	}

	plan.FlattenIPAddress(ctx, ip, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete linode_networking_ip")
	var state NetworkingIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	// Regular assigned ephemeral IP address
	if !state.Reserved.ValueBool() {
		// This is a regular ephemeral IP address
		linodeID := helper.FrameworkSafeInt64ToInt(state.LinodeID.ValueInt64(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		// Proceed with deleting the IP if it's not the only one
		err := client.DeleteInstanceIPAddress(ctx, linodeID, state.Address.ValueString())
		if err != nil {
			if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
				resp.Diagnostics.AddError(
					"Failed to Delete IP",
					fmt.Sprintf(
						"failed to delete instance (%d) ip (%s): %s",
						linodeID, state.Address.ValueString(), err.Error(),
					),
				)
			}
		}
	} else {
		// Reserved IP address
		// If the IP is currently assigned (reserved but used)
		if state.LinodeID.ValueInt64() != 0 {
			// It's an assigned reserved IP, we can delete it regardless of being the only IP
			linodeID := helper.FrameworkSafeInt64ToInt(state.LinodeID.ValueInt64(), &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}

			// Delete the reserved IP (this will turn it into an ephemeral IP if it's the only IP)
			err := client.DeleteReservedIPAddress(ctx, state.Address.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to Delete Assigned Reserved IP",
					fmt.Sprintf(
						"failed to delete assigned reserved ip (%s) from linode (%d): %s",
						state.Address.ValueString(), linodeID, err.Error(),
					),
				)
			}
			return
		} else {
			// Reserved IP (unassigned) that needs to be deleted
			// If it's a reserved IP but it is not assigned to a Linode, proceed with deletion
			err := client.DeleteReservedIPAddress(ctx, state.Address.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to Delete Reserved IP",
					fmt.Sprintf(
						"failed to delete reserved ip (%s): %s",
						state.Address.ValueString(), err.Error(),
					),
				)
			}
		}
	}
}
