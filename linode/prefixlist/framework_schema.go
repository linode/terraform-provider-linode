package prefixlist

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDatasourceSchema = schema.Schema{
	Description: "Provides details about a specific Linode Prefix List.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of the Prefix List.",
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the Prefix List (e.g. pl:system:object-storage:us-iad).",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "A description for this Prefix List.",
			Computed:    true,
		},
		"visibility": schema.StringAttribute{
			Description: "The visibility of the Prefix List (e.g. system, user).",
			Computed:    true,
		},
		"source_prefixlist_id": schema.Int64Attribute{
			Description: "If this Prefix List was cloned, the ID of the source Prefix List.",
			Computed:    true,
		},
		"ipv4": schema.ListAttribute{
			Description: "The IPv4 prefixes in this Prefix List.",
			ElementType: types.StringType,
			Computed:    true,
		},
		"ipv6": schema.ListAttribute{
			Description: "The IPv6 prefixes in this Prefix List.",
			ElementType: types.StringType,
			Computed:    true,
		},
		"version": schema.Int64Attribute{
			Description: "The version of this Prefix List.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When the Prefix List was created.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "When the Prefix List was last updated.",
			Computed:    true,
		},
	},
}
