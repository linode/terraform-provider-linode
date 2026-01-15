package vpcsubnet

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/customtypes"
)

/*
Shared Implementation
*/

type BaseModel struct {
	ID        types.String      `tfsdk:"id"`
	VPCId     types.Int64       `tfsdk:"vpc_id"`
	Label     types.String      `tfsdk:"label"`
	IPv4      types.String      `tfsdk:"ipv4"`
	Linodes   types.List        `tfsdk:"linodes"`
	Created   timetypes.RFC3339 `tfsdk:"created"`
	Updated   timetypes.RFC3339 `tfsdk:"updated"`
	Databases types.List        `tfsdk:"databases"`
}

func FlattenSubnetLinodeInterface(iface linodego.VPCSubnetLinodeInterface) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(LinodeInterfaceObjectType.AttrTypes, map[string]attr.Value{
		"id":        types.Int64Value(int64(iface.ID)),
		"config_id": types.Int64PointerValue(helper.IntPtrToInt64Ptr(iface.ConfigID)),
		"active":    types.BoolValue(iface.Active),
	})
}

func FlattenSubnetDatabaseIPv6Range(
	ipv6Range linodego.VPCIPv6Range,
) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(IPV6RangeObjectType.AttrTypes, map[string]attr.Value{
		"range": types.StringValue(ipv6Range.Range),
	})
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
	result := map[string]attr.Value{
		"id":         types.Int64Value(int64(db.ID)),
		"ipv4_range": types.StringPointerValue(db.IPv4Range),
	}

	ipv6Ranges := make([]types.Object, len(db.IPv6Ranges))

	for i, ipv6Range := range db.IPv6Ranges {
		ipv6RangeObj, d := FlattenSubnetDatabaseIPv6Range(linodego.VPCIPv6Range{Range: ipv6Range})
		if d.HasError() {
			return nil, d
		}

		ipv6Ranges[i] = ipv6RangeObj
	}

	ipv6RangesList, d := types.ListValueFrom(ctx, IPV6RangeObjectType, ipv6Ranges)
	if d.HasError() {
		return nil, d
	}

	result["ipv6_ranges"] = ipv6RangesList

	resultObject, d := types.ObjectValue(LinodeObjectType.AttrTypes, result)
	return &resultObject, d
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

func (m *BaseModel) CopyFrom(other BaseModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.VPCId = helper.KeepOrUpdateValue(m.VPCId, other.VPCId, preserveKnown)

	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Linodes = helper.KeepOrUpdateValue(m.Linodes, other.Linodes, preserveKnown)
	m.IPv4 = helper.KeepOrUpdateValue(m.IPv4, other.IPv4, preserveKnown)
	m.Databases = helper.KeepOrUpdateValue(m.Databases, other.Databases, preserveKnown)
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

	dbsList, diags := FlattenSubnetDatabases(ctx, subnet.Databases)
	if diags.HasError() {
		return diags
	}
	d.Databases = helper.KeepOrUpdateValue(d.Databases, *dbsList, preserveKnown)

	return nil
}

/*
Resource-specific Implementation
*/

type ResourceModel struct {
	BaseModel
	IPv6 types.List `tfsdk:"ipv6"`
}

type ResourceModelIPv6 struct {
	Range          customtypes.LinodeAutoAllocRangeValue `tfsdk:"range"`
	AllocatedRange types.String                          `tfsdk:"allocated_range"`
}

func (m *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	m.BaseModel.CopyFrom(other.BaseModel, preserveKnown)
	m.IPv6 = helper.KeepOrUpdateValue(m.IPv6, other.IPv6, preserveKnown)
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

	return diags
}

/*
Data Source-specific Implementation
*/

type DataSourceModel struct {
	BaseModel
	IPv6 types.List `tfsdk:"ipv6"`
}

type DataSourceModelIPv6 struct {
	Range customtypes.LinodeAutoAllocRangeValue `tfsdk:"range"`
}

func (m *DataSourceModel) CopyFrom(other DataSourceModel, preserveKnown bool) {
	m.BaseModel.CopyFrom(other.BaseModel, preserveKnown)
	m.IPv6 = helper.KeepOrUpdateValue(m.IPv6, other.IPv6, preserveKnown)
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

	return diags
}
