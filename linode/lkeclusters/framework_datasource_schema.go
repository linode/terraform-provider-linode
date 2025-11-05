package lkeclusters

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"k8s_version": {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"label":       {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"region":      {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"tags":        {APIFilterable: true, TypeFunc: helper.FilterTypeString},

	"created": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"updated": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"status":  {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"tier":    {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"lke_clusters": schema.ListNestedAttribute{
			Description: "The returned list of LKE clusters available on the account.",
			NestedObject: schema.NestedAttributeObject{
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
					"apl_enabled": schema.BoolAttribute{
						Description: "Enables the App Platform Layer for this cluster. " +
							"Note: v4beta only and may not currently be available to all users.",
						Computed: true,
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
					"tier": schema.StringAttribute{
						Computed:    true,
						Description: "The desired Kubernetes tier.",
					},
					"subnet_id": schema.Int64Attribute{
						Computed:    true,
						Description: "The ID of the VPC subnet to use for the Kubernetes cluster. This subnet must be dual stack (IPv4 and IPv6 should both be enabled). ",
					},
					"vpc_id": schema.Int64Attribute{
						Computed:    true,
						Description: "The ID of the VPC to use for the Kubernetes cluster.",
					},
					"stack_type": schema.StringAttribute{
						Computed:    true,
						Description: "The networking stack type of the Kubernetes cluster.",
					},
					"control_plane": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Defines settings for the Kubernetes Control Plane.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"high_availability": schema.BoolAttribute{
									Description: "Defines whether High Availability is enabled for the Control Plane Components of the cluster.",
									Computed:    true,
								},
								"audit_logs_enabled": schema.BoolAttribute{
									Description: "Enables audit logs on the cluster's control plane.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
