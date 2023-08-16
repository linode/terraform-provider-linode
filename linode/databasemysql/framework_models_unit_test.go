package databasemysql

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseMySQLDatabase(t *testing.T) {
	mockDB := &linodego.MySQLDatabase{
		ID:     123,
		Status: "active",
		Label:  "example-db",
		Hosts: linodego.DatabaseHost{
			Primary:   "lin-123-456-mysql-mysql-primary.servers.linodedb.net",
			Secondary: "lin-123-456-mysql-primary-private.servers.linodedb.net",
		},
		Region:          "us-east",
		Type:            "g6-dedicated-2",
		Engine:          "mysql",
		Version:         "8.0.26",
		ClusterSize:     3,
		ReplicationType: "semi_synch",
		SSLConnection:   true,
		Encrypted:       false,
		AllowList:       []string{"203.0.113.1/32", "192.0.1.0/24"},
		Created:         &time.Time{},
		Updated:         &time.Time{},
		Updates: linodego.DatabaseMaintenanceWindow{
			DayOfWeek:   1,
			Duration:    3,
			Frequency:   "weekly",
			HourOfDay:   0,
			WeekOfMonth: nil,
		},
	}

	data := &DataSourceModel{}
	diagnostics := data.parseMySQLDatabase(context.Background(), mockDB)

	assert.False(t, diagnostics.HasError(), "Expected no error")

	assert.Equal(t, types.Int64Value(123), data.DatabaseID)
	assert.Equal(t, types.StringValue("active"), data.Status)
	assert.Equal(t, types.StringValue("example-db"), data.Label)
	assert.Equal(t, types.StringValue("lin-123-456-mysql-mysql-primary.servers.linodedb.net"), data.HostPrimary)
	assert.Equal(t, types.StringValue("lin-123-456-mysql-primary-private.servers.linodedb.net"), data.HostSecondary)
	assert.Equal(t, types.StringValue("us-east"), data.Region)
	assert.Equal(t, types.StringValue("g6-dedicated-2"), data.Type)
	assert.Equal(t, types.StringValue("mysql"), data.Engine)
	assert.Equal(t, types.StringValue("8.0.26"), data.Version)
	assert.Equal(t, types.Int64Value(3), data.ClusterSize)
	assert.Equal(t, types.StringValue("semi_synch"), data.ReplicationType)
	assert.Equal(t, types.BoolValue(true), data.SSLConnection)
	assert.Equal(t, types.BoolValue(false), data.Encrypted)
	assert.Contains(t, data.AllowList.String(), "203.0.113.1/32")
	assert.Contains(t, data.AllowList.String(), "192.0.1.0/24")

	assert.Contains(t, data.Updates.String(), "sunday")
	assert.Contains(t, data.Updates.String(), "3")
	assert.Contains(t, data.Updates.String(), "weekly")
	assert.Contains(t, data.Updates.String(), "0")
}

func TestParseMySQLDatabaseSSL(t *testing.T) {
	mockDBSSL := &linodego.MySQLDatabaseSSL{
		CACertificate: []byte("-----BEGIN CERTIFICATE-----\nMIIDdTCCAl2gAwIBAgIUT01...\n...u4QIDAQABo1MwUTAdBgNV...\n-----END CERTIFICATE-----"),
	}

	data := &DataSourceModel{}
	data.parseMySQLDatabaseSSL(mockDBSSL)

	assert.Equal(t, types.StringValue("-----BEGIN CERTIFICATE-----\nMIIDdTCCAl2gAwIBAgIUT01...\n...u4QIDAQABo1MwUTAdBgNV...\n-----END CERTIFICATE-----"), data.CACert)
}

func TestParseMySQLDatabaseCredentials(t *testing.T) {
	mockDBCred := &linodego.MySQLDatabaseCredential{
		Username: "linode_sqldb_user",
		Password: "password123",
	}

	data := &DataSourceModel{}
	data.parseMySQLDatabaseCredentials(mockDBCred)

	assert.Equal(t, types.StringValue("linode_sqldb_user"), data.RootUsername)
	assert.Equal(t, types.StringValue("password123"), data.RootPassword)
}
