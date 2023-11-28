package databasebackups

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type DatabaseBackupModel struct {
	Created types.String `tfsdk:"created"`
	Label   types.String `tfsdk:"label"`
	ID      types.Int64  `tfsdk:"id"`
	Type    types.String `tfsdk:"type"`
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

func (m *DatabaseBackupModel) ParseMySQLBackup(backup linodego.MySQLDatabaseBackup) {
	m.ID = types.Int64Value(int64(backup.ID))
	m.Label = types.StringValue(backup.Label)
	m.Type = types.StringValue(backup.Type)
	m.Created = types.StringValue(backup.Created.Format(time.RFC3339))
}

func (model *DatabaseBackupFilterModel) parseMySQLBackups(backups []linodego.MySQLDatabaseBackup) {
	result := make([]DatabaseBackupModel, len(backups))

	for i, backup := range backups {
		result[i].ParseMySQLBackup(backup)
	}

	model.Backups = result
}

func (m *DatabaseBackupModel) ParsePostgresSQLBackup(backup linodego.PostgresDatabaseBackup) {
	m.ID = types.Int64Value(int64(backup.ID))
	m.Label = types.StringValue(backup.Label)
	m.Type = types.StringValue(backup.Type)
	m.Created = types.StringValue(backup.Created.Format(time.RFC3339))
}

func (model *DatabaseBackupFilterModel) parsePostgresSQLBackups(backups []linodego.PostgresDatabaseBackup) {
	result := make([]DatabaseBackupModel, len(backups))

	for i, backup := range backups {
		result[i].ParsePostgresSQLBackup(backup)
	}

	model.Backups = result
}
