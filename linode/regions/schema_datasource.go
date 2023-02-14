package regions

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/region"
)

var filterConfig = helper.FilterConfig{
	"capabilities": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"country":      {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"status":       {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = map[string]*schema.Schema{
	"filter": filterConfig.FilterSchema(),
	"regions": {
		Type:        schema.TypeList,
		Description: "The returned list of regions.",
		Computed:    true,
		Elem:        region.DataSource(),
	},
}
