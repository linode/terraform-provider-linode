package nbconfig

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
				Name:   "linode_nodebalancer_config",
				IDType: types.StringType,
				Schema: &frameworkResourceSchemaV1,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &frameworkResourceSchemaV0,
			StateUpgrader: upgradeNodeBalancerConfigStateV0toV1,
		},
	}
}

func upgradeNodeBalancerConfigStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	var stateV0 ResourceModelV0
	var stateV1 ResourceModelV1
	resp.Diagnostics.Append(req.State.Get(ctx, &stateV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(stateV1.UpgradeFromV0(ctx, stateV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateV1)...)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var plan ResourceModelV1
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeBalancerID := helper.FrameworkSafeInt64ToInt(
		plan.NodeBalancerID.ValueInt64(), &resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := plan.GetNodeBalancerConfigCreateOptions(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.CreateNodeBalancerConfig(...)", map[string]any{
		"options": createOpts,
	})

	config, err := client.CreateNodeBalancerConfig(ctx, nodeBalancerID, *createOpts)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Create Node Balancer Config", err.Error())
		return
	}

	ctx = tflog.SetField(ctx, "config_id", config.ID)

	resp.Diagnostics.Append(plan.FlattenNodeBalancerConfig(config, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(strconv.Itoa(config.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	client := r.Meta.Client
	var state ResourceModelV1

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	populateLogAttributes(ctx, state)

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	id, nodeBalancerID := getIDs(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := client.GetNodeBalancerConfig(ctx, nodeBalancerID, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf(
					"NodeBalancer Config  %q No Longer Exists",
					state.ID.ValueString(),
				),
				"Removing the NodeBalancer config from the Terraform "+
					"state because it no longer exists.",
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get NodeBalancer Config %d", id),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(state.FlattenNodeBalancerConfig(config, false)...)
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

	var state ResourceModelV1
	var plan ResourceModelV1

	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	id, nodeBalancerID := getIDs(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpts := plan.GetNodeBalancerConfigUpdateOptions(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.UpdateNodeBalancerConfig(...)", map[string]any{
		"options": updateOpts,
	})

	config, err := client.UpdateNodeBalancerConfig(ctx, nodeBalancerID, id, *updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Update the NodeBalancer Config %v", id),
			err.Error(),
		)
		return
	}

	plan.FlattenNodeBalancerConfig(config, true)
	plan.CopyFrom(state, true)

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)
	var state ResourceModelV1
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	id, nodeBalancerID := getIDs(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteNodeBalancerConfig(...)")

	err := client.DeleteNodeBalancerConfig(ctx, nodeBalancerID, id)
	if err != nil {
		if !linodego.IsNotFound(err) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Delete the NodeBalancer Config %d", id),
				err.Error(),
			)
		}
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import "+r.Config.Name)

	helper.ImportStateWithMultipleIDs(
		ctx,
		req,
		resp,
		[]helper.ImportableID{
			{
				Name:          "nodebalancer_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "id",
				TypeConverter: helper.IDTypeConverterString,
			},
		})
}

func populateLogAttributes(ctx context.Context, data ResourceModelV1) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"nodebalancer_id": data.NodeBalancerID.ValueInt64(),
		"id":              data.ID.ValueString(),
	})
}

func getIDs(data ResourceModelV1, diags *diag.Diagnostics) (int, int) {
	id := helper.StringToInt(data.ID.ValueString(), diags)
	nodeBalancerID := helper.FrameworkSafeInt64ToInt(
		data.NodeBalancerID.ValueInt64(), diags,
	)
	return id, nodeBalancerID
}
