package natgatewaytypes

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

// DataSourceModel describes a single NAT Gateway type in the `types` list.
type DataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	Price        types.List   `tfsdk:"price"`
	RegionPrices types.List   `tfsdk:"region_prices"`
	Transfer     types.Int64  `tfsdk:"transfer"`
}

// NATGatewayTypeFilterModel describes the Terraform data source model to match the
// list data source schema.
type NATGatewayTypeFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Types   []DataSourceModel                `tfsdk:"types"`
}

func (model *NATGatewayTypeFilterModel) ParseNATGatewayTypes(natGatewayTypes []linodego.NATGatewayType) diag.Diagnostics {
	var diags diag.Diagnostics

	result := make([]DataSourceModel, len(natGatewayTypes))
	for i := range natGatewayTypes {
		var m DataSourceModel
		diags.Append(m.parseNATGatewayType(&natGatewayTypes[i])...)
		result[i] = m
	}
	model.Types = result

	return diags
}

func (data *DataSourceModel) parseNATGatewayType(natGatewayType *linodego.NATGatewayType) diag.Diagnostics {
	data.ID = types.StringValue(natGatewayType.ID)
	data.Label = types.StringValue(natGatewayType.Label)

	price, diags := flattenPrice(natGatewayType.Price.Hourly, natGatewayType.Price.Monthly)
	if diags.HasError() {
		return diags
	}
	data.Price = *price

	regionPrices, diags := flattenRegionPrices(natGatewayType.RegionPrices)
	if diags.HasError() {
		return diags
	}
	data.RegionPrices = *regionPrices

	data.Transfer = types.Int64Value(int64(natGatewayType.Transfer))

	return nil
}

func flattenRegionPrices(prices []linodego.NATGatewayTypeRegionPrice) (*types.List, diag.Diagnostics) {
	result := make([]attr.Value, len(prices))
	for i, price := range prices {
		obj, diags := types.ObjectValue(helper.RegionPriceObjectType.AttrTypes, map[string]attr.Value{
			"id":      types.StringValue(price.ID),
			"hourly":  types.Float64Value(price.Hourly),
			"monthly": types.Float64Value(price.Monthly),
		})
		if diags.HasError() {
			return nil, diags
		}

		result[i] = obj
	}

	list, diags := types.ListValue(helper.RegionPriceObjectType, result)
	if diags.HasError() {
		return nil, diags
	}

	return &list, nil
}

func flattenPrice(hourly, monthly float64) (*types.List, diag.Diagnostics) {
	obj, diags := types.ObjectValue(helper.PriceObjectType.AttrTypes, map[string]attr.Value{
		"hourly":  types.Float64Value(hourly),
		"monthly": types.Float64Value(monthly),
	})
	if diags.HasError() {
		return nil, diags
	}

	result, diags := types.ListValue(helper.PriceObjectType, []attr.Value{obj})
	if diags.HasError() {
		return nil, diags
	}

	return &result, nil
}
