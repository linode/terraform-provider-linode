package nbnode

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type BaseModel struct {
	NodeBalancerID types.Int64  `tfsdk:"nodebalancer_id"`
	ConfigID       types.Int64  `tfsdk:"config_id"`
	Label          types.String `tfsdk:"label"`
	Weight         types.Int64  `tfsdk:"weight"`
	Mode           types.String `tfsdk:"mode"`
	Address        types.String `tfsdk:"address"`
	Status         types.String `tfsdk:"status"`
}

type DataSourceModel struct {
	ID types.Int64 `tfsdk:"id"`
	BaseModel
}

func (data *DataSourceModel) ParseNodeBalancerNode(nbnode *linodego.NodeBalancerNode) {
	data.ID = types.Int64Value(int64(nbnode.ID))
	data.NodeBalancerID = types.Int64Value(int64(nbnode.NodeBalancerID))
	data.ConfigID = types.Int64Value(int64(nbnode.ConfigID))
	data.Label = types.StringValue(nbnode.Label)
	data.Weight = types.Int64Value(int64(nbnode.Weight))
	data.Mode = types.StringValue(string(nbnode.Mode))
	data.Address = types.StringValue(nbnode.Address)
	data.Status = types.StringValue(nbnode.Status)
}

// TODO: consider merging two models when resource's ID change to int type
type ResourceModel struct {
	ID types.String `tfsdk:"id"`
	BaseModel
}

func (data *ResourceModel) FlattenNodeBalancerNode(
	nbnode *linodego.NodeBalancerNode, preserveKnown bool,
) {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(nbnode.ID), preserveKnown)
	data.NodeBalancerID = helper.KeepOrUpdateInt64(data.NodeBalancerID, int64(nbnode.NodeBalancerID), preserveKnown)
	data.ConfigID = helper.KeepOrUpdateInt64(data.ConfigID, int64(nbnode.ConfigID), preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, nbnode.Label, preserveKnown)
	data.Weight = helper.KeepOrUpdateInt64(data.Weight, int64(nbnode.Weight), preserveKnown)
	data.Mode = helper.KeepOrUpdateString(data.Mode, string(nbnode.Mode), preserveKnown)
	data.Address = helper.KeepOrUpdateString(data.Address, nbnode.Address, preserveKnown)
	data.Status = helper.KeepOrUpdateString(data.Status, nbnode.Status, preserveKnown)
}

func (data *ResourceModel) GetIDs(diags *diag.Diagnostics) (int, int, int) {
	id := helper.StringToInt(data.ID.ValueString(), diags)
	nodeBalancerID := helper.FrameworkSafeInt64ToInt(data.NodeBalancerID.ValueInt64(), diags)
	configID := helper.FrameworkSafeInt64ToInt(data.ConfigID.ValueInt64(), diags)

	return id, nodeBalancerID, configID
}

func (data *ResourceModel) GetCreateParameters(diags *diag.Diagnostics) (int, int, linodego.NodeBalancerNodeCreateOptions) {
	nodeBalancerID := helper.FrameworkSafeInt64ToInt(data.NodeBalancerID.ValueInt64(), diags)
	configID := helper.FrameworkSafeInt64ToInt(data.ConfigID.ValueInt64(), diags)
	return nodeBalancerID, configID, data.GetCreateOptions(diags)
}

func (plan *ResourceModel) GetCreateOptions(diags *diag.Diagnostics) linodego.NodeBalancerNodeCreateOptions {
	weight := helper.FrameworkSafeInt64ToInt(plan.Weight.ValueInt64(), diags)
	return linodego.NodeBalancerNodeCreateOptions{
		Address: plan.Address.ValueString(),
		Label:   plan.Label.ValueString(),
		Weight:  weight,
		Mode:    linodego.NodeMode(plan.Mode.ValueString()),
	}
}

func (plan *ResourceModel) GetUpdateOptions(
	state ResourceModel, diags *diag.Diagnostics,
) (result linodego.NodeBalancerNodeUpdateOptions) {
	if !plan.Address.Equal(state.Address) {
		result.Address = plan.Address.ValueString()
	}

	if !plan.Label.Equal(state.Label) {
		result.Label = plan.Label.ValueString()
	}

	if !plan.Weight.Equal(state.Weight) {
		weight := helper.FrameworkSafeInt64ToInt(plan.Weight.ValueInt64(), diags)
		if diags.HasError() {
			return
		}
		result.Weight = weight
	}

	if !plan.Mode.Equal(state.Mode) {
		result.Mode = linodego.NodeMode(plan.Mode.ValueString())
	}

	return
}

func (plan *ResourceModel) CopyFrom(state ResourceModel, preserveKnown bool) {
	plan.ID = helper.KeepOrUpdateValue(plan.ID, state.ID, preserveKnown)
	plan.NodeBalancerID = helper.KeepOrUpdateValue(plan.NodeBalancerID, state.NodeBalancerID, preserveKnown)
	plan.ConfigID = helper.KeepOrUpdateValue(plan.ConfigID, state.ConfigID, preserveKnown)
	plan.Label = helper.KeepOrUpdateValue(plan.Label, state.Label, preserveKnown)
	plan.Weight = helper.KeepOrUpdateValue(plan.Weight, state.Weight, preserveKnown)
	plan.Mode = helper.KeepOrUpdateValue(plan.Mode, state.Mode, preserveKnown)
	plan.Address = helper.KeepOrUpdateValue(plan.Address, state.Address, preserveKnown)
	plan.Status = helper.KeepOrUpdateValue(plan.Status, state.Status, preserveKnown)
}
