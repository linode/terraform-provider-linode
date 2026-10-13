package vpcdefaultranges

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"default_ipv4_ranges": schema.SetAttribute{
			ElementType: types.StringType,
			Description: "The default IPv4 address ranges for all VPCs.",
			Computed:    true,
		},
		"forbidden_ipv4_ranges": schema.SetAttribute{
			ElementType: types.StringType,
			Description: "List of the IPv4 address ranges that are not allowed to be used.",
			Computed:    true,
		},
	},
}
