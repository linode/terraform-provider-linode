package objendpoints

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"endpoint_type": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"region":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"s3_endpoint":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"country":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"region_label":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"endpoints": schema.ListNestedAttribute{
			Description: "The returned list of endpoints for the Linode Object Storage.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"endpoint_type": schema.StringAttribute{
						Description: "The type of `s3_endpoint` available to the active " +
							"user in this region.",
						Computed: true,
					},
					"region": schema.StringAttribute{
						Description: "The Akamai cloud computing region, represented by " +
							"its slug value.",
						Computed: true,
					},
					"s3_endpoint": schema.StringAttribute{
						Description: "Your s3 endpoint URL, based on the `endpoint_type` " +
							"and region. Shown as `null` if you haven't assigned an endpoint " +
							"for your user.",
						Computed: true,
					},
					"supported_key_classes": schema.ListAttribute{
						Description: "The key classes (`E`, `IA`) that can be created in this " +
							"endpoint's region. Shown as `null` if the account has no access " +
							"to Infrequent Access.",
						ElementType: types.StringType,
						Computed:    true,
					},
					"country": schema.StringAttribute{
						Description: "The country of the endpoint's region. Shown as `null` " +
							"if the account has no access to Infrequent Access.",
						Computed: true,
					},
					"region_label": schema.StringAttribute{
						Description: "The display label of the endpoint's region. Shown as " +
							"`null` if the account has no access to Infrequent Access.",
						Computed: true,
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
