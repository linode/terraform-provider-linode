package consumerimagesharegrouptoken

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"valid_for_sharegroup_uuid": schema.StringAttribute{
			Description: "The UUID of the Image Share Group this token is for.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label of the token.",
			Optional:    true,
		},
		"token": schema.StringAttribute{
			Description: "The one-time-use token to be provided to the Share Group Producer.",
			Computed:    true,
			Sensitive:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"token_uuid": schema.StringAttribute{
			Description: "The UUID of the token.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The status of the token.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this token was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "When this token was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"expiry": schema.StringAttribute{
			Description: "When this token will expire.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"sharegroup_uuid": schema.StringAttribute{
			Description: "The UUID of the Image Share Group this token is for.",
			Computed:    true,
		},
		"sharegroup_label": schema.StringAttribute{
			Description: "The label of the Image Share Group this token is for.",
			Computed:    true,
		},
	},
}
