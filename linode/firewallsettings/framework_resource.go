package firewallsettings

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
				Name:   "linode_firewall_settings",
				Schema: &FrameworkResourceSchema,
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
	var plan FirewallSettingsResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateFirewallSettings(ctx, client, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate UUID v7 for the resource
	id, err := uuid.NewV7()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to generate UUID",
			fmt.Sprintf("An error occurred while generating a UUID for firewall settings resource: %s", err.Error()),
		)
		return
	}
	plan.ID = types.StringValue(id.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	client := r.Meta.Client

	var state FirewallSettingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetFirewallSettings(...)")
	firewallSettings, err := client.GetFirewallSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall settings",
			fmt.Sprintf(
				"An error occurred while retrieving the firewall settings: %s",
				err.Error(),
			),
		)
		return
	}

	state.FlattenFirewallSettings(ctx, *firewallSettings, false, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	client := r.Meta.Client
	var plan FirewallSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateFirewallSettings(ctx, client, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func updateFirewallSettings(
	ctx context.Context,
	client *linodego.Client,
	plan *FirewallSettingsResourceModel,
	diags *diag.Diagnostics,
) {
	tflog.Debug(ctx, "Updating firewall settings")

	updateOptions := plan.GetUpdateOptions(ctx, diags)
	if diags.HasError() {
		return
	}

	tflog.Debug(ctx, "client.UpdateFirewallSettings(...)", map[string]any{
		"options": updateOptions,
	})
	firewallSettings, err := client.UpdateFirewallSettings(ctx, updateOptions)
	if err != nil {
		diags.AddError(
			"Failed to update firewall settings",
			fmt.Sprintf("An error occurred while updating the firewall settings: %s", err.Error()),
		)
		return
	}

	plan.FlattenFirewallSettings(ctx, *firewallSettings, true, diags)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)
	tflog.Info(
		ctx, "Firewall settings cannot be deleted. "+
			"The TF state has been deleted, but the "+
			"firewall settings will remain in Linode's system.",
	)
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resp.Diagnostics.AddError(
		"Resource Import Not Supported",
		"Importing state for this resource is not supported. "+
			"Please use the resource to manage firewall settings directly for your account.",
	)
}
