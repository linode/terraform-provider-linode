package image

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ImageAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique ID assigned to this Image.",
		Required:    true,
	},
	"label": schema.StringAttribute{
		Description: "A short description of the Image. Labels cannot contain special characters.",
		Computed:    true,
	},
	"capabilities": schema.SetAttribute{
		Description: "The capabilities of this Image.",
		ElementType: types.StringType,
		Computed:    true,
	},
	"description": schema.StringAttribute{
		Description: "A detailed description of this Image.",
		Computed:    true,
	},
	"created": schema.StringAttribute{
		Description: "When this Image was created.",
		Computed:    true,
	},
	"created_by": schema.StringAttribute{
		Description: "The name of the User who created this Image.",
		Computed:    true,
	},
	"deprecated": schema.BoolAttribute{
		Description: "Whether or not this Image is deprecated. Will only be True for deprecated public Images.",
		Computed:    true,
	},
	"is_public": schema.BoolAttribute{
		Description: "True if the Image is public.",
		Computed:    true,
	},
	"size": schema.Int64Attribute{
		Description: "The minimum size this Image needs to deploy. Size is in MB.",
		Computed:    true,
	},
	"status": schema.StringAttribute{
		Description: "The current status of this Image.",
		Computed:    true,
	},
	"type": schema.StringAttribute{
		Description: "How the Image was created. 'Manual' Images can be created at any time. 'Automatic' " +
			"images are created automatically from a deleted Linode.",
		Computed: true,
	},
	"expiry": schema.StringAttribute{
		Description: "Only Images created automatically (from a deleted Linode; type=automatic) will expire.",
		Computed:    true,
	},
	"vendor": schema.StringAttribute{
		Description: "The upstream distribution vendor. Nil for private Images.",
		Computed:    true,
	},
	"tags": schema.ListAttribute{
		Description: "The customized tags for the image.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"total_size": schema.Int64Attribute{
		Description: "The total size of the image in all available regions.",
		Computed:    true,
	},
}

var ReplicationsBlock = map[string]schema.Block{
	"replications": schema.ListNestedBlock{
		Description: "A list of image replications region and corresponding status.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					Description: "The region of an image replica.",
					Computed:    true,
				},
				"status": schema.StringAttribute{
					Description: "The status of an image replica.",
					Computed:    true,
				},
			},
		},
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: ImageAttributes,
	Blocks:     ReplicationsBlock,
}
