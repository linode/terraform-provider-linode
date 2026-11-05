package databasevalkeyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/databaseshared"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"engine_id": schema.StringAttribute{
			Required:    true,
			Description: "The unique ID of the Valkey engine and version to use (for example, valkey/8.1).",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Required:    true,
			Description: "A unique, user-defined label for the Managed Valkey Database.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Required:    true,
			Description: "The region where the Managed Valkey Database will be deployed.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"type": schema.StringAttribute{
			Required:    true,
			Description: "The Linode plan used by the Managed Valkey Database nodes.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"allow_list": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Description: "IP addresses or CIDR ranges allowed to access the Managed Valkey Database.",
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		},
		"ca_cert": schema.StringAttribute{
			Computed:      true,
			Sensitive:     true,
			Description:   "The base64-encoded SSL CA certificate for the Managed Valkey Database.",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown()},
		},
		"cluster_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of nodes in the Managed Valkey Database cluster.",
			Validators:  []validator.Int64{int64validator.AtLeast(1)},
			Default:     int64default.StaticInt64(1),
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"fork_restore_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			Description: "A Valkey RDB snapshot time to restore from. It must be one of the source database's available_restore_times; when omitted, the latest snapshot is used.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplaceIf(
					func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
						resp.RequiresReplace = !helper.CompareRFC3339TimeStrings(
							req.PlanValue.ValueString(),
							req.StateValue.ValueString(),
						)
					},
					"Triggers replacement when `fork_restore_time` changes",
					"Changing `fork_restore_time` forces a new resource.",
				),
			},
		},
		"fork_source": schema.Int64Attribute{
			Optional:    true,
			Description: "The ID of a Valkey database to fork.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
				int64planmodifier.RequiresReplace(),
			},
		},
		"suspended": schema.BoolAttribute{
			Optional:      true,
			Computed:      true,
			Default:       booldefault.StaticBool(false),
			Description:   "Whether the Managed Valkey Database is suspended.",
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
		},
		"updates": databaseshared.ResourceAttributeUpdates,
		"created": schema.StringAttribute{
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			Description: "When the Managed Valkey Database was created.",
		},
		"encrypted": schema.BoolAttribute{Computed: true, Description: "Whether the Managed Valkey Database is encrypted."},
		"engine":    schema.StringAttribute{Computed: true, Description: "The database engine."},
		"host_primary": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: databaseshared.HostStringPlanModifiers,
			Description:   "The primary Valkey connection host.",
		},
		"host_secondary": schema.StringAttribute{
			Computed:           true,
			DeprecationMessage: "Use host_standby instead.",
			Description:        "The secondary Valkey connection host.",
		},
		"host_standby": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: databaseshared.HostStringPlanModifiers,
			Description:   "The standby Valkey connection host.",
		},
		"members":         schema.MapAttribute{Computed: true, ElementType: types.StringType, Description: "A mapping of cluster members to their roles."},
		"platform":        schema.StringAttribute{Computed: true, Description: "The backend platform for the Managed Valkey Database."},
		"port":            schema.Int64Attribute{Computed: true, Description: "The port used to connect to Valkey."},
		"private_network": databaseshared.ResourceAttributePrivateNetwork,
		"root_password": schema.StringAttribute{
			Computed:      true,
			Sensitive:     true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown()},
			Description:   "The generated root password. This value is stored in Terraform state.",
		},
		"root_username": schema.StringAttribute{
			Computed:      true,
			Sensitive:     true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown()},
			Description:   "The root username. This value is stored in Terraform state.",
		},
		"ssl_connection": schema.BoolAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			Description:   "Whether SSL is required for database connections.",
		},
		"status": schema.StringAttribute{Computed: true, Description: "The current status of the Managed Valkey Database."},
		"updated": schema.StringAttribute{
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			Description: "When the Managed Valkey Database was last updated.",
		},
		"version": schema.StringAttribute{Computed: true, Description: "The Valkey server version."},
		"available_restore_times": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Available RDB snapshot timestamps that can be used as discrete restore points when forking this database.",
		},
		"pending_updates": databaseshared.ResourceAttributePendingUpdates,

		"engine_config_backup_hour":        valkeyIntConfig("The UTC hour when a backup starts.", 0, 23),
		"engine_config_backup_minute":      valkeyIntConfig("The minute within the hour when a backup starts.", 0, 59),
		"engine_config_frequent_snapshots": valkeyBoolConfig("Whether frequent local RDB snapshots are enabled."),
		"engine_config_valkey_acl_channels_default": valkeyStringConfig(
			"The default pub/sub channel ACL for new users.",
			stringvalidator.OneOf("allchannels", "resetchannels"),
		),
		"engine_config_valkey_active_expire_effort": valkeyIntConfig("The effort used to reclaim expired keys in the background.", 1, 10),
		"engine_config_valkey_activedefrag":         valkeyBoolConfig("Whether active memory defragmentation is enabled."),
		"engine_config_valkey_lfu_decay_time":       valkeyIntConfig("The LFU counter decay time in minutes.", 1, 120),
		"engine_config_valkey_lfu_log_factor":       valkeyIntConfig("The logarithm factor for LFU counters.", 0, 100),
		"engine_config_valkey_maxmemory_policy": valkeyStringConfig(
			"The policy Valkey uses when memory is full.",
			stringvalidator.OneOf(
				"noeviction",
				"allkeys-lru",
				"volatile-lru",
				"allkeys-random",
				"volatile-random",
				"volatile-ttl",
				"volatile-lfu",
				"allkeys-lfu",
			),
		),
		"engine_config_valkey_number_of_databases": valkeyIntConfig(
			"The number of logical Valkey databases. Changing this setting restarts the service.",
			1,
			128,
		),
		"engine_config_valkey_persistence": valkeyStringConfig(
			"The persistence mode. `off` disables backups and forking.",
			stringvalidator.OneOf("off", "rdb"),
		),
		"engine_config_valkey_pubsub_client_output_buffer_limit": valkeyIntConfig("The hard output-buffer limit for pub/sub clients in MB.", 32, 262144),
		"engine_config_valkey_timeout":                           valkeyIntConfig("The idle connection timeout in seconds.", 0, 2073600),
	},
}

func valkeyIntConfig(description string, min, max int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional: true, Computed: true, Description: description,
		Validators:    []validator.Int64{int64validator.Between(min, max)},
		PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
	}
}

func valkeyBoolConfig(description string) schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional: true, Computed: true, Description: description,
		PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
	}
}

func valkeyStringConfig(description string, valueValidator validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Optional: true, Computed: true, Description: description,
		Validators:    []validator.String{valueValidator},
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
}
