package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var PriceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"hourly":  types.Float64Type,
		"monthly": types.Float64Type,
	},
}

var RegionPriceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":      types.StringType,
		"hourly":  types.Float64Type,
		"monthly": types.Float64Type,
	},
}

func GetPricingTypeAttributes(typeName string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID assigned to this " + typeName + ".",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The " + typeName + "'s label.",
			Computed:    true,
			Optional:    true,
		},
		"price": schema.ListAttribute{
			Description: "Cost in US dollars, broken down into hourly and monthly charges.",
			Computed:    true,
			ElementType: PriceObjectType,
		},
		"region_prices": schema.ListAttribute{
			Description: "A list of region-specific prices for this " + typeName + ".",
			Computed:    true,
			ElementType: RegionPriceObjectType,
		},
		"transfer": schema.Int64Attribute{
			Description: "The monthly outbound transfer amount, in MB.",
			Computed:    true,
		},
	}
}
