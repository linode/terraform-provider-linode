package sshkeys

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/sshkey"
)

var filterConfig = frameworkfilter.Config{
	"created": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"id":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"ssh_key": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
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
		"sshkeys": schema.ListNestedBlock{
			Description: "The returned list of SSH Keys.",
			NestedObject: schema.NestedBlockObject{
				Attributes: sshkey.SSHKeyAttributes,
			},
		},
	},
}
