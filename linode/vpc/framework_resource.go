package vpc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_vpc",
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
	var data VPCResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnetCreateOptsList := make([]linodego.VPCSubnetCreateOptions, len(data.SubnetsCreateOptions))

	for i, subnetOpts := range data.SubnetsCreateOptions {
		subnetCreateOpts := linodego.VPCSubnetCreateOptions{
			Label: subnetOpts.Label.ValueString(),
			IPv4:  subnetOpts.IPv4.ValueString(),
		}
		subnetCreateOptsList[i] = subnetCreateOpts
	}

	vpcCreateOpts := linodego.VPCCreateOptions{
		Label:       data.Label.ValueString(),
		Region:      data.Region.ValueString(),
		Description: data.Description.ValueString(),
		Subnets:     subnetCreateOptsList,
	}

	vpc, err := client.CreateVPC(ctx, vpcCreateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create VPC.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseComputedAttributes(ctx, vpc)...)
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
	var data VPCResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vpc, err := client.GetVPC(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"VPC no longer exists.",
				fmt.Sprintf(
					"Removing Linode VPC with ID %v from state because it no longer exists.",
					data.ID.ValueInt64(),
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the VPC.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseVPC(ctx, vpc)...)
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
	client := r.Meta.Client
	var plan, state VPCResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var updateOpts linodego.VPCUpdateOptions
	shouldUpdate := false

	if !state.Description.Equal(plan.Description) {
		shouldUpdate = true
		updateOpts.Description = plan.Description.ValueString()
	}

	if !state.Label.Equal(plan.Label) {
		shouldUpdate = true
		updateOpts.Label = plan.Label.ValueString()
	}

	if shouldUpdate {
		vpc, err := client.UpdateVPC(ctx, int(plan.ID.ValueInt64()), updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update VPC (%d).", plan.ID.ValueInt64()),
				err.Error(),
			)
			return
		}
		resp.Diagnostics.Append(plan.parseComputedAttributes(ctx, vpc)...)
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
	client := r.Meta.Client
	var data VPCResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := client.DeleteVPC(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); (ok && lerr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete the VPC (%d)", data.ID.ValueInt64()),
				err.Error(),
			)
		}
		return
	}
}
