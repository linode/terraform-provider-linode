package networkingips

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/instancenetworking"
)

var filterConfig = frameworkfilter.Config{
	"type":    {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"region":  {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"rdns":    {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"address": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"prefix":  {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},

	"gateway":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"subnet_mask": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"public":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"linode_id":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"reserved":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"ip_addresses": schema.ListNestedBlock{
			Description: "The returned list of Images.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Description: "The IP address.",
						Computed:    true,
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
					"reserved": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether this IP is reserved or not.",
					},
					"vpc_nat_1_1": schema.ObjectAttribute{
						Description:    "Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.",
						Computed:       true,
						AttributeTypes: instancenetworking.VPCNAT1To1Type.AttrTypes,
					},
				},
			},
		},
	},
}
