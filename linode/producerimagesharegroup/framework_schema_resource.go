package producerimagesharegroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var imageShareGroupImage = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of an image to share in this Image Share Group.",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label for the im_ImageShare row associated with this shared image.",
			Optional:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description for the im_ImageShare row associated with this shared image.",
			Optional:    true,
		},
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The ID of the Image Share Group.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"uuid": schema.StringAttribute{
			Description: "The UUID of the Image Share Group.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label of the Image Share Group.",
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: "The label of the Image Share Group.",
			Optional:    true,
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
		"images": schema.ListNestedAttribute{
			Description:  "The images to be shared using this Image Share Group.",
			Optional:     true,
			NestedObject: imageShareGroupImage,
		},
	},
}
