package vpcsubnet

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The id of the VPC Subnet.",
			Required:    true,
		},
		"vpc_id": schema.Int64Attribute{
			Description: "The id of the parent VPC for this VPC Subnet",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the VPC Subnet.",
			Computed:    true,
		},
		"ipv4": schema.StringAttribute{
			Description: "The IPv4 range of this subnet in CIDR format.",
			Computed:    true,
		},
		"linodes": schema.ListAttribute{
			ElementType: types.Int64Type,
			Description: "A list of Linode IDs that added to this subnet.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "The date and time when the VPC Subnet was created.",
			Computed:    true,
			CustomType:  customtypes.RFC3339TimeStringType{},
		},
		"updated": schema.StringAttribute{
			Description: "The date and time when the VPC Subnet was updated.",
			Computed:    true,
			CustomType:  customtypes.RFC3339TimeStringType{},
		},
	},
}
