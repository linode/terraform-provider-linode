package databasebackups

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

type DatabaseBackupModel struct {
	Created types.String `tfsdk:"created"`
	Label   types.String `tfsdk:"engine"`
	ID      types.String `tfsdk:"id"`
	Type    types.String `tfsdk:"version"`
}

type DatabaseBackupFilterModel struct {
	DatabaseID   types.Int64                      `tfsdk:"database_id"`
	ID           types.Int64                      `tfsdk:"id"`
	DatabaseType types.String                     `tfsdk:"database_type"`
	Latest       types.Bool                       `tfsdk:"latest"`
	Order        types.String                     `tfsdk:"order"`
	OrderBy      types.String                     `tfsdk:"order_by"`
	Filters      frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Backups      []DatabaseBackupModel            `tfsdk:"backups"`
}

// func (model *DatabaseBackupFilterModel) parseBackups(backups []linodego.Database) {
// 	parseEngine := func(engine linodego.DatabaseEngine) DatabaseEngineModel {
// 		var m DatabaseEngineModel

// 		m.ID = types.StringValue(engine.ID)
// 		m.Engine = types.StringValue(engine.Engine)
// 		m.Version = types.StringValue(engine.Version)

// 		return m
// 	}

// 	result := make([]DatabaseEngineModel, len(engines))

// 	for i, engine := range engines {
// 		result[i] = parseEngine(engine)
// 	}

// 	model.Engines = result
// }
