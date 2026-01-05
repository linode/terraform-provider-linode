package linodeinterface

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
)

// DataSourceModel is the data source-specific model for linode_interface
type DataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	LinodeID     types.Int64  `tfsdk:"linode_id"`
	DefaultRoute types.Object `tfsdk:"default_route"`
	Public       types.Object `tfsdk:"public"`
	VLAN         types.Object `tfsdk:"vlan"`
	VPC          types.Object `tfsdk:"vpc"`
}

// DataSourceDefaultRouteModel is the data source model for default route
type DataSourceDefaultRouteModel struct {
	IPv4 types.Bool `tfsdk:"ipv4"`
	IPv6 types.Bool `tfsdk:"ipv6"`
}

// DataSourcePublicModel is the data source model for public interface
type DataSourcePublicModel struct {
	IPv4 types.Object `tfsdk:"ipv4"`
	IPv6 types.Object `tfsdk:"ipv6"`
}

// DataSourcePublicIPv4Model is the data source model for public IPv4
type DataSourcePublicIPv4Model struct {
	Addresses types.Set `tfsdk:"addresses"`
	Shared    types.Set `tfsdk:"shared"`
}

// DataSourcePublicIPv6Model is the data source model for public IPv6
type DataSourcePublicIPv6Model struct {
	Ranges types.Set `tfsdk:"ranges"`
	Shared types.Set `tfsdk:"shared"`
	SLAAC  types.Set `tfsdk:"slaac"`
}

// DataSourcePublicIPv4AddressModel represents a public IPv4 address
type DataSourcePublicIPv4AddressModel struct {
	Address types.String `tfsdk:"address"`
	Primary types.Bool   `tfsdk:"primary"`
}

// DataSourceSharedPublicIPv4AddressModel represents a shared public IPv4 address
type DataSourceSharedPublicIPv4AddressModel struct {
	Address  types.String `tfsdk:"address"`
	LinodeID types.Int64  `tfsdk:"linode_id"`
}

// DataSourcePublicIPv6RangeModel represents a public IPv6 range
type DataSourcePublicIPv6RangeModel struct {
	Range       cidrtypes.IPv6Prefix `tfsdk:"range"`
	RouteTarget types.String         `tfsdk:"route_target"`
}

// DataSourcePublicIPv6SLAACModel represents IPv6 SLAAC configuration
type DataSourcePublicIPv6SLAACModel struct {
	Address iptypes.IPv6Address `tfsdk:"address"`
	Prefix  types.Int64         `tfsdk:"prefix"`
}

// DataSourceVLANModel is the data source model for VLAN interface
type DataSourceVLANModel struct {
	VLANLabel   types.String         `tfsdk:"vlan_label"`
	IPAMAddress cidrtypes.IPv4Prefix `tfsdk:"ipam_address"`
}

// DataSourceVPCModel is the data source model for VPC interface
type DataSourceVPCModel struct {
	SubnetID types.Int64  `tfsdk:"subnet_id"`
	IPv4     types.Object `tfsdk:"ipv4"`
	IPv6     types.Object `tfsdk:"ipv6"`
}

// DataSourceVPCIPv4Model is the data source model for VPC IPv4
type DataSourceVPCIPv4Model struct {
	Addresses types.Set `tfsdk:"addresses"`
	Ranges    types.Set `tfsdk:"ranges"`
}

// DataSourceVPCIPv6Model is the data source model for VPC IPv6
type DataSourceVPCIPv6Model struct {
	IsPublic types.Bool `tfsdk:"is_public"`
	SLAAC    types.Set  `tfsdk:"slaac"`
	Ranges   types.Set  `tfsdk:"ranges"`
}

// DataSourceVPCIPv4AddressModel represents a VPC IPv4 address
type DataSourceVPCIPv4AddressModel struct {
	Address      iptypes.IPv4Address `tfsdk:"address"`
	Primary      types.Bool          `tfsdk:"primary"`
	NAT11Address iptypes.IPv4Address `tfsdk:"nat_1_1_address"`
}

// DataSourceVPCIPv4RangeModel represents a VPC IPv4 range
type DataSourceVPCIPv4RangeModel struct {
	Range cidrtypes.IPv4Prefix `tfsdk:"range"`
}

// DataSourceVPCIPv6SLAACModel represents a VPC IPv6 SLAAC configuration
type DataSourceVPCIPv6SLAACModel struct {
	Range   cidrtypes.IPv6Prefix `tfsdk:"range"`
	Address iptypes.IPv6Address  `tfsdk:"address"`
}

// DataSourceVPCIPv6RangeModel represents a VPC IPv6 range
type DataSourceVPCIPv6RangeModel struct {
	Range cidrtypes.IPv6Prefix `tfsdk:"range"`
}

// FlattenInterface flattens a linodego.LinodeInterface into the data source model
func (data *DataSourceModel) FlattenInterface(
	ctx context.Context,
	i linodego.LinodeInterface,
	diags *diag.Diagnostics,
) {
	tflog.Trace(ctx, "Enter DataSourceModel.FlattenInterface")

	data.ID = types.StringValue(strconv.Itoa(i.ID))

	// Flatten default route
	if i.DefaultRoute != nil {
		defaultRoute := DataSourceDefaultRouteModel{
			IPv4: types.BoolPointerValue(i.DefaultRoute.IPv4),
			IPv6: types.BoolPointerValue(i.DefaultRoute.IPv6),
		}
		defaultRouteObj, defaultRouteDiags := types.ObjectValueFrom(ctx, dataSourceDefaultRouteAttribute.GetType().(types.ObjectType).AttrTypes, defaultRoute)
		diags.Append(defaultRouteDiags...)
		if diags.HasError() {
			return
		}
		data.DefaultRoute = defaultRouteObj
	}

	// Flatten VLAN
	if i.VLAN != nil {
		vlan := DataSourceVLANModel{
			VLANLabel:   types.StringValue(i.VLAN.VLANLabel),
			IPAMAddress: cidrtypes.NewIPv4PrefixPointerValue(i.VLAN.IPAMAddress),
		}
		vlanObj, vlanDiags := types.ObjectValueFrom(ctx, dataSourceVLANAttribute.GetType().(types.ObjectType).AttrTypes, vlan)
		diags.Append(vlanDiags...)
		if diags.HasError() {
			return
		}
		data.VLAN = vlanObj
	}

	// Flatten Public interface
	if i.Public != nil {
		data.flattenPublicInterface(ctx, *i.Public, diags)
		if diags.HasError() {
			return
		}
	}

	// Flatten VPC interface
	if i.VPC != nil {
		data.flattenVPCInterface(ctx, *i.VPC, diags)
		if diags.HasError() {
			return
		}
	}
}

func (data *DataSourceModel) flattenPublicInterface(
	ctx context.Context,
	publicInterface linodego.PublicInterface,
	diags *diag.Diagnostics,
) {
	tflog.Trace(ctx, "Enter DataSourceModel.flattenPublicInterface")

	var publicIPv4 DataSourcePublicIPv4Model
	var publicIPv6 DataSourcePublicIPv6Model

	// Flatten IPv4
	if publicInterface.IPv4 != nil {
		addresses := make([]DataSourcePublicIPv4AddressModel, len(publicInterface.IPv4.Addresses))
		for i, addr := range publicInterface.IPv4.Addresses {
			addresses[i] = DataSourcePublicIPv4AddressModel{
				Address: types.StringValue(addr.Address),
				Primary: types.BoolValue(addr.Primary),
			}
		}

		addressesSet, addressesDiags := types.SetValueFrom(ctx, dataSourcePublicIPv4AddressAttribute.Type(), addresses)
		diags.Append(addressesDiags...)
		if diags.HasError() {
			return
		}
		publicIPv4.Addresses = addressesSet

		shared := make([]DataSourceSharedPublicIPv4AddressModel, len(publicInterface.IPv4.Shared))
		for i, addr := range publicInterface.IPv4.Shared {
			shared[i] = DataSourceSharedPublicIPv4AddressModel{
				Address:  types.StringValue(addr.Address),
				LinodeID: types.Int64Value(int64(addr.LinodeID)),
			}
		}

		sharedSet, sharedDiags := types.SetValueFrom(ctx, dataSourceSharedPublicIPv4AddressAttribute.Type(), shared)
		diags.Append(sharedDiags...)
		if diags.HasError() {
			return
		}
		publicIPv4.Shared = sharedSet
	}

	// Flatten IPv6
	if publicInterface.IPv6 != nil {
		ranges := make([]DataSourcePublicIPv6RangeModel, len(publicInterface.IPv6.Ranges))
		for i, r := range publicInterface.IPv6.Ranges {
			ranges[i] = DataSourcePublicIPv6RangeModel{
				Range:       cidrtypes.NewIPv6PrefixValue(r.Range),
				RouteTarget: types.StringPointerValue(r.RouteTarget),
			}
		}

		rangesSet, rangesDiags := types.SetValueFrom(ctx, dataSourcePublicIPv6RangeAttribute.Type(), ranges)
		diags.Append(rangesDiags...)
		if diags.HasError() {
			return
		}
		publicIPv6.Ranges = rangesSet

		shared := make([]DataSourcePublicIPv6RangeModel, len(publicInterface.IPv6.Shared))
		for i, r := range publicInterface.IPv6.Shared {
			shared[i] = DataSourcePublicIPv6RangeModel{
				Range:       cidrtypes.NewIPv6PrefixValue(r.Range),
				RouteTarget: types.StringPointerValue(r.RouteTarget),
			}
		}

		sharedSet, sharedDiags := types.SetValueFrom(ctx, dataSourcePublicIPv6RangeAttribute.Type(), shared)
		diags.Append(sharedDiags...)
		if diags.HasError() {
			return
		}
		publicIPv6.Shared = sharedSet

		slaac := make([]DataSourcePublicIPv6SLAACModel, len(publicInterface.IPv6.SLAAC))
		for i, s := range publicInterface.IPv6.SLAAC {
			slaac[i] = DataSourcePublicIPv6SLAACModel{
				Address: iptypes.NewIPv6AddressValue(s.Address),
				Prefix:  types.Int64Value(int64(s.Prefix)),
			}
		}

		slaacSet, slaacDiags := types.SetValueFrom(ctx, dataSourcePublicIPv6SLAACAttribute.Type(), slaac)
		diags.Append(slaacDiags...)
		if diags.HasError() {
			return
		}
		publicIPv6.SLAAC = slaacSet
	}

	// Convert IPv4 and IPv6 models to objects
	ipv4Obj, ipv4Diags := types.ObjectValueFrom(ctx, dataSourcePublicIPv4Attribute.GetType().(types.ObjectType).AttrTypes, publicIPv4)
	diags.Append(ipv4Diags...)
	if diags.HasError() {
		return
	}

	ipv6Obj, ipv6Diags := types.ObjectValueFrom(ctx, dataSourcePublicIPv6Attribute.GetType().(types.ObjectType).AttrTypes, publicIPv6)
	diags.Append(ipv6Diags...)
	if diags.HasError() {
		return
	}

	publicModel := DataSourcePublicModel{
		IPv4: ipv4Obj,
		IPv6: ipv6Obj,
	}

	publicObj, publicDiags := types.ObjectValueFrom(ctx, dataSourcePublicAttribute.GetType().(types.ObjectType).AttrTypes, publicModel)
	diags.Append(publicDiags...)
	if diags.HasError() {
		return
	}

	data.Public = publicObj
}

func (data *DataSourceModel) flattenVPCInterface(
	ctx context.Context,
	vpcInterface linodego.VPCInterface,
	diags *diag.Diagnostics,
) {
	tflog.Trace(ctx, "Enter DataSourceModel.flattenVPCInterface")

	var vpcIPv4 DataSourceVPCIPv4Model
	var vpcIPv6 DataSourceVPCIPv6Model

	// Flatten IPv4 addresses
	addresses := make([]DataSourceVPCIPv4AddressModel, len(vpcInterface.IPv4.Addresses))
	for i, addr := range vpcInterface.IPv4.Addresses {
		addresses[i] = DataSourceVPCIPv4AddressModel{
			Address:      iptypes.NewIPv4AddressValue(addr.Address),
			Primary:      types.BoolValue(addr.Primary),
			NAT11Address: iptypes.NewIPv4AddressPointerValue(addr.NAT1To1Address),
		}
	}

	addressesSet, addressesDiags := types.SetValueFrom(ctx, dataSourceVPCIPv4AddressAttribute.Type(), addresses)
	diags.Append(addressesDiags...)
	if diags.HasError() {
		return
	}
	vpcIPv4.Addresses = addressesSet

	// Flatten IPv4 ranges
	ranges := make([]DataSourceVPCIPv4RangeModel, len(vpcInterface.IPv4.Ranges))
	for i, r := range vpcInterface.IPv4.Ranges {
		ranges[i] = DataSourceVPCIPv4RangeModel{
			Range: cidrtypes.NewIPv4PrefixValue(r.Range),
		}
	}

	rangesSet, rangesDiags := types.SetValueFrom(ctx, dataSourceVPCIPv4RangeAttribute.Type(), ranges)
	diags.Append(rangesDiags...)
	if diags.HasError() {
		return
	}
	vpcIPv4.Ranges = rangesSet

	// Convert IPv4 model to object
	ipv4Obj, ipv4Diags := types.ObjectValueFrom(ctx, dataSourceVPCIPv4Attribute.GetType().(types.ObjectType).AttrTypes, vpcIPv4)
	diags.Append(ipv4Diags...)
	if diags.HasError() {
		return
	}

	// Flatten IPv6
	vpcIPv6.IsPublic = types.BoolPointerValue(vpcInterface.IPv6.IsPublic)

	// Flatten IPv6 SLAAC
	slaac := make([]DataSourceVPCIPv6SLAACModel, len(vpcInterface.IPv6.SLAAC))
	for i, s := range vpcInterface.IPv6.SLAAC {
		slaac[i] = DataSourceVPCIPv6SLAACModel{
			Range:   cidrtypes.NewIPv6PrefixValue(s.Range),
			Address: iptypes.NewIPv6AddressValue(s.Address),
		}
	}

	slaacSet, slaacDiags := types.SetValueFrom(ctx, dataSourceVPCIPv6SLAACAttribute.Type(), slaac)
	diags.Append(slaacDiags...)
	if diags.HasError() {
		return
	}
	vpcIPv6.SLAAC = slaacSet

	// Flatten IPv6 ranges
	ipv6Ranges := make([]DataSourceVPCIPv6RangeModel, len(vpcInterface.IPv6.Ranges))
	for i, r := range vpcInterface.IPv6.Ranges {
		ipv6Ranges[i] = DataSourceVPCIPv6RangeModel{
			Range: cidrtypes.NewIPv6PrefixValue(r.Range),
		}
	}

	ipv6RangesSet, ipv6RangesDiags := types.SetValueFrom(ctx, dataSourceVPCIPv6RangeAttribute.Type(), ipv6Ranges)
	diags.Append(ipv6RangesDiags...)
	if diags.HasError() {
		return
	}
	vpcIPv6.Ranges = ipv6RangesSet

	// Convert IPv6 model to object
	ipv6Obj, ipv6Diags := types.ObjectValueFrom(ctx, dataSourceVPCIPv6Attribute.GetType().(types.ObjectType).AttrTypes, vpcIPv6)
	diags.Append(ipv6Diags...)
	if diags.HasError() {
		return
	}

	vpcModel := DataSourceVPCModel{
		SubnetID: types.Int64Value(int64(vpcInterface.SubnetID)),
		IPv4:     ipv4Obj,
		IPv6:     ipv6Obj,
	}

	vpcObj, vpcDiags := types.ObjectValueFrom(ctx, dataSourceVPCAttribute.GetType().(types.ObjectType).AttrTypes, vpcModel)
	diags.Append(vpcDiags...)
	if diags.HasError() {
		return
	}

	data.VPC = vpcObj
}
