package databasemysqlconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ConnectTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var DefaultTimeZoneObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.StringType,
		"maxLength":        types.Int64Type,
		"minLength":        types.Int64Type,
		"pattern":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var GroupConcatMaxLenObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Float64Type,
		"maximum":          types.Float64Type,
		"minimum":          types.Float64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InformationSchemaStatsExpiryObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBChangeBufferMaxSizeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBFlushNeighborsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBFTMinTokenSizeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBFTServerStopwordTableObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.StringType,
		"maxLength":        types.Int64Type,
		"pattern":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.ListType{ElemType: types.StringType},
	},
}

var InnoDBLockWaitTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBLogBufferSizeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBOnlineAlterLogMaxSizeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBReadIOThreadsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBRollbackOnTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.BoolType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBThreadConcurrencyObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InnoDBWriteIOThreadsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InteractiveTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var InternalTmpMemStorageEngineObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"enum":             types.ListType{ElemType: types.StringType},
		"example":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxAllowedPacketObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MaxHeapTableSizeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var NetBufferLengthObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var NetReadTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var NetWriteTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var SortBufferSizeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var SQLModeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.StringType,
		"maxLength":        types.Int64Type,
		"pattern":          types.StringType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var SQLRequirePrimaryKeyObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.BoolType,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var TmpTableSizeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var WaitTimeoutObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var BinlogRetentionPeriodObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"description":      types.StringType,
		"example":          types.Int64Type,
		"maximum":          types.Int64Type,
		"minimum":          types.Int64Type,
		"requires_restart": types.BoolType,
		"type":             types.StringType,
	},
}

var MySQLConfigMySQLObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"connect_timeout":                  ConnectTimeoutObjectType,
		"default_time_zone":                DefaultTimeZoneObjectType,
		"group_concat_max_len":             GroupConcatMaxLenObjectType,
		"information_schema_stats_expiry":  InformationSchemaStatsExpiryObjectType,
		"innodb_change_buffer_max_size":    InnoDBChangeBufferMaxSizeObjectType,
		"innodb_flush_neighbors":           InnoDBFlushNeighborsObjectType,
		"innodb_ft_min_token_size":         InnoDBFTMinTokenSizeObjectType,
		"innodb_ft_server_stopword_table":  InnoDBFTServerStopwordTableObjectType,
		"innodb_lock_wait_timeout":         InnoDBLockWaitTimeoutObjectType,
		"innodb_log_buffer_size":           InnoDBLogBufferSizeObjectType,
		"innodb_online_alter_log_max_size": InnoDBOnlineAlterLogMaxSizeObjectType,
		"innodb_read_io_threads":           InnoDBReadIOThreadsObjectType,
		"innodb_rollback_on_timeout":       InnoDBRollbackOnTimeoutObjectType,
		"innodb_thread_concurrency":        InnoDBThreadConcurrencyObjectType,
		"innodb_write_io_threads":          InnoDBWriteIOThreadsObjectType,
		"interactive_timeout":              InteractiveTimeoutObjectType,
		"internal_tmp_mem_storage_engine":  InternalTmpMemStorageEngineObjectType,
		"max_allowed_packet":               MaxAllowedPacketObjectType,
		"max_heap_table_size":              MaxHeapTableSizeObjectType,
		"net_buffer_length":                NetBufferLengthObjectType,
		"net_read_timeout":                 NetReadTimeoutObjectType,
		"net_write_timeout":                NetWriteTimeoutObjectType,
		"sort_buffer_size":                 SortBufferSizeObjectType,
		"sql_mode":                         SQLModeObjectType,
		"sql_require_primary_key":          SQLRequirePrimaryKeyObjectType,
		"tmp_table_size":                   TmpTableSizeObjectType,
		"wait_timeout":                     WaitTimeoutObjectType,
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Unique identifier for this DataSource.",
			Computed:    true,
		},
		"mysql": schema.ListAttribute{
			Description: "MySQL configuration settings.",
			Computed:    true,
			ElementType: MySQLConfigMySQLObjectType,
		},
		"binlog_retention_period": schema.ListAttribute{
			Description: "The minimum amount of time in seconds to keep binlog entries before deletion." +
				"This may be extended for services that require binlog entries for longer than the default" +
				"for example if using the MySQL Debezium Kafka connector.",
			Computed:    true,
			ElementType: BinlogRetentionPeriodObjectType,
		},
	},
}
