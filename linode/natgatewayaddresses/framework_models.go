package natgatewayaddresses

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

// NATGatewayAddressesModel describes the Terraform data source model to match the
// list data source schema.
type NATGatewayAddressesModel struct {
	ID           types.String                     `tfsdk:"id"`
	NATGatewayID types.Int64                      `tfsdk:"natgateway_id"`
	Filters      frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order        types.String                     `tfsdk:"order"`
	OrderBy      types.String                     `tfsdk:"order_by"`
	Addresses    []NATGatewayAddressItemModel     `tfsdk:"addresses"`
}

// NATGatewayAddressItemModel describes a single entry of the `addresses` attribute.
type NATGatewayAddressItemModel struct {
	Address            types.String `tfsdk:"address"`
	InUse              types.Bool   `tfsdk:"in_use"`
	InterfaceCount     types.Int64  `tfsdk:"interface_count"`
	InterfaceURL       types.String `tfsdk:"interface_url"`
	PortsetAssignments types.Int64  `tfsdk:"portset_assignments"`
	PortsetCapacity    types.Int64  `tfsdk:"portset_capacity"`
}

func (model *NATGatewayAddressesModel) ParseNATGatewayAddresses(addresses []linodego.NATGatewayAddressObject) diag.Diagnostics {
	result := make([]NATGatewayAddressItemModel, len(addresses))
	for i := range addresses {
		addr := &addresses[i]
		result[i] = NATGatewayAddressItemModel{
			Address:            types.StringValue(addr.Address),
			InUse:              types.BoolValue(addr.InUse),
			InterfaceCount:     types.Int64Value(int64(addr.InterfaceCount)),
			InterfaceURL:       types.StringValue(addr.InterfaceURL),
			PortsetAssignments: types.Int64Value(int64(addr.PortsetAssignments)),
			PortsetCapacity:    types.Int64Value(int64(addr.PortsetCapacity)),
		}
	}
	model.Addresses = result

	return nil
}
