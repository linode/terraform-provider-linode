package vpcdefaultranges

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

type DataSourceModel struct {
	DefaultIPv4Ranges   types.Set `tfsdk:"default_ipv4_ranges"`
	ForbiddenIPv4Ranges types.Set `tfsdk:"forbidden_ipv4_ranges"`
}

func (data *DataSourceModel) parseVPCDefaultRanges(
	ctx context.Context,
	ranges *linodego.VPCDefaultRanges,
) diag.Diagnostics {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	data.DefaultIPv4Ranges, d = types.SetValueFrom(ctx, types.StringType, ranges.DefaultIPV4Ranges)
	diags.Append(d...)

	data.ForbiddenIPv4Ranges, d = types.SetValueFrom(ctx, types.StringType, ranges.ForbiddenIPV4Ranges)
	diags.Append(d...)

	return diags
}
