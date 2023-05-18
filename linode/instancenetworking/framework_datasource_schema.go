package instancenetworking

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var networkObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"address":     types.StringType,
		"gateway":     types.StringType,
		"prefix":      types.Int64Type,
		"rdns":        types.StringType,
		"region":      types.StringType,
		"subnet_mask": types.StringType,
		"type":        types.StringType,
	},
}

var ipv4ObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"private":  types.ListType{ElemType: networkObjectType},
		"public":   types.ListType{ElemType: networkObjectType},
		"reserved": types.ListType{ElemType: networkObjectType},
		"shared":   types.ListType{ElemType: networkObjectType},
	},
}

var globalObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"prefix":       types.Int64Type,
		"range":        types.StringType,
		"region":       types.StringType,
		"route_target": types.StringType,
	},
}

var ipv6ObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"global":     types.ListType{ElemType: globalObjectType},
		"link_local": types.ObjectType{AttrTypes: networkObjectType.AttrTypes},
		"slaac":      types.ObjectType{AttrTypes: networkObjectType.AttrTypes},
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode for network info.",
			Required:    true,
		},
		"ipv4": schema.ListAttribute{
			Description: "Information about this Linode’s IPv4 addresses.",
			Computed:    true,
			ElementType: ipv4ObjectType,
		},
		"ipv6": schema.ListAttribute{
			Description: "Information about this Linode’s IPv6 addresses.",
			Computed:    true,
			ElementType: ipv6ObjectType,
		},
		"id": schema.StringAttribute{
			Description: "Unique identifier for this DataSource.",
			Computed:    true,
		},
	},
}
