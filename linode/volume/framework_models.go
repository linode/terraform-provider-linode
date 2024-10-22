package volume

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type VolumeDataSourceModel struct {
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
	Encryption     types.String `tfsdk:"encryption"`
}

func (data *VolumeDataSourceModel) ParseComputedAttributes(
	ctx context.Context,
	volume *linodego.Volume,
) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.Int64Value(int64(volume.ID))
	data.Status = types.StringValue(string(volume.Status))
	data.Region = types.StringValue(volume.Region)
	data.Size = types.Int64Value(int64(volume.Size))

	// Future breaking change:
	// data.LinodeID = types.Int64PointerValue(int64(*volume.LinodeID))
	data.LinodeID = helper.IntPointerValueWithDefault(volume.LinodeID)

	data.FilesystemPath = types.StringValue(volume.FilesystemPath)
	data.Created = types.StringValue(volume.Created.Format(time.RFC3339))
	data.Updated = types.StringValue(volume.Updated.Format(time.RFC3339))
	data.Encryption = types.StringValue(volume.Encryption)

	return diags
}

func (data *VolumeDataSourceModel) ParseNonComputedAttributes(
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

type VolumeResourceModel struct {
	SourceVolumeID types.Int64    `tfsdk:"source_volume_id"`
	ID             types.String   `tfsdk:"id"`
	Label          types.String   `tfsdk:"label"`
	Region         types.String   `tfsdk:"region"`
	Size           types.Int64    `tfsdk:"size"`
	LinodeID       types.Int64    `tfsdk:"linode_id"`
	FilesystemPath types.String   `tfsdk:"filesystem_path"`
	Tags           types.Set      `tfsdk:"tags"`
	Status         types.String   `tfsdk:"status"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
	Encryption     types.String   `tfsdk:"encryption"`
}

func (data *VolumeResourceModel) FlattenVolume(volume *linodego.Volume, preserveKnown bool) diag.Diagnostics {
	var diags diag.Diagnostics
	if volume == nil {
		diags.AddError(
			"Volume is nil",
			"The pointer to linodego.Volume received by FlattenVolume function is nil.",
		)
		return diags
	}
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(volume.ID), preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, volume.Label, preserveKnown)
	data.Region = helper.KeepOrUpdateString(data.Region, volume.Region, preserveKnown)
	data.Size = helper.KeepOrUpdateInt64(data.Size, int64(volume.Size), preserveKnown)
	data.Encryption = helper.KeepOrUpdateString(data.Encryption, volume.Encryption, preserveKnown)

	// planned breaking change:
	// 	data.LinodeID = helper.KeepOrUpdateIntPointer(data.LinodeID, volume.LinodeID, preserveKnown)
	data.LinodeID = helper.KeepOrUpdateValue(
		data.LinodeID, helper.IntPointerValueWithDefault(volume.LinodeID), preserveKnown,
	)

	data.FilesystemPath = helper.KeepOrUpdateString(data.FilesystemPath, volume.FilesystemPath, preserveKnown)

	tagValues := make([]attr.Value, len(volume.Tags))
	for i, tag := range volume.Tags {
		tagValues[i] = types.StringValue(tag)
	}

	tagsSetValue, diags := types.SetValue(types.StringType, tagValues)
	diags.Append(diags...)
	if diags.HasError() {
		return diags
	}

	data.Tags = helper.KeepOrUpdateValue(data.Tags, tagsSetValue, preserveKnown)

	data.Status = helper.KeepOrUpdateString(data.Status, string(volume.Status), preserveKnown)

	return diags
}

func (data *VolumeResourceModel) CopyFrom(other VolumeResourceModel, preserveKnown bool) {
	data.SourceVolumeID = helper.KeepOrUpdateValue(data.SourceVolumeID, other.SourceVolumeID, preserveKnown)
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.Label = helper.KeepOrUpdateValue(data.Label, other.Label, preserveKnown)
	data.Region = helper.KeepOrUpdateValue(data.Region, other.Region, preserveKnown)
	data.Encryption = helper.KeepOrUpdateValue(data.Encryption, other.Encryption, preserveKnown)
	data.Size = helper.KeepOrUpdateValue(data.Size, other.Size, preserveKnown)
	data.LinodeID = helper.KeepOrUpdateValue(data.LinodeID, other.LinodeID, preserveKnown)
	data.FilesystemPath = helper.KeepOrUpdateValue(data.FilesystemPath, other.FilesystemPath, preserveKnown)
	data.Tags = helper.KeepOrUpdateValue(data.Tags, other.Tags, preserveKnown)
	data.Status = helper.KeepOrUpdateValue(data.Status, other.Status, preserveKnown)
	data.Timeouts = helper.KeepOrUpdateValue(data.Timeouts, other.Timeouts, preserveKnown)
}
