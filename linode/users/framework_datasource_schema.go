package users

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/user"
)

var filterConfig = frameworkfilter.Config{
	"username":              {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"email":                 {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"restricted":            {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
	"password_created":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"tfa_enabled":           {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
	"verified_phone_number": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"users": schema.ListNestedBlock{
			Description: "The returned list of Users.",
			NestedObject: schema.NestedBlockObject{
				Attributes: user.UserAttributes,
			},
		},
	},
}
