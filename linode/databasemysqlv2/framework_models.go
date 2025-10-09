package databasemysqlv2

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
	EngineConfigBinlogRetentionPeriod             types.Int64   `tfsdk:"engine_config_binlog_retention_period"`
	EngineConfigMySQLConnectTimeout               types.Int64   `tfsdk:"engine_config_mysql_connect_timeout"`
	EngineConfigMySQLDefaultTimeZone              types.String  `tfsdk:"engine_config_mysql_default_time_zone"`
	EngineConfigMySQLGroupConcatMaxLen            types.Float64 `tfsdk:"engine_config_mysql_group_concat_max_len"`
	EngineConfigMySQLInformationSchemaStatsExpiry types.Int64   `tfsdk:"engine_config_mysql_information_schema_stats_expiry"`
	EngineConfigMySQLInnoDBChangeBufferMaxSize    types.Int64   `tfsdk:"engine_config_mysql_innodb_change_buffer_max_size"`
	EngineConfigMySQLInnoDBFlushNeighbors         types.Int64   `tfsdk:"engine_config_mysql_innodb_flush_neighbors"`
	EngineConfigMySQLInnoDBFTMinTokenSize         types.Int64   `tfsdk:"engine_config_mysql_innodb_ft_min_token_size"`
	EngineConfigMySQLInnoDBFTServerStopwordTable  types.String  `tfsdk:"engine_config_mysql_innodb_ft_server_stopword_table"`
	EngineConfigMySQLInnoDBLockWaitTimeout        types.Int64   `tfsdk:"engine_config_mysql_innodb_lock_wait_timeout"`
	EngineConfigMySQLInnoDBLogBufferSize          types.Int64   `tfsdk:"engine_config_mysql_innodb_log_buffer_size"`
	EngineConfigMySQLInnoDBOnlineAlterLogMaxSize  types.Int64   `tfsdk:"engine_config_mysql_innodb_online_alter_log_max_size"`
	EngineConfigMySQLInnoDBReadIOThreads          types.Int64   `tfsdk:"engine_config_mysql_innodb_read_io_threads"`
	EngineConfigMySQLInnoDBRollbackOnTimeout      types.Bool    `tfsdk:"engine_config_mysql_innodb_rollback_on_timeout"`
	EngineConfigMySQLInnoDBThreadConcurrency      types.Int64   `tfsdk:"engine_config_mysql_innodb_thread_concurrency"`
	EngineConfigMySQLInnoDBWriteIOThreads         types.Int64   `tfsdk:"engine_config_mysql_innodb_write_io_threads"`
	EngineConfigMySQLInteractiveTimeout           types.Int64   `tfsdk:"engine_config_mysql_interactive_timeout"`
	EngineConfigMySQLInternalTmpMemStorageEngine  types.String  `tfsdk:"engine_config_mysql_internal_tmp_mem_storage_engine"`
	EngineConfigMySQLMaxAllowedPacket             types.Int64   `tfsdk:"engine_config_mysql_max_allowed_packet"`
	EngineConfigMySQLMaxHeapTableSize             types.Int64   `tfsdk:"engine_config_mysql_max_heap_table_size"`
	EngineConfigMySQLNetBufferLength              types.Int64   `tfsdk:"engine_config_mysql_net_buffer_length"`
	EngineConfigMySQLNetReadTimeout               types.Int64   `tfsdk:"engine_config_mysql_net_read_timeout"`
	EngineConfigMySQLNetWriteTimeout              types.Int64   `tfsdk:"engine_config_mysql_net_write_timeout"`
	EngineConfigMySQLSortBufferSize               types.Int64   `tfsdk:"engine_config_mysql_sort_buffer_size"`
	EngineConfigMySQLSQLMode                      types.String  `tfsdk:"engine_config_mysql_sql_mode"`
	EngineConfigMySQLSQLRequirePrimaryKey         types.Bool    `tfsdk:"engine_config_mysql_sql_require_primary_key"`
	EngineConfigMySQLTmpTableSize                 types.Int64   `tfsdk:"engine_config_mysql_tmp_table_size"`
	EngineConfigMySQLWaitTimeout                  types.Int64   `tfsdk:"engine_config_mysql_wait_timeout"`
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
		return d
	}

	var ssl *linodego.MySQLDatabaseSSL
	var creds *linodego.MySQLDatabaseCredential

	if !helper.DatabaseStatusIsSuspended(db.Status) {
		// SSL and credentials endpoints return 400s while a DB is suspended

		tflog.Debug(ctx, "client.GetMySQLDatabaseSSL(...)")
		ssl, err = client.GetMySQLDatabaseSSL(ctx, dbID)
		if err != nil {
			d.AddError("Failed to refresh MySQL database SSL", err.Error())
			return d
		}

		tflog.Debug(ctx, "client.GetMySQLDatabaseCredentials(...)")
		creds, err = client.GetMySQLDatabaseCredentials(ctx, dbID)
		if err != nil {
			d.AddError("Failed to refresh MySQL database credentials", err.Error())
			return d
		}
	}

	m.Flatten(ctx, db, ssl, creds, preserveKnown)
	return d
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

	m.EngineConfigBinlogRetentionPeriod = helper.KeepOrUpdateIntPointer(
		m.EngineConfigBinlogRetentionPeriod,
		db.EngineConfig.BinlogRetentionPeriod,
		preserveKnown,
	)
	m.EngineConfigMySQLConnectTimeout = helper.KeepOrUpdateIntPointer(m.EngineConfigMySQLConnectTimeout, db.EngineConfig.MySQL.ConnectTimeout, preserveKnown)
	m.EngineConfigMySQLDefaultTimeZone = helper.KeepOrUpdateStringPointer(
		m.EngineConfigMySQLDefaultTimeZone,
		db.EngineConfig.MySQL.DefaultTimeZone,
		preserveKnown,
	)
	m.EngineConfigMySQLGroupConcatMaxLen = helper.KeepOrUpdateFloat64Pointer(
		m.EngineConfigMySQLGroupConcatMaxLen,
		db.EngineConfig.MySQL.GroupConcatMaxLen,
		preserveKnown,
	)
	m.EngineConfigMySQLInformationSchemaStatsExpiry = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInformationSchemaStatsExpiry,
		db.EngineConfig.MySQL.InformationSchemaStatsExpiry,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBChangeBufferMaxSize = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInnoDBChangeBufferMaxSize,
		db.EngineConfig.MySQL.InnoDBChangeBufferMaxSize,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBFlushNeighbors = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInnoDBFlushNeighbors,
		db.EngineConfig.MySQL.InnoDBFlushNeighbors,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBFTMinTokenSize = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInnoDBFTMinTokenSize,
		db.EngineConfig.MySQL.InnoDBFTMinTokenSize,
		preserveKnown,
	)

	var stopwordTable *string
	if db.EngineConfig.MySQL != nil && db.EngineConfig.MySQL.InnoDBFTServerStopwordTable != nil {
		stopwordTable = *db.EngineConfig.MySQL.InnoDBFTServerStopwordTable
	}

	m.EngineConfigMySQLInnoDBFTServerStopwordTable = helper.KeepOrUpdateStringPointer(
		m.EngineConfigMySQLInnoDBFTServerStopwordTable,
		stopwordTable,
		preserveKnown,
	)

	m.EngineConfigMySQLInnoDBLockWaitTimeout = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInnoDBLockWaitTimeout,
		db.EngineConfig.MySQL.InnoDBLockWaitTimeout,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBLogBufferSize = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInnoDBLogBufferSize,
		db.EngineConfig.MySQL.InnoDBLogBufferSize,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBOnlineAlterLogMaxSize = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInnoDBOnlineAlterLogMaxSize,
		db.EngineConfig.MySQL.InnoDBOnlineAlterLogMaxSize,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBReadIOThreads = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInnoDBReadIOThreads,
		db.EngineConfig.MySQL.InnoDBReadIOThreads,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBRollbackOnTimeout = helper.KeepOrUpdateBoolPointer(
		m.EngineConfigMySQLInnoDBRollbackOnTimeout,
		db.EngineConfig.MySQL.InnoDBRollbackOnTimeout,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBThreadConcurrency = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInnoDBThreadConcurrency,
		db.EngineConfig.MySQL.InnoDBThreadConcurrency,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBWriteIOThreads = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInnoDBWriteIOThreads,
		db.EngineConfig.MySQL.InnoDBWriteIOThreads,
		preserveKnown,
	)
	m.EngineConfigMySQLInteractiveTimeout = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLInteractiveTimeout,
		db.EngineConfig.MySQL.InteractiveTimeout,
		preserveKnown,
	)
	m.EngineConfigMySQLInternalTmpMemStorageEngine = helper.KeepOrUpdateStringPointer(
		m.EngineConfigMySQLInternalTmpMemStorageEngine,
		db.EngineConfig.MySQL.InternalTmpMemStorageEngine,
		preserveKnown,
	)
	m.EngineConfigMySQLMaxAllowedPacket = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLMaxAllowedPacket,
		db.EngineConfig.MySQL.MaxAllowedPacket,
		preserveKnown,
	)
	m.EngineConfigMySQLMaxHeapTableSize = helper.KeepOrUpdateIntPointer(
		m.EngineConfigMySQLMaxHeapTableSize,
		db.EngineConfig.MySQL.MaxHeapTableSize,
		preserveKnown,
	)
	m.EngineConfigMySQLNetBufferLength = helper.KeepOrUpdateIntPointer(m.EngineConfigMySQLNetBufferLength, db.EngineConfig.MySQL.NetBufferLength, preserveKnown)
	m.EngineConfigMySQLNetReadTimeout = helper.KeepOrUpdateIntPointer(m.EngineConfigMySQLNetReadTimeout, db.EngineConfig.MySQL.NetReadTimeout, preserveKnown)
	m.EngineConfigMySQLNetWriteTimeout = helper.KeepOrUpdateIntPointer(m.EngineConfigMySQLNetWriteTimeout, db.EngineConfig.MySQL.NetWriteTimeout, preserveKnown)
	m.EngineConfigMySQLSortBufferSize = helper.KeepOrUpdateIntPointer(m.EngineConfigMySQLSortBufferSize, db.EngineConfig.MySQL.SortBufferSize, preserveKnown)
	m.EngineConfigMySQLSQLMode = helper.KeepOrUpdateStringPointer(m.EngineConfigMySQLSQLMode, db.EngineConfig.MySQL.SQLMode, preserveKnown)
	m.EngineConfigMySQLSQLRequirePrimaryKey = helper.KeepOrUpdateBoolPointer(
		m.EngineConfigMySQLSQLRequirePrimaryKey,
		db.EngineConfig.MySQL.SQLRequirePrimaryKey,
		preserveKnown,
	)
	m.EngineConfigMySQLTmpTableSize = helper.KeepOrUpdateIntPointer(m.EngineConfigMySQLTmpTableSize, db.EngineConfig.MySQL.TmpTableSize, preserveKnown)
	m.EngineConfigMySQLWaitTimeout = helper.KeepOrUpdateIntPointer(m.EngineConfigMySQLWaitTimeout, db.EngineConfig.MySQL.WaitTimeout, preserveKnown)

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

	m.EngineConfigBinlogRetentionPeriod = helper.KeepOrUpdateValue(m.EngineConfigBinlogRetentionPeriod, other.EngineConfigBinlogRetentionPeriod, preserveKnown)
	m.EngineConfigMySQLConnectTimeout = helper.KeepOrUpdateValue(m.EngineConfigMySQLConnectTimeout, other.EngineConfigMySQLConnectTimeout, preserveKnown)
	m.EngineConfigMySQLDefaultTimeZone = helper.KeepOrUpdateValue(m.EngineConfigMySQLDefaultTimeZone, other.EngineConfigMySQLDefaultTimeZone, preserveKnown)
	m.EngineConfigMySQLGroupConcatMaxLen = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLGroupConcatMaxLen,
		other.EngineConfigMySQLGroupConcatMaxLen,
		preserveKnown,
	)
	m.EngineConfigMySQLInformationSchemaStatsExpiry = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInformationSchemaStatsExpiry,
		other.EngineConfigMySQLInformationSchemaStatsExpiry,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBChangeBufferMaxSize = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBChangeBufferMaxSize,
		other.EngineConfigMySQLInnoDBChangeBufferMaxSize,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBFlushNeighbors = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBFlushNeighbors,
		other.EngineConfigMySQLInnoDBFlushNeighbors,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBFTMinTokenSize = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBFTMinTokenSize,
		other.EngineConfigMySQLInnoDBFTMinTokenSize,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBFTServerStopwordTable = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBFTServerStopwordTable,
		other.EngineConfigMySQLInnoDBFTServerStopwordTable,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBLockWaitTimeout = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBLockWaitTimeout,
		other.EngineConfigMySQLInnoDBLockWaitTimeout,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBLogBufferSize = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBLogBufferSize,
		other.EngineConfigMySQLInnoDBLogBufferSize,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBOnlineAlterLogMaxSize = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBOnlineAlterLogMaxSize,
		other.EngineConfigMySQLInnoDBOnlineAlterLogMaxSize,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBReadIOThreads = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBReadIOThreads,
		other.EngineConfigMySQLInnoDBReadIOThreads,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBRollbackOnTimeout = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBRollbackOnTimeout,
		other.EngineConfigMySQLInnoDBRollbackOnTimeout,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBThreadConcurrency = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBThreadConcurrency,
		other.EngineConfigMySQLInnoDBThreadConcurrency,
		preserveKnown,
	)
	m.EngineConfigMySQLInnoDBWriteIOThreads = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInnoDBWriteIOThreads,
		other.EngineConfigMySQLInnoDBWriteIOThreads,
		preserveKnown,
	)
	m.EngineConfigMySQLInteractiveTimeout = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInteractiveTimeout,
		other.EngineConfigMySQLInteractiveTimeout,
		preserveKnown,
	)
	m.EngineConfigMySQLInternalTmpMemStorageEngine = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLInternalTmpMemStorageEngine,
		other.EngineConfigMySQLInternalTmpMemStorageEngine,
		preserveKnown,
	)
	m.EngineConfigMySQLMaxAllowedPacket = helper.KeepOrUpdateValue(m.EngineConfigMySQLMaxAllowedPacket, other.EngineConfigMySQLMaxAllowedPacket, preserveKnown)
	m.EngineConfigMySQLMaxHeapTableSize = helper.KeepOrUpdateValue(m.EngineConfigMySQLMaxHeapTableSize, other.EngineConfigMySQLMaxHeapTableSize, preserveKnown)
	m.EngineConfigMySQLNetBufferLength = helper.KeepOrUpdateValue(m.EngineConfigMySQLNetBufferLength, other.EngineConfigMySQLNetBufferLength, preserveKnown)
	m.EngineConfigMySQLNetReadTimeout = helper.KeepOrUpdateValue(m.EngineConfigMySQLNetReadTimeout, other.EngineConfigMySQLNetReadTimeout, preserveKnown)
	m.EngineConfigMySQLNetWriteTimeout = helper.KeepOrUpdateValue(m.EngineConfigMySQLNetWriteTimeout, other.EngineConfigMySQLNetWriteTimeout, preserveKnown)
	m.EngineConfigMySQLSortBufferSize = helper.KeepOrUpdateValue(m.EngineConfigMySQLSortBufferSize, other.EngineConfigMySQLSortBufferSize, preserveKnown)
	m.EngineConfigMySQLSQLMode = helper.KeepOrUpdateValue(m.EngineConfigMySQLSQLMode, other.EngineConfigMySQLSQLMode, preserveKnown)
	m.EngineConfigMySQLSQLRequirePrimaryKey = helper.KeepOrUpdateValue(
		m.EngineConfigMySQLSQLRequirePrimaryKey,
		other.EngineConfigMySQLSQLRequirePrimaryKey,
		preserveKnown,
	)
	m.EngineConfigMySQLTmpTableSize = helper.KeepOrUpdateValue(m.EngineConfigMySQLTmpTableSize, other.EngineConfigMySQLTmpTableSize, preserveKnown)
	m.EngineConfigMySQLWaitTimeout = helper.KeepOrUpdateValue(m.EngineConfigMySQLWaitTimeout, other.EngineConfigMySQLWaitTimeout, preserveKnown)
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

// GetEngineConfig returns a pointer to the linodego.MySQLDatabaseEngineConfig for this model if specified, else nil.
func (m *Model) GetEngineConfig(d diag.Diagnostics) *linodego.MySQLDatabaseEngineConfig {
	var engineConfig linodego.MySQLDatabaseEngineConfig
	var engineConfigMySQL linodego.MySQLDatabaseEngineConfigMySQL

	if !m.EngineConfigBinlogRetentionPeriod.IsUnknown() {
		engineConfig.BinlogRetentionPeriod = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigBinlogRetentionPeriod.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLConnectTimeout.IsUnknown() {
		engineConfigMySQL.ConnectTimeout = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLConnectTimeout.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLDefaultTimeZone.IsUnknown() {
		engineConfigMySQL.DefaultTimeZone = m.EngineConfigMySQLDefaultTimeZone.ValueStringPointer()
	}

	if !m.EngineConfigMySQLGroupConcatMaxLen.IsUnknown() {
		engineConfigMySQL.GroupConcatMaxLen = m.EngineConfigMySQLGroupConcatMaxLen.ValueFloat64Pointer()
	}

	if !m.EngineConfigMySQLInformationSchemaStatsExpiry.IsUnknown() {
		engineConfigMySQL.InformationSchemaStatsExpiry = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigMySQLInformationSchemaStatsExpiry.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigMySQLInnoDBChangeBufferMaxSize.IsUnknown() {
		engineConfigMySQL.InnoDBChangeBufferMaxSize = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigMySQLInnoDBChangeBufferMaxSize.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigMySQLInnoDBFlushNeighbors.IsUnknown() {
		engineConfigMySQL.InnoDBFlushNeighbors = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLInnoDBFlushNeighbors.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLInnoDBFTMinTokenSize.IsUnknown() {
		engineConfigMySQL.InnoDBFTMinTokenSize = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLInnoDBFTMinTokenSize.ValueInt64Pointer(), &d)
	}

	if m.EngineConfigMySQLInnoDBFTServerStopwordTable.IsUnknown() {
		var nilString *string = nil
		engineConfigMySQL.InnoDBFTServerStopwordTable = linodego.Pointer(nilString)
	} else {
		engineConfigMySQL.InnoDBFTServerStopwordTable = linodego.Pointer(m.EngineConfigMySQLInnoDBFTServerStopwordTable.ValueStringPointer())
	}

	if !m.EngineConfigMySQLInnoDBLockWaitTimeout.IsUnknown() {
		engineConfigMySQL.InnoDBLockWaitTimeout = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLInnoDBLockWaitTimeout.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLInnoDBLogBufferSize.IsUnknown() {
		engineConfigMySQL.InnoDBLogBufferSize = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLInnoDBLogBufferSize.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLInnoDBOnlineAlterLogMaxSize.IsUnknown() {
		engineConfigMySQL.InnoDBOnlineAlterLogMaxSize = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigMySQLInnoDBOnlineAlterLogMaxSize.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigMySQLInnoDBReadIOThreads.IsUnknown() {
		engineConfigMySQL.InnoDBReadIOThreads = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLInnoDBReadIOThreads.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLInnoDBRollbackOnTimeout.IsUnknown() {
		engineConfigMySQL.InnoDBRollbackOnTimeout = m.EngineConfigMySQLInnoDBRollbackOnTimeout.ValueBoolPointer()
	}

	if !m.EngineConfigMySQLInnoDBThreadConcurrency.IsUnknown() {
		engineConfigMySQL.InnoDBThreadConcurrency = helper.FrameworkSafeInt64PointerToIntPointer(
			m.EngineConfigMySQLInnoDBThreadConcurrency.ValueInt64Pointer(),
			&d,
		)
	}

	if !m.EngineConfigMySQLInnoDBWriteIOThreads.IsUnknown() {
		engineConfigMySQL.InnoDBWriteIOThreads = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLInnoDBWriteIOThreads.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLInteractiveTimeout.IsUnknown() {
		engineConfigMySQL.InteractiveTimeout = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLInteractiveTimeout.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLInternalTmpMemStorageEngine.IsUnknown() {
		engineConfigMySQL.InternalTmpMemStorageEngine = m.EngineConfigMySQLInternalTmpMemStorageEngine.ValueStringPointer()
	}

	if !m.EngineConfigMySQLMaxAllowedPacket.IsUnknown() {
		engineConfigMySQL.MaxAllowedPacket = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLMaxAllowedPacket.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLMaxHeapTableSize.IsUnknown() {
		engineConfigMySQL.MaxHeapTableSize = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLMaxHeapTableSize.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLNetBufferLength.IsUnknown() {
		engineConfigMySQL.NetBufferLength = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLNetBufferLength.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLNetReadTimeout.IsUnknown() {
		engineConfigMySQL.NetReadTimeout = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLNetReadTimeout.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLNetWriteTimeout.IsUnknown() {
		engineConfigMySQL.NetWriteTimeout = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLNetWriteTimeout.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLSortBufferSize.IsUnknown() {
		engineConfigMySQL.SortBufferSize = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLSortBufferSize.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLSQLMode.IsUnknown() {
		engineConfigMySQL.SQLMode = m.EngineConfigMySQLSQLMode.ValueStringPointer()
	}

	if !m.EngineConfigMySQLSQLRequirePrimaryKey.IsUnknown() {
		engineConfigMySQL.SQLRequirePrimaryKey = m.EngineConfigMySQLSQLRequirePrimaryKey.ValueBoolPointer()
	}

	if !m.EngineConfigMySQLTmpTableSize.IsUnknown() {
		engineConfigMySQL.TmpTableSize = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLTmpTableSize.ValueInt64Pointer(), &d)
	}

	if !m.EngineConfigMySQLWaitTimeout.IsUnknown() {
		engineConfigMySQL.WaitTimeout = helper.FrameworkSafeInt64PointerToIntPointer(m.EngineConfigMySQLWaitTimeout.ValueInt64Pointer(), &d)
	}

	engineConfig.MySQL = &engineConfigMySQL
	return &engineConfig
}
