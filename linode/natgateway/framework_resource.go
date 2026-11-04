package natgateway

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_natgateway",
				IDType: types.Int64Type,
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

	var data NATGatewayModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts, diags := data.ExpandCreateOptions(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.CreateNATGateway(...)", map[string]any{
		"options": createOpts,
	})
	natGateway, err := client.CreateNATGateway(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create NAT Gateway.",
			err.Error(),
		)
		return
	}

	ctx = tflog.SetField(ctx, "id", natGateway.ID)

	resp.Diagnostics.Append(data.ParseNATGateway(natGateway)...)
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

	var data NATGatewayModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "id", data.ID.ValueInt64())

	if data.ID.IsNull() || data.ID.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Removing Resource From State",
			"This resource is being implicitly removed from the Terraform state because "+
				"its ID attribute is empty.",
		)
		resp.State.RemoveResource(ctx)
		return
	}

	id := int(data.ID.ValueInt64())

	natGateway, err := client.GetNATGateway(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"NAT Gateway no longer exists.",
				fmt.Sprintf(
					"Removing Linode NAT Gateway with ID %d from state because it no longer exists.",
					id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the NAT Gateway.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseNATGateway(natGateway)...)
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

	client := r.Meta.Client
	var plan, state NATGatewayModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(state.ID.ValueInt64())
	ctx = tflog.SetField(ctx, "id", id)

	// Create-only fields (region, default_ports_per_interface, use_autoscaling,
	// vpc_subnet_id) force replacement when they change, so they never reach
	// Update with a real change. They are left untouched here.

	if !plan.Label.Equal(state.Label) {
		label := plan.Label.ValueString()
		tflog.Debug(ctx, "client.UpdateNATGateway(...)", map[string]any{
			"label": label,
		})
		if _, err := client.UpdateNATGateway(ctx, id, linodego.NATGatewayUpdateOptions{
			Label: &label,
		}); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update NAT Gateway (%d) label.", id),
				err.Error(),
			)
			return
		}
	}

	// An unknown planned addresses set (not configured) is not treated as a change.
	// Whether addresses may be modified (use_autoscaling) is enforced by the API.
	if !plan.Addresses.IsUnknown() && !plan.Addresses.Equal(state.Addresses) {
		resp.Diagnostics.Append(r.updateAddresses(ctx, id, plan.Addresses, state.Addresses)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	natGateway, err := client.GetNATGateway(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to refresh NAT Gateway (%d) after update.", id),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(plan.ParseNATGateway(natGateway)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// updateAddresses removes addresses present in state but not in plan, then
// adds addresses present in plan but not in state. Removals run first so
// addresses are freed before new ones are assigned.
func (r *Resource) updateAddresses(
	ctx context.Context,
	id int,
	plan, state types.Set,
) diag.Diagnostics {
	var diags diag.Diagnostics

	planAddrs, d := natGatewayAddressStrings(ctx, plan)
	diags.Append(d...)
	stateAddrs, d := natGatewayAddressStrings(ctx, state)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	client := r.Meta.Client

	for _, addr := range stateAddrs {
		if slices.Contains(planAddrs, addr) {
			continue
		}
		tflog.Debug(ctx, "client.NATGatewayDeleteAddress(...)", map[string]any{
			"address": addr,
		})
		if err := client.NATGatewayDeleteAddress(ctx, id, addr); err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to remove address %s from NAT Gateway (%d).", addr, id),
				err.Error(),
			)
			return diags
		}
	}

	for _, addr := range planAddrs {
		if slices.Contains(stateAddrs, addr) {
			continue
		}
		tflog.Debug(ctx, "client.NATGatewayAddAddress(...)", map[string]any{
			"address": addr,
		})
		if _, err := client.NATGatewayAddAddress(ctx, id, linodego.NATGatewayAddAddressOptions{
			Address: addr,
		}); err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to add address %s to NAT Gateway (%d).", addr, id),
				err.Error(),
			)
			return diags
		}
	}

	return diags
}

// natGatewayAddressStrings returns the address strings in the set.
// A null or unknown set is treated as empty.
func natGatewayAddressStrings(ctx context.Context, s types.Set) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if s.IsNull() || s.IsUnknown() {
		return nil, diags
	}

	var models []natGatewayAddressModel
	diags.Append(s.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	addrs := make([]string, 0, len(models))
	for _, m := range models {
		addrs = append(addrs, m.Address.ValueString())
	}
	return addrs, diags
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	client := r.Meta.Client
	var data NATGatewayModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.ID.ValueInt64())
	ctx = tflog.SetField(ctx, "id", id)

	// The API rejects deletion while the gateway is attached to a subnet, so
	// clear the attachment this resource created before deleting.
	resp.Diagnostics.Append(r.detachVPCSubnet(ctx, id, data.VPCSubnetID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteNATGateway(...)")
	if err := client.DeleteNATGateway(ctx, id); err != nil {
		if linodego.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete NAT Gateway (%d).", id),
			err.Error(),
		)
	}
}

const (
	detachPollAttempts = 10
	detachPollInterval = 3 * time.Second
)

// detachVPCSubnet refreshes the NAT Gateway's subnet association and, when the
// attached subnet is the one set in vpc_subnet_id, detaches it and waits until
// the API reports the gateway as detached. Attachments made through other
// resources are left alone.
func (r *Resource) detachVPCSubnet(
	ctx context.Context,
	id int,
	configuredSubnetID types.Int64,
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := r.Meta.Client

	natGateway, err := client.GetNATGateway(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			return diags
		}
		diags.AddError(
			fmt.Sprintf("Failed to refresh NAT Gateway (%d) before deletion.", id),
			err.Error(),
		)
		return diags
	}

	subnetID := natGateway.VPCSubnet.ID
	if subnetID == 0 {
		return diags
	}

	if configuredSubnetID.IsNull() || configuredSubnetID.IsUnknown() ||
		int(configuredSubnetID.ValueInt64()) != subnetID {
		return diags
	}

	tflog.Debug(ctx, "client.UpdateVPCSubnet(...) to detach NAT Gateway", map[string]any{
		"subnet_id": subnetID,
	})

	// A nil inner pointer serializes as JSON null, which clears the association.
	var nilID *int
	if _, err := client.UpdateVPCSubnet(
		ctx,
		natGateway.VPCSubnet.VPCID,
		subnetID,
		linodego.VPCSubnetUpdateOptions{
			Label:      natGateway.VPCSubnet.Label,
			NATGateway: &linodego.VPCSubnetUpdateOptionsNATGateway{ID: &nilID},
		},
	); err != nil {
		diags.AddError(
			fmt.Sprintf("Failed to detach NAT Gateway (%d) from subnet (%d).", id, subnetID),
			err.Error(),
		)
		return diags
	}

	for attempt := 1; ; attempt++ {
		current, err := client.GetNATGateway(ctx, id)
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to refresh NAT Gateway (%d) after detaching.", id),
				err.Error(),
			)
			return diags
		}
		if current.VPCSubnet.ID == 0 {
			return diags
		}
		if attempt >= detachPollAttempts {
			diags.AddError(
				fmt.Sprintf("NAT Gateway (%d) is still attached to subnet (%d).", id, subnetID),
				"The NAT Gateway cannot be deleted while attached to a subnet.",
			)
			return diags
		}

		select {
		case <-ctx.Done():
			diags.AddError(
				fmt.Sprintf("Context canceled while waiting to detach NAT Gateway (%d).", id),
				ctx.Err().Error(),
			)
			return diags
		case <-time.After(detachPollInterval):
		}
	}
}
