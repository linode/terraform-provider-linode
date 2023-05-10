package rdns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *linodego.Client
}

func (data *ResourceModel) parseIP(ip *linodego.InstanceIP) {
	data.Address = types.StringValue(ip.Address)
	data.RDNS = types.StringValue(ip.RDNS)

	id, _ := json.Marshal(ip)

	data.ID = types.StringValue(string(id))
}

func (r *Resource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetResourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	r.client = meta.Client
}

type ResourceModel struct {
	Address          types.String `tfsdk:"address"`
	RDNS             types.String `tfsdk:"rdns"`
	WaitForAvailable types.Bool   `tfsdk:"wait_for_available"`
	ID               types.String `tfsdk:"id"`
}

func (r *Resource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "linode_rdns"
}

func (r *Resource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = frameworkResourceSchema
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("address"), req, resp)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data ResourceModel
	client := r.client

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
	client := r.client

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Println(data.Address.ValueString())

	ip, err := client.GetIPAddress(ctx, data.Address.ValueString())
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

	client := r.client
	ip, err := client.GetIPAddress(ctx, data.Address.ValueString())
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

	client := r.client

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
