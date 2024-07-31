package instancedisk

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the Linode disk.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The Disk;s label is for display purposes only.",
			Required:    true,
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode to assign this disk to.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"size": schema.Int64Attribute{
			Description: "The ID of the token.",
			Required:    true,
		},
		"authorized_keys": schema.SetAttribute{
			Description: "A list of public SSH keys that will be automatically " +
				"appended to the root user's ~/.ssh/authorized_keys file when " +
				"deploying from an Image.",
			Optional:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.RequiresReplace(),
			},
			Validators: []validator.Set{
				setvalidator.AlsoRequires(path.MatchRoot("image")),
			},
		},
		"authorized_users": schema.SetAttribute{
			Description: "A list of usernames. If the usernames have associated SSH " +
				"keys, the keys will be appended to the root users ~/.ssh/authorized_keys " +
				"file automatically when deploying from an Image.",
			Optional:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.RequiresReplace(),
			},
			Validators: []validator.Set{
				setvalidator.AlsoRequires(path.MatchRoot("image")),
			},
		},
		"filesystem": schema.StringAttribute{
			Description: "The filesystem of this disk.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("raw", "swap", "ext3", "ext4", "initrd"),
			},
		},
		"image": schema.StringAttribute{
			Description: "An Image ID to deploy the Linode Disk from.",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"root_pass": schema.StringAttribute{
			Description: "This sets the root user's password on a " +
				"newly-created Linode Disk when deploying from an Image.",
			Optional:  true,
			Sensitive: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.LengthBetween(
					helper.RootPassMinimumCharacters,
					helper.RootPassMaximumCharacters,
				),
			},
		},
		"stackscript_data": schema.MapAttribute{
			Description: "An object containing responses to any User Defined " +
				"Fields present in the StackScript being deployed to this Disk. " +
				"Only accepted if 'stackscript_id' is given. The required values " +
				"depend on the StackScript being deployed.",
			ElementType: types.StringType,
			Optional:    true,
			Sensitive:   true,
			PlanModifiers: []planmodifier.Map{
				mapplanmodifier.RequiresReplace(),
			},
			Validators: []validator.Map{
				mapvalidator.AlsoRequires(path.MatchRoot("image")),
			},
		},
		"stackscript_id": schema.Int64Attribute{
			Description: "A StackScript ID that will cause the referenced " +
				"StackScript to be run during deployment of this Linode.",
			Optional: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
			Validators: []validator.Int64{
				int64validator.AlsoRequires(path.MatchRoot("image")),
			},
		},
		"created": schema.StringAttribute{
			Description: "When this disk was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "When this disk was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"status": schema.StringAttribute{
			Description: "A brief description of this Disk's current state.",
			Computed:    true,
		},
		"disk_encryption": schema.StringAttribute{
			Description: "The disk encryption policy for this disk's parent Linode. " +
				"NOTE: Disk encryption may not currently be available to all users.",
			Computed: true,
		},
	},
}
