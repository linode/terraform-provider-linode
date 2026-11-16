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
	ID    types.String `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
	Price types.List   `tfsdk:"price"`
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

	return nil
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
