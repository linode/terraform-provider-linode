package account_logins

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/account_login"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var filterConfig = helper.FilterConfig{
	"ip":         {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"restricted": {APIFilterable: false, TypeFunc: helper.FilterTypeBool},
	"username":   {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = map[string]*schema.Schema{
	"order_by": filterConfig.OrderBySchema(),
	"order":    filterConfig.OrderSchema(),
	"filter":   filterConfig.FilterSchema(),
	"logins": {
		Type:        schema.TypeList,
		Description: "The returned list of account logins.",
		Computed:    true,
		Elem:        account_login.DataSource(),
	},
}
