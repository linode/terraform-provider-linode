package region

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ObjectTypePGLimits = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"maximum_pgs_per_customer": types.Int64Type,
		"maximum_linodes_per_pg":   types.Int64Type,
	},
}

var DataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"country": schema.StringAttribute{
			Description: "The country where this Region resides.",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "The unique ID of this Region.",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "Detailed location information for this Region, including city, state or region, and country.",
			Computed:    true,
		},
		"site_type": schema.StringAttribute{
			Description: "The type of this Region.",
			Computed:    true,
		},
		"capabilities": schema.SetAttribute{
			Description: "A list of capabilities of this region.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"status": schema.StringAttribute{
			Description: "This region’s current operational status.",
			Computed:    true,
		},
		"placement_group_limits": schema.ListAttribute{
			Description: "Information about placement groups limits for this region.",
			Computed:    true,
			ElementType: ObjectTypePGLimits,
		},
	},
	Blocks: map[string]schema.Block{
		"resolvers": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"ipv4": schema.StringAttribute{
						Description: "The IPv4 addresses for this region’s DNS resolvers, separated by commas.",
						Computed:    true,
					},
					"ipv6": schema.StringAttribute{
						Description: "The IPv6 addresses for this region’s DNS resolvers, separated by commas.",
						Computed:    true,
					},
				},
			},
		},
	},
}
