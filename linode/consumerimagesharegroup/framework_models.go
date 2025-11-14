package consumerimagesharegroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	TokenUUID   types.String      `tfsdk:"token_uuid"`
	ID          types.Int64       `tfsdk:"id"`
	UUID        types.String      `tfsdk:"uuid"`
	Label       types.String      `tfsdk:"label"`
	Description types.String      `tfsdk:"description"`
	IsSuspended types.Bool        `tfsdk:"is_suspended"`
	Created     timetypes.RFC3339 `tfsdk:"created"`
	Updated     timetypes.RFC3339 `tfsdk:"updated"`
}

func (data *DataSourceModel) ParseConsumerImageShareGroup(m *linodego.ConsumerImageShareGroup,
) diag.Diagnostics {
	// Do not touch TokenUUID as it is not returned by the API and must be preserved

	data.ID = types.Int64Value(int64(m.ID))
	data.UUID = types.StringValue(m.UUID)
	data.Label = types.StringValue(m.Label)
	data.Description = types.StringValue(m.Description)
	data.IsSuspended = types.BoolValue(m.IsSuspended)
	data.Created = timetypes.NewRFC3339TimePointerValue(m.Created)
	data.Updated = timetypes.NewRFC3339TimePointerValue(m.Updated)

	return nil
}
