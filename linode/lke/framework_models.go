package lke

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// LKEDataModel describes the Terraform resource data model to match the
// data source schema.
type LKEDataModel struct {
	// LKE Cluster
	ID           types.Int64       `tfsdk:"id"`
	Created      types.String      `tfsdk:"created"`
	Updated      types.String      `tfsdk:"updated"`
	Label        types.String      `tfsdk:"label"`
	Region       types.String      `tfsdk:"region"`
	Status       types.String      `tfsdk:"status"`
	K8sVersion   types.String      `tfsdk:"k8s_version"`
	Tags         types.Set         `tfsdk:"tags"`
	ControlPlane []LKEControlPlane `tfsdk:"control_plane"`

	// LKE Node Pools
	Pools []LKENodePool `tfsdk:"pools"`

	// LKE Cluster Kubeconfig
	Kubeconfig types.String `tfsdk:"kubeconfig"`

	// LKE Cluster API endpoints
	APIEndpoints types.List `tfsdk:"api_endpoints"`

	// LKE Cluster Dashboard
	DashboardURL types.String `tfsdk:"dashboard_url"`
}

type LKEControlPlane struct {
	HighAvailability types.Bool `tfsdk:"high_availability"`
}

type LKENodePool struct {
	ID         types.Int64             `tfsdk:"id"`
	Count      types.Int64             `tfsdk:"count"`
	Type       types.String            `tfsdk:"type"`
	Tags       types.List              `tfsdk:"tags"`
	Disks      []LKENodePoolDisk       `tfsdk:"disks"`
	Nodes      []LKENodePoolNode       `tfsdk:"nodes"`
	Autoscaler []LKENodePoolAutoscaler `tfsdk:"autoscaler"`
}

type LKENodePoolDisk struct {
	Size types.Int64  `tfsdk:"size"`
	Type types.String `tfsdk:"type"`
}

type LKENodePoolNode struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	Status     types.String `tfsdk:"status"`
}

type LKENodePoolAutoscaler struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Min     types.Int64 `tfsdk:"min"`
	Max     types.Int64 `tfsdk:"max"`
}

func (data *LKEDataModel) parseLKEAttributes(
	ctx context.Context,
	cluster *linodego.LKECluster,
	pools []linodego.LKENodePool,
	kubeconfig *linodego.LKEClusterKubeconfig,
	endpoints []linodego.LKEClusterAPIEndpoint,
	dashboard *linodego.LKEClusterDashboard,
) diag.Diagnostics {
	data.Created = types.StringValue(cluster.Created.Format(helper.TIME_FORMAT))
	data.Updated = types.StringValue(cluster.Updated.Format(helper.TIME_FORMAT))
	data.Label = types.StringValue(cluster.Label)
	data.Region = types.StringValue(cluster.Region)
	data.Status = types.StringValue(string(cluster.Status))
	data.K8sVersion = types.StringValue(cluster.K8sVersion)

	tags, diags := types.SetValueFrom(ctx, types.StringType, cluster.Tags)
	if diags != nil {
		return diags
	}
	data.Tags = tags

	data.ControlPlane = []LKEControlPlane{ParseControlPlane(cluster.ControlPlane)}

	parseLKEPools := func() ([]LKENodePool, diag.Diagnostics) {
		lkePools := make([]LKENodePool, len(pools))

		for i, p := range pools {
			var pool LKENodePool
			pool.ID = types.Int64Value(int64(p.ID))
			pool.Count = types.Int64Value(int64(p.Count))
			pool.Type = types.StringValue(p.Type)
			tags, diags := types.ListValueFrom(ctx, types.StringType, p.Tags)
			if diags != nil {
				return nil, diags
			}
			pool.Tags = tags

			poolNodes := make([]LKENodePoolNode, len(p.Linodes))
			for j, n := range p.Linodes {
				var node LKENodePoolNode
				node.ID = types.StringValue(n.ID)
				node.InstanceID = types.Int64Value(int64(n.InstanceID))
				node.Status = types.StringValue(string(n.Status))

				poolNodes[j] = node
			}
			pool.Nodes = poolNodes

			// Only parse the autoscaler when it's enabled in order to keep returning
			// the same list result of SDKv2.
			if p.Autoscaler.Enabled {
				var autoscaler LKENodePoolAutoscaler
				autoscaler.Enabled = types.BoolValue(p.Autoscaler.Enabled)
				autoscaler.Min = types.Int64Value(int64(p.Autoscaler.Min))
				autoscaler.Max = types.Int64Value(int64(p.Autoscaler.Max))
				pool.Autoscaler = []LKENodePoolAutoscaler{autoscaler}
			}

			poolDisks := make([]LKENodePoolDisk, len(p.Disks))
			for k, d := range p.Disks {
				var poolDisk LKENodePoolDisk
				poolDisk.Size = types.Int64Value(int64(d.Size))
				poolDisk.Type = types.StringValue(d.Type)
				poolDisks[k] = poolDisk
			}
			pool.Disks = poolDisks

			lkePools[i] = pool
		}
		return lkePools, nil
	}

	lkePools, diags := parseLKEPools()
	if diags != nil {
		return diags
	}
	data.Pools = lkePools

	data.Kubeconfig = types.StringValue(kubeconfig.KubeConfig)

	var urls []string
	for _, e := range endpoints {
		urls = append(urls, e.Endpoint)
	}

	apiEndpoints, diags := types.ListValueFrom(ctx, types.StringType, urls)
	if diags != nil {
		return diags
	}
	data.APIEndpoints = apiEndpoints

	data.DashboardURL = types.StringValue(dashboard.URL)

	return nil
}

func ParseControlPlane(
	controlPlane linodego.LKEClusterControlPlane,
) LKEControlPlane {
	var cp LKEControlPlane
	cp.HighAvailability = types.BoolValue(controlPlane.HighAvailability)

	return cp
}
