package ipv6range

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSourceModel struct {
	Range   types.String `tfsdk:"range"`
	IsBGP   types.Bool   `tfsdk:"is_bgp"`
	Linodes types.Set    `tfsdk:"linodes"`
	Prefix  types.Int64  `tfsdk:"prefix"`
	Region  types.String `tfsdk:"region"`
	ID      types.String `tfsdk:"id"`
}

type ResourceModel struct {
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	LinodeId     types.Int64  `tfsdk:"linode_id"`
	RouteTarget  types.String `tfsdk:"route_target"`
	Range        types.String `tfsdk:"range"`
	IsBGP        types.Bool   `tfsdk:"is_bgp"`
	Linodes      types.Set    `tfsdk:"linodes"`
	Region       types.String `tfsdk:"region"`
	ID           types.String `tfsdk:"id"`
}

func (data *DataSourceModel) parseIPv6RangeDataSource(
	ctx context.Context, ipv6Range *linodego.IPv6Range,
) diag.Diagnostics {
	data.Range = types.StringValue(ipv6Range.Range)
	data.IsBGP = types.BoolValue(ipv6Range.IsBGP)

	linodes, diag := types.SetValueFrom(ctx, types.Int64Type, ipv6Range.Linodes)
	if diag.HasError() {
		return diag
	}
	data.Linodes = linodes

	data.Prefix = types.Int64Value(int64(ipv6Range.Prefix))
	data.Region = types.StringValue(ipv6Range.Region)

	id, _ := json.Marshal(ipv6Range)

	data.ID = types.StringValue(string(id))

	return nil
}

func (r *ResourceModel) FlattenIPv6Range(
	ctx context.Context,
	ipv6Range *linodego.IPv6Range,
	preserveKnown bool,
) diag.Diagnostics {
	linodes, diag := types.SetValueFrom(ctx, types.Int64Type, ipv6Range.Linodes)
	if diag.HasError() {
		return diag
	}

	r.ID = helper.KeepOrUpdateString(r.Range, ipv6Range.Range, preserveKnown)
	r.IsBGP = helper.KeepOrUpdateBool(r.IsBGP, ipv6Range.IsBGP, preserveKnown)
	r.Linodes = helper.KeepOrUpdateValue(r.Linodes, linodes, preserveKnown)
	r.Range = helper.KeepOrUpdateString(r.Range, ipv6Range.Range, preserveKnown)
	r.Region = helper.KeepOrUpdateString(r.Region, ipv6Range.Region, preserveKnown)
	r.PrefixLength = helper.KeepOrUpdateInt64(
		r.PrefixLength, int64(ipv6Range.Prefix), preserveKnown,
	)

	return nil
}

func (r *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	r.IsBGP = helper.KeepOrUpdateValue(r.IsBGP, other.IsBGP, preserveKnown)
	r.Linodes = helper.KeepOrUpdateValue(r.Linodes, other.Linodes, preserveKnown)
	r.Range = helper.KeepOrUpdateValue(r.Range, other.Range, preserveKnown)
	r.Region = helper.KeepOrUpdateValue(r.Region, other.Region, preserveKnown)
	r.PrefixLength = helper.KeepOrUpdateValue(
		r.PrefixLength, other.PrefixLength, preserveKnown,
	)
	r.RouteTarget = helper.KeepOrUpdateValue(r.RouteTarget, other.RouteTarget, preserveKnown)
}
