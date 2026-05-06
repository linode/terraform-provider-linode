package monitoralertchannels

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":           {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"type":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"channel_type": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var monitorAlertChannelAttributes = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The unique identifier for the notification channel.",
		Computed:    true,
	},
	"label": schema.StringAttribute{
		Description: "The name of the notification channel for identification purposes.",
		Computed:    true,
	},
	"type": schema.StringAttribute{
		Description: "The type of notification channel. Valid values are system and user.",
		Computed:    true,
	},
	"channel_type": schema.StringAttribute{
		Description: "The delivery mechanism used by the channel. Currently, only email is supported.",
		Computed:    true,
	},
	"created": schema.StringAttribute{
		Description: "When the notification channel was created.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Description: "When the notification channel was last updated. If never updated, this equals created.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"created_by": schema.StringAttribute{
		Description: "For user channels, the account user who created it; for system channels, this is system.",
		Computed:    true,
	},
	"updated_by": schema.StringAttribute{
		Description: "For user channels, the account user who last updated it; for system channels, this is system.",
		Computed:    true,
	},
	"alerts": schema.SingleNestedAttribute{
		Description: "Details about the alerts where this notification channel is applied.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "The API URL for the associated alerts.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The alert type associated with this channel.",
				Computed:    true,
			},
			"alert_count": schema.Int64Attribute{
				Description: "The number of alerts associated with this channel.",
				Computed:    true,
			},
		},
	},
	"details": schema.SingleNestedAttribute{
		Description: "The notification channel configuration details.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"email": schema.SingleNestedAttribute{
				Description: "Email delivery configuration for the notification channel.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"usernames": schema.ListAttribute{
						Description: "Usernames on the account that receive the alert for user channels.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"recipient_type": schema.StringAttribute{
						Description: "Recipient selection strategy. For system channels this is read_write_users; for user channels this is user.",
						Computed:    true,
					},
				},
			},
		},
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"monitor_alert_channels": schema.ListNestedAttribute{
			Description: "The returned list of alert notification channels.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: monitorAlertChannelAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
