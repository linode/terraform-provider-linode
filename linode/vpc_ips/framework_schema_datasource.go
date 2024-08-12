package vpc_ips

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var VPCIPAttrs = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		Description: "The IP address in CIDR format.",
		Computed:    true,
	},
	"gateway": schema.StringAttribute{
		Description: "The default gateway for this address.",
		Computed:    true,
	},
	"linode_id": schema.Int64Attribute{
		Description: "The ID of the Linode this address currently belongs to. " +
			"For IPv4 addresses, this defaults to the Linode that this address was assigned to on creation.",
		Computed: true,
	},
	"prefix": schema.Int64Attribute{
		Description: "The number of bits set in the subnet mask.",
		Computed:    true,
	},
	"region": schema.StringAttribute{
		Description: "The Region this IP address resides in.",
		Computed:    true,
	},
	"subnet_mask": schema.StringAttribute{
		Description: "The mask that separates host bits from network bits for this address.",
		Computed:    true,
	},
	"nat_1_1": schema.StringAttribute{
		Description: "IPv4 address configured as a 1:1 NAT for this Interface. " +
			"If no address is configured as a 1:1 NAT, null is returned." +
			"Note: Only allowed for vpc type Interfaces.",
		Computed: true,
	},
	"subnet_id": schema.Int64Attribute{
		Description: "The ID of the subnet this IP address currently belongs to.",
		Computed:    true,
	},
	"config_id": schema.Int64Attribute{
		Description: "The ID of the config this IP address is associated with.",
		Computed:    true,
	},
	"interface_id": schema.Int64Attribute{
		Description: "The ID of the interface this IP address is associated with.",
		Computed:    true,
	},
	"address_range": schema.StringAttribute{
		Description: "The IP address range that this IP address is associated with.",
		Computed:    true,
	},
	"vpc_id": schema.Int64Attribute{
		Description: "The ID of the VPC this IP address is associated with.",
		Computed:    true,
	},
	"active": schema.BoolAttribute{
		Description: "Indicates whether this IP address is active or not.",
		Computed:    true,
	},
}

var filterConfig = frameworkfilter.Config{
	"address": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"prefix":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"region":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"vpc_id": schema.Int64Attribute{
			Description: "The ID of the VPC that the list of IP addresses is associated with.",
			Optional:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"vpc_ips": schema.ListNestedBlock{
			Description: "The returned list of IP addresses that exist in Linode's system, either IPv4 or IPv6.",
			NestedObject: schema.NestedBlockObject{
				Attributes: VPCIPAttrs,
			},
		},
	},
}
