package reservedip

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
	var data ReservedIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, &data)

	client := r.Meta.Client
	reserveIP, err := client.ReserveIPAddress(ctx, linodego.ReserveIPOptions{
		Region: data.Region.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to reserve IP address",
			err.Error(),
		)
		return
	}

	if reserveIP == nil {
		resp.Diagnostics.AddError("nil Pointer", "received nil pointer of the reserved ip")
		return
	}

	data.ID = types.StringValue(reserveIP.Address)
	tflog.Debug(ctx, "Setting ID for reserved IP", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	diags := data.FlattenReservedIP(ctx, *reserveIP, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read linode_reserved_ip")
	var data ReservedIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	ctx = populateLogAttributes(ctx, &data)

	client := r.Meta.Client
	address := data.ID.ValueString()

	reservedIP, err := client.GetReservedIPAddress(ctx, address)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Reserved IP No Longer Exists",
				fmt.Sprintf(
					"Removing reserved IP %s from state because it no longer exists",
					data.ID.ValueString(),
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

	resp.Diagnostics.Append(data.FlattenReservedIP(ctx, *reservedIP, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
