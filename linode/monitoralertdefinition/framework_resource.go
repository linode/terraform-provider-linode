package monitoralertdefinition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

const (
	DefaultAlertDefinitionTimeout = 30 * time.Minute
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_monitor_alert_definition",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create "+r.Config.Name)
	var data AlertDefinitionResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceType := data.ServiceType.ValueString()

	severity, err := helper.SafeInt64ToInt(data.Severity.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting severity",
			err.Error(),
		)
		return
	}

	createOpts := linodego.AlertDefinitionCreateOptions{
		Description: data.Description.ValueStringPointer(),
		Label:       data.Label.ValueString(),
		Severity:    severity,
	}

	if !data.TriggerConditions.IsUnknown() || !data.TriggerConditions.IsNull() {
		createOpts.TriggerConditions = flattenTriggerConditions(ctx, data.TriggerConditions, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.RuleCriteria.IsUnknown() || !data.RuleCriteria.IsNull() {
		createOpts.RuleCriteria = flattenRuleCriteriaOpts(ctx, data.RuleCriteria, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ChannelIDs.IsUnknown() && !data.ChannelIDs.IsNull() {
		resp.Diagnostics.Append(data.ChannelIDs.ElementsAs(ctx, &createOpts.ChannelIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !data.EntityIDs.IsUnknown() && !data.EntityIDs.IsNull() {
		resp.Diagnostics.Append(data.EntityIDs.ElementsAs(ctx, &createOpts.EntityIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, "client.CreateMonitorAlertDefinition(...)", map[string]any{
		"service_type": serviceType,
		"options":      createOpts,
	})

	alertdefinition, err := client.CreateMonitorAlertDefinition(ctx, serviceType, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating a Monitor Service Alert Definition",
			err.Error(),
		)
		return
	}

	if data.WaitFor.ValueBool() {
		// wait for the alert definition to be ready (not in progress)
		alertdefinition, err = client.WaitForAlertDefinitionStatusReady(ctx, serviceType, alertdefinition.ID, int(DefaultAlertDefinitionTimeout.Seconds()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for Monitor Service Alert Definition to be ready",
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(data.flattenResourceModel(ctx, alertdefinition, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var data AlertDefinitionResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)
	id := data.ID.ValueInt64()

	alertDefinition, err := client.GetMonitorAlertDefinition(ctx, data.ServiceType.ValueString(), int(id))
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Alert Definition No Longer Exists",
				fmt.Sprintf("Removing Alert Definition ID %v from state because it no longer exists", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Get Alert Definition %v for service type %s", id, data.ServiceType.ValueString()),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.flattenResourceModel(ctx, alertDefinition, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	var plan, state AlertDefinitionResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	ctx = populateLogAttributes(ctx, state)

	id := state.ID.ValueInt64()

	isEqual := state.Label.Equal(plan.Label) &&
		state.Description.Equal(plan.Description) &&
		state.Severity.Equal(plan.Severity) &&
		state.TriggerConditions.Equal(plan.TriggerConditions) &&
		state.RuleCriteria.Equal(plan.RuleCriteria) &&
		state.EntityIDs.Equal(plan.EntityIDs) &&
		state.ChannelIDs.Equal(plan.ChannelIDs) &&
		state.Status.Equal(plan.Status)

	if !isEqual {
		severity := helper.FrameworkSafeInt64ToInt(plan.Severity.ValueInt64(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		status := linodego.AlertDefinitionStatus(plan.Status.ValueString())

		updateOpts := linodego.AlertDefinitionUpdateOptions{
			Label:       plan.Label.ValueString(),
			Description: plan.Description.ValueStringPointer(),
			Severity:    severity,
			Status:      &status,
		}

		if !plan.EntityIDs.IsUnknown() || !plan.EntityIDs.IsNull() {
			resp.Diagnostics.Append(plan.EntityIDs.ElementsAs(ctx, &updateOpts.EntityIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if !plan.ChannelIDs.IsUnknown() || !plan.ChannelIDs.IsNull() {
			resp.Diagnostics.Append(plan.ChannelIDs.ElementsAs(ctx, &updateOpts.ChannelIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if !plan.TriggerConditions.IsUnknown() || !plan.TriggerConditions.IsNull() {
			updateOpts.TriggerConditions = flattenTriggerConditions(ctx, plan.TriggerConditions, &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if !plan.RuleCriteria.IsUnknown() || !plan.RuleCriteria.IsNull() {
			updateOpts.RuleCriteria = flattenRuleCriteriaOpts(ctx, plan.RuleCriteria, &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		tflog.Debug(ctx, "client.UpdateMonitorAlertDefinition(...)", map[string]any{
			"options": updateOpts,
		})

		alertDefinition, err := client.UpdateMonitorAlertDefinition(ctx, plan.ServiceType.ValueString(), int(id), updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Update Alert Definition %v for service type %s", id, plan.ServiceType.ValueString()),
				err.Error(),
			)
			return
		}

		if plan.WaitFor.ValueBool() {
			// wait for the alert definition to be ready (not in progress)
			alertDefinition, err = client.WaitForAlertDefinitionStatusReady(
				ctx,
				alertDefinition.ServiceType,
				alertDefinition.ID,
				int(DefaultAlertDefinitionTimeout.Seconds()),
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error waiting for Monitor Service Alert Definition to be ready",
					err.Error(),
				)
				return
			}
		}

		resp.Diagnostics.Append(plan.flattenResourceModel(ctx, alertDefinition, true)...)
	}

	plan.CopyFrom(state, true)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	var data AlertDefinitionResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueInt64()

	ctx = populateLogAttributes(ctx, data)
	tflog.Debug(ctx, "client.DeleteMonitorAlertDefinition(...)")

	err := client.DeleteMonitorAlertDefinition(ctx, data.ServiceType.ValueString(), int(id))
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Alert Definition No Longer Exists",
				fmt.Sprintf("Alert Definition %v does not exist, removing it from state.", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Failed to Delete Alert Definition",
			err.Error(),
		)
		return
	}
}

func populateLogAttributes(ctx context.Context, model AlertDefinitionResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"alert_definition_id": model.ID.ValueInt64(),
	})
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import "+r.Config.Name)

	helper.ImportStateWithMultipleIDs(
		ctx,
		req,
		resp,
		[]helper.ImportableID{
			{
				Name:          "id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "service_type",
				TypeConverter: helper.IDTypeConverterString,
			},
		},
	)
}
