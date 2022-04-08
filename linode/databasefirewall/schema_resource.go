package databasefirewall

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceSchema = map[string]*schema.Schema{
	// Required fields
	"database_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the MySQL database to manage the allow list for.",
		Required:    true,
		ForceNew:    true,
	},
	"allow_list": {
		Type: schema.TypeSet,
		Description: "A list of IP addresses that can access the Managed Database. " +
			"Each item can be a single IP address or a range in CIDR format.",
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
}
