package domainzonefile

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"domain_id": schema.Int64Attribute{
			Description: "The domain's ID.",
			Required:    true,
		},
		"zone_file": schema.ListAttribute{
			Description: "Lines of the zone file for the last rendered zone for this domain.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"id": schema.StringAttribute{
			Description: "The unique ID for this DataSource",
			Computed:    true,
		},
	},
}
