package vpcsubnet

import (
	"context"

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

func (d *VPCSubnetModel) parseComputedAttributes(
	ctx context.Context,
	subnet *linodego.VPCSubnet,
) diag.Diagnostics {
	d.ID = types.Int64Value(int64(subnet.ID))

	linodes, diag := types.ListValueFrom(ctx, types.Int64Type, subnet.Linodes)
	if diag != nil {
		return diag
	}
	d.Linodes = linodes

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
