package consumerimagesharegrouptoken

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var Attributes = map[string]schema.Attribute{
	"token_uuid": schema.StringAttribute{
		Description: "The UUID of the token.",
		Required:    true,
	},
	"label": schema.StringAttribute{
		Description: "The label of the token.",
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
	"valid_for_sharegroup_uuid": schema.StringAttribute{
		Description: "The UUID of the Image Share Group this token is for.",
		Computed:    true,
	},
	"sharegroup_uuid": schema.StringAttribute{
		Description: "The UUID of the Image Share Group this token is for.",
		Computed:    true,
	},
	"sharegroup_label": schema.StringAttribute{
		Description: "The label of the Image Share Group this token is for.",
		Computed:    true,
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: Attributes,
}
