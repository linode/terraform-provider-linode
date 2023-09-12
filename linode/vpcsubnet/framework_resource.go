package vpcsubnet

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_vpc_subnet",
				IDType: types.Int64Type,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: vpc_id,id. Got: %q", req.ID),
		)
		return
	}

	vpcID, err := strconv.ParseInt(idParts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to convert VPC ID attribute",
			err.Error(),
		)
		return
	}

	id, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to convert ID attribute",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vpc_id"), vpcID)...)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	println("STATE CREATING")
	var data VPCSubnetModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.VPCSubnetCreateOptions{
		Label: data.Label.ValueString(),
		IPv4:  data.IPv4.ValueString(),
	}

	subnet, err := client.CreateVPCSubnet(ctx, createOpts, int(data.VPCId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create VPC subnet.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseComputedAttributes(ctx, subnet)...)
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
	client := r.Meta.Client
	var data VPCSubnetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnet, err := client.GetVPCSubnet(ctx, int(data.VPCId.ValueInt64()), int(data.ID.ValueInt64()))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"VPC subnet does not exist.",
				fmt.Sprintf(
					"Removing VPC subnet with ID %v from state because it no longer exists",
					data.ID.ValueInt64(),
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the VPC subnet.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseVPCSubnet(ctx, subnet)...)
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
	var plan, state VPCSubnetModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var updateOpts linodego.VPCSubnetUpdateOptions
	shouldUpdate := false

	if !state.Label.Equal(plan.Label) {
		updateOpts.Label = plan.Label.ValueString()
		shouldUpdate = true
	}

	if shouldUpdate {
		subnet, err := client.UpdateVPCSubnet(ctx, int(plan.VPCId.ValueInt64()), int(plan.ID.ValueInt64()), updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update VPC subnet (%d).", plan.ID.ValueInt64()),
				err.Error(),
			)
		}
		resp.Diagnostics.Append(plan.parseComputedAttributes(ctx, subnet)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data VPCSubnetModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := client.DeleteVPCSubnet(ctx, int(data.VPCId.ValueInt64()), int(data.ID.ValueInt64()))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); (ok && lerr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete the VPC subnet (%d)", data.ID.ValueInt64()),
				err.Error(),
			)
		}
		return
	}
}
