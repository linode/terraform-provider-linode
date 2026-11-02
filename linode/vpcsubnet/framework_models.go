package vpcsubnet

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/customtypes"
)

/*
Shared Implementation
*/

type BaseModel struct {
	ID            types.String      `tfsdk:"id"`
	VPCId         types.Int64       `tfsdk:"vpc_id"`
	VPCType       types.String      `tfsdk:"vpc_type"`
	Label         types.String      `tfsdk:"label"`
	IPv4          types.String      `tfsdk:"ipv4"`
	Linodes       types.List        `tfsdk:"linodes"`
	Created       timetypes.RFC3339 `tfsdk:"created"`
	Updated       timetypes.RFC3339 `tfsdk:"updated"`
	Databases     types.List        `tfsdk:"databases"`
	Nodebalancers types.List        `tfsdk:"nodebalancers"`
}

type SubnetDatabaseIPv6RangeModel struct {
	Range types.String `tfsdk:"range"`
}

type SubnetDatabaseModel struct {
	ID         types.Int64  `tfsdk:"id"`
	IPv4Range  types.String `tfsdk:"ipv4_range"`
	IPv6Ranges types.List   `tfsdk:"ipv6_ranges"`
}

type SubnetNodebalancerIPv6RangeModel struct {
	Range types.String `tfsdk:"range"`
}

type SubnetNodebalancerModel struct {
	ID         types.Int64  `tfsdk:"id"`
	IPv4Range  types.String `tfsdk:"ipv4_range"`
	IPv6Ranges types.List   `tfsdk:"ipv6_ranges"`
}

func FlattenSubnetLinodeInterface(iface linodego.VPCSubnetLinodeInterface) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(LinodeInterfaceObjectType.AttrTypes, map[string]attr.Value{
		"id":        types.Int64Value(int64(iface.ID)),
		"config_id": types.Int64PointerValue(helper.IntPtrToInt64Ptr(iface.ConfigID)),
		"active":    types.BoolValue(iface.Active),
	})
}

func flattenSubnetDatabaseIPv6Range(
	ipv6Range string,
) SubnetDatabaseIPv6RangeModel {
	return SubnetDatabaseIPv6RangeModel{
		Range: types.StringValue(ipv6Range),
	}
}

func FlattenSubnetLinode(
	ctx context.Context,
	linode linodego.VPCSubnetLinode,
) (*types.Object, diag.Diagnostics) {
	result := map[string]attr.Value{
		"id": types.Int64Value(int64(linode.ID)),
	}

	ifaces := make([]types.Object, len(linode.Interfaces))

	for i, iface := range linode.Interfaces {
		ifaceObj, d := FlattenSubnetLinodeInterface(iface)
		if d.HasError() {
			return nil, d
		}

		ifaces[i] = ifaceObj
	}

	ifacesList, d := types.ListValueFrom(ctx, LinodeInterfaceObjectType, ifaces)
	if d.HasError() {
		return nil, d
	}

	result["interfaces"] = ifacesList

	resultObject, d := types.ObjectValue(LinodeObjectType.AttrTypes, result)
	return &resultObject, d
}

func FlattenSubnetLinodes(
	ctx context.Context,
	subnetLinodes []linodego.VPCSubnetLinode,
) (*types.List, diag.Diagnostics) {
	result := make([]types.Object, len(subnetLinodes))

	for i, inst := range subnetLinodes {
		linodeObj, diags := FlattenSubnetLinode(ctx, inst)
		if diags.HasError() {
			return nil, diags
		}
		result[i] = *linodeObj
	}

	linodesList, diags := types.ListValueFrom(ctx, LinodeObjectType, result)
	return &linodesList, diags
}

func FlattenSubnetDatabase(
	ctx context.Context,
	db linodego.VPCSubnetDatabase,
) (*types.Object, diag.Diagnostics) {
	ipv6RangeModels := helper.MapSlice(
		db.IPv6Ranges,
		func(r string) SubnetDatabaseIPv6RangeModel {
			return flattenSubnetDatabaseIPv6Range(r)
		},
	)

	ipv6RangesList, diags := types.ListValueFrom(
		ctx,
		IPV6RangeObjectType,
		ipv6RangeModels,
	)
	if diags.HasError() {
		return nil, diags
	}

	model := SubnetDatabaseModel{
		ID:         types.Int64Value(int64(db.ID)),
		IPv4Range:  types.StringPointerValue(db.IPv4Range),
		IPv6Ranges: ipv6RangesList,
	}

	obj, diags := types.ObjectValueFrom(
		ctx,
		DatabaseObjectType.AttrTypes,
		model,
	)

	return &obj, diags
}

func FlattenSubnetDatabases(
	ctx context.Context,
	subnetDatabases []linodego.VPCSubnetDatabase,
) (*types.List, diag.Diagnostics) {
	result := make([]types.Object, len(subnetDatabases))

	for i, db := range subnetDatabases {
		dbObj, diags := FlattenSubnetDatabase(ctx, db)
		if diags.HasError() {
			return nil, diags
		}
		result[i] = *dbObj
	}

	dbsList, diags := types.ListValueFrom(ctx, DatabaseObjectType, result)
	return &dbsList, diags
}

func flattenSubnetNodebalancerIPv6Range(
	ipv6Range string,
) SubnetNodebalancerIPv6RangeModel {
	return SubnetNodebalancerIPv6RangeModel{
		Range: types.StringValue(ipv6Range),
	}
}

func FlattenSubnetNodebalancer(
	ctx context.Context,
	nb linodego.VPCSubnetNodebalancers,
) (*types.Object, diag.Diagnostics) {
	ipv6RangeModels := helper.MapSlice(
		nb.Ipv6Ranges,
		func(r linodego.VPCSubnetNodebalancersRanges) SubnetNodebalancerIPv6RangeModel {
			return flattenSubnetNodebalancerIPv6Range(r.Range)
		},
	)

	ipv6RangesList, diags := types.ListValueFrom(
		ctx,
		IPV6RangeObjectType,
		ipv6RangeModels,
	)
	if diags.HasError() {
		return nil, diags
	}

	model := SubnetNodebalancerModel{
		ID:         types.Int64Value(int64(nb.ID)),
		IPv4Range:  types.StringValue(nb.Ipv4Range),
		IPv6Ranges: ipv6RangesList,
	}

	obj, diags := types.ObjectValueFrom(
		ctx,
		NodebalancerObjectType.AttrTypes,
		model,
	)

	return &obj, diags
}

func FlattenSubnetNodebalancers(
	ctx context.Context,
	subnetNodebalancers []linodego.VPCSubnetNodebalancers,
) (*types.List, diag.Diagnostics) {
	result := make([]types.Object, len(subnetNodebalancers))

	for i, nb := range subnetNodebalancers {
		nbObj, diags := FlattenSubnetNodebalancer(ctx, nb)
		if diags.HasError() {
			return nil, diags
		}
		result[i] = *nbObj
	}

	nbsList, diags := types.ListValueFrom(ctx, NodebalancerObjectType, result)
	return &nbsList, diags
}

var NATGatewayObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id": types.Int64Type,
	},
}

var NATGatewayPortsetPortObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"start": types.Int64Type,
		"end":   types.Int64Type,
	},
}

var NATGatewayPortsetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"address": types.StringType,
		"ports": types.ListType{
			ElemType: NATGatewayPortsetPortObjectType,
		},
	},
}

var NATGatewayDetailsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                  types.Int64Type,
		"label":               types.StringType,
		"addresses":           types.ListType{ElemType: types.StringType},
		"portset_assignments": types.Int64Type,
		"portset_capacity":    types.Int64Type,
		"portsets": types.ListType{
			ElemType: NATGatewayPortsetObjectType,
		},
	},
}

// flattenNATGatewaySetting converts a subnet's NAT Gateway into the settable natgateway block.
func flattenNATGatewaySetting(gw *linodego.VPCSubnetNATGateway) (types.Object, diag.Diagnostics) {
	if gw == nil {
		return types.ObjectNull(NATGatewayObjectType.AttrTypes), nil
	}

	return types.ObjectValue(NATGatewayObjectType.AttrTypes, map[string]attr.Value{
		"id": types.Int64Value(int64(gw.ID)),
	})
}

// flattenNATGatewayDetails converts a subnet's NAT Gateway into the read-only natgateway details block.
func flattenNATGatewayDetails(
	ctx context.Context,
	gw *linodego.VPCSubnetNATGateway,
) (types.Object, diag.Diagnostics) {
	if gw == nil {
		return types.ObjectNull(NATGatewayDetailsObjectType.AttrTypes), nil
	}

	addresses, diags := types.ListValueFrom(ctx, types.StringType, gw.Addresses)
	if diags.HasError() {
		return types.ObjectNull(NATGatewayDetailsObjectType.AttrTypes), diags
	}

	portsets := types.ListNull(NATGatewayPortsetObjectType)
	if gw.Portsets != nil {
		portsetObjs := make([]types.Object, len(gw.Portsets))
		for i, ps := range gw.Portsets {
			portObjs := make([]types.Object, len(ps.Ports))
			for j, p := range ps.Ports {
				portObj, d := types.ObjectValue(NATGatewayPortsetPortObjectType.AttrTypes, map[string]attr.Value{
					"start": types.Int64Value(int64(p.Start)),
					"end":   types.Int64Value(int64(p.End)),
				})
				if d.HasError() {
					return types.ObjectNull(NATGatewayDetailsObjectType.AttrTypes), d
				}
				portObjs[j] = portObj
			}

			portsList, d := types.ListValueFrom(ctx, NATGatewayPortsetPortObjectType, portObjs)
			if d.HasError() {
				return types.ObjectNull(NATGatewayDetailsObjectType.AttrTypes), d
			}

			portsetObj, d := types.ObjectValue(NATGatewayPortsetObjectType.AttrTypes, map[string]attr.Value{
				"address": types.StringValue(ps.Address),
				"ports":   portsList,
			})
			if d.HasError() {
				return types.ObjectNull(NATGatewayDetailsObjectType.AttrTypes), d
			}
			portsetObjs[i] = portsetObj
		}

		portsets, diags = types.ListValueFrom(ctx, NATGatewayPortsetObjectType, portsetObjs)
		if diags.HasError() {
			return types.ObjectNull(NATGatewayDetailsObjectType.AttrTypes), diags
		}
	}

	return types.ObjectValue(NATGatewayDetailsObjectType.AttrTypes, map[string]attr.Value{
		"id":                  types.Int64Value(int64(gw.ID)),
		"label":               types.StringValue(gw.Label),
		"addresses":           addresses,
		"portset_assignments": types.Int64Value(int64(gw.PortsetAssignments)),
		"portset_capacity":    types.Int64Value(int64(gw.PortsetCapacity)),
		"portsets":            portsets,
	})
}

// FlattenNATGatewayDetails exposes the read-only NAT Gateway details block for use by other packages.
func FlattenNATGatewayDetails(
	ctx context.Context,
	gw *linodego.VPCSubnetNATGateway,
) (types.Object, diag.Diagnostics) {
	return flattenNATGatewayDetails(ctx, gw)
}

// isNATGatewayIDNull reports whether a natgateway block is set but its id is null ({ id = null }).
func isNATGatewayIDNull(obj types.Object) bool {
	if obj.IsNull() || obj.IsUnknown() {
		return false
	}

	id, ok := obj.Attributes()["id"]
	return ok && id.IsNull()
}

func (m *BaseModel) CopyFrom(other BaseModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.VPCId = helper.KeepOrUpdateValue(m.VPCId, other.VPCId, preserveKnown)
	m.VPCType = helper.KeepOrUpdateValue(m.VPCType, other.VPCType, preserveKnown)

	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Linodes = helper.KeepOrUpdateValue(m.Linodes, other.Linodes, preserveKnown)
	m.IPv4 = helper.KeepOrUpdateValue(m.IPv4, other.IPv4, preserveKnown)
	m.Databases = helper.KeepOrUpdateValue(m.Databases, other.Databases, preserveKnown)
	m.Nodebalancers = helper.KeepOrUpdateValue(m.Nodebalancers, other.Nodebalancers, preserveKnown)
}

func (d *BaseModel) FlattenSubnet(
	ctx context.Context,
	subnet *linodego.VPCSubnet,
	preserveKnown bool,
) diag.Diagnostics {
	d.ID = helper.KeepOrUpdateString(d.ID, strconv.Itoa(subnet.ID), preserveKnown)

	linodesList, diags := FlattenSubnetLinodes(ctx, subnet.Linodes)
	if diags.HasError() {
		return diags
	}
	d.Linodes = helper.KeepOrUpdateValue(d.Linodes, *linodesList, preserveKnown)

	d.Created = helper.KeepOrUpdateValue(
		d.Created,
		timetypes.NewRFC3339TimePointerValue(subnet.Created),
		preserveKnown,
	)
	d.Updated = helper.KeepOrUpdateValue(
		d.Updated,
		timetypes.NewRFC3339TimePointerValue(subnet.Updated),
		preserveKnown,
	)
	d.Label = helper.KeepOrUpdateString(d.Label, subnet.Label, preserveKnown)

	d.IPv4 = helper.KeepOrUpdateString(d.IPv4, subnet.IPv4, preserveKnown)
	d.VPCType = helper.KeepOrUpdateString(d.VPCType, string(subnet.VPCType), preserveKnown)

	dbsList, diags := FlattenSubnetDatabases(ctx, subnet.Databases)
	if diags.HasError() {
		return diags
	}
	d.Databases = helper.KeepOrUpdateValue(d.Databases, *dbsList, preserveKnown)

	nbsList, diags := FlattenSubnetNodebalancers(ctx, subnet.Nodebalancers)
	if diags.HasError() {
		return diags
	}
	d.Nodebalancers = helper.KeepOrUpdateValue(d.Nodebalancers, *nbsList, preserveKnown)

	return nil
}

/*
Resource-specific Implementation
*/

type ResourceModel struct {
	BaseModel
	IPv6              types.List   `tfsdk:"ipv6"`
	NATGateway        types.Object `tfsdk:"natgateway"`
	NATGatewayDetails types.Object `tfsdk:"natgateway_details"`
}

type ResourceModelIPv6 struct {
	Range          customtypes.LinodeAutoAllocRangeValue `tfsdk:"range"`
	AllocatedRange types.String                          `tfsdk:"allocated_range"`
}

func (m *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	m.BaseModel.CopyFrom(other.BaseModel, preserveKnown)
	m.IPv6 = helper.KeepOrUpdateValue(m.IPv6, other.IPv6, preserveKnown)
	m.NATGateway = helper.KeepOrUpdateValue(m.NATGateway, other.NATGateway, preserveKnown)
	m.NATGatewayDetails = helper.KeepOrUpdateValue(m.NATGatewayDetails, other.NATGatewayDetails, preserveKnown)
}

func (m *ResourceModel) FlattenSubnet(
	ctx context.Context,
	subnet *linodego.VPCSubnet,
	preserveKnown bool,
) diag.Diagnostics {
	m.BaseModel.FlattenSubnet(ctx, subnet, preserveKnown)

	ipv6AddressModels := helper.MapSlice(
		subnet.IPv6,
		func(subnet linodego.VPCIPv6Range) ResourceModelIPv6 {
			return ResourceModelIPv6{
				Range:          customtypes.LinodeAutoAllocRangeValue{StringValue: types.StringValue(subnet.Range)},
				AllocatedRange: types.StringValue(subnet.Range),
			}
		},
	)

	ipv6AddressesList, diags := types.ListValueFrom(ctx, ResourceSchemaIPv6NestedObject.Type(), ipv6AddressModels)
	if diags.HasError() {
		return diags
	}

	m.IPv6 = helper.KeepOrUpdateValue(
		m.IPv6,
		ipv6AddressesList,
		false,
	)

	natGatewayDetails, natGatewayDiags := flattenNATGatewayDetails(ctx, subnet.NATGateway)
	if natGatewayDiags.HasError() {
		return natGatewayDiags
	}
	m.NATGatewayDetails = helper.KeepOrUpdateValue(m.NATGatewayDetails, natGatewayDetails, preserveKnown)

	natGateway, natGatewayDiags := flattenNATGatewaySetting(subnet.NATGateway)
	if natGatewayDiags.HasError() {
		return natGatewayDiags
	}

	// Both an omitted block and { id = null } mean "no NAT Gateway". Keep a configured
	// { id = null } so Read does not return a value that conflicts with the config.
	if subnet.NATGateway == nil && isNATGatewayIDNull(m.NATGateway) {
		natGateway = m.NATGateway
	}
	m.NATGateway = natGateway

	return diags
}

/*
Data Source-specific Implementation
*/

type DataSourceModel struct {
	BaseModel
	IPv6              types.List   `tfsdk:"ipv6"`
	NATGatewayDetails types.Object `tfsdk:"natgateway"`
}

type DataSourceModelIPv6 struct {
	Range customtypes.LinodeAutoAllocRangeValue `tfsdk:"range"`
}

func (m *DataSourceModel) CopyFrom(other DataSourceModel, preserveKnown bool) {
	m.BaseModel.CopyFrom(other.BaseModel, preserveKnown)
	m.IPv6 = helper.KeepOrUpdateValue(m.IPv6, other.IPv6, preserveKnown)
	m.NATGatewayDetails = helper.KeepOrUpdateValue(m.NATGatewayDetails, other.NATGatewayDetails, preserveKnown)
}

func (m *DataSourceModel) FlattenSubnet(
	ctx context.Context,
	subnet *linodego.VPCSubnet,
	preserveKnown bool,
) diag.Diagnostics {
	m.BaseModel.FlattenSubnet(ctx, subnet, preserveKnown)

	ipv6AddressModels := helper.MapSlice(
		subnet.IPv6,
		func(subnet linodego.VPCIPv6Range) DataSourceModelIPv6 {
			return DataSourceModelIPv6{
				Range: customtypes.LinodeAutoAllocRangeValue{StringValue: types.StringValue(subnet.Range)},
			}
		},
	)

	ipv6AddressesList, diags := types.ListValueFrom(ctx, DataSourceSchemaIPv6NestedObject.Type(), ipv6AddressModels)
	if diags.HasError() {
		return diags
	}

	m.IPv6 = helper.KeepOrUpdateValue(
		m.IPv6,
		ipv6AddressesList,
		preserveKnown,
	)

	natGatewayDetails, natGatewayDiags := flattenNATGatewayDetails(ctx, subnet.NATGateway)
	if natGatewayDiags.HasError() {
		return natGatewayDiags
	}
	m.NATGatewayDetails = helper.KeepOrUpdateValue(m.NATGatewayDetails, natGatewayDetails, preserveKnown)

	return diags
}
