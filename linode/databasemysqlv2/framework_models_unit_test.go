//go:build unit

package databasemysqlv2_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/databasemysqlv2"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/databaseshared"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/unit"
	"github.com/stretchr/testify/require"
)

var (
	currentTime        = time.Now()
	currentTimeFWValue = timetypes.NewRFC3339TimePointerValue(&currentTime)

	testDB = linodego.MySQLDatabase{
		ID:        12345,
		Status:    linodego.DatabaseStatusProvisioning,
		Label:     "foobar",
		Region:    "us-mia",
		Type:      "g6-nanode-1",
		Engine:    "mysql",
		Version:   "8",
		Encrypted: true,
		AllowList: []string{"0.0.0.0/0", "10.0.0.1/32"},

		// TODO: sdf
		// Port:          1234,

		SSLConnection: true,
		ClusterSize:   3,
		Hosts: linodego.DatabaseHost{
			Primary: "1.2.3.4",
			Standby: "4.3.2.1",
		},
		Updates: linodego.DatabaseMaintenanceWindow{
			DayOfWeek: 1,
			Duration:  1,
			Frequency: linodego.DatabaseMaintenanceFrequencyWeekly,
			HourOfDay: 1,
			Pending: []linodego.DatabaseMaintenanceWindowPending{
				{
					Deadline:    &currentTime,
					Description: "foobar",
					PlannedFor:  &currentTime,
				},
			},
		},
		Created: &currentTime,
		Updated: &currentTime,
		Fork: &linodego.DatabaseFork{
			Source:      12345,
			RestoreTime: &currentTime,
		},
		OldestRestoreTime: &currentTime,
		Platform:          "foobar",
		EngineConfig: linodego.MySQLDatabaseEngineConfig{
			BinlogRetentionPeriod: linodego.Pointer(600),
			MySQL: &linodego.MySQLDatabaseEngineConfigMySQL{
				ConnectTimeout:       linodego.Pointer(10),
				DefaultTimeZone:      linodego.Pointer("+03:00"),
				GroupConcatMaxLen:    linodego.Pointer(float64(1024)),
				SQLMode:              linodego.Pointer("test"),
				SQLRequirePrimaryKey: linodego.Pointer(true),
			},
		},
	}

	testDBSSL = linodego.MySQLDatabaseSSL{CACertificate: []byte("Zm9vYmFy")}

	testDBCreds = linodego.MySQLDatabaseCredential{
		Username: "foobar",
		Password: "barfoo",
	}
)

func TestModel_Flatten(t *testing.T) {
	var model databasemysqlv2.Model

	model.Flatten(context.Background(), &testDB, &testDBSSL, &testDBCreds, false)

	updates := unit.FrameworkObjectAs[databaseshared.ModelUpdates](t, model.Updates)

	require.Equal(t, "12345", model.ID.ValueString())

	require.Equal(t, "provisioning", model.Status.ValueString())
	require.Equal(t, "foobar", model.Label.ValueString())
	require.Equal(t, "us-mia", model.Region.ValueString())
	require.Equal(t, "g6-nanode-1", model.Type.ValueString())
	require.Equal(t, "mysql/8", model.EngineID.ValueString())
	require.Equal(t, "mysql", model.Engine.ValueString())
	require.Equal(t, "8", model.Version.ValueString())
	require.Equal(t, true, model.Encrypted.ValueBool())
	require.Equal(t, "foobar", model.Platform.ValueString())

	// TODO
	require.Equal(t, int64(0), model.Port.ValueInt64())

	require.Equal(t, true, model.SSLConnection.ValueBool())
	require.Equal(t, "Zm9vYmFy", model.CACert.ValueString())
	require.Equal(t, int64(12345), model.ForkSource.ValueInt64())
	require.Equal(t, currentTimeFWValue, model.ForkRestoreTime)
	require.Equal(t, "1.2.3.4", model.HostPrimary.ValueString())
	require.Equal(t, "4.3.2.1", model.HostSecondary.ValueString())
	require.Equal(t, "foobar", model.RootUsername.ValueString())
	require.Equal(t, "barfoo", model.RootPassword.ValueString())
	require.Equal(t, currentTimeFWValue, model.Created)
	require.Equal(t, currentTimeFWValue, model.Updated)
	require.Equal(t, currentTimeFWValue, model.OldestRestoreTime)

	require.Equal(t, false, model.Suspended.ValueBool())

	require.Equal(t, int64(1), updates.DayOfWeek.ValueInt64())
	require.Equal(t, int64(1), updates.Duration.ValueInt64())
	require.Equal(t, "weekly", updates.Frequency.ValueString())
	require.Equal(t, int64(1), updates.HourOfDay.ValueInt64())

	allowListElements := model.AllowList.Elements()
	require.Contains(t, allowListElements, types.StringValue("0.0.0.0/0"))
	require.Contains(t, allowListElements, types.StringValue("10.0.0.1/32"))

	expectedPendingElement, d := types.ObjectValue(
		map[string]attr.Type{
			"deadline":    timetypes.RFC3339Type{},
			"description": types.StringType,
			"planned_for": timetypes.RFC3339Type{},
		},
		map[string]attr.Value{
			"deadline":    currentTimeFWValue,
			"description": types.StringValue("foobar"),
			"planned_for": currentTimeFWValue,
		},
	)
	require.False(t, d.HasError(), d.Errors())

	require.Equal(t, "test", model.EngineConfigMySQLSQLMode.ValueString())
	require.Equal(t, true, model.EngineConfigMySQLSQLRequirePrimaryKey.ValueBool())
	require.Equal(t, int64(600), model.EngineConfigBinlogRetentionPeriod.ValueInt64())
	require.Equal(t, int64(10), model.EngineConfigMySQLConnectTimeout.ValueInt64())
	require.Equal(t, "+03:00", model.EngineConfigMySQLDefaultTimeZone.ValueString())
	require.Equal(t, float64(1024), model.EngineConfigMySQLGroupConcatMaxLen.ValueFloat64())

	require.True(t, model.PendingUpdates.Elements()[0].Equal(expectedPendingElement))
}

func TestModel_Copy(t *testing.T) {
	var modelOld, modelNew databasemysqlv2.Model
	modelOld.Flatten(context.Background(), &testDB, &testDBSSL, &testDBCreds, false)

	modelNew.CopyFrom(&modelOld, false)

	require.Equal(t, modelOld, modelNew)
}
