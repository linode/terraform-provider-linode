package networkingips

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var updatedIPObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"address":     types.StringType,
		"region":      types.StringType,
		"gateway":     types.StringType,
		"subnet_mask": types.StringType,
		"prefix":      types.Int64Type,
		"type":        types.StringType,
		"public":      types.BoolType,
		"rdns":        types.StringType,
		"linode_id":   types.Int64Type,
		"reserved":    types.BoolType,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Description: "The IP address.",
			// Required:    true,
			Optional: true,
		},
		"gateway": schema.StringAttribute{
			Description: "The default gateway for this address.",
			Computed:    true,
		},
		"subnet_mask": schema.StringAttribute{
			Description: "The mask that separates host bits from network bits for this address.",
			Computed:    true,
		},
		"prefix": schema.Int64Attribute{
			Description: "The number of bits set in the subnet mask.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of address this is (ipv4, ipv6, ipv6/pool, ipv6/range).",
			Computed:    true,
		},
		"public": schema.BoolAttribute{
			Description: "Whether this is a public or private IP address.",
			Computed:    true,
		},
		"rdns": schema.StringAttribute{
			Description: "The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to " +
				"a default value provided by Linode if not explicitly set.",
			Computed: true,
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode this address currently belongs to.",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The Region this IP address resides in.",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "A unique identifier for this datasource.",
			Computed:    true,
		},
		"reserved": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this IP is reserved or not.",
		},
		"ip_addresses": schema.ListAttribute{
			Description: "A list of all IPs.",
			Computed:    true,
			ElementType: updatedIPObjectType,
		},
		"filter_reserved": schema.BoolAttribute{
			Description: "Filter IPs by reserved status.",
			Optional:    true,
		},
	},
}
