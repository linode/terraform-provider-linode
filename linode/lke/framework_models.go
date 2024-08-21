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
	HighAvailability types.Bool           `tfsdk:"high_availability"`
	ACL              []LKEControlPlaneACL `tfsdk:"acl"`
}

type LKEControlPlaneACL struct {
	Enabled   types.Bool                    `tfsdk:"enabled"`
	Addresses []LKEControlPlaneACLAddresses `tfsdk:"addresses"`
}

type LKEControlPlaneACLAddresses struct {
	IPv4 types.Set `tfsdk:"ipv4"`
	IPv6 types.Set `tfsdk:"ipv6"`
}

type LKENodePool struct {
	ID             types.Int64             `tfsdk:"id"`
	Count          types.Int64             `tfsdk:"count"`
	Type           types.String            `tfsdk:"type"`
	Tags           types.List              `tfsdk:"tags"`
	DiskEncryption types.String            `tfsdk:"disk_encryption"`
	Disks          []LKENodePoolDisk       `tfsdk:"disks"`
	Nodes          []LKENodePoolNode       `tfsdk:"nodes"`
	Autoscaler     []LKENodePoolAutoscaler `tfsdk:"autoscaler"`
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
	acl *linodego.LKEClusterControlPlaneACLResponse,
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

	cp, diags := parseControlPlane(ctx, cluster.ControlPlane, acl)
	if diags.HasError() {
		return diags
	}

	data.ControlPlane = []LKEControlPlane{cp}

	parseLKEPools := func() ([]LKENodePool, diag.Diagnostics) {
		lkePools := make([]LKENodePool, len(pools))

		for i, p := range pools {
			var pool LKENodePool
			pool.ID = types.Int64Value(int64(p.ID))
			pool.Count = types.Int64Value(int64(p.Count))
			pool.Type = types.StringValue(p.Type)
			pool.DiskEncryption = types.StringValue(string(p.DiskEncryption))

			tags, diags := types.ListValueFrom(ctx, types.StringType, p.Tags)
			if diags != nil {
				return nil, diags
			}
			pool.Tags = tags

			pool.Nodes = make([]LKENodePoolNode, len(p.Linodes))
			for i, linode := range p.Linodes {
				pool.Nodes[i].ID = types.StringValue(linode.ID)
				pool.Nodes[i].InstanceID = types.Int64Value(int64(linode.InstanceID))
				pool.Nodes[i].Status = types.StringValue(string(linode.Status))
			}

			// Only parse the autoscaler when it's enabled in order to keep returning
			// the same list result of SDKv2.
			if p.Autoscaler.Enabled {
				pool.Autoscaler = []LKENodePoolAutoscaler{
					{
						Enabled: types.BoolValue(p.Autoscaler.Enabled),
						Min:     types.Int64Value(int64(p.Autoscaler.Min)),
						Max:     types.Int64Value(int64(p.Autoscaler.Max)),
					},
				}
			}

			pool.Disks = make([]LKENodePoolDisk, len(p.Disks))
			for i, d := range p.Disks {
				pool.Disks[i].Size = types.Int64Value(int64(d.Size))
				pool.Disks[i].Type = types.StringValue(d.Type)
			}

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

func parseControlPlane(
	ctx context.Context,
	controlPlane linodego.LKEClusterControlPlane,
	aclResp *linodego.LKEClusterControlPlaneACLResponse,
) (LKEControlPlane, diag.Diagnostics) {
	var cp LKEControlPlane

	if aclResp != nil {
		acl := aclResp.ACL
		var aclAddresses LKEControlPlaneACLAddresses

		ipv4, diags := types.SetValueFrom(ctx, types.StringType, acl.Addresses.IPv4)
		if diags.HasError() {
			return cp, diags
		}
		aclAddresses.IPv4 = ipv4

		ipv6, diags := types.SetValueFrom(ctx, types.StringType, acl.Addresses.IPv6)
		if diags.HasError() {
			return cp, diags
		}
		aclAddresses.IPv6 = ipv6

		var cpACL LKEControlPlaneACL
		cpACL.Enabled = types.BoolValue(acl.Enabled)
		cpACL.Addresses = []LKEControlPlaneACLAddresses{aclAddresses}
		cp.ACL = []LKEControlPlaneACL{cpACL}
	} else {
		cp.ACL = []LKEControlPlaneACL{}
	}

	cp.HighAvailability = types.BoolValue(controlPlane.HighAvailability)

	return cp, nil
}
