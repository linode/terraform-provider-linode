package reservediptypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

// PriceModel is the tfsdk-tagged struct for a single price entry.
type PriceModel struct {
	Hourly  types.Float64 `tfsdk:"hourly"`
	Monthly types.Float64 `tfsdk:"monthly"`
}

// RegionPriceModel is the tfsdk-tagged struct for a region-specific price entry.
type RegionPriceModel struct {
	ID      types.String  `tfsdk:"id"`
	Hourly  types.Float64 `tfsdk:"hourly"`
	Monthly types.Float64 `tfsdk:"monthly"`
}

type DataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	Price        types.List   `tfsdk:"price"`
	RegionPrices types.List   `tfsdk:"region_prices"`
}

func (data *DataSourceModel) parseReservedIPType(ctx context.Context, ipType *linodego.ReservedIPType) diag.Diagnostics {
	data.ID = types.StringValue(ipType.ID)
	data.Label = types.StringValue(ipType.Label)

	price, diags := types.ListValueFrom(ctx, helper.PriceObjectType, []PriceModel{
		{
			Hourly:  types.Float64Value(ipType.Price.Hourly),
			Monthly: types.Float64Value(ipType.Price.Monthly),
		},
	})
	if diags.HasError() {
		return diags
	}
	data.Price = price

	regionPrices := make([]RegionPriceModel, len(ipType.RegionPrices))
	for i, rp := range ipType.RegionPrices {
		regionPrices[i] = RegionPriceModel{
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

type ReservedIPTypeFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Types   []DataSourceModel                `tfsdk:"types"`
}

func (model *ReservedIPTypeFilterModel) parseReservedIPTypes(ctx context.Context, ipTypes []linodego.ReservedIPType) diag.Diagnostics {
	result := make([]DataSourceModel, len(ipTypes))

	for i := range ipTypes {
		var m DataSourceModel

		diags := m.parseReservedIPType(ctx, &ipTypes[i])
		if diags.HasError() {
			return diags
		}

		result[i] = m
	}

	model.Types = result

	return nil
}
