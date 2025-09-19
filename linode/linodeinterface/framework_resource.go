package linodeinterface

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
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
				Name:   "linode_interface",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func AddInterfaceResource(ctx context.Context, i linodego.LinodeInterface, resp *resource.CreateResponse, plan LinodeInterfaceModel) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(strconv.Itoa(i.ID)))
	resp.State.SetAttribute(ctx, path.Root("linode_id"), plan.LinodeID)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var plan LinodeInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	helper.SetLogFieldBulk(ctx, map[string]any{"linode_id": plan.LinodeID})
	client := r.Meta.Client

	opts, linodeID := plan.GetCreateOptions(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	i, err := client.CreateInterface(ctx, linodeID, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Create a Network Interface for Linode Instance %d", linodeID),
			err.Error(),
		)
		return
	}

	// Add resource to TF states earlier to prevent dangling resources
	// (resources created but not managed by TF) when a later step fails.
	AddInterfaceResource(ctx, *i, resp, plan)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(strconv.Itoa(i.ID))

	plan.FlattenInterface(ctx, *i, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var state LinodeInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	linodeID, id := state.GetIDs(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	linodeInterface, err := client.GetInterface(ctx, linodeID, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Linode Interface No Longer Exists",
				fmt.Sprintf("Linode Interface %v does not exist, removing it from state.", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Get Linode Interface %v", id),
			err.Error(),
		)
		return
	}

	state.FlattenInterface(ctx, *linodeInterface, false, &resp.Diagnostics)
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

	var plan, state LinodeInterfaceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	ctx = populateLogAttributes(ctx, state)

	linodeID, id := state.GetIDs(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	opts := plan.GetUpdateOptions(ctx, state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	i, err := client.UpdateInterface(ctx, linodeID, id, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Update Linode Interface %d", id),
			err.Error(),
		)
		return
	}

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	plan.FlattenInterface(ctx, *i, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	// plan.CopyFrom(ctx, state, &resp.Diagnostics, true)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	var state LinodeInterfaceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Meta.Client

	linodeID, id := state.GetIDs(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := client.DeleteInterface(ctx, linodeID, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Linode Interface No Longer Exists",
				fmt.Sprintf("Linode Interface %v does not exist, removing it from state.", id),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to Delete Linode Interface",
			err.Error(),
		)
		return
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import "+r.Config.Name)
	helper.ImportStateWithMultipleIDs(
		ctx, req, resp,
		[]helper.ImportableID{
			{
				Name:          "id",
				TypeConverter: helper.IDTypeConverterString,
			},
			{
				Name:          "linode_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
		},
	)
}

func populateLogAttributes(ctx context.Context, model LinodeInterfaceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"linode_id": model.LinodeID.ValueInt64(),
		"id":        model.ID.ValueString(),
	})
}
