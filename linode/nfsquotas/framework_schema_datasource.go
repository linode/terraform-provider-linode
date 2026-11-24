package nfsquotas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsquota"
)

var filterConfig = frameworkfilter.Config{
	"id":                 {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"path":               {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"is_root":            {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
	"status":             {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"quota_exceeded":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
	"max_capacity_bytes": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"max_file_count":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"space_id": schema.Int64Attribute{
			Description: "The ID of the NFS Space the filesystem belongs to.",
			Required:    true,
		},
		"filesystem_id": schema.Int64Attribute{
			Description: "The ID of the NFS Filesystem to list quotas for.",
			Required:    true,
		},
		"nfs_quotas": schema.ListNestedAttribute{
			Description: "The returned list of NFS Quotas.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: nfsquota.DataSourceSchema.Attributes,
			},
		},
		"filter": filterConfig.NestedAttributeSchema(),
	},
}
