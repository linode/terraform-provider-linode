package maintenancepolicies

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"slug": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var maintenancePolicyAttributes = map[string]schema.Attribute{
	"slug": schema.StringAttribute{
		Computed:    true,
		Description: "Unique identifier for the policy.",
	},
	"label": schema.StringAttribute{
		Computed:    true,
		Description: "The label for the policy.",
	},
	"description": schema.StringAttribute{
		Computed:    true,
		Description: "Description of the policy.",
	},
	"type": schema.StringAttribute{
		Computed:    true,
		Description: "Type of action taken during maintenance.",
	},
	"notification_period_sec": schema.Int64Attribute{
		Computed:    true,
		Description: "Notification lead time in seconds.",
	},
	"is_default": schema.BoolAttribute{
		Computed:    true,
		Description: "Whether this is the default policy.",
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Unique identification field for this list of Maintenance Policies.",
		},
		"maintenance_policies": schema.ListNestedAttribute{
			Description: "The returned list list of available maintenance policies.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: maintenancePolicyAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
