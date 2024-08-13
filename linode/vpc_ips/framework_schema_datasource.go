package vpc_ips

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var VPCIPAttrs = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		Description: "An IPv4 address configured for this VPC interface. These follow the RFC 1918 private address format. Displayed as null if an address_range.",
		Computed:    true,
	},
	"gateway": schema.StringAttribute{
		Description: "The default gateway for the VPC subnet that the IP or IP range belongs to.",
		Computed:    true,
	},
	"linode_id": schema.Int64Attribute{
		Description: "The identifier for the Linode the VPC interface currently belongs to.",
		Computed:    true,
	},
	"prefix": schema.Int64Attribute{
		Description: "The number of bits set in the subnet_mask.",
		Computed:    true,
	},
	"region": schema.StringAttribute{
		Description: "The region of the VPC.",
		Computed:    true,
	},
	"subnet_mask": schema.StringAttribute{
		Description: "The mask that separates host bits from network bits for the address or address_range.",
		Computed:    true,
	},
	"nat_1_1": schema.StringAttribute{
		Description: "The public IP address used for NAT 1:1 with the VPC. This is empty if NAT 1:1 isn't used.",
		Computed:    true,
	},
	"subnet_id": schema.Int64Attribute{
		Description: "The id of the VPC Subnet for this interface.",
		Computed:    true,
	},
	"config_id": schema.Int64Attribute{
		Description: "The globally general entity identifier for the Linode configuration profile where the VPC is included.",
		Computed:    true,
	},
	"interface_id": schema.Int64Attribute{
		Description: "The globally general API entity identifier for the Linode interface.",
		Computed:    true,
	},
	"address_range": schema.StringAttribute{
		Description: "A range of IPv4 addresses configured for this VPC interface. Displayed as null if a single address.",
		Computed:    true,
	},
	"vpc_id": schema.Int64Attribute{
		Description: "The unique globally general API entity identifier for the VPC.",
		Computed:    true,
	},
	"active": schema.BoolAttribute{
		Description: "Returns true if the VPC interface is in use, meaning that the Linode was powered on using the config_id to which the interface belongs. Otherwise returns false",
		Computed:    true,
	},
}

var filterConfig = frameworkfilter.Config{
	"active":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
	"config_id": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"linode_id": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"region":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"vpc_id":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
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
