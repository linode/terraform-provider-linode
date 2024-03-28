package image

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type ResourceModel struct {
	ID           types.String      `tfsdk:"id"`
	Label        types.String      `tfsdk:"label"`
	DiskID       types.Int64       `tfsdk:"disk_id"`
	LinodeID     types.Int64       `tfsdk:"linode_id"`
	FilePath     types.String      `tfsdk:"file_path"`
	Region       types.String      `tfsdk:"region"`
	FileHash     types.String      `tfsdk:"file_hash"`
	Description  types.String      `tfsdk:"description"`
	CloudInit    types.Bool        `tfsdk:"cloud_init"`
	Capabilities types.List        `tfsdk:"capabilities"`
	Created      timetypes.RFC3339 `tfsdk:"created"`
	CreatedBy    types.String      `tfsdk:"created_by"`
	Deprecated   types.Bool        `tfsdk:"deprecated"`
	IsPublic     types.Bool        `tfsdk:"is_public"`
	Size         types.Int64       `tfsdk:"size"`
	Status       types.String      `tfsdk:"status"`
	Type         types.String      `tfsdk:"type"`
	Expiry       timetypes.RFC3339 `tfsdk:"expiry"`
	Vendor       types.String      `tfsdk:"vendor"`
	Timeouts     timeouts.Value    `tfsdk:"timeouts"`
}

func (data *ResourceModel) FlattenImage(image *linodego.Image, preserveKnown bool, diags *diag.Diagnostics) {
	data.ID = helper.KeepOrUpdateString(data.ID, image.ID, preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, image.Label, preserveKnown)
	data.Description = helper.KeepOrUpdateString(
		data.Description, image.Description, preserveKnown,
	)

	newCapabilities, newDiags := types.ListValue(
		types.StringType, helper.StringSliceToFrameworkValueSlice(image.Capabilities),
	)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	data.Capabilities = helper.KeepOrUpdateValue(data.Capabilities, newCapabilities, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(
		data.Created, timetypes.NewRFC3339TimePointerValue(image.Created), preserveKnown,
	)

	data.CreatedBy = helper.KeepOrUpdateString(data.CreatedBy, image.CreatedBy, preserveKnown)
	data.Deprecated = helper.KeepOrUpdateBool(data.Deprecated, image.Deprecated, preserveKnown)
	data.IsPublic = helper.KeepOrUpdateBool(data.IsPublic, image.IsPublic, preserveKnown)
	data.Size = helper.KeepOrUpdateInt64(data.Size, int64(image.Size), preserveKnown)
	data.Status = helper.KeepOrUpdateString(data.Status, string(image.Status), preserveKnown)
	data.Type = helper.KeepOrUpdateString(data.Type, image.Type, preserveKnown)
	data.Expiry = helper.KeepOrUpdateValue(
		data.Expiry, timetypes.NewRFC3339TimePointerValue(image.Expiry), preserveKnown,
	)
	data.Vendor = helper.KeepOrUpdateString(data.Vendor, image.Vendor, preserveKnown)
}

func (data *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.Label = helper.KeepOrUpdateValue(data.Label, other.Label, preserveKnown)
	data.DiskID = helper.KeepOrUpdateValue(data.DiskID, other.DiskID, preserveKnown)
	data.LinodeID = helper.KeepOrUpdateValue(data.LinodeID, other.LinodeID, preserveKnown)
	data.FilePath = helper.KeepOrUpdateValue(data.FilePath, other.FilePath, preserveKnown)
	data.Region = helper.KeepOrUpdateValue(data.Region, other.Region, preserveKnown)
	data.FileHash = helper.KeepOrUpdateValue(data.FileHash, other.FileHash, preserveKnown)
	data.Description = helper.KeepOrUpdateValue(data.Description, other.Description, preserveKnown)
	data.CloudInit = helper.KeepOrUpdateValue(data.CloudInit, other.CloudInit, preserveKnown)
	data.Capabilities = helper.KeepOrUpdateValue(data.Capabilities, other.Capabilities, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(data.Created, other.Created, preserveKnown)
	data.CreatedBy = helper.KeepOrUpdateValue(data.CreatedBy, other.CreatedBy, preserveKnown)
	data.Deprecated = helper.KeepOrUpdateValue(data.Deprecated, other.Deprecated, preserveKnown)
	data.IsPublic = helper.KeepOrUpdateValue(data.IsPublic, other.IsPublic, preserveKnown)
	data.Size = helper.KeepOrUpdateValue(data.Size, other.Size, preserveKnown)
	data.Status = helper.KeepOrUpdateValue(data.Status, other.Status, preserveKnown)
	data.Type = helper.KeepOrUpdateValue(data.Type, other.Type, preserveKnown)
	data.Expiry = helper.KeepOrUpdateValue(data.Expiry, other.Expiry, preserveKnown)
	data.Vendor = helper.KeepOrUpdateValue(data.Vendor, other.Vendor, preserveKnown)
	data.Timeouts = helper.KeepOrUpdateValue(data.Timeouts, other.Timeouts, preserveKnown)
}

// ImageModel describes the Terraform resource data model to match the
// resource schema.
type ImageModel struct {
	ID           types.String   `tfsdk:"id"`
	Label        types.String   `tfsdk:"label"`
	Description  types.String   `tfsdk:"description"`
	Capabilities []types.String `tfsdk:"capabilities"`
	Created      types.String   `tfsdk:"created"`
	CreatedBy    types.String   `tfsdk:"created_by"`
	Deprecated   types.Bool     `tfsdk:"deprecated"`
	IsPublic     types.Bool     `tfsdk:"is_public"`
	Size         types.Int64    `tfsdk:"size"`
	Status       types.String   `tfsdk:"status"`
	Type         types.String   `tfsdk:"type"`
	Expiry       types.String   `tfsdk:"expiry"`
	Vendor       types.String   `tfsdk:"vendor"`
}

func (data *ImageModel) ParseImage(
	image *linodego.Image,
) {
	data.ID = types.StringValue(image.ID)
	data.Label = types.StringValue(image.Label)

	data.Description = types.StringValue(image.Description)
	if image.Created != nil {
		data.Created = types.StringValue(image.Created.Format(time.RFC3339))
	} else {
		data.Created = types.StringNull()
	}
	if image.Expiry != nil {
		data.Expiry = types.StringValue(image.Expiry.Format(time.RFC3339))
	} else {
		data.Expiry = types.StringNull()
	}
	data.Capabilities = helper.StringSliceToFramework(image.Capabilities)
	data.CreatedBy = types.StringValue(image.CreatedBy)
	data.Deprecated = types.BoolValue(image.Deprecated)
	data.IsPublic = types.BoolValue(image.IsPublic)
	data.Size = types.Int64Value(int64(image.Size))
	data.Status = types.StringValue(string(image.Status))
	data.Type = types.StringValue(image.Type)
	data.Vendor = types.StringValue(image.Vendor)
}
