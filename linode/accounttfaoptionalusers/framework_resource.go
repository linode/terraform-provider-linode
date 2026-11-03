package accounttfaoptionalusers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

var _ resource.ResourceWithImportState = &Resource{}

type Resource struct {
	helper.BaseResource
}

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(helper.BaseResourceConfig{
			Name:   "linode_account_tfa_optional_users",
			IDType: types.StringType,
			Schema: &frameworkResourceSchema,
		}),
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.write(ctx, plan, &resp.State, &resp.Diagnostics)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var data ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	usernames, diags := readUsernames(ctx, r.Meta.Client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(singletonID)
	data.Usernames = usernames
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	var plan, state ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	if plan.Usernames.Equal(state.Usernames) {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	r.write(ctx, plan, &resp.State, &resp.Diagnostics)
}

func (r *Resource) write(ctx context.Context, plan ResourceModel, state *tfsdk.State, diags *diag.Diagnostics) {
	usernames := make([]string, 0, len(plan.Usernames.Elements()))
	diags.Append(plan.Usernames.ElementsAs(ctx, &usernames, false)...)
	if diags.HasError() {
		return
	}

	tflog.Trace(ctx, "client.UpdateTFAOptionalUsers(...)")
	_, err := r.Meta.Client.UpdateTFAOptionalUsers(ctx, linodego.TFAOptionalUsersUpdateOptions{Usernames: usernames})
	if err != nil {
		diags.AddError("Failed to update TFA optional users", err.Error())
		return
	}

	// PUT returns only a count. Keep the successful write tracked even if the subsequent GET fails.
	plan.ID = types.StringValue(singletonID)
	diags.Append(state.Set(ctx, &plan)...)
	if diags.HasError() {
		return
	}

	actual, readDiags := readUsernames(ctx, r.Meta.Client)
	diags.Append(readDiags...)
	if diags.HasError() {
		return
	}

	matchesPlan := actual.Equal(plan.Usernames)
	plan.Usernames = actual
	diags.Append(state.Set(ctx, &plan)...)
	if !matchesPlan {
		diags.AddError("TFA optional users did not match the plan",
			"The API returned a different exemption set after the update. The observed set has been saved in state; refresh and retry.")
	}
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)
	tflog.Info(ctx, "Leaving TFA exemptions unchanged; removing only Terraform state.")
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError("Invalid import identifier", "Use account to import the account-wide TFA exemption set.")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &ResourceModel{
		ID:        types.StringValue(singletonID),
		Usernames: types.SetNull(types.StringType),
	})...)
}
