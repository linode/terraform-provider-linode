package vpcsubnet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type VPCSubnetModel struct {
	ID      types.Int64       `tfsdk:"id"`
	VPCId   types.Int64       `tfsdk:"vpc_id"`
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

func FlattenSubnetLinode(ctx context.Context, linode linodego.VPCSubnetLinode) (*types.Object, diag.Diagnostics) {
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

func FlattenSubnetLinodes(ctx context.Context, subnetLinodes []linodego.VPCSubnetLinode) (*types.List, diag.Diagnostics) {
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

func (d *VPCSubnetModel) FlattenSubnet(
	ctx context.Context,
	subnet *linodego.VPCSubnet,
	preserveKnown bool,
) diag.Diagnostics {
	d.ID = helper.KeepOrUpdateInt64(d.ID, int64(subnet.ID), preserveKnown)

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

	return nil
}
