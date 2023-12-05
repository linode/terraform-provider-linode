package accountavailability

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: Attributes,
}

var Attributes = map[string]schema.Attribute{
	"region": schema.StringAttribute{
		Description: "The id of a region.",
		Required:    true,
	},
	"unavailable": schema.ListAttribute{
		Description: "The unavailable resources in a region to a specific customer.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
