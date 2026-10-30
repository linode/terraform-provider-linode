package natgatewayaddresses

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

// TODO: Replace with API filterable fields once they are available
var filterConfig = frameworkfilter.Config{
	"address": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"in_use":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"natgateway_id": schema.Int64Attribute{
			Description: "The ID of the NAT Gateway whose addresses are listed.",
			Required:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"addresses": schema.ListNestedAttribute{
			Description: "The returned list of addresses assigned to the NAT Gateway.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Description: "The public IPv4 address assigned to the NAT Gateway.",
						Computed:    true,
					},
					"in_use": schema.BoolAttribute{
						Description: "Whether the address is currently in use by one or more interfaces.",
						Computed:    true,
					},
					"interface_count": schema.Int64Attribute{
						Description: "The number of interfaces using this address.",
						Computed:    true,
					},
					"interface_url": schema.StringAttribute{
						Description: "The API URL of the interfaces using this address.",
						Computed:    true,
					},
					"portset_assignments": schema.Int64Attribute{
						Description: "The number of portsets currently assigned to interfaces using this address, excluding system-reserved ports.",
						Computed:    true,
					},
					"portset_capacity": schema.Int64Attribute{
						Description: "The total number of portsets available to allocate for this address, excluding system-reserved ports.",
						Computed:    true,
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
