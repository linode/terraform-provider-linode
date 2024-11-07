package networkipassignment

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type NetworkingIPModel struct {
	ID          types.String      `tfsdk:"id"`
	LinodeID    types.Int64       `tfsdk:"linode_id"`
	Reserved    types.Bool        `tfsdk:"reserved"`
	Region      types.String      `tfsdk:"region"`
	Public      types.Bool        `tfsdk:"public"`
	Address     types.String      `tfsdk:"address"`
	Type        types.String      `tfsdk:"type"`
	Assignments []AssignmentModel `tfsdk:"assignments"`
}

type AssignmentModel struct {
	Address  types.String `tfsdk:"address"`
	LinodeID types.Int64  `tfsdk:"linode_id"`
}

func (m *NetworkingIPModel) FlattenIPAddress(ip *linodego.InstanceIP) {
	m.ID = types.StringValue(ip.Address)
	if ip.LinodeID != 0 {
		m.LinodeID = types.Int64Value(int64(ip.LinodeID))
	} else {
		m.LinodeID = types.Int64Null()
	}
	m.Reserved = types.BoolValue(ip.Reserved)
	m.Region = types.StringValue(ip.Region)
	m.Public = types.BoolValue(ip.Public)
	m.Address = types.StringValue(ip.Address)
	m.Type = types.StringValue(string(ip.Type))
}
