package accountsettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *linodego.Client
}

func (r *Resource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetResourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	r.client = meta.Client
}

func (r *Resource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "linode_account_settings"
}

func (r *Resource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = frameworkResourceSchema
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
	plan, err := r.updateAccountSettings(ctx, plan)
	resp.Diagnostics.Append(
		err...,
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
	client := r.client

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

	resp.Diagnostics.Append(data.parseAccountSettings(
		ctx,
		account.Email,
		settings,
	)...)
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
	var plan AccountSettingsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the account
	plan, err := r.updateAccountSettings(ctx, plan)
	resp.Diagnostics.Append(
		err...,
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
	plan AccountSettingsModel,
) (AccountSettingsModel, diag.Diagnostics) {
	client := r.client
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
			return plan, diagnostics
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
		return plan, diagnostics
	}

	diagnostics.Append(
		plan.parseAccountSettings(
			ctx,
			plan.ID.ValueString(),
			settings,
		)...,
	)

	return plan, diagnostics
}
