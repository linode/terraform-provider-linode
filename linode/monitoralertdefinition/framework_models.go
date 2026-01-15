package monitoralertdefinition

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

// AlertDefinitionResourceModel describes the Terraform resource model to match the resource schema.
type AlertDefinitionResourceModel struct {
	AlertDefinitionDataSourceModel
	WaitFor types.Bool `tfsdk:"wait_for"`
}

// AlertDefinitionDataSourceModel describes the Terraform data source model to match the data source schema.
type AlertDefinitionDataSourceModel struct {
	ID                types.Int64       `tfsdk:"id"`
	ServiceType       types.String      `tfsdk:"service_type"`
	ChannelIDs        types.List        `tfsdk:"channel_ids"`
	Description       types.String      `tfsdk:"description"`
	EntityIDs         types.List        `tfsdk:"entity_ids"`
	Label             types.String      `tfsdk:"label"`
	Status            types.String      `tfsdk:"status"`
	Severity          types.Int64       `tfsdk:"severity"`
	RuleCriteria      types.Object      `tfsdk:"rule_criteria"`
	TriggerConditions types.Object      `tfsdk:"trigger_conditions"`
	Type              types.String      `tfsdk:"type"`
	HasMoreResources  types.Bool        `tfsdk:"has_more_resources"`
	AlertChannels     types.List        `tfsdk:"alert_channels"`
	Created           timetypes.RFC3339 `tfsdk:"created"`
	Updated           timetypes.RFC3339 `tfsdk:"updated"`
	CreatedBy         types.String      `tfsdk:"created_by"`
	UpdatedBy         types.String      `tfsdk:"updated_by"`
	Class             types.String      `tfsdk:"class"`
}

type RuleCriteriaModel struct {
	Rules types.List `tfsdk:"rules"`
}

type RuleModel struct {
	AggregateFunction types.String  `tfsdk:"aggregate_function"`
	DimensionFilters  types.List    `tfsdk:"dimension_filters"`
	Metric            types.String  `tfsdk:"metric"`
	Operator          types.String  `tfsdk:"operator"`
	Threshold         types.Float64 `tfsdk:"threshold"`
	Label             types.String  `tfsdk:"label"`
	Unit              types.String  `tfsdk:"unit"`
}

type DimensionFilterModel struct {
	DimensionLabel types.String `tfsdk:"dimension_label"`
	Operator       types.String `tfsdk:"operator"`
	Value          types.String `tfsdk:"value"`
	Label          types.String `tfsdk:"label"`
}

type AlertChannelModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
	Type  types.String `tfsdk:"type"`
	Url   types.String `tfsdk:"url"`
}

type TriggerConditionsModel struct {
	CriteriaCondition       types.String `tfsdk:"criteria_condition"`
	EvaluationPeriodSeconds types.Int64  `tfsdk:"evaluation_period_seconds"`
	PollingIntervalSeconds  types.Int64  `tfsdk:"polling_interval_seconds"`
	TriggerOccurrences      types.Int64  `tfsdk:"trigger_occurrences"`
}

var alertChannelObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":    types.Int64Type,
		"label": types.StringType,
		"type":  types.StringType,
		"url":   types.StringType,
	},
}

var ruleCriteriaObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"rules": types.ListType{
			ElemType: ruleObjectType,
		},
	},
}

var ruleObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"aggregate_function": types.StringType,
		"dimension_filters": types.ListType{
			ElemType: dimensionFilterObjectType,
		},
		"metric":    types.StringType,
		"operator":  types.StringType,
		"threshold": types.Float64Type,
		"label":     types.StringType,
		"unit":      types.StringType,
	},
}

var dimensionFilterObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"dimension_label": types.StringType,
		"operator":        types.StringType,
		"value":           types.StringType,
		"label":           types.StringType,
	},
}

var triggerConditionsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"criteria_condition":        types.StringType,
		"evaluation_period_seconds": types.Int64Type,
		"polling_interval_seconds":  types.Int64Type,
		"trigger_occurrences":       types.Int64Type,
	},
}

func (data *AlertDefinitionDataSourceModel) FlattenDataSourceModel(
	ctx context.Context,
	alertdefinition *linodego.AlertDefinition,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics
	data.ID = helper.KeepOrUpdateInt64(data.ID, int64(alertdefinition.ID), preserveKnown)
	data.ServiceType = helper.KeepOrUpdateString(data.ServiceType, alertdefinition.ServiceType, preserveKnown)
	data.Description = helper.KeepOrUpdateString(data.Description, alertdefinition.Description, preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, alertdefinition.Label, preserveKnown)
	data.Status = helper.KeepOrUpdateString(data.Status, string(alertdefinition.Status), preserveKnown)
	data.Severity = helper.KeepOrUpdateInt64(data.Severity, int64(alertdefinition.Severity), preserveKnown)
	data.Type = helper.KeepOrUpdateString(data.Type, alertdefinition.Type, preserveKnown)
	data.HasMoreResources = helper.KeepOrUpdateBool(data.HasMoreResources, alertdefinition.HasMoreResources, preserveKnown)
	data.CreatedBy = helper.KeepOrUpdateString(data.CreatedBy, alertdefinition.CreatedBy, preserveKnown)
	data.UpdatedBy = helper.KeepOrUpdateString(data.UpdatedBy, alertdefinition.UpdatedBy, preserveKnown)
	data.Class = helper.KeepOrUpdateString(data.Class, alertdefinition.Class, preserveKnown)

	data.EntityIDs = helper.KeepOrUpdateList(
		types.StringType,
		data.EntityIDs,
		helper.StringSliceToFrameworkValueSlice(alertdefinition.EntityIDs),
		preserveKnown,
		&diags,
	)

	if diags.HasError() {
		return diags
	}

	data.Created = helper.KeepOrUpdateValue(
		data.Created, timetypes.NewRFC3339TimePointerValue(alertdefinition.Created), preserveKnown,
	)
	data.Updated = helper.KeepOrUpdateValue(
		data.Updated, timetypes.NewRFC3339TimePointerValue(alertdefinition.Updated), preserveKnown,
	)

	channels, channelIDs, d := flattenAlertChannels(ctx, alertdefinition.AlertChannels)
	if d != nil {
		diags.Append(d...)
		return diags
	}

	cIDs, newDiags := types.ListValue(types.Int64Type, helper.IntSliceToFrameworkValueSlice(channelIDs))
	diags.Append(newDiags...)
	if diags.HasError() {
		return diags
	}

	data.ChannelIDs = helper.KeepOrUpdateValue(
		data.ChannelIDs,
		cIDs,
		preserveKnown,
	)

	data.AlertChannels = helper.KeepOrUpdateValue(
		data.AlertChannels,
		channels,
		preserveKnown,
	)

	ruleCriteria := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx,
		data.RuleCriteria,
		ruleCriteriaObjectType.AttrTypes,
		false,
		&diags,
		func(m *RuleCriteriaModel, isNull *bool, p bool, d *diag.Diagnostics) {
			m.flattenRuleCriteria(ctx, alertdefinition.RuleCriteria, p, d)
		},
	)

	if ruleCriteria != nil {
		data.RuleCriteria = *ruleCriteria
	} else {
		data.RuleCriteria = types.ObjectNull(ruleCriteriaObjectType.AttrTypes)
	}

	triggerConditions := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx,
		data.TriggerConditions,
		triggerConditionsObjectType.AttrTypes,
		preserveKnown,
		&diags,
		func(triggerConditionsModel *TriggerConditionsModel, _ *bool, preserveKnown bool, _ *diag.Diagnostics) {
			triggerConditionsModel.flattenTriggerConditionsModel(&alertdefinition.TriggerConditions, preserveKnown)
		},
	)

	if triggerConditions != nil {
		data.TriggerConditions = *triggerConditions
	} else {
		data.TriggerConditions = types.ObjectNull(triggerConditionsObjectType.AttrTypes)
	}

	return diags
}

func (data *AlertDefinitionResourceModel) flattenResourceModel(
	ctx context.Context,
	alertdefinition *linodego.AlertDefinition,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics
	data.ID = helper.KeepOrUpdateInt64(data.ID, int64(alertdefinition.ID), preserveKnown)
	data.ServiceType = helper.KeepOrUpdateString(data.ServiceType, alertdefinition.ServiceType, preserveKnown)
	data.Description = helper.KeepOrUpdateString(data.Description, alertdefinition.Description, preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, alertdefinition.Label, preserveKnown)
	data.Status = helper.KeepOrUpdateString(data.Status, string(alertdefinition.Status), preserveKnown)
	data.Severity = helper.KeepOrUpdateInt64(data.Severity, int64(alertdefinition.Severity), preserveKnown)
	data.Type = helper.KeepOrUpdateString(data.Type, alertdefinition.Type, preserveKnown)
	data.HasMoreResources = helper.KeepOrUpdateBool(data.HasMoreResources, alertdefinition.HasMoreResources, preserveKnown)
	data.CreatedBy = helper.KeepOrUpdateString(data.CreatedBy, alertdefinition.CreatedBy, preserveKnown)
	data.UpdatedBy = helper.KeepOrUpdateString(data.UpdatedBy, alertdefinition.UpdatedBy, preserveKnown)
	data.Class = helper.KeepOrUpdateString(data.Class, alertdefinition.Class, preserveKnown)
	data.WaitFor = helper.KeepOrUpdateBool(data.WaitFor, data.WaitFor.ValueBool(), preserveKnown)
	data.EntityIDs = helper.KeepOrUpdateList(
		types.StringType,
		data.EntityIDs,
		helper.StringSliceToFrameworkValueSlice(alertdefinition.EntityIDs),
		preserveKnown,
		&diags,
	)

	if diags.HasError() {
		return diags
	}

	data.Created = helper.KeepOrUpdateValue(
		data.Created, timetypes.NewRFC3339TimePointerValue(alertdefinition.Created), preserveKnown,
	)
	data.Updated = helper.KeepOrUpdateValue(
		data.Updated, timetypes.NewRFC3339TimePointerValue(alertdefinition.Updated), preserveKnown,
	)

	channels, channelIDs, d := flattenAlertChannels(ctx, alertdefinition.AlertChannels)
	if d != nil {
		diags.Append(d...)
		return diags
	}

	cIDs, newDiags := types.ListValue(types.Int64Type, helper.IntSliceToFrameworkValueSlice(channelIDs))
	diags.Append(newDiags...)
	if diags.HasError() {
		return diags
	}

	data.ChannelIDs = helper.KeepOrUpdateValue(
		data.ChannelIDs,
		cIDs,
		preserveKnown,
	)

	data.AlertChannels = helper.KeepOrUpdateValue(
		data.AlertChannels,
		channels,
		preserveKnown,
	)

	ruleCriteria := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx,
		data.RuleCriteria,
		ruleCriteriaObjectType.AttrTypes,
		false,
		&diags,
		func(m *RuleCriteriaModel, isNull *bool, p bool, d *diag.Diagnostics) {
			m.flattenRuleCriteria(ctx, alertdefinition.RuleCriteria, p, d)
		},
	)

	if ruleCriteria != nil {
		data.RuleCriteria = *ruleCriteria
	}

	triggerConditions := helper.KeepOrUpdateSingleNestedAttributes(
		ctx,
		data.TriggerConditions,
		preserveKnown,
		&diags,
		func(triggerConditionsModel *TriggerConditionsModel, _ *bool, preserveKnown bool, _ *diag.Diagnostics) {
			triggerConditionsModel.flattenTriggerConditionsModel(&alertdefinition.TriggerConditions, preserveKnown)
		},
	)

	if triggerConditions != nil {
		data.TriggerConditions = *triggerConditions
	}

	return diags
}

func (m *TriggerConditionsModel) flattenTriggerConditionsModel(
	triggerConditions *linodego.TriggerConditions,
	preserveKnown bool,
) {
	if triggerConditions == nil {
		return
	}

	m.CriteriaCondition = helper.KeepOrUpdateString(
		m.CriteriaCondition,
		triggerConditions.CriteriaCondition,
		preserveKnown,
	)
	m.EvaluationPeriodSeconds = helper.KeepOrUpdateInt64(
		m.EvaluationPeriodSeconds,
		int64(triggerConditions.EvaluationPeriodSeconds),
		preserveKnown,
	)
	m.PollingIntervalSeconds = helper.KeepOrUpdateInt64(
		m.PollingIntervalSeconds,
		int64(triggerConditions.PollingIntervalSeconds),
		preserveKnown,
	)
	m.TriggerOccurrences = helper.KeepOrUpdateInt64(
		m.TriggerOccurrences,
		int64(triggerConditions.TriggerOccurrences),
		preserveKnown,
	)
}

func (m *RuleCriteriaModel) flattenRuleCriteria(
	ctx context.Context,
	ruleCriteria linodego.RuleCriteria,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	if m.Rules.IsNull() {
		m.Rules = types.ListNull(ruleObjectType)
	}

	ruleCriteriaModels := make([]RuleModel, len(ruleCriteria.Rules))

	if !m.Rules.IsNull() && !m.Rules.IsUnknown() {
		diags.Append(m.Rules.ElementsAs(ctx, &ruleCriteriaModels, false)...)
		if diags.HasError() {
			return
		}
	}

	for i, rule := range ruleCriteria.Rules {
		if ruleCriteriaModels[i].DimensionFilters.IsNull() {
			ruleCriteriaModels[i].DimensionFilters = types.ListNull(dimensionFilterObjectType)
		}

		dimensionFilterModels := make([]DimensionFilterModel, len(rule.DimensionFilters))
		if !ruleCriteriaModels[i].DimensionFilters.IsNull() && !ruleCriteriaModels[i].DimensionFilters.IsUnknown() {
			diags.Append(ruleCriteriaModels[i].DimensionFilters.ElementsAs(ctx, &dimensionFilterModels, false)...)
			if diags.HasError() {
				return
			}
		}

		for j, filter := range rule.DimensionFilters {
			dimensionFilterModels[j].DimensionLabel = helper.KeepOrUpdateString(dimensionFilterModels[j].DimensionLabel, filter.DimensionLabel, preserveKnown)
			dimensionFilterModels[j].Operator = helper.KeepOrUpdateString(dimensionFilterModels[j].Operator, filter.Operator, preserveKnown)
			dimensionFilterModels[j].Value = helper.KeepOrUpdateString(dimensionFilterModels[j].Value, filter.Value, preserveKnown)
			dimensionFilterModels[j].Label = helper.KeepOrUpdateString(dimensionFilterModels[j].Label, filter.Label, preserveKnown)
		}

		dimFiltersList, d := types.ListValueFrom(ctx, dimensionFilterObjectType, dimensionFilterModels)
		if d.HasError() {
			diags.Append(d...)
			return
		}

		ruleCriteriaModels[i].DimensionFilters = helper.KeepOrUpdateValue(
			ruleCriteriaModels[i].DimensionFilters,
			dimFiltersList,
			preserveKnown,
		)

		if diags.HasError() {
			return
		}

		ruleCriteriaModels[i].AggregateFunction = helper.KeepOrUpdateString(ruleCriteriaModels[i].AggregateFunction, rule.AggregateFunction, preserveKnown)
		ruleCriteriaModels[i].Metric = helper.KeepOrUpdateString(ruleCriteriaModels[i].Metric, rule.Metric, preserveKnown)
		ruleCriteriaModels[i].Operator = helper.KeepOrUpdateString(ruleCriteriaModels[i].Operator, rule.Operator, preserveKnown)
		ruleCriteriaModels[i].Threshold = helper.KeepOrUpdateFloat64(ruleCriteriaModels[i].Threshold, rule.Threshold, preserveKnown)
		ruleCriteriaModels[i].Label = helper.KeepOrUpdateString(ruleCriteriaModels[i].Label, rule.Label, preserveKnown)
		ruleCriteriaModels[i].Unit = helper.KeepOrUpdateString(ruleCriteriaModels[i].Unit, rule.Unit, preserveKnown)
	}

	rulesList, d := types.ListValueFrom(ctx, ruleObjectType, ruleCriteriaModels)
	if d.HasError() {
		diags.Append(d...)
		return
	}

	m.Rules = helper.KeepOrUpdateValue(
		m.Rules,
		rulesList,
		preserveKnown,
	)
}

func flattenAlertChannels(
	ctx context.Context,
	alertChannels []linodego.AlertChannelEnvelope,
) (types.List, []int, diag.Diagnostics) {
	aChannelModels := make([]AlertChannelModel, len(alertChannels))
	channelIDs := make([]int, len(alertChannels))

	for i, channel := range alertChannels {
		aChannelModels[i].ID = types.Int64Value(int64(channel.ID))
		aChannelModels[i].Label = types.StringValue(channel.Label)
		aChannelModels[i].Type = types.StringValue(channel.Type)
		aChannelModels[i].Url = types.StringValue(channel.URL)

		// Collect channel IDs
		channelIDs[i] = channel.ID
	}

	aChannelList, diag := types.ListValueFrom(ctx, alertChannelObjectType, aChannelModels)

	if diag != nil {
		return types.ListNull(alertChannelObjectType), nil, diag
	}

	return aChannelList, channelIDs, nil
}

func flattenRuleCriteriaOpts(
	ctx context.Context,
	ruleCriteria types.Object,
	diags *diag.Diagnostics,
) *linodego.RuleCriteriaOptions {
	var ruleCriteriaModel RuleCriteriaModel
	diags.Append(ruleCriteria.As(ctx, &ruleCriteriaModel, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)

	if diags.HasError() {
		return nil
	}

	var ruleModels []RuleModel
	diags.Append(ruleCriteriaModel.Rules.ElementsAs(ctx, &ruleModels, false)...)

	if diags.HasError() {
		return nil
	}

	ruleOpts := make([]linodego.RuleOptions, len(ruleModels))

	for i, rule := range ruleModels {
		var dFilterModels []DimensionFilterModel
		diags.Append(rule.DimensionFilters.ElementsAs(ctx, &dFilterModels, false)...)
		if diags.HasError() {
			return nil
		}
		dimensionFilters := make([]linodego.DimensionFilterOptions, len(dFilterModels))
		for j, filter := range dFilterModels {
			dimensionFilters[j] = linodego.DimensionFilterOptions{
				DimensionLabel: filter.DimensionLabel.ValueString(),
				Operator:       filter.Operator.ValueString(),
				Value:          filter.Value.ValueString(),
			}
		}

		ruleOpts[i].DimensionFilters = dimensionFilters
		ruleOpts[i].AggregateFunction = rule.AggregateFunction.ValueString()
		ruleOpts[i].Metric = rule.Metric.ValueString()
		ruleOpts[i].Operator = rule.Operator.ValueString()
		ruleOpts[i].Threshold = rule.Threshold.ValueFloat64()
	}

	return &linodego.RuleCriteriaOptions{
		Rules: ruleOpts,
	}
}

func flattenTriggerConditions(
	ctx context.Context,
	triggerConditions types.Object,
	diags *diag.Diagnostics,
) *linodego.TriggerConditions {
	var triggerConditionsModel TriggerConditionsModel
	var triggerConditionOpts linodego.TriggerConditions

	diags.Append(triggerConditions.As(ctx, &triggerConditionsModel, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)

	if diags.HasError() {
		return nil
	}

	if !triggerConditionsModel.CriteriaCondition.IsUnknown() {
		triggerConditionOpts.CriteriaCondition = triggerConditionsModel.CriteriaCondition.ValueString()
	}
	if !triggerConditionsModel.EvaluationPeriodSeconds.IsUnknown() {
		evaluationPeriodSec, err := helper.SafeInt64ToInt(triggerConditionsModel.EvaluationPeriodSeconds.ValueInt64())
		if err != nil {
			diags.AddError(
				"Error converting evaluation period seconds",
				err.Error(),
			)
		}
		triggerConditionOpts.EvaluationPeriodSeconds = evaluationPeriodSec
	}

	if !triggerConditionsModel.PollingIntervalSeconds.IsUnknown() {
		pollingIntervalSec, err := helper.SafeInt64ToInt(triggerConditionsModel.PollingIntervalSeconds.ValueInt64())
		if err != nil {
			diags.AddError(
				"Error converting polling interval seconds",
				err.Error(),
			)
		}
		triggerConditionOpts.PollingIntervalSeconds = pollingIntervalSec
	}

	if !triggerConditionsModel.TriggerOccurrences.IsUnknown() {
		triggerOccurrences, err := helper.SafeInt64ToInt(triggerConditionsModel.TriggerOccurrences.ValueInt64())
		if err != nil {
			diags.AddError(
				"Error converting trigger occurrences",
				err.Error(),
			)
		}
		triggerConditionOpts.TriggerOccurrences = triggerOccurrences
	}

	return &triggerConditionOpts
}

func (data *AlertDefinitionResourceModel) CopyFrom(other AlertDefinitionResourceModel, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.ServiceType = helper.KeepOrUpdateValue(data.ServiceType, other.ServiceType, preserveKnown)
	data.Description = helper.KeepOrUpdateValue(data.Description, other.Description, preserveKnown)
	data.Label = helper.KeepOrUpdateValue(data.Label, other.Label, preserveKnown)
	data.Status = helper.KeepOrUpdateValue(data.Status, other.Status, preserveKnown)
	data.Severity = helper.KeepOrUpdateValue(data.Severity, other.Severity, preserveKnown)
	data.Type = helper.KeepOrUpdateValue(data.Type, other.Type, preserveKnown)
	data.HasMoreResources = helper.KeepOrUpdateValue(data.HasMoreResources, other.HasMoreResources, preserveKnown)
	data.CreatedBy = helper.KeepOrUpdateValue(data.CreatedBy, other.CreatedBy, preserveKnown)
	data.UpdatedBy = helper.KeepOrUpdateValue(data.UpdatedBy, other.UpdatedBy, preserveKnown)
	data.Class = helper.KeepOrUpdateValue(data.Class, other.Class, preserveKnown)
	data.EntityIDs = helper.KeepOrUpdateValue(data.EntityIDs, other.EntityIDs, preserveKnown)
	data.ChannelIDs = helper.KeepOrUpdateValue(data.ChannelIDs, other.ChannelIDs, preserveKnown)
	data.AlertChannels = helper.KeepOrUpdateValue(data.AlertChannels, other.AlertChannels, preserveKnown)
	data.RuleCriteria = helper.KeepOrUpdateValue(data.RuleCriteria, other.RuleCriteria, preserveKnown)
	data.TriggerConditions = helper.KeepOrUpdateValue(data.TriggerConditions, other.TriggerConditions, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(data.Created, other.Created, preserveKnown)
	data.Updated = helper.KeepOrUpdateValue(data.Updated, other.Updated, preserveKnown)
	data.WaitFor = helper.KeepOrUpdateValue(data.WaitFor, other.WaitFor, preserveKnown)
}
