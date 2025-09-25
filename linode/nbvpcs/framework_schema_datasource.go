package nbvpcs

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/nbvpc"
)

var filterConfig = frameworkfilter.Config{
	"id":              {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"ipv4_range":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"nodebalancer_id": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"subnet_id":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"vpc_id":          {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"nodebalancer_id": schema.Int64Attribute{
			Description: "The ID of the NodeBalancer to list VPC configurations for.",
			Required:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"vpc_configs": schema.ListNestedAttribute{
			Description: "The returned list of NodeBalancer-VPC configurations.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: nbvpc.DataSourceSchema.Attributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
