package accountavailabilities

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type AccountAvailabilityModel struct {
	Region      types.String `tfsdk:"region"`
	Unavailable types.Set    `tfsdk:"unavailable"`
}

type AccountAvailabilityFilterModel struct {
	ID             types.String                     `tfsdk:"id"`
	Filters        frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Availabilities []AccountAvailabilityModel       `tfsdk:"availabilities"`
}

func (model *AccountAvailabilityFilterModel) parseAvailabilities(ctx context.Context, availabilities []linodego.AccountAvailability) diag.Diagnostics {
	parseAvailability := func(entry linodego.AccountAvailability) (AccountAvailabilityModel, diag.Diagnostics) {
		var m AccountAvailabilityModel

		m.Region = types.StringValue(entry.Region)

		unavailable, d := types.SetValueFrom(ctx, types.StringType, entry.Unavailable)
		if d.HasError() {
			return m, d
		}

		m.Unavailable = unavailable

		return m, nil
	}

	result := make([]AccountAvailabilityModel, len(availabilities))

	for i, entry := range availabilities {
		parsedEntry, d := parseAvailability(entry)
		if d.HasError() {
			return d
		}

		result[i] = parsedEntry
	}

	model.Availabilities = result

	return nil
}
