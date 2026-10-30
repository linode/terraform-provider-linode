//go:build unit

package databases

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestParseDatabases(t *testing.T) {
	currentTime := time.Now()
	currentTimeFWValue := timetypes.NewRFC3339TimePointerValue(&currentTime)
	restoreTime := currentTime.Add(-time.Hour)

	mockDB1 := linodego.Database{
		ID:          123,
		Status:      linodego.DatabaseStatusActive,
		Label:       "example-db-1",
		Region:      "us-east",
		Type:        "g6-standard-1",
		Engine:      "mysql",
		Version:     "8.0",
		Encrypted:   false,
		AllowList:   []string{"203.0.113.1/32", "192.0.1.0/24"},
		Hosts:       linodego.DatabaseHost{Primary: "primary.example.com", Standby: "secondary.example.com"},
		InstanceURI: "mysql://user:pass@primary.example.com:3306/db",
		Created:     &time.Time{},
		Updated:     &time.Time{},
		Fork: &linodego.DatabaseFork{
			Source:      12345,
			RestoreTime: &currentTime,
		},
		OldestRestoreTime:     &currentTime,
		AvailableRestoreTimes: []time.Time{restoreTime},
	}

	mockDB2 := linodego.Database{
		ID:          456,
		Status:      linodego.DatabaseStatusProvisioning,
		Label:       "example-db-2",
		Region:      "us-central",
		Type:        "g6-standard-2",
		Engine:      "postgresql",
		Version:     "13",
		Encrypted:   true,
		AllowList:   []string{"10.0.0.1/32"},
		Hosts:       linodego.DatabaseHost{Primary: "primary-pg.example.com", Standby: "secondary-pg.example.com"},
		InstanceURI: "postgresql://user:pass@primary-pg.example.com:5432/db",
		Created:     &time.Time{},
		Updated:     &time.Time{},
		Fork: &linodego.DatabaseFork{
			Source:      54321,
			RestoreTime: &currentTime,
		},
		OldestRestoreTime: &currentTime,
	}

	mockDB3 := linodego.Database{
		ID:                    789,
		Status:                linodego.DatabaseStatusActive,
		Label:                 "example-valkey-empty",
		Region:                "us-mia",
		Type:                  "g7-dedicated-4-2",
		Engine:                "valkey",
		Version:               "8.1.0",
		AvailableRestoreTimes: []time.Time{},
	}

	mockDB4 := linodego.Database{
		ID:                    987,
		Status:                linodego.DatabaseStatusActive,
		Label:                 "example-valkey",
		Region:                "us-mia",
		Type:                  "g7-dedicated-4-2",
		Engine:                "valkey",
		Version:               "8.1.0",
		AvailableRestoreTimes: []time.Time{restoreTime, currentTime},
	}

	mockDatabases := []linodego.Database{mockDB1, mockDB2, mockDB3, mockDB4}

	model := &DatabaseFilterModel{}
	model.parseDatabases(t.Context(), mockDatabases)

	assert.Len(t, model.Databases, 4)

	// Database 1 Assertions
	assert.Equal(t, types.Int64Value(123), model.Databases[0].ID)
	assert.Equal(t, types.StringValue("active"), model.Databases[0].Status)
	assert.Equal(t, types.StringValue("example-db-1"), model.Databases[0].Label)
	assert.Equal(t, currentTimeFWValue, model.Databases[0].ForkRestoreTime)
	assert.Equal(t, currentTimeFWValue, model.Databases[0].OldestRestoreTime)
	assert.Equal(t, []types.String{types.StringValue(restoreTime.Format(time.RFC3339))}, model.Databases[0].AvailableRestoreTimes)
	assert.Equal(t, int64(12345), model.Databases[0].ForkSource.ValueInt64())

	// Database 2 Assertions
	assert.Equal(t, types.Int64Value(456), model.Databases[1].ID)
	assert.Equal(t, types.StringValue("provisioning"), model.Databases[1].Status)
	assert.Equal(t, types.StringValue("example-db-2"), model.Databases[1].Label)
	assert.Equal(t, currentTimeFWValue, model.Databases[1].ForkRestoreTime)
	assert.Equal(t, currentTimeFWValue, model.Databases[1].OldestRestoreTime)
	assert.Nil(t, model.Databases[1].AvailableRestoreTimes)
	assert.Equal(t, int64(54321), model.Databases[1].ForkSource.ValueInt64())

	// Valkey preserves the distinction between an explicitly empty and a nil list.
	assert.Equal(t, types.StringValue("valkey"), model.Databases[2].Engine)
	assert.NotNil(t, model.Databases[2].AvailableRestoreTimes)
	assert.Empty(t, model.Databases[2].AvailableRestoreTimes)
	assert.Equal(t, types.StringValue("valkey"), model.Databases[3].Engine)
	assert.Equal(t, []types.String{
		types.StringValue(restoreTime.Format(time.RFC3339)),
		types.StringValue(currentTime.Format(time.RFC3339)),
	}, model.Databases[3].AvailableRestoreTimes)
}
