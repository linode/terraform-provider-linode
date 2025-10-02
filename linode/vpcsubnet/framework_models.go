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
)

type VPCSubnetModel struct {
	ID      types.String      `tfsdk:"id"`
	VPCID   types.Int64       `tfsdk:"vpc_id"`
	Label   types.String      `tfsdk:"label"`
	IPv4    types.String      `tfsdk:"ipv4"`
	Linodes types.List        `tfsdk:"linodes"`
	Created timetypes.RFC3339 `tfsdk:"created"`
	Updated timetypes.RFC3339 `tfsdk:"updated"`
}

func FlattenSubnetLinodeInterface(iface linodego.VPCSubnetLinodeInterface) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(LinodeInterfaceObjectType.AttrTypes, map[string]attr.Value{
		"id":     types.Int64Value(int64(iface.ID)),
		"active": types.BoolValue(iface.Active),
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

func (v *VPCSubnetModel) FlattenSubnet(
	ctx context.Context,
	subnet *linodego.VPCSubnet,
	preserveKnown bool,
) diag.Diagnostics {
	v.ID = helper.KeepOrUpdateString(v.ID, strconv.Itoa(subnet.ID), preserveKnown)

	linodesList, diags := FlattenSubnetLinodes(ctx, subnet.Linodes)
	if diags.HasError() {
		return diags
	}
	v.Linodes = helper.KeepOrUpdateValue(v.Linodes, *linodesList, preserveKnown)

	v.Created = helper.KeepOrUpdateValue(
		v.Created,
		timetypes.NewRFC3339TimePointerValue(subnet.Created),
		preserveKnown,
	)
	v.Updated = helper.KeepOrUpdateValue(
		v.Updated,
		timetypes.NewRFC3339TimePointerValue(subnet.Updated),
		preserveKnown,
	)
	v.Label = helper.KeepOrUpdateString(v.Label, subnet.Label, preserveKnown)
	v.IPv4 = helper.KeepOrUpdateString(v.IPv4, subnet.IPv4, preserveKnown)

	return nil
}

func (v *VPCSubnetModel) CopyFrom(other VPCSubnetModel, preserveKnown bool) {
	v.ID = helper.KeepOrUpdateValue(v.ID, other.ID, preserveKnown)
	v.VPCID = helper.KeepOrUpdateValue(v.VPCID, other.VPCID, preserveKnown)
	v.Label = helper.KeepOrUpdateValue(v.Label, other.Label, preserveKnown)
	v.IPv4 = helper.KeepOrUpdateValue(v.IPv4, other.IPv4, preserveKnown)
	v.Linodes = helper.KeepOrUpdateValue(v.Linodes, other.Linodes, preserveKnown)
	v.Created = helper.KeepOrUpdateValue(v.Created, other.Created, preserveKnown)
	v.Updated = helper.KeepOrUpdateValue(v.Updated, other.Updated, preserveKnown)
}
