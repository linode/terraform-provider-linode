package producerimagesharegroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var Attributes = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the Image Share Group.",
		Required:    true,
	},
	"uuid": schema.StringAttribute{
		Description: "The UUID of the Image Share Group.",
		Computed:    true,
	},
	"label": schema.StringAttribute{
		Description: "The label of the Image Share Group.",
		Computed:    true,
	},
	"description": schema.StringAttribute{
		Description: "The label of the Image Share Group.",
		Computed:    true,
	},
	"is_suspended": schema.BoolAttribute{
		Description: "Whether or not the Image Share Group is suspended.",
		Computed:    true,
	},
	"images_count": schema.Int64Attribute{
		Description: "The number of images in the Image Share Group.",
		Computed:    true,
	},
	"members_count": schema.Int64Attribute{
		Description: "The number of members in the Image Share Group.",
		Computed:    true,
	},
	"created": schema.StringAttribute{
		Description: "When this Image Share Group was created.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Description: "When this Image Share Group was last updated.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"expiry": schema.StringAttribute{
		Description: "When this Image Share Group will expire.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: Attributes,
}
