package lkeversion

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var Attributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique ID assigned to this LKE version.",
		Required:    true,
	},
	"tier": schema.StringAttribute{
		Description: "The tier of the LKE version, either standard or enterprise.",
		Computed:    true,
		Optional:    true,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: Attributes,
}
