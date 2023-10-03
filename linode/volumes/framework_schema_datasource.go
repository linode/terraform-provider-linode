package volumes

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/volume"
)

var filterConfig = frameworkfilter.Config{
	"label":           {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"tags":            {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"filesystem_path": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"hardware_type":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"linode_id":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"linode_label":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"region":          {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"status":          {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"size":            {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"created":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"updated":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"volumes": schema.ListNestedBlock{
			Description: "The return list of Volumes.",
			NestedObject: schema.NestedBlockObject{
				Attributes: volume.VolumeAttributes,
			},
		},
	},
}
