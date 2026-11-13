package natgatewayinterfaces

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

// Model describes the Terraform data source model for linode_natgateway_interfaces.
type Model struct {
	ID           types.String     `tfsdk:"id"`
	NATGatewayID types.Int64      `tfsdk:"natgateway_id"`
	Address      types.String     `tfsdk:"address"`
	Interfaces   []interfaceModel `tfsdk:"interfaces"`
}

type interfaceModel struct {
	ID        types.Int64    `tfsdk:"id"`
	Linode    linodeModel    `tfsdk:"linode"`
	Addresses []string       `tfsdk:"addresses"`
	Portsets  []portsetModel `tfsdk:"portsets"`
}

type linodeModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
	Type  types.String `tfsdk:"type"`
	URL   types.String `tfsdk:"url"`
}

type portsetModel struct {
	Address types.String       `tfsdk:"address"`
	Ports   []portsetPortModel `tfsdk:"ports"`
}

type portsetPortModel struct {
	Start types.Int64 `tfsdk:"start"`
	End   types.Int64 `tfsdk:"end"`
}

// ParseInterfaces maps the API response to the `interfaces` model attribute.
func (m *Model) ParseInterfaces(interfaces []linodego.NATGatewayInterface) {
	result := make([]interfaceModel, len(interfaces))
	for i, iface := range interfaces {
		result[i] = parseInterface(iface)
	}
	m.Interfaces = result
}

func parseInterface(iface linodego.NATGatewayInterface) interfaceModel {
	addresses := make([]string, len(iface.Addresses))
	copy(addresses, iface.Addresses)

	portsets := make([]portsetModel, len(iface.Portsets))
	for i, ps := range iface.Portsets {
		ports := make([]portsetPortModel, len(ps.Ports))
		for j, p := range ps.Ports {
			ports[j] = portsetPortModel{
				Start: types.Int64Value(int64(p.Start)),
				End:   types.Int64Value(int64(p.End)),
			}
		}
		portsets[i] = portsetModel{
			Address: types.StringValue(ps.Address),
			Ports:   ports,
		}
	}

	return interfaceModel{
		ID: types.Int64Value(int64(iface.ID)),
		Linode: linodeModel{
			ID:    types.Int64Value(int64(iface.Linode.ID)),
			Label: types.StringValue(iface.Linode.Label),
			Type:  types.StringValue(iface.Linode.Type),
			URL:   types.StringValue(iface.Linode.URL),
		},
		Addresses: addresses,
		Portsets:  portsets,
	}
}
