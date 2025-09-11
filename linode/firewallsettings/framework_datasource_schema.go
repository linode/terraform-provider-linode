package firewallsettings

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"default_firewall_ids": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "The default firewall ID for a linode, nodebalancer, public_interface, or vpc_interface.",
			Attributes: map[string]schema.Attribute{
				"linode":           schema.Int64Attribute{Computed: true},
				"nodebalancer":     schema.Int64Attribute{Computed: true},
				"public_interface": schema.Int64Attribute{Computed: true},
				"vpc_interface":    schema.Int64Attribute{Computed: true},
			},
		},
	},
}
