package regions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

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
