package nfsspaces

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsspace"
)

var filterConfig = frameworkfilter.Config{
	"id":                  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":               {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"description":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"status":              {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"max_capacity_bytes":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"used_capacity_bytes": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"nfs_spaces": schema.ListNestedAttribute{
			Description: "The returned list of NFS Spaces.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: nfsspace.SpaceAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
