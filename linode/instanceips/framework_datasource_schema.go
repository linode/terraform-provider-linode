package instanceips

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the Linode instance.",
			Required:    true,
		},
		"ipv4": schema.ObjectAttribute{
			Computed: true,
			AttributeTypes: map[string]attr.Type{
				"public":   types.ListType{ElemType: types.StringType},
				"private":  types.ListType{ElemType: types.StringType},
				"shared":   types.ListType{ElemType: types.StringType},
				"reserved": types.ListType{ElemType: types.StringType},
				"vpc":      types.ListType{ElemType: types.StringType},
			},
			Description: "The IPv4 addresses of the instance.",
		},
		"ipv6": schema.ObjectAttribute{
			Computed: true,
			AttributeTypes: map[string]attr.Type{
				"link_local": types.StringType,
				"slaac":      types.StringType,
				"global":     types.ListType{ElemType: types.StringType},
			},
			Description: "The IPv6 addresses of the instance.",
		},
	},
}
