package lke

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID of this LKE Cluster.",
			Required:    true,
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
		"api_endpoints": schema.ListAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "The API endpoints for the cluster.",
		},
		"kubeconfig": schema.StringAttribute{
			Computed:    true,
			Sensitive:   true,
			Description: "The Base64-encoded Kubeconfig for the cluster.",
		},
		"dashboard_url": schema.StringAttribute{
			Computed:    true,
			Description: "The dashboard URL of the cluster.",
		},
		"status": schema.StringAttribute{
			Computed:    true,
			Description: "The status of the cluster.",
		},
		"created": schema.StringAttribute{
			Computed:    true,
			Description: "When this Kubernetes cluster was created.",
		},
		"updated": schema.StringAttribute{
			Computed:    true,
			Description: "When this Kubernetes cluster was updated.",
		},
	},
	Blocks: map[string]schema.Block{
		"control_plane": schema.ListNestedBlock{
			Description: "Defines settings for the Kubernetes Control Plane.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"high_availability": schema.BoolAttribute{
						Description: "Defines whether High Availability is enabled for the Control Plane Components of the cluster.",
						Computed:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"acl": schema.ListNestedBlock{
						Description: "The ACL configuration for an LKE cluster's control plane.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "The default policy. A value of true means a default policy of DENY. A value of false means default policy of ALLOW.",
									Computed:    true,
								},
							},
							Blocks: map[string]schema.Block{
								"addresses": schema.ListNestedBlock{
									Description: "A list of ip addresses allowed.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"ipv4": schema.SetAttribute{
												Description: "A set of individual ipv4 addresses or CIDRs allowed.",
												Computed:    true,
												ElementType: types.StringType,
											},
											"ipv6": schema.SetAttribute{
												Description: "A set of individual ipv6 addresses or CIDRs allowed.",
												Computed:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"pools": schema.ListNestedBlock{
			Description: "All active Node Pools on the cluster.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "The ID of the Node Pool.",
					},
					"count": schema.Int64Attribute{
						Computed:    true,
						Description: "The number of nodes in the Node Pool.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "A Linode Type for all of the nodes in the Node Pool.",
					},
					"tags": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
					},
					"disk_encryption": schema.StringAttribute{
						Computed: true,
						Description: "The disk encryption policy for the nodes in this pool. " +
							"NOTE: Disk encryption may not currently be available to all users.",
					},
				},
				Blocks: map[string]schema.Block{
					"nodes": schema.ListNestedBlock{
						Description: "The nodes in the node pool.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "The ID of the node.",
								},
								"instance_id": schema.Int64Attribute{
									Computed:    true,
									Description: "The ID of the underlying Linode instance.",
								},
								"status": schema.StringAttribute{
									Computed:    true,
									Description: "The status of the node.",
								},
							},
						},
					},
					"autoscaler": schema.ListNestedBlock{
						Description: "When specified, the number of nodes autoscales within " +
							"the defined minimum and maximum values.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"min": schema.Int64Attribute{
									Description: "The minimum number of nodes to autoscale to. Defaults to the Node Pool’s count.",
									Computed:    true,
								},
								"max": schema.Int64Attribute{
									Description: "The maximum number of nodes to autoscale to. Defaults to the Node Pool’s count.",
									Computed:    true,
								},
								"enabled": schema.BoolAttribute{
									Description: "Whether autoscaling is enabled for this Node Pool. Defaults to false.",
									Computed:    true,
								},
							},
						},
					},
					"disks": schema.ListNestedBlock{
						Description: "This Node Pool’s custom disk layout.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"size": schema.Int64Attribute{
									Description: "The size of this custom disk partition in MB.",
									Computed:    true,
								},
								"type": schema.StringAttribute{
									Description: "This custom disk partition’s filesystem type.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	},
}
