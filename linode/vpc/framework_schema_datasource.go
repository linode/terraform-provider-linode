package vpc

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

var subnetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":      types.Int64Type,
		"label":   types.StringType,
		"ipv4":    types.StringType,
		"linodes": types.ListType{ElemType: types.Int64Type},
		"created": types.StringType,
		"updated": types.StringType,
	},
}

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
		"subnets": schema.ListAttribute{
			Description: "A list of subnets under this VPC.",
			Computed:    true,
			ElementType: subnetObjectType,
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
