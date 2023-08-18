//go:build unit

package databasebackups

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseMySQLBackup(t *testing.T) {
	mockBackup := linodego.MySQLDatabaseBackup{
		ID:      123,
		Label:   "Scheduled - 02/04/22 11:11 UTC-XcCRmI",
		Type:    "manual",
		Created: &time.Time{},
	}

	model := DatabaseBackupModel{}

	model.ParseMySQLBackup(mockBackup)

	assert.Equal(t, types.Int64Value(123), model.ID)
	assert.Equal(t, types.StringValue("Scheduled - 02/04/22 11:11 UTC-XcCRmI"), model.Label)
	assert.Equal(t, types.StringValue("manual"), model.Type)
	assert.NotNil(t, model.Created) // Created field should not be nil

	expectedFormattedTime := mockBackup.Created.Format(time.RFC3339)
	assert.Equal(t, types.StringValue(expectedFormattedTime), model.Created)
}

func TestParseMySQLBackups(t *testing.T) {
	mockBackups := []linodego.MySQLDatabaseBackup{
		{
			ID:      1,
			Label:   "Scheduled - 02/07/22 11:18 UTC-XcCRmI",
			Type:    "manual",
			Created: &time.Time{},
		},
		{
			ID:      2,
			Label:   "Scheduled - 02/07/22 11:18 UTC-XcCRmI",
			Type:    "auto",
			Created: &time.Time{},
		},
	}

	model := DatabaseBackupFilterModel{}

	model.parseMySQLBackups(mockBackups)

	assert.Len(t, model.Backups, len(mockBackups))

	assert.Equal(t, types.Int64Value(1), model.Backups[0].ID)
	assert.Equal(t, types.StringValue("Scheduled - 02/07/22 11:18 UTC-XcCRmI"), model.Backups[0].Label)
	assert.Equal(t, types.StringValue("manual"), model.Backups[0].Type)
	assert.NotNil(t, model.Backups[0].Created)

	assert.Equal(t, types.Int64Value(2), model.Backups[1].ID)
	assert.Equal(t, types.StringValue("Scheduled - 02/07/22 11:18 UTC-XcCRmI"), model.Backups[1].Label)
	assert.Equal(t, types.StringValue("auto"), model.Backups[1].Type)
	assert.NotNil(t, model.Backups[1].Created)
}

func TestParsePostgresSQLBackup(t *testing.T) {
	mockBackup := linodego.PostgresDatabaseBackup{
		ID:      123,
		Label:   "Postgres Backup",
		Type:    "auto",
		Created: &time.Time{},
	}

	model := DatabaseBackupModel{}

	model.ParsePostgresSQLBackup(mockBackup)

	assert.Equal(t, types.Int64Value(123), model.ID)
	assert.Equal(t, types.StringValue("Postgres Backup"), model.Label)
	assert.Equal(t, types.StringValue("auto"), model.Type)
	assert.NotNil(t, model.Created) // Created field should not be nil

	expectedFormattedTime := mockBackup.Created.Format(time.RFC3339)
	assert.Equal(t, types.StringValue(expectedFormattedTime), model.Created)
}

func TestParsePostgresSQLBackups(t *testing.T) {
	mockBackups := []linodego.PostgresDatabaseBackup{
		{
			ID:      1,
			Label:   "Postgres Backup 1",
			Type:    "manual",
			Created: &time.Time{},
		},
		{
			ID:      2,
			Label:   "Postgres Backup 2",
			Type:    "auto",
			Created: &time.Time{},
		},
	}

	model := DatabaseBackupFilterModel{}

	model.parsePostgresSQLBackups(mockBackups)

	assert.Len(t, model.Backups, len(mockBackups))

	assert.Equal(t, types.Int64Value(1), model.Backups[0].ID)
	assert.Equal(t, types.StringValue("Postgres Backup 1"), model.Backups[0].Label)
	assert.Equal(t, types.StringValue("manual"), model.Backups[0].Type)
	assert.NotNil(t, model.Backups[0].Created)

	assert.Equal(t, types.Int64Value(2), model.Backups[1].ID)
	assert.Equal(t, types.StringValue("Postgres Backup 2"), model.Backups[1].Label)
	assert.Equal(t, types.StringValue("auto"), model.Backups[1].Type)
	assert.NotNil(t, model.Backups[1].Created)
}
