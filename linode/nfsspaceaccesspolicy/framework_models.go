package nfsspaceaccesspolicy

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

var vpcObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":      types.Int64Type,
		"subnets": types.SetType{ElemType: types.Int64Type},
	},
}

var vpcACLSubnetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":    types.Int64Type,
		"label": types.StringType,
		"url":   types.StringType,
		"range": types.StringType,
	},
}

var vpcACLObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":      types.Int64Type,
		"label":   types.StringType,
		"url":     types.StringType,
		"range":   types.StringType,
		"subnets": types.ListType{ElemType: vpcACLSubnetObjectType},
	},
}

// ResourceModel represents the access policy of an NFS Space.
type ResourceModel struct {
	ID         types.String      `tfsdk:"id"`
	SpaceID    types.Int64       `tfsdk:"space_id"`
	Label      types.String      `tfsdk:"label"`
	Enabled    types.Bool        `tfsdk:"enabled"`
	VPCs       types.Set         `tfsdk:"vpcs"`
	MTLSCACert types.String      `tfsdk:"mtls_ca_cert"`
	MTLSMode   types.String      `tfsdk:"mtls_mode"`
	VPCACL     types.List        `tfsdk:"vpc_acl"`
	Status     types.String      `tfsdk:"status"`
	Created    timetypes.RFC3339 `tfsdk:"created"`
	Updated    timetypes.RFC3339 `tfsdk:"updated"`
}

// VPCModel represents a VPC (and optionally specific subnets) permitted by the
// space-level access policy.
type VPCModel struct {
	ID      types.Int64 `tfsdk:"id"`
	Subnets types.Set   `tfsdk:"subnets"`
}

// VPCACLModel represents a resolved VPC ACL entry.
type VPCACLModel struct {
	ID      types.Int64  `tfsdk:"id"`
	Label   types.String `tfsdk:"label"`
	URL     types.String `tfsdk:"url"`
	Range   types.String `tfsdk:"range"`
	Subnets types.List   `tfsdk:"subnets"`
}

// VPCACLSubnetModel represents a resolved subnet within a VPC ACL entry.
type VPCACLSubnetModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
	URL   types.String `tfsdk:"url"`
	Range types.String `tfsdk:"range"`
}

func (m *ResourceModel) Flatten(
	ctx context.Context,
	spaceID int,
	policy *linodego.NFSSpaceAccessPolicy,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(spaceID), preserveKnown)
	m.SpaceID = helper.KeepOrUpdateInt64(m.SpaceID, int64(spaceID), preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, policy.Label, preserveKnown)
	m.Enabled = helper.KeepOrUpdateBool(m.Enabled, policy.Enabled, preserveKnown)
	m.MTLSCACert = helper.KeepOrUpdateStringPointer(m.MTLSCACert, policy.MTLSCACert, preserveKnown)
	m.MTLSMode = helper.KeepOrUpdateString(m.MTLSMode, string(policy.MTLSMode), preserveKnown)
	m.Status = helper.KeepOrUpdateString(m.Status, string(policy.Status), preserveKnown)
	m.Created = helper.KeepOrUpdateValue(
		m.Created,
		timetypes.NewRFC3339TimePointerValue(policy.Created),
		preserveKnown,
	)
	m.Updated = helper.KeepOrUpdateValue(
		m.Updated,
		timetypes.NewRFC3339TimePointerValue(policy.Updated),
		preserveKnown,
	)

	vpcs := make([]VPCModel, len(policy.VPCACL))
	acl := make([]VPCACLModel, len(policy.VPCACL))

	for i, entry := range policy.VPCACL {
		subnetIDs := make([]int, len(entry.Subnets))
		aclSubnets := make([]VPCACLSubnetModel, len(entry.Subnets))

		for j, subnet := range entry.Subnets {
			subnetIDs[j] = subnet.ID
			aclSubnets[j] = VPCACLSubnetModel{
				ID:    types.Int64Value(int64(subnet.ID)),
				Label: types.StringPointerValue(subnet.Label),
				URL:   types.StringPointerValue(subnet.URL),
				Range: types.StringPointerValue(subnet.Range),
			}
		}

		// The API only returns subnets when specific subnets were requested,
		// so an empty list maps back to an unset `subnets` attribute.
		subnets := types.SetNull(types.Int64Type)
		if len(subnetIDs) > 0 {
			var d diag.Diagnostics
			subnets, d = types.SetValue(types.Int64Type, helper.IntSliceToFrameworkValueSlice(subnetIDs))
			diags.Append(d...)
		}

		vpcs[i] = VPCModel{
			ID:      types.Int64Value(int64(entry.ID)),
			Subnets: subnets,
		}

		aclSubnetList, d := types.ListValueFrom(ctx, vpcACLSubnetObjectType, aclSubnets)
		diags.Append(d...)

		acl[i] = VPCACLModel{
			ID:      types.Int64Value(int64(entry.ID)),
			Label:   types.StringPointerValue(entry.Label),
			URL:     types.StringPointerValue(entry.URL),
			Range:   types.StringPointerValue(entry.Range),
			Subnets: aclSubnetList,
		}
	}

	if diags.HasError() {
		return
	}

	vpcSet, d := types.SetValueFrom(ctx, vpcObjectType, vpcs)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	m.VPCs = helper.KeepOrUpdateValue(m.VPCs, vpcSet, preserveKnown)

	aclList, d := types.ListValueFrom(ctx, vpcACLObjectType, acl)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	m.VPCACL = helper.KeepOrUpdateValue(m.VPCACL, aclList, preserveKnown)
}

// GetUpdateOptions builds the access policy update options from the planned
// values that differ from the prior state. When state is nil, every known,
// non-null planned value is included.
func (m *ResourceModel) GetUpdateOptions(
	ctx context.Context,
	state *ResourceModel,
	diags *diag.Diagnostics,
) (opts linodego.NFSSpaceAccessPolicyUpdateOptions, shouldUpdate bool) {
	if shouldSetField(m.Enabled, state, func(s *ResourceModel) attr.Value { return s.Enabled }) {
		opts.Enabled = m.Enabled.ValueBoolPointer()
		shouldUpdate = true
	}

	if shouldSetField(m.MTLSCACert, state, func(s *ResourceModel) attr.Value { return s.MTLSCACert }) {
		opts.MTLSCACert = m.MTLSCACert.ValueStringPointer()
		shouldUpdate = true
	}

	if shouldSetField(m.MTLSMode, state, func(s *ResourceModel) attr.Value { return s.MTLSMode }) {
		opts.MTLSMode = linodego.Pointer(linodego.NFSMTLSMode(m.MTLSMode.ValueString()))
		shouldUpdate = true
	}

	if shouldSetField(m.VPCs, state, func(s *ResourceModel) attr.Value { return s.VPCs }) {
		var vpcModels []VPCModel
		diags.Append(m.VPCs.ElementsAs(ctx, &vpcModels, false)...)
		if diags.HasError() {
			return opts, false
		}

		vpcs := make([]linodego.NFSSpaceAccessPolicyVPCOptions, len(vpcModels))
		for i, vpc := range vpcModels {
			vpcs[i].ID = helper.FrameworkSafeInt64ToInt(vpc.ID.ValueInt64(), diags)

			if !vpc.Subnets.IsNull() && !vpc.Subnets.IsUnknown() {
				var subnetIDs []int64
				diags.Append(vpc.Subnets.ElementsAs(ctx, &subnetIDs, false)...)
				for _, subnetID := range subnetIDs {
					vpcs[i].Subnets = append(vpcs[i].Subnets, helper.FrameworkSafeInt64ToInt(subnetID, diags))
				}
			}
		}
		if diags.HasError() {
			return opts, false
		}

		opts.VPCs = &vpcs
		shouldUpdate = true
	}

	return opts, shouldUpdate
}

func shouldSetField(
	planned attr.Value,
	state *ResourceModel,
	getStateValue func(*ResourceModel) attr.Value,
) bool {
	if planned.IsNull() || planned.IsUnknown() {
		return false
	}

	return state == nil || !planned.Equal(getStateValue(state))
}
