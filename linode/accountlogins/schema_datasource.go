package accountlogins

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/accountlogin"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var filterConfig = helper.FilterConfig{
	"ip":         {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"restricted": {APIFilterable: false, TypeFunc: helper.FilterTypeBool},
	"username":   {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"status":     {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = map[string]*schema.Schema{
	"filter": filterConfig.FilterSchema(),
	"logins": {
		Type:        schema.TypeList,
		Description: "The returned list of account logins.",
		Computed:    true,
		Elem:        accountlogin.DataSource(),
	},
}
