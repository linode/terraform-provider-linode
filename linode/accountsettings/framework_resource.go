package accountsettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_account_settings",
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
	client := r.Meta.Client

	var data AccountSettingsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := client.GetAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Linode Account",
			err.Error(),
		)
		return
	}

	settings, err := client.GetAccountSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Linode Account Settings",
			err.Error(),
		)
		return
	}

	data.parseAccountSettings(
		account.Email,
		settings,
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
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

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
}

func (r *Resource) updateAccountSettings(
	ctx context.Context,
	plan *AccountSettingsModel,
) diag.Diagnostics {
	client := r.Meta.Client
	var diagnostics diag.Diagnostics

	// Longview Plan update functionality has been moved
	if !plan.LongviewSubscription.IsNull() {
		_, err := client.UpdateLongviewPlan(ctx, linodego.LongviewPlanUpdateOptions{
			LongviewSubscription: plan.LongviewSubscription.ValueString(),
		})
		if err != nil {
			diagnostics.AddError(
				"Failed to update Linode Longview Plan",
				err.Error(),
			)
			return diagnostics
		}
	}

	updateOpts := linodego.AccountSettingsUpdateOptions{
		BackupsEnabled: plan.BackupsEnabed.ValueBoolPointer(),
		NetworkHelper:  plan.NetworkHelper.ValueBoolPointer(),
	}

	settings, err := client.UpdateAccountSettings(ctx, updateOpts)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("Failed to update Linode Account Settings"),
			err.Error(),
		)
		return diagnostics
	}

	plan.parseAccountSettings(
		plan.ID.ValueString(),
		settings,
	)

	return diagnostics
}
