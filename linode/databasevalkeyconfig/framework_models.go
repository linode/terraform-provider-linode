package databasevalkeyconfig

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

type DataSourceModel struct {
	ID                                  types.String `tfsdk:"id"`
	BackupHour                          types.List   `tfsdk:"backup_hour"`
	BackupMinute                        types.List   `tfsdk:"backup_minute"`
	FrequentSnapshots                   types.List   `tfsdk:"frequent_snapshots"`
	ValkeyACLChannelsDefault            types.List   `tfsdk:"valkey_acl_channels_default"`
	ValkeyActiveExpireEffort            types.List   `tfsdk:"valkey_active_expire_effort"`
	ValkeyActiveDefrag                  types.List   `tfsdk:"valkey_activedefrag"`
	ValkeyLFUDecayTime                  types.List   `tfsdk:"valkey_lfu_decay_time"`
	ValkeyLFULogFactor                  types.List   `tfsdk:"valkey_lfu_log_factor"`
	ValkeyMaxmemoryPolicy               types.List   `tfsdk:"valkey_maxmemory_policy"`
	ValkeyNumberOfDatabases             types.List   `tfsdk:"valkey_number_of_databases"`
	ValkeyPersistence                   types.List   `tfsdk:"valkey_persistence"`
	ValkeyPubsubClientOutputBufferLimit types.List   `tfsdk:"valkey_pubsub_client_output_buffer_limit"`
	ValkeyTimeout                       types.List   `tfsdk:"valkey_timeout"`
}

func flattenOption(ctx context.Context, name, description string, restart bool, paramType linodego.DatabaseConfigParamType,
	fields map[string]attr.Value, diags *diag.Diagnostics,
) types.List {
	fields["description"] = types.StringValue(description)
	fields["requires_restart"] = types.BoolValue(restart)
	typeList, diagnostics := types.ListValueFrom(ctx, types.StringType, []string(paramType))
	diags.Append(diagnostics...)
	fields["type"] = typeList
	objectType := frameworkDataSourceSchema.Attributes[name].GetType().(types.ListType).ElemType.(types.ObjectType)
	object, diagnostics := types.ObjectValue(objectType.AttrTypes, fields)
	diags.Append(diagnostics...)
	result, diagnostics := types.ListValue(objectType, []attr.Value{object})
	diags.Append(diagnostics...)
	return result
}

func intMetadata(valueField string, value, minimum, maximum int) map[string]attr.Value {
	return map[string]attr.Value{
		valueField: types.Int64Value(int64(value)),
		"minimum":  types.Int64Value(int64(minimum)),
		"maximum":  types.Int64Value(int64(maximum)),
	}
}

func (data *DataSourceModel) ParseValkeyConfig(ctx context.Context, config *linodego.ValkeyDatabaseConfigInfo, diags *diag.Diagnostics) {
	data.BackupHour = flattenOption(ctx, "backup_hour", config.BackupHour.Description, config.BackupHour.RequiresRestart, config.BackupHour.Type,
		intMetadata("example", config.BackupHour.Example, config.BackupHour.Minimum, config.BackupHour.Maximum), diags)
	data.BackupMinute = flattenOption(ctx, "backup_minute", config.BackupMinute.Description, config.BackupMinute.RequiresRestart, config.BackupMinute.Type,
		intMetadata("example", config.BackupMinute.Example, config.BackupMinute.Minimum, config.BackupMinute.Maximum), diags)
	data.FrequentSnapshots = flattenOption(ctx, "frequent_snapshots", config.FrequentSnapshots.Description, config.FrequentSnapshots.RequiresRestart,
		config.FrequentSnapshots.Type, map[string]attr.Value{"default": types.BoolValue(config.FrequentSnapshots.Default)}, diags)
	aclEnum, diagnostics := types.ListValueFrom(ctx, types.StringType, config.ValkeyACLChannelsDefault.Enum)
	diags.Append(diagnostics...)
	data.ValkeyACLChannelsDefault = flattenOption(ctx, "valkey_acl_channels_default", config.ValkeyACLChannelsDefault.Description,
		config.ValkeyACLChannelsDefault.RequiresRestart, config.ValkeyACLChannelsDefault.Type, map[string]attr.Value{"enum": aclEnum}, diags)
	data.ValkeyActiveExpireEffort = flattenOption(
		ctx,
		"valkey_active_expire_effort",
		config.ValkeyActiveExpireEffort.Description,
		config.ValkeyActiveExpireEffort.RequiresRestart,
		config.ValkeyActiveExpireEffort.Type,
		intMetadata(
			"default",
			config.ValkeyActiveExpireEffort.Default,
			config.ValkeyActiveExpireEffort.Minimum,
			config.ValkeyActiveExpireEffort.Maximum,
		),
		diags,
	)
	data.ValkeyActiveDefrag = flattenOption(ctx, "valkey_activedefrag", config.ValkeyActiveDefrag.Description, config.ValkeyActiveDefrag.RequiresRestart,
		config.ValkeyActiveDefrag.Type, map[string]attr.Value{"default": types.BoolValue(config.ValkeyActiveDefrag.Default)}, diags)
	data.ValkeyLFUDecayTime = flattenOption(
		ctx,
		"valkey_lfu_decay_time",
		config.ValkeyLFUDecayTime.Description,
		config.ValkeyLFUDecayTime.RequiresRestart,
		config.ValkeyLFUDecayTime.Type,
		intMetadata("default", config.ValkeyLFUDecayTime.Default, config.ValkeyLFUDecayTime.Minimum, config.ValkeyLFUDecayTime.Maximum),
		diags,
	)
	data.ValkeyLFULogFactor = flattenOption(
		ctx,
		"valkey_lfu_log_factor",
		config.ValkeyLFULogFactor.Description,
		config.ValkeyLFULogFactor.RequiresRestart,
		config.ValkeyLFULogFactor.Type,
		intMetadata("default", config.ValkeyLFULogFactor.Default, config.ValkeyLFULogFactor.Minimum, config.ValkeyLFULogFactor.Maximum),
		diags,
	)
	policyEnum, diagnostics := types.ListValueFrom(ctx, types.StringType, config.ValkeyMaxmemoryPolicy.Enum)
	diags.Append(diagnostics...)
	data.ValkeyMaxmemoryPolicy = flattenOption(ctx, "valkey_maxmemory_policy", config.ValkeyMaxmemoryPolicy.Description,
		config.ValkeyMaxmemoryPolicy.RequiresRestart, config.ValkeyMaxmemoryPolicy.Type,
		map[string]attr.Value{"default": types.StringValue(config.ValkeyMaxmemoryPolicy.Default), "enum": policyEnum}, diags)
	data.ValkeyNumberOfDatabases = flattenOption(ctx, "valkey_number_of_databases", config.ValkeyNumberOfDatabases.Description,
		config.ValkeyNumberOfDatabases.RequiresRestart, config.ValkeyNumberOfDatabases.Type,
		intMetadata("example", config.ValkeyNumberOfDatabases.Example, config.ValkeyNumberOfDatabases.Minimum, config.ValkeyNumberOfDatabases.Maximum), diags)
	persistenceEnum, diagnostics := types.ListValueFrom(ctx, types.StringType, config.ValkeyPersistence.Enum)
	diags.Append(diagnostics...)
	data.ValkeyPersistence = flattenOption(ctx, "valkey_persistence", config.ValkeyPersistence.Description, config.ValkeyPersistence.RequiresRestart,
		config.ValkeyPersistence.Type, map[string]attr.Value{"enum": persistenceEnum}, diags)
	data.ValkeyPubsubClientOutputBufferLimit = flattenOption(
		ctx,
		"valkey_pubsub_client_output_buffer_limit",
		config.ValkeyPubsubClientOutputBufferLimit.Description,
		config.ValkeyPubsubClientOutputBufferLimit.RequiresRestart,
		config.ValkeyPubsubClientOutputBufferLimit.Type,
		intMetadata("example", config.ValkeyPubsubClientOutputBufferLimit.Example, config.ValkeyPubsubClientOutputBufferLimit.Minimum,
			config.ValkeyPubsubClientOutputBufferLimit.Maximum),
		diags,
	)
	data.ValkeyTimeout = flattenOption(ctx, "valkey_timeout", config.ValkeyTimeout.Description, config.ValkeyTimeout.RequiresRestart, config.ValkeyTimeout.Type,
		intMetadata("default", config.ValkeyTimeout.Default, config.ValkeyTimeout.Minimum, config.ValkeyTimeout.Maximum), diags)
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		diags.AddError("Error marshaling json", err.Error())
		return
	}
	hash := sha256.Sum256(jsonBytes)
	data.ID = types.StringValue(hex.EncodeToString(hash[:]))
}
