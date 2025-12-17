package locks

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":           {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"lock_type":    {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"entity_id":    {APIFilterable: false, TypeFunc: helper.FilterTypeInt},
	"entity_type":  {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"entity_label": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var lockAttributes = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The unique ID of the Lock.",
		Computed:    true,
	},
	"entity_id": schema.Int64Attribute{
		Description: "The ID of the locked entity.",
		Computed:    true,
	},
	"entity_type": schema.StringAttribute{
		Description: "The type of the locked entity.",
		Computed:    true,
	},
	"lock_type": schema.StringAttribute{
		Description: "The type of lock. Possible values: 'cannot_delete' (prevents deletion, rebuild, and transfer) or 'cannot_delete_with_subresources' (also prevents deletion of subresources).",
		Computed:    true,
	},
	"entity_label": schema.StringAttribute{
		Description: "The label of the locked entity.",
		Computed:    true,
	},
	"entity_url": schema.StringAttribute{
		Description: "The URL of the locked entity.",
		Computed:    true,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Description: "Provides information about Linode Locks that match a set of filters.",
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
		"locks": schema.ListNestedBlock{
			Description: "The returned list of Locks.",
			NestedObject: schema.NestedBlockObject{
				Attributes: lockAttributes,
			},
		},
	},
}
