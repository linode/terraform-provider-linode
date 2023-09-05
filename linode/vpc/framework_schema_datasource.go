package vpc

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
	"github.com/linode/terraform-provider-linode/linode/vpcsubnet"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The id of the VPC.",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the VPC.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "The user-defined description of this VPC.",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region of the VPC.",
			Computed:    true,
		},
		"subnets": schema.ListNestedAttribute{
			Description: "A list of subnets under this VPC.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: vpcsubnet.VPCSubnetAttrs,
			},
		},
		"created": schema.StringAttribute{
			Description: "The date and time when the VPC was created.",
			Computed:    true,
			CustomType:  customtypes.RFC3339TimeStringType{},
		},
		"updated": schema.StringAttribute{
			Description: "The date and time when the VPC was updated.",
			Computed:    true,
			CustomType:  customtypes.RFC3339TimeStringType{},
		},
	},
}
