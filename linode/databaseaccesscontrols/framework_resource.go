package databaseaccesscontrols

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			"linode_database_access_controls",
			frameworkResourceSchema,
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
	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbID := int(data.DatabaseID.ValueInt64())
	dbType := data.DatabaseType.ValueString()

	if err := updateDBAllowListByEngine(
		ctx,
		client,
		dbType,
		dbID,
		helper.FrameworkSliceToString(data.AllowList),
	); err != nil {
		resp.Diagnostics.AddError(
			"Failed to set DB allow list",
			err.Error(),
		)
		return
	}

	data.ID = types.StringValue(formatID(dbID, dbType))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client := r.Meta.Client

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbID, dbType, err := parseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to parse implicit DB ID %s", data.ID),
			err.Error(),
		)
		return
	}

	allowList, err := getDBAllowListByEngine(ctx, client, dbType, dbID)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Removing allow_list from state because it no longer exists",
				err.Error(),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get DB with ID %d and type %s", dbID, dbType),
			err.Error(),
		)
		return
	}

	data.DatabaseID = types.Int64Value(int64(dbID))
	data.DatabaseType = types.StringValue(dbType)
	data.AllowList = helper.StringSliceToFramework(allowList)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if err := updateDBAllowListByEngine(
		ctx,
		r.Meta.Client,
		state.DatabaseType.ValueString(),
		int(state.DatabaseID.ValueInt64()),
		helper.FrameworkSliceToString(plan.AllowList),
	); err != nil {
		resp.Diagnostics.AddError(
			"Failed to set DB allow list",
			err.Error(),
		)
		return
	}

	state.AllowList = plan.AllowList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := updateDBAllowListByEngine(
		ctx,
		r.Meta.Client,
		data.DatabaseType.ValueString(),
		int(data.DatabaseID.ValueInt64()),
		[]string{},
	); err != nil {
		resp.Diagnostics.AddError(
			"Failed to set DB allow list",
			err.Error(),
		)
		return
	}
}
