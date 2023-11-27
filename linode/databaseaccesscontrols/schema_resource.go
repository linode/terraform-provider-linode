package databaseaccesscontrols

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var resourceSchema = map[string]*schema.Schema{
	// Required fields
	"database_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the database to manage the allow list for.",
		Required:    true,
		ForceNew:    true,
	},
	"database_type": {
		Type:        schema.TypeString,
		Description: "The type of the  database to manage the allow list for.",
		Required:    true,
		ForceNew:    true,
		ValidateDiagFunc: validation.ToDiagFunc(
			validation.StringInSlice(helper.ValidDatabaseTypes, true)),
	},

	"allow_list": {
		Type: schema.TypeSet,
		Description: "A list of IP addresses that can access the Managed Database. " +
			"Each item can be a single IP address or a range in CIDR format.",
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
}
