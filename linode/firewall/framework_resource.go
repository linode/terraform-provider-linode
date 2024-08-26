package firewall

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_firewall",
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

	var plan FirewallResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := plan.getCreateOptions(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	firewall, err := client.CreateFirewall(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Create Firewall", err.Error())
		return
	}

	addFirewallResource(ctx, resp, strconv.Itoa(firewall.ID))

	if plan.Disabled.ValueBool() {
		firewall = disableFirewall(ctx, firewall.ID, client, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	plan.flattenFirewallForResource(firewall, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshDevices(ctx, client, firewall.ID, &plan, &resp.Diagnostics, true)
	refreshRules(ctx, client, firewall.ID, &plan, &resp.Diagnostics, true)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(strconv.Itoa(firewall.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	client := r.Meta.Client

	var state FirewallResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	firewall, err := client.GetFirewall(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Removing Firewall %d from State", id),
				"Removing the Linode Firewall from state because it no longer exists",
			)
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed to Get Firewall %d", id), err.Error())
		}
		return
	}

	state.flattenFirewallForResource(firewall, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshRules(ctx, client, id, &state, &resp.Diagnostics, false)
	refreshDevices(ctx, client, id, &state, &resp.Diagnostics, false)

	// TODO: cleanup when Crossplane fixes it
	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
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
	var plan, state FirewallResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	ctx = populateLogAttributes(ctx, state)

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpts, shouldUpdate := plan.getUpdateOptions(ctx, state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if shouldUpdate {
		firewall, err := client.UpdateFirewall(ctx, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed to Update Firewall %d", id), err.Error())
		}

		plan.flattenFirewallForResource(firewall, true, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if state.RulesHaveChanges(ctx, plan, &resp.Diagnostics) {
		if resp.Diagnostics.HasError() {
			return
		}

		ruleSet := plan.ExpandFirewallRuleSet(ctx, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		firewallRuleSet, err := client.UpdateFirewallRules(ctx, id, ruleSet)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Update Rules for Firewall %d", id), err.Error(),
			)
		}

		plan.flattenRules(ctx, firewallRuleSet, true, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if state.LinodesOrNodeBalancersHaveChanges(ctx, plan) {
		linodeIDs := helper.ExpandFwInt64Set(plan.Linodes, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		nodeBalancerIDs := helper.ExpandFwInt64Set(plan.NodeBalancers, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		assignments := make([]firewallDeviceAssignment, 0, len(linodeIDs)+len(nodeBalancerIDs))
		for _, entityID := range linodeIDs {
			assignments = append(assignments, firewallDeviceAssignment{
				ID:   entityID,
				Type: linodego.FirewallDeviceLinode,
			})
		}

		for _, entityID := range nodeBalancerIDs {
			assignments = append(assignments, firewallDeviceAssignment{
				ID:   entityID,
				Type: linodego.FirewallDeviceNodeBalancer,
			})
		}

		if err := fwUpdateFirewallDevices(ctx, *client, id, assignments); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Update Devices for Firewall %d", id), err.Error(),
			)
			return
		}
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

	var state FirewallResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := client.DeleteFirewall(ctx, id); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to Delete Firewall %d", id), err.Error())
		return
	}
}

func populateLogAttributes(ctx context.Context, model FirewallResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"token_id": model.ID.ValueString(),
	})
}

func addFirewallResource(
	ctx context.Context, resp *resource.CreateResponse, id string,
) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))
}
