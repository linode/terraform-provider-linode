package regions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var filterConfig = helper.FrameworkFilterConfig{
	"capabilities": {APIFilterable: false, Type: types.StringType},
	"country":      {APIFilterable: false, Type: types.StringType},
	"status":       {APIFilterable: false, Type: types.StringType},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"filter": helper.FrameworkFilterSchema,
	},
}
