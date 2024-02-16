package nbs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/nb"
)

// NodeBalancerFilterModel describes the Terraform resource data model to match the
// resource schema.
type NodeBalancerFilterModel struct {
	ID            types.String                     `tfsdk:"id"`
	Filters       frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order         types.String                     `tfsdk:"order"`
	OrderBy       types.String                     `tfsdk:"order_by"`
	NodeBalancers []NodeBalancerModel              `tfsdk:"nodebalancers"`
}

type NodeBalancerModel struct {
	ID                 types.Int64       `tfsdk:"id"`
	Label              types.String      `tfsdk:"label"`
	Region             types.String      `tfsdk:"region"`
	ClientConnThrottle types.Int64       `tfsdk:"client_conn_throttle"`
	Hostname           types.String      `tfsdk:"hostname"`
	Ipv4               types.String      `tfsdk:"ipv4"`
	Ipv6               types.String      `tfsdk:"ipv6"`
	Created            timetypes.RFC3339 `tfsdk:"created"`
	Updated            timetypes.RFC3339 `tfsdk:"updated"`
	Transfer           types.List        `tfsdk:"transfer"`
	Tags               types.Set         `tfsdk:"tags"`
}

func (data *NodeBalancerFilterModel) parseNodeBalancers(
	ctx context.Context,
	nodebalancers []linodego.NodeBalancer,
) {
	result := make([]NodeBalancerModel, len(nodebalancers))
	for i := range nodebalancers {
		var nbData NodeBalancerModel
		nbData.flattenNodeBalancer(ctx, &nodebalancers[i])
		result[i] = nbData
	}

	data.NodeBalancers = result
}

func (data *NodeBalancerModel) flattenNodeBalancer(
	ctx context.Context,
	nodebalancer *linodego.NodeBalancer,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(nodebalancer.ID))
	data.Label = types.StringPointerValue(nodebalancer.Label)
	data.ID = types.Int64Value(int64(nodebalancer.ID))
	data.Region = types.StringValue(nodebalancer.Region)
	data.ClientConnThrottle = types.Int64Value(int64(nodebalancer.ClientConnThrottle))
	data.Hostname = types.StringPointerValue(nodebalancer.Hostname)
	data.Ipv4 = types.StringPointerValue(nodebalancer.IPv4)
	data.Ipv6 = types.StringPointerValue(nodebalancer.IPv6)
	data.Created = timetypes.NewRFC3339TimePointerValue(nodebalancer.Created)
	data.Updated = timetypes.NewRFC3339TimePointerValue(nodebalancer.Updated)

	transfer, diags := nb.FlattenTransfer(ctx, nodebalancer.Transfer)
	if diags.HasError() {
		return diags
	}
	data.Transfer = *transfer

	tags, diags := types.SetValueFrom(
		ctx,
		types.StringType,
		helper.StringSliceToFramework(nodebalancer.Tags),
	)
	if diags.HasError() {
		return diags
	}
	data.Tags = tags

	return nil
}
