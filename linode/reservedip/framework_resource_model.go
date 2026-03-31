package reservedip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/instancenetworking"
)

type ResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Region         types.String `tfsdk:"region"`
	Address        types.String `tfsdk:"address"`
	Gateway        types.String `tfsdk:"gateway"`
	SubnetMask     types.String `tfsdk:"subnet_mask"`
	Prefix         types.Int64  `tfsdk:"prefix"`
	Type           types.String `tfsdk:"type"`
	Public         types.Bool   `tfsdk:"public"`
	RDNS           types.String `tfsdk:"rdns"`
	LinodeID       types.Int64  `tfsdk:"linode_id"`
	Reserved       types.Bool   `tfsdk:"reserved"`
	Tags           types.Set    `tfsdk:"tags"`
	VPCNAT1To1     types.List   `tfsdk:"vpc_nat_1_1"`
	AssignedEntity types.List   `tfsdk:"assigned_entity"`
}

func (m *ResourceModel) flatten(
	ctx context.Context,
	ip linodego.InstanceIP,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = helper.KeepOrUpdateString(m.ID, ip.Address, preserveKnown)
	m.Address = helper.KeepOrUpdateString(m.Address, ip.Address, preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, ip.Region, preserveKnown)
	m.Gateway = helper.KeepOrUpdateString(m.Gateway, ip.Gateway, preserveKnown)
	m.SubnetMask = helper.KeepOrUpdateString(m.SubnetMask, ip.SubnetMask, preserveKnown)
	m.Prefix = helper.KeepOrUpdateInt64(m.Prefix, int64(ip.Prefix), preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, string(ip.Type), preserveKnown)
	m.Public = helper.KeepOrUpdateBool(m.Public, ip.Public, preserveKnown)
	m.RDNS = helper.KeepOrUpdateString(m.RDNS, ip.RDNS, preserveKnown)
	if ip.LinodeID != 0 {
		m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, types.Int64Value(int64(ip.LinodeID)), preserveKnown)
	} else {
		m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, types.Int64Null(), preserveKnown)
	}
	m.Reserved = helper.KeepOrUpdateBool(m.Reserved, ip.Reserved, preserveKnown)

	// Tags: if the API returns nil (endpoint omits the field or returns null),
	// preserve the existing state/plan value rather than clobbering with null.
	// When the existing value is also unknown or null (e.g. fresh create with
	// no tags configured), default to an empty set.
	if ip.Tags != nil {
		tags, tagDiags := types.SetValueFrom(ctx, types.StringType, ip.Tags)
		diags.Append(tagDiags...)
		m.Tags = tags
	} else if m.Tags.IsNull() || m.Tags.IsUnknown() {
		m.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// vpc_nat_1_1
	var vpcNATList types.List
	if ip.VPCNAT1To1 == nil {
		vpcNATList = types.ListNull(instancenetworking.VPCNAT1To1Type)
	} else {
		vpcObj, flatDiags := instancenetworking.FlattenIPVPCNAT1To1(ip.VPCNAT1To1)
		diags.Append(flatDiags...)
		var listDiags diag.Diagnostics
		vpcNATList, listDiags = types.ListValue(
			instancenetworking.VPCNAT1To1Type,
			[]attr.Value{vpcObj},
		)
		diags.Append(listDiags...)
	}
	m.VPCNAT1To1 = helper.KeepOrUpdateValue(m.VPCNAT1To1, vpcNATList, preserveKnown)

	// assigned_entity
	var assignedList types.List
	if ip.AssignedEntity == nil {
		assignedList = types.ListNull(assignedEntityType)
	} else {
		var listDiags diag.Diagnostics
		assignedList, listDiags = types.ListValueFrom(
			ctx,
			assignedEntityType,
			[]AssignedEntityModel{{
				ID:    types.Int64Value(int64(ip.AssignedEntity.ID)),
				Label: types.StringValue(ip.AssignedEntity.Label),
				Type:  types.StringValue(ip.AssignedEntity.Type),
				URL:   types.StringValue(ip.AssignedEntity.URL),
			}},
		)
		diags.Append(listDiags...)
	}
	m.AssignedEntity = helper.KeepOrUpdateValue(m.AssignedEntity, assignedList, preserveKnown)

	return diags
}
