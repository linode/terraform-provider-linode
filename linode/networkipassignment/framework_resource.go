package networkipassignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_networking_assign_ip",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create linode_networking_assign_ip")
	var plan NetworkingIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	apiAssignments := make([]linodego.LinodeIPAssignment, len(plan.Assignments))
	for i, assignment := range plan.Assignments {
		apiAssignments[i] = linodego.LinodeIPAssignment{
			Address:  assignment.Address.ValueString(),
			LinodeID: int(assignment.LinodeID.ValueInt64()),
		}
	}

	assignOpts := linodego.LinodesAssignIPsOptions{
		Region:      plan.Region.ValueString(),
		Assignments: apiAssignments,
	}

	err := client.InstancesAssignIPs(ctx, assignOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error assigning IP Addresses",
			fmt.Sprintf("Could not assign IP addresses: %s", err),
		)
		return
	}

	// Generate a unique ID for this resource
	plan.ID = types.StringValue(fmt.Sprintf("%s-%d", plan.Region.ValueString(), len(plan.Assignments)))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read linode_networking_assign_ip")
	var state NetworkingIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	for i, assignment := range state.Assignments {
		ip, err := client.GetIPAddress(ctx, assignment.Address.ValueString())
		if err != nil {
			if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
				// IP not found, remove it from state
				state.Assignments = append(state.Assignments[:i], state.Assignments[i+1:]...)
				continue
			}
			resp.Diagnostics.AddError(
				"Error reading IP Address",
				fmt.Sprintf("Could not read IP address %s: %s", assignment.Address.ValueString(), err),
			)
			return
		}

		state.Assignments[i] = AssignmentModel{
			Address:  types.StringValue(ip.Address),
			LinodeID: types.Int64Value(int64(ip.LinodeID)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update linode_networking_assign_ip")
	var plan, state NetworkingIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	apiAssignments := make([]linodego.LinodeIPAssignment, len(plan.Assignments))
	for i, assignment := range plan.Assignments {
		apiAssignments[i] = linodego.LinodeIPAssignment{
			Address:  assignment.Address.ValueString(),
			LinodeID: int(assignment.LinodeID.ValueInt64()),
		}
	}

	assignOpts := linodego.LinodesAssignIPsOptions{
		Region:      plan.Region.ValueString(),
		Assignments: apiAssignments,
	}

	err := client.InstancesAssignIPs(ctx, assignOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating IP Assignments",
			fmt.Sprintf("Could not update IP assignments: %s", err),
		)
		return
	}

	// Update the ID to reflect any changes in the number of assignments
	plan.ID = types.StringValue(fmt.Sprintf("%s-%d", plan.Region.ValueString(), len(plan.Assignments)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete linode_networking_assign_ip")
	var state NetworkingIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No need to do anything for delete as the IPs will be automatically unassigned when the Linodes are deleted
}
