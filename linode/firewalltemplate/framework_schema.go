package firewalltemplate

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/firewall"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The computed ID of the data source should have the same value as the `slug` attribute.",
			Computed:    true,
		},
		"slug": schema.StringAttribute{
			Description: "The slug of the firewall template.",
			Required:    true,
		},
		"inbound": schema.ListAttribute{
			ElementType: firewall.RuleObjectType,
			Description: "A firewall rule that specifies what inbound network traffic is allowed.",
			Computed:    true,
		},
		"inbound_policy": schema.StringAttribute{
			Description: "The default behavior for inbound traffic. This setting can be overridden by updating " +
				"the inbound.action property for an individual Firewall Rule.",
			Computed: true,
		},
		"outbound": schema.ListAttribute{
			ElementType: firewall.RuleObjectType,
			Description: "A firewall rule that specifies what outbound network traffic is allowed.",
			Computed:    true,
		},
		"outbound_policy": schema.StringAttribute{
			Description: "The default behavior for outbound traffic. This setting can be overridden by updating " +
				"the outbound.action property for an individual Firewall Rule.",
			Computed: true,
		},
	},
}
