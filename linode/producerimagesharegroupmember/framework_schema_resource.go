package producerimagesharegroupmember

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"sharegroup_id": schema.Int64Attribute{
			Description: "The ID of the Image Share Group the member belongs to.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"token": schema.StringAttribute{
			Description: "The one-time-use token provided by the prospective member.",
			Required:    true,
			Sensitive:   true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the member.",
			Required:    true,
		},
		"token_uuid": schema.StringAttribute{
			Description: "The UUID of member's token.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The status of the member.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this member was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "When this member was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"expiry": schema.StringAttribute{
			Description: "When this member will expire.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
