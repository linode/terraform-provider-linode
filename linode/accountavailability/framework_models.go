package accountavailability

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	Region      types.String `tfsdk:"region"`
	Unavailable types.List   `tfsdk:"unavailable"`
}

func (d *DataSourceModel) parseAvailability(
	ctx context.Context,
	availability *linodego.AccountAvailability,
) diag.Diagnostics {
	d.Region = types.StringValue(availability.Region)

	unavailable, diags := types.ListValueFrom(ctx, types.StringType, availability.Unavailable)
	if diags != nil {
		return diags
	}
	d.Unavailable = unavailable

	return nil
}
