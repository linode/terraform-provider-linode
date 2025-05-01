package databasemysqlv2

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

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

type ModelEngineConfig struct {
	BinlogRetentionPeriod types.Int64             `tfsdk:"binlog_retention_period"`
	MySQL                 *ModelEngineConfigMySQL `tfsdk:"mysql"`
}

type ModelEngineConfigMySQL struct {
	ConnectTimeout               types.Int64   `tfsdk:"connect_timeout"`
	DefaultTimeZone              types.String  `tfsdk:"default_time_zone"`
	GroupConcatMaxLen            types.Float64 `tfsdk:"group_concat_max_len"`
	InformationSchemaStatsExpiry types.Int64   `tfsdk:"information_schema_stats_expiry"`
	InnoDBChangeBufferMaxSize    types.Int64   `tfsdk:"innodb_change_buffer_max_size"`
	InnoDBFlushNeighbors         types.Int64   `tfsdk:"innodb_flush_neighbors"`
	InnoDBFTMinTokenSize         types.Int64   `tfsdk:"innodb_ft_min_token_size"`
	InnoDBFTServerStopwordTable  types.String  `tfsdk:"innodb_ft_server_stopword_table"`
	InnoDBLockWaitTimeout        types.Int64   `tfsdk:"innodb_lock_wait_timeout"`
	InnoDBLogBufferSize          types.Int64   `tfsdk:"innodb_log_buffer_size"`
	InnoDBOnlineAlterLogMaxSize  types.Int64   `tfsdk:"innodb_online_alter_log_max_size"`
	InnoDBReadIOThreads          types.Int64   `tfsdk:"innodb_read_io_threads"`
	InnoDBRollbackOnTimeout      types.Bool    `tfsdk:"innodb_rollback_on_timeout"`
	InnoDBThreadConcurrency      types.Int64   `tfsdk:"innodb_thread_concurrency"`
	InnoDBWriteIOThreads         types.Int64   `tfsdk:"innodb_write_io_threads"`
	InteractiveTimeout           types.Int64   `tfsdk:"interactive_timeout"`
	InternalTmpMemStorageEngine  types.String  `tfsdk:"internal_tmp_mem_storage_engine"`
	MaxAllowedPacket             types.Int64   `tfsdk:"max_allowed_packet"`
	MaxHeapTableSize             types.Int64   `tfsdk:"max_heap_table_size"`
	NetBufferLength              types.Int64   `tfsdk:"net_buffer_length"`
	NetReadTimeout               types.Int64   `tfsdk:"net_read_timeout"`
	NetWriteTimeout              types.Int64   `tfsdk:"net_write_timeout"`
	SortBufferSize               types.Int64   `tfsdk:"sort_buffer_size"`
	SQLMode                      types.String  `tfsdk:"sql_mode"`
	SQLRequirePrimaryKey         types.Bool    `tfsdk:"sql_require_primary_key"`
	TmpTableSize                 types.Int64   `tfsdk:"tmp_table_size"`
	WaitTimeout                  types.Int64   `tfsdk:"wait_timeout"`
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

	Updates        types.Object `tfsdk:"updates"`
	PendingUpdates types.Set    `tfsdk:"pending_updates"`

	EngineConfig *ModelEngineConfig `tfsdk:"engine_config"`
}

func (m *Model) Refresh(
	ctx context.Context,
	client *linodego.Client,
	dbID int,
	preserveKnown bool,
) (d diag.Diagnostics) {
	tflog.SetField(ctx, "id", dbID)

	tflog.Debug(ctx, "Refreshing the MySQL database...")

	tflog.Debug(ctx, "client.GetMySQLDatabase(...)")
	db, err := client.GetMySQLDatabase(ctx, dbID)
	if err != nil {
		d.AddError("Failed to refresh MySQL database", err.Error())
		return
	}

	var ssl *linodego.MySQLDatabaseSSL
	var creds *linodego.MySQLDatabaseCredential

	if !helper.DatabaseStatusIsSuspended(db.Status) {
		// SSL and credentials endpoints return 400s while a DB is suspended

		tflog.Debug(ctx, "client.GetMySQLDatabaseSSL(...)")
		ssl, err = client.GetMySQLDatabaseSSL(ctx, dbID)
		if err != nil {
			d.AddError("Failed to refresh MySQL database SSL", err.Error())
			return
		}

		tflog.Debug(ctx, "client.GetMySQLDatabaseCredentials(...)")
		creds, err = client.GetMySQLDatabaseCredentials(ctx, dbID)
		if err != nil {
			d.AddError("Failed to refresh MySQL database credentials", err.Error())
			return
		}
	}

	m.Flatten(ctx, db, ssl, creds, preserveKnown)
	return
}

func (m *Model) Flatten(
	ctx context.Context,
	db *linodego.MySQLDatabase,
	ssl *linodego.MySQLDatabaseSSL,
	creds *linodego.MySQLDatabaseCredential,
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
		return
	}

	membersCasted := helper.MapMap(
		db.Members,
		func(key string, value linodego.DatabaseMemberType) (string, string) {
			return key, string(value)
		},
	)

	m.Members = helper.KeepOrUpdateStringMap(ctx, m.Members, membersCasted, preserveKnown, &d)
	if d.HasError() {
		return
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

	engineConfigModel := ModelEngineConfig{
		BinlogRetentionPeriod: int64OrNull(db.EngineConfig.BinlogRetentionPeriod),
		MySQL: &ModelEngineConfigMySQL{
			ConnectTimeout:               int64OrNull(db.EngineConfig.MySQL.ConnectTimeout),
			DefaultTimeZone:              stringOrNull(db.EngineConfig.MySQL.DefaultTimeZone),
			GroupConcatMaxLen:            float64OrNull(db.EngineConfig.MySQL.GroupConcatMaxLen),
			InformationSchemaStatsExpiry: int64OrNull(db.EngineConfig.MySQL.InformationSchemaStatsExpiry),
			InnoDBChangeBufferMaxSize:    int64OrNull(db.EngineConfig.MySQL.InnoDBChangeBufferMaxSize),
			InnoDBFlushNeighbors:         int64OrNull(db.EngineConfig.MySQL.InnoDBFlushNeighbors),
			InnoDBFTMinTokenSize:         int64OrNull(db.EngineConfig.MySQL.InnoDBFTMinTokenSize),
			InnoDBFTServerStopwordTable:  stringOrNull(db.EngineConfig.MySQL.InnoDBFTServerStopwordTable),
			InnoDBLockWaitTimeout:        int64OrNull(db.EngineConfig.MySQL.InnoDBLockWaitTimeout),
			InnoDBLogBufferSize:          int64OrNull(db.EngineConfig.MySQL.InnoDBLogBufferSize),
			InnoDBOnlineAlterLogMaxSize:  int64OrNull(db.EngineConfig.MySQL.InnoDBOnlineAlterLogMaxSize),
			InnoDBReadIOThreads:          int64OrNull(db.EngineConfig.MySQL.InnoDBReadIOThreads),
			InnoDBRollbackOnTimeout:      boolOrNull(db.EngineConfig.MySQL.InnoDBRollbackOnTimeout),
			InnoDBThreadConcurrency:      int64OrNull(db.EngineConfig.MySQL.InnoDBThreadConcurrency),
			InnoDBWriteIOThreads:         int64OrNull(db.EngineConfig.MySQL.InnoDBWriteIOThreads),
			InteractiveTimeout:           int64OrNull(db.EngineConfig.MySQL.InteractiveTimeout),
			InternalTmpMemStorageEngine:  stringOrNull(db.EngineConfig.MySQL.InternalTmpMemStorageEngine),
			MaxAllowedPacket:             int64OrNull(db.EngineConfig.MySQL.MaxAllowedPacket),
			MaxHeapTableSize:             int64OrNull(db.EngineConfig.MySQL.MaxHeapTableSize),
			NetBufferLength:              int64OrNull(db.EngineConfig.MySQL.NetBufferLength),
			NetReadTimeout:               int64OrNull(db.EngineConfig.MySQL.NetReadTimeout),
			NetWriteTimeout:              int64OrNull(db.EngineConfig.MySQL.NetWriteTimeout),
			SortBufferSize:               int64OrNull(db.EngineConfig.MySQL.SortBufferSize),
			SQLMode:                      stringOrNull(db.EngineConfig.MySQL.SQLMode),
			SQLRequirePrimaryKey:         boolOrNull(db.EngineConfig.MySQL.SQLRequirePrimaryKey),
			TmpTableSize:                 int64OrNull(db.EngineConfig.MySQL.TmpTableSize),
			WaitTimeout:                  int64OrNull(db.EngineConfig.MySQL.WaitTimeout),
		},
	}

	m.EngineConfig = &engineConfigModel

	pendingObjects := helper.MapSlice(
		db.Updates.Pending,
		func(pending linodego.DatabaseMaintenanceWindowPending) types.Object {
			result, rd := types.ObjectValueFrom(
				ctx,
				pendingUpdateAttributes,
				&ModelPendingUpdate{
					Deadline:    timetypes.NewRFC3339TimePointerValue(pending.Deadline),
					Description: types.StringValue(pending.Description),
					PlannedFor:  timetypes.NewRFC3339TimePointerValue(pending.PlannedFor),
				},
			)
			d.Append(rd...)

			return result
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

	m.PendingUpdates = helper.KeepOrUpdateValue(m.PendingUpdates, pendingSet, preserveKnown)

	return nil
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

	if !preserveKnown {
		m.EngineConfig = other.EngineConfig
	}
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

// GetAttributeCapable abstracts over tfsdk.Plan and tfsdk.State, allowing access to attribute values via the GetAttribute method
type GetAttributeCapable interface {
	GetAttribute(ctx context.Context, path path.Path, target interface{}) diag.Diagnostics
}

// GetEngineConfig returns the ModelEngineConfig for this model if specified, else nil.
func (m *ModelEngineConfig) GetEngineConfig(d diag.Diagnostics, data GetAttributeCapable) *linodego.MySQLDatabaseEngineConfig {
	var engineConfig linodego.MySQLDatabaseEngineConfig

	if m == nil || ((m.BinlogRetentionPeriod.IsNull() || m.BinlogRetentionPeriod.IsUnknown()) && m.MySQL == nil) {
		return nil
	}

	binlogRetentionPeriod := helper.FrameworkSafeInt64PointerToIntPointer(m.BinlogRetentionPeriod.ValueInt64Pointer(), &d)
	engineConfig.BinlogRetentionPeriod = binlogRetentionPeriod

	var engineConfigMySQL linodego.MySQLDatabaseEngineConfigMySQL

	connectTimeout := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.ConnectTimeout.ValueInt64Pointer(), &d)
	engineConfigMySQL.ConnectTimeout = connectTimeout

	defaultTimeZone := m.MySQL.DefaultTimeZone.ValueStringPointer()
	engineConfigMySQL.DefaultTimeZone = defaultTimeZone

	groupConcatMaxLen := m.MySQL.GroupConcatMaxLen.ValueFloat64Pointer()
	engineConfigMySQL.GroupConcatMaxLen = groupConcatMaxLen

	informationSchemaStatsExpiry := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InformationSchemaStatsExpiry.ValueInt64Pointer(), &d)
	engineConfigMySQL.InformationSchemaStatsExpiry = informationSchemaStatsExpiry

	innodbChangeBufferMaxSize := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InnoDBChangeBufferMaxSize.ValueInt64Pointer(), &d)
	engineConfigMySQL.InnoDBChangeBufferMaxSize = innodbChangeBufferMaxSize

	innodbFlushNeighbors := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InnoDBFlushNeighbors.ValueInt64Pointer(), &d)
	engineConfigMySQL.InnoDBFlushNeighbors = innodbFlushNeighbors

	innodbFTMinTokenSize := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InnoDBFTMinTokenSize.ValueInt64Pointer(), &d)
	engineConfigMySQL.InnoDBFTMinTokenSize = innodbFTMinTokenSize

	innodbFTServerStopwordTable := m.MySQL.InnoDBFTServerStopwordTable.ValueStringPointer()
	engineConfigMySQL.InnoDBFTServerStopwordTable = innodbFTServerStopwordTable

	innodbLockWaitTimeout := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InnoDBLockWaitTimeout.ValueInt64Pointer(), &d)
	engineConfigMySQL.InnoDBLockWaitTimeout = innodbLockWaitTimeout

	innodbLogBufferSize := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InnoDBLogBufferSize.ValueInt64Pointer(), &d)
	engineConfigMySQL.InnoDBLogBufferSize = innodbLogBufferSize

	innodbOnlineAlterLogMaxSize := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InnoDBOnlineAlterLogMaxSize.ValueInt64Pointer(), &d)
	engineConfigMySQL.InnoDBOnlineAlterLogMaxSize = innodbOnlineAlterLogMaxSize

	innodbReadIOThreads := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InnoDBReadIOThreads.ValueInt64Pointer(), &d)
	engineConfigMySQL.InnoDBReadIOThreads = innodbReadIOThreads

	innodbRollbackOnTimeout := m.MySQL.InnoDBRollbackOnTimeout.ValueBoolPointer()
	engineConfigMySQL.InnoDBRollbackOnTimeout = innodbRollbackOnTimeout

	innodbThreadConcurrency := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InnoDBThreadConcurrency.ValueInt64Pointer(), &d)
	engineConfigMySQL.InnoDBThreadConcurrency = innodbThreadConcurrency

	innodbWriteIOThreads := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InnoDBWriteIOThreads.ValueInt64Pointer(), &d)
	engineConfigMySQL.InnoDBWriteIOThreads = innodbWriteIOThreads

	interactiveTimeout := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.InteractiveTimeout.ValueInt64Pointer(), &d)
	engineConfigMySQL.InteractiveTimeout = interactiveTimeout

	internalTmpMemStorageEngine := m.MySQL.InternalTmpMemStorageEngine.ValueStringPointer()
	engineConfigMySQL.InternalTmpMemStorageEngine = internalTmpMemStorageEngine

	maxAllowedPacket := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.MaxAllowedPacket.ValueInt64Pointer(), &d)
	engineConfigMySQL.MaxAllowedPacket = maxAllowedPacket

	maxHeapTableSize := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.MaxHeapTableSize.ValueInt64Pointer(), &d)
	engineConfigMySQL.MaxHeapTableSize = maxHeapTableSize

	netBufferLength := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.NetBufferLength.ValueInt64Pointer(), &d)
	engineConfigMySQL.NetBufferLength = netBufferLength

	netReadTimeout := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.NetReadTimeout.ValueInt64Pointer(), &d)
	engineConfigMySQL.NetReadTimeout = netReadTimeout

	netWriteTimeout := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.NetWriteTimeout.ValueInt64Pointer(), &d)
	engineConfigMySQL.NetWriteTimeout = netWriteTimeout

	sortBufferSize := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.SortBufferSize.ValueInt64Pointer(), &d)
	engineConfigMySQL.SortBufferSize = sortBufferSize

	var sqlMode types.String
	diags := data.GetAttribute(context.Background(), path.Root("engine_config").AtName("mysql").AtName("sql_mode"), &sqlMode)
	d.Append(diags...)
	if !sqlMode.IsNull() && !sqlMode.IsUnknown() {
		engineConfigMySQL.SQLMode = sqlMode.ValueStringPointer()
	}

	var sqlRequirePrimaryKey types.Bool
	diags = data.GetAttribute(context.Background(), path.Root("engine_config").AtName("mysql").AtName("sql_require_primary_key"), &sqlRequirePrimaryKey)
	d.Append(diags...)
	if !sqlRequirePrimaryKey.IsNull() && !sqlRequirePrimaryKey.IsUnknown() {
		engineConfigMySQL.SQLRequirePrimaryKey = sqlRequirePrimaryKey.ValueBoolPointer()
	}

	tmpTableSize := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.TmpTableSize.ValueInt64Pointer(), &d)
	engineConfigMySQL.TmpTableSize = tmpTableSize

	waitTimeout := helper.FrameworkSafeInt64PointerToIntPointer(m.MySQL.WaitTimeout.ValueInt64Pointer(), &d)
	engineConfigMySQL.WaitTimeout = waitTimeout

	engineConfig.MySQL = &engineConfigMySQL

	return &engineConfig
}

func int64OrNull(v *int) types.Int64 {
	if v != nil {
		return types.Int64Value(int64(*v))
	}
	return types.Int64Null()
}

func float64OrNull(v *float64) types.Float64 {
	if v != nil {
		return types.Float64Value(*v)
	}
	return types.Float64Null()
}

func stringOrNull(v *string) types.String {
	if v != nil {
		return types.StringValue(*v)
	}
	return types.StringNull()
}

func boolOrNull(v *bool) types.Bool {
	if v != nil {
		return types.BoolValue(*v)
	}
	return types.BoolNull()
}
