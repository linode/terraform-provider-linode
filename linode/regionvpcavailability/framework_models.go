package regionvpcavailability

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type RegionVPCAvailabilityModel struct {
	Available                  types.Bool   `tfsdk:"available"`
	AvailableIPV6PrefixLengths types.List   `tfsdk:"available_ipv6_prefix_lengths"`
	ID                         types.String `tfsdk:"id"`
}

func (model *RegionVPCAvailabilityModel) ParseRegionVPCAvailability(
	ctx context.Context,
	regionVPCAvailability *linodego.RegionVPCAvailability,
) diag.Diagnostics {
	availableIPV6PrefixLengths, diags := types.ListValueFrom(
		ctx,
		types.Int64Type,
		regionVPCAvailability.AvailableIPV6PrefixLengths)
	if diags.HasError() {
		return diags
	}

	model.AvailableIPV6PrefixLengths = availableIPV6PrefixLengths
	model.Available = types.BoolValue(regionVPCAvailability.Available)
	model.ID = types.StringValue(regionVPCAvailability.Region)

	return nil
}
