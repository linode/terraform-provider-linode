package placementgroup

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type PlacementGroupModel struct {
	ID           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	Region       types.String `tfsdk:"region"`
	AffinityType types.String `tfsdk:"affinity_type"`
	IsStrict     types.Bool   `tfsdk:"is_strict"`
	IsCompliant  types.Bool   `tfsdk:"is_compliant"`

	Members types.Set `tfsdk:"members"`
}

func (m *PlacementGroupModel) FlattenPlacementGroup(
	pg *linodego.PlacementGroup,
	preserveKnown bool,
) (resultDiags diag.Diagnostics) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(pg.ID), preserveKnown)

	m.Label = helper.KeepOrUpdateString(m.Label, pg.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, pg.Region, preserveKnown)
	m.AffinityType = helper.KeepOrUpdateString(m.AffinityType, string(pg.AffinityType), preserveKnown)
	m.IsStrict = helper.KeepOrUpdateBool(m.IsStrict, pg.IsStrict, preserveKnown)
	m.IsCompliant = helper.KeepOrUpdateBool(m.IsCompliant, pg.IsCompliant, preserveKnown)

	members := make([]attr.Value, len(pg.Members))
	for i, member := range pg.Members {
		flattenedMember, d := flattenPlacementGroupMember(member)
		resultDiags.Append(d...)

		if resultDiags.HasError() {
			return
		}

		members[i] = flattenedMember
	}

	m.Members = helper.KeepOrUpdateSet(
		pgMemberObjectType,
		m.Members,
		members,
		preserveKnown,
		&resultDiags,
	)
	return
}

func (m *PlacementGroupModel) CopyFrom(other PlacementGroupModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)

	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
	m.AffinityType = helper.KeepOrUpdateValue(m.AffinityType, other.AffinityType, preserveKnown)
	m.IsStrict = helper.KeepOrUpdateValue(m.IsStrict, other.IsStrict, preserveKnown)
	m.IsCompliant = helper.KeepOrUpdateValue(m.IsCompliant, other.IsCompliant, preserveKnown)

	m.Members = other.Members
}

func flattenPlacementGroupMember(member linodego.PlacementGroupMember) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		pgMemberObjectType.AttrTypes,
		map[string]attr.Value{
			"linode_id":    types.Int64Value(int64(member.LinodeID)),
			"is_compliant": types.BoolValue(member.IsCompliant),
		},
	)
}
