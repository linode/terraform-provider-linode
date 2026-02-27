package lkenodepool

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_lke_node_pool",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+r.Config.Name)

	var data nodePoolDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := helper.FrameworkSafeInt64ToInt(data.ClusterID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "cluster_id", clusterID)

	nodePoolID := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "nodepool_id", nodePoolID)

	client := r.Meta.Client
	lkeNodePool, err := client.GetLKENodePool(ctx, clusterID, nodePoolID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get LKE Node Pool %d for LKE Cluster %d", nodePoolID, clusterID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseLKENodePool(ctx, lkeNodePool)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
