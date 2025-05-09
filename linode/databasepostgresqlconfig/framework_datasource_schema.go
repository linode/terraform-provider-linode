package databasepostgresqlconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AutovacuumAnalyzeScaleFactorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Float64Type,
		"minimum":          types.Float64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var AutovacuumAnalyzeThresholdObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int32Type,
		"minimum":          types.Int32Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var AutovacuumMaxWorkersObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var AutovacuumNaptimeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var AutovacuumVacuumCostDelayObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var AutovacuumVacuumCostLimitObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var AutovacuumVacuumScaleFactorObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Float64Type,
		"minimum":          types.Float64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var AutovacuumVacuumThresholdObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int32Type,
		"minimum":          types.Int32Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var BGWriterDelayObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var BGWriterFlushAfterObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var BGWriterLRUMaxPagesObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var BGWriterLRUMultiplierObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Float64Type,
		"maximum":          types.Float64Type,
		"minimum":          types.Float64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var DeadlockTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var DefaultToastCompressionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"enum":             types.ListType{ElemType: types.StringType},
		"example":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var IdleInTransactionSessionTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var JITObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.BoolType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxFilesPerProcessObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxLocksPerTransactionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxLogicalReplicationWorkersObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxParallelWorkersObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxParallelWorkersPerGatherObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxPredLocksPerTransactionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxReplicationSlotsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxSlotWALKeepSizeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int32Type,
		"minimum":          types.Int32Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxStackDepthObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxStandbyArchiveDelayObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxStandbyStreamingDelayObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxWALSendersObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxWorkerProcessesObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var PasswordEncryptionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"enum":             types.ListType{ElemType: types.StringType},
		"example":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.ListType{ElemType: types.StringType},
	},
}

var PGPartmanBGWIntervalObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var PGPartmanBGWRoleObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.StringType,
		"maxLength":        types.Int64Type,
		"pattern":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var PGStatMonitorPGSMEnableQueryPlanObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.BoolType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var PGStatMonitorPGSMMaxBucketsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var PGStatStatementsTrackObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"enum":             types.ListType{ElemType: types.StringType},
		"requires_restart": types.BoolType,
		"type":             types.ListType{ElemType: types.StringType},
	},
}

var TempFileLimitObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int32Type,
		"maximum":          types.Int32Type,
		"minimum":          types.Int32Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var TimezoneObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.StringType,
		"maxLength":        types.Int64Type,
		"pattern":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var TrackActivityQuerySizeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var TrackCommitTimestampObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"enum":             types.ListType{ElemType: types.StringType},
		"example":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var TrackFunctionsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"enum":             types.ListType{ElemType: types.StringType},
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var TrackIOTimingObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"enum":             types.ListType{ElemType: types.StringType},
		"example":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var WALSenderTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var WALWriterDelayObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var PGStatMonitorEnableObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var PGLookoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"max_failover_replication_time_lag": PGLookoutMaxFailoverReplicationTimeLagObjectType,
	},
}

var PGLookoutMaxFailoverReplicationTimeLagObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var SharedBuffersPercentageObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Float64Type,
		"maximum":          types.Float64Type,
		"minimum":          types.Float64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var WorkMemObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var PostgreSQLConfigPGObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"autovacuum_analyze_scale_factor":        AutovacuumAnalyzeScaleFactorObjectType,
		"autovacuum_analyze_threshold":           AutovacuumAnalyzeThresholdObjectType,
		"autovacuum_max_workers":                 AutovacuumMaxWorkersObjectType,
		"autovacuum_naptime":                     AutovacuumNaptimeObjectType,
		"autovacuum_vacuum_cost_delay":           AutovacuumVacuumCostDelayObjectType,
		"autovacuum_vacuum_cost_limit":           AutovacuumVacuumCostLimitObjectType,
		"autovacuum_vacuum_scale_factor":         AutovacuumVacuumScaleFactorObjectType,
		"autovacuum_vacuum_threshold":            AutovacuumVacuumThresholdObjectType,
		"bgwriter_delay":                         BGWriterDelayObjectType,
		"bgwriter_flush_after":                   BGWriterFlushAfterObjectType,
		"bgwriter_lru_maxpages":                  BGWriterLRUMaxPagesObjectType,
		"bgwriter_lru_multiplier":                BGWriterLRUMultiplierObjectType,
		"deadlock_timeout":                       DeadlockTimeoutObjectType,
		"default_toast_compression":              DefaultToastCompressionObjectType,
		"idle_in_transaction_session_timeout":    IdleInTransactionSessionTimeoutObjectType,
		"jit":                                    JITObjectType,
		"max_files_per_process":                  MaxFilesPerProcessObjectType,
		"max_locks_per_transaction":              MaxLocksPerTransactionObjectType,
		"max_logical_replication_workers":        MaxLogicalReplicationWorkersObjectType,
		"max_parallel_workers":                   MaxParallelWorkersObjectType,
		"max_parallel_workers_per_gather":        MaxParallelWorkersPerGatherObjectType,
		"max_pred_locks_per_transaction":         MaxPredLocksPerTransactionObjectType,
		"max_replication_slots":                  MaxReplicationSlotsObjectType,
		"max_slot_wal_keep_size":                 MaxSlotWALKeepSizeObjectType,
		"max_stack_depth":                        MaxStackDepthObjectType,
		"max_standby_archive_delay":              MaxStandbyArchiveDelayObjectType,
		"max_standby_streaming_delay":            MaxStandbyStreamingDelayObjectType,
		"max_wal_senders":                        MaxWALSendersObjectType,
		"max_worker_processes":                   MaxWorkerProcessesObjectType,
		"password_encryption":                    PasswordEncryptionObjectType,
		"pg_partman_bgw.interval":                PGPartmanBGWIntervalObjectType,
		"pg_partman_bgw.role":                    PGPartmanBGWRoleObjectType,
		"pg_stat_monitor.pgsm_enable_query_plan": PGStatMonitorPGSMEnableQueryPlanObjectType,
		"pg_stat_monitor.pgsm_max_buckets":       PGStatMonitorPGSMMaxBucketsObjectType,
		"pg_stat_statements.track":               PGStatStatementsTrackObjectType,
		"temp_file_limit":                        TempFileLimitObjectType,
		"timezone":                               TimezoneObjectType,
		"track_activity_query_size":              TrackActivityQuerySizeObjectType,
		"track_commit_timestamp":                 TrackCommitTimestampObjectType,
		"track_functions":                        TrackFunctionsObjectType,
		"track_io_timing":                        TrackIOTimingObjectType,
		"wal_sender_timeout":                     WALSenderTimeoutObjectType,
		"wal_writer_delay":                       WALWriterDelayObjectType,
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Unique identifier for this DataSource.",
			Computed:    true,
		},
		"pg": schema.ListAttribute{
			Description: "PostgreSQL configuration settings.",
			Computed:    true,
			ElementType: PostgreSQLConfigPGObjectType,
		},
		"pg_stat_monitor_enable": schema.ListAttribute{
			Description: "Enable the pg_stat_monitor extension. Enabling this extension will cause the cluster to be restarted.When this extension is enabled, pg_stat_statements results for utility commands are unreliable.",
			Computed:    true,
			ElementType: PGStatMonitorEnableObjectType,
		},
		"pglookout": schema.ListAttribute{
			Description: "System-wide settings for pglookout.",
			Computed:    true,
			ElementType: PGLookoutObjectType,
		},
		"shared_buffers_percentage": schema.ListAttribute{
			Description: "Percentage of total RAM that the database server uses for shared memory buffers. Valid range is 20-60 (float), which corresponds to 20% - 60%. This setting adjusts the shared_buffers configuration value.",
			Computed:    true,
			ElementType: SharedBuffersPercentageObjectType,
		},
		"work_mem": schema.ListAttribute{
			Description: "Sets the maximum amount of memory to be used by a query operation (such as a sort or hash table) before writing to temporary disk files, in MB. Default is 1MB + 0.075% of total RAM (up to 32MB).",
			Computed:    true,
			ElementType: WorkMemObjectType,
		},
	},
}
