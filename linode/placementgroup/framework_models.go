package placementgroup

import (
	"context"
	"strconv"

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

type PlacementGroupMemberModel struct {
	LinodeID    types.Int64 `tfsdk:"linode_id"`
	IsCompliant types.Bool  `tfsdk:"is_compliant"`
}

func (m *PlacementGroupMemberModel) FlattenMember(member linodego.PlacementGroupMember) {
	m.LinodeID = types.Int64Value(int64(member.LinodeID))
	m.IsCompliant = types.BoolValue(member.IsCompliant)
}

func (m *PlacementGroupModel) FlattenPlacementGroup(
	ctx context.Context,
	pg *linodego.PlacementGroup,
	preserveKnown bool,
) (resultDiags diag.Diagnostics) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(pg.ID), preserveKnown)

	m.Label = helper.KeepOrUpdateString(m.Label, pg.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, pg.Region, preserveKnown)
	m.AffinityType = helper.KeepOrUpdateString(m.AffinityType, string(pg.AffinityType), preserveKnown)
	m.IsStrict = helper.KeepOrUpdateBool(m.IsStrict, pg.IsStrict, preserveKnown)
	m.IsCompliant = helper.KeepOrUpdateBool(m.IsCompliant, pg.IsCompliant, preserveKnown)

	members := make([]PlacementGroupMemberModel, len(pg.Members))
	for i, member := range pg.Members {
		memberModel := PlacementGroupMemberModel{}
		memberModel.FlattenMember(member)
		members[i] = memberModel
	}

	membersSet, diags := types.SetValueFrom(ctx, pgMemberObjectType, members)
	resultDiags.Append(diags...)
	if resultDiags.HasError() {
		return
	}

	m.Members = membersSet
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
