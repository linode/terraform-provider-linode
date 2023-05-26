package ipv6range

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "The IPv6 range to retrieve information about.",
			Required:    true,
		},
		"is_bgp": schema.BoolAttribute{
			Description: "Whether this IPv6 range is shared.",
			Computed:    true,
		},
		"linodes": schema.SetAttribute{
			ElementType: types.Int64Type,
			Description: "The IDs of Linodes to apply this firewall to.",
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
		"id": schema.StringAttribute{
			Description: "The unique ID for this DataSource",
			Computed:    true,
		},
	},
}
