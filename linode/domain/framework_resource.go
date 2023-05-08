package domain

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *linodego.Client
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

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data DomainModel
	client := r.client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var masterIPs, axfrIPs, tags []string

	resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.AXFRIPs.ElementsAs(ctx, &axfrIPs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.MasterIPs.ElementsAs(ctx, &masterIPs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.DomainCreateOptions{
		Domain:      data.Domain.ValueString(),
		Type:        linodego.DomainType(data.Type.ValueString()),
		Group:       data.Group.ValueString(),
		Status:      linodego.DomainStatus(data.Status.ValueString()),
		Description: data.Description.ValueString(),
		TTLSec:      int(data.TTLSec.ValueInt64()),
		RetrySec:    int(data.RetrySec.ValueInt64()),
		ExpireSec:   int(data.ExpireSec.ValueInt64()),
		RefreshSec:  int(data.RefreshSec.ValueInt64()),
		SOAEmail:    data.SOAEmail.ValueString(),
		MasterIPs:   masterIPs,
		AXfrIPs:     axfrIPs,
		Tags:        tags,
	}

	domain, err := client.CreateDomain(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Linode Domain",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseDomain(ctx, domain)...)
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
	client := r.client

	var data DomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.StringToInt64(data.ID.ValueString(), resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := client.GetDomain(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Domain no longer exists.",
				fmt.Sprintf(
					"Removing Linode Domain with ID %v from state because it no longer exists",
					data.ID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the Linode Domain",
			fmt.Sprintf(
				"Error finding the specified Linode Domain: %s",
				err.Error(),
			),
		)
		return
	}

	resp.Diagnostics.Append(data.parseDomain(ctx, domain)...)
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
	client := r.client

	var state DomainModel
	var plan DomainModel

	// Get the state & plan
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check whether there were any changes
	shouldUpdate, err := helper.IsModelUpdated(state, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to check is resource needs to be updated",
			err.Error(),
		)
		return
	}

	if !shouldUpdate {
		return
	}

	// Get the ID from the plan
	domainID := int(helper.StringToInt64(state.ID.ValueString(), resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpts := linodego.DomainUpdateOptions{
		Domain:      plan.Domain.ValueString(),
		Type:        linodego.DomainType(plan.Type.ValueString()),
		Group:       plan.Group.ValueString(),
		Status:      linodego.DomainStatus(plan.Status.ValueString()),
		Description: plan.Description.ValueString(),
		SOAEmail:    plan.SOAEmail.ValueString(),
		RetrySec:    int(plan.RetrySec.ValueInt64()),
		ExpireSec:   int(plan.ExpireSec.ValueInt64()),
		RefreshSec:  int(plan.RefreshSec.ValueInt64()),
		TTLSec:      int(plan.TTLSec.ValueInt64()),
	}

	domain, err := client.UpdateDomain(ctx, domainID, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update domain",
			err.Error(),
		)
	}

	resp.Diagnostics.Append(plan.parseDomain(ctx, domain)...)
	if resp.Diagnostics.HasError() {
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

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID := int(helper.StringToInt64(data.ID.ValueString(), resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client
	err := client.DeleteDomain(ctx, domainID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete the Domain with id %v", domainID),
			err.Error(),
		)
		return
	}
}
