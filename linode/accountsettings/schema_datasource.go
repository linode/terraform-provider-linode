package accountsettings

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"backups_enabled": {
		Type:        schema.TypeBool,
		Description: "Account-wide backups default.",
		Computed:    true,
	},
	"longview_subscription": {
		Type:        schema.TypeString,
		Description: "The Longview Pro tier you are currently subscribed to.",
		Computed:    true,
	},
	"managed": {
		Type:        schema.TypeBool,
		Description: "Enables monitoring for connectivity, response, and total request time.",
		Computed:    true,
	},
	"network_helper": {
		Type:        schema.TypeBool,
		Description: "Enables network helper across all users by default for new Linodes and Linode Configs.",
		Computed:    true,
	},
	"object_storage": {
		Type:        schema.TypeString,
		Description: "A string describing the status of this account’s Object Storage service enrollment.",
		Computed:    true,
	},
}
