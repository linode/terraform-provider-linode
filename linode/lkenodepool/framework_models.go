package lkenodepool

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type NodePoolModel struct {
	ID             types.String              `tfsdk:"id"`
	ClusterID      types.Int64               `tfsdk:"cluster_id"`
	Count          types.Int64               `tfsdk:"node_count"`
	Type           types.String              `tfsdk:"type"`
	DiskEncryption types.String              `tfsdk:"disk_encryption"`
	Tags           types.Set                 `tfsdk:"tags"`
	Nodes          types.List                `tfsdk:"nodes"`
	Autoscaler     []NodePoolAutoscalerModel `tfsdk:"autoscaler"`
	Taints         []NodePoolTaintModel      `tfsdk:"taint"`
	Labels         types.Map                 `tfsdk:"labels"`
	K8sVersion     types.String              `tfsdk:"k8s_version"`
	UpdateStrategy types.String              `tfsdk:"update_strategy"`
}

type NodePoolAutoscalerModel struct {
	Min types.Int64 `tfsdk:"min"`
	Max types.Int64 `tfsdk:"max"`
}

type NodePoolTaintModel struct {
	Effect types.String `tfsdk:"effect"`
	Key    types.String `tfsdk:"key"`
	Value  types.String `tfsdk:"value"`
}

func flattenLKENodePoolLinode(node linodego.LKENodePoolLinode) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["id"] = types.StringValue(node.ID)
	result["instance_id"] = types.Int64Value(int64(node.InstanceID))
	result["status"] = types.StringValue(string(node.Status))

	obj, errors := types.ObjectValue(nodeObjectType.AttrTypes, result)
	if errors.HasError() {
		return nil, errors
	}
	return &obj, nil
}

func flattenLKENodePoolLinodeList(nodes []linodego.LKENodePoolLinode,
) (*basetypes.ListValue, diag.Diagnostics) {
	resultList := make([]attr.Value, len(nodes))
	for i, node := range nodes {
		result, errors := flattenLKENodePoolLinode(node)
		if errors.HasError() {
			return nil, errors
		}

		resultList[i] = result
	}
	result, errors := basetypes.NewListValue(
		nodeObjectType,
		resultList,
	)
	if errors.HasError() {
		return nil, errors
	}

	return &result, nil
}

func (taint *NodePoolTaintModel) FlattenLKENodePoolTaint(t linodego.LKENodePoolTaint, preserveKnown bool) {
	taint.Effect = helper.KeepOrUpdateString(taint.Effect, string(t.Effect), preserveKnown)
	taint.Key = helper.KeepOrUpdateString(taint.Key, t.Key, preserveKnown)
	taint.Value = helper.KeepOrUpdateString(taint.Value, t.Value, preserveKnown)
}

func (pool *NodePoolModel) FlattenLKENodePoolTaints(taints []linodego.LKENodePoolTaint, preserveKnown bool) {
	// taints block can't be computed and can't be modified if known values are preserved.
	if preserveKnown {
		return
	}

	pool.Taints = make([]NodePoolTaintModel, len(taints))
	for i := range pool.Taints {
		pool.Taints[i].FlattenLKENodePoolTaint(taints[i], preserveKnown)
	}
}

func (pool *NodePoolModel) FlattenLKENodePool(
	ctx context.Context, p *linodego.LKENodePool, preserveKnown bool, diags *diag.Diagnostics,
) {
	pool.ID = helper.KeepOrUpdateString(pool.ID, strconv.Itoa(p.ID), preserveKnown)
	pool.Count = helper.KeepOrUpdateInt64(pool.Count, int64(p.Count), preserveKnown)
	pool.Type = helper.KeepOrUpdateString(pool.Type, p.Type, preserveKnown)
	pool.DiskEncryption = helper.KeepOrUpdateString(pool.DiskEncryption, string(p.DiskEncryption), preserveKnown)
	pool.Tags = helper.KeepOrUpdateStringSet(pool.Tags, p.Tags, preserveKnown, diags)
	if diags.HasError() {
		return
	}

	pool.Labels = helper.KeepOrUpdateStringMap(ctx, pool.Labels, p.Labels, preserveKnown, diags)

	if !preserveKnown {
		if p.Autoscaler.Enabled {
			pool.Autoscaler = []NodePoolAutoscalerModel{
				{
					Min: types.Int64Value(int64(p.Autoscaler.Min)),
					Max: types.Int64Value(int64(p.Autoscaler.Max)),
				},
			}
		}
		pool.FlattenLKENodePoolTaints(p.Taints, preserveKnown)
	}

	nodePoolLinodes, errs := flattenLKENodePoolLinodeList(p.Linodes)
	if errs != nil {
		diags.Append(errs...)
	}
	pool.Nodes = helper.KeepOrUpdateValue(pool.Nodes, *nodePoolLinodes, preserveKnown)

	pool.K8sVersion = helper.KeepOrUpdateStringPointer(pool.K8sVersion, p.K8sVersion, preserveKnown)

	if p.UpdateStrategy != nil {
		pool.UpdateStrategy = helper.KeepOrUpdateString(pool.UpdateStrategy, string(*p.UpdateStrategy), preserveKnown)
	} else {
		pool.UpdateStrategy = helper.KeepOrUpdateString(pool.UpdateStrategy, "", preserveKnown)
	}
}

func (pool *NodePoolModel) SetNodePoolCreateOptions(
	ctx context.Context,
	p *linodego.LKENodePoolCreateOptions,
	diags *diag.Diagnostics,
	tier string,
) {
	p.Count = helper.FrameworkSafeInt64ToInt(
		pool.Count.ValueInt64(),
		diags,
	)
	p.Type = pool.Type.ValueString()

	if !pool.Tags.IsNull() {
		diags.Append(pool.Tags.ElementsAs(ctx, &p.Tags, false)...)
	}

	p.Autoscaler = pool.getLKENodePoolAutoscaler(p.Count, diags)
	if p.Autoscaler.Enabled && p.Count == 0 {
		p.Count = p.Autoscaler.Min
	}

	p.Taints = pool.getLKENodePoolTaints()

	pool.Labels.ElementsAs(ctx, &p.Labels, false)

	if tier == "enterprise" {
		if !pool.K8sVersion.IsNull() && !pool.K8sVersion.IsUnknown() {
			p.K8sVersion = pool.K8sVersion.ValueStringPointer()
		}
		if !pool.UpdateStrategy.IsNull() && !pool.UpdateStrategy.IsUnknown() {
			p.UpdateStrategy = linodego.Pointer(linodego.LKENodePoolUpdateStrategy(pool.UpdateStrategy.ValueString()))
		}
	}
}

func (pool *NodePoolModel) SetNodePoolUpdateOptions(
	ctx context.Context,
	p *linodego.LKENodePoolUpdateOptions,
	diags *diag.Diagnostics,
	state *NodePoolModel,
	tier string,
) bool {
	var shouldUpdate bool

	if !state.Count.Equal(pool.Count) {
		p.Count = helper.FrameworkSafeInt64ToInt(
			pool.Count.ValueInt64(),
			diags,
		)
		if diags.HasError() {
			return false
		}

		shouldUpdate = true
	}

	if !state.Tags.Equal(pool.Tags) {
		if !pool.Tags.IsNull() {
			diags.Append(pool.Tags.ElementsAs(ctx, &p.Tags, false)...)
			if diags.HasError() {
				return false
			}
		}

		shouldUpdate = true
	}

	autoscaler, asNeedsUpdate := pool.shouldUpdateLKENodePoolAutoscaler(p.Count, state, diags)
	if diags.HasError() {
		return false
	}

	if asNeedsUpdate {
		p.Autoscaler = autoscaler
		if p.Autoscaler.Enabled && p.Count == 0 {
			p.Count = p.Autoscaler.Min
		}
	}

	shouldUpdate = shouldUpdate || asNeedsUpdate

	if !(len(state.Taints) == 0 && len(pool.Taints) == 0) {
		taints := pool.getLKENodePoolTaints()
		p.Taints = &taints
		shouldUpdate = true
	}

	if !state.Labels.Equal(pool.Labels) {
		pool.Labels.ElementsAs(ctx, &p.Labels, false)
		shouldUpdate = true
	}

	if tier == "enterprise" {
		if !state.K8sVersion.Equal(pool.K8sVersion) {
			p.K8sVersion = pool.K8sVersion.ValueStringPointer()
			shouldUpdate = true
		}

		if !state.UpdateStrategy.Equal(pool.UpdateStrategy) {
			p.UpdateStrategy = linodego.Pointer(linodego.LKENodePoolUpdateStrategy(pool.UpdateStrategy.ValueString()))
			shouldUpdate = true
		}
	}

	return shouldUpdate
}

func (pool *NodePoolModel) ExtractClusterAndNodePoolIDs(diags *diag.Diagnostics) (int, int) {
	clusterID := helper.FrameworkSafeInt64ToInt(pool.ClusterID.ValueInt64(), diags)
	poolID, err := strconv.Atoi(pool.ID.ValueString())
	if err != nil {
		diags.AddError("Failed to parse poolID", err.Error())
	}
	return clusterID, poolID
}

func (pool *NodePoolModel) getLKENodePoolAutoscaler(count int, diags *diag.Diagnostics) *linodego.LKENodePoolAutoscaler {
	var autoscaler linodego.LKENodePoolAutoscaler
	if len(pool.Autoscaler) > 0 {
		autoscaler.Enabled = true
		autoscaler.Min = helper.FrameworkSafeInt64ToInt(pool.Autoscaler[0].Min.ValueInt64(), diags)
		autoscaler.Max = helper.FrameworkSafeInt64ToInt(pool.Autoscaler[0].Max.ValueInt64(), diags)
	} else {
		autoscaler.Enabled = false
		autoscaler.Min = count
		autoscaler.Max = count
	}
	return &autoscaler
}

func (pool *NodePoolModel) shouldUpdateLKENodePoolAutoscaler(
	count int,
	state *NodePoolModel,
	diags *diag.Diagnostics,
) (*linodego.LKENodePoolAutoscaler, bool) {
	var autoscaler linodego.LKENodePoolAutoscaler
	var shouldUpdate bool

	if len(pool.Autoscaler) > 0 {
		if len(state.Autoscaler) == 0 ||
			(len(state.Autoscaler) > 0 && (!state.Autoscaler[0].Min.Equal(pool.Autoscaler[0].Min) ||
				!state.Autoscaler[0].Max.Equal(pool.Autoscaler[0].Max))) {
			autoscaler.Enabled = true
			autoscaler.Min = helper.FrameworkSafeInt64ToInt(pool.Autoscaler[0].Min.ValueInt64(), diags)
			autoscaler.Max = helper.FrameworkSafeInt64ToInt(pool.Autoscaler[0].Max.ValueInt64(), diags)

			shouldUpdate = true
		}
	} else if len(state.Autoscaler) > 0 {
		autoscaler.Enabled = false
		autoscaler.Min = count
		autoscaler.Max = count

		shouldUpdate = true
	}

	return &autoscaler, shouldUpdate
}

func (taint NodePoolTaintModel) getLKENodePoolTaint() linodego.LKENodePoolTaint {
	return linodego.LKENodePoolTaint{
		Effect: linodego.LKENodePoolTaintEffect(taint.Effect.ValueString()),
		Key:    taint.Key.ValueString(),
		Value:  taint.Value.ValueString(),
	}
}

func (pool *NodePoolModel) getLKENodePoolTaints() []linodego.LKENodePoolTaint {
	taints := make([]linodego.LKENodePoolTaint, len(pool.Taints))

	for i := range pool.Taints {
		taints[i] = pool.Taints[i].getLKENodePoolTaint()
	}

	return taints
}

func (data *NodePoolModel) CopyFrom(other NodePoolModel, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.ClusterID = helper.KeepOrUpdateValue(data.ClusterID, other.ClusterID, preserveKnown)
	data.Count = helper.KeepOrUpdateValue(data.Count, other.Count, preserveKnown)
	data.Type = helper.KeepOrUpdateValue(data.Type, other.Type, preserveKnown)
	data.DiskEncryption = helper.KeepOrUpdateValue(data.DiskEncryption, other.DiskEncryption, preserveKnown)
	data.Tags = helper.KeepOrUpdateValue(data.Tags, other.Tags, preserveKnown)
	data.Nodes = helper.KeepOrUpdateValue(data.Nodes, other.Nodes, preserveKnown)
	data.Labels = helper.KeepOrUpdateValue(data.Labels, other.Labels, preserveKnown)
	data.K8sVersion = helper.KeepOrUpdateValue(data.K8sVersion, other.K8sVersion, preserveKnown)
	data.UpdateStrategy = helper.KeepOrUpdateValue(data.UpdateStrategy, other.UpdateStrategy, preserveKnown)

	if !preserveKnown {
		data.Autoscaler = other.Autoscaler
		data.Taints = other.Taints
	}
}
