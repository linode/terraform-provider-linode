package ipv6ranges

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type IPv6ListEntryModel struct {
	Range       types.String `tfsdk:"range"`
	RouteTarget types.String `tfsdk:"route_target"`
	Prefix      types.Int64  `tfsdk:"prefix"`
	Region      types.String `tfsdk:"region"`
}

func (m *IPv6ListEntryModel) ParseRange(r linodego.IPv6Range) {
	m.Range = types.StringValue(r.Range)
	m.RouteTarget = types.StringValue(r.RouteTarget)
	m.Prefix = types.Int64Value(int64(r.Prefix))
	m.Region = types.StringValue(r.Region)
}

type IPv6RangeFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Ranges  []IPv6ListEntryModel             `tfsdk:"ranges"`
}

func (data *IPv6RangeFilterModel) parseRanges(
	ranges []linodego.IPv6Range,
) diag.Diagnostics {
	result := make([]IPv6ListEntryModel, len(ranges))

	for i, r := range ranges {
		var mod IPv6ListEntryModel
		mod.ParseRange(r)
		result[i] = mod
	}

	data.Ranges = result

	return nil
}
