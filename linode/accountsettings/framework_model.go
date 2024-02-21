package accountsettings

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// AccountSettingsModel describes the Terraform resource data model to match the
// resource schema.
type AccountSettingsModel struct {
	ID                   types.String `tfsdk:"id"`
	LongviewSubscription types.String `tfsdk:"longview_subscription"`
	ObjectStorage        types.String `tfsdk:"object_storage"`
	BackupsEnabled       types.Bool   `tfsdk:"backups_enabled"`
	Managed              types.Bool   `tfsdk:"managed"`
	NetworkHelper        types.Bool   `tfsdk:"network_helper"`
}

func (data *AccountSettingsModel) FlattenAccountSettings(
	email string,
	settings *linodego.AccountSettings,
	preserveKnown bool,
) {
	data.ID = helper.KeepOrUpdateString(data.ID, email, preserveKnown)

	// These use empty strings ("") rather than StringNull to maintain backwards compatibility
	// with the SDKv2 version of this resource.
	data.LongviewSubscription = helper.KeepOrUpdateValue(
		data.LongviewSubscription,
		helper.GetStringPtrWithDefault(settings.LongviewSubscription, ""),
		preserveKnown,
	)
	data.ObjectStorage = helper.KeepOrUpdateValue(
		data.ObjectStorage,
		helper.GetStringPtrWithDefault(settings.ObjectStorage, ""),
		preserveKnown,
	)

	data.Managed = helper.KeepOrUpdateBool(data.Managed, settings.Managed, preserveKnown)
	data.BackupsEnabled = helper.KeepOrUpdateBool(
		data.BackupsEnabled, settings.BackupsEnabled, preserveKnown,
	)
	data.NetworkHelper = helper.KeepOrUpdateBool(
		data.NetworkHelper, settings.NetworkHelper, preserveKnown,
	)
}

func (data *AccountSettingsModel) CopyFrom(other AccountSettingsModel, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.LongviewSubscription = helper.KeepOrUpdateValue(
		data.LongviewSubscription, other.LongviewSubscription, preserveKnown,
	)
	data.ObjectStorage = helper.KeepOrUpdateValue(
		data.ObjectStorage, other.ObjectStorage, preserveKnown,
	)
	data.BackupsEnabled = helper.KeepOrUpdateValue(
		data.BackupsEnabled, other.BackupsEnabled, preserveKnown,
	)
	data.Managed = helper.KeepOrUpdateValue(data.Managed, other.Managed, preserveKnown)
	data.NetworkHelper = helper.KeepOrUpdateValue(
		data.NetworkHelper, other.NetworkHelper, preserveKnown,
	)
}
