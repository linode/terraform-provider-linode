package vpcsubnet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type VPCSubnetInterfaceModel struct {
	ID     types.Int64 `tfsdk:"id"`
	Active types.Bool  `tfsdk:"active"`
}

type VPCSubnetLinodeModel struct {
	ID         types.Int64               `tfsdk:"id"`
	Interfaces []VPCSubnetInterfaceModel `tfsdk:"interfaces"`
}

type VPCSubnetModel struct {
	ID      types.Int64       `tfsdk:"id"`
	VPCId   types.Int64       `tfsdk:"vpc_id"`
	Label   types.String      `tfsdk:"label"`
	IPv4    types.String      `tfsdk:"ipv4"`
	Linodes types.List        `tfsdk:"linodes"`
	Created timetypes.RFC3339 `tfsdk:"created"`
	Updated timetypes.RFC3339 `tfsdk:"updated"`
}

func parseLinode(linode linodego.VPCSubnetLinode) VPCSubnetLinodeModel {
	result := VPCSubnetLinodeModel{
		ID: types.Int64Value(int64(linode.ID)),
	}

	result.Interfaces = make([]VPCSubnetInterfaceModel, len(linode.Interfaces))
	for i, v := range linode.Interfaces {
		result.Interfaces[i] = VPCSubnetInterfaceModel{
			ID:     types.Int64Value(int64(v.ID)),
			Active: types.BoolValue(v.Active),
		}
	}

	return result
}

func (d *VPCSubnetModel) parseComputedAttributes(
	ctx context.Context,
	subnet *linodego.VPCSubnet,
) diag.Diagnostics {
	d.ID = types.Int64Value(int64(subnet.ID))

	linodes := make([]VPCSubnetLinodeModel, len(subnet.Linodes))
	for i, v := range subnet.Linodes {
		linodes[i] = parseLinode(v)
	}

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
