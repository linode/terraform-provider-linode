package networkreservedips

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
				Name:   "linode_reserved_ip",
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
	tflog.Debug(ctx, "Starting Create for linode_reserved_ip")
	var plan ReservedIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting plan", map[string]interface{}{
			"error": resp.Diagnostics.Errors(),
		})
		return
	}

	ctx = populateLogAttributes(ctx, &plan)

	client := r.Meta.Client
	reserveIP, err := client.ReserveIPAddress(ctx, linodego.ReserveIPOptions{
		Region: plan.Region.ValueString(),
	})
	if err != nil {
		tflog.Error(ctx, "Failed to reserve IP address", map[string]interface{}{
			"error": err.Error(),
		})
		resp.Diagnostics.AddError(
			"Failed to reserve IP address",
			err.Error(),
		)
		return
	}

	if reserveIP == nil {
		tflog.Error(ctx, "Received nil pointer for reserved IP")
		resp.Diagnostics.AddError("nil Pointer", "received nil pointer of the reserved ip")
		return
	}

	tflog.Debug(ctx, "Successfully reserved IP address", map[string]interface{}{
		"address": reserveIP.Address,
		"region":  reserveIP.Region,
	})

	plan.ID = types.StringValue(reserveIP.Address)
	tflog.Debug(ctx, "Setting ID for reserved IP", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	diags := plan.FlattenReservedIP(ctx, *reserveIP, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error flattening reserved IP", map[string]interface{}{
			"error": resp.Diagnostics.Errors(),
		})
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error setting state for reserved IP", map[string]interface{}{
			"error": resp.Diagnostics.Errors(),
		})
	} else {
		tflog.Debug(ctx, "Successfully set state for reserved IP", map[string]interface{}{
			"id": plan.ID.ValueString(),
		})
	}
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read linode_reserved_ip")
	var state ReservedIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	ctx = populateLogAttributes(ctx, &state)

	client := r.Meta.Client
	address := state.Address.ValueString()

	reservedIP, err := client.GetReservedIPAddress(ctx, address)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Reserved IP No Longer Exists",
				fmt.Sprintf(
					"Removing reserved IP %s from state because it no longer exists",
					state.ID.ValueString(),
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Refresh the Reserved IP",
			fmt.Sprintf(
				"Error finding the specified Reserved IP: %s",
				err.Error(),
			),
		)
		return
	}

	if reservedIP == nil {
		resp.Diagnostics.AddError("nil Pointer", "received nil pointer of the reserved ip")
		return
	}

	resp.Diagnostics.Append(state.FlattenReservedIP(ctx, *reservedIP, false)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Reserved IPs cannot be updated, so this method is left empty
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ReservedIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	address := state.Address.ValueString()

	tflog.Debug(ctx, "client.DeleteReservedIPAddress(...)")
	if err := client.DeleteReservedIPAddress(ctx, address); err != nil {
		if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				"Failed to Delete Reserved IP",
				fmt.Sprintf(
					"failed to delete reserved ip (%s): %s",
					address, err.Error(),
				),
			)
		}
	}
}

func populateLogAttributes(ctx context.Context, data *ReservedIPModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"region":  data.Region.ValueString(),
		"address": data.ID.ValueString(),
	})
}
