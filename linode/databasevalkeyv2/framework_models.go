package databasevalkeyv2

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/databaseshared"
)

type ResourceModel struct {
	Model
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type Model struct {
	ID                    types.String      `tfsdk:"id"`
	AllowList             types.Set         `tfsdk:"allow_list"`
	CACert                types.String      `tfsdk:"ca_cert"`
	ClusterSize           types.Int64       `tfsdk:"cluster_size"`
	Created               timetypes.RFC3339 `tfsdk:"created"`
	Encrypted             types.Bool        `tfsdk:"encrypted"`
	Engine                types.String      `tfsdk:"engine"`
	EngineID              types.String      `tfsdk:"engine_id"`
	ForkRestoreTime       timetypes.RFC3339 `tfsdk:"fork_restore_time"`
	ForkSource            types.Int64       `tfsdk:"fork_source"`
	HostPrimary           types.String      `tfsdk:"host_primary"`
	HostSecondary         types.String      `tfsdk:"host_secondary"`
	HostStandby           types.String      `tfsdk:"host_standby"`
	Label                 types.String      `tfsdk:"label"`
	Members               types.Map         `tfsdk:"members"`
	Platform              types.String      `tfsdk:"platform"`
	Port                  types.Int64       `tfsdk:"port"`
	PrivateNetwork        types.Object      `tfsdk:"private_network"`
	Region                types.String      `tfsdk:"region"`
	RootPassword          types.String      `tfsdk:"root_password"`
	RootUsername          types.String      `tfsdk:"root_username"`
	SSLConnection         types.Bool        `tfsdk:"ssl_connection"`
	Status                types.String      `tfsdk:"status"`
	Suspended             types.Bool        `tfsdk:"suspended"`
	Type                  types.String      `tfsdk:"type"`
	Updated               timetypes.RFC3339 `tfsdk:"updated"`
	Updates               types.Object      `tfsdk:"updates"`
	Version               types.String      `tfsdk:"version"`
	AvailableRestoreTimes types.List        `tfsdk:"available_restore_times"`
	PendingUpdates        types.Set         `tfsdk:"pending_updates"`

	EngineConfigBackupHour                          types.Int64  `tfsdk:"engine_config_backup_hour"`
	EngineConfigBackupMinute                        types.Int64  `tfsdk:"engine_config_backup_minute"`
	EngineConfigFrequentSnapshots                   types.Bool   `tfsdk:"engine_config_frequent_snapshots"`
	EngineConfigValkeyACLChannelsDefault            types.String `tfsdk:"engine_config_valkey_acl_channels_default"`
	EngineConfigValkeyActiveExpireEffort            types.Int64  `tfsdk:"engine_config_valkey_active_expire_effort"`
	EngineConfigValkeyActiveDefrag                  types.Bool   `tfsdk:"engine_config_valkey_activedefrag"`
	EngineConfigValkeyLFUDecayTime                  types.Int64  `tfsdk:"engine_config_valkey_lfu_decay_time"`
	EngineConfigValkeyLFULogFactor                  types.Int64  `tfsdk:"engine_config_valkey_lfu_log_factor"`
	EngineConfigValkeyMaxmemoryPolicy               types.String `tfsdk:"engine_config_valkey_maxmemory_policy"`
	EngineConfigValkeyNumberOfDatabases             types.Int64  `tfsdk:"engine_config_valkey_number_of_databases"`
	EngineConfigValkeyPersistence                   types.String `tfsdk:"engine_config_valkey_persistence"`
	EngineConfigValkeyPubsubClientOutputBufferLimit types.Int64  `tfsdk:"engine_config_valkey_pubsub_client_output_buffer_limit"`
	EngineConfigValkeyTimeout                       types.Int64  `tfsdk:"engine_config_valkey_timeout"`
}

func (m *Model) Refresh(ctx context.Context, client *linodego.Client, dbID int, preserveKnown bool) diag.Diagnostics {
	tflog.SetField(ctx, "id", dbID)
	tflog.Debug(ctx, "Refreshing the Valkey database")

	db, err := client.GetValkeyDatabase(ctx, dbID)
	if err != nil {
		var d diag.Diagnostics
		d.AddError("Failed to refresh Valkey database", err.Error())
		return d
	}

	var ssl *linodego.ValkeyDatabaseSSL
	var creds *linodego.ValkeyDatabaseCredential
	if !databaseshared.StatusIsSuspended(db.Status) {
		ssl, err = client.GetValkeyDatabaseSSL(ctx, dbID)
		if err != nil {
			var d diag.Diagnostics
			d.AddError("Failed to refresh Valkey database SSL", err.Error())
			return d
		}

		creds, err = client.GetValkeyDatabaseCredentials(ctx, dbID)
		if err != nil {
			var d diag.Diagnostics
			d.AddError("Failed to refresh Valkey database credentials", err.Error())
			return d
		}
	}

	return m.Flatten(ctx, db, ssl, creds, preserveKnown)
}

func (m *Model) Flatten(
	ctx context.Context,
	db *linodego.ValkeyDatabase,
	ssl *linodego.ValkeyDatabaseSSL,
	creds *linodego.ValkeyDatabaseCredential,
	preserveKnown bool,
) (d diag.Diagnostics) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(db.ID), preserveKnown)
	m.AllowList = helper.KeepOrUpdateSet(types.StringType, m.AllowList, helper.StringSliceToFrameworkValueSlice(db.AllowList), preserveKnown, &d)
	m.ClusterSize = helper.KeepOrUpdateInt64(m.ClusterSize, int64(db.ClusterSize), preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, timetypes.NewRFC3339TimePointerValue(db.Created), preserveKnown)
	m.Encrypted = helper.KeepOrUpdateBool(m.Encrypted, db.Encrypted, preserveKnown)
	m.Engine = helper.KeepOrUpdateString(m.Engine, db.Engine, preserveKnown)
	m.EngineID = helper.KeepOrUpdateString(m.EngineID, databaseshared.CreateDatabaseEngineSlug(db.Engine, db.Version), preserveKnown)
	var forkRestoreTime *time.Time
	if db.Fork != nil {
		forkRestoreTime = db.Fork.RestoreTime
	}
	m.ForkRestoreTime = helper.KeepOrUpdateValue(m.ForkRestoreTime, timetypes.NewRFC3339TimePointerValue(forkRestoreTime), preserveKnown)
	m.HostPrimary = helper.KeepOrUpdateString(m.HostPrimary, db.Hosts.Primary, preserveKnown)
	m.HostSecondary = helper.KeepOrUpdateString(m.HostSecondary, db.Hosts.Standby, preserveKnown)
	m.HostStandby = helper.KeepOrUpdateString(m.HostStandby, db.Hosts.Standby, preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, db.Label, preserveKnown)
	m.Platform = helper.KeepOrUpdateString(m.Platform, string(db.Platform), preserveKnown)
	m.Port = helper.KeepOrUpdateInt64(m.Port, int64(db.Port), preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, db.Region, preserveKnown)
	m.SSLConnection = helper.KeepOrUpdateBool(m.SSLConnection, db.SSLConnection, preserveKnown)
	m.Status = helper.KeepOrUpdateString(m.Status, string(db.Status), preserveKnown)
	m.Suspended = helper.KeepOrUpdateBool(m.Suspended, databaseshared.StatusIsSuspended(db.Status), preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, db.Type, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, timetypes.NewRFC3339TimePointerValue(db.Updated), preserveKnown)
	m.Version = helper.KeepOrUpdateString(m.Version, db.Version, preserveKnown)

	if db.Fork == nil {
		m.ForkSource = helper.KeepOrUpdateValue(m.ForkSource, types.Int64Null(), preserveKnown)
	} else {
		m.ForkSource = helper.KeepOrUpdateInt64(m.ForkSource, int64(db.Fork.Source), preserveKnown)
	}

	if db.AvailableRestoreTimes == nil {
		m.AvailableRestoreTimes = types.ListNull(types.StringType)
	} else {
		restoreTimes := make([]string, len(db.AvailableRestoreTimes))
		for i, restoreTime := range db.AvailableRestoreTimes {
			restoreTimes[i] = restoreTime.Format(time.RFC3339)
		}
		restoreTimesValue, rd := types.ListValueFrom(ctx, types.StringType, restoreTimes)
		d.Append(rd...)
		m.AvailableRestoreTimes = restoreTimesValue
	}

	if ssl != nil {
		m.CACert = helper.KeepOrUpdateString(m.CACert, string(ssl.CACertificate), preserveKnown)
	} else {
		m.CACert = helper.KeepOrUpdateValue(m.CACert, types.StringNull(), true)
	}
	if creds != nil {
		m.RootPassword = helper.KeepOrUpdateString(m.RootPassword, creds.Password, preserveKnown)
		m.RootUsername = helper.KeepOrUpdateString(m.RootUsername, creds.Username, preserveKnown)
	} else {
		m.RootPassword = helper.KeepOrUpdateValue(m.RootPassword, types.StringNull(), true)
		m.RootUsername = helper.KeepOrUpdateValue(m.RootUsername, types.StringNull(), true)
	}

	members := helper.MapMap(db.Members, func(key string, value linodego.DatabaseMemberType) (string, string) {
		return key, string(value)
	})
	m.Members = helper.KeepOrUpdateStringMap(ctx, m.Members, members, preserveKnown, &d)
	if d.HasError() {
		return d
	}

	if db.PrivateNetwork == nil {
		m.PrivateNetwork = helper.KeepOrUpdateValue(m.PrivateNetwork, types.ObjectNull(databaseshared.ObjectTypePrivateNetwork.AttrTypes), preserveKnown)
	} else {
		privateNetwork, rd := databaseshared.FlattenPrivateNetwork(ctx, *db.PrivateNetwork)
		d.Append(rd...)
		m.PrivateNetwork = helper.KeepOrUpdateValue(m.PrivateNetwork, privateNetwork, preserveKnown)
	}

	updates, rd := databaseshared.FlattenUpdates(ctx, db.Updates)
	d.Append(rd...)
	m.Updates = helper.KeepOrUpdateValue(m.Updates, updates, preserveKnown)
	pendingUpdates, rd := databaseshared.FlattenPendingUpdates(ctx, db.Updates.Pending)
	d.Append(rd...)
	m.PendingUpdates = helper.KeepOrUpdateValue(m.PendingUpdates, pendingUpdates, preserveKnown)

	config := db.EngineConfig
	var backupHour, backupMinute *int
	if config.BackupHour != nil {
		backupHour = *config.BackupHour
	}
	if config.BackupMinute != nil {
		backupMinute = *config.BackupMinute
	}
	m.EngineConfigBackupHour = helper.KeepOrUpdateIntPointer(m.EngineConfigBackupHour, backupHour, preserveKnown)
	m.EngineConfigBackupMinute = helper.KeepOrUpdateIntPointer(m.EngineConfigBackupMinute, backupMinute, preserveKnown)
	m.EngineConfigFrequentSnapshots = helper.KeepOrUpdateBoolPointer(m.EngineConfigFrequentSnapshots, config.FrequentSnapshots, preserveKnown)
	m.EngineConfigValkeyACLChannelsDefault = helper.KeepOrUpdateStringPointer(
		m.EngineConfigValkeyACLChannelsDefault,
		config.ValkeyACLChannelsDefault,
		preserveKnown,
	)
	m.EngineConfigValkeyActiveExpireEffort = helper.KeepOrUpdateIntPointer(
		m.EngineConfigValkeyActiveExpireEffort,
		config.ValkeyActiveExpireEffort,
		preserveKnown,
	)
	m.EngineConfigValkeyActiveDefrag = helper.KeepOrUpdateBoolPointer(m.EngineConfigValkeyActiveDefrag, config.ValkeyActiveDefrag, preserveKnown)
	m.EngineConfigValkeyLFUDecayTime = helper.KeepOrUpdateIntPointer(m.EngineConfigValkeyLFUDecayTime, config.ValkeyLFUDecayTime, preserveKnown)
	m.EngineConfigValkeyLFULogFactor = helper.KeepOrUpdateIntPointer(m.EngineConfigValkeyLFULogFactor, config.ValkeyLFULogFactor, preserveKnown)
	var maxmemoryPolicy *string
	if config.ValkeyMaxmemoryPolicy != nil {
		maxmemoryPolicy = *config.ValkeyMaxmemoryPolicy
	}
	m.EngineConfigValkeyMaxmemoryPolicy = helper.KeepOrUpdateStringPointer(m.EngineConfigValkeyMaxmemoryPolicy, maxmemoryPolicy, preserveKnown)
	m.EngineConfigValkeyNumberOfDatabases = helper.KeepOrUpdateIntPointer(m.EngineConfigValkeyNumberOfDatabases, config.ValkeyNumberOfDatabases, preserveKnown)
	m.EngineConfigValkeyPersistence = helper.KeepOrUpdateStringPointer(m.EngineConfigValkeyPersistence, config.ValkeyPersistence, preserveKnown)
	m.EngineConfigValkeyPubsubClientOutputBufferLimit = helper.KeepOrUpdateIntPointer(
		m.EngineConfigValkeyPubsubClientOutputBufferLimit,
		config.ValkeyPubsubClientOutputBufferLimit,
		preserveKnown,
	)
	m.EngineConfigValkeyTimeout = helper.KeepOrUpdateIntPointer(m.EngineConfigValkeyTimeout, config.ValkeyTimeout, preserveKnown)

	return d
}

func (m *Model) CopyFrom(other *Model, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.AllowList = helper.KeepOrUpdateValue(m.AllowList, other.AllowList, preserveKnown)
	m.CACert = helper.KeepOrUpdateValue(m.CACert, other.CACert, preserveKnown)
	m.ClusterSize = helper.KeepOrUpdateValue(m.ClusterSize, other.ClusterSize, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Encrypted = helper.KeepOrUpdateValue(m.Encrypted, other.Encrypted, preserveKnown)
	m.Engine = helper.KeepOrUpdateValue(m.Engine, other.Engine, preserveKnown)
	m.EngineID = helper.KeepOrUpdateValue(m.EngineID, other.EngineID, preserveKnown)
	m.ForkRestoreTime = helper.KeepOrUpdateValue(m.ForkRestoreTime, other.ForkRestoreTime, preserveKnown)
	m.ForkSource = helper.KeepOrUpdateValue(m.ForkSource, other.ForkSource, preserveKnown)
	m.HostPrimary = helper.KeepOrUpdateValue(m.HostPrimary, other.HostPrimary, preserveKnown)
	m.HostSecondary = helper.KeepOrUpdateValue(m.HostSecondary, other.HostSecondary, preserveKnown)
	m.HostStandby = helper.KeepOrUpdateValue(m.HostStandby, other.HostStandby, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Members = helper.KeepOrUpdateValue(m.Members, other.Members, preserveKnown)
	m.Platform = helper.KeepOrUpdateValue(m.Platform, other.Platform, preserveKnown)
	m.Port = helper.KeepOrUpdateValue(m.Port, other.Port, preserveKnown)
	m.PrivateNetwork = helper.KeepOrUpdateValue(m.PrivateNetwork, other.PrivateNetwork, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
	m.RootPassword = helper.KeepOrUpdateValue(m.RootPassword, other.RootPassword, preserveKnown)
	m.RootUsername = helper.KeepOrUpdateValue(m.RootUsername, other.RootUsername, preserveKnown)
	m.SSLConnection = helper.KeepOrUpdateValue(m.SSLConnection, other.SSLConnection, preserveKnown)
	m.Status = helper.KeepOrUpdateValue(m.Status, other.Status, preserveKnown)
	m.Suspended = helper.KeepOrUpdateValue(m.Suspended, other.Suspended, preserveKnown)
	m.Type = helper.KeepOrUpdateValue(m.Type, other.Type, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Updates = helper.KeepOrUpdateValue(m.Updates, other.Updates, preserveKnown)
	m.Version = helper.KeepOrUpdateValue(m.Version, other.Version, preserveKnown)
	m.AvailableRestoreTimes = helper.KeepOrUpdateValue(m.AvailableRestoreTimes, other.AvailableRestoreTimes, preserveKnown)
	m.PendingUpdates = helper.KeepOrUpdateValue(m.PendingUpdates, other.PendingUpdates, preserveKnown)
	m.EngineConfigBackupHour = helper.KeepOrUpdateValue(m.EngineConfigBackupHour, other.EngineConfigBackupHour, preserveKnown)
	m.EngineConfigBackupMinute = helper.KeepOrUpdateValue(m.EngineConfigBackupMinute, other.EngineConfigBackupMinute, preserveKnown)
	m.EngineConfigFrequentSnapshots = helper.KeepOrUpdateValue(m.EngineConfigFrequentSnapshots, other.EngineConfigFrequentSnapshots, preserveKnown)
	m.EngineConfigValkeyACLChannelsDefault = helper.KeepOrUpdateValue(
		m.EngineConfigValkeyACLChannelsDefault,
		other.EngineConfigValkeyACLChannelsDefault,
		preserveKnown,
	)
	m.EngineConfigValkeyActiveExpireEffort = helper.KeepOrUpdateValue(
		m.EngineConfigValkeyActiveExpireEffort,
		other.EngineConfigValkeyActiveExpireEffort,
		preserveKnown,
	)
	m.EngineConfigValkeyActiveDefrag = helper.KeepOrUpdateValue(m.EngineConfigValkeyActiveDefrag, other.EngineConfigValkeyActiveDefrag, preserveKnown)
	m.EngineConfigValkeyLFUDecayTime = helper.KeepOrUpdateValue(m.EngineConfigValkeyLFUDecayTime, other.EngineConfigValkeyLFUDecayTime, preserveKnown)
	m.EngineConfigValkeyLFULogFactor = helper.KeepOrUpdateValue(m.EngineConfigValkeyLFULogFactor, other.EngineConfigValkeyLFULogFactor, preserveKnown)
	m.EngineConfigValkeyMaxmemoryPolicy = helper.KeepOrUpdateValue(m.EngineConfigValkeyMaxmemoryPolicy, other.EngineConfigValkeyMaxmemoryPolicy, preserveKnown)
	m.EngineConfigValkeyNumberOfDatabases = helper.KeepOrUpdateValue(
		m.EngineConfigValkeyNumberOfDatabases,
		other.EngineConfigValkeyNumberOfDatabases,
		preserveKnown,
	)
	m.EngineConfigValkeyPersistence = helper.KeepOrUpdateValue(m.EngineConfigValkeyPersistence, other.EngineConfigValkeyPersistence, preserveKnown)
	m.EngineConfigValkeyPubsubClientOutputBufferLimit = helper.KeepOrUpdateValue(
		m.EngineConfigValkeyPubsubClientOutputBufferLimit,
		other.EngineConfigValkeyPubsubClientOutputBufferLimit,
		preserveKnown,
	)
	m.EngineConfigValkeyTimeout = helper.KeepOrUpdateValue(m.EngineConfigValkeyTimeout, other.EngineConfigValkeyTimeout, preserveKnown)
}

func (m *Model) GetFork(d *diag.Diagnostics) *linodego.DatabaseFork {
	var fork linodego.DatabaseFork
	isSpecified := false
	if !m.ForkSource.IsUnknown() && !m.ForkSource.IsNull() {
		fork.Source = helper.FrameworkSafeInt64ToInt(m.ForkSource.ValueInt64(), d)
		isSpecified = true
	}
	if !m.ForkRestoreTime.IsUnknown() && !m.ForkRestoreTime.IsNull() {
		if m.ForkSource.IsNull() || m.ForkSource.IsUnknown() {
			d.AddError("Missing Valkey fork source", "fork_source must be set when fork_restore_time is specified.")
			return nil
		}
		restoreTime, rd := m.ForkRestoreTime.ValueRFC3339Time()
		d.Append(rd...)
		fork.RestoreTime = &restoreTime
		isSpecified = true
	}
	if d.HasError() || !isSpecified {
		return nil
	}
	return &fork
}

func (m *Model) GetAllowList(ctx context.Context, d *diag.Diagnostics) []string {
	if m.AllowList.IsUnknown() || m.AllowList.IsNull() {
		return nil
	}
	var allowList []string
	d.Append(m.AllowList.ElementsAs(ctx, &allowList, false)...)
	return allowList
}

func (m *Model) GetUpdates(ctx context.Context, d *diag.Diagnostics) *databaseshared.ModelUpdates {
	if m.Updates.IsUnknown() || m.Updates.IsNull() {
		return nil
	}
	var updates databaseshared.ModelUpdates
	d.Append(m.Updates.As(ctx, &updates, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true})...)
	return &updates
}

func (m *Model) GetPrivateNetwork(ctx context.Context, d *diag.Diagnostics) *databaseshared.PrivateNetworkModel {
	if m.PrivateNetwork.IsUnknown() || m.PrivateNetwork.IsNull() {
		return nil
	}
	var privateNetwork databaseshared.PrivateNetworkModel
	d.Append(m.PrivateNetwork.As(ctx, &privateNetwork, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true})...)
	return &privateNetwork
}

func (m *Model) GetEngineConfig(d *diag.Diagnostics) *linodego.ValkeyDatabaseEngineConfig {
	var config linodego.ValkeyDatabaseEngineConfig
	if !m.EngineConfigBackupHour.IsUnknown() && !m.EngineConfigBackupHour.IsNull() {
		value := helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigBackupHour.ValueInt64Pointer(), d)
		if value != nil {
			config.BackupHour = linodego.Pointer(value)
		}
	}
	if !m.EngineConfigBackupMinute.IsUnknown() && !m.EngineConfigBackupMinute.IsNull() {
		value := helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigBackupMinute.ValueInt64Pointer(), d)
		if value != nil {
			config.BackupMinute = linodego.Pointer(value)
		}
	}
	if !m.EngineConfigFrequentSnapshots.IsUnknown() {
		config.FrequentSnapshots = m.EngineConfigFrequentSnapshots.ValueBoolPointer()
	}
	if !m.EngineConfigValkeyACLChannelsDefault.IsUnknown() {
		config.ValkeyACLChannelsDefault = m.EngineConfigValkeyACLChannelsDefault.ValueStringPointer()
	}
	if !m.EngineConfigValkeyActiveExpireEffort.IsUnknown() {
		config.ValkeyActiveExpireEffort = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigValkeyActiveExpireEffort.ValueInt64Pointer(), d)
	}
	if !m.EngineConfigValkeyActiveDefrag.IsUnknown() {
		config.ValkeyActiveDefrag = m.EngineConfigValkeyActiveDefrag.ValueBoolPointer()
	}
	if !m.EngineConfigValkeyLFUDecayTime.IsUnknown() {
		config.ValkeyLFUDecayTime = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigValkeyLFUDecayTime.ValueInt64Pointer(), d)
	}
	if !m.EngineConfigValkeyLFULogFactor.IsUnknown() {
		config.ValkeyLFULogFactor = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigValkeyLFULogFactor.ValueInt64Pointer(), d)
	}
	if !m.EngineConfigValkeyMaxmemoryPolicy.IsUnknown() && !m.EngineConfigValkeyMaxmemoryPolicy.IsNull() {
		config.ValkeyMaxmemoryPolicy = linodego.Pointer(m.EngineConfigValkeyMaxmemoryPolicy.ValueStringPointer())
	}
	if !m.EngineConfigValkeyNumberOfDatabases.IsUnknown() {
		config.ValkeyNumberOfDatabases = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigValkeyNumberOfDatabases.ValueInt64Pointer(), d)
	}
	if !m.EngineConfigValkeyPersistence.IsUnknown() {
		config.ValkeyPersistence = m.EngineConfigValkeyPersistence.ValueStringPointer()
	}
	if !m.EngineConfigValkeyPubsubClientOutputBufferLimit.IsUnknown() {
		config.ValkeyPubsubClientOutputBufferLimit = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigValkeyPubsubClientOutputBufferLimit.ValueInt64Pointer(),
			d,
		)
	}
	if !m.EngineConfigValkeyTimeout.IsUnknown() {
		config.ValkeyTimeout = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigValkeyTimeout.ValueInt64Pointer(), d)
	}
	if config.BackupHour == nil && config.BackupMinute == nil && config.FrequentSnapshots == nil &&
		config.ValkeyACLChannelsDefault == nil && config.ValkeyActiveExpireEffort == nil &&
		config.ValkeyActiveDefrag == nil && config.ValkeyLFUDecayTime == nil && config.ValkeyLFULogFactor == nil &&
		config.ValkeyMaxmemoryPolicy == nil && config.ValkeyNumberOfDatabases == nil && config.ValkeyPersistence == nil &&
		config.ValkeyPubsubClientOutputBufferLimit == nil && config.ValkeyTimeout == nil {
		return nil
	}
	return &config
}
