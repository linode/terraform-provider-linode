package volumetypes

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var priceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"hourly":  types.Float64Type,
		"monthly": types.Float64Type,
	},
}

var regionPriceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":      types.StringType,
		"hourly":  types.Float64Type,
		"monthly": types.Float64Type,
	},
}

var Attributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique ID assigned to this Volume type.",
		Required:    true,
	},
	"label": schema.StringAttribute{
		Description: "The Volume Type's label.",
		Computed:    true,
		Optional:    true,
	},
	"price": schema.ListAttribute{
		Description: "Cost in US dollars, broken down into hourly and monthly charges.",
		Computed:    true,
		ElementType: priceObjectType,
	},
	"region_prices": schema.ListAttribute{
		Description: "A list of region-specific prices for this Volume type.",
		Computed:    true,
		ElementType: regionPriceObjectType,
	},
	"transfer": schema.Int64Attribute{
		Description: "The monthly outbound transfer amount, in MB.",
		Computed:    true,
	},
}

var volumeTypeSchema = schema.NestedBlockObject{
	Attributes: Attributes,
}

var filterConfig = frameworkfilter.Config{
	"id":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"label": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order_by": filterConfig.OrderBySchema(),
		"order":    filterConfig.OrderSchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"types": schema.ListNestedBlock{
			Description:  "The returned list of Volume types.",
			NestedObject: volumeTypeSchema,
		},
	},
}
