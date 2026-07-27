package objglobalquotas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"quota_id":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"quota_name":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"description":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"quota_limit":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"resource_metric": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"quota_type":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"has_usage":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"quotas": schema.ListNestedAttribute{
			Description: "The returned list of Object Storage global quotas.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"quota_id": schema.StringAttribute{
						Description: "The ID of the Object Storage global quota.",
						Required:    true,
					},
					"quota_name": schema.StringAttribute{
						Description: "The name of the Object Storage global quota.",
						Computed:    true,
					},
					"description": schema.StringAttribute{
						Description: "The description of the Object Storage global quota.",
						Computed:    true,
					},
					"quota_limit": schema.Int64Attribute{
						Description: "The maximum quantity of the `resource_metric` allowed by the quota.",
						Computed:    true,
					},
					"resource_metric": schema.StringAttribute{
						Description: "The specific Object Storage resource for the quota.",
						Computed:    true,
					},
					"quota_type": schema.StringAttribute{
						Description: "The type of the Object Storage global quota.",
						Computed:    true,
					},
					"has_usage": schema.BoolAttribute{
						Description: "Whether usage data is available for this Object Storage global quota.",
						Computed:    true,
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
