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

	if len(plan.Assignments) > 0 {
		tflog.Info(ctx, "Assigning IP addresses", map[string]interface{}{
			"assignments": plan.Assignments,
		})

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

		plan.ID = types.StringValue(plan.Assignments[0].Address.ValueString())
		plan.Region = types.StringValue(assignOpts.Region)
	}

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

	// Initialize a slice to store all assignments
	allAssignments := []AssignmentModel{}

	// Iterate through all assignments in the current state
	for _, assignment := range state.Assignments {
		ip, err := client.GetIPAddress(ctx, assignment.Address.ValueString())
		if err != nil {
			if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
				// If the IP is not found, skip it
				continue
			}
			resp.Diagnostics.AddError(
				"Error reading IP Address",
				fmt.Sprintf("Could not read IP address %s: %s", assignment.Address.ValueString(), err),
			)
			return
		}

		// Add the assignment to our slice
		allAssignments = append(allAssignments, AssignmentModel{
			Address:  types.StringValue(ip.Address),
			LinodeID: types.Int64Value(int64(ip.LinodeID)),
		})

		// If this is the first assignment, use it to set the main resource attributes
		if len(allAssignments) == 1 {
			state.ID = types.StringValue(ip.Address)
			state.Address = types.StringValue(ip.Address)
			state.LinodeID = types.Int64Value(int64(ip.LinodeID))
			state.Region = types.StringValue(ip.Region)
			state.Public = types.BoolValue(ip.Public)
			state.Type = types.StringValue(string(ip.Type))
			state.Reserved = types.BoolValue(ip.Reserved)
		}
	}

	// Update the state with all assignments
	state.Assignments = allAssignments

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

	if len(plan.Assignments) > 0 {
		// Handle IP assignment updates
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

		// Update plan with new assignment details
		plan.ID = types.StringValue(plan.Assignments[0].Address.ValueString())
		plan.Address = plan.Assignments[0].Address
	}

	// Re-read the state to get the latest information
	readResp := resource.ReadResponse{State: resp.State}
	if readResp.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResp.Diagnostics...)
		return
	}
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete linode_networking_assign_ip")
	var state NetworkingIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	// Check if this is an assigned IP
	if len(state.Assignments) > 0 {
		// For assigned IPs, we need to unassign them
		for _, assignment := range state.Assignments {
			linodeID := int(assignment.LinodeID.ValueInt64())
			ipAddress := assignment.Address.ValueString()

			err := client.DeleteInstanceIPAddress(ctx, linodeID, ipAddress)
			if err != nil {
				if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
					resp.Diagnostics.AddError(
						"Failed to Unassign IP",
						fmt.Sprintf(
							"failed to unassign ip (%s) from instance (%d): %s",
							ipAddress, linodeID, err.Error(),
						),
					)
				}
			}
		}
	}
}
