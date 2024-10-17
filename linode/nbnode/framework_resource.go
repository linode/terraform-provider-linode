package nbnode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_nodebalancer_node",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func AddNodeResource(ctx context.Context, node linodego.NodeBalancerNode, resp *resource.CreateResponse, plan ResourceModel) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(strconv.Itoa(node.ID)))
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

	nodeBalancerID, configID, createOpts := plan.GetCreateParameters(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.CreateNodeBalancerNode(...)", map[string]any{
		"options": createOpts,
	})
	node, err := client.CreateNodeBalancerNode(ctx, nodeBalancerID, configID, createOpts)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Create a Linode NodeBalancerNode", err.Error())
		return
	}

	// Add resource to TF states earlier to prevent
	// dangling resources (resources created but not managed by TF)
	AddNodeResource(ctx, *node, resp, plan)

	ctx = tflog.SetField(ctx, "node_id", node.ID)

	plan.FlattenNodeBalancerNode(node, true)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(strconv.Itoa(node.ID))

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

	ctx = populateLogAttributes(ctx, state)

	// TODO: cleanup when Crossplane fixes it
	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	id, nodeBalancerID, configID := state.GetIDs(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	node, err := client.GetNodeBalancerNode(ctx, nodeBalancerID, configID, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"The NodeBalancer Node No Longer Exists",
				fmt.Sprintf(
					"Removing Linode Token with ID %v from state because it no longer exists", id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Refresh the NodeBalancer Node",
			fmt.Sprintf(
				"Error finding the specified Linode Token: %s",
				err.Error(),
			),
		)
		return
	}

	state.FlattenNodeBalancerNode(node, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	ctx = populateLogAttributes(ctx, state)
	client := r.Meta.Client

	id, nodeBalancerID, configID := plan.GetIDs(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpts := plan.GetUpdateOptions(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.UpdateNodeBalancerNode(...)", map[string]any{
		"options": updateOpts,
	})
	node, err := client.UpdateNodeBalancerNode(ctx, nodeBalancerID, configID, id, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Failed to Update Linode NodeBalancer %d Config %d Node %d",
				nodeBalancerID, configID, id,
			),
			err.Error(),
		)
		return
	}

	plan.FlattenNodeBalancerNode(node, true)

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

	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	client := r.Meta.Client

	id, nodeBalancerID, configID := state.GetIDs(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.DeleteNodeBalancerNode(...)")
	err := client.DeleteNodeBalancerNode(ctx, nodeBalancerID, configID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Failed to Delete Linode NodeBalancer %d Config %d Node %d",
				nodeBalancerID, configID, id,
			),
			err.Error(),
		)
		return
	}
}

func populateLogAttributes(ctx context.Context, model ResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"nodebalancer_id": model.NodeBalancerID.ValueInt64(),
		"config_id":       model.ConfigID.ValueInt64(),
		"id":              model.ID.ValueString(),
	})
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
				Name:          "config_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "id",
				TypeConverter: helper.IDTypeConverterString,
			},
		},
	)
}
