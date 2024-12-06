package reservedip

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the reserved IP address, which will be the IP address itself.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Description: "The Region in which to reserve the IP address.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"address": schema.StringAttribute{
			Description: "The reserved IP address.",
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
		"vpc_nat_1_1": schema.ListAttribute{
			Description: "Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.",
			Computed:    true,
			ElementType: instancenetworking.VPCNAT1To1Type,
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
			},
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
		"reserved": schema.BoolAttribute{
			Description: "Whether this IP is reserved or not.",
			Computed:    true,
		},
	},
}
