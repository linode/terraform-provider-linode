package regionvpcavailability

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var FrameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the Region.",
			Required:    true,
		},
		"available": schema.BoolAttribute{
			Description: "Whether VPCs can be created in the region.",
			Computed:    true,
		},
		"available_ipv6_prefix_lengths": schema.ListAttribute{
			ElementType: types.Int64Type,
			Description: "The IPv6 prefix lengths allowed when creating a VPC in the region.",
			Computed:    true,
		},
	},
}
