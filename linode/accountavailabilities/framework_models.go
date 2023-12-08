package accountavailabilities

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/accountavailability"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type AccountAvailabilityFilterModel struct {
	ID             types.String                                   `tfsdk:"id"`
	Filters        frameworkfilter.FiltersModelType               `tfsdk:"filter"`
	Availabilities []accountavailability.AccountAvailabilityModel `tfsdk:"availabilities"`
}

func (model *AccountAvailabilityFilterModel) parseAvailabilities(ctx context.Context, availabilities []linodego.AccountAvailability) diag.Diagnostics {
	result := make([]accountavailability.AccountAvailabilityModel, len(availabilities))

	for i := range availabilities {
		var availabilityModel accountavailability.AccountAvailabilityModel

		d := availabilityModel.ParseAvailability(ctx, &availabilities[i])
		if d.HasError() {
			return d
		}

		result[i] = availabilityModel
	}

	model.Availabilities = result

	return nil
}
