package accountavailability

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

type AccountAvailabilityModel struct {
	Region      types.String `tfsdk:"region"`
	Unavailable types.Set    `tfsdk:"unavailable"`
	Available   types.Set    `tfsdk:"available"`
}

func (d *AccountAvailabilityModel) ParseAvailability(
	ctx context.Context,
	availability *linodego.AccountAvailability,
) diag.Diagnostics {
	d.Region = types.StringValue(availability.Region)

	unavailable, diags := types.SetValueFrom(ctx, types.StringType, availability.Unavailable)
	if diags.HasError() {
		return diags
	}
	d.Unavailable = unavailable

	available, diags := types.SetValueFrom(ctx, types.StringType, availability.Available)
	if diags.HasError() {
		return diags
	}
	d.Available = available

	return nil
}
