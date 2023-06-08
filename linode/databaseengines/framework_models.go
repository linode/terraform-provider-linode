package databaseengines

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

type DatabaseEngineModel struct {
	Engine  types.String `tfsdk:"engine"`
	ID      types.String `tfsdk:"id"`
	Version types.String `tfsdk:"version"`
}

type DatabaseEngineFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Latest  types.Bool                       `tfsdk:"latest"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Engines []DatabaseEngineModel            `tfsdk:"engines"`
}

func (model *DatabaseEngineFilterModel) parseEngines(engines []linodego.DatabaseEngine) {
	parseEngine := func(engine linodego.DatabaseEngine) DatabaseEngineModel {
		var m DatabaseEngineModel

		m.ID = types.StringValue(engine.ID)
		m.Engine = types.StringValue(engine.Engine)
		m.Version = types.StringValue(engine.Version)

		return m
	}

	result := make([]DatabaseEngineModel, len(engines))

	for i, engine := range engines {
		result[i] = parseEngine(engine)
	}

	model.Engines = result
}
