package firewalltemplates

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/firewall"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"slug": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var templatesAttributes = map[string]schema.Attribute{
	"slug": schema.StringAttribute{
		Description: "The slug of the firewall template.",
		Computed:    true,
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
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"firewall_templates": schema.ListNestedAttribute{
			Description: "The returned list of firewall templates.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: templatesAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
