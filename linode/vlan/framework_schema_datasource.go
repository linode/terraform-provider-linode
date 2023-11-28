package vlan

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"label":  {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"region": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},

	"linodes": {APIFilterable: false, TypeFunc: helper.FilterTypeInt},
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

		"vlans": schema.ListNestedBlock{
			Description:  "The returned list of VLANs.",
			NestedObject: vlanSchema,
		},
	},
}

var vlanSchema = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description: "The unique label of this VLAN.",
			Computed:    true,
		},
		"linodes": schema.SetAttribute{
			Description: "The Linodes currently attached to this VLAN.",
			ElementType: types.Int64Type,
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region this VLAN is located in.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this VLAN was created.",
			Computed:    true,
		},
	},
}
