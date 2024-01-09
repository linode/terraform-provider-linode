package ipv6ranges

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"route_target": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"region":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"range":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"prefix":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"ranges": schema.ListNestedBlock{
			Description: "The return list of IPv6 ranges.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"range": schema.StringAttribute{
						Description: "The IPv6 address of this range.",
						Required:    true,
					},
					"route_target": schema.StringAttribute{
						Description: "The IPv6 SLAAC address.",
						Computed:    true,
					},
					"prefix": schema.Int64Attribute{
						Description: "The prefix length of the address, denoting how many addresses can be assigned from this range.",
						Computed:    true,
					},
					"region": schema.StringAttribute{
						Description: "The region for this range of IPv6 addresses.",
						Computed:    true,
					},
				},
			},
		},
	},
}
