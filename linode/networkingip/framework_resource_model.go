package networkingip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type NetworkingIPModel struct {
	ID       types.String `tfsdk:"id"`
	LinodeID types.Int64  `tfsdk:"linode_id"`
	Reserved types.Bool   `tfsdk:"reserved"`
	Region   types.String `tfsdk:"region"`
	Public   types.Bool   `tfsdk:"public"`
	Address  types.String `tfsdk:"address"`
	Type     types.String `tfsdk:"type"`
}

func (m *NetworkingIPModel) FlattenIPAddress(ctx context.Context, ip *linodego.InstanceIP, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateString(m.ID, ip.Address, preserveKnown)

	if ip.LinodeID != 0 {
		m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, types.Int64Value(int64(ip.LinodeID)), preserveKnown)
	} else {
		m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, types.Int64Null(), preserveKnown)
	}

	m.Reserved = helper.KeepOrUpdateValue(m.Reserved, types.BoolValue(ip.Reserved), preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, ip.Region, preserveKnown)
	m.Public = helper.KeepOrUpdateValue(m.Public, types.BoolValue(ip.Public), preserveKnown)
	m.Address = helper.KeepOrUpdateString(m.Address, ip.Address, preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, string(ip.Type), preserveKnown)
}
