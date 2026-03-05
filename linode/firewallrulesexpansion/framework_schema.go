package firewallrulesexpansion

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/firewall"
)

var frameworkDatasourceSchema = schema.Schema{
	Description: "Provides the expanded (fully resolved) firewall rules for a Linode Firewall. " +
		"Rule Set references and prefix list tokens are expanded into concrete rules.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"firewall_id": schema.Int64Attribute{
			Description: "The ID of the Firewall to get expanded rules for.",
			Required:    true,
		},
		"inbound": schema.ListAttribute{
			ElementType: firewall.RuleObjectType,
			Description: "The expanded inbound firewall rules.",
			Computed:    true,
		},
		"inbound_policy": schema.StringAttribute{
			Description: "The default behavior for inbound traffic.",
			Computed:    true,
		},
		"outbound": schema.ListAttribute{
			ElementType: firewall.RuleObjectType,
			Description: "The expanded outbound firewall rules.",
			Computed:    true,
		},
		"outbound_policy": schema.StringAttribute{
			Description: "The default behavior for outbound traffic.",
			Computed:    true,
		},
		"version": schema.Int64Attribute{
			Description: "The version of the firewall rules.",
			Computed:    true,
		},
	},
}
