package monitorlogsdestinations

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

// LogsDestinationModel is an item in the list data source.
type LogsDestinationModel struct {
	ID        types.Int64       `tfsdk:"id"`
	Label     types.String      `tfsdk:"label"`
	Type      types.String      `tfsdk:"type"`
	Status    types.String      `tfsdk:"status"`
	CreatedBy types.String      `tfsdk:"created_by"`
	UpdatedBy types.String      `tfsdk:"updated_by"`
	Created   timetypes.RFC3339 `tfsdk:"created"`
	Updated   timetypes.RFC3339 `tfsdk:"updated"`
	Version   types.Int64       `tfsdk:"version"`
}

// LogsDestinationFilterModel is the model for the list data source.
type LogsDestinationFilterModel struct {
	ID           types.String                     `tfsdk:"id"`
	Filters      frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order        types.String                     `tfsdk:"order"`
	OrderBy      types.String                     `tfsdk:"order_by"`
	Destinations []LogsDestinationModel           `tfsdk:"destinations"`
}

func (m *LogsDestinationFilterModel) ParseLogsDestinations(destinations []linodego.LogsDestination) {
	result := make([]LogsDestinationModel, len(destinations))
	for i, dest := range destinations {
		result[i] = LogsDestinationModel{
			ID:        types.Int64Value(int64(dest.ID)),
			Label:     types.StringValue(dest.Label),
			Type:      types.StringValue(string(dest.Type)),
			Status:    types.StringValue(string(dest.Status)),
			CreatedBy: types.StringValue(dest.CreatedBy),
			UpdatedBy: types.StringValue(dest.UpdatedBy),
			Created:   timetypes.NewRFC3339TimePointerValue(dest.Created),
			Updated:   timetypes.NewRFC3339TimePointerValue(dest.Updated),
			Version:   types.Int64Value(int64(dest.Version)),
		}
	}
	m.Destinations = result
}
