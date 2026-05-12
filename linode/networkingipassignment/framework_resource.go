package networkingipassignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/instancenetworking"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_networking_ip_assignment",
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
	tflog.Debug(ctx, "Create "+r.Config.Name)
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

	// Computed fields (reserved, tags, assigned_entity) are left null here and will be
	// populated by Read on the next refresh.
	for i := range plan.Assignments {
		plan.Assignments[i].Reserved = types.BoolNull()
		plan.Assignments[i].Tags = types.SetNull(types.StringType)
		plan.Assignments[i].AssignedEntity = types.ObjectNull(instancenetworking.AssignedEntityObjectType.AttrTypes)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read "+r.Config.Name)
	var state NetworkingIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	retained := make([]AssignmentModel, 0, len(state.Assignments))
	for _, assignment := range state.Assignments {
		ip, err := client.GetIPAddress(ctx, assignment.Address.ValueString())
		if err != nil {
			if linodego.IsNotFound(err) {
				// IP not found; drop it from state.
				continue
			}
			resp.Diagnostics.AddError(
				"Error reading IP Address",
				fmt.Sprintf("Could not read IP address %s: %s", assignment.Address.ValueString(), err),
			)
			return
		}

		retained = append(retained, AssignmentModel{
			Address:        types.StringValue(ip.Address),
			LinodeID:       types.Int64Value(int64(ip.LinodeID)),
			Reserved:       types.BoolValue(ip.Reserved),
			Tags:           flattenTagsList(ip.Tags, &resp.Diagnostics),
			AssignedEntity: flattenAssignedEntity(ip.AssignedEntity, &resp.Diagnostics),
		})
	}
	state.Assignments = retained

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)
	var state NetworkingIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	// Iterate through all assignments and unassign each IP
	for _, assignment := range state.Assignments {
		linodeID := helper.FrameworkSafeInt64ToInt(assignment.LinodeID.ValueInt64(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		ipAddress := assignment.Address.ValueString()
		if ipAddress == "" {
			resp.Diagnostics.AddWarning(
				"Invalid IP Address",
				"An IP address in the assignment list is empty or invalid. Skipping.",
			)
			continue
		}

		// Unassign the IP address from the Linode
		err := client.DeleteInstanceIPAddress(ctx, linodeID, ipAddress)
		if err != nil {
			// Log and continue with the next IP address if an error occurs
			resp.Diagnostics.AddError(
				"Failed to Unassign IP Address",
				fmt.Sprintf("Error unassigning IP address %s from Linode %d: %s", ipAddress, linodeID, err),
			)
			continue
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully unassigned IP %s from Linode %d", ipAddress, linodeID))
	}

	resp.State.RemoveResource(ctx)
}

func flattenTagsList(tags []string, diags *diag.Diagnostics) types.Set {
	elems := make([]attr.Value, len(tags))
	for i, t := range tags {
		elems[i] = types.StringValue(t)
	}
	set, d := types.SetValue(types.StringType, elems)
	diags.Append(d...)
	return set
}

func flattenAssignedEntity(entity *linodego.ReservedIPAssignedEntity, diags *diag.Diagnostics) types.Object {
	result, d := instancenetworking.FlattenAssignedEntity(entity)
	diags.Append(d...)
	return result
}
