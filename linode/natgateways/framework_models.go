package natgateways

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v4/linode/natgateway"
)

// NATGatewayFilterModel describes the Terraform data source model to match the
// list data source schema.
type NATGatewayFilterModel struct {
	ID          types.String                           `tfsdk:"id"`
	Filters     frameworkfilter.FiltersModelType       `tfsdk:"filter"`
	Order       types.String                           `tfsdk:"order"`
	OrderBy     types.String                           `tfsdk:"order_by"`
	NATGateways []natgateway.NATGatewayDataSourceModel `tfsdk:"natgateways"`
}

// ParseNATGateways parses the given NAT Gateways into the `natgateways` model attribute.
func (model *NATGatewayFilterModel) ParseNATGateways(gateways []linodego.NATGateway) diag.Diagnostics {
	var diags diag.Diagnostics

	result := make([]natgateway.NATGatewayDataSourceModel, len(gateways))
	for i := range gateways {
		var m natgateway.NATGatewayDataSourceModel
		diags.Append(m.ParseNATGateway(&gateways[i])...)
		result[i] = m
	}
	model.NATGateways = result

	return diags
}
