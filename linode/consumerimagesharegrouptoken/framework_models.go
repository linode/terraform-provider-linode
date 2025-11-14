package consumerimagesharegrouptoken

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type ResourceModel struct {
	ValidForShareGroupUUID types.String      `tfsdk:"valid_for_sharegroup_uuid"`
	Label                  types.String      `tfsdk:"label"`
	Token                  types.String      `tfsdk:"token"`
	TokenUUID              types.String      `tfsdk:"token_uuid"`
	Status                 types.String      `tfsdk:"status"`
	Created                timetypes.RFC3339 `tfsdk:"created"`
	Updated                timetypes.RFC3339 `tfsdk:"updated"`
	Expiry                 timetypes.RFC3339 `tfsdk:"expiry"`
	ShareGroupUUID         types.String      `tfsdk:"sharegroup_uuid"`
	ShareGroupLabel        types.String      `tfsdk:"sharegroup_label"`
}

func (data *ResourceModel) FlattenImageShareGroupCreateToken(
	resp *linodego.ImageShareGroupCreateTokenResponse,
) {
	data.ValidForShareGroupUUID = types.StringValue(resp.ValidForShareGroupUUID)

	if resp.Label != "" {
		data.Label = types.StringValue(resp.Label)
	} else {
		data.Label = types.StringNull()
	}

	data.TokenUUID = types.StringValue(resp.TokenUUID)
	data.Status = types.StringValue(resp.Status)
	data.Created = timetypes.NewRFC3339TimePointerValue(resp.Created)
	data.Updated = timetypes.NewRFC3339TimePointerValue(resp.Updated)
	data.Expiry = timetypes.NewRFC3339TimePointerValue(resp.Expiry)
	data.ShareGroupUUID = types.StringPointerValue(resp.ShareGroupUUID)
	data.ShareGroupLabel = types.StringPointerValue(resp.ShareGroupLabel)

	// Token is only present in the API response during creation
	data.Token = types.StringValue(resp.Token)
}

func (data *ResourceModel) FlattenImageShareGroupToken(
	token *linodego.ImageShareGroupToken,
	preserveKnown bool,
) {
	// Do not touch Token here since it’s only returned at create time

	data.ValidForShareGroupUUID = helper.KeepOrUpdateString(data.ValidForShareGroupUUID, token.ValidForShareGroupUUID, preserveKnown)

	if token.Label != "" {
		data.Label = helper.KeepOrUpdateString(data.Label, token.Label, preserveKnown)
	} else if !preserveKnown {
		data.Label = types.StringNull()
	}

	data.TokenUUID = helper.KeepOrUpdateString(data.TokenUUID, token.TokenUUID, preserveKnown)
	data.Status = helper.KeepOrUpdateString(data.Status, token.Status, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(
		data.Created, timetypes.NewRFC3339TimePointerValue(token.Created), preserveKnown,
	)
	data.Updated = helper.KeepOrUpdateValue(
		data.Updated, timetypes.NewRFC3339TimePointerValue(token.Updated), preserveKnown,
	)
	data.Expiry = helper.KeepOrUpdateValue(
		data.Expiry, timetypes.NewRFC3339TimePointerValue(token.Expiry), preserveKnown,
	)
	data.ShareGroupUUID = helper.KeepOrUpdateStringPointer(data.ShareGroupUUID, token.ShareGroupUUID, preserveKnown)
	data.ShareGroupLabel = helper.KeepOrUpdateStringPointer(data.ShareGroupLabel, token.ShareGroupLabel, preserveKnown)
}

func (m *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	m.ValidForShareGroupUUID = helper.KeepOrUpdateValue(m.ValidForShareGroupUUID, other.ValidForShareGroupUUID, preserveKnown)
	m.Token = helper.KeepOrUpdateValue(m.Token, other.Token, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.TokenUUID = helper.KeepOrUpdateValue(m.TokenUUID, other.TokenUUID, preserveKnown)
	m.Status = helper.KeepOrUpdateValue(m.Status, other.Status, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Expiry = helper.KeepOrUpdateValue(m.Expiry, other.Expiry, preserveKnown)
	m.ShareGroupUUID = helper.KeepOrUpdateValue(m.ShareGroupUUID, other.ShareGroupUUID, preserveKnown)
	m.ShareGroupLabel = helper.KeepOrUpdateValue(m.ShareGroupLabel, other.ShareGroupLabel, preserveKnown)
}

type DataSourceModel struct {
	TokenUUID              types.String      `tfsdk:"token_uuid"`
	Label                  types.String      `tfsdk:"label"`
	Status                 types.String      `tfsdk:"status"`
	Created                timetypes.RFC3339 `tfsdk:"created"`
	Updated                timetypes.RFC3339 `tfsdk:"updated"`
	Expiry                 timetypes.RFC3339 `tfsdk:"expiry"`
	ValidForShareGroupUUID types.String      `tfsdk:"valid_for_sharegroup_uuid"`
	ShareGroupUUID         types.String      `tfsdk:"sharegroup_uuid"`
	ShareGroupLabel        types.String      `tfsdk:"sharegroup_label"`
}

func (data *DataSourceModel) ParseImageShareGroupToken(m *linodego.ImageShareGroupToken,
) diag.Diagnostics {
	data.TokenUUID = types.StringValue(m.TokenUUID)
	data.ValidForShareGroupUUID = types.StringValue(m.ValidForShareGroupUUID)
	data.Status = types.StringValue(m.Status)
	data.Label = types.StringValue(m.Label)
	data.Created = timetypes.NewRFC3339TimePointerValue(m.Created)
	data.Updated = timetypes.NewRFC3339TimePointerValue(m.Updated)
	data.Expiry = timetypes.NewRFC3339TimePointerValue(m.Expiry)
	data.ShareGroupUUID = types.StringPointerValue(m.ShareGroupUUID)
	data.ShareGroupLabel = types.StringPointerValue(m.ShareGroupLabel)

	return nil
}
