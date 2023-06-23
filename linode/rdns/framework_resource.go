package rdns

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			"linode_rdns",
			frameworkResourceSchema,
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
	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: data.RDNS.ValueStringPointer(),
	}

	ip, err := client.UpdateIPAddress(ctx, data.Address.ValueString(), updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Linode RDNS",
			err.Error(),
		)
		return
	}

	data.parseIP(ip)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client := r.Meta.Client

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ip, err := client.GetIPAddress(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read the Linode RDNS", err.Error(),
		)
		return
	}

	data.parseIP(ip)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	ip, err := client.GetIPAddress(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get IP Address: %s", err.Error(),
		)
		return
	}

	updateOpts := ip.GetUpdateOptions()
	plannedRDNS := data.RDNS.ValueString()

	resourceUpdated := false

	if *updateOpts.RDNS != plannedRDNS {
		updateOpts.RDNS = &plannedRDNS
		resourceUpdated = true
	}

	if resourceUpdated {
		ip, err = client.UpdateIPAddress(ctx, data.Address.ValueString(), updateOpts)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to update the Linode RDNS",
				err.Error(),
			)
			return
		}

		data.parseIP(ip)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: nil,
	}

	ip, err := client.UpdateIPAddress(ctx, data.Address.ValueString(), updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete Linode RDNS",
			err.Error(),
		)
		return
	}

	data.parseIP(ip)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
