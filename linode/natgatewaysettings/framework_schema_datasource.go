package natgatewaysettings

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"allowed_ports_per_interface": schema.ListAttribute{
			Description: "The list of port counts allowed per NAT Gateway interface.",
			Computed:    true,
			ElementType: types.Int64Type,
		},
		"maximum_autoscaling_addresses_per_natgateway": schema.Int64Attribute{
			Description: "The maximum number of autoscaling addresses allowed per NAT Gateway.",
			Computed:    true,
		},
		"maximum_reserved_addresses_per_natgateway": schema.Int64Attribute{
			Description: "The maximum number of reserved addresses allowed per NAT Gateway.",
			Computed:    true,
		},
	},
}
