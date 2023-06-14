package accountsettings

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The email of the current account.",
			Computed:    true,
		},
		"backups_enabled": schema.BoolAttribute{
			Description: "Account-wide backups default.",
			Computed:    true,
		},
		"network_helper": schema.BoolAttribute{
			Description: "Enables network helper across all users by default for new Linodes and Linode Configs.",
			Computed:    true,
		},
		"managed": schema.BoolAttribute{
			Description: "Enables monitoring for connectivity, response, and total request time.",
			Computed:    true,
		},
		"longview_subscription": schema.StringAttribute{
			Description: "The Longview Pro tier you are currently subscribed to.",
			Computed:    true,
		},
		"object_storage": schema.StringAttribute{
			Description: "A string describing the status of this account's Object Storage service enrollment.",
			Computed:    true,
		},
	},
}
