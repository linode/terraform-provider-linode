package natgatewayaddress

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"natgateway_id": schema.Int64Attribute{
			Description: "The ID of the NAT Gateway the address is assigned to.",
			Required:    true,
		},
		"address": schema.StringAttribute{
			Description: "The public IPv4 address assigned to the NAT Gateway.",
			Required:    true,
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
}
