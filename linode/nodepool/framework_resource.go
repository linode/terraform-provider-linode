package nodepool

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/lke"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_nodepool",
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
	tflog.Debug(ctx, "Read linode_nodepool")
	var data NodePoolModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID, poolID := data.ExtractClusterAndNodePoolIDs(&resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetLKENodePool(...)")
	nodePool, err := client.GetLKENodePool(ctx, clusterID, poolID)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Error reading Linode Node Pool",
				fmt.Sprintf("Removing Linode Node Pool %d, cluster %d, from state because it no longer exists", poolID, clusterID),
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

	data.ParseNodePool(ctx, clusterID, nodePool, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "Read linode_nodepool done")
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create linode_nodepool")
	var data NodePoolModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var createOpts linodego.LKENodePoolCreateOptions

	data.SetNodePoolCreateOptions(ctx, &createOpts, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := helper.FrameworkSafeInt64ToInt(
		data.ClusterID.ValueInt64(),
		&resp.Diagnostics,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.CreateLKENodePool(...)")
	pool, err := client.CreateLKENodePool(ctx, clusterID, createOpts)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Linode Node Pool", err.Error())
		return
	}

	tflog.Debug(ctx, "waiting for node pool to enter ready status")
	if err := lke.WaitForNodePoolReady(
		ctx,
		*client,
		int(r.Meta.Config.EventPollMilliseconds.ValueInt64()),
		clusterID,
		pool.ID,
	); err != nil {
		resp.Diagnostics.AddWarning("Linode Node Pool is not ready after create", err.Error())
	}

	data.ParseNodePool(ctx, clusterID, pool, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "Create linode_nodepool done")
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_nodepool")
	var plan NodePoolModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var updateOpts linodego.LKENodePoolUpdateOptions

	plan.SetNodePoolUpdateOptions(ctx, &updateOpts, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID, poolID := plan.ExtractClusterAndNodePoolIDs(&resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.UpdateLKENodePool(...)")
	pool, err := client.UpdateLKENodePool(ctx, clusterID, poolID, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError("Error updating a Linode Node Pool", err.Error())
		return
	}

	tflog.Debug(ctx, "waiting for node pool to enter ready status")
	if err := lke.WaitForNodePoolReady(
		ctx,
		*client,
		int(r.Meta.Config.EventPollMilliseconds.ValueInt64()),
		clusterID,
		pool.ID,
	); err != nil {
		resp.Diagnostics.AddWarning("Linode Node Pool is not ready after update", err.Error())
	}

	plan.ParseNodePool(ctx, clusterID, pool, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Update linode_nodepool done")
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_nodepool")
	var data NodePoolModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID, poolID := data.ExtractClusterAndNodePoolIDs(&resp.Diagnostics)

	tflog.Trace(ctx, "client.DeleteLKENodePool(...)")
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
	tflog.Trace(ctx, "Delete linode_nodepool done")
}
