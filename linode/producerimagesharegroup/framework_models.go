package producerimagesharegroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type ResourceModel struct {
	ID           types.Int64       `tfsdk:"id"`
	UUID         types.String      `tfsdk:"uuid"`
	Label        types.String      `tfsdk:"label"`
	Description  types.String      `tfsdk:"description"`
	IsSuspended  types.Bool        `tfsdk:"is_suspended"`
	ImagesCount  types.Int64       `tfsdk:"images_count"`
	MembersCount types.Int64       `tfsdk:"members_count"`
	Created      timetypes.RFC3339 `tfsdk:"created"`
	Updated      timetypes.RFC3339 `tfsdk:"updated"`
	Expiry       timetypes.RFC3339 `tfsdk:"expiry"`
	Images       types.List        `tfsdk:"images"`
}

type ImageShareAttributesModel struct {
	ID          types.String `tfsdk:"id"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
}

func (data *ResourceModel) FlattenImageShareGroup(
	imageShareGroup *linodego.ProducerImageShareGroup,
	preserveKnown bool,
) {
	data.ID = helper.KeepOrUpdateInt64(data.ID, int64(imageShareGroup.ID), preserveKnown)
	data.UUID = helper.KeepOrUpdateString(data.UUID, imageShareGroup.UUID, preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, imageShareGroup.Label, preserveKnown)
	data.Description = helper.KeepOrUpdateString(
		data.Description, imageShareGroup.Description, preserveKnown,
	)
	data.IsSuspended = helper.KeepOrUpdateBool(data.IsSuspended, imageShareGroup.IsSuspended, preserveKnown)
	data.ImagesCount = helper.KeepOrUpdateInt64(data.ImagesCount, int64(imageShareGroup.ImagesCount), preserveKnown)
	data.MembersCount = helper.KeepOrUpdateInt64(data.MembersCount, int64(imageShareGroup.MembersCount), preserveKnown)
	data.Created = helper.KeepOrUpdateValue(
		data.Created, timetypes.NewRFC3339TimePointerValue(imageShareGroup.Created), preserveKnown,
	)
	data.Updated = helper.KeepOrUpdateValue(
		data.Updated, timetypes.NewRFC3339TimePointerValue(imageShareGroup.Updated), preserveKnown,
	)
	data.Expiry = helper.KeepOrUpdateValue(
		data.Expiry, timetypes.NewRFC3339TimePointerValue(imageShareGroup.Expiry), preserveKnown,
	)

	// Images will persist in state across CRUD operations even though it is not returned by the API. It is maintained
	// manually as a part of Create, Update, and Read. We do not need to set it here because it defaults to a
	// properly typed empty list when omitted from the config.
}

type DataSourceModel struct {
	ID           types.Int64       `tfsdk:"id"`
	UUID         types.String      `tfsdk:"uuid"`
	Label        types.String      `tfsdk:"label"`
	Description  types.String      `tfsdk:"description"`
	IsSuspended  types.Bool        `tfsdk:"is_suspended"`
	ImagesCount  types.Int64       `tfsdk:"images_count"`
	MembersCount types.Int64       `tfsdk:"members_count"`
	Created      timetypes.RFC3339 `tfsdk:"created"`
	Updated      timetypes.RFC3339 `tfsdk:"updated"`
	Expiry       timetypes.RFC3339 `tfsdk:"expiry"`
}

func (data *DataSourceModel) ParseImageShareGroup(isg *linodego.ProducerImageShareGroup,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(isg.ID))
	data.UUID = types.StringValue(isg.UUID)
	data.Label = types.StringValue(isg.Label)
	data.Description = types.StringValue(isg.Description)
	data.IsSuspended = types.BoolValue(isg.IsSuspended)
	data.ImagesCount = types.Int64Value(int64(isg.ImagesCount))
	data.MembersCount = types.Int64Value(int64(isg.MembersCount))
	data.Created = timetypes.NewRFC3339TimePointerValue(isg.Created)
	data.Updated = timetypes.NewRFC3339TimePointerValue(isg.Updated)
	data.Expiry = timetypes.NewRFC3339TimePointerValue(isg.Expiry)

	return nil
}
