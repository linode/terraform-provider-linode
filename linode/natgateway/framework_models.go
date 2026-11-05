package natgateway

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

var (
	natGatewayAddressAttrTypes = map[string]attr.Type{
		"address": types.StringType,
	}
	natGatewayVPCSubnetAttrTypes = map[string]attr.Type{
		"id":        types.StringType,
		"type":      types.StringType,
		"url":       types.StringType,
		"label":     types.StringType,
		"vpc_id":    types.StringType,
		"vpc_label": types.StringType,
	}
)

type NATGatewayModel struct {
	ID                       types.Int64       `tfsdk:"id"`
	Region                   types.String      `tfsdk:"region"`
	Label                    types.String      `tfsdk:"label"`
	DefaultPortsPerInterface types.Int64       `tfsdk:"default_ports_per_interface"`
	UseAutoscaling           types.Bool        `tfsdk:"use_autoscaling"`
	AddressAutoscaleMax      types.Int64       `tfsdk:"address_autoscale_max"`
	PortsetAssignments       types.Int64       `tfsdk:"portset_assignments"`
	PortsetCapacity          types.Int64       `tfsdk:"portset_capacity"`
	VPCSubnet                types.Object      `tfsdk:"vpc_subnet"`
	Created                  timetypes.RFC3339 `tfsdk:"created"`
	Updated                  timetypes.RFC3339 `tfsdk:"updated"`
	Addresses                types.Set         `tfsdk:"addresses"`
}

func (m *NATGatewayModel) ParseNATGateway(gw *linodego.NATGateway) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.Int64Value(int64(gw.ID))
	m.Region = types.StringValue(gw.Region)
	m.Label = types.StringValue(gw.Label)
	m.DefaultPortsPerInterface = types.Int64Value(int64(gw.DefaultPortsPerInterface))
	m.AddressAutoscaleMax = types.Int64Value(int64(gw.AddressAutoscaleMax))
	m.PortsetAssignments = types.Int64Value(int64(gw.PortsetAssignments))
	m.PortsetCapacity = types.Int64Value(int64(gw.PortsetCapacity))
	m.Created = timetypes.NewRFC3339TimePointerValue(gw.Created)
	m.Updated = timetypes.NewRFC3339TimePointerValue(gw.Updated)

	if gw.VPCSubnet.ID == 0 {
		m.VPCSubnet = types.ObjectNull(natGatewayVPCSubnetAttrTypes)
	} else {
		subnet, d := types.ObjectValue(natGatewayVPCSubnetAttrTypes, map[string]attr.Value{
			"id":        types.StringValue(strconv.Itoa(gw.VPCSubnet.ID)),
			"type":      types.StringValue(gw.VPCSubnet.Type),
			"url":       types.StringValue(gw.VPCSubnet.URL),
			"label":     types.StringValue(gw.VPCSubnet.Label),
			"vpc_id":    types.StringValue(strconv.Itoa(gw.VPCSubnet.VPCID)),
			"vpc_label": types.StringValue(gw.VPCSubnet.VPCLabel),
		})
		diags.Append(d...)
		m.VPCSubnet = subnet
	}

	addresses := make([]attr.Value, 0, len(gw.Addresses))
	for _, a := range gw.Addresses {
		obj, d := types.ObjectValue(natGatewayAddressAttrTypes, map[string]attr.Value{
			"address": types.StringValue(a.Address),
		})
		diags.Append(d...)
		addresses = append(addresses, obj)
	}
	addressSet, d := types.SetValue(types.ObjectType{AttrTypes: natGatewayAddressAttrTypes}, addresses)
	diags.Append(d...)
	m.Addresses = addressSet

	return diags
}

// natGatewayAddressModel is the Terraform representation of a single entry
// in the addresses set.
type natGatewayAddressModel struct {
	Address types.String `tfsdk:"address"`
}

// ExpandCreateOptions builds the NAT Gateway create request from the model.
// Optional fields are only sent when they are known and not null, so the API
// applies its own defaults otherwise. addresses is always sent (as [] when
// empty), since the API requires it to be empty when use_autoscaling is true.
func (m *NATGatewayModel) ExpandCreateOptions(ctx context.Context) (linodego.NATGatewayCreateOptions, diag.Diagnostics) {
	var diags diag.Diagnostics

	opts := linodego.NATGatewayCreateOptions{
		Region:    m.Region.ValueString(),
		Label:     m.Label.ValueString(),
		Addresses: []linodego.NATGatewayAddressCreateOptions{},
	}

	if !m.Addresses.IsNull() && !m.Addresses.IsUnknown() {
		var addresses []natGatewayAddressModel
		diags.Append(m.Addresses.ElementsAs(ctx, &addresses, false)...)
		for _, a := range addresses {
			opts.Addresses = append(opts.Addresses, linodego.NATGatewayAddressCreateOptions{
				Address: a.Address.ValueString(),
			})
		}
	}

	if !m.DefaultPortsPerInterface.IsNull() && !m.DefaultPortsPerInterface.IsUnknown() {
		v := int(m.DefaultPortsPerInterface.ValueInt64())
		opts.DefaultPortsPerInterface = &v
	}

	if !m.UseAutoscaling.IsNull() && !m.UseAutoscaling.IsUnknown() {
		v := m.UseAutoscaling.ValueBool()
		opts.UseAutoscaling = &v
	}

	return opts, diags
}

// NATGatewayDataSourceModel is the state model for the linode_natgateway data source.
// It omits the create-only field (use_autoscaling), which the data
// source schema does not define.
type NATGatewayDataSourceModel struct {
	ID                       types.Int64       `tfsdk:"id"`
	Region                   types.String      `tfsdk:"region"`
	Label                    types.String      `tfsdk:"label"`
	DefaultPortsPerInterface types.Int64       `tfsdk:"default_ports_per_interface"`
	AddressAutoscaleMax      types.Int64       `tfsdk:"address_autoscale_max"`
	PortsetAssignments       types.Int64       `tfsdk:"portset_assignments"`
	PortsetCapacity          types.Int64       `tfsdk:"portset_capacity"`
	VPCSubnet                types.Object      `tfsdk:"vpc_subnet"`
	Created                  timetypes.RFC3339 `tfsdk:"created"`
	Updated                  timetypes.RFC3339 `tfsdk:"updated"`
	Addresses                types.Set         `tfsdk:"addresses"`
}

// ParseNATGateway maps the API response to the data source model by reusing the
// resource parse and copying the shared fields.
func (m *NATGatewayDataSourceModel) ParseNATGateway(gw *linodego.NATGateway) diag.Diagnostics {
	var full NATGatewayModel
	diags := full.ParseNATGateway(gw)

	m.ID = full.ID
	m.Region = full.Region
	m.Label = full.Label
	m.DefaultPortsPerInterface = full.DefaultPortsPerInterface
	m.AddressAutoscaleMax = full.AddressAutoscaleMax
	m.PortsetAssignments = full.PortsetAssignments
	m.PortsetCapacity = full.PortsetCapacity
	m.VPCSubnet = full.VPCSubnet
	m.Created = full.Created
	m.Updated = full.Updated
	m.Addresses = full.Addresses

	return diags
}
