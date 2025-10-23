package consumerimagesharegrouptoken

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
				Name:   "linode_consumer_image_share_group_token",
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

	createOpts := linodego.ImageShareGroupCreateTokenOptions{
		ValidForShareGroupUUID: plan.ValidForShareGroupUUID.ValueString(),
		Label:                  plan.Label.ValueStringPointer(),
	}

	tflog.Debug(ctx, "client.ImageShareGroupCreateToken(...)", map[string]any{
		"options": createOpts,
	})

	token, err := client.ImageShareGroupCreateToken(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Image Share Group Token.",
			err.Error(),
		)
		return
	}

	plan.FlattenImageShareGroupCreateToken(token)
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

	tokenUUID := state.TokenUUID.ValueString()

	token, err := client.ImageShareGroupGetToken(ctx, tokenUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read Image Share Group Token.",
			err.Error(),
		)
		return
	}

	state.FlattenImageShareGroupToken(token, true)
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

	tokenUUID := state.TokenUUID.ValueString()

	var updateOpts linodego.ImageShareGroupUpdateTokenOptions
	shouldUpdate := false

	if !state.Label.Equal(plan.Label) {
		shouldUpdate = true
		updateOpts.Label = plan.Label.ValueString()
	}

	if shouldUpdate {
		tflog.Debug(ctx, "client.ImageShareGroupUpdateToken(...)", map[string]any{
			"options": updateOpts,
		})

		token, err := client.ImageShareGroupUpdateToken(ctx, tokenUUID, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update Image Share Group Token (%d).", tokenUUID),
				err.Error(),
			)
			return
		}

		plan.FlattenImageShareGroupToken(token, false)
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

	tokenUUID := state.TokenUUID.ValueString()

	client := r.Meta.Client

	err := client.ImageShareGroupRemoveToken(ctx, tokenUUID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Delete Image Share Group Token.", err.Error())
		return
	}
}
