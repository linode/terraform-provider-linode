package rdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var resourceSchema = map[string]*schema.Schema{
	"address": {
		Type:         schema.TypeString,
		Description:  "The public Linode IPv4 or IPv6 address to operate on.",
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.IsIPAddress,
	},
	"rdns": {
		Type: schema.TypeString,
		Description: "The reverse DNS assigned to this address. For public IPv4 addresses, this will be set " +
			"to a default value provided by Linode if not explicitly set.",
		Required:     true,
		ValidateFunc: validation.StringLenBetween(3, 254),
	},
}
