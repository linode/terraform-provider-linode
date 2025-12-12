package producerimagesharegroupmember

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type ResourceModel struct {
	ShareGroupID types.Int64       `tfsdk:"sharegroup_id"`
	Token        types.String      `tfsdk:"token"`
	Label        types.String      `tfsdk:"label"`
	TokenUUID    types.String      `tfsdk:"token_uuid"`
	Status       types.String      `tfsdk:"status"`
	Created      timetypes.RFC3339 `tfsdk:"created"`
	Updated      timetypes.RFC3339 `tfsdk:"updated"`
	Expiry       timetypes.RFC3339 `tfsdk:"expiry"`
}

func (data *ResourceModel) FlattenImageShareGroupMember(
	member *linodego.ImageShareGroupMember,
	preserveKnown bool,
) {
	// Do not touch ShareGroupID or Token as they are not returned by the API and must be preserved

	data.Label = helper.KeepOrUpdateString(data.Label, member.Label, preserveKnown)
	data.TokenUUID = helper.KeepOrUpdateString(data.TokenUUID, member.TokenUUID, preserveKnown)
	data.Status = helper.KeepOrUpdateString(data.Status, member.Status, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(
		data.Created, timetypes.NewRFC3339TimePointerValue(member.Created), preserveKnown,
	)
	data.Updated = helper.KeepOrUpdateValue(
		data.Updated, timetypes.NewRFC3339TimePointerValue(member.Updated), preserveKnown,
	)
	data.Expiry = helper.KeepOrUpdateValue(
		data.Expiry, timetypes.NewRFC3339TimePointerValue(member.Expiry), preserveKnown,
	)
}

func (m *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	m.ShareGroupID = helper.KeepOrUpdateValue(m.ShareGroupID, other.ShareGroupID, preserveKnown)
	m.Token = helper.KeepOrUpdateValue(m.Token, other.Token, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.TokenUUID = helper.KeepOrUpdateValue(m.TokenUUID, other.TokenUUID, preserveKnown)
	m.Status = helper.KeepOrUpdateValue(m.Status, other.Status, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Expiry = helper.KeepOrUpdateValue(m.Expiry, other.Expiry, preserveKnown)
}

type DataSourceModel struct {
	ShareGroupID types.Int64       `tfsdk:"sharegroup_id"`
	TokenUUID    types.String      `tfsdk:"token_uuid"`
	Label        types.String      `tfsdk:"label"`
	Status       types.String      `tfsdk:"status"`
	Created      timetypes.RFC3339 `tfsdk:"created"`
	Updated      timetypes.RFC3339 `tfsdk:"updated"`
	Expiry       timetypes.RFC3339 `tfsdk:"expiry"`
}

func (data *DataSourceModel) ParseImageShareGroupMember(m *linodego.ImageShareGroupMember,
) diag.Diagnostics {
	// Do not touch ShareGroupID as it is not returned by the API and must be preserved

	data.TokenUUID = types.StringValue(m.TokenUUID)
	data.Status = types.StringValue(m.Status)
	data.Label = types.StringValue(m.Label)
	data.Created = timetypes.NewRFC3339TimePointerValue(m.Created)
	data.Updated = timetypes.NewRFC3339TimePointerValue(m.Updated)
	data.Expiry = timetypes.NewRFC3339TimePointerValue(m.Expiry)

	return nil
}
