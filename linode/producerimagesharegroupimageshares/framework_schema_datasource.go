package producerimagesharegroupimageshares

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"label": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var Attributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique ID assigned to this Image Share.",
		Computed:    true,
	},
	"label": schema.StringAttribute{
		Description: "The label of the Image Share.",
		Computed:    true,
	},
	"capabilities": schema.SetAttribute{
		Description: "The capabilities of the Image represented by the Image Share.",
		ElementType: types.StringType,
		Computed:    true,
	},
	"description": schema.StringAttribute{
		Description: "A description of the Image Share.",
		Computed:    true,
	},
	"created": schema.StringAttribute{
		Description: "When this Image Share was created.",
		Computed:    true,
	},
	"deprecated": schema.BoolAttribute{
		Description: "Whether or not this Image is deprecated.",
		Computed:    true,
	},
	"is_public": schema.BoolAttribute{
		Description: "True if the Image is public.",
		Computed:    true,
	},
	"image_sharing": schema.SingleNestedAttribute{
		Description: "Details about image sharing, including who the image is shared with and by.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"shared_with": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"sharegroup_count": schema.Int64Attribute{
						Description: "The number of sharegroups the private image is present in.",
						Computed:    true,
					},
					"sharegroup_list_url": schema.StringAttribute{
						Description: "The GET api url to view the sharegroups in which the image is shared.",
						Computed:    true,
					},
				},
			},
			"shared_by": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"sharegroup_id": schema.Int64Attribute{
						Description: "The sharegroup_id from the im_ImageShare row.",
						Computed:    true,
					},
					"sharegroup_uuid": schema.StringAttribute{
						Description: "The sharegroup_uuid from the im_ImageShare row.",
						Computed:    true,
					},
					"sharegroup_label": schema.StringAttribute{
						Description: "The label from the associated im_ImageShareGroup row.",
						Computed:    true,
					},
					"source_image_id": schema.StringAttribute{
						Description: "The image id of the base image (will only be shown to producers, will be None for consumers).",
						Computed:    true,
					},
				},
			},
		},
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

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"sharegroup_id": schema.Int64Attribute{
			Description: "The ID of the Image Share Group to list Image Shares for.",
			Required:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"image_shares": schema.ListNestedBlock{
			Description: "The returned list of Image Shares.",
			NestedObject: schema.NestedBlockObject{
				Attributes: Attributes,
			},
		},
	},
}
