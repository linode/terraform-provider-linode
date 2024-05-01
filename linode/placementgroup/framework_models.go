package placementgroup

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type PlacementGroupModel struct {
	ID           types.Int64            `tfsdk:"id"`
	Label        types.String           `tfsdk:"label"`
	Region       types.String           `tfsdk:"region"`
	AffinityType types.String           `tfsdk:"affinity_type"`
	IsCompliant  types.Bool             `tfsdk:"is_compliant"`
	IsStrict     types.Bool             `tfsdk:"is_strict"`
	Members      []PlacementGroupMember `tfsdk:"members"`
}

type PlacementGroupMember struct {
	LinodeID    types.Int64 `tfsdk:"linode_id"`
	IsCompliant types.Bool  `tfsdk:"is_compliant"`
}

func (data *PlacementGroupModel) parsePlacementGroup(
	pg *linodego.PlacementGroup,
) {
	data.Label = types.StringValue(pg.Label)
	data.Region = types.StringValue(pg.Region)
	data.AffinityType = types.StringValue(string(pg.AffinityType))
	data.IsCompliant = types.BoolValue(pg.IsCompliant)
	data.IsStrict = types.BoolValue(pg.IsStrict)

	members := make([]PlacementGroupMember, len(pg.Members))

	for i, member := range pg.Members {
		var m PlacementGroupMember
		m.LinodeID = types.Int64Value(int64(member.LinodeID))
		m.IsCompliant = types.BoolValue(member.IsCompliant)

		members[i] = m
	}

	data.Members = members
}
