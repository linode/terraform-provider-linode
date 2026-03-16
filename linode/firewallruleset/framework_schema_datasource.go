package firewallruleset

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/firewall"
)

var frameworkDatasourceSchema = schema.Schema{
	Description: "Provides details about a specific Linode Firewall Rule Set.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of the Rule Set.",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label for the Rule Set.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "A description for this Rule Set.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "Whether the rules in this set are inbound or outbound.",
			Computed:    true,
		},
		"rules": schema.ListAttribute{
			ElementType: firewall.RuleObjectType,
			Description: "The ordered list of firewall rules in this rule set.",
			Computed:    true,
		},
		"is_service_defined": schema.BoolAttribute{
			Description: "Whether this Rule Set is read-only and defined by a Linode service.",
			Computed:    true,
		},
		"version": schema.Int64Attribute{
			Description: "The version of this Rule Set.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When the Rule Set was created.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "When the Rule Set was last updated.",
			Computed:    true,
		},
	},
}
