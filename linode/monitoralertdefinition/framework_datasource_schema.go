package monitoralertdefinition

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var dimensionFilterDataSourceNestedObj = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"dimension_label": schema.StringAttribute{
			Computed:    true,
			Description: "The name of the dimension to be used in the filter.",
		},
		"operator": schema.StringAttribute{
			Computed:    true,
			Description: "The operator to apply to the dimension filter. Allowed values: eq, neq, startswith, endswith.",
		},
		"value": schema.StringAttribute{
			Computed:    true,
			Description: "The value to compare the dimension_label against.",
		},
		"label": schema.StringAttribute{
			Computed:    true,
			Description: "The name of the dimension filter. Used for display purposes.",
		},
	},
}

var alertChannelDataSourceNestedObj = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed:    true,
			Description: "The unique identifier for this alert channel on your account.",
		},
		"label": schema.StringAttribute{
			Computed:    true,
			Description: "The name of the channel, used to display it in Akamai Cloud Manager.",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: "The type of notification used with the channel. For a user alert definition, only email is supported.",
		},
		"url": schema.StringAttribute{
			Computed:    true,
			Description: "The URL for the channel that ends in the channel's id.",
		},
	},
}

var ruleDataSourceNestedObj = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"aggregate_function": schema.StringAttribute{
			Computed:    true,
			Description: "The aggregation function applied to the metric.",
		},
		"dimension_filters": schema.ListNestedAttribute{
			NestedObject: dimensionFilterDataSourceNestedObj,
			Description:  "Individual objects that define dimension filters for the rule.",
			Computed:     true,
		},
		"metric": schema.StringAttribute{
			Computed:    true,
			Description: "The metric to query.",
		},
		"operator": schema.StringAttribute{
			Computed:    true,
			Description: "The operator to apply to the metric. Allowed values: eq, gt, lt, gte, lte.",
		},
		"threshold": schema.Float64Attribute{
			Computed:    true,
			Description: "The predefined value or condition that triggers an alert when met or exceeded.",
		},
		"label": schema.StringAttribute{
			Computed:    true,
			Description: "The name of the individual rule. This is used for display purposes in Akamai Cloud Manager.",
		},
		"unit": schema.StringAttribute{
			Computed:    true,
			Description: "The unit of the metric. This can be values like percent for percentage or GB for gigabyte.",
		},
	},
}

var FrameworkDataSourceSchema = schema.Schema{
	Attributes: AlertDefinitionAttributes,
}

var AlertDefinitionAttributes = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Required:    true,
		Description: "The unique identifier assigned to the alert definition.",
	},
	"service_type": schema.StringAttribute{
		Required:    true,
		Description: "The Akamai Cloud Computing service being monitored.",
	},
	"channel_ids": schema.ListAttribute{
		ElementType: types.Int64Type,
		Computed:    true,
		Description: "The identifiers for the alert channels to use for the alert.",
	},
	"description": schema.StringAttribute{
		Computed:    true,
		Description: "An additional description for the alert definition.",
	},
	"entity_ids": schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "The id for each individual entity from a service_type.",
	},
	"label": schema.StringAttribute{
		Computed:    true,
		Description: "The name of the alert definition.",
	},
	"status": schema.StringAttribute{
		Computed:    true,
		Description: "The current status of the alert.",
	},
	"severity": schema.Int64Attribute{
		Computed:    true,
		Description: "The severity of the alert. Supported values: 3 for info, 2 for low, 1 for medium, 0 for severe.",
	},
	"rule_criteria": schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"rules": schema.ListNestedAttribute{
				NestedObject: ruleDataSourceNestedObj,
				Computed:     true,
				Description:  "The individual rules that make up the alert definition.",
			},
		},
		Computed:    true,
		Description: "Details for the rules required to trigger the alert.",
	},
	"trigger_conditions": schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"criteria_condition": schema.StringAttribute{
				Computed:    true,
				Description: "The logical operation applied. Currently only 'ALL' allowed.",
			},
			"evaluation_period_seconds": schema.Int64Attribute{
				Computed:    true,
				Description: "Time period over which data is collected before evaluating the threshold.",
			},
			"polling_interval_seconds": schema.Int64Attribute{
				Computed:    true,
				Description: "Frequency at which the metric is checked.",
			},
			"trigger_occurrences": schema.Int64Attribute{
				Computed:    true,
				Description: "Minimum consecutive polling periods for threshold breach.",
			},
		},
		Computed:    true,
		Description: "The conditions that need to be met to send a notification for the alert.",
	},
	"type": schema.StringAttribute{
		Computed:    true,
		Description: "The type of alert.",
	},
	"has_more_resources": schema.BoolAttribute{
		Computed:    true,
		Description: "Whether there are additional entity_ids associated with the alert.",
	},
	"alert_channels": schema.ListNestedAttribute{
		NestedObject: alertChannelDataSourceNestedObj,
		Computed:     true,
		Description:  "The alert channels set up for use with this alert.",
	},
	"created": schema.StringAttribute{
		Computed:    true,
		Description: "When this Alert Definition was created.",
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Computed:    true,
		Description: "When this Alert Definition was last updated.",
		CustomType:  timetypes.RFC3339Type{},
	},
	"created_by": schema.StringAttribute{
		Computed:    true,
		Description: "The user (or system) that created this alert definition.",
	},
	"updated_by": schema.StringAttribute{
		Computed:    true,
		Description: "The user (or system) that last updated this alert definition.",
	},
	"class": schema.StringAttribute{
		Computed:    true,
		Description: "Plan type for Managed Database clusters (shared or dedicated).",
	},
}
