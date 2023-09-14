//go:build unit

package databasepostgresql

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParsePostgresDatabase(t *testing.T) {
	mockDatabase := linodego.PostgresDatabase{
		ID:                    123,
		Status:                "active",
		Label:                 "example-db",
		Region:                "us-east",
		Type:                  "g6-dedicated-2",
		Engine:                "postgresql",
		Version:               "13.2",
		Encrypted:             false,
		AllowList:             []string{"203.0.113.1/32", "192.0.1.0/24"},
		Port:                  3306,
		SSLConnection:         true,
		ClusterSize:           3,
		ReplicationCommitType: "local",
		ReplicationType:       "async",
		Hosts: linodego.DatabaseHost{
			Primary:   "lin-0000-000-pgsql-primary.servers.linodedb.net",
			Secondary: "lin-0000-000-pgsql-primary-private.servers.linodedb.net",
		},
		Updates: linodego.DatabaseMaintenanceWindow{
			DayOfWeek:   1,
			Duration:    3,
			Frequency:   "weekly",
			HourOfDay:   0,
			WeekOfMonth: nil,
		},
		Created: &time.Time{},
		Updated: &time.Time{},
	}

	data := &DataSourceModel{}
	diagnostics := data.parsePostgresDatabase(context.Background(), &mockDatabase)
	assert.False(t, diagnostics.HasError(), "Expected no error")

	assert.Equal(t, types.Int64Value(123), data.DatabaseID)
	assert.Equal(t, types.StringValue("active"), data.Status)
	assert.Equal(t, types.StringValue("example-db"), data.Label)

	assert.Contains(t, data.AllowList.String(), "203.0.113.1/32")
	assert.Contains(t, data.AllowList.String(), "192.0.1.0/24")

	assert.Equal(t, types.StringValue("lin-0000-000-pgsql-primary.servers.linodedb.net"), data.HostPrimary)
	assert.Equal(t, types.StringValue("lin-0000-000-pgsql-primary-private.servers.linodedb.net"), data.HostSecondary)

	assert.Contains(t, data.Updates.String(), "monday")
	assert.Contains(t, data.Updates.String(), "3")
	assert.Contains(t, data.Updates.String(), "weekly")
	assert.Contains(t, data.Updates.String(), "0")
}

func TestParsePostgresDatabaseSSL(t *testing.T) {
	mockSSL := linodego.PostgresDatabaseSSL{
		CACertificate: []byte("-----BEGIN CERTIFICATE-----\nMIIDdTCCAl2gAwIBAgIUT01...\n...u4QIDAQABo1MwUTAdBgNV...\n-----END CERTIFICATE-----"),
	}

	data := &DataSourceModel{}
	data.parsePostgresDatabaseSSL(&mockSSL)

	assert.Equal(t, types.StringValue("-----BEGIN CERTIFICATE-----\nMIIDdTCCAl2gAwIBAgIUT01...\n...u4QIDAQABo1MwUTAdBgNV...\n-----END CERTIFICATE-----"), data.CACert)
}

func TestParsePostgresDatabaseCredentials(t *testing.T) {
	mockCredential := linodego.PostgresDatabaseCredential{
		Username: "linode_postgresql_user",
		Password: "password123",
	}

	data := &DataSourceModel{}
	data.parsePostgresDatabaseCredentials(&mockCredential)

	assert.Equal(t, types.StringValue("linode_postgresql_user"), data.RootUsername)
	assert.Equal(t, types.StringValue("password123"), data.RootPassword)
}
