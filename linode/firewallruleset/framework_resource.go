package firewallruleset

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_firewall_ruleset",
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

	var plan RuleSetResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := plan.GetCreateOptions(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.CreateFirewallRuleSet(...)")
	ruleset, err := client.CreateFirewallRuleSet(ctx, createOpts)
	if ruleset != nil && ruleset.ID != 0 {
		resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(strconv.Itoa(ruleset.ID)))
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to Create Firewall RuleSet", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(ruleset.ID))
	plan.FlattenRuleSet(ctx, *ruleset, &resp.Diagnostics, true)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	client := r.Meta.Client
	var state RuleSetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetFirewallRuleSet(...)", map[string]any{
		"id": id,
	})
	ruleset, err := client.GetFirewallRuleSet(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Removing RuleSet %d from State", id),
				"Removing the Firewall RuleSet from state because it no longer exists",
			)
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Get Firewall RuleSet %d", id),
				err.Error(),
			)
		}
		return
	}

	state.FlattenRuleSet(ctx, *ruleset, &resp.Diagnostics, false)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	client := r.Meta.Client
	var plan, state RuleSetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpts := plan.GetUpdateOptions(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.UpdateFirewallRuleSet(...)", map[string]any{
		"id": id,
	})
	ruleset, err := client.UpdateFirewallRuleSet(ctx, id, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Update Firewall RuleSet %d", id),
			err.Error(),
		)
		return
	}

	plan.ID = state.ID
	plan.FlattenRuleSet(ctx, *ruleset, &resp.Diagnostics, true)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	var state RuleSetResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.DeleteFirewallRuleSet(...)", map[string]any{
		"id": id,
	})
	if err := client.DeleteFirewallRuleSet(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Delete Firewall RuleSet %d", id),
			err.Error(),
		)
		return
	}
}
