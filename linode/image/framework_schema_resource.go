package image

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var replicationObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"region": types.StringType,
		"status": types.StringType,
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the Linode image.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "A short description of the Image. Labels cannot contain special characters.",
			Required:    true,
		},
		"disk_id": schema.Int64Attribute{
			Description: "The ID of the Linode Disk that this Image will be created from.",
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
			Validators: []validator.Int64{
				int64validator.ConflictsWith(path.MatchRoot("file_path")),
				int64validator.AlsoRequires(path.MatchRoot("linode_id")),
				int64validator.ExactlyOneOf(
					path.MatchRoot("disk_id"),
					path.MatchRoot("file_path"),
				),
			},
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode that this Image will be created from.",
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
			Validators: []validator.Int64{
				int64validator.ConflictsWith(path.MatchRoot("file_path")),
				int64validator.AlsoRequires(path.MatchRoot("disk_id")),
				int64validator.ExactlyOneOf(
					path.MatchRoot("linode_id"),
					path.MatchRoot("file_path"),
				),
			},
		},
		"file_path": schema.StringAttribute{
			Description: "The name of the file to upload to this image.",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.ConflictsWith(
					path.MatchRoot("linode_id"),
					path.MatchRoot("disk_id"),
				),
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("file_path"),
					path.MatchRoot("linode_id"),
				),
				stringvalidator.AlsoRequires(path.MatchRoot("region")),
			},
		},
		"region": schema.StringAttribute{
			Description: "The region to upload to.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(
					path.MatchRoot("linode_id"),
					path.MatchRoot("disk_id"),
				),
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("region"),
					path.MatchRoot("linode_id"),
				),
				stringvalidator.AlsoRequires(path.MatchRoot("file_path")),
			},
		},
		"file_hash": schema.StringAttribute{
			Description: "The MD5 hash of the image file.",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"description": schema.StringAttribute{
			Description: "A detailed description of this Image.",
			Optional:    true,
		},
		"cloud_init": schema.BoolAttribute{
			Description: "Whether this image supports cloud-init.",
			Computed:    true,
			Default:     booldefault.StaticBool(false),
			Optional:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
		"created": schema.StringAttribute{
			Description: "When this Image was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			Description: "The name of the User who created this Image.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"deprecated": schema.BoolAttribute{
			Description: "Whether or not this Image is deprecated. Will only be True for deprecated public Images.",
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"is_public": schema.BoolAttribute{
			Description: "True if the Image is public.",
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"size": schema.Int64Attribute{
			Description: "The minimum size this Image needs to deploy. Size is in MB.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"type": schema.StringAttribute{
			Description: "How the Image was created. 'Manual' Images can be created at any time. 'Automatic' " +
				"images are created automatically from a deleted Linode.",
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"expiry": schema.StringAttribute{
			Description: "Only Images created automatically (from a deleted Linode; type=automatic) will expire.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"vendor": schema.StringAttribute{
			Description: "The upstream distribution vendor. Nil for private Images.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"status": schema.StringAttribute{
			Description: "The current status of this Image.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"capabilities": schema.ListAttribute{
			Description: "The capabilities of this Image.",
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"tags": schema.ListAttribute{
			Description: "The customized tags for the image.",
			Computed:    true,
			Optional:    true,
			ElementType: types.StringType,
		},
		"total_size": schema.Int64Attribute{
			Description: "The total size of the image in all available regions.",
			Computed:    true,
		},
		"replica_regions": schema.ListAttribute{
			Description: "A list of regions that customer wants to replicate this image in. " +
				"At least one available region is required and only core regions allowed. " +
				"Existing images in the regions not passed will be removed.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"replications": schema.ListAttribute{
			Description: "A list of image replications region and corresponding status.",
			Computed:    true,
			ElementType: replicationObjType,
		},
		"wait_for_replications": schema.BoolAttribute{
			Description: "Whether to wait for all image replications become `available`.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
	},
}
