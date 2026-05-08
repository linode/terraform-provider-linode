package reservedip

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
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

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create "+r.Config.Name)
	var plan ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"region": plan.Region.ValueString(),
	})

	createOpts := linodego.ReserveIPOptions{
		Region: plan.Region.ValueString(),
	}

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &createOpts.Tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, "client.ReserveIPAddress(...)", map[string]any{
		"options": createOpts,
	})

	client := r.Meta.Client

	ip, err := client.ReserveIPAddress(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Reserve IP Address",
			fmt.Sprintf("failed to reserve IP in region %s: %s", plan.Region.ValueString(), err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(plan.flatten(ctx, *ip, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(ip.Address)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	address := state.ID.ValueString()

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"address": address,
	})

	client := r.Meta.Client

	ip, err := client.GetReservedIPAddress(ctx, address)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Reserved IP No Longer Exists",
				fmt.Sprintf("Removing reserved IP %s from state because it no longer exists", address),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Read Reserved IP",
			fmt.Sprintf("error reading reserved IP %s: %s", address, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(state.flatten(ctx, *ip, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)
	var plan ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	address := plan.ID.ValueString()

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"address": address,
	})

	var newTags []string
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &newTags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	updateOpts := linodego.UpdateReservedIPOptions{
		Tags: newTags,
	}

	tflog.Debug(ctx, "client.UpdateReservedIPAddress(...)", map[string]any{
		"options": updateOpts,
	})

	client := r.Meta.Client

	ip, err := client.UpdateReservedIPAddress(ctx, address, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Reserved IP",
			fmt.Sprintf("failed to update reserved IP %s: %s", address, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(plan.flatten(ctx, *ip, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	address := state.ID.ValueString()

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"address": address,
	})

	tflog.Debug(ctx, "client.DeleteReservedIPAddress(...)")

	client := r.Meta.Client

	if err := client.DeleteReservedIPAddress(ctx, address); err != nil {
		if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				"Failed to Delete Reserved IP",
				fmt.Sprintf("failed to delete reserved IP %s: %s", address, err.Error()),
			)
		}
	}
}
