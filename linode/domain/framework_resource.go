package domain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			"linode_domain",
			frameworkResourceSchema,
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client := r.Meta.Client

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

	data.parseDomain(ctx, domain)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client := r.Meta.Client
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
	}

	if !data.MasterIPs.IsNull() && !data.MasterIPs.IsUnknown() {
		resp.Diagnostics.Append(data.MasterIPs.ElementsAs(ctx, &createOpts.MasterIPs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !data.AXFRIPs.IsNull() && !data.AXFRIPs.IsUnknown() {
		resp.Diagnostics.Append(data.AXFRIPs.ElementsAs(ctx, &createOpts.AXfrIPs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &createOpts.Tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	domain, err := client.CreateDomain(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Domain creation error",
			err.Error(),
		)
		return
	}

	data.parseDomain(ctx, domain)
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
	if resp.Diagnostics.HasError() {
		return
	}

	domainID := int(state.ID.ValueInt64())

	if !domainDeepEqual(plan, state) {
		r.updateDomain(ctx, resp, plan, domainID)
	}
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

	client := r.Meta.Client
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

func (r *Resource) updateDomain(
	ctx context.Context,
	resp *resource.UpdateResponse,
	plan DomainModel,
	domainID int,
) {
	client := r.Meta.Client
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
	if !plan.MasterIPs.IsNull() && !plan.MasterIPs.IsUnknown() {
		resp.Diagnostics.Append(plan.MasterIPs.ElementsAs(ctx, &updateOpts.MasterIPs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.AXFRIPs.IsNull() && !plan.AXFRIPs.IsUnknown() {
		resp.Diagnostics.Append(plan.AXFRIPs.ElementsAs(ctx, &updateOpts.AXfrIPs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &updateOpts.Tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	domain, err := client.UpdateDomain(ctx, domainID, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update domain: %v", domainID),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(plan.parseDomain(ctx, domain)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
