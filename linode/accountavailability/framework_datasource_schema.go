package accountavailability

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: AccountAvailabilityAttributes,
}

var AccountAvailabilityAttributes = map[string]schema.Attribute{
	"region": schema.StringAttribute{
		Description: "The region of this availability entry.",
		Required:    true,
	},
	"unavailable": schema.SetAttribute{
		ElementType: types.StringType,
		Description: "A set of unavailable services for the current account in this region.",
		Computed:    true,
	},
	"available": schema.SetAttribute{
		ElementType: types.StringType,
		Description: "A set of available services for the current account in this region.",
		Computed:    true,
	},
}
