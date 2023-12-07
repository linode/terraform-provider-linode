package accountavailabilities

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"region":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"unavailable": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

// TODO: Merge with singular data source schema
var accountAvailabilitySchema = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"region": schema.StringAttribute{
			Description: "The region of this availability entry.",
			Computed:    true,
		},
		"unavailable": schema.SetAttribute{
			ElementType: types.StringType,
			Description: "A list of unavailable services for the current account in this region.",
			Computed:    true,
		},
	},
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
			Description:  "The returned list of account availabilities.",
			NestedObject: accountAvailabilitySchema,
		},
	},
}
