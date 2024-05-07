package accountsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_account_settings",
				IDType: types.StringType,
				Schema: frameworkResourceSchema,
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
	tflog.Debug(ctx, "Create linode_account_settings")

	var plan AccountSettingsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the account
	resp.Diagnostics.Append(
		r.updateAccountSettings(ctx, &plan)...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply the state changes
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_account_settings")

	client := r.Meta.Client

	var data AccountSettingsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	tflog.Trace(ctx, "client.GetAccount(...)")
	account, err := client.GetAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Linode Account",
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "client.GetAccountSettings(...)")
	settings, err := client.GetAccountSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Linode Account Settings",
			err.Error(),
		)
		return
	}

	data.FlattenAccountSettings(account.Email, settings, false)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	data.ID = types.StringValue(account.Email)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_account_settings")

	var plan, state AccountSettingsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the account
	resp.Diagnostics.Append(
		r.updateAccountSettings(ctx, &plan)...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.CopyFrom(state, true)

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	// Apply the state changes
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_account_settings")
}

func (r *Resource) updateAccountSettings(
	ctx context.Context,
	plan *AccountSettingsModel,
) diag.Diagnostics {
	client := r.Meta.Client
	var diagnostics diag.Diagnostics

	// Longview Plan update functionality has been moved
	if !plan.LongviewSubscription.IsNull() {
		options := linodego.LongviewPlanUpdateOptions{
			LongviewSubscription: plan.LongviewSubscription.ValueString(),
		}

		tflog.Debug(ctx, "client.UpdateLongviewPlan(...)", map[string]any{
			"options": options,
		})

		_, err := client.UpdateLongviewPlan(ctx, options)
		if err != nil {
			diagnostics.AddError(
				"Failed to update Linode Longview Plan",
				err.Error(),
			)
			return diagnostics
		}
	}

	updateOpts := linodego.AccountSettingsUpdateOptions{
		BackupsEnabled: plan.BackupsEnabled.ValueBoolPointer(),
		NetworkHelper:  plan.NetworkHelper.ValueBoolPointer(),
	}

	tflog.Debug(ctx, "client.UpdateAccountSettings(...)", map[string]any{
		"options": updateOpts,
	})

	settings, err := client.UpdateAccountSettings(ctx, updateOpts)
	if err != nil {
		diagnostics.AddError("Failed to update Linode Account Settings", err.Error())
		return diagnostics
	}

	plan.FlattenAccountSettings(plan.ID.ValueString(), settings, true)
	return diagnostics
}
