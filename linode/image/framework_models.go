package image

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type ResourceModel struct {
	ID                  types.String      `tfsdk:"id"`
	Label               types.String      `tfsdk:"label"`
	DiskID              types.Int64       `tfsdk:"disk_id"`
	LinodeID            types.Int64       `tfsdk:"linode_id"`
	FilePath            types.String      `tfsdk:"file_path"`
	Region              types.String      `tfsdk:"region"`
	FileHash            types.String      `tfsdk:"file_hash"`
	Description         types.String      `tfsdk:"description"`
	CloudInit           types.Bool        `tfsdk:"cloud_init"`
	Capabilities        types.List        `tfsdk:"capabilities"`
	Created             timetypes.RFC3339 `tfsdk:"created"`
	CreatedBy           types.String      `tfsdk:"created_by"`
	Deprecated          types.Bool        `tfsdk:"deprecated"`
	IsPublic            types.Bool        `tfsdk:"is_public"`
	IsShared            types.Bool        `tfsdk:"is_shared"`
	ImageSharing        types.Object      `tfsdk:"image_sharing"`
	Size                types.Int64       `tfsdk:"size"`
	Status              types.String      `tfsdk:"status"`
	Type                types.String      `tfsdk:"type"`
	Expiry              timetypes.RFC3339 `tfsdk:"expiry"`
	Vendor              types.String      `tfsdk:"vendor"`
	Timeouts            timeouts.Value    `tfsdk:"timeouts"`
	Tags                types.List        `tfsdk:"tags"`
	TotalSize           types.Int64       `tfsdk:"total_size"`
	ReplicaRegions      types.List        `tfsdk:"replica_regions"`
	Replications        types.List        `tfsdk:"replications"`
	WaitForReplications types.Bool        `tfsdk:"wait_for_replications"`
}

type ImageSharingModel struct {
	SharedWith types.Object `tfsdk:"shared_with"`
	SharedBy   types.Object `tfsdk:"shared_by"`
}

var imageSharingObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"shared_with": imageSharingSharedWithObjectType,
		"shared_by":   imageSharingSharedByObjectType,
	},
}

type ImageSharingSharedWithAttributesModel struct {
	ShareGroupCount   types.Int64  `tfsdk:"sharegroup_count"`
	ShareGroupListURL types.String `tfsdk:"sharegroup_list_url"`
}

var imageSharingSharedWithObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"sharegroup_count":    types.Int64Type,
		"sharegroup_list_url": types.StringType,
	},
}

type ImageSharingSharedByAttributesModel struct {
	ShareGroupID    types.Int64  `tfsdk:"sharegroup_id"`
	ShareGroupUUID  types.String `tfsdk:"sharegroup_uuid"`
	ShareGroupLabel types.String `tfsdk:"sharegroup_label"`
	SourceImageID   types.String `tfsdk:"source_image_id"`
}

var imageSharingSharedByObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"sharegroup_id":    types.Int64Type,
		"sharegroup_uuid":  types.StringType,
		"sharegroup_label": types.StringType,
		"source_image_id":  types.StringType,
	},
}

func (data *ResourceModel) FlattenImage(
	ctx context.Context,
	image *linodego.Image,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
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
	data.IsShared = helper.KeepOrUpdateBool(data.IsShared, image.IsShared, preserveKnown)

	imageSharing := helper.KeepOrUpdateSingleNestedAttributeWithTypes(
		ctx,
		data.ImageSharing,
		imageSharingObjectType.AttrTypes,
		preserveKnown,
		diags,
		func(model *ImageSharingModel, isNull *bool, preserveKnown bool, diags *diag.Diagnostics) {
			model.FlattenImageSharing(ctx, image.ImageSharing, preserveKnown, diags)
		},
	)

	if imageSharing != nil {
		data.ImageSharing = *imageSharing
	}

	data.Size = helper.KeepOrUpdateInt64(data.Size, int64(image.Size), preserveKnown)
	data.Status = helper.KeepOrUpdateString(data.Status, string(image.Status), preserveKnown)
	data.Type = helper.KeepOrUpdateString(data.Type, image.Type, preserveKnown)
	data.Expiry = helper.KeepOrUpdateValue(
		data.Expiry, timetypes.NewRFC3339TimePointerValue(image.Expiry), preserveKnown,
	)
	data.Vendor = helper.KeepOrUpdateString(data.Vendor, image.Vendor, preserveKnown)
	data.TotalSize = helper.KeepOrUpdateInt64(data.TotalSize, int64(image.TotalSize), preserveKnown)

	tags, newDiags := types.ListValue(types.StringType, helper.StringSliceToFrameworkValueSlice(image.Tags))
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	data.Tags = helper.KeepOrUpdateValue(data.Tags, tags, preserveKnown)

	replications, newDiags := flattenReplications(ctx, image.Regions)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	data.Replications = helper.KeepOrUpdateValue(data.Replications, *replications, preserveKnown)
}

func (m *ImageSharingModel) FlattenImageSharing(ctx context.Context, imageSharing linodego.ImageSharing, preserveKnown bool, diags *diag.Diagnostics) {
	if m.SharedWith.IsNull() {
		m.SharedWith = types.ObjectNull(imageSharingSharedWithObjectType.AttrTypes)
	}
	if m.SharedBy.IsNull() {
		m.SharedBy = types.ObjectNull(imageSharingSharedByObjectType.AttrTypes)
	}

	if imageSharing.SharedWith != nil {
		var sharedWithModel ImageSharingSharedWithAttributesModel
		sharedWithModel.FlattenImageSharingSharedWith(*imageSharing.SharedWith, preserveKnown)
		objVal, d := types.ObjectValueFrom(ctx, imageSharingSharedWithObjectType.AttrTypes, sharedWithModel)
		diags.Append(d...)
		m.SharedWith = objVal
	} else if !preserveKnown {
		m.SharedWith = types.ObjectNull(imageSharingSharedWithObjectType.AttrTypes)
	}

	if imageSharing.SharedBy != nil {
		var sharedByModel ImageSharingSharedByAttributesModel
		sharedByModel.FlattenImageSharingSharedBy(*imageSharing.SharedBy, preserveKnown)
		objVal, d := types.ObjectValueFrom(ctx, imageSharingSharedByObjectType.AttrTypes, sharedByModel)
		diags.Append(d...)
		m.SharedBy = objVal
	} else if !preserveKnown {
		m.SharedBy = types.ObjectNull(imageSharingSharedByObjectType.AttrTypes)
	}
}

func (m *ImageSharingSharedWithAttributesModel) FlattenImageSharingSharedWith(sharedWith linodego.ImageSharingSharedWith, preserveKnown bool) {
	m.ShareGroupCount = helper.KeepOrUpdateInt64(m.ShareGroupCount, int64(sharedWith.ShareGroupCount), preserveKnown)
	m.ShareGroupListURL = helper.KeepOrUpdateString(m.ShareGroupListURL, sharedWith.ShareGroupListURL, preserveKnown)
}

func (m *ImageSharingSharedByAttributesModel) FlattenImageSharingSharedBy(sharedBy linodego.ImageSharingSharedBy, preserveKnown bool) {
	m.ShareGroupID = helper.KeepOrUpdateInt64(m.ShareGroupID, int64(sharedBy.ShareGroupID), preserveKnown)
	m.ShareGroupUUID = helper.KeepOrUpdateString(m.ShareGroupUUID, sharedBy.ShareGroupUUID, preserveKnown)
	m.ShareGroupLabel = helper.KeepOrUpdateString(m.ShareGroupLabel, sharedBy.ShareGroupLabel, preserveKnown)
	m.SourceImageID = helper.KeepOrUpdateStringPointer(m.SourceImageID, sharedBy.SourceImageID, preserveKnown)
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
	data.IsShared = helper.KeepOrUpdateValue(data.IsShared, other.IsShared, preserveKnown)
	data.ImageSharing = helper.KeepOrUpdateValue(data.ImageSharing, other.ImageSharing, preserveKnown)
	data.Size = helper.KeepOrUpdateValue(data.Size, other.Size, preserveKnown)
	data.Status = helper.KeepOrUpdateValue(data.Status, other.Status, preserveKnown)
	data.Type = helper.KeepOrUpdateValue(data.Type, other.Type, preserveKnown)
	data.Expiry = helper.KeepOrUpdateValue(data.Expiry, other.Expiry, preserveKnown)
	data.Vendor = helper.KeepOrUpdateValue(data.Vendor, other.Vendor, preserveKnown)
	data.Timeouts = helper.KeepOrUpdateValue(data.Timeouts, other.Timeouts, preserveKnown)
	data.Tags = helper.KeepOrUpdateValue(data.Tags, other.Tags, preserveKnown)
	data.TotalSize = helper.KeepOrUpdateValue(data.TotalSize, other.TotalSize, preserveKnown)
	data.ReplicaRegions = helper.KeepOrUpdateValue(data.ReplicaRegions, other.ReplicaRegions, preserveKnown)
	data.Replications = helper.KeepOrUpdateValue(data.Replications, other.Replications, preserveKnown)
	data.WaitForReplications = helper.KeepOrUpdateValue(data.WaitForReplications, other.WaitForReplications, preserveKnown)
}

// ImageModel describes the Terraform resource data model to match the
// resource schema.
type ImageModel struct {
	ID           types.String                 `tfsdk:"id"`
	Label        types.String                 `tfsdk:"label"`
	Description  types.String                 `tfsdk:"description"`
	Capabilities []types.String               `tfsdk:"capabilities"`
	Created      types.String                 `tfsdk:"created"`
	CreatedBy    types.String                 `tfsdk:"created_by"`
	Deprecated   types.Bool                   `tfsdk:"deprecated"`
	IsPublic     types.Bool                   `tfsdk:"is_public"`
	IsShared     types.Bool                   `tfsdk:"is_shared"`
	ImageSharing *ImageSharingDataSourceModel `tfsdk:"image_sharing"`
	Size         types.Int64                  `tfsdk:"size"`
	Status       types.String                 `tfsdk:"status"`
	Type         types.String                 `tfsdk:"type"`
	Expiry       types.String                 `tfsdk:"expiry"`
	Vendor       types.String                 `tfsdk:"vendor"`
	Tags         types.List                   `tfsdk:"tags"`
	TotalSize    types.Int64                  `tfsdk:"total_size"`
	Replications []ReplicationModel           `tfsdk:"replications"`
}

// ReplicationModel describes an image replication.
type ReplicationModel struct {
	Region types.String `tfsdk:"region"`
	Status types.String `tfsdk:"status"`
}

type ImageSharingDataSourceModel struct {
	SharedWith *ImageSharingSharedWithAttributesModel `tfsdk:"shared_with"`
	SharedBy   *ImageSharingSharedByAttributesModel   `tfsdk:"shared_by"`
}

func (data *ImageModel) ParseImage(
	ctx context.Context,
	image *linodego.Image,
) diag.Diagnostics {
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
	data.IsShared = types.BoolValue(image.IsShared)
	data.Size = types.Int64Value(int64(image.Size))
	data.Status = types.StringValue(string(image.Status))
	data.Type = types.StringValue(image.Type)
	data.Vendor = types.StringValue(image.Vendor)
	data.TotalSize = types.Int64Value(int64(image.TotalSize))

	tags, diags := types.ListValueFrom(ctx, types.StringType, image.Tags)
	if diags.HasError() {
		return diags
	}
	data.Tags = tags

	data.Replications = parseReplicationModels(image.Regions)
	data.ImageSharing = parseImageSharingDataSourceModel(&image.ImageSharing)

	return nil
}

func parseImageSharingDataSourceModel(
	imageSharing *linodego.ImageSharing,
) *ImageSharingDataSourceModel {
	if imageSharing == nil {
		return nil
	}

	var sharedWith *ImageSharingSharedWithAttributesModel
	if sw := imageSharing.SharedWith; sw != nil {
		sharedWith = &ImageSharingSharedWithAttributesModel{
			ShareGroupCount:   types.Int64Value(int64(sw.ShareGroupCount)),
			ShareGroupListURL: types.StringValue(sw.ShareGroupListURL),
		}
	}

	var sharedBy *ImageSharingSharedByAttributesModel
	if sb := imageSharing.SharedBy; sb != nil {
		sharedBy = &ImageSharingSharedByAttributesModel{
			ShareGroupID:    types.Int64Value(int64(sb.ShareGroupID)),
			ShareGroupUUID:  types.StringValue(sb.ShareGroupUUID),
			ShareGroupLabel: types.StringValue(sb.ShareGroupLabel),
			SourceImageID:   types.StringPointerValue(sb.SourceImageID),
		}
	}

	return &ImageSharingDataSourceModel{
		SharedWith: sharedWith,
		SharedBy:   sharedBy,
	}
}

func parseReplicationModels(
	regions []linodego.ImageRegion,
) []ReplicationModel {
	replications := make([]ReplicationModel, len(regions))

	for i, r := range regions {
		replications[i].Region = types.StringValue(r.Region)
		replications[i].Status = types.StringValue(string(r.Status))
	}

	return replications
}

func flattenReplications(
	ctx context.Context,
	regions []linodego.ImageRegion,
) (*types.List, diag.Diagnostics) {
	replications := parseReplicationModels(regions)

	result, diags := types.ListValueFrom(ctx, replicationObjType, replications)
	if diags.HasError() {
		return nil, diags
	}

	return &result, nil
}
