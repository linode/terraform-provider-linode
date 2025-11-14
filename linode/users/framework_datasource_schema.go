package users

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/user"
)

var filterConfig = frameworkfilter.Config{
	"username":              {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"email":                 {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"restricted":            {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
	"user_type":             {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
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
		"users": schema.ListNestedAttribute{
			Description: "The returned list of Users.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: user.UserAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
