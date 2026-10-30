package natgatewaysettings

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

// DataSourceModel describes the Terraform data source model for NAT Gateway settings.
type DataSourceModel struct {
	AllowedPortsPerInterface                 []types.Int64 `tfsdk:"allowed_ports_per_interface"`
	MaximumAutoscalingAddressesPerNATGateway types.Int64   `tfsdk:"maximum_autoscaling_addresses_per_natgateway"`
	MaximumReservedAddressesPerNATGateway    types.Int64   `tfsdk:"maximum_reserved_addresses_per_natgateway"`
}

func (model *DataSourceModel) ParseNATGatewaySettings(settings *linodego.NATGatewaySettings) {
	ports := make([]types.Int64, len(settings.AllowedPortsPerInterface))
	for i, port := range settings.AllowedPortsPerInterface {
		ports[i] = types.Int64Value(int64(port))
	}

	model.AllowedPortsPerInterface = ports
	model.MaximumAutoscalingAddressesPerNATGateway = types.Int64Value(int64(settings.MaximumAutoscalingAddressesPerNATGateway))
	model.MaximumReservedAddressesPerNATGateway = types.Int64Value(int64(settings.MaximumReservedAddressesPerNATGateway))
}
