package monitoralertchannels

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

// MonitorAlertChannelFilterModel describes the Terraform data source data model
// to match the data source schema.
type MonitorAlertChannelFilterModel struct {
	ID                   types.String                     `tfsdk:"id"`
	Filters              frameworkfilter.FiltersModelType `tfsdk:"filter"`
	MonitorAlertChannels []MonitorAlertChannelModel       `tfsdk:"monitor_alert_channels"`
}

type MonitorAlertChannelModel struct {
	ID          types.Int64                 `tfsdk:"id"`
	Label       types.String                `tfsdk:"label"`
	Type        types.String                `tfsdk:"type"`
	ChannelType types.String                `tfsdk:"channel_type"`
	Created     timetypes.RFC3339           `tfsdk:"created"`
	Updated     timetypes.RFC3339           `tfsdk:"updated"`
	CreatedBy   types.String                `tfsdk:"created_by"`
	UpdatedBy   types.String                `tfsdk:"updated_by"`
	Alerts      *MonitorAlertChannelAlerts  `tfsdk:"alerts"`
	Content     *MonitorAlertChannelContent `tfsdk:"content"`
	Details     *MonitorAlertChannelDetails `tfsdk:"details"`
}

type MonitorAlertChannelAlerts struct {
	URL        types.String `tfsdk:"url"`
	Type       types.String `tfsdk:"type"`
	AlertCount types.Int64  `tfsdk:"alert_count"`
}

type MonitorAlertChannelContent struct {
	Email *MonitorAlertChannelEmailContent `tfsdk:"email"`
}

type MonitorAlertChannelEmailContent struct {
	EmailAddresses types.List `tfsdk:"email_addresses"`
}

type MonitorAlertChannelDetails struct {
	Email *MonitorAlertChannelEmailDetails `tfsdk:"email"`
}

type MonitorAlertChannelEmailDetails struct {
	Usernames     types.List   `tfsdk:"usernames"`
	RecipientType types.String `tfsdk:"recipient_type"`
}

func (data *MonitorAlertChannelFilterModel) parseMonitorAlertChannels(channels []linodego.AlertChannel) {
	result := make([]MonitorAlertChannelModel, len(channels))
	for i := range channels {
		result[i] = flattenMonitorAlertChannel(context.Background(), &channels[i])
	}

	data.MonitorAlertChannels = result
}

func flattenMonitorAlertChannel(ctx context.Context, channel *linodego.AlertChannel) MonitorAlertChannelModel {
	model := MonitorAlertChannelModel{
		ID:          types.Int64Value(int64(channel.ID)),
		Label:       types.StringValue(channel.Label),
		Type:        types.StringValue(string(channel.Type)),
		ChannelType: types.StringValue(string(channel.ChannelType)),
		Created:     timetypes.NewRFC3339TimePointerValue(channel.Created),
		Updated:     timetypes.NewRFC3339TimePointerValue(channel.Updated),
		CreatedBy:   types.StringValue(channel.CreatedBy),
		UpdatedBy:   types.StringValue(channel.UpdatedBy),
		Alerts: &MonitorAlertChannelAlerts{
			URL:        types.StringValue(channel.Alerts.URL),
			Type:       types.StringValue(channel.Alerts.Type),
			AlertCount: types.Int64Value(int64(channel.Alerts.AlertCount)),
		},
	}

	if channel.Content.Email != nil {
		emails, _ := types.ListValueFrom(ctx, types.StringType, channel.Content.Email.EmailAddresses)
		model.Content = &MonitorAlertChannelContent{
			Email: &MonitorAlertChannelEmailContent{
				EmailAddresses: emails,
			},
		}
	}

	if channel.Details.Email != nil {
		usernames, _ := types.ListValueFrom(ctx, types.StringType, channel.Details.Email.Usernames)
		model.Details = &MonitorAlertChannelDetails{
			Email: &MonitorAlertChannelEmailDetails{
				Usernames:     usernames,
				RecipientType: types.StringValue(channel.Details.Email.RecipientType),
			},
		}
	}

	return model
}
