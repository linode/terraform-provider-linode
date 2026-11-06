package natgatewayaddress

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

// NATGatewayAddressDataSourceModel describes the Terraform data source model for a single
// NAT Gateway address.
type NATGatewayAddressDataSourceModel struct {
	NATGatewayID       types.Int64  `tfsdk:"natgateway_id"`
	Address            types.String `tfsdk:"address"`
	InUse              types.Bool   `tfsdk:"in_use"`
	InterfaceCount     types.Int64  `tfsdk:"interface_count"`
	InterfaceURL       types.String `tfsdk:"interface_url"`
	PortsetAssignments types.Int64  `tfsdk:"portset_assignments"`
	PortsetCapacity    types.Int64  `tfsdk:"portset_capacity"`
}

func (m *NATGatewayAddressDataSourceModel) ParseNATGatewayAddress(addr *linodego.NATGatewayAddressObject) {
	m.Address = types.StringValue(addr.Address)
	m.InUse = types.BoolValue(addr.InUse)
	m.InterfaceCount = types.Int64Value(int64(addr.InterfaceCount))
	m.InterfaceURL = types.StringValue(addr.InterfaceURL)
	m.PortsetAssignments = types.Int64Value(int64(addr.PortsetAssignments))
	m.PortsetCapacity = types.Int64Value(int64(addr.PortsetCapacity))
}
