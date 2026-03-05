package prefixlists

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"name":       {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"visibility": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
}

var prefixListAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "The name of the Prefix List.",
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
}

var frameworkDatasourceSchema = schema.Schema{
	Description: "Provides a list of Linode Prefix Lists matching optional filters.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"prefix_lists": schema.ListNestedAttribute{
			Description: "The returned list of prefix lists.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: prefixListAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
