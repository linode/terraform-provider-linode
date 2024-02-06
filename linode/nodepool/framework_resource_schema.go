package nodepool

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Version: 1,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Compound ID of the Node Pool. The ID format is <clusterID>:<PoolID>.",
			Computed:    true,
		},
		"pool_id": schema.Int64Attribute{
			Description: "The ID of the Node Pool",
			Computed:    true,
		},
		"cluster_id": schema.Int64Attribute{
			Description: "The ID of the cluster to associate this node pool with.",
			Required:    true,
		},
		"node_count": schema.Int64Attribute{
			Description: "The number of nodes in the node pool.",
			Required:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of node pool.",
			Required:    true,
		},
		"tags": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		},
		"autoscaler": schema.ObjectAttribute{
			Description: "When specified, the number of nodes autoscales within the defined minimum and maximum values.",
			Optional:    true,
			AttributeTypes: map[string]attr.Type{
				"min": types.Int64Type,
				"max": types.Int64Type,
			},
		},
		"nodes": schema.ListAttribute{
			Description: "A list of nodes in the node pool.",
			Computed:    true,
			ElementType: nodeObjectType,
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
