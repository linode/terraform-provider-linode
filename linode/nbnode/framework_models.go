package nbnode

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type BaseModel struct {
	NodeBalancerID types.Int64  `tfsdk:"nodebalancer_id"`
	ConfigID       types.Int64  `tfsdk:"config_id"`
	Label          types.String `tfsdk:"label"`
	Weight         types.Int64  `tfsdk:"weight"`
	Mode           types.String `tfsdk:"mode"`
	Address        types.String `tfsdk:"address"`
	Status         types.String `tfsdk:"status"`
	SubnetID       types.Int64  `tfsdk:"subnet_id"`
	VPCConfigID    types.Int64  `tfsdk:"vpc_config_id"`
}

func (data *BaseModel) Flatten(
	node *linodego.NodeBalancerNode,
	vpcConfig *linodego.NodeBalancerVPCConfig,
	preserveKnown bool,
) (diags diag.Diagnostics) {
	data.NodeBalancerID = helper.KeepOrUpdateInt64(data.NodeBalancerID, int64(node.NodeBalancerID), preserveKnown)
	data.ConfigID = helper.KeepOrUpdateInt64(data.ConfigID, int64(node.ConfigID), preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, node.Label, preserveKnown)
	data.Weight = helper.KeepOrUpdateInt64(data.Weight, int64(node.Weight), preserveKnown)
	data.Mode = helper.KeepOrUpdateString(data.Mode, string(node.Mode), preserveKnown)
	data.Address = helper.KeepOrUpdateString(data.Address, node.Address, preserveKnown)
	data.Status = helper.KeepOrUpdateString(data.Status, node.Status, preserveKnown)

	if vpcConfig != nil {
		data.SubnetID = helper.KeepOrUpdateValue(
			data.SubnetID,
			types.Int64Value(int64(vpcConfig.SubnetID)),
			preserveKnown,
		)
		fmt.Println(vpcConfig.ID)
		data.VPCConfigID = helper.KeepOrUpdateValue(
			data.VPCConfigID,
			types.Int64Value(int64(vpcConfig.ID)),
			preserveKnown,
		)
	} else {
		data.SubnetID = helper.KeepOrUpdateValue(data.SubnetID, types.Int64Null(), preserveKnown)
		data.VPCConfigID = helper.KeepOrUpdateValue(data.SubnetID, types.Int64Null(), preserveKnown)
	}

	return diags
}

type DataSourceModel struct {
	ID types.Int64 `tfsdk:"id"`
	BaseModel
}

func (data *DataSourceModel) Flatten(
	node *linodego.NodeBalancerNode,
	vpcConfig *linodego.NodeBalancerVPCConfig,
) (diags diag.Diagnostics) {
	data.ID = helper.KeepOrUpdateInt64(data.ID, int64(node.ID), false)
	diags.Append(data.BaseModel.Flatten(node, vpcConfig, false)...)

	return diags
}

// TODO: consider merging two models when resource's ID change to int type
type ResourceModel struct {
	ID types.String `tfsdk:"id"`
	BaseModel
}

func (data *ResourceModel) Flatten(
	node *linodego.NodeBalancerNode,
	vpcConfig *linodego.NodeBalancerVPCConfig,
	preserveKnown bool,
) (diags diag.Diagnostics) {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(node.ID), preserveKnown)
	diags.Append(data.BaseModel.Flatten(node, vpcConfig, preserveKnown)...)

	return nil
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
	subnetID := helper.FrameworkSafeInt64ToInt(plan.SubnetID.ValueInt64(), diags)

	return linodego.NodeBalancerNodeCreateOptions{
		Address:  plan.Address.ValueString(),
		Label:    plan.Label.ValueString(),
		Weight:   weight,
		Mode:     linodego.NodeMode(plan.Mode.ValueString()),
		SubnetID: subnetID,
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
			return result
		}
		result.Weight = weight
	}

	if !plan.Mode.Equal(state.Mode) {
		result.Mode = linodego.NodeMode(plan.Mode.ValueString())
	}

	return result
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
	plan.SubnetID = helper.KeepOrUpdateValue(plan.SubnetID, state.SubnetID, preserveKnown)
	plan.VPCConfigID = helper.KeepOrUpdateValue(plan.VPCConfigID, state.VPCConfigID, preserveKnown)
}
