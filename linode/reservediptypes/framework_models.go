package reservediptypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

// priceModel is the tfsdk-tagged struct for a single price entry.
type priceModel struct {
	Hourly  types.Float64 `tfsdk:"hourly"`
	Monthly types.Float64 `tfsdk:"monthly"`
}

// regionPriceModel is the tfsdk-tagged struct for a region-specific price entry.
type regionPriceModel struct {
	ID      types.String  `tfsdk:"id"`
	Hourly  types.Float64 `tfsdk:"hourly"`
	Monthly types.Float64 `tfsdk:"monthly"`
}

type dataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	Price        types.List   `tfsdk:"price"`
	RegionPrices types.List   `tfsdk:"region_prices"`
}

func (data *dataSourceModel) parseReservedIPType(ctx context.Context, ipType *linodego.ReservedIPType) diag.Diagnostics {
	data.ID = types.StringValue(ipType.ID)
	data.Label = types.StringValue(ipType.Label)

	price, diags := types.ListValueFrom(ctx, helper.PriceObjectType, []priceModel{
		{
			Hourly:  types.Float64Value(ipType.Price.Hourly),
			Monthly: types.Float64Value(ipType.Price.Monthly),
		},
	})
	if diags.HasError() {
		return diags
	}
	data.Price = price

	regionPrices := make([]regionPriceModel, len(ipType.RegionPrices))
	for i, rp := range ipType.RegionPrices {
		regionPrices[i] = regionPriceModel{
			ID:      types.StringValue(rp.ID),
			Hourly:  types.Float64Value(rp.Hourly),
			Monthly: types.Float64Value(rp.Monthly),
		}
	}
	rpList, diags := types.ListValueFrom(ctx, helper.RegionPriceObjectType, regionPrices)
	if diags.HasError() {
		return diags
	}
	data.RegionPrices = rpList

	return nil
}

type reservedIPTypeFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Types   []dataSourceModel                `tfsdk:"types"`
}

func (model *reservedIPTypeFilterModel) parseReservedIPTypes(ctx context.Context, ipTypes []linodego.ReservedIPType) diag.Diagnostics {
	result := make([]dataSourceModel, len(ipTypes))

	for i := range ipTypes {
		var m dataSourceModel

		diags := m.parseReservedIPType(ctx, &ipTypes[i])
		if diags.HasError() {
			return diags
		}

		result[i] = m
	}

	model.Types = result

	return nil
}
