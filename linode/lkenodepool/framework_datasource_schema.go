package lkenodepool

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "ID of the Pool to look up.",
			Required:    true,
		},
		"cluster_id": schema.Int64Attribute{
			Description: "ID of the Kubernetes cluster to look up.",
			Required:    true,
		},
		"autoscaler": schema.SingleNestedAttribute{
			Description: "When enabled, the number of nodes autoscales within the defined minimum and maximum values.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Description: "Whether autoscaling is enabled for this node pool.",
					Computed:    true,
				},
				"max": schema.Int64Attribute{
					Description: "The maximum number of nodes to autoscale to.",
					Computed:    true,
				},
				"min": schema.Int64Attribute{
					Description: "The minimum number of nodes to autoscale to.",
					Computed:    true,
				},
			},
		},
		"disk_encryption": schema.StringAttribute{
			Description: "Indicates the local disk encryption setting for this LKE node pool.",
			Computed:    true,
		},
		"disks": schema.ListNestedAttribute{
			Description: "This node pool's custom disk layout.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"size": schema.Int64Attribute{
						Description: "The size of this custom disk partition in MB.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "This custom disk partition's filesystem type.",
						Computed:    true,
					},
				},
			},
		},
		"firewall_id": schema.Int64Attribute{
			Description: "The ID of the Cloud Firewall assigned to this node pool.",
			Computed:    true,
		},
		"k8s_version": schema.StringAttribute{
			Description: "The Kubernetes version used for the worker nodes within this node pool.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The optional label defined for this node pool.",
			Computed:    true,
		},
		"labels": schema.MapAttribute{
			Description: "Key-value pairs added as labels to nodes in the node pool.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"node_count": schema.Int64Attribute{
			Description: "The number of nodes in the node pool.",
			Computed:    true,
		},
		"nodes": schema.ListNestedAttribute{
			Description: "Status information for the nodes that are members of this node pool. " +
				"If a Linode has not been provisioned for a given node slot, the 'instance_id' is null",
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The Node's ID.",
						Computed:    true,
					},
					"instance_id": schema.Int64Attribute{
						Description: "The Linode's ID. When no Linode is currently provisioned for this node, this is null.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "The creation status of this node.",
						Computed:    true,
					},
				},
			},
		},
		"tags": schema.ListAttribute{
			Description: "An array of tags applied to this object.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"taints": schema.ListNestedAttribute{
			Description: "Kubernetes taints added to nodes in the node pool.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"effect": schema.StringAttribute{
						Description: "The Kubernetes taint effect.",
						Computed:    true,
					},
					"key": schema.StringAttribute{
						Description: "The Kubernetes taint key.",
						Computed:    true,
					},
					"value": schema.StringAttribute{
						Description: "The Kubernetes taint value.",
						Computed:    true,
					},
				},
			},
		},
		"type": schema.StringAttribute{
			Description: "The Linode type for all of the nodes in the node pool.",
			Computed:    true,
		},
		"update_strategy": schema.StringAttribute{
			Description: "Determines when the worker nodes within this node pool upgrade to the latest selected Kubernetes version.",
			Computed:    true,
		},
	},
}
