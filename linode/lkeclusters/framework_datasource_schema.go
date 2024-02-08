package lkeclusters

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"k8s_version": {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"label":       {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"region":      {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"tags":        {APIFilterable: true, TypeFunc: helper.FilterTypeString},

	"created": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"updated": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"status":  {APIFilterable: false, TypeFunc: helper.FilterTypeString},
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
		"lke_clusters": schema.ListNestedBlock{
			Description: "The returned list of LKE clusters available on the account.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "This Kubernetes cluster’s unique ID.",
					},
					"created": schema.StringAttribute{
						Computed:    true,
						Description: "When this Kubernetes cluster was created.",
					},
					"updated": schema.StringAttribute{
						Computed:    true,
						Description: "When this Kubernetes cluster was updated.",
					},
					"label": schema.StringAttribute{
						Computed:    true,
						Description: "The unique label for the cluster.",
					},
					"k8s_version": schema.StringAttribute{
						Computed: true,
						Description: "The desired Kubernetes version for this Kubernetes cluster in the format of <major>.<minor>. " +
							"The latest supported patch version will be deployed.",
					},
					"tags": schema.SetAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
					},
					"region": schema.StringAttribute{
						Computed:    true,
						Description: "This cluster's location.",
					},
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "The status of the cluster.",
					},
				},
				Blocks: map[string]schema.Block{
					"control_plane": schema.SingleNestedBlock{
						Description: "Defines settings for the Kubernetes Control Plane.",
						Attributes: map[string]schema.Attribute{
							"high_availability": schema.BoolAttribute{
								Description: "Defines whether High Availability is enabled for the Control Plane Components of the cluster.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	},
}
