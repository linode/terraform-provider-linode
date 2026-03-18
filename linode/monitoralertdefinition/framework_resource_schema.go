package monitoralertdefinition

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var dimensionFilterNestedObj = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"dimension_label": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The name of the dimension to be used in the filter.",
		},
		"operator": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The operator to apply to the dimension filter. Allowed values: eq, neq, startswith, endswith.",
			Validators: []validator.String{
				stringvalidator.OneOf("eq", "neq", "startswith", "endswith"),
			},
		},
		"value": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The value to compare the dimension_label against.",
		},
		"label": schema.StringAttribute{
			Computed:    true,
			Description: "The name of the dimension filter. Used for display purposes.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}

var ruleNestedObj = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"aggregate_function": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The aggregation function applied to the metric.",
			Validators: []validator.String{
				stringvalidator.OneOf("avg", "sum", "min", "max"),
			},
		},
		"dimension_filters": schema.ListNestedAttribute{
			NestedObject: dimensionFilterNestedObj,
			Description:  "Individual objects that define dimension filters for the rule.",
			Optional:     true,
			Computed:     true,
		},
		"metric": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The metric to query.",
		},
		"operator": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The operator to apply to the metric. Allowed values: eq, gt, lt, gte, lte.",
			Validators: []validator.String{
				stringvalidator.OneOf("eq", "gt", "lt", "gte", "lte"),
			},
		},
		"threshold": schema.Float64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The predefined value or condition that triggers an alert when met or exceeded.",
		},
		"label": schema.StringAttribute{
			Description: "The name of the individual rule. This is used for display purposes in Akamai Cloud Manager.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Computed: true,
		},
		"unit": schema.StringAttribute{
			Description: "The unit of the metric. This can be values like percent for percentage or GB for gigabyte.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Computed: true,
		},
	},
}

var alertChannelNestedObj = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Description: "The unique identifier for this alert channel on your account.",
		},
		"label": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "The name of the channel, used to display it in Akamai Cloud Manager.",
		},
		"type": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "The type of notification used with the channel. For a user alert definition, only email is supported.",
		},
		"url": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "The URL for the channel that ends in the channel's id.",
		},
	},
}

var triggerConditionsObj = map[string]schema.Attribute{
	"criteria_condition": schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "The logical operation applied. Currently only 'ALL' allowed.",
		Validators: []validator.String{
			stringvalidator.OneOf("ALL"),
		},
	},
	"evaluation_period_seconds": schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Description: "Time period over which data is collected before evaluating the threshold.",
	},
	"polling_interval_seconds": schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Description: "Frequency at which the metric is checked.",
	},
	"trigger_occurrences": schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Description: "Minimum consecutive polling periods for threshold breach.",
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed:    true,
			Description: "The unique identifier assigned to the alert definition.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"service_type": schema.StringAttribute{
			Required:    true,
			Description: "The Akamai Cloud Computing service being monitored.",
		},
		"channel_ids": schema.ListAttribute{
			ElementType: types.Int64Type,
			Required:    true,
			Description: "The identifiers for the alert channels to use for the alert.",
		},
		"description": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "An additional description for the alert definition.",
		},
		"entity_ids": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			Description: "The id for each individual entity from a service_type.",
		},
		"label": schema.StringAttribute{
			Required:    true,
			Description: "The name of the alert definition.",
		},
		"status": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The current status of the alert.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"wait_for": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Wait for the alert definition ready (not in progress).",
			Default:     booldefault.StaticBool(false),
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"severity": schema.Int64Attribute{
			Required:    true,
			Description: "The severity of the alert. Supported values: 3 for info, 2 for low, 1 for medium, 0 for severe.",
			Validators: []validator.Int64{
				int64validator.OneOf(0, 1, 2, 3),
			},
		},
		"rule_criteria": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"rules": schema.ListNestedAttribute{
					NestedObject: ruleNestedObj,
					Required:     true,
					Description:  "The individual rules that make up the alert definition.",
				},
			},
			Required:    true,
			Description: "Details for the rules required to trigger the alert.",
		},
		"trigger_conditions": schema.SingleNestedAttribute{
			Attributes:  triggerConditionsObj,
			Required:    true,
			Description: "The conditions that need to be met to send a notification for the alert.",
		},
		"type": schema.StringAttribute{
			Computed: true,
			Description: "The type of alert. This can be either user for an alert specific to the current user, " +
				"or system for one that applies to all users on your account.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"has_more_resources": schema.BoolAttribute{
			Computed: true,
			Description: "Whether there are additional entity_ids associated with the alert for which " +
				"the user doesn't have at least read-only access.",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"alert_channels": schema.ListNestedAttribute{
			NestedObject: alertChannelNestedObj,
			Computed:     true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			Description: "The alert channels set up for use with this alert. Run the List alert channels operation to see all of the available channels.",
		},
		"created": schema.StringAttribute{
			Description: "When this Alert Definition was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "When this Alert Definition was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			Computed: true,
			Description: "For a user alert definition, this is the user on your account that created it. " +
				"For a system alert definition, this is returned as system.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_by": schema.StringAttribute{
			Computed: true,
			Description: "For a user alert definition, this is the user on your account that last updated it. " +
				"For a system alert definition, this is returned as system. If it hasn't been updated, this value is the same as created_by.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"class": schema.StringAttribute{
			Computed: true,
			Description: "The plan type for the Managed Database cluster, either shared or dedicated. " +
				"This only applies to a system alert for a service_type of dbaas (Managed Databases). " +
				"For user alerts for dbaas, this is returned as null.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
