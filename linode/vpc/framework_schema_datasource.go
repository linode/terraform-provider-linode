package vpc

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var VPCAttrs = map[string]schema.Attribute{
	"id": schema.StringAttribute{
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
	"ipv6": schema.SetNestedAttribute{
		Description: "The IPv6 configuration of this VPC.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"range": schema.StringAttribute{
					Description: "The IPv6 range assigned to this VPC.",
					Computed:    true,
				},
			},
		},
	},
	"created": schema.StringAttribute{
		Description: "The date and time when the VPC was created.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Description: "The date and time when the VPC was updated.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: VPCAttrs,
}
