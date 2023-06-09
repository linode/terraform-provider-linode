package databasebackups

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

type DatabaseBackupModel struct {
	Created types.String `tfsdk:"created"`
	Label   types.String `tfsdk:"engine"`
	ID      types.Int64  `tfsdk:"id"`
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

func (model *DatabaseBackupFilterModel) parseMySQLBackups(backups []linodego.MySQLDatabaseBackup) {
	parseBackup := func(backup linodego.MySQLDatabaseBackup) DatabaseBackupModel {
		var m DatabaseBackupModel

		m.ID = types.Int64Value(int64(backup.ID))
		m.Label = types.StringValue(backup.Label)
		m.Type = types.StringValue(backup.Type)
		m.Created = types.StringValue(backup.Created.Format(time.RFC3339))

		return m
	}

	result := make([]DatabaseBackupModel, len(backups))

	for i, backup := range backups {
		result[i] = parseBackup(backup)
	}

	model.Backups = result
}

func (model *DatabaseBackupFilterModel) parsePostgresSQLBackups(backups []linodego.PostgresDatabaseBackup) {
	parseBackup := func(backup linodego.PostgresDatabaseBackup) DatabaseBackupModel {
		var m DatabaseBackupModel

		m.ID = types.Int64Value(int64(backup.ID))
		m.Label = types.StringValue(backup.Label)
		m.Type = types.StringValue(backup.Type)
		m.Created = types.StringValue(backup.Created.Format(time.RFC3339))

		return m
	}

	result := make([]DatabaseBackupModel, len(backups))

	for i, backup := range backups {
		result[i] = parseBackup(backup)
	}

	model.Backups = result
}
