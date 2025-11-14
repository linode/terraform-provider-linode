package consumerimagesharegroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"token_uuid": schema.StringAttribute{
			Description: "The UUID of the token that has been accepted into this Image Share Group.",
			Required:    true,
		},
		"id": schema.Int64Attribute{
			Description: "The id of the Image Share Group.",
			Computed:    true,
		},
		"uuid": schema.StringAttribute{
			Description: "The uuid of the Image Share Group.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the Image Share Group.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description of the Image Share Group.",
			Computed:    true,
		},
		"is_suspended": schema.BoolAttribute{
			Description: "Whether or not the Image Share Group is suspended..",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When the Image Share Group was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "When the Image Share Group was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
