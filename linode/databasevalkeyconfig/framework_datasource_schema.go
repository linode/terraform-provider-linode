package databasevalkeyconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func configOption(description string, fields map[string]attr.Type) schema.ListAttribute {
	fields["description"] = types.StringType
	fields["requires_restart"] = types.BoolType
	fields["type"] = types.ListType{ElemType: types.StringType}

	return schema.ListAttribute{
		Description: description,
		Computed:    true,
		ElementType: types.ObjectType{AttrTypes: fields},
	}
}

func intConfigOption(description, valueField string) schema.ListAttribute {
	return configOption(description, map[string]attr.Type{
		valueField: types.Int64Type,
		"minimum":  types.Int64Type,
		"maximum":  types.Int64Type,
	})
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Unique identifier for this DataSource.",
			Computed:    true,
		},
		"backup_hour":   intConfigOption("The UTC hour when a backup starts.", "example"),
		"backup_minute": intConfigOption("The minute within the hour when a backup starts.", "example"),
		"frequent_snapshots": configOption("Whether frequent local RDB snapshots are enabled.", map[string]attr.Type{
			"default": types.BoolType,
		}),
		"valkey_acl_channels_default": configOption("The default pub/sub channel ACL for new users.", map[string]attr.Type{
			"enum": types.ListType{ElemType: types.StringType},
		}),
		"valkey_active_expire_effort": intConfigOption("The effort used to reclaim expired keys.", "default"),
		"valkey_activedefrag": configOption("Whether active memory defragmentation is enabled.", map[string]attr.Type{
			"default": types.BoolType,
		}),
		"valkey_lfu_decay_time": intConfigOption("The LFU counter decay time in minutes.", "default"),
		"valkey_lfu_log_factor": intConfigOption("The logarithm factor for LFU counters.", "default"),
		"valkey_maxmemory_policy": configOption("The memory eviction policy.", map[string]attr.Type{
			"default": types.StringType,
			"enum":    types.ListType{ElemType: types.StringType},
		}),
		"valkey_number_of_databases": intConfigOption("The number of logical Valkey databases.", "example"),
		"valkey_persistence": configOption("The persistence mode.", map[string]attr.Type{
			"enum": types.ListType{ElemType: types.StringType},
		}),
		"valkey_pubsub_client_output_buffer_limit": intConfigOption("The hard output-buffer limit for pub/sub clients in MB.", "example"),
		"valkey_timeout": intConfigOption("The idle connection timeout in seconds.", "default"),
	},
}
