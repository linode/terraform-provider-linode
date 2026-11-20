package nfsfilesystemaccesspolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

var linodeACLObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":    types.Int64Type,
		"label": types.StringType,
		"url":   types.StringType,
		"ipv6":  types.StringType,
	},
}

// ResourceModel represents the access policy of an NFS Filesystem.
type ResourceModel struct {
	ID           types.String      `tfsdk:"id"`
	SpaceID      types.Int64       `tfsdk:"space_id"`
	FilesystemID types.Int64       `tfsdk:"filesystem_id"`
	Label        types.String      `tfsdk:"label"`
	Enabled      types.Bool        `tfsdk:"enabled"`
	Inherit      types.Bool        `tfsdk:"inherit"`
	LinodeIDs    types.Set         `tfsdk:"linode_ids"`
	SquashPolicy types.String      `tfsdk:"squash_policy"`
	Protocols    types.Set         `tfsdk:"protocols"`
	LinodeACL    types.List        `tfsdk:"linode_acl"`
	Status       types.String      `tfsdk:"status"`
	Created      timetypes.RFC3339 `tfsdk:"created"`
	Updated      timetypes.RFC3339 `tfsdk:"updated"`
}

// LinodeACLModel represents a resolved Linode ACL entry.
type LinodeACLModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
	URL   types.String `tfsdk:"url"`
	IPv6  types.String `tfsdk:"ipv6"`
}

func (m *ResourceModel) Flatten(
	ctx context.Context,
	spaceID, filesystemID int,
	policy *linodego.NFSFilesystemAccessPolicy,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	m.ID = helper.KeepOrUpdateString(m.ID, fmt.Sprintf("%d,%d", spaceID, filesystemID), preserveKnown)
	m.SpaceID = helper.KeepOrUpdateInt64(m.SpaceID, int64(spaceID), preserveKnown)
	m.FilesystemID = helper.KeepOrUpdateInt64(m.FilesystemID, int64(filesystemID), preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, policy.Label, preserveKnown)
	m.Enabled = helper.KeepOrUpdateBool(m.Enabled, policy.Enabled, preserveKnown)
	m.Inherit = helper.KeepOrUpdateBool(m.Inherit, policy.Inherit, preserveKnown)
	m.SquashPolicy = helper.KeepOrUpdateString(m.SquashPolicy, string(policy.SquashPolicy), preserveKnown)
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

	protocols := helper.MapSlice(policy.Protocols, func(p linodego.NFSProtocolVersion) string {
		return string(p)
	})
	m.Protocols = helper.KeepOrUpdateStringSet(m.Protocols, protocols, preserveKnown, diags)

	linodeIDs := make([]int, len(policy.LinodeACL))
	acl := make([]LinodeACLModel, len(policy.LinodeACL))

	for i, entry := range policy.LinodeACL {
		linodeIDs[i] = entry.ID
		acl[i] = LinodeACLModel{
			ID:    types.Int64Value(int64(entry.ID)),
			Label: types.StringPointerValue(entry.Label),
			URL:   types.StringPointerValue(entry.URL),
			IPv6:  types.StringPointerValue(entry.IPv6),
		}
	}

	m.LinodeIDs = helper.KeepOrUpdateIntSet(m.LinodeIDs, linodeIDs, preserveKnown, diags)

	aclList, d := types.ListValueFrom(ctx, linodeACLObjectType, acl)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	m.LinodeACL = helper.KeepOrUpdateValue(m.LinodeACL, aclList, preserveKnown)
}

// GetUpdateOptions builds the access policy update options from the planned
// values that differ from the prior state. When state is nil, every known,
// non-null planned value is included.
func (m *ResourceModel) GetUpdateOptions(
	ctx context.Context,
	state *ResourceModel,
	diags *diag.Diagnostics,
) (opts linodego.NFSFilesystemAccessPolicyUpdateOptions, shouldUpdate bool) {
	if shouldSetField(m.Enabled, state, func(s *ResourceModel) attr.Value { return s.Enabled }) {
		opts.Enabled = m.Enabled.ValueBoolPointer()
		shouldUpdate = true
	}

	if shouldSetField(m.Inherit, state, func(s *ResourceModel) attr.Value { return s.Inherit }) {
		opts.Inherit = m.Inherit.ValueBoolPointer()
		shouldUpdate = true
	}

	if shouldSetField(m.SquashPolicy, state, func(s *ResourceModel) attr.Value { return s.SquashPolicy }) {
		opts.SquashPolicy = linodego.Pointer(linodego.NFSSquashPolicy(m.SquashPolicy.ValueString()))
		shouldUpdate = true
	}

	if shouldSetField(m.LinodeIDs, state, func(s *ResourceModel) attr.Value { return s.LinodeIDs }) {
		var ids []int64
		diags.Append(m.LinodeIDs.ElementsAs(ctx, &ids, false)...)
		if diags.HasError() {
			return opts, false
		}

		linodeIDs := make([]int, len(ids))
		for i, id := range ids {
			linodeIDs[i] = helper.FrameworkSafeInt64ToInt(id, diags)
		}
		if diags.HasError() {
			return opts, false
		}

		opts.LinodeIDs = &linodeIDs
		shouldUpdate = true
	}

	if shouldSetField(m.Protocols, state, func(s *ResourceModel) attr.Value { return s.Protocols }) {
		var protocols []string
		diags.Append(m.Protocols.ElementsAs(ctx, &protocols, false)...)
		if diags.HasError() {
			return opts, false
		}

		protocolValues := make([]linodego.NFSProtocolVersion, len(protocols))
		for i, p := range protocols {
			protocolValues[i] = linodego.NFSProtocolVersion(p)
		}
		opts.Protocols = &protocolValues
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
