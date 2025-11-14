package producerimagesharegroupmember

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
				Name:   "linode_producer_image_share_group_member",
				IDType: types.Int64Type,
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
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.ImageShareGroupAddMemberOptions{
		Token: plan.Token.ValueString(),
		Label: plan.Label.ValueString(),
	}

	tflog.Debug(ctx, "client.ImageShareGroupAddMember(...)", map[string]any{
		"options": createOpts,
	})

	shareGroupID := helper.FrameworkSafeInt64ToInt(plan.ShareGroupID.ValueInt64(), &resp.Diagnostics)

	member, err := client.ImageShareGroupAddMember(ctx, shareGroupID, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Image Share Group Member.",
			err.Error(),
		)
		return
	}

	plan.FlattenImageShareGroupMember(member, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	client := r.Meta.Client

	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shareGroupID := helper.FrameworkSafeInt64ToInt(state.ShareGroupID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tokenUUID := state.TokenUUID.ValueString()

	member, err := client.ImageShareGroupGetMember(ctx, shareGroupID, tokenUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read Image Share Group Member.",
			err.Error(),
		)
		return
	}

	state.FlattenImageShareGroupMember(member, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	client := r.Meta.Client

	var plan ResourceModel
	var state ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shareGroupID := helper.FrameworkSafeInt64ToInt(state.ShareGroupID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tokenUUID := state.TokenUUID.ValueString()

	var updateOpts linodego.ImageShareGroupUpdateMemberOptions
	shouldUpdate := false

	if !state.Label.Equal(plan.Label) {
		shouldUpdate = true
		updateOpts.Label = plan.Label.ValueString()
	}

	if shouldUpdate {
		tflog.Debug(ctx, "client.ImageShareGroupUpdateMember(...)", map[string]any{
			"options": updateOpts,
		})

		member, err := client.ImageShareGroupUpdateMember(ctx, shareGroupID, tokenUUID, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update Image Share Group Member (%d).", shareGroupID),
				err.Error(),
			)
			return
		}

		plan.FlattenImageShareGroupMember(member, false)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	plan.CopyFrom(state, true)

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

	shareGroupID := helper.FrameworkSafeInt64ToInt(state.ShareGroupID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tokenUUID := state.TokenUUID.ValueString()

	client := r.Meta.Client

	err := client.ImageShareGroupRemoveMember(ctx, shareGroupID, tokenUUID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Delete Image Share Group Member.", err.Error())
		return
	}
}
