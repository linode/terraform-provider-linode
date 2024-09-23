package nbtypes

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type DataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	Price        types.List   `tfsdk:"price"`
	RegionPrices types.List   `tfsdk:"region_prices"`
}

func (data *DataSourceModel) ParseNodeBalancerType(nbType *linodego.NodeBalancerType,
) diag.Diagnostics {
	data.ID = types.StringValue(nbType.ID)

	price, diags := FlattenPrice(nbType.Price)
	if diags.HasError() {
		return diags
	}
	data.Price = *price

	data.Label = types.StringValue(nbType.Label)

	regionPrices, d := FlattenRegionPrices(nbType.RegionPrices)
	if d.HasError() {
		return d
	}
	data.RegionPrices = *regionPrices

	return nil
}

func FlattenPrice(price linodego.NodeBalancerTypePrice) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	result["hourly"] = types.Float64Value(float64(price.Hourly))
	result["monthly"] = types.Float64Value(float64(price.Monthly))

	obj, diag := types.ObjectValue(priceObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	objList := []attr.Value{obj}

	resultList, diag := types.ListValue(
		priceObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

func FlattenRegionPrices(prices []linodego.NodeBalancerTypeRegionPrice) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make([]attr.Value, len(prices))

	for i, price := range prices {
		obj, d := types.ObjectValue(regionPriceObjectType.AttrTypes, map[string]attr.Value{
			"id":      types.StringValue(price.ID),
			"hourly":  types.Float64Value(float64(price.Hourly)),
			"monthly": types.Float64Value(float64(price.Monthly)),
		})
		if d.HasError() {
			return nil, d
		}

		result[i] = obj
	}

	priceList, d := basetypes.NewListValue(
		regionPriceObjectType,
		result,
	)
	return &priceList, d
}

type NodeBalancerTypeFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Types   []DataSourceModel                `tfsdk:"types"`
}

func (model *NodeBalancerTypeFilterModel) parseNodeBalancerTypes(nbTypes []linodego.NodeBalancerType,
) diag.Diagnostics {
	result := make([]DataSourceModel, len(nbTypes))

	for i := range nbTypes {
		var m DataSourceModel

		diags := m.ParseNodeBalancerType(&nbTypes[i])
		if diags.HasError() {
			return diags
		}

		result[i] = m
	}

	model.Types = result

	return nil
}
