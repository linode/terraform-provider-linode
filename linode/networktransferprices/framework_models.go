package networktransferprices

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type DataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	Price        types.List   `tfsdk:"price"`
	RegionPrices types.List   `tfsdk:"region_prices"`
	Transfer     types.Int64  `tfsdk:"transfer"`
}

func (data *DataSourceModel) parseNetworkTransferPrice(networkTransferPrice *linodego.NetworkTransferPrice,
) diag.Diagnostics {
	data.ID = types.StringValue(networkTransferPrice.ID)

	price, diags := flattenPrice(networkTransferPrice.Price)
	if diags.HasError() {
		return diags
	}
	data.Price = *price

	data.Label = types.StringValue(networkTransferPrice.Label)

	regionPrices, d := flattenRegionPrices(networkTransferPrice.RegionPrices)
	if d.HasError() {
		return d
	}
	data.RegionPrices = *regionPrices

	data.Transfer = types.Int64Value(int64(networkTransferPrice.Transfer))

	return nil
}

func flattenPrice(price linodego.NetworkTransferTypePrice) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	result["hourly"] = types.Float64Value(float64(price.Hourly))
	result["monthly"] = types.Float64Value(float64(price.Monthly))

	obj, diag := types.ObjectValue(helper.PriceObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	objList := []attr.Value{obj}

	resultList, diag := types.ListValue(
		helper.PriceObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

func flattenRegionPrices(prices []linodego.NetworkTransferTypeRegionPrice) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make([]attr.Value, len(prices))

	for i, price := range prices {
		obj, d := types.ObjectValue(helper.RegionPriceObjectType.AttrTypes, map[string]attr.Value{
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
		helper.RegionPriceObjectType,
		result,
	)
	return &priceList, d
}

type NetworkTransferPriceFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Types   []DataSourceModel                `tfsdk:"types"`
}

func (model *NetworkTransferPriceFilterModel) parseNetworkTransferPrices(networkTransferPrices []linodego.NetworkTransferPrice,
) diag.Diagnostics {
	result := make([]DataSourceModel, len(networkTransferPrices))

	for i := range networkTransferPrices {
		var m DataSourceModel

		diags := m.parseNetworkTransferPrice(&networkTransferPrices[i])
		if diags.HasError() {
			return diags
		}

		result[i] = m
	}

	model.Types = result

	return nil
}
