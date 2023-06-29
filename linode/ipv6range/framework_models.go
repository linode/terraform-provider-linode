package ipv6range

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
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

func (data *DataSourceModel) parseIPv6Range(
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
