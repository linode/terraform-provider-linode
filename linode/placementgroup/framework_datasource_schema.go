package placementgroup

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
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
		"affinity_type": schema.StringAttribute{
			Description: "The affinity policy for Linodes in a placement group",
			Computed:    true,
		},
		"is_compliant": schema.BoolAttribute{
			Description: "Whether all Linodes in this group are currently compliant with the group's affinity policy.",
			Computed:    true,
		},
		"is_strict": schema.BoolAttribute{
			Description: "Whether Linodes must be able to become compliant during assignment.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"members": schema.ListNestedBlock{
			Description: "A list of Linodes assigned to a placement group.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"linode_id": schema.Int64Attribute{
						Description: "The ID of the Linode.",
						Computed:    true,
					},
					"is_compliant": schema.BoolAttribute{
						Description: "Whether this Linode is currently compliant with the group's affinity policy.",
						Computed:    true,
					},
				},
			},
		},
	},
}
