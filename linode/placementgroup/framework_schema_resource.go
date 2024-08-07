package placementgroup

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var pgMemberObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"linode_id":    types.Int64Type,
		"is_compliant": types.BoolType,
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the Placement Group.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"label": schema.StringAttribute{
			Description: "The label of the Placement Group.",
			Required:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region of the Placement Group.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"placement_group_type": schema.StringAttribute{
			Description: "The placement group type for Linodes in this Placement Group.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			// Should we validate on the choices here or defer it to the API?
		},
		"placement_group_policy": schema.StringAttribute{
			Description: "Whether this Placement Group has a strict compliance policy.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("strict"),
			Validators: []validator.String{
				stringvalidator.OneOf("strict", "flexible"),
			},
		},
		"is_compliant": schema.BoolAttribute{
			Description: "Whether all Linodes in this Placement Group are currently compliant.",
			Computed:    true,
		},
		"members": schema.SetAttribute{
			Description: "A set of Linodes currently assigned to this Placement Group.",
			Computed:    true,
			ElementType: pgMemberObjectType,
		},
	},
}
