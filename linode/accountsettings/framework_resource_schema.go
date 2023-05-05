package accountsettings

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The email of the current account.",
			Computed:    true,
		},
		"backups_enabled": schema.BoolAttribute{
			Description: "Account-wide backups default.",
			Computed:    true,
			Optional:    true,
		},
		"network_helper": schema.BoolAttribute{
			Description: "Enables network helper across all users by default for new Linodes and Linode Configs.",
			Optional:    true,
			Computed:    true,
		},

		"managed": schema.BoolAttribute{
			Description: "Enables monitoring for connectivity, response, and total request time.",
			Computed:    true,
		},
		"longview_subscription": schema.StringAttribute{
			Description: "The Longview Pro tier you are currently subscribed to.",
			Computed:    true,
			Optional:    true,
		},
		"object_storage": schema.StringAttribute{
			Description: "A string describing the status of this account's Object Storage service enrollment.",
			Computed:    true,
		},
	},
}
