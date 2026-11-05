package databasevalkeyconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func configOption(fields map[string]schema.Attribute) schema.ListNestedAttribute {
	fields["description"] = schema.StringAttribute{Computed: true, Description: "A description of the configuration option."}
	fields["requires_restart"] = schema.BoolAttribute{Computed: true, Description: "Whether changing this option requires a restart."}
	fields["type"] = schema.ListAttribute{Computed: true, ElementType: types.StringType, Description: "The accepted API types, including null when supported."}
	return schema.ListNestedAttribute{
		Computed:     true,
		Description:  "Metadata for this Valkey configuration option.",
		NestedObject: schema.NestedAttributeObject{Attributes: fields},
	}
}

func intConfigOption(valueField string) schema.ListNestedAttribute {
	return configOption(map[string]schema.Attribute{
		valueField: schema.Int64Attribute{Computed: true, Description: "The API-provided " + valueField + " value."},
		"minimum":  schema.Int64Attribute{Computed: true, Description: "The minimum accepted value."},
		"maximum":  schema.Int64Attribute{Computed: true, Description: "The maximum accepted value."},
	})
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id":            schema.StringAttribute{Computed: true, Description: "A content-derived identifier for the configuration catalog."},
		"backup_hour":   intConfigOption("example"),
		"backup_minute": intConfigOption("example"),
		"frequent_snapshots": configOption(map[string]schema.Attribute{
			"default": schema.BoolAttribute{Computed: true, Description: "The default value."},
		}),
		"valkey_acl_channels_default": configOption(map[string]schema.Attribute{
			"enum": schema.ListAttribute{Computed: true, ElementType: types.StringType, Description: "The accepted values."},
		}),
		"valkey_active_expire_effort": intConfigOption("default"),
		"valkey_activedefrag": configOption(map[string]schema.Attribute{
			"default": schema.BoolAttribute{Computed: true, Description: "The default value."},
		}),
		"valkey_lfu_decay_time": intConfigOption("default"),
		"valkey_lfu_log_factor": intConfigOption("default"),
		"valkey_maxmemory_policy": configOption(map[string]schema.Attribute{
			"default": schema.StringAttribute{Computed: true, Description: "The default value."},
			"enum":    schema.ListAttribute{Computed: true, ElementType: types.StringType, Description: "The accepted values."},
		}),
		"valkey_number_of_databases": intConfigOption("example"),
		"valkey_persistence": configOption(map[string]schema.Attribute{
			"enum": schema.ListAttribute{Computed: true, ElementType: types.StringType, Description: "The accepted values."},
		}),
		"valkey_pubsub_client_output_buffer_limit": intConfigOption("example"),
		"valkey_timeout": intConfigOption("default"),
	},
}
