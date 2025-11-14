package consumerimagesharegroupimageshares

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

type DataSourceModel struct {
	ID           types.String                 `tfsdk:"id"`
	Label        types.String                 `tfsdk:"label"`
	Description  types.String                 `tfsdk:"description"`
	Capabilities []types.String               `tfsdk:"capabilities"`
	Created      types.String                 `tfsdk:"created"`
	Deprecated   types.Bool                   `tfsdk:"deprecated"`
	IsPublic     types.Bool                   `tfsdk:"is_public"`
	ImageSharing *ImageSharingDataSourceModel `tfsdk:"image_sharing"`
	Size         types.Int64                  `tfsdk:"size"`
	Status       types.String                 `tfsdk:"status"`
	Type         types.String                 `tfsdk:"type"`
	Tags         types.List                   `tfsdk:"tags"`
	TotalSize    types.Int64                  `tfsdk:"total_size"`
}

type ImageSharingDataSourceModel struct {
	SharedWith *ImageSharingSharedWithAttributesModel `tfsdk:"shared_with"`
	SharedBy   *ImageSharingSharedByAttributesModel   `tfsdk:"shared_by"`
}

type ImageSharingSharedWithAttributesModel struct {
	ShareGroupCount   types.Int64  `tfsdk:"sharegroup_count"`
	ShareGroupListURL types.String `tfsdk:"sharegroup_list_url"`
}

type ImageSharingSharedByAttributesModel struct {
	ShareGroupID    types.Int64  `tfsdk:"sharegroup_id"`
	ShareGroupUUID  types.String `tfsdk:"sharegroup_uuid"`
	ShareGroupLabel types.String `tfsdk:"sharegroup_label"`
	SourceImageID   types.String `tfsdk:"source_image_id"`
}

func (data *DataSourceModel) ParseImageShare(
	ctx context.Context,
	imageShare *linodego.ImageShareEntry,
) diag.Diagnostics {
	data.ID = types.StringValue(imageShare.ID)
	data.Label = types.StringValue(imageShare.Label)

	data.Description = types.StringValue(imageShare.Description)
	if imageShare.Created != nil {
		data.Created = types.StringValue(imageShare.Created.Format(time.RFC3339))
	} else {
		data.Created = types.StringNull()
	}
	data.Capabilities = helper.StringSliceToFramework(imageShare.Capabilities)
	data.Deprecated = types.BoolValue(imageShare.Deprecated)
	data.IsPublic = types.BoolValue(imageShare.IsPublic)
	data.Size = types.Int64Value(int64(imageShare.Size))
	data.Status = types.StringValue(string(imageShare.Status))
	data.Type = types.StringValue(imageShare.Type)
	data.TotalSize = types.Int64Value(int64(imageShare.TotalSize))

	tags, diags := types.ListValueFrom(ctx, types.StringType, imageShare.Tags)
	if diags.HasError() {
		return diags
	}
	data.Tags = tags

	data.ImageSharing = parseImageSharingDataSourceModel(&imageShare.ImageSharing)

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

type ImageShareGroupImageShareFilterModel struct {
	ID          types.String                     `tfsdk:"id"`
	TokenUUID   types.String                     `tfsdk:"token_uuid"`
	Filters     frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order       types.String                     `tfsdk:"order"`
	OrderBy     types.String                     `tfsdk:"order_by"`
	ImageShares []DataSourceModel                `tfsdk:"image_shares"`
}

func (data *ImageShareGroupImageShareFilterModel) parseImageShares(
	ctx context.Context,
	imageShares []linodego.ImageShareEntry,
) diag.Diagnostics {
	result := make([]DataSourceModel, len(imageShares))
	for i := range imageShares {
		var imgShareData DataSourceModel
		diags := imgShareData.ParseImageShare(ctx, &imageShares[i])
		if diags.HasError() {
			return diags
		}
		result[i] = imgShareData
	}

	data.ImageShares = result

	return nil
}
