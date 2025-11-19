package placementgroup

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var pgMigrationInstanceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"linode_id": types.Int64Type,
	},
}

var DataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The ID of the placement group.",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the placement group.",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region that the placement group is in.",
			Computed:    true,
		},
		"placement_group_type": schema.StringAttribute{
			Description: "The placement group type for Linodes in a placement group",
			Computed:    true,
		},
		"is_compliant": schema.BoolAttribute{
			Description: "Whether all Linodes in this group are currently compliant with the group's placement group type.",
			Computed:    true,
		},
		"placement_group_policy": schema.StringAttribute{
			Description: "Whether Linodes must be able to become compliant during assignment.",
			Computed:    true,
		},
		"members": schema.SetNestedAttribute{
			Description: "A list of Linodes assigned to a placement group.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"linode_id": schema.Int64Attribute{
						Description: "The ID of the Linode.",
						Computed:    true,
					},
					"is_compliant": schema.BoolAttribute{
						Description: "Whether this Linode is currently compliant with the group's placement group type.",
						Computed:    true,
					},
				},
			},
		},
		"migrations": schema.SingleNestedAttribute{
			Description: "An object representing migrations associated with the placement group. Contains inbound and outbound migration lists.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"inbound": schema.ListAttribute{
					Description: "The compute instances being migrated into the placement group.",
					Computed:    true,
					ElementType: pgMigrationInstanceObjectType,
				},
				"outbound": schema.ListAttribute{
					Description: "The compute instances being migrated out of the placement group.",
					Computed:    true,
					ElementType: pgMigrationInstanceObjectType,
				},
			},
		},
	},
}
