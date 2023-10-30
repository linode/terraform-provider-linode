package vpcsubnet

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
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

func parseLinodeInterface(iface linodego.VPCSubnetLinodeInterface) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(LinodeInterfaceObjectType.AttrTypes, map[string]attr.Value{
		"id":     types.Int64Value(int64(iface.ID)),
		"active": types.BoolValue(iface.Active),
	})

}

func ParseLinode(ctx context.Context, linode linodego.VPCSubnetLinode) (*types.Object, diag.Diagnostics) {
	result := map[string]attr.Value{
		"id": types.Int64Value(int64(linode.ID)),
	}

	ifaces := make([]types.Object, len(linode.Interfaces))

	for i, iface := range linode.Interfaces {
		ifaceObj, d := parseLinodeInterface(iface)
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

func (d *VPCSubnetModel) parseComputedAttributes(
	ctx context.Context,
	subnet *linodego.VPCSubnet,
) diag.Diagnostics {
	d.ID = types.Int64Value(int64(subnet.ID))

	linodes := make([]types.Object, len(subnet.Linodes))

	for i, inst := range subnet.Linodes {
		linodeObj, d := ParseLinode(ctx, inst)
		if d.HasError() {
			return d
		}

		linodes[i] = *linodeObj
	}

	linodesList, dg := types.ListValueFrom(ctx, LinodeObjectType, linodes)
	if dg.HasError() {
		return dg
	}
	d.Linodes = linodesList

	d.Created = timetypes.NewRFC3339TimePointerValue(subnet.Created)
	d.Updated = timetypes.NewRFC3339TimePointerValue(subnet.Updated)

	return nil
}

func (d *VPCSubnetModel) parseVPCSubnet(
	ctx context.Context,
	subnet *linodego.VPCSubnet,
) diag.Diagnostics {
	d.Label = types.StringValue(subnet.Label)
	d.IPv4 = types.StringValue(subnet.IPv4)

	return d.parseComputedAttributes(ctx, subnet)
}
