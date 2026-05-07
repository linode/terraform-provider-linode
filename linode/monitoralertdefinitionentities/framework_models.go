package monitoralertdefinitionentities

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

// EntityFilterModel describes the Terraform data source model.
type EntityFilterModel struct {
	ID          types.String                     `tfsdk:"id"`
	ServiceType types.String                     `tfsdk:"service_type"`
	AlertID     types.Int64                      `tfsdk:"alert_id"`
	Filters     frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Entities    []EntityDataSourceModel          `tfsdk:"entities"`
}

// EntityDataSourceModel describes a single entity.
type EntityDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
	URL   types.String `tfsdk:"url"`
	Type  types.String `tfsdk:"type"`
}

func (data *EntityFilterModel) parseEntities(
	entities []linodego.AlertDefinitionEntity,
) diag.Diagnostics {
	result := make([]EntityDataSourceModel, len(entities))
	for i := range entities {
		result[i] = EntityDataSourceModel{
			ID:    types.StringValue(entities[i].ID),
			Label: types.StringValue(entities[i].Label),
			URL:   types.StringValue(entities[i].URL),
			Type:  types.StringValue(entities[i].Type),
		}
	}

	data.Entities = result

	return nil
}
