package nb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var _ resource.ResourceWithUpgradeState = &Resource{}

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_nodebalancer",
				IDType: types.StringType,
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
	tflog.Debug(ctx, "Create linode_nodebalancer")
	var data NodeBalancerModel
	client := r.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientConnThrottle := helper.FrameworkSafeInt64ToInt(
		data.ClientConnThrottle.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.NodeBalancerCreateOptions{
		Region:             data.Region.ValueString(),
		Label:              data.Label.ValueStringPointer(),
		ClientConnThrottle: &clientConnThrottle,
	}

	if !data.FirewallID.IsNull() {
		createOpts.FirewallID = helper.FrameworkSafeInt64ToInt(
			data.FirewallID.ValueInt64(),
			&resp.Diagnostics,
		)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !data.Tags.IsNull() {
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &createOpts.Tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, "client.CreateNodeBalancer(...)", map[string]any{
		"options": createOpts,
	})

	nodebalancer, err := client.CreateNodeBalancer(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating a Linode NodeBalancer",
			err.Error(),
		)
		return
	}

	firewalls, err := client.ListNodeBalancerFirewalls(ctx, nodebalancer.ID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to list firewalls assigned to NodeBalancer %d", nodebalancer.ID),
			err.Error(),
		)
	}

	resp.Diagnostics.Append(data.FlattenNodeBalancer(ctx, nodebalancer, firewalls, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	data.ID = types.StringValue(strconv.Itoa(nodebalancer.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_nodebalancer")

	var data NodeBalancerModel
	client := r.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.StringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)
	tflog.Trace(ctx, "client.GetNodeBalancer(...)")

	nodeBalancer, err := client.GetNodeBalancer(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"NodeBalancer No Longer Exists",
				fmt.Sprintf("Removing Linode NodeBalancer ID %v from state because it no longer exists", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Get NodeBalancer %v", id),
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "client.ListNodeBalancerFirewalls(...)")

	firewalls, err := client.ListNodeBalancerFirewalls(ctx, id, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to list firewalls assigned to NodeBalancer %d", id),
			err.Error(),
		)
	}

	resp.Diagnostics.Append(data.FlattenNodeBalancer(ctx, nodeBalancer, firewalls, false)...)
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
	tflog.Debug(ctx, "Update linode_nodebalancer")

	var plan, state NodeBalancerModel
	client := r.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	ctx = populateLogAttributes(ctx, state)

	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.StringToInt(plan.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	isEqual := state.Label.Equal(plan.Label) &&
		state.ClientConnThrottle.Equal(plan.ClientConnThrottle) &&
		state.Tags.Equal(plan.Tags)

	if !isEqual {
		clientConnThrottle := helper.FrameworkSafeInt64ToInt(
			plan.ClientConnThrottle.ValueInt64(),
			&resp.Diagnostics,
		)
		if resp.Diagnostics.HasError() {
			return
		}
		updateOpts := linodego.NodeBalancerUpdateOptions{
			Label:              plan.Label.ValueStringPointer(),
			ClientConnThrottle: &clientConnThrottle,
		}
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &updateOpts.Tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "client.UpdateNodeBalancer(...)", map[string]any{
			"options": updateOpts,
		})

		nodeBalancer, err := client.UpdateNodeBalancer(ctx, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Update NodeBalancer %v", id),
				err.Error(),
			)
			return
		}

		tflog.Trace(ctx, "client.ListNodeBalancerFirewalls(...)")

		firewalls, err := client.ListNodeBalancerFirewalls(ctx, id, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to list firewalls assigned to NodeBalancer %d", id),
				err.Error(),
			)
		}

		resp.Diagnostics.Append(plan.FlattenNodeBalancer(ctx, nodeBalancer, firewalls, true)...)
	}

	plan.CopyFrom(state, true)

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_nodebalancer")

	var data NodeBalancerModel
	client := r.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.StringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)
	tflog.Debug(ctx, "client.DeleteNodeBalancer(...)")

	err := client.DeleteNodeBalancer(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"NodeBalancer No Longer Exists",
				fmt.Sprintf("NodeBalancer %v does not exist, removing it from state.", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to Delete NodeBalancer",
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

	nbDataV1 := NodeBalancerModel{
		ID:                 nbDataV0.ID,
		Label:              nbDataV0.Label,
		Region:             nbDataV0.Region,
		ClientConnThrottle: nbDataV0.ClientConnThrottle,
		Hostname:           nbDataV0.Hostname,
		IPv4:               nbDataV0.IPv4,
		IPv6:               nbDataV0.IPv6,
		Created:            timetypes.RFC3339{StringValue: nbDataV0.Created},
		Updated:            timetypes.RFC3339{StringValue: nbDataV0.Updated},
		Tags:               nbDataV0.Tags,
		Firewalls:          types.ListNull(firewallObjType),
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

func populateLogAttributes(ctx context.Context, model NodeBalancerModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"nodebalancer_id": model.ID.ValueString(),
	})
}
