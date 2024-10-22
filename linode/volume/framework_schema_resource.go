package volume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	linodeplanmodifiers "github.com/linode/terraform-provider-linode/v2/linode/helper/planmodifiers"
)

const RequireReplacementWhenNewSourceVolumeIDIsNotNull = "When source_volume_id is set to a non-null new value, a replacement will be required."

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the volume.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"source_volume_id": schema.Int64Attribute{
			Description: "The ID of a volume to clone.",
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplaceIf(
					func(
						ctx context.Context,
						sr planmodifier.Int64Request,
						rrifr *int64planmodifier.RequiresReplaceIfFuncResponse,
					) {
						rrifr.RequiresReplace = !sr.PlanValue.IsNull()
					},
					RequireReplacementWhenNewSourceVolumeIDIsNotNull,
					RequireReplacementWhenNewSourceVolumeIDIsNotNull,
				),
			},
			Validators: []validator.Int64{
				int64validator.AtLeastOneOf(
					path.MatchRoot("region"),
					path.MatchRoot("source_volume_id"),
				),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label of the Linode Volume.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 32),
			},
		},
		"status": schema.StringAttribute{
			Description: "The status of the volume, indicating the current readiness state.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Description: "The region where this volume will be deployed.",
			Optional:    true,
			Computed:    true,
			Validators: []validator.String{
				stringvalidator.AtLeastOneOf(
					path.MatchRoot("region"),
					path.MatchRoot("source_volume_id"),
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"size": schema.Int64Attribute{
			Description: "Size of the Volume in GB",
			Optional:    true,
			Computed:    true,
			Default:     int64default.StaticInt64(20),
		},
		"linode_id": schema.Int64Attribute{
			Description: "The Linode ID where the Volume should be attached.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"filesystem_path": schema.StringAttribute{
			Description: "The full filesystem path for the Volume based on the Volume's label. Path is " +
				"/dev/disk/by-id/scsi-0Linode_Volume_ + Volume label.",
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"tags": schema.SetAttribute{
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Set{
				linodeplanmodifiers.CaseInsensitiveSet(),
			},
			Default: helper.EmptySetDefault(types.StringType),
		},
		"encryption": schema.StringAttribute{
			Description: "Whether Block Storage Disk Encryption is enabled or disabled on this Volume. " +
				"Note: Block Storage Disk Encryption is not currently available to all users.",
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("disabled"),
			Validators: []validator.String{
				stringvalidator.OneOf("enabled", "disabled"),
			},
		},
	},
}
