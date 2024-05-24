package lke

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_lke_cluster",
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
	tflog.Debug(ctx, "Read data.linode_lke_cluster")

	client := r.Meta.Client

	var data LKEDataModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cluster_id", clusterId)

	cluster, err := client.GetLKECluster(ctx, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get LKE cluster %d", clusterId),
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "client.ListLKENodePools(...)")

	pools, err := client.ListLKENodePools(ctx, clusterId, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get pools for LKE cluster %d", clusterId),
			err.Error(),
		)
		return
	}

	kubeconfig, err := client.GetLKEClusterKubeconfig(ctx, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get kubeconfig for LKE cluster %d", clusterId),
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "client.ListLKEClusterAPIEndpoints(...)")

	endpoints, err := client.ListLKEClusterAPIEndpoints(ctx, clusterId, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get API endpoints for LKE cluster %d", clusterId),
			err.Error(),
		)
		return
	}

	dashboard, err := client.GetLKEClusterDashboard(ctx, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get dashboard URL for LKE cluster %d", clusterId),
			err.Error(),
		)
		return
	}

	acl, err := client.GetLKEClusterControlPlaneACL(ctx, clusterId)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok &&
			lerr.Code == 400 && strings.Contains(lerr.Message, "Cluster does not support Control Plane ACL") {
			// The cluster does not have a Gateway. Nothing to do here.
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to get Control Plane ACL for LKE cluster %d", clusterId),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(data.parseLKEAttributes(ctx, cluster, pools, kubeconfig, endpoints, dashboard, acl)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
