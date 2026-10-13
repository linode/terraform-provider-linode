package sshkeys

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"created": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"id":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"ssh_key": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var sshKeyAttributes = map[string]schema.Attribute{
	"label": schema.StringAttribute{
		Description: "The label of the Linode SSH Key.",
		Computed:    true,
	},
	"ssh_key": schema.StringAttribute{
		Description: "The public SSH Key, which is used to authenticate to the root user of the Linodes you deploy.",
		Computed:    true,
	},
	"created": schema.StringAttribute{
		CustomType:  timetypes.RFC3339Type{},
		Description: "The date this key was added.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The ID of the SSH Key.",
		Computed:    true,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"sshkeys": schema.ListNestedAttribute{
			Description: "The returned list of SSH Keys.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: sshKeyAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
