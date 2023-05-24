package accountlogin

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"datetime": {
		Type:        schema.TypeString,
		Description: "The time when the login was initiated.",
		Computed:    true,
	},
	"id": {
		Type:        schema.TypeInt,
		Description: "The unique ID of this login object.",
		Required:    true,
	},
	"ip": {
		Type:        schema.TypeString,
		Description: "The remote IP address that requested the login.",
		Computed:    true,
	},
	"restricted": {
		Type:        schema.TypeBool,
		Description: "True if the User that was logged into was a restricted User, false otherwise.",
		Computed:    true,
	},
	"username": {
		Type:        schema.TypeString,
		Description: "The username of the User that was logged into.",
		Computed:    true,
	},
	"status": {
		Type:        schema.TypeString,
		Description: "Whether the login attempt succeeded or failed.",
		Computed:    true,
	},
}
