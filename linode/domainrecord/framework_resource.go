package domainrecord

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
				Name:   "linode_domain_record",
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
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var plan ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = populateLogAttributes(ctx, plan)

	priority := helper.FrameworkSafeInt64PointerToIntPointer(plan.Priority.ValueInt64Pointer(), &resp.Diagnostics)
	weight := helper.FrameworkSafeInt64PointerToIntPointer(plan.Weight.ValueInt64Pointer(), &resp.Diagnostics)
	port := helper.FrameworkSafeInt64PointerToIntPointer(plan.Port.ValueInt64Pointer(), &resp.Diagnostics)
	ttlSec := helper.FrameworkSafeInt64ToInt(plan.TTLSec.ValueInt64(), &resp.Diagnostics)
	domainID := helper.FrameworkSafeInt64ToInt(plan.DomainID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.DomainRecordCreateOptions{
		Type:     linodego.DomainRecordType(plan.RecordType.ValueString()),
		Name:     plan.Name.ValueString(),
		Target:   plan.Target.ValueString(),
		Priority: priority,
		Weight:   weight,
		Port:     port,
		Service:  plan.Service.ValueStringPointer(),
		Protocol: plan.Protocol.ValueStringPointer(),
		TTLSec:   ttlSec,
		Tag:      plan.Tag.ValueStringPointer(),
	}

	client := r.Meta.Client

	tflog.Debug(ctx, "client.CreateDomainRecord(...)", map[string]interface{}{
		"options": createOpts,
	})

	domainRecord, err := client.CreateDomainRecord(ctx, domainID, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating a Linode DomainRecord: %s", err), err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(plan.FlattenDomainRecord(ctx, client, domainRecord, true)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	domainID, domainRecordID := getDomainAndRecordIDs(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "domain_id", domainID)

	tflog.Trace(ctx, "client.GetDomainRecord(...)")

	record, err := client.GetDomainRecord(ctx, domainID, domainRecordID)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Domain Record does not exist.",
				fmt.Sprintf(
					"Removing Domain Record with ID %v from state because it no longer exists",
					domainRecordID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the Domain Record.",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(state.FlattenDomainRecord(ctx, client, record, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)
	var plan, state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	domainID, domainRecordID := getDomainAndRecordIDs(plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	priority := helper.FrameworkSafeInt64PointerToIntPointer(plan.Priority.ValueInt64Pointer(), &resp.Diagnostics)
	weight := helper.FrameworkSafeInt64PointerToIntPointer(plan.Weight.ValueInt64Pointer(), &resp.Diagnostics)
	port := helper.FrameworkSafeInt64PointerToIntPointer(plan.Port.ValueInt64Pointer(), &resp.Diagnostics)
	ttlSec := helper.FrameworkSafeInt64ToInt(plan.TTLSec.ValueInt64(), &resp.Diagnostics)

	updateOpts := linodego.DomainRecordUpdateOptions{
		Type:     linodego.DomainRecordType(plan.RecordType.ValueString()),
		Name:     plan.Name.ValueString(),
		Target:   plan.Target.ValueString(),
		Priority: priority,
		Weight:   weight,
		Port:     port,
		Service:  plan.Service.ValueStringPointer(),
		Protocol: plan.Protocol.ValueStringPointer(),
		TTLSec:   ttlSec,
		Tag:      plan.Tag.ValueStringPointer(),
	}

	tflog.Debug(ctx, "client.UpdateDomainRecord(...)", map[string]interface{}{
		"options": updateOpts,
	})
	domainRecord, err := client.UpdateDomainRecord(ctx, domainID, domainRecordID, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Domain Record", err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(plan.FlattenDomainRecord(ctx, client, domainRecord, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.CopyFrom(state, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	var state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	domainID, domainRecordID := getDomainAndRecordIDs(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteDomainRecord(...)")
	err := client.DeleteDomainRecord(ctx, domainID, domainRecordID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting Linode DomainRecord %d", domainRecordID),
			err.Error(),
		)
		return
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import "+r.Config.Name)
	helper.ImportStateWithMultipleIDs(ctx, req, resp,
		[]helper.ImportableID{
			{
				Name:          "domain_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "id",
				TypeConverter: helper.IDTypeConverterString,
			},
		})
}

func populateLogAttributes(ctx context.Context, data ResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"domain_record_id": data.ID.ValueString(),
		"domain_id":        data.DomainID.ValueInt64(),
	})
}

func getDomainAndRecordIDs(data ResourceModel, diags *diag.Diagnostics) (int, int) {
	domainID := helper.FrameworkSafeInt64ToInt(
		data.DomainID.ValueInt64(), diags,
	)
	domainRecordID := helper.FrameworkSafeStringToInt(
		data.ID.ValueString(), diags,
	)
	return domainID, domainRecordID
}
