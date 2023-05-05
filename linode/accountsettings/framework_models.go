package accountsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
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
	ctx context.Context,
	email string,
	settings *linodego.AccountSettings,
) diag.Diagnostics {
	data.ID = types.StringValue(email)

	longviewSubscription := ""
	if settings.LongviewSubscription != nil {
		longviewSubscription = *settings.LongviewSubscription
	}

	objectStorage := ""
	if settings.ObjectStorage != nil {
		objectStorage = *settings.ObjectStorage
	}

	data.LongviewSubscription = types.StringValue(longviewSubscription)
	data.ObjectStorage = types.StringValue(objectStorage)
	data.BackupsEnabed = types.BoolValue(settings.BackupsEnabled)
	data.Managed = types.BoolValue(settings.Managed)
	data.NetworkHelper = types.BoolValue(settings.NetworkHelper)

	return nil
}
