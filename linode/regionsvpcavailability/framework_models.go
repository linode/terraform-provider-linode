package regionsvpcavailability

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/regionvpcavailability"
)

type regionsVPCAvailabilityModel struct {
	ID                     types.String                                       `tfsdk:"id"`
	RegionsVPCAvailability []regionvpcavailability.RegionVPCAvailabilityModel `tfsdk:"regions_vpc_availability"`
}

func (model *regionsVPCAvailabilityModel) parseRegionsVPCAvailability(
	ctx context.Context,
	regionsVPCAvailability []linodego.RegionVPCAvailability,
) diag.Diagnostics {
	result := make([]regionvpcavailability.RegionVPCAvailabilityModel, len(regionsVPCAvailability))

	for i, regionVPCAvailability := range regionsVPCAvailability {
		regionVPCAvailabilityModel := regionvpcavailability.RegionVPCAvailabilityModel{}
		regionVPCAvailabilityModel.ParseRegionVPCAvailability(ctx, &regionVPCAvailability)
		result[i] = regionVPCAvailabilityModel
	}
	model.RegionsVPCAvailability = result

	var diags diag.Diagnostics
	id, err := json.Marshal(regionsVPCAvailability)
	if err != nil {
		diags.AddError("Error marshalling json: %s", err.Error())
		return diags
	}
	model.ID = types.StringValue(string(id))

	return nil
}
