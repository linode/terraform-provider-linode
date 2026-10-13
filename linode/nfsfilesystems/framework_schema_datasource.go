package nfsfilesystems

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsfilesystem"
)

var filterConfig = frameworkfilter.Config{
	"id":                 {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":              {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"region":             {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"status":             {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"max_capacity_bytes": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"space_id": schema.Int64Attribute{
			Description: "The ID of the NFS Space to list filesystems for.",
			Required:    true,
		},
		"nfs_filesystems": schema.ListNestedAttribute{
			Description: "The returned list of NFS Filesystems.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: nfsfilesystem.FilesystemAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
