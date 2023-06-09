package databases

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"label":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"region":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"status":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"type":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"version": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},

	"engine":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"allow_list":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"cluster_size":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"created":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"encrypted":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
	"host_primary":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"host_secondary": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"id":             {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"instance_uri":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"updated":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order_by": filterConfig.OrderBySchema(),
		"order":    filterConfig.OrderSchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"databases": schema.ListNestedBlock{
			Description: "The returned list of databases.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "A unique ID that can be used to identify and reference the Managed Database.",
						Computed:    true,
					},
					"allow_list": schema.SetAttribute{
						Description: "A list of IP addresses that can access the Managed Database. " +
							"Each item can be a single IP address or a range in CIDR format.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"cluster_size": schema.Int64Attribute{
						Description: "The number of Linode Instance nodes deployed to the Managed Database.",
						Computed:    true,
					},
					"created": schema.StringAttribute{
						Description: "When this Managed Database was created.",
						Computed:    true,
					},
					"encrypted": schema.BoolAttribute{
						Description: "Whether the Managed Databases is encrypted.",
						Computed:    true,
					},
					"engine": schema.StringAttribute{
						Description: "The Managed Database engine type.",
						Computed:    true,
					},
					"host_primary": schema.StringAttribute{
						Description: "The primary host for the Managed Database.",
						Computed:    true,
					},
					"host_secondary": schema.StringAttribute{
						Description: "The secondary/private host for the Managed Database.",
						Computed:    true,
					},
					"instance_uri": schema.StringAttribute{
						Description: "he API route for the database instance.",
						Computed:    true,
					},
					"label": schema.StringAttribute{
						Description: "A unique, user-defined string referring to the Managed Database.",
						Computed:    true,
					},
					"region": schema.StringAttribute{
						Description: "The Region ID for the Managed Database.",
						Computed:    true,
					},
					"replication_type": schema.StringAttribute{
						Description: "The replication method used for the Managed Database.",
						Computed:    true,
					},
					"ssl_connection": schema.BoolAttribute{
						Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "The operating status of the Managed Database.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The Linode Instance type used by the Managed Database for its nodes.",
						Computed:    true,
					},
					"updated": schema.StringAttribute{
						Description: "When this Managed Database was last updated.",
						Computed:    true,
					},
					"version": schema.StringAttribute{
						Description: "The Managed Database engine version.",
						Computed:    true,
					},
				},
			},
		},
	},
}
