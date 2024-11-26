package reservedip

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"region": schema.StringAttribute{
			Description: "The Region in which to reserve the IP address.",
			Optional:    true,
		},
		"address": schema.StringAttribute{
			Description: "The reserved IP address.",
			Required:    true,
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
		"vpc_nat_1_1": schema.ListAttribute{
			Description: "Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.",
			Computed:    true,
			ElementType: instancenetworking.VPCNAT1To1Type,
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
			},
		},
		"rdns": schema.StringAttribute{
			Description: "The reverse DNS assigned to this address.",
			Computed:    true,
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode this address currently belongs to.",
			Computed:    true,
		},
		"reserved": schema.BoolAttribute{
			Description: "Whether this IP is reserved or not.",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "The unique ID of the reserved IP address.",
			Computed:    true,
		},
	},
}
