package monitoralertdefinitions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/monitoralertdefinition"
)

var filterConfig = frameworkfilter.Config{
	"id":                 {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"service_type":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"channel_ids":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"description":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"entity_ids":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"label":              {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"status":             {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"severity":           {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"type":               {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"has_more_resources": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
	"created":            {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"updated":            {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"created_by":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"updated_by":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"class":              {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"service_type": schema.StringAttribute{
			Description: "If provided, will restrict results to only Alert Definitions for the given service type.",
			Optional:    true,
		},
		"alert_definitions": schema.ListNestedAttribute{
			Description: "The returned list of Alert Definitions.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: monitoralertdefinition.AlertDefinitionAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
