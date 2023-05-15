package stackscript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var udfObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"label":   types.StringType,
		"name":    types.StringType,
		"example": types.StringType,
		"one_of":  types.StringType,
		"many_of": types.StringType,
		"default": types.StringType,
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The StackScript's unique ID.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"label": schema.StringAttribute{
			Description: "The StackScript's label is for display purposes only.",
			Required:    true,
		},
		"script": schema.StringAttribute{
			Description: "The script to execute when provisioning a new Linode with this StackScript.",
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: "A description for the StackScript.",
			Required:    true,
		},
		"rev_note": schema.StringAttribute{
			Description: "This field allows you to add notes for the set of revisions made to this StackScript.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"is_public": schema.BoolAttribute{
			Description: "This determines whether other users can use your StackScript. Once a StackScript is " +
				"made public, it cannot be made private.",
			Default:  booldefault.StaticBool(false),
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplaceIf(
					func(
						ctx context.Context,
						request planmodifier.BoolRequest,
						response *boolplanmodifier.RequiresReplaceIfFuncResponse,
					) {
						if !request.PlanValue.ValueBool() && request.StateValue.ValueBool() {
							response.RequiresReplace = true
						}
					},
					"Replaces should only be required when attempting to make a StackScript private.",
					"Replaces should only be required when attempting to make a StackScript private.",
				),
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"images": schema.ListAttribute{
			Description: "An array of Image IDs representing the Images that this StackScript is compatible for " +
				"deploying with.",
			ElementType: types.StringType,
			Required:    true,
		},

		"deployments_active": schema.Int64Attribute{
			Description: "Count of currently active, deployed Linodes created from this StackScript.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"user_gravatar_id": schema.StringAttribute{
			Description: "The Gravatar ID for the User who created the StackScript.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"deployments_total": schema.Int64Attribute{
			Description: "The total number of times this StackScript has been deployed.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"username": schema.StringAttribute{
			Description: "The User who created the StackScript.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created": schema.StringAttribute{
			Description: "The date this StackScript was created.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "The date this StackScript was updated.",
			Computed:    true,
		},

		"user_defined_fields": schema.ListAttribute{
			Description: "This is a list of fields defined with a special syntax inside this " +
				"StackScript that allow for supplying customized parameters during deployment.",
			Computed:    true,
			ElementType: udfObjectType,
		},
	},
}
