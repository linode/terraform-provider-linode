package vpc

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

var VPCAttrs = map[string]schema.Attribute{
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
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: VPCAttrs,
}
