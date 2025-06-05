package firewallsettings

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var DefaultFirewallIDs = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"linode":           types.Int64Type,
		"nodebalancer":     types.Int64Type,
		"public_interface": types.Int64Type,
		"vpc_interface":    types.Int64Type,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"default_firewall_ids": schema.ObjectAttribute{
			Description: "The default firewall ID for a linode, nodebalancer, public_interface, or vpc_interface.",
			Computed:    true,
			AttributeTypes: DefaultFirewallIDs.AttrTypes,
		},
	},
}
