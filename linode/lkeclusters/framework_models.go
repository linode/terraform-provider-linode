package lkeclusters

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

// LKEClusterFilterModel describes the Terraform resource data model to match the
// resource schema.
type LKEClusterFilterModel struct {
	ID          types.String                     `tfsdk:"id"`
	Filters     frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order       types.String                     `tfsdk:"order"`
	OrderBy     types.String                     `tfsdk:"order_by"`
	LKEClusters []LKEClusterModel                `tfsdk:"lke_clusters"`
}

type LKEClusterModel struct {
	ID           types.Int64          `tfsdk:"id"`
	Created      types.String         `tfsdk:"created"`
	Updated      types.String         `tfsdk:"updated"`
	Label        types.String         `tfsdk:"label"`
	Region       types.String         `tfsdk:"region"`
	Status       types.String         `tfsdk:"status"`
	K8sVersion   types.String         `tfsdk:"k8s_version"`
	Tags         types.Set            `tfsdk:"tags"`
	ControlPlane LKEControlPlaneModel `tfsdk:"control_plane"`
	Tier         types.String         `tfsdk:"tier"`
	APLEnabled   types.Bool           `tfsdk:"apl_enabled"`
	SubnetID     types.Int64          `tfsdk:"subnet_id"`
	VpcID        types.Int64          `tfsdk:"vpc_id"`
	StackType    types.String         `tfsdk:"stack_type"`
}

type LKEControlPlaneModel struct {
	HighAvailability types.Bool `tfsdk:"high_availability"`
	AuditLogsEnabled types.Bool `tfsdk:"audit_logs_enabled"`
}

func (data *LKEClusterFilterModel) parseLKEClusters(
	ctx context.Context,
	clusters []linodego.LKECluster,
) diag.Diagnostics {
	result := make([]LKEClusterModel, len(clusters))
	for i := range clusters {
		var lkeCluster LKEClusterModel
		diags := lkeCluster.parseLKECluster(ctx, &clusters[i])
		if diags != nil {
			return diags
		}
		result[i] = lkeCluster
	}

	data.LKEClusters = result
	return nil
}

func (data *LKEClusterModel) parseLKECluster(
	ctx context.Context,
	cluster *linodego.LKECluster,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(cluster.ID))
	data.Created = types.StringValue(cluster.Created.Format(helper.TIME_FORMAT))
	data.Updated = types.StringValue(cluster.Updated.Format(helper.TIME_FORMAT))
	data.Label = types.StringValue(cluster.Label)
	data.Region = types.StringValue(cluster.Region)
	data.Status = types.StringValue(string(cluster.Status))
	data.K8sVersion = types.StringValue(cluster.K8sVersion)
	data.Tier = types.StringValue(cluster.Tier)
	data.SubnetID = types.Int64Value(int64(cluster.SubnetID))
	data.VpcID = types.Int64Value(int64(cluster.VpcID))
	data.StackType = types.StringValue(string(cluster.StackType))

	tags, diags := types.SetValueFrom(ctx, types.StringType, cluster.Tags)
	if diags != nil {
		return diags
	}
	data.Tags = tags

	parseControlPlane := func() LKEControlPlaneModel {
		var cp LKEControlPlaneModel
		cp.HighAvailability = types.BoolValue(cluster.ControlPlane.HighAvailability)
		cp.AuditLogsEnabled = types.BoolValue(cluster.ControlPlane.AuditLogsEnabled)

		return cp
	}

	data.ControlPlane = parseControlPlane()

	return nil
}
