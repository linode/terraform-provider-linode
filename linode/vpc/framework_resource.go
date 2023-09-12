package vpc

import (
	"context"

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
	var data VPCModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnetCreateOptsList := make([]linodego.VPCSubnetCreateOptions, len(data.Subnets))

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
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
}
