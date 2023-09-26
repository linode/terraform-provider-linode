package nb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

var _ resource.ResourceWithUpgradeState = &Resource{}

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_nodebalancer",
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
	var data NodebalancerModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientConnThrottle := int(data.ClientConnThrottle.ValueInt64())
	createOpts := linodego.NodeBalancerCreateOptions{
		Region:             data.Region.ValueString(),
		Label:              data.Label.ValueStringPointer(),
		ClientConnThrottle: &clientConnThrottle,
	}

	if !data.FirewallID.IsNull() {
		createOpts.FirewallID = int(data.FirewallID.ValueInt64())
	}

	if !data.Tags.IsNull() {
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &createOpts.Tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	nodebalancer, err := client.CreateNodeBalancer(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating a Linode NodeBalancer",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseComputedAttrs(ctx, nodebalancer)...)
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
	var data NodebalancerModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodebalancer, err := client.GetNodeBalancer(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Nodebalancer",
				fmt.Sprintf("removing Linode NodeBalancer ID %v from state because it no longer exists", data.ID.ValueInt64()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get nodebalancer %v", data.ID.ValueInt64()),
			err.Error(),
		)
	}

	resp.Diagnostics.Append(data.ParseComputedAttrs(ctx, nodebalancer)...)
	resp.Diagnostics.Append(data.ParseNonComputedAttrs(ctx, nodebalancer)...)
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
	var plan, state NodebalancerModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	isEqual := state.Label.Equal(plan.Label) &&
		state.ClientConnThrottle.Equal(plan.ClientConnThrottle) &&
		state.Tags.Equal(plan.Tags)

	if !isEqual {
		clientConnThrottle := int(plan.ClientConnThrottle.ValueInt64())
		updateOpts := linodego.NodeBalancerUpdateOptions{
			Label:              plan.Label.ValueStringPointer(),
			ClientConnThrottle: &clientConnThrottle,
		}
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &updateOpts.Tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		nodebalancer, err := client.UpdateNodeBalancer(ctx, int(plan.ID.ValueInt64()), updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update Nodebalancer %v", plan.ID.ValueInt64()),
				err.Error(),
			)
			return
		}

		resp.Diagnostics.Append(plan.ParseComputedAttrs(ctx, nodebalancer)...)
	} else {
		req.State.GetAttribute(ctx, path.Root("updated"), &plan.Updated)
		req.State.GetAttribute(ctx, path.Root("transfer"), &plan.Transfer)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data NodebalancerModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := client.DeleteNodeBalancer(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Nodebalancer does not exist.",
				fmt.Sprintf("Nodebalancer %v does not exist, removing from state.", data.ID.ValueInt64()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to delete Nodebalancer",
			err.Error(),
		)
		return
	}
}

func (r *Resource) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &resourceNodebalancerV0,
			StateUpgrader: upgradeNodebalancerResourceStateV0toV1,
		},
	}
}

func upgradeNodebalancerResourceStateV0toV1(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	var nbDataV0 nbModelV0

	resp.Diagnostics.Append(req.State.Get(ctx, &nbDataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nbDataV1 := NodebalancerModel{
		ID:                 nbDataV0.ID,
		Label:              nbDataV0.Label,
		Region:             nbDataV0.Region,
		ClientConnThrottle: nbDataV0.ClientConnThrottle,
		Hostname:           nbDataV0.Hostname,
		Ipv4:               nbDataV0.Ipv4,
		Ipv6:               nbDataV0.Ipv6,
		Created:            customtypes.RFC3339TimeStringValue{StringValue: nbDataV0.Created},
		Updated:            customtypes.RFC3339TimeStringValue{StringValue: nbDataV0.Updated},
		Tags:               nbDataV0.Tags,
	}

	var transferMap map[string]string
	resp.Diagnostics.Append(nbDataV0.Transfer.ElementsAs(ctx, &transferMap, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result := make(map[string]attr.Value)
	in, diag := UpgradeResourceStateValue(transferMap["in"])
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	result["in"] = in

	out, diag := UpgradeResourceStateValue(transferMap["out"])
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	result["out"] = out

	total, diag := UpgradeResourceStateValue(transferMap["total"])
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	result["total"] = total

	transferObj, diags := types.ObjectValue(TransferObjectType.AttrTypes, result)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resultList, diags := types.ListValueFrom(ctx, TransferObjectType, []attr.Value{transferObj})

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	nbDataV1.Transfer = resultList

	resp.Diagnostics.Append(resp.State.Set(ctx, &nbDataV1)...)
}
