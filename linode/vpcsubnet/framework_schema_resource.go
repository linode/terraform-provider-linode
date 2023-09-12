package vpcsubnet

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

var VPCSubnetResourceAttrs = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The id of the VPC Subnet.",
		Computed:    true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	},
	"vpc_id": schema.Int64Attribute{
		Description: "The id of the parent VPC for this VPC Subnet",
		Required:    true,
	},
	"label": schema.StringAttribute{
		Description: "The label of the VPC subnet.",
		Required:    true,
	},
	"ipv4": schema.StringAttribute{
		Description: "The IPv4 range of this subnet in CIDR format.",
		Required:    true,
	},
	"linodes": schema.ListAttribute{
		ElementType: types.Int64Type,
		Description: "A list of Linode IDs that added to this subnet.",
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
	},
	"created": schema.StringAttribute{
		Description: "The date and time when the VPC Subnet was created.",
		Computed:    true,
		CustomType:  customtypes.RFC3339TimeStringType{},
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	},
	"updated": schema.StringAttribute{
		Description: "The date and time when the VPC Subnet was updated.",
		Computed:    true,
		CustomType:  customtypes.RFC3339TimeStringType{},
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: VPCSubnetResourceAttrs,
}
