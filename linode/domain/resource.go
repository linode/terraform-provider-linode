package domain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type Resource struct {
	client *linodego.Client
}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetResourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	r.client = meta.Client
}

func (r *Resource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "linode_domain"
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client := r.client

	var data DomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueInt64()
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := client.GetDomain(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Domain No Longer Exists",
				fmt.Sprintf(
					"Removing Domain with ID %v from state because it no longer exists",
					data.ID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error finding the specified Domain",
			err.Error(),
		)
		return
	}

	data.parseDomain(domain)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client := r.client
	var data DomainModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.DomainCreateOptions{
		Domain:      data.Domain.ValueString(),
		Type:        linodego.DomainType(data.Type.ValueString()),
		Group:       data.Group.ValueString(),
		Description: data.Description.ValueString(),
		SOAEmail:    data.SOAEmail.ValueString(),
		RetrySec:    int(data.RetrySec.ValueInt64()),
		ExpireSec:   int(data.ExpireSec.ValueInt64()),
		RefreshSec:  int(data.RefreshSec.ValueInt64()),
		TTLSec:      int(data.TTLSec.ValueInt64()),
		MasterIPs:   helper.FrameworkToStringSlice(data.MasterIPs),
		AXfrIPs:     helper.FrameworkToStringSlice(data.AXFRIPs),
		Tags:        helper.FrameworkToStringSlice(data.Tags),
	}

	domain, err := client.CreateDomain(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Domain creation error",
			err.Error(),
		)
		return
	}

	data.parseDomain(domain)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state DomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var masterIPs []string
	if !state.MasterIPs.Equal(plan.MasterIPs) {
		masterIPs = helper.FrameworkToStringSlice(plan.MasterIPs)
	}
	var axfrIPs []string
	if !state.AXFRIPs.Equal(plan.AXFRIPs) {
		axfrIPs = helper.FrameworkToStringSlice(plan.AXFRIPs)
	}
	var tags []string
	if !state.Tags.Equal(plan.Tags) {
		tags = helper.FrameworkToStringSlice(plan.Tags)
	}

	ops := linodego.DomainUpdateOptions{
		Domain:      plan.Domain.ValueString(),
		Type:        linodego.DomainType(plan.Type.ValueString()),
		Group:       plan.Group.ValueString(),
		Description: plan.Description.ValueString(),
		SOAEmail:    plan.SOAEmail.ValueString(),
		RetrySec:    int(plan.RetrySec.ValueInt64()),
		ExpireSec:   int(plan.ExpireSec.ValueInt64()),
		RefreshSec:  int(plan.RefreshSec.ValueInt64()),
		TTLSec:      int(plan.TTLSec.ValueInt64()),
		MasterIPs:   masterIPs,
		AXfrIPs:     axfrIPs,
		Tags:        tags,
	}
	id := plan.ID.ValueInt64()
	client := r.client

	_, err := client.UpdateDomain(ctx, int(id), ops)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update domain: %v", id),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data DomainModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := data.ID.ValueInt64()
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client
	err := client.DeleteDomain(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); (ok && lerr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete domain with id %v", id),
				err.Error(),
			)
		}
		return
	}
}
