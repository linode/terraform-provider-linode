package volume

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type VolumeModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Label          types.String `tfsdk:"label"`
	Region         types.String `tfsdk:"region"`
	Size           types.Int64  `tfsdk:"size"`
	LinodeID       types.Int64  `tfsdk:"linode_id"`
	FilesystemPath types.String `tfsdk:"filesystem_path"`
	Tags           types.Set    `tfsdk:"tags"`
	Status         types.String `tfsdk:"status"`
	Created        types.String `tfsdk:"created"`
	Updated        types.String `tfsdk:"updated"`
}

func (data *VolumeModel) ParseComputedAttributes(
	ctx context.Context,
	volume *linodego.Volume,
) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.Int64Value(int64(volume.ID))
	data.Status = types.StringValue(string(volume.Status))
	data.Region = types.StringValue(volume.Region)
	data.Size = types.Int64Value(int64(volume.Size))

	// Future breaking change:
	// if volume.LinodeID != nil {
	// 	data.LinodeID = types.Int64Value(int64(*volume.LinodeID))
	// } else {
	// 	data.LinodeID = types.Int64Null()
	// }
	data.LinodeID = helper.IntPointerValueWithDefault(volume.LinodeID)

	data.FilesystemPath = types.StringValue(volume.FilesystemPath)
	data.Created = types.StringValue(volume.Created.Format(time.RFC3339))
	data.Updated = types.StringValue(volume.Updated.Format(time.RFC3339))

	return diags
}

func (data *VolumeModel) ParseNonComputedAttributes(
	ctx context.Context,
	volume *linodego.Volume,
) diag.Diagnostics {
	var diags diag.Diagnostics

	data.Tags, diags = types.SetValueFrom(ctx, types.StringType, volume.Tags)
	diags.Append(diags...)
	if diags.HasError() {
		return diags
	}

	data.Label = types.StringValue(volume.Label)

	return diags
}
