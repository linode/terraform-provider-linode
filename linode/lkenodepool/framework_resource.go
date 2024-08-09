package lkenodepool

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
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
				Name:   "linode_lke_node_pool",
				IDType: types.StringType,
				Schema: &resourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_lke_node_pool")
	var data NodePoolModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"cluster_id": data.ClusterID.ValueInt64(),
		"pool_id":    data.ID.ValueString(),
	})

	clusterID, poolID := data.ExtractClusterAndNodePoolIDs(&resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	nodePool, err := client.GetLKENodePool(ctx, clusterID, poolID)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Error reading Linode Node Pool",
				fmt.Sprintf("Removing Linode Node Pool %d in cluster %d, from state because it no longer exists", poolID, clusterID),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading Linode Node Pool %d, cluster %d", poolID, clusterID),
			err.Error(),
		)
		return
	}

	data.FlattenLKENodePool(ctx, nodePool, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "Read linode_lke_node_pool done")
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create linode_lke_node_pool")
	var plan NodePoolModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var createOpts linodego.LKENodePoolCreateOptions

	plan.SetNodePoolCreateOptions(ctx, &createOpts, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := helper.FrameworkSafeInt64ToInt(
		plan.ClusterID.ValueInt64(),
		&resp.Diagnostics,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.CreateLKENodePool(...)", map[string]any{
		"cluster_id": clusterID,
		"options":    createOpts,
	})
	pool, err := client.CreateLKENodePool(ctx, clusterID, createOpts)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Linode Node Pool", err.Error())
		return
	}

	// set cluster ID and pool ID right after pool creation to
	// prevent resource leak when the waiting fails
	AddPoolResource(ctx, pool, resp, plan)

	tflog.Debug(ctx, "waiting for node pool to enter ready status")
	readyPool, err := WaitForNodePoolReady(ctx,
		*client,
		int(r.Meta.Config.EventPollMilliseconds.ValueInt64()),
		clusterID,
		pool.ID,
	)
	if err != nil {
		resp.Diagnostics.AddError("Linode Node Pool is not ready after create", err.Error())
		return
	}

	plan.FlattenLKENodePool(ctx, readyPool, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(strconv.Itoa(readyPool.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Trace(ctx, "Create linode_lke_node_pool done")
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_lke_node_pool")
	var state, plan NodePoolModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateOpts linodego.LKENodePoolUpdateOptions

	plan.SetNodePoolUpdateOptions(ctx, &updateOpts, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID, poolID := state.ExtractClusterAndNodePoolIDs(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.UpdateLKENodePool(...)", map[string]any{
		"cluster_id": clusterID,
		"options":    updateOpts,
	})
	pool, err := client.UpdateLKENodePool(ctx, clusterID, poolID, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError("Error updating a Linode Node Pool", err.Error())
		return
	}

	tflog.Debug(ctx, "waiting for node pool to enter ready status")
	readyPool, err := WaitForNodePoolReady(ctx,
		*client,
		int(r.Meta.Config.EventPollMilliseconds.ValueInt64()),
		clusterID,
		pool.ID,
	)
	if err != nil {
		resp.Diagnostics.AddError("Linode Node Pool is not ready after update", err.Error())
		return
	}

	plan.FlattenLKENodePool(ctx, readyPool, true, &resp.Diagnostics)

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Update linode_lke_node_pool done")
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_lke_node_pool")
	var data NodePoolModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID, poolID := data.ExtractClusterAndNodePoolIDs(&resp.Diagnostics)

	tflog.Debug(ctx, "client.DeleteLKENodePool(...)", map[string]any{
		"cluster_id": clusterID,
		"pool_id":    poolID,
	})
	err := client.DeleteLKENodePool(ctx, clusterID, poolID)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Node Pool does not exist.",
				fmt.Sprintf("Node Pool %v does not exist in cluster %v, removing from state.", poolID, clusterID),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to delete Node Pool",
			err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "Delete linode_lke_node_pool done")
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import linode_lke_node_pool")

	helper.ImportStateWithMultipleIDs(
		ctx,
		req,
		resp,
		[]helper.ImportableID{
			{
				Name:          "cluster_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "id",
				TypeConverter: helper.IDTypeConverterString,
			},
		})
}

func AddPoolResource(
	ctx context.Context, p *linodego.LKENodePool, resp *resource.CreateResponse, plan NodePoolModel,
) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(strconv.Itoa(p.ID)))
	resp.State.SetAttribute(ctx, path.Root("cluster_id"), plan.ClusterID)
}
