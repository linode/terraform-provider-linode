package accountlogins

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/accountlogin"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"ip":         {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"restricted": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
	"username":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"status":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"logins": schema.ListNestedAttribute{
			Description: "The returned list of account logins.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: accountlogin.Attributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
