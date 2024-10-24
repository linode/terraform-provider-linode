package instancereservedip

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the IPv4 address, which will be IPv4 address itself.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode to allocate an IPv4 address for.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"public": schema.BoolAttribute{
			Description: "Whether the IPv4 address is public or private.",
			Default:     booldefault.StaticBool(true),
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},

		"address": schema.StringAttribute{
			Description: "The resulting IPv4 address.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"gateway": schema.StringAttribute{
			Description: "The default gateway for this address",
			Computed:    true,
		},
		"prefix": schema.Int64Attribute{
			Description: "The number of bits set in the subnet mask.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"rdns": schema.StringAttribute{
			Description: "The reverse DNS assigned to this address.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Description: "The region this IP resides in.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"subnet_mask": schema.StringAttribute{
			Description: "The mask that separates host bits from network bits for this address.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"type": schema.StringAttribute{
			Description: "The type of IP address.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"vpc_nat_1_1": schema.ListAttribute{
			Description: "Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.",
			Computed:    true,
			ElementType: instancenetworking.VPCNAT1To1Type,
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
			},
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},

		"apply_immediately": schema.BoolAttribute{
			Description: "If true, the instance will be rebooted to update network interfaces. " +
				"This functionality is not affected by the `skip_implicit_reboots` provider argument.",
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
		},
	},
}
