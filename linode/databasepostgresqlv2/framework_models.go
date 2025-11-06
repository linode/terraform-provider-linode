package databasepostgresqlv2

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type ModelPrivateNetwork struct {
	VPCID        types.Int64 `tfsdk:"vpc_id"`
	SubnetID     types.Int64 `tfsdk:"subnet_id"`
	PublicAccess types.Bool  `tfsdk:"public_access"`
}

func (m ModelPrivateNetwork) ToLinodego(d diag.Diagnostics) *linodego.DatabasePrivateNetwork {
	return &linodego.DatabasePrivateNetwork{
		VPCID:        helper.FrameworkSafeInt64ToInt(m.VPCID.ValueInt64(), &d),
		SubnetID:     helper.FrameworkSafeInt64ToInt(m.SubnetID.ValueInt64(), &d),
		PublicAccess: m.PublicAccess.ValueBool(),
	}
}

type ModelHosts struct {
	Primary   types.String `tfsdk:"primary"`
	Secondary types.String `tfsdk:"secondary"`
}

type ModelUpdates struct {
	DayOfWeek types.Int64  `tfsdk:"day_of_week"`
	Duration  types.Int64  `tfsdk:"duration"`
	Frequency types.String `tfsdk:"frequency"`
	HourOfDay types.Int64  `tfsdk:"hour_of_day"`
}

func (m ModelUpdates) ToLinodego(d diag.Diagnostics) *linodego.DatabaseMaintenanceWindow {
	return &linodego.DatabaseMaintenanceWindow{
		DayOfWeek: linodego.DatabaseDayOfWeek(helper.FrameworkSafeInt64ToInt(m.DayOfWeek.ValueInt64(), &d)),
		Duration:  helper.FrameworkSafeInt64ToInt(m.Duration.ValueInt64(), &d),
		Frequency: linodego.DatabaseMaintenanceFrequency(m.Frequency.ValueString()),
		HourOfDay: helper.FrameworkSafeInt64ToInt(m.HourOfDay.ValueInt64(), &d),
	}
}

type ModelPendingUpdate struct {
	Deadline    timetypes.RFC3339 `tfsdk:"deadline"`
	Description types.String      `tfsdk:"description"`
	PlannedFor  timetypes.RFC3339 `tfsdk:"planned_for"`
}

type ResourceModel struct {
	Model
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type Model struct {
	ID types.String `tfsdk:"id"`

	AllowList     types.Set         `tfsdk:"allow_list"`
	CACert        types.String      `tfsdk:"ca_cert"`
	ClusterSize   types.Int64       `tfsdk:"cluster_size"`
	Created       timetypes.RFC3339 `tfsdk:"created"`
	Encrypted     types.Bool        `tfsdk:"encrypted"`
	Engine        types.String      `tfsdk:"engine"`
	EngineID      types.String      `tfsdk:"engine_id"`
	HostPrimary   types.String      `tfsdk:"host_primary"`
	HostSecondary types.String      `tfsdk:"host_secondary"`
	Label         types.String      `tfsdk:"label"`
	Members       types.Map         `tfsdk:"members"`
	Platform      types.String      `tfsdk:"platform"`
	Port          types.Int64       `tfsdk:"port"`
	Region        types.String      `tfsdk:"region"`
	RootPassword  types.String      `tfsdk:"root_password"`
	RootUsername  types.String      `tfsdk:"root_username"`
	SSLConnection types.Bool        `tfsdk:"ssl_connection"`
	Status        types.String      `tfsdk:"status"`
	Suspended     types.Bool        `tfsdk:"suspended"`
	Type          types.String      `tfsdk:"type"`
	Updated       timetypes.RFC3339 `tfsdk:"updated"`
	Version       types.String      `tfsdk:"version"`

	// Fork-specific fields
	OldestRestoreTime timetypes.RFC3339 `tfsdk:"oldest_restore_time"`
	ForkSource        types.Int64       `tfsdk:"fork_source"`
	ForkRestoreTime   timetypes.RFC3339 `tfsdk:"fork_restore_time"`

	PrivateNetwork types.Object `tfsdk:"private_network"`
	Updates        types.Object `tfsdk:"updates"`
	PendingUpdates types.Set    `tfsdk:"pending_updates"`

	// EngineConfig-specific fields
	EngineConfigPGAutovacuumAnalyzeScaleFactor         types.Float64 `tfsdk:"engine_config_pg_autovacuum_analyze_scale_factor"`
	EngineConfigPGAutovacuumAnalyzeThreshold           types.Int32   `tfsdk:"engine_config_pg_autovacuum_analyze_threshold"`
	EngineConfigPGAutovacuumMaxWorkers                 types.Int64   `tfsdk:"engine_config_pg_autovacuum_max_workers"`
	EngineConfigPGAutovacuumNaptime                    types.Int64   `tfsdk:"engine_config_pg_autovacuum_naptime"`
	EngineConfigPGAutovacuumVacuumCostDelay            types.Int64   `tfsdk:"engine_config_pg_autovacuum_vacuum_cost_delay"`
	EngineConfigPGAutovacuumVacuumCostLimit            types.Int64   `tfsdk:"engine_config_pg_autovacuum_vacuum_cost_limit"`
	EngineConfigPGAutovacuumVacuumScaleFactor          types.Float64 `tfsdk:"engine_config_pg_autovacuum_vacuum_scale_factor"`
	EngineConfigPGAutovacuumVacuumThreshold            types.Int32   `tfsdk:"engine_config_pg_autovacuum_vacuum_threshold"`
	EngineConfigPGBGWriterDelay                        types.Int64   `tfsdk:"engine_config_pg_bgwriter_delay"`
	EngineConfigPGBGWriterFlushAfter                   types.Int64   `tfsdk:"engine_config_pg_bgwriter_flush_after"`
	EngineConfigPGBGWriterLRUMaxpages                  types.Int64   `tfsdk:"engine_config_pg_bgwriter_lru_maxpages"`
	EngineConfigPGBGWriterLRUMultiplier                types.Float64 `tfsdk:"engine_config_pg_bgwriter_lru_multiplier"`
	EngineConfigPGDeadlockTimeout                      types.Int64   `tfsdk:"engine_config_pg_deadlock_timeout"`
	EngineConfigPGDefaultToastCompression              types.String  `tfsdk:"engine_config_pg_default_toast_compression"`
	EngineConfigPGIdleInTransactionSessionTimeout      types.Int64   `tfsdk:"engine_config_pg_idle_in_transaction_session_timeout"`
	EngineConfigPGJIT                                  types.Bool    `tfsdk:"engine_config_pg_jit"`
	EngineConfigPGMaxFilesPerProcess                   types.Int64   `tfsdk:"engine_config_pg_max_files_per_process"`
	EngineConfigPGMaxLocksPerTransaction               types.Int64   `tfsdk:"engine_config_pg_max_locks_per_transaction"`
	EngineConfigPGMaxLogicalReplicationWorkers         types.Int64   `tfsdk:"engine_config_pg_max_logical_replication_workers"`
	EngineConfigPGMaxParallelWorkers                   types.Int64   `tfsdk:"engine_config_pg_max_parallel_workers"`
	EngineConfigPGMaxParallelWorkersPerGather          types.Int64   `tfsdk:"engine_config_pg_max_parallel_workers_per_gather"`
	EngineConfigPGMaxPredLocksPerTransaction           types.Int64   `tfsdk:"engine_config_pg_max_pred_locks_per_transaction"`
	EngineConfigPGMaxReplicationSlots                  types.Int64   `tfsdk:"engine_config_pg_max_replication_slots"`
	EngineConfigPGMaxSlotWALKeepSize                   types.Int32   `tfsdk:"engine_config_pg_max_slot_wal_keep_size"`
	EngineConfigPGMaxStackDepth                        types.Int64   `tfsdk:"engine_config_pg_max_stack_depth"`
	EngineConfigPGMaxStandbyArchiveDelay               types.Int64   `tfsdk:"engine_config_pg_max_standby_archive_delay"`
	EngineConfigPGMaxStandbyStreamingDelay             types.Int64   `tfsdk:"engine_config_pg_max_standby_streaming_delay"`
	EngineConfigPGMaxWALSenders                        types.Int64   `tfsdk:"engine_config_pg_max_wal_senders"`
	EngineConfigPGMaxWorkerProcesses                   types.Int64   `tfsdk:"engine_config_pg_max_worker_processes"`
	EngineConfigPGPasswordEncryption                   types.String  `tfsdk:"engine_config_pg_password_encryption"`
	EngineConfigPGPGPartmanBGWInterval                 types.Int64   `tfsdk:"engine_config_pg_pg_partman_bgw_interval"`
	EngineConfigPGPGPartmanBGWRole                     types.String  `tfsdk:"engine_config_pg_pg_partman_bgw_role"`
	EngineConfigPGPGStatMonitorPGSMEnableQueryPlan     types.Bool    `tfsdk:"engine_config_pg_pg_stat_monitor_pgsm_enable_query_plan"`
	EngineConfigPGPGStatMonitorPGSMMaxBuckets          types.Int64   `tfsdk:"engine_config_pg_pg_stat_monitor_pgsm_max_buckets"`
	EngineConfigPGPGStatStatementsTrack                types.String  `tfsdk:"engine_config_pg_pg_stat_statements_track"`
	EngineConfigPGTempFileLimit                        types.Int32   `tfsdk:"engine_config_pg_temp_file_limit"`
	EngineConfigPGTimezone                             types.String  `tfsdk:"engine_config_pg_timezone"`
	EngineConfigPGTrackActivityQuerySize               types.Int64   `tfsdk:"engine_config_pg_track_activity_query_size"`
	EngineConfigPGTrackCommitTimestamp                 types.String  `tfsdk:"engine_config_pg_track_commit_timestamp"`
	EngineConfigPGTrackFunctions                       types.String  `tfsdk:"engine_config_pg_track_functions"`
	EngineConfigPGTrackIOTiming                        types.String  `tfsdk:"engine_config_pg_track_io_timing"`
	EngineConfigPGWALSenderTimeout                     types.Int64   `tfsdk:"engine_config_pg_wal_sender_timeout"`
	EngineConfigPGWALWriterDelay                       types.Int64   `tfsdk:"engine_config_pg_wal_writer_delay"`
	EngineConfigPGStatMonitorEnable                    types.Bool    `tfsdk:"engine_config_pg_stat_monitor_enable"`
	EngineConfigPGLookoutMaxFailoverReplicationTimeLag types.Int64   `tfsdk:"engine_config_pglookout_max_failover_replication_time_lag"`
	EngineConfigSharedBuffersPercentage                types.Float64 `tfsdk:"engine_config_shared_buffers_percentage"`
	EngineConfigWorkMem                                types.Int64   `tfsdk:"engine_config_work_mem"`
}

func (m *Model) Refresh(
	ctx context.Context,
	client *linodego.Client,
	dbID int,
	preserveKnown bool,
) (d diag.Diagnostics) {
	tflog.SetField(ctx, "id", dbID)

	tflog.Debug(ctx, "Refreshing the PostgreSQL database...")

	tflog.Debug(ctx, "client.GetPostgresDatabase(...)")
	db, err := client.GetPostgresDatabase(ctx, dbID)
	if err != nil {
		d.AddError("Failed to refresh PostgreSQL database", err.Error())
		return d
	}

	var ssl *linodego.PostgresDatabaseSSL
	var creds *linodego.PostgresDatabaseCredential

	if !helper.DatabaseStatusIsSuspended(db.Status) {
		// SSL and credentials endpoints return 400s while a DB is suspended

		tflog.Debug(ctx, "client.GetPostgresDatabaseSSL(...)")
		ssl, err = client.GetPostgresDatabaseSSL(ctx, dbID)
		if err != nil {
			d.AddError("Failed to refresh PostgreSQL database SSL", err.Error())
			return d
		}

		tflog.Debug(ctx, "client.GetPostgresDatabaseCredentials(...)")
		creds, err = client.GetPostgresDatabaseCredentials(ctx, dbID)
		if err != nil {
			d.AddError("Failed to refresh PostgreSQL database credentials", err.Error())
			return d
		}
	}

	m.Flatten(ctx, db, ssl, creds, preserveKnown)
	return d
}

func (m *Model) Flatten(
	ctx context.Context,
	db *linodego.PostgresDatabase,
	ssl *linodego.PostgresDatabaseSSL,
	creds *linodego.PostgresDatabaseCredential,
	preserveKnown bool,
) (d diag.Diagnostics) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(db.ID), preserveKnown)

	m.ClusterSize = helper.KeepOrUpdateInt64(m.ClusterSize, int64(db.ClusterSize), preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, timetypes.NewRFC3339TimePointerValue(db.Created), preserveKnown)
	m.Encrypted = helper.KeepOrUpdateBool(m.Encrypted, db.Encrypted, preserveKnown)
	m.Engine = helper.KeepOrUpdateString(m.Engine, db.Engine, preserveKnown)
	m.EngineID = helper.KeepOrUpdateString(
		m.EngineID,
		helper.CreateDatabaseEngineSlug(db.Engine, db.Version),
		preserveKnown,
	)
	m.HostPrimary = helper.KeepOrUpdateString(m.HostPrimary, db.Hosts.Primary, preserveKnown)
	m.HostSecondary = helper.KeepOrUpdateString(m.HostSecondary, db.Hosts.Secondary, preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, db.Label, preserveKnown)
	m.OldestRestoreTime = helper.KeepOrUpdateValue(m.OldestRestoreTime, timetypes.NewRFC3339TimePointerValue(db.OldestRestoreTime), preserveKnown)
	m.Platform = helper.KeepOrUpdateString(m.Platform, string(db.Platform), preserveKnown)
	m.Port = helper.KeepOrUpdateInt64(m.Port, int64(db.Port), preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, db.Region, preserveKnown)
	m.SSLConnection = helper.KeepOrUpdateBool(m.SSLConnection, db.SSLConnection, preserveKnown)
	m.Status = helper.KeepOrUpdateString(m.Status, string(db.Status), preserveKnown)
	m.Suspended = helper.KeepOrUpdateBool(m.Suspended, helper.DatabaseStatusIsSuspended(db.Status), preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, db.Type, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, timetypes.NewRFC3339TimePointerValue(db.Updated), preserveKnown)
	m.Version = helper.KeepOrUpdateString(m.Version, db.Version, preserveKnown)

	// SSL and credentials may be nil if the database is suspended
	if ssl != nil {
		m.CACert = helper.KeepOrUpdateString(m.CACert, string(ssl.CACertificate), preserveKnown)
	} else {
		// We always enable perserveKnown here because it will otherwise
		// result in an inconsistent state when a database is suspended.
		// The alternative would be to make these fields not UseStateForUnknown(),
		// but that would prevent users from using these fields to initialize
		// a database provider during the same apply as a DB update under
		// certain circumstances.
		m.CACert = helper.KeepOrUpdateValue(m.CACert, types.StringNull(), true)
	}

	if creds != nil {
		m.RootPassword = helper.KeepOrUpdateString(m.RootPassword, creds.Password, preserveKnown)
		m.RootUsername = helper.KeepOrUpdateString(m.RootUsername, creds.Username, preserveKnown)
	} else {
		m.RootPassword = helper.KeepOrUpdateValue(m.RootPassword, types.StringNull(), true)
		m.RootUsername = helper.KeepOrUpdateValue(m.RootUsername, types.StringNull(), true)
	}

	m.AllowList = helper.KeepOrUpdateSet(
		types.StringType,
		m.AllowList,
		helper.StringSliceToFrameworkValueSlice(db.AllowList),
		preserveKnown,
		&d,
	)
	if d.HasError() {
		return d
	}

	membersCasted := helper.MapMap(
		db.Members,
		func(key string, value linodego.DatabaseMemberType) (string, string) {
			return key, string(value)
		},
	)

	m.Members = helper.KeepOrUpdateStringMap(ctx, m.Members, membersCasted, preserveKnown, &d)
	if d.HasError() {
		return d
	}

	if db.Fork != nil {
		m.ForkSource = helper.KeepOrUpdateInt64(
			m.ForkSource,
			int64(db.Fork.Source),
			preserveKnown,
		)

		m.ForkRestoreTime = helper.KeepOrUpdateValue(
			m.ForkRestoreTime,
			timetypes.NewRFC3339TimePointerValue(db.Fork.RestoreTime),
			preserveKnown,
		)

	} else {
		m.ForkSource = helper.KeepOrUpdateValue(
			m.ForkSource,
			types.Int64Null(),
			preserveKnown,
		)

		m.ForkRestoreTime = helper.KeepOrUpdateValue(
			m.ForkRestoreTime,
			timetypes.NewRFC3339Null(),
			preserveKnown,
		)
	}

	if db.PrivateNetwork != nil {
		privateNetworkObject, rd := types.ObjectValueFrom(
			ctx,
			privateNetworkAttributes,
			&ModelPrivateNetwork{
				VPCID:        types.Int64Value(int64(db.PrivateNetwork.VPCID)),
				SubnetID:     types.Int64Value(int64(db.PrivateNetwork.SubnetID)),
				PublicAccess: types.BoolValue(db.PrivateNetwork.PublicAccess),
			},
		)
		d.Append(rd...)
		m.PrivateNetwork = helper.KeepOrUpdateValue(m.PrivateNetwork, privateNetworkObject, preserveKnown)
	}

	updatesObject, rd := types.ObjectValueFrom(
		ctx,
		updatesAttributes,
		&ModelUpdates{
			DayOfWeek: types.Int64Value(int64(db.Updates.DayOfWeek)),
			Duration:  types.Int64Value(int64(db.Updates.Duration)),
			Frequency: types.StringValue(string(db.Updates.Frequency)),
			HourOfDay: types.Int64Value(int64(db.Updates.HourOfDay)),
		},
	)
	d.Append(rd...)
	m.Updates = helper.KeepOrUpdateValue(m.Updates, updatesObject, preserveKnown)

	m.EngineConfigPGAutovacuumAnalyzeScaleFactor = helper.KeepOrUpdateFloat64Pointer(
		m.EngineConfigPGAutovacuumAnalyzeScaleFactor,
		db.EngineConfig.PG.AutovacuumAnalyzeScaleFactor,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumAnalyzeThreshold = helper.KeepOrUpdateInt32Pointer(
		m.EngineConfigPGAutovacuumAnalyzeThreshold,
		db.EngineConfig.PG.AutovacuumAnalyzeThreshold,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumMaxWorkers = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGAutovacuumMaxWorkers,
		db.EngineConfig.PG.AutovacuumMaxWorkers,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumNaptime = helper.KeepOrUpdateIntPointer(m.EngineConfigPGAutovacuumNaptime, db.EngineConfig.PG.AutovacuumNaptime, preserveKnown)
	m.EngineConfigPGAutovacuumVacuumCostDelay = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGAutovacuumVacuumCostDelay,
		db.EngineConfig.PG.AutovacuumVacuumCostDelay,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumVacuumCostLimit = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGAutovacuumVacuumCostLimit,
		db.EngineConfig.PG.AutovacuumVacuumCostLimit,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumVacuumScaleFactor = helper.KeepOrUpdateFloat64Pointer(
		m.EngineConfigPGAutovacuumVacuumScaleFactor,
		db.EngineConfig.PG.AutovacuumVacuumScaleFactor,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumVacuumThreshold = helper.KeepOrUpdateInt32Pointer(
		m.EngineConfigPGAutovacuumVacuumThreshold,
		db.EngineConfig.PG.AutovacuumVacuumThreshold,
		preserveKnown,
	)
	m.EngineConfigPGBGWriterDelay = helper.KeepOrUpdateIntPointer(m.EngineConfigPGBGWriterDelay, db.EngineConfig.PG.BGWriterDelay, preserveKnown)
	m.EngineConfigPGBGWriterFlushAfter = helper.KeepOrUpdateIntPointer(m.EngineConfigPGBGWriterFlushAfter, db.EngineConfig.PG.BGWriterFlushAfter, preserveKnown)
	m.EngineConfigPGBGWriterLRUMaxpages = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGBGWriterLRUMaxpages,
		db.EngineConfig.PG.BGWriterLRUMaxPages,
		preserveKnown,
	)
	m.EngineConfigPGBGWriterLRUMultiplier = helper.KeepOrUpdateFloat64Pointer(
		m.EngineConfigPGBGWriterLRUMultiplier,
		db.EngineConfig.PG.BGWriterLRUMultiplier,
		preserveKnown,
	)
	m.EngineConfigPGDeadlockTimeout = helper.KeepOrUpdateIntPointer(m.EngineConfigPGDeadlockTimeout, db.EngineConfig.PG.DeadlockTimeout, preserveKnown)
	m.EngineConfigPGDefaultToastCompression = helper.KeepOrUpdateStringPointer(
		m.EngineConfigPGDefaultToastCompression,
		db.EngineConfig.PG.DefaultToastCompression,
		preserveKnown,
	)
	m.EngineConfigPGIdleInTransactionSessionTimeout = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGIdleInTransactionSessionTimeout,
		db.EngineConfig.PG.IdleInTransactionSessionTimeout,
		preserveKnown,
	)
	m.EngineConfigPGJIT = helper.KeepOrUpdateBoolPointer(m.EngineConfigPGJIT, db.EngineConfig.PG.JIT, preserveKnown)
	m.EngineConfigPGMaxFilesPerProcess = helper.KeepOrUpdateIntPointer(m.EngineConfigPGMaxFilesPerProcess, db.EngineConfig.PG.MaxFilesPerProcess, preserveKnown)
	m.EngineConfigPGMaxLocksPerTransaction = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGMaxLocksPerTransaction,
		db.EngineConfig.PG.MaxLocksPerTransaction,
		preserveKnown,
	)
	m.EngineConfigPGMaxLogicalReplicationWorkers = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGMaxLogicalReplicationWorkers,
		db.EngineConfig.PG.MaxLogicalReplicationWorkers,
		preserveKnown,
	)
	m.EngineConfigPGMaxParallelWorkers = helper.KeepOrUpdateIntPointer(m.EngineConfigPGMaxParallelWorkers, db.EngineConfig.PG.MaxParallelWorkers, preserveKnown)
	m.EngineConfigPGMaxParallelWorkersPerGather = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGMaxParallelWorkersPerGather,
		db.EngineConfig.PG.MaxParallelWorkersPerGather,
		preserveKnown,
	)
	m.EngineConfigPGMaxPredLocksPerTransaction = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGMaxPredLocksPerTransaction,
		db.EngineConfig.PG.MaxPredLocksPerTransaction,
		preserveKnown,
	)
	m.EngineConfigPGMaxReplicationSlots = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGMaxReplicationSlots,
		db.EngineConfig.PG.MaxReplicationSlots,
		preserveKnown,
	)
	m.EngineConfigPGMaxSlotWALKeepSize = helper.KeepOrUpdateInt32Pointer(
		m.EngineConfigPGMaxSlotWALKeepSize,
		db.EngineConfig.PG.MaxSlotWALKeepSize,
		preserveKnown,
	)
	m.EngineConfigPGMaxStackDepth = helper.KeepOrUpdateIntPointer(m.EngineConfigPGMaxStackDepth, db.EngineConfig.PG.MaxStackDepth, preserveKnown)
	m.EngineConfigPGMaxStandbyArchiveDelay = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGMaxStandbyArchiveDelay,
		db.EngineConfig.PG.MaxStandbyArchiveDelay,
		preserveKnown,
	)
	m.EngineConfigPGMaxStandbyStreamingDelay = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGMaxStandbyStreamingDelay,
		db.EngineConfig.PG.MaxStandbyStreamingDelay,
		preserveKnown,
	)
	m.EngineConfigPGMaxWALSenders = helper.KeepOrUpdateIntPointer(m.EngineConfigPGMaxWALSenders, db.EngineConfig.PG.MaxWALSenders, preserveKnown)
	m.EngineConfigPGMaxWorkerProcesses = helper.KeepOrUpdateIntPointer(m.EngineConfigPGMaxWorkerProcesses, db.EngineConfig.PG.MaxWorkerProcesses, preserveKnown)
	m.EngineConfigPGPasswordEncryption = helper.KeepOrUpdateStringPointer(
		m.EngineConfigPGPasswordEncryption,
		db.EngineConfig.PG.PasswordEncryption,
		preserveKnown,
	)
	m.EngineConfigPGPGPartmanBGWInterval = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGPGPartmanBGWInterval,
		db.EngineConfig.PG.PGPartmanBGWInterval,
		preserveKnown,
	)
	m.EngineConfigPGPGPartmanBGWRole = helper.KeepOrUpdateStringPointer(m.EngineConfigPGPGPartmanBGWRole, db.EngineConfig.PG.PGPartmanBGWRole, preserveKnown)
	m.EngineConfigPGPGStatMonitorPGSMEnableQueryPlan = helper.KeepOrUpdateBoolPointer(
		m.EngineConfigPGPGStatMonitorPGSMEnableQueryPlan,
		db.EngineConfig.PG.PGStatMonitorPGSMEnableQueryPlan,
		preserveKnown,
	)
	m.EngineConfigPGPGStatMonitorPGSMMaxBuckets = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGPGStatMonitorPGSMMaxBuckets,
		db.EngineConfig.PG.PGStatMonitorPGSMMaxBuckets,
		preserveKnown,
	)
	m.EngineConfigPGPGStatStatementsTrack = helper.KeepOrUpdateStringPointer(
		m.EngineConfigPGPGStatStatementsTrack,
		db.EngineConfig.PG.PGStatStatementsTrack,
		preserveKnown,
	)
	m.EngineConfigPGTempFileLimit = helper.KeepOrUpdateInt32Pointer(m.EngineConfigPGTempFileLimit, db.EngineConfig.PG.TempFileLimit, preserveKnown)
	m.EngineConfigPGTimezone = helper.KeepOrUpdateStringPointer(m.EngineConfigPGTimezone, db.EngineConfig.PG.Timezone, preserveKnown)
	m.EngineConfigPGTrackActivityQuerySize = helper.KeepOrUpdateIntPointer(
		m.EngineConfigPGTrackActivityQuerySize,
		db.EngineConfig.PG.TrackActivityQuerySize,
		preserveKnown,
	)
	m.EngineConfigPGTrackCommitTimestamp = helper.KeepOrUpdateStringPointer(
		m.EngineConfigPGTrackCommitTimestamp,
		db.EngineConfig.PG.TrackCommitTimestamp,
		preserveKnown,
	)
	m.EngineConfigPGTrackFunctions = helper.KeepOrUpdateStringPointer(m.EngineConfigPGTrackFunctions, db.EngineConfig.PG.TrackFunctions, preserveKnown)
	m.EngineConfigPGTrackIOTiming = helper.KeepOrUpdateStringPointer(m.EngineConfigPGTrackIOTiming, db.EngineConfig.PG.TrackIOTiming, preserveKnown)
	m.EngineConfigPGWALSenderTimeout = helper.KeepOrUpdateIntPointer(m.EngineConfigPGWALSenderTimeout, db.EngineConfig.PG.WALSenderTimeout, preserveKnown)
	m.EngineConfigPGWALWriterDelay = helper.KeepOrUpdateIntPointer(m.EngineConfigPGWALWriterDelay, db.EngineConfig.PG.WALWriterDelay, preserveKnown)
	m.EngineConfigPGStatMonitorEnable = helper.KeepOrUpdateBoolPointer(m.EngineConfigPGStatMonitorEnable, db.EngineConfig.PGStatMonitorEnable, preserveKnown)
	m.EngineConfigPGLookoutMaxFailoverReplicationTimeLag = helper.KeepOrUpdateInt64Pointer(
		m.EngineConfigPGLookoutMaxFailoverReplicationTimeLag,
		db.EngineConfig.PGLookout.MaxFailoverReplicationTimeLag,
		preserveKnown,
	)
	m.EngineConfigSharedBuffersPercentage = helper.KeepOrUpdateFloat64Pointer(
		m.EngineConfigSharedBuffersPercentage,
		db.EngineConfig.SharedBuffersPercentage,
		preserveKnown,
	)
	m.EngineConfigWorkMem = helper.KeepOrUpdateIntPointer(m.EngineConfigWorkMem, db.EngineConfig.WorkMem, preserveKnown)

	pendingObjects := helper.MapSlice(
		db.Updates.Pending,
		func(pending linodego.DatabaseMaintenanceWindowPending) ModelPendingUpdate {
			return ModelPendingUpdate{
				Deadline:    timetypes.NewRFC3339TimePointerValue(pending.Deadline),
				Description: types.StringValue(pending.Description),
				PlannedFor:  timetypes.NewRFC3339TimePointerValue(pending.PlannedFor),
			}
		},
	)

	pendingSet, rd := types.SetValueFrom(
		ctx,
		types.ObjectType{
			AttrTypes: pendingUpdateAttributes,
		},
		pendingObjects,
	)
	d.Append(rd...)
	if d.HasError() {
		return d
	}

	m.PendingUpdates = helper.KeepOrUpdateValue(m.PendingUpdates, pendingSet, preserveKnown)

	return d
}

func (m *Model) CopyFrom(other *Model, preserveKnown bool) {
	m.ForkSource = helper.KeepOrUpdateValue(m.ForkSource, other.ForkSource, preserveKnown)
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)

	m.AllowList = helper.KeepOrUpdateValue(m.AllowList, other.AllowList, preserveKnown)
	m.CACert = helper.KeepOrUpdateValue(m.CACert, other.CACert, preserveKnown)
	m.ClusterSize = helper.KeepOrUpdateValue(m.ClusterSize, other.ClusterSize, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Encrypted = helper.KeepOrUpdateValue(m.Encrypted, other.Encrypted, preserveKnown)
	m.Engine = helper.KeepOrUpdateValue(m.Engine, other.Engine, preserveKnown)
	m.EngineID = helper.KeepOrUpdateValue(m.EngineID, other.EngineID, preserveKnown)
	m.ForkRestoreTime = helper.KeepOrUpdateValue(m.ForkRestoreTime, other.ForkRestoreTime, preserveKnown)
	m.HostPrimary = helper.KeepOrUpdateValue(m.HostPrimary, other.HostPrimary, preserveKnown)
	m.HostSecondary = helper.KeepOrUpdateValue(m.HostSecondary, other.HostSecondary, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Members = helper.KeepOrUpdateValue(m.Members, other.Members, preserveKnown)
	m.OldestRestoreTime = helper.KeepOrUpdateValue(m.OldestRestoreTime, other.OldestRestoreTime, preserveKnown)
	m.PendingUpdates = helper.KeepOrUpdateValue(m.PendingUpdates, other.PendingUpdates, preserveKnown)
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

	m.EngineConfigPGAutovacuumAnalyzeScaleFactor = helper.KeepOrUpdateValue(
		m.EngineConfigPGAutovacuumAnalyzeScaleFactor,
		other.EngineConfigPGAutovacuumAnalyzeScaleFactor,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumAnalyzeThreshold = helper.KeepOrUpdateValue(
		m.EngineConfigPGAutovacuumAnalyzeThreshold,
		other.EngineConfigPGAutovacuumAnalyzeThreshold,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumMaxWorkers = helper.KeepOrUpdateValue(
		m.EngineConfigPGAutovacuumMaxWorkers,
		other.EngineConfigPGAutovacuumMaxWorkers,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumNaptime = helper.KeepOrUpdateValue(m.EngineConfigPGAutovacuumNaptime, other.EngineConfigPGAutovacuumNaptime, preserveKnown)
	m.EngineConfigPGAutovacuumVacuumCostDelay = helper.KeepOrUpdateValue(
		m.EngineConfigPGAutovacuumVacuumCostDelay,
		other.EngineConfigPGAutovacuumVacuumCostDelay,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumVacuumCostLimit = helper.KeepOrUpdateValue(
		m.EngineConfigPGAutovacuumVacuumCostLimit,
		other.EngineConfigPGAutovacuumVacuumCostLimit,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumVacuumScaleFactor = helper.KeepOrUpdateValue(
		m.EngineConfigPGAutovacuumVacuumScaleFactor,
		other.EngineConfigPGAutovacuumVacuumScaleFactor,
		preserveKnown,
	)
	m.EngineConfigPGAutovacuumVacuumThreshold = helper.KeepOrUpdateValue(
		m.EngineConfigPGAutovacuumVacuumThreshold,
		other.EngineConfigPGAutovacuumVacuumThreshold,
		preserveKnown,
	)
	m.EngineConfigPGBGWriterDelay = helper.KeepOrUpdateValue(m.EngineConfigPGBGWriterDelay, other.EngineConfigPGBGWriterDelay, preserveKnown)
	m.EngineConfigPGBGWriterFlushAfter = helper.KeepOrUpdateValue(m.EngineConfigPGBGWriterFlushAfter, other.EngineConfigPGBGWriterFlushAfter, preserveKnown)
	m.EngineConfigPGBGWriterLRUMaxpages = helper.KeepOrUpdateValue(m.EngineConfigPGBGWriterLRUMaxpages, other.EngineConfigPGBGWriterLRUMaxpages, preserveKnown)
	m.EngineConfigPGBGWriterLRUMultiplier = helper.KeepOrUpdateValue(
		m.EngineConfigPGBGWriterLRUMultiplier,
		other.EngineConfigPGBGWriterLRUMultiplier,
		preserveKnown,
	)
	m.EngineConfigPGDeadlockTimeout = helper.KeepOrUpdateValue(m.EngineConfigPGDeadlockTimeout, other.EngineConfigPGDeadlockTimeout, preserveKnown)
	m.EngineConfigPGDefaultToastCompression = helper.KeepOrUpdateValue(
		m.EngineConfigPGDefaultToastCompression,
		other.EngineConfigPGDefaultToastCompression,
		preserveKnown,
	)
	m.EngineConfigPGIdleInTransactionSessionTimeout = helper.KeepOrUpdateValue(
		m.EngineConfigPGIdleInTransactionSessionTimeout,
		other.EngineConfigPGIdleInTransactionSessionTimeout,
		preserveKnown,
	)
	m.EngineConfigPGJIT = helper.KeepOrUpdateValue(m.EngineConfigPGJIT, other.EngineConfigPGJIT, preserveKnown)
	m.EngineConfigPGMaxFilesPerProcess = helper.KeepOrUpdateValue(m.EngineConfigPGMaxFilesPerProcess, other.EngineConfigPGMaxFilesPerProcess, preserveKnown)
	m.EngineConfigPGMaxLocksPerTransaction = helper.KeepOrUpdateValue(
		m.EngineConfigPGMaxLocksPerTransaction,
		other.EngineConfigPGMaxLocksPerTransaction,
		preserveKnown,
	)
	m.EngineConfigPGMaxLogicalReplicationWorkers = helper.KeepOrUpdateValue(
		m.EngineConfigPGMaxLogicalReplicationWorkers,
		other.EngineConfigPGMaxLogicalReplicationWorkers,
		preserveKnown,
	)
	m.EngineConfigPGMaxParallelWorkers = helper.KeepOrUpdateValue(m.EngineConfigPGMaxParallelWorkers, other.EngineConfigPGMaxParallelWorkers, preserveKnown)
	m.EngineConfigPGMaxParallelWorkersPerGather = helper.KeepOrUpdateValue(
		m.EngineConfigPGMaxParallelWorkersPerGather,
		other.EngineConfigPGMaxParallelWorkersPerGather,
		preserveKnown,
	)
	m.EngineConfigPGMaxPredLocksPerTransaction = helper.KeepOrUpdateValue(
		m.EngineConfigPGMaxPredLocksPerTransaction,
		other.EngineConfigPGMaxPredLocksPerTransaction,
		preserveKnown,
	)
	m.EngineConfigPGMaxReplicationSlots = helper.KeepOrUpdateValue(m.EngineConfigPGMaxReplicationSlots, other.EngineConfigPGMaxReplicationSlots, preserveKnown)
	m.EngineConfigPGMaxSlotWALKeepSize = helper.KeepOrUpdateValue(m.EngineConfigPGMaxSlotWALKeepSize, other.EngineConfigPGMaxSlotWALKeepSize, preserveKnown)
	m.EngineConfigPGMaxStackDepth = helper.KeepOrUpdateValue(m.EngineConfigPGMaxStackDepth, other.EngineConfigPGMaxStackDepth, preserveKnown)
	m.EngineConfigPGMaxStandbyArchiveDelay = helper.KeepOrUpdateValue(
		m.EngineConfigPGMaxStandbyArchiveDelay,
		other.EngineConfigPGMaxStandbyArchiveDelay,
		preserveKnown,
	)
	m.EngineConfigPGMaxStandbyStreamingDelay = helper.KeepOrUpdateValue(
		m.EngineConfigPGMaxStandbyStreamingDelay,
		other.EngineConfigPGMaxStandbyStreamingDelay,
		preserveKnown,
	)
	m.EngineConfigPGMaxWALSenders = helper.KeepOrUpdateValue(m.EngineConfigPGMaxWALSenders, other.EngineConfigPGMaxWALSenders, preserveKnown)
	m.EngineConfigPGMaxWorkerProcesses = helper.KeepOrUpdateValue(m.EngineConfigPGMaxWorkerProcesses, other.EngineConfigPGMaxWorkerProcesses, preserveKnown)
	m.EngineConfigPGPasswordEncryption = helper.KeepOrUpdateValue(m.EngineConfigPGPasswordEncryption, other.EngineConfigPGPasswordEncryption, preserveKnown)
	m.EngineConfigPGPGPartmanBGWInterval = helper.KeepOrUpdateValue(
		m.EngineConfigPGPGPartmanBGWInterval,
		other.EngineConfigPGPGPartmanBGWInterval,
		preserveKnown,
	)
	m.EngineConfigPGPGPartmanBGWRole = helper.KeepOrUpdateValue(m.EngineConfigPGPGPartmanBGWRole, other.EngineConfigPGPGPartmanBGWRole, preserveKnown)
	m.EngineConfigPGPGStatMonitorPGSMEnableQueryPlan = helper.KeepOrUpdateValue(
		m.EngineConfigPGPGStatMonitorPGSMEnableQueryPlan,
		other.EngineConfigPGPGStatMonitorPGSMEnableQueryPlan,
		preserveKnown,
	)
	m.EngineConfigPGPGStatMonitorPGSMMaxBuckets = helper.KeepOrUpdateValue(
		m.EngineConfigPGPGStatMonitorPGSMMaxBuckets,
		other.EngineConfigPGPGStatMonitorPGSMMaxBuckets,
		preserveKnown,
	)
	m.EngineConfigPGPGStatStatementsTrack = helper.KeepOrUpdateValue(
		m.EngineConfigPGPGStatStatementsTrack,
		other.EngineConfigPGPGStatStatementsTrack,
		preserveKnown,
	)
	m.EngineConfigPGTempFileLimit = helper.KeepOrUpdateValue(m.EngineConfigPGTempFileLimit, other.EngineConfigPGTempFileLimit, preserveKnown)
	m.EngineConfigPGTimezone = helper.KeepOrUpdateValue(m.EngineConfigPGTimezone, other.EngineConfigPGTimezone, preserveKnown)
	m.EngineConfigPGTrackActivityQuerySize = helper.KeepOrUpdateValue(
		m.EngineConfigPGTrackActivityQuerySize,
		other.EngineConfigPGTrackActivityQuerySize,
		preserveKnown,
	)
	m.EngineConfigPGTrackCommitTimestamp = helper.KeepOrUpdateValue(
		m.EngineConfigPGTrackCommitTimestamp,
		other.EngineConfigPGTrackCommitTimestamp,
		preserveKnown,
	)
	m.EngineConfigPGTrackFunctions = helper.KeepOrUpdateValue(m.EngineConfigPGTrackFunctions, other.EngineConfigPGTrackFunctions, preserveKnown)
	m.EngineConfigPGTrackIOTiming = helper.KeepOrUpdateValue(m.EngineConfigPGTrackIOTiming, other.EngineConfigPGTrackIOTiming, preserveKnown)
	m.EngineConfigPGWALSenderTimeout = helper.KeepOrUpdateValue(m.EngineConfigPGWALSenderTimeout, other.EngineConfigPGWALSenderTimeout, preserveKnown)
	m.EngineConfigPGWALWriterDelay = helper.KeepOrUpdateValue(m.EngineConfigPGWALWriterDelay, other.EngineConfigPGWALWriterDelay, preserveKnown)
	m.EngineConfigPGStatMonitorEnable = helper.KeepOrUpdateValue(m.EngineConfigPGStatMonitorEnable, other.EngineConfigPGStatMonitorEnable, preserveKnown)
	m.EngineConfigPGLookoutMaxFailoverReplicationTimeLag = helper.KeepOrUpdateValue(
		m.EngineConfigPGLookoutMaxFailoverReplicationTimeLag,
		other.EngineConfigPGLookoutMaxFailoverReplicationTimeLag,
		preserveKnown,
	)
	m.EngineConfigSharedBuffersPercentage = helper.KeepOrUpdateValue(
		m.EngineConfigSharedBuffersPercentage,
		other.EngineConfigSharedBuffersPercentage,
		preserveKnown,
	)
	m.EngineConfigWorkMem = helper.KeepOrUpdateValue(m.EngineConfigWorkMem, other.EngineConfigWorkMem, preserveKnown)
}

// GetFork returns the linodego.DatabaseFork for this model if specified, else nil.
func (m *Model) GetFork(d diag.Diagnostics) *linodego.DatabaseFork {
	var result linodego.DatabaseFork

	isSpecified := false

	if !m.ForkSource.IsUnknown() && !m.ForkSource.IsNull() {
		isSpecified = true

		result.Source = helper.FrameworkSafeInt64ToInt(m.ForkSource.ValueInt64(), &d)
	}

	if !m.ForkRestoreTime.IsUnknown() && !m.ForkRestoreTime.IsNull() {
		isSpecified = true

		restoreTime, rd := m.ForkRestoreTime.ValueRFC3339Time()
		d.Append(rd...)

		result.RestoreTime = &restoreTime
	}

	if d.HasError() || !isSpecified {
		return nil
	}

	return &result
}

// GetAllowList returns the allow list slice for this model if specified, else nil.
func (m *Model) GetAllowList(ctx context.Context, d diag.Diagnostics) []string {
	if m.AllowList.IsUnknown() || m.AllowList.IsNull() {
		return nil
	}

	var result []string

	d.Append(
		m.AllowList.ElementsAs(
			ctx,
			&result,
			false,
		)...,
	)

	return result
}

// GetPrivateNetwork returns the ModelPrivateNetwork for this model if specified, else nil.
func (m *Model) GetPrivateNetwork(ctx context.Context, d diag.Diagnostics) *ModelPrivateNetwork {
	if m.PrivateNetwork.IsUnknown() || m.PrivateNetwork.IsNull() {
		return nil
	}

	var result ModelPrivateNetwork

	d.Append(
		m.PrivateNetwork.As(
			ctx,
			&result,
			basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true},
		)...,
	)

	return &result
}

// GetUpdates returns the ModelUpdates for this model if specified, else nil.
func (m *Model) GetUpdates(ctx context.Context, d diag.Diagnostics) *ModelUpdates {
	if m.Updates.IsUnknown() || m.Updates.IsNull() {
		return nil
	}

	var result ModelUpdates

	d.Append(
		m.Updates.As(
			ctx,
			&result,
			basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true},
		)...,
	)

	return &result
}

// GetEngineConfig returns a pointer to the linodego.PostgreSQLDatabaseEngineConfig for this model if specified, else nil.
func (m *Model) GetEngineConfig(d diag.Diagnostics) *linodego.PostgresDatabaseEngineConfig {
	var engineConfig linodego.PostgresDatabaseEngineConfig
	var engineConfigPG linodego.PostgresDatabaseEngineConfigPG
	var engineConfigPGLookout linodego.PostgresDatabaseEngineConfigPGLookout

	if !m.EngineConfigPGAutovacuumAnalyzeScaleFactor.IsUnknown() {
		engineConfigPG.AutovacuumAnalyzeScaleFactor = m.EngineConfigPGAutovacuumAnalyzeScaleFactor.ValueFloat64Pointer()
	}

	if !m.EngineConfigPGAutovacuumAnalyzeThreshold.IsUnknown() {
		engineConfigPG.AutovacuumAnalyzeThreshold = m.EngineConfigPGAutovacuumAnalyzeThreshold.ValueInt32Pointer()
	}

	if !m.EngineConfigPGAutovacuumMaxWorkers.IsUnknown() {
		engineConfigPG.AutovacuumMaxWorkers = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGAutovacuumMaxWorkers.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGAutovacuumNaptime.IsUnknown() {
		engineConfigPG.AutovacuumNaptime = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGAutovacuumNaptime.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGAutovacuumVacuumCostDelay.IsUnknown() {
		engineConfigPG.AutovacuumVacuumCostDelay = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigPGAutovacuumVacuumCostDelay.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigPGAutovacuumVacuumCostLimit.IsUnknown() {
		engineConfigPG.AutovacuumVacuumCostLimit = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigPGAutovacuumVacuumCostLimit.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigPGAutovacuumVacuumScaleFactor.IsUnknown() {
		engineConfigPG.AutovacuumVacuumScaleFactor = m.EngineConfigPGAutovacuumVacuumScaleFactor.ValueFloat64Pointer()
	}

	if !m.EngineConfigPGAutovacuumVacuumThreshold.IsUnknown() {
		engineConfigPG.AutovacuumVacuumThreshold = m.EngineConfigPGAutovacuumVacuumThreshold.ValueInt32Pointer()
	}

	if !m.EngineConfigPGBGWriterDelay.IsUnknown() {
		engineConfigPG.BGWriterDelay = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGBGWriterDelay.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGBGWriterFlushAfter.IsUnknown() {
		engineConfigPG.BGWriterFlushAfter = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGBGWriterFlushAfter.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGBGWriterLRUMaxpages.IsUnknown() {
		engineConfigPG.BGWriterLRUMaxPages = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGBGWriterLRUMaxpages.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGBGWriterLRUMultiplier.IsUnknown() {
		engineConfigPG.BGWriterLRUMultiplier = m.EngineConfigPGBGWriterLRUMultiplier.ValueFloat64Pointer()
	}

	if !m.EngineConfigPGDeadlockTimeout.IsUnknown() {
		engineConfigPG.DeadlockTimeout = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGDeadlockTimeout.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGDefaultToastCompression.IsUnknown() {
		engineConfigPG.DefaultToastCompression = m.EngineConfigPGDefaultToastCompression.ValueStringPointer()
	}

	if !m.EngineConfigPGIdleInTransactionSessionTimeout.IsUnknown() {
		engineConfigPG.IdleInTransactionSessionTimeout = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigPGIdleInTransactionSessionTimeout.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigPGJIT.IsUnknown() {
		engineConfigPG.JIT = m.EngineConfigPGJIT.ValueBoolPointer()
	}

	if !m.EngineConfigPGMaxFilesPerProcess.IsUnknown() {
		engineConfigPG.MaxFilesPerProcess = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGMaxFilesPerProcess.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGMaxLocksPerTransaction.IsUnknown() {
		engineConfigPG.MaxLocksPerTransaction = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGMaxLocksPerTransaction.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGMaxLogicalReplicationWorkers.IsUnknown() {
		engineConfigPG.MaxLogicalReplicationWorkers = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigPGMaxLogicalReplicationWorkers.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigPGMaxParallelWorkers.IsUnknown() {
		engineConfigPG.MaxParallelWorkers = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGMaxParallelWorkers.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGMaxParallelWorkersPerGather.IsUnknown() {
		engineConfigPG.MaxParallelWorkersPerGather = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigPGMaxParallelWorkersPerGather.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigPGMaxPredLocksPerTransaction.IsUnknown() {
		engineConfigPG.MaxPredLocksPerTransaction = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigPGMaxPredLocksPerTransaction.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigPGMaxReplicationSlots.IsUnknown() {
		engineConfigPG.MaxReplicationSlots = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGMaxReplicationSlots.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGMaxSlotWALKeepSize.IsUnknown() {
		engineConfigPG.MaxSlotWALKeepSize = m.EngineConfigPGMaxSlotWALKeepSize.ValueInt32Pointer()
	}

	if !m.EngineConfigPGMaxStackDepth.IsUnknown() {
		engineConfigPG.MaxStackDepth = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGMaxStackDepth.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGMaxStandbyArchiveDelay.IsUnknown() {
		engineConfigPG.MaxStandbyArchiveDelay = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGMaxStandbyArchiveDelay.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGMaxStandbyStreamingDelay.IsUnknown() {
		engineConfigPG.MaxStandbyStreamingDelay = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGMaxStandbyStreamingDelay.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGMaxWALSenders.IsUnknown() {
		engineConfigPG.MaxWALSenders = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGMaxWALSenders.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGMaxWorkerProcesses.IsUnknown() {
		engineConfigPG.MaxWorkerProcesses = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGMaxWorkerProcesses.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGPasswordEncryption.IsUnknown() {
		engineConfigPG.PasswordEncryption = m.EngineConfigPGPasswordEncryption.ValueStringPointer()
	}

	if !m.EngineConfigPGPGPartmanBGWInterval.IsUnknown() {
		engineConfigPG.PGPartmanBGWInterval = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGPGPartmanBGWInterval.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGPGPartmanBGWRole.IsUnknown() {
		engineConfigPG.PGPartmanBGWRole = m.EngineConfigPGPGPartmanBGWRole.ValueStringPointer()
	}

	if !m.EngineConfigPGPGStatMonitorPGSMEnableQueryPlan.IsUnknown() {
		engineConfigPG.PGStatMonitorPGSMEnableQueryPlan = m.EngineConfigPGPGStatMonitorPGSMEnableQueryPlan.ValueBoolPointer()
	}

	if !m.EngineConfigPGPGStatMonitorPGSMMaxBuckets.IsUnknown() {
		engineConfigPG.PGStatMonitorPGSMMaxBuckets = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigPGPGStatMonitorPGSMMaxBuckets.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigPGPGStatStatementsTrack.IsUnknown() {
		engineConfigPG.PGStatStatementsTrack = m.EngineConfigPGPGStatStatementsTrack.ValueStringPointer()
	}

	if !m.EngineConfigPGTempFileLimit.IsUnknown() {
		engineConfigPG.TempFileLimit = m.EngineConfigPGTempFileLimit.ValueInt32Pointer()
	}

	if !m.EngineConfigPGTimezone.IsUnknown() {
		engineConfigPG.Timezone = m.EngineConfigPGTimezone.ValueStringPointer()
	}

	if !m.EngineConfigPGTrackActivityQuerySize.IsUnknown() {
		engineConfigPG.TrackActivityQuerySize = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGTrackActivityQuerySize.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGTrackCommitTimestamp.IsUnknown() {
		engineConfigPG.TrackCommitTimestamp = m.EngineConfigPGTrackCommitTimestamp.ValueStringPointer()
	}

	if !m.EngineConfigPGTrackFunctions.IsUnknown() {
		engineConfigPG.TrackFunctions = m.EngineConfigPGTrackFunctions.ValueStringPointer()
	}
	if !m.EngineConfigPGTrackIOTiming.IsUnknown() {
		engineConfigPG.TrackIOTiming = m.EngineConfigPGTrackIOTiming.ValueStringPointer()
	}

	if !m.EngineConfigPGWALSenderTimeout.IsUnknown() {
		engineConfigPG.WALSenderTimeout = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGWALSenderTimeout.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGWALWriterDelay.IsUnknown() {
		engineConfigPG.WALWriterDelay = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigPGWALWriterDelay.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigPGStatMonitorEnable.IsUnknown() {
		engineConfig.PGStatMonitorEnable = m.EngineConfigPGStatMonitorEnable.ValueBoolPointer()
	}

	if !m.EngineConfigPGLookoutMaxFailoverReplicationTimeLag.IsUnknown() {
		engineConfigPGLookout.MaxFailoverReplicationTimeLag = m.EngineConfigPGLookoutMaxFailoverReplicationTimeLag.ValueInt64Pointer()
	}

	if !m.EngineConfigSharedBuffersPercentage.IsUnknown() {
		engineConfig.SharedBuffersPercentage = m.EngineConfigSharedBuffersPercentage.ValueFloat64Pointer()
	}

	if !m.EngineConfigWorkMem.IsUnknown() {
		engineConfig.WorkMem = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigWorkMem.ValueInt64Pointer(), &d)
	}

	engineConfig.PG = &engineConfigPG
	engineConfig.PGLookout = &engineConfigPGLookout
	return &engineConfig
}
