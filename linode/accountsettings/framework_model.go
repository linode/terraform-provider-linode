package accountsettings

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

// AccountSettingsModel describes the Terraform resource data model to match the
// resource schema.
type AccountSettingsModel struct {
	ID                   types.String `tfsdk:"id"`
	LongviewSubscription types.String `tfsdk:"longview_subscription"`
	ObjectStorage        types.String `tfsdk:"object_storage"`
	BackupsEnabed        types.Bool   `tfsdk:"backups_enabled"`
	Managed              types.Bool   `tfsdk:"managed"`
	NetworkHelper        types.Bool   `tfsdk:"network_helper"`
}

func (data *AccountSettingsModel) parseAccountSettings(
	email string,
	settings *linodego.AccountSettings,
) {
	data.ID = types.StringValue(email)

	// These use empty strings ("") rather than StringNull to maintain backwards compatibility
	// with the SDKv2 version of this resource.
	data.LongviewSubscription = helper.GetStringPtrWithDefault(settings.LongviewSubscription, "")
	data.ObjectStorage = helper.GetStringPtrWithDefault(settings.ObjectStorage, "")

	data.Managed = types.BoolValue(settings.Managed)
	data.BackupsEnabed = types.BoolValue(settings.BackupsEnabled)
	data.NetworkHelper = types.BoolValue(settings.NetworkHelper)
}
