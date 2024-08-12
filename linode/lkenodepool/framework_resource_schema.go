package lkenodepool

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	linodeplanmodifiers "github.com/linode/terraform-provider-linode/v2/linode/helper/planmodifiers"
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
			Description: "The disk encryption policy for nodes in this pool. " +
				"NOTE: Disk encryption may not currently be available to all users.",
			Computed: true,
		},
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Default:     helper.EmptySetDefault(types.StringType),
			PlanModifiers: []planmodifier.Set{
				linodeplanmodifiers.CaseInsensitiveSet(),
			},
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		},
		"nodes": schema.ListAttribute{
			Description: "A list of nodes in the node pool.",
			Computed:    true,
			ElementType: nodeObjectType,
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
						Required: true,
					},
					"max": schema.Int64Attribute{
						Required: true,
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
