package natgateway

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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

var _ resource.ResourceWithValidateConfig = (*Resource)(nil)

// ValidateConfig rejects configurations where `addresses` is declared for an
// autoscaling NAT Gateway. Autoscaling gateways allocate their own addresses,
// which Read stores in `addresses`. A declared set (including an empty one)
// would then permanently conflict with that state, and Update would try to
// delete addresses the API refuses to remove.
func (r *Resource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config NATGatewayModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// An omitted `addresses` is the only valid configuration for an
	// autoscaling gateway, so there is nothing left to check.
	if config.Addresses.IsNull() || config.Addresses.IsUnknown() {
		return
	}

	if config.UseAutoscaling.IsUnknown() {
		return
	}

	// `use_autoscaling` defaults to true when `addresses` is empty.
	autoscaling := config.UseAutoscaling.ValueBool() ||
		(config.UseAutoscaling.IsNull() && len(config.Addresses.Elements()) == 0)
	if !autoscaling {
		return
	}

	detail := "An autoscaling NAT Gateway allocates its own addresses, which are " +
		"reported in the computed `addresses` attribute. Remove the `addresses` " +
		"attribute from this resource."

	if config.UseAutoscaling.IsNull() {
		detail += " To manage addresses explicitly, set `use_autoscaling = false`."
	}

	resp.Diagnostics.AddAttributeError(
		path.Root("addresses"),
		"`addresses` cannot be configured when autoscaling is enabled.",
		detail,
	)
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

	resp.Diagnostics.Append(data.FlattenNATGateway(natGateway, true)...)
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

	resp.Diagnostics.Append(data.FlattenNATGateway(natGateway, false)...)
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

	// Create-only fields (region, default_ports_per_interface, use_autoscaling)
	// force replacement when they change, so they never reach
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
		addressDiags := r.updateAddresses(ctx, id, plan.Addresses, state.Addresses)
		resp.Diagnostics.Append(addressDiags...)
		if addressDiags.HasError() {
			// Addresses may have been partially modified, so persist the
			// gateway's actual configuration rather than leaving stale state.
			if natGateway, err := client.GetNATGateway(ctx, id); err == nil {
				resp.Diagnostics.Append(state.FlattenNATGateway(natGateway, false)...)
				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			}
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

	resp.Diagnostics.Append(plan.FlattenNATGateway(natGateway, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// addressOp is a single pending address assignment or unassignment.
type addressOp struct {
	address string
	add     bool
}

// updateAddresses reconciles the NAT Gateway's addresses with the plan.
//
// Additions are ordered ahead of removals whenever the gateway has headroom
// under maximum_reserved_addresses_per_natgateway, so a failed addition never
// leaves the gateway without its working egress addresses. When the gateway is
// already at the limit, a removal is interleaved first to free a slot. Any
// operation that fails causes the operations already applied to be undone, so
// a failed replacement preserves the working configuration.
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

	toAdd := addressesMissingFrom(planAddrs, stateAddrs)
	toRemove := addressesMissingFrom(stateAddrs, planAddrs)
	if len(toAdd) == 0 && len(toRemove) == 0 {
		return diags
	}

	headroom, d := r.addressHeadroom(ctx, len(stateAddrs))
	diags.Append(d...)

	ops := planAddressOps(toAdd, toRemove, headroom)

	applied := make([]addressOp, 0, len(ops))
	for _, op := range ops {
		if err := r.applyAddressOp(ctx, id, op); err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to update the addresses of NAT Gateway (%d).", id),
				err.Error(),
			)
			diags.Append(r.rollbackAddressOps(ctx, id, applied)...)
			return diags
		}
		applied = append(applied, op)
	}

	return diags
}

// addressHeadroom returns how many additional addresses may be assigned to a
// NAT Gateway that currently has assigned addresses. If the settings cannot be
// read, zero is returned so the gateway is treated as being at its limit.
func (r *Resource) addressHeadroom(ctx context.Context, assigned int) (int, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Debug(ctx, "client.NATGatewayGetSettings(...)")
	settings, err := r.Meta.Client.NATGatewayGetSettings(ctx)
	if err != nil {
		diags.AddWarning(
			"Failed to read NAT Gateway settings.",
			"Address changes will be applied one slot at a time because the maximum number "+
				"of reserved addresses per NAT Gateway could not be determined: "+err.Error(),
		)
		return 0, diags
	}

	return max(settings.MaximumReservedAddressesPerNATGateway-assigned, 0), diags
}

// planAddressOps orders the pending address changes so the limit on reserved
// addresses is never exceeded while keeping as many working addresses assigned
// as possible. Additions are applied while there is headroom; otherwise a
// removal is applied first to free a slot.
func planAddressOps(toAdd, toRemove []string, headroom int) []addressOp {
	ops := make([]addressOp, 0, len(toAdd)+len(toRemove))

	for len(toAdd) > 0 || len(toRemove) > 0 {
		switch {
		case len(toAdd) > 0 && headroom > 0:
			ops = append(ops, addressOp{address: toAdd[0], add: true})
			toAdd = toAdd[1:]
			headroom--
		case len(toRemove) > 0:
			ops = append(ops, addressOp{address: toRemove[0]})
			toRemove = toRemove[1:]
			headroom++
		default:
			// No headroom and nothing left to remove. Attempt the addition
			// anyway so the API reports the limit to the user.
			ops = append(ops, addressOp{address: toAdd[0], add: true})
			toAdd = toAdd[1:]
		}
	}

	return ops
}

// applyAddressOp assigns or unassigns a single address.
func (r *Resource) applyAddressOp(ctx context.Context, id int, op addressOp) error {
	if op.add {
		tflog.Debug(ctx, "client.NATGatewayAddAddress(...)", map[string]any{
			"address": op.address,
		})
		if _, err := r.Meta.Client.NATGatewayAddAddress(ctx, id, linodego.NATGatewayAddAddressOptions{
			Address: op.address,
		}); err != nil {
			return fmt.Errorf("failed to add address %s: %w", op.address, err)
		}
		return nil
	}

	tflog.Debug(ctx, "client.NATGatewayDeleteAddress(...)", map[string]any{
		"address": op.address,
	})
	if err := r.Meta.Client.NATGatewayDeleteAddress(ctx, id, op.address); err != nil {
		return fmt.Errorf("failed to remove address %s: %w", op.address, err)
	}

	return nil
}

// rollbackAddressOps undoes the applied operations in reverse order on a
// best-effort basis. Failures are reported as warnings because the original
// error is the one the user needs to act on.
func (r *Resource) rollbackAddressOps(ctx context.Context, id int, applied []addressOp) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, op := range slices.Backward(applied) {
		op.add = !op.add

		if err := r.applyAddressOp(ctx, id, op); err != nil {
			diags.AddWarning(
				"Failed to roll back a NAT Gateway address change.",
				fmt.Sprintf(
					"NAT Gateway (%d) may be left with addresses that do not match the "+
						"previous configuration and should be reconciled manually: %s",
					id, err.Error(),
				),
			)
			return diags
		}
	}

	return diags
}

// addressesMissingFrom returns the addresses in a that are not in b.
func addressesMissingFrom(a, b []string) []string {
	missing := make([]string, 0, len(a))
	for _, addr := range a {
		if !slices.Contains(b, addr) {
			missing = append(missing, addr)
		}
	}
	return missing
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
