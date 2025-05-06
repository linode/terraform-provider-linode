package objquotas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"quota_name":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"s3_endpoint":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"quota_id":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"endpoint_type":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"description":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"quota_limit":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"resource_metric": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"quotas": schema.ListNestedBlock{
			Description: "The returned list of Object Storage quotas.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"quota_id": schema.StringAttribute{
						Description: "The ID of the Object Storage quota.",
						Required:    true,
					},
					"quota_name": schema.StringAttribute{
						Description: "The name of the Object Storage quota. ",
						Computed:    true,
					},
					"endpoint_type": schema.StringAttribute{
						Description: "The type of the S3 endpoint of the Object Storage.",
						Computed:    true,
					},
					"s3_endpoint": schema.StringAttribute{
						Description: "The S3 endpoint URL of the Object Storage, based on the `endpoint_type` and `region`.",
						Computed:    true,
					},
					"description": schema.StringAttribute{
						Description: "The description of the Object Storage quota.",
						Computed:    true,
					},
					"quota_limit": schema.Int64Attribute{
						Description: "The maximum quantity of the `resource_metric` allowed by the quota.",
						Computed:    true,
					},
					"resource_metric": schema.StringAttribute{
						Description: "The specific Object Storage resource for the quota. ",
						Computed:    true,
					},
				},
			},
		},
	},
}
