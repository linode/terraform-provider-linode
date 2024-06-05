package childaccounts

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/account"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"euuid":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"email":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"first_name":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"last_name":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"company":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"address_1":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"address_2":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"phone":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"city":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"state":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"country":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"zip":          {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"capabilities": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"active_since": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"child_accounts": schema.ListNestedBlock{
			Description: "The returned list of Child Accounts.",
			NestedObject: schema.NestedBlockObject{
				Attributes: account.DataSourceSchema().Attributes,
			},
		},
	},
}
