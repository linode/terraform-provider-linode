package accountavailabilities

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/accountavailability"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"region":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"unavailable": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"available":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"availabilities": schema.ListNestedBlock{
			Description: "The returned list of account availabilities.",
			NestedObject: schema.NestedBlockObject{
				Attributes: accountavailability.AccountAvailabilityAttributes,
			},
		},
	},
}
