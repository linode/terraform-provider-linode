package lkenodepool

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	linodesetplanmodifiers "github.com/linode/terraform-provider-linode/v3/linode/helper/setplanmodifiers"
)

var resourceSchema = schema.Schema{
	Version: 0,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "ID of the Node Pool.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"cluster_id": schema.Int64Attribute{
			Description: "The ID of the cluster to associate this node pool with.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label of the Node Pool.",
			Optional:    true,
			Default:     stringdefault.StaticString(""),
			Computed:    true,
		},
		"firewall_id": schema.Int64Attribute{
			Description: "The ID of the Firewall to attach to nodes in this node pool.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"node_count": schema.Int64Attribute{
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
				int64validator.AtLeastOneOf(path.MatchRoot("autoscaler")),
			},
			Description: "The number of nodes in the Node Pool.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"type": schema.StringAttribute{
			Description: "The type of node pool.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"disk_encryption": schema.StringAttribute{
			Description: "The disk encryption policy for nodes in this pool.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.InstanceDiskEncryptionEnabled),
					string(linodego.InstanceDiskEncryptionDisabled),
				),
			},
		},
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Default:     helper.EmptySetDefault(types.StringType),
			PlanModifiers: []planmodifier.Set{
				linodesetplanmodifiers.CaseInsensitiveSet(),
			},
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		},
		"nodes": schema.ListAttribute{
			Description: "A list of nodes in the node pool.",
			Computed:    true,
			ElementType: nodeObjectType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"labels": schema.MapAttribute{
			Description: "Key-value pairs added as labels to nodes in the node pool. " +
				"Labels help classify your nodes and to easily select subsets of objects.",
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Default:     helper.EmptyMapDefault(types.StringType),
		},
		"k8s_version": schema.StringAttribute{
			Description: "The k8s version of the nodes in this node pool. " +
				"For LKE enterprise only and may not currently available to all users.",
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"update_strategy": schema.StringAttribute{
			Description: "The strategy for updating the node pool k8s version. " +
				"For LKE enterprise only and may not currently available to all users.",
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.LKENodePoolOnRecycle),
					string(linodego.LKENodePoolRollingUpdate),
				),
			},
		},
	},
	Blocks: map[string]schema.Block{
		"autoscaler": schema.ListNestedBlock{
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
				listvalidator.AtLeastOneOf(path.MatchRoot("node_count")),
			},
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"min": schema.Int64Attribute{
						Optional:    true,
						Description: "The minimum number of nodes to automatically scale to.",
					},
					"max": schema.Int64Attribute{
						Optional:    true,
						Description: "The maximum number of nodes to automatically scale to.",
					},
				},
			},
		},

		"taint": schema.SetNestedBlock{
			Description: "Kubernetes taints to add to node pool nodes. Taints help control how " +
				"pods are scheduled onto nodes, specifically allowing them to repel certain pods.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"effect": schema.StringAttribute{
						Description: "The Kubernetes taint effect.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								string(linodego.LKENodePoolTaintEffectNoExecute),
								string(linodego.LKENodePoolTaintEffectNoSchedule),
								string(linodego.LKENodePoolTaintEffectPreferNoSchedule),
							),
						},
						Required: true,
					},
					"key": schema.StringAttribute{
						Description: "The Kubernetes taint key.",
						Required:    true,
					},
					"value": schema.StringAttribute{
						Description: "The Kubernetes taint value.",
						Required:    true,
					},
				},
			},
		},

		"isolation": schema.ListNestedBlock{
			Description: "Network isolation settings for the node pool. " +
				"Controls whether nodes have public IPv4/IPv6 addresses.",
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
			},
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"public_ipv4": schema.BoolAttribute{
						Description: "Whether nodes in this pool have public IPv4 addresses.",
						Optional:    true,
						Computed:    true,
					},
					"public_ipv6": schema.BoolAttribute{
						Description: "Whether nodes in this pool have public IPv6 addresses.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	},
}

var nodeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"instance_id": types.Int64Type,
		"status":      types.StringType,
	},
}
