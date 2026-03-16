package firewallrulesets

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/firewall"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"label": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"type":  {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
}

var rulesetAttributes = map[string]schema.Attribute{
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
}

var frameworkDatasourceSchema = schema.Schema{
	Description: "Provides a list of Linode Firewall Rule Sets matching optional filters.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"rulesets": schema.ListNestedAttribute{
			Description: "The returned list of firewall rule sets.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: rulesetAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
