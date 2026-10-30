//go:build unit

package databasevalkeyv2

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModel_FlattenAvailableRestoreTimes(t *testing.T) {
	timestamp := time.Date(2026, time.October, 5, 12, 30, 0, 0, time.UTC)
	tests := []struct {
		name           string
		availableTimes []time.Time
		wantNull       bool
		want           []string
	}{
		{name: "nil is null", wantNull: true},
		{name: "empty stays empty", availableTimes: []time.Time{}, want: []string{}},
		{name: "populated", availableTimes: []time.Time{timestamp}, want: []string{timestamp.Format(time.RFC3339)}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			model := &Model{}
			database := &linodego.ValkeyDatabase{AvailableRestoreTimes: test.availableTimes}

			diags := model.Flatten(t.Context(), database, nil, nil, false)
			require.False(t, diags.HasError(), "%v", diags)
			if test.wantNull {
				assert.True(t, model.AvailableRestoreTimes.IsNull())
				return
			}

			var got []string
			require.False(t, model.AvailableRestoreTimes.ElementsAs(t.Context(), &got, false).HasError())
			assert.Equal(t, test.want, got)
		})
	}
}

func TestModel_GetFork(t *testing.T) {
	restoreTime := time.Date(2026, time.October, 5, 12, 30, 0, 0, time.UTC)
	tests := []struct {
		name       string
		model      Model
		wantNil    bool
		wantSource int
		wantTime   *time.Time
	}{
		{name: "no fork inputs", model: Model{}, wantNil: true},
		{name: "source only", model: Model{ForkSource: types.Int64Value(123)}, wantSource: 123},
		{
			name:       "source and snapshot",
			model:      Model{ForkSource: types.Int64Value(123), ForkRestoreTime: timetypes.NewRFC3339TimeValue(restoreTime)},
			wantSource: 123,
			wantTime:   &restoreTime,
		},
		{name: "snapshot requires source", model: Model{ForkRestoreTime: timetypes.NewRFC3339TimeValue(restoreTime)}, wantNil: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var diags diag.Diagnostics
			fork := test.model.GetFork(&diags)
			if test.name == "snapshot requires source" {
				require.True(t, diags.HasError())
				assert.Nil(t, fork)
				return
			}
			require.False(t, diags.HasError(), "%v", diags)
			if test.wantNil {
				assert.Nil(t, fork)
				return
			}
			require.NotNil(t, fork)
			assert.Equal(t, test.wantSource, fork.Source)
			if test.wantTime == nil {
				assert.Nil(t, fork.RestoreTime)
			} else {
				require.NotNil(t, fork.RestoreTime)
				assert.True(t, test.wantTime.Equal(*fork.RestoreTime))
			}
		})
	}
}

func TestModel_GetEngineConfigPreservesExplicitZeroValuesAndOmitsNull(t *testing.T) {
	model := Model{
		EngineConfigBackupHour:            types.Int64Value(0),
		EngineConfigBackupMinute:          types.Int64Null(),
		EngineConfigFrequentSnapshots:     types.BoolValue(false),
		EngineConfigValkeyMaxmemoryPolicy: types.StringNull(),
	}
	var diags diag.Diagnostics
	config := model.GetEngineConfig(&diags)
	require.False(t, diags.HasError(), "%v", diags)
	require.NotNil(t, config.BackupHour)
	require.NotNil(t, *config.BackupHour)
	assert.Equal(t, 0, **config.BackupHour)
	assert.Nil(t, config.BackupMinute)
	require.NotNil(t, config.FrequentSnapshots)
	assert.False(t, *config.FrequentSnapshots)
	assert.Nil(t, config.ValkeyMaxmemoryPolicy)

	encoded, err := json.Marshal(config)
	require.NoError(t, err)
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(encoded, &fields))
	assert.JSONEq(t, "0", string(fields["backup_hour"]))
	assert.JSONEq(t, "false", string(fields["frequent_snapshots"]))
	assert.NotContains(t, fields, "backup_minute")
	assert.NotContains(t, fields, "valkey_maxmemory_policy")
}

func TestModel_GetEngineConfigOmitsUnsetConfig(t *testing.T) {
	var diags diag.Diagnostics
	assert.Nil(t, (&Model{}).GetEngineConfig(&diags))
	assert.False(t, diags.HasError(), "%v", diags)
}

func TestModel_GetEngineConfigOmitsNullAndUnknownBackupTimes(t *testing.T) {
	model := Model{
		EngineConfigBackupHour:   types.Int64Null(),
		EngineConfigBackupMinute: types.Int64Unknown(),
	}
	var diags diag.Diagnostics
	assert.Nil(t, model.GetEngineConfig(&diags))
	assert.False(t, diags.HasError(), "%v", diags)
}

func TestModel_GetEngineConfigOmitsUnknownValues(t *testing.T) {
	model := Model{
		EngineConfigBackupHour:                          types.Int64Unknown(),
		EngineConfigBackupMinute:                        types.Int64Unknown(),
		EngineConfigFrequentSnapshots:                   types.BoolUnknown(),
		EngineConfigValkeyACLChannelsDefault:            types.StringUnknown(),
		EngineConfigValkeyActiveExpireEffort:            types.Int64Unknown(),
		EngineConfigValkeyActiveDefrag:                  types.BoolUnknown(),
		EngineConfigValkeyLFUDecayTime:                  types.Int64Unknown(),
		EngineConfigValkeyLFULogFactor:                  types.Int64Unknown(),
		EngineConfigValkeyMaxmemoryPolicy:               types.StringUnknown(),
		EngineConfigValkeyNumberOfDatabases:             types.Int64Unknown(),
		EngineConfigValkeyPersistence:                   types.StringUnknown(),
		EngineConfigValkeyPubsubClientOutputBufferLimit: types.Int64Unknown(),
		EngineConfigValkeyTimeout:                       types.Int64Unknown(),
	}
	var diags diag.Diagnostics
	assert.Nil(t, model.GetEngineConfig(&diags))
	assert.False(t, diags.HasError(), "%v", diags)
}

func TestModel_CopyFromCopiesRestoreTimes(t *testing.T) {
	list, diags := types.ListValueFrom(t.Context(), types.StringType, []string{"2026-10-05T12:30:00Z"})
	require.False(t, diags.HasError())
	original := Model{AvailableRestoreTimes: list}
	var copied Model
	copied.CopyFrom(&original, false)
	assert.Equal(t, original.AvailableRestoreTimes, copied.AvailableRestoreTimes)
}

func TestModel_FlattenRefreshesAvailableRestoreTimes(t *testing.T) {
	previous, diags := types.ListValueFrom(t.Context(), types.StringType, []string{"2026-10-05T12:30:00Z"})
	require.False(t, diags.HasError())

	currentTime := time.Date(2026, time.October, 6, 12, 30, 0, 0, time.UTC)
	model := Model{AvailableRestoreTimes: previous}
	diags = model.Flatten(t.Context(), &linodego.ValkeyDatabase{
		AvailableRestoreTimes: []time.Time{currentTime},
	}, nil, nil, true)
	require.False(t, diags.HasError(), "%v", diags)

	var got []string
	require.False(t, model.AvailableRestoreTimes.ElementsAs(t.Context(), &got, false).HasError())
	assert.Equal(t, []string{currentTime.Format(time.RFC3339)}, got)
}

func TestModel_Flatten(t *testing.T) {
	created := time.Date(2026, time.October, 5, 12, 30, 0, 0, time.UTC)
	updated := created.Add(time.Hour)
	restoreTime := created.Add(-time.Hour)
	database := &linodego.ValkeyDatabase{
		ID:            123,
		Status:        linodego.DatabaseStatusActive,
		Label:         "valkey-test",
		Region:        "us-mia",
		Type:          "g7-dedicated-4-2",
		Engine:        "valkey",
		Version:       "8.1.0",
		ClusterSize:   3,
		Platform:      "rdbms-default",
		SSLConnection: true,
		Encrypted:     true,
		AllowList:     []string{"10.0.0.1/32"},
		Created:       &created,
		Updated:       &updated,
		Fork: &linodego.DatabaseFork{
			Source:      456,
			RestoreTime: &restoreTime,
		},
		AvailableRestoreTimes: []time.Time{restoreTime, created},
		Port:                  6379,
		EngineConfig: linodego.ValkeyDatabaseEngineConfig{
			BackupHour:                          linodego.Pointer(linodego.Pointer(1)),
			BackupMinute:                        linodego.Pointer(linodego.Pointer(2)),
			FrequentSnapshots:                   linodego.Pointer(true),
			ValkeyACLChannelsDefault:            linodego.Pointer("allchannels"),
			ValkeyActiveExpireEffort:            linodego.Pointer(3),
			ValkeyActiveDefrag:                  linodego.Pointer(true),
			ValkeyLFUDecayTime:                  linodego.Pointer(4),
			ValkeyLFULogFactor:                  linodego.Pointer(5),
			ValkeyMaxmemoryPolicy:               linodego.Pointer(linodego.Pointer("allkeys-lru")),
			ValkeyNumberOfDatabases:             linodego.Pointer(6),
			ValkeyPersistence:                   linodego.Pointer("rdb"),
			ValkeyPubsubClientOutputBufferLimit: linodego.Pointer(7),
			ValkeyTimeout:                       linodego.Pointer(8),
		},
	}
	ssl := &linodego.ValkeyDatabaseSSL{CACertificate: []byte("certificate")}
	credentials := &linodego.ValkeyDatabaseCredential{Username: "default", Password: "secret"}
	var model Model

	diags := model.Flatten(t.Context(), database, ssl, credentials, false)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "123", model.ID.ValueString())
	assert.Equal(t, "valkey/8.1", model.EngineID.ValueString())
	assert.Equal(t, "2026-10-05T11:30:00Z", model.ForkRestoreTime.ValueString())
	assert.Equal(t, int64(456), model.ForkSource.ValueInt64())
	assert.Equal(t, "certificate", model.CACert.ValueString())
	assert.Equal(t, "default", model.RootUsername.ValueString())
	assert.Equal(t, "secret", model.RootPassword.ValueString())
	assert.Equal(t, int64(1), model.EngineConfigBackupHour.ValueInt64())
	assert.Equal(t, int64(2), model.EngineConfigBackupMinute.ValueInt64())
	assert.True(t, model.EngineConfigFrequentSnapshots.ValueBool())
	assert.Equal(t, "allchannels", model.EngineConfigValkeyACLChannelsDefault.ValueString())
	assert.Equal(t, int64(3), model.EngineConfigValkeyActiveExpireEffort.ValueInt64())
	assert.True(t, model.EngineConfigValkeyActiveDefrag.ValueBool())
	assert.Equal(t, int64(4), model.EngineConfigValkeyLFUDecayTime.ValueInt64())
	assert.Equal(t, int64(5), model.EngineConfigValkeyLFULogFactor.ValueInt64())
	assert.Equal(t, "allkeys-lru", model.EngineConfigValkeyMaxmemoryPolicy.ValueString())
	assert.Equal(t, int64(6), model.EngineConfigValkeyNumberOfDatabases.ValueInt64())
	assert.Equal(t, "rdb", model.EngineConfigValkeyPersistence.ValueString())
	assert.Equal(t, int64(7), model.EngineConfigValkeyPubsubClientOutputBufferLimit.ValueInt64())
	assert.Equal(t, int64(8), model.EngineConfigValkeyTimeout.ValueInt64())

	var restoreTimes []string
	require.False(t, model.AvailableRestoreTimes.ElementsAs(t.Context(), &restoreTimes, false).HasError())
	assert.Equal(t, []string{restoreTime.Format(time.RFC3339), created.Format(time.RFC3339)}, restoreTimes)

	var copied Model
	copied.CopyFrom(&model, false)
	assert.Equal(t, model, copied)
}

func TestModel_FlattenPreservesKnownEngineConfig(t *testing.T) {
	database := &linodego.ValkeyDatabase{EngineConfig: linodego.ValkeyDatabaseEngineConfig{
		BackupHour:                          linodego.Pointer(linodego.Pointer(1)),
		BackupMinute:                        linodego.Pointer(linodego.Pointer(2)),
		FrequentSnapshots:                   linodego.Pointer(false),
		ValkeyACLChannelsDefault:            linodego.Pointer("allchannels"),
		ValkeyActiveExpireEffort:            linodego.Pointer(3),
		ValkeyActiveDefrag:                  linodego.Pointer(false),
		ValkeyLFUDecayTime:                  linodego.Pointer(4),
		ValkeyLFULogFactor:                  linodego.Pointer(5),
		ValkeyMaxmemoryPolicy:               linodego.Pointer(linodego.Pointer("allkeys-lru")),
		ValkeyNumberOfDatabases:             linodego.Pointer(6),
		ValkeyPersistence:                   linodego.Pointer("rdb"),
		ValkeyPubsubClientOutputBufferLimit: linodego.Pointer(7),
		ValkeyTimeout:                       linodego.Pointer(8),
	}}
	model := Model{
		EngineConfigBackupHour:                          types.Int64Value(10),
		EngineConfigBackupMinute:                        types.Int64Value(20),
		EngineConfigFrequentSnapshots:                   types.BoolValue(true),
		EngineConfigValkeyACLChannelsDefault:            types.StringValue("resetchannels"),
		EngineConfigValkeyActiveExpireEffort:            types.Int64Value(30),
		EngineConfigValkeyActiveDefrag:                  types.BoolValue(true),
		EngineConfigValkeyLFUDecayTime:                  types.Int64Value(40),
		EngineConfigValkeyLFULogFactor:                  types.Int64Value(50),
		EngineConfigValkeyMaxmemoryPolicy:               types.StringValue("volatile-lru"),
		EngineConfigValkeyNumberOfDatabases:             types.Int64Value(60),
		EngineConfigValkeyPersistence:                   types.StringValue("off"),
		EngineConfigValkeyPubsubClientOutputBufferLimit: types.Int64Value(70),
		EngineConfigValkeyTimeout:                       types.Int64Value(80),
	}

	diags := model.Flatten(t.Context(), database, nil, nil, true)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, int64(10), model.EngineConfigBackupHour.ValueInt64())
	assert.Equal(t, int64(20), model.EngineConfigBackupMinute.ValueInt64())
	assert.True(t, model.EngineConfigFrequentSnapshots.ValueBool())
	assert.Equal(t, "resetchannels", model.EngineConfigValkeyACLChannelsDefault.ValueString())
	assert.Equal(t, int64(30), model.EngineConfigValkeyActiveExpireEffort.ValueInt64())
	assert.True(t, model.EngineConfigValkeyActiveDefrag.ValueBool())
	assert.Equal(t, int64(40), model.EngineConfigValkeyLFUDecayTime.ValueInt64())
	assert.Equal(t, int64(50), model.EngineConfigValkeyLFULogFactor.ValueInt64())
	assert.Equal(t, "volatile-lru", model.EngineConfigValkeyMaxmemoryPolicy.ValueString())
	assert.Equal(t, int64(60), model.EngineConfigValkeyNumberOfDatabases.ValueInt64())
	assert.Equal(t, "off", model.EngineConfigValkeyPersistence.ValueString())
	assert.Equal(t, int64(70), model.EngineConfigValkeyPubsubClientOutputBufferLimit.ValueInt64())
	assert.Equal(t, int64(80), model.EngineConfigValkeyTimeout.ValueInt64())
}

func TestModel_RefreshSuspendedPreservesCredentials(t *testing.T) {
	requestedSensitiveEndpoints := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/ssl") || strings.HasSuffix(r.URL.Path, "/credentials") {
			requestedSensitiveEndpoints = true
			http.Error(w, "endpoint unavailable while suspended", http.StatusConflict)
			return
		}
		if !strings.HasSuffix(r.URL.Path, "/instances/123") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":123,"status":"suspended","label":"valkey-test","engine":"valkey","version":"8.1.0"}`))
	}))
	defer server.Close()

	clientValue, err := linodego.NewClient(server.Client())
	require.NoError(t, err)
	clientValue.SetBaseURL(server.URL)
	model := Model{
		CACert:       types.StringValue("existing-certificate"),
		RootUsername: types.StringValue("existing-user"),
		RootPassword: types.StringValue("existing-secret"),
	}

	diags := model.Refresh(t.Context(), &clientValue, 123, true)
	require.False(t, diags.HasError(), "%v", diags)
	assert.False(t, requestedSensitiveEndpoints)
	assert.Equal(t, "existing-certificate", model.CACert.ValueString())
	assert.Equal(t, "existing-user", model.RootUsername.ValueString())
	assert.Equal(t, "existing-secret", model.RootPassword.ValueString())
}
