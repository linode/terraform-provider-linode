package rdns

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
				Name:   "linode_rdns",
				Schema: &frameworkResourceSchema,
				IDType: types.StringType,
				TimeoutOpts: &timeouts.Opts{
					Update: true,
					Create: true,
				},
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create linode_rdns")

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	updateOpts := linodego.IPAddressUpdateOptions{}

	if !plan.RDNS.IsNull() {
		updateOpts.RDNS = plan.RDNS.ValueStringPointer()
	}

	if !plan.Reserved.IsNull() {
		reserved := plan.Reserved.ValueBool()
		updateOpts.Reserved = reserved
	}

	ip, err := updateIPAddress(
		ctx,
		client,
		plan.Address.ValueString(),
		updateOpts,
		plan.WaitForAvailable.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create/update IP Address", err.Error())
		return
	}

	plan.FlattenInstanceIP(ip, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read linode_rdns")

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	ip, err := client.GetIPAddress(ctx, state.ID.ValueString())
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read IP Address", err.Error())
		return
	}

	state.FlattenInstanceIP(ip, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update linode_rdns")

	var plan, state ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	updateOpts := linodego.IPAddressUpdateOptions{}

	if !plan.RDNS.IsNull() {
		updateOpts.RDNS = plan.RDNS.ValueStringPointer()
	}

	if !plan.Reserved.IsNull() {
		reserved := plan.Reserved.ValueBool()
		updateOpts.Reserved = reserved
	}

	ip, err := updateIPAddress(
		ctx,
		client,
		plan.Address.ValueString(),
		updateOpts,
		plan.WaitForAvailable.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update IP Address", err.Error())
		return
	}

	plan.FlattenInstanceIP(ip, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete linode_rdns")

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS:     nil,
		Reserved: false,
	}

	_, err := client.UpdateIPAddress(ctx, state.Address.ValueString(), updateOpts)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("Failed to delete IP Address reservation", err.Error())
	}
}
