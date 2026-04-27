package monitoralertdefinitionentities

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"label": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"type":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"url":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var entityAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Computed:    true,
		Description: "The unique identifier for this entity.",
	},
	"label": schema.StringAttribute{
		Computed:    true,
		Description: "The label of this entity.",
	},
	"url": schema.StringAttribute{
		Computed:    true,
		Description: "The URL for this entity.",
	},
	"type": schema.StringAttribute{
		Computed:    true,
		Description: "The type of this entity.",
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"service_type": schema.StringAttribute{
			Description: "The service type for the alert definition.",
			Required:    true,
		},
		"alert_id": schema.Int64Attribute{
			Description: "The unique identifier for the alert definition.",
			Required:    true,
		},
		"order_by": schema.StringAttribute{
			Description: "The attribute to order the results by.",
			Optional:    true,
		},
		"order": schema.StringAttribute{
			Description: "The order in which results are returned (asc or desc).",
			Optional:    true,
		},
		"entities": schema.ListNestedAttribute{
			Description: "The returned list of entities associated with the alert definition.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: entityAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
