package databasemysqlconfig

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSourceModel struct {
	ID                    types.String `tfsdk:"id"`
	MySQL                 types.List   `tfsdk:"mysql"`
	BinlogRetentionPeriod types.List   `tfsdk:"binlog_retention_period"`
}

func (data *DataSourceModel) ParseMySQLConfig(
	config *linodego.MySQLDatabaseConfigInfo, diags *diag.Diagnostics,
) {
	mySQL := flattenMySQL(&config.MySQL, diags)
	if diags.HasError() {
		return
	}
	data.MySQL = *mySQL

	binlogRetentionPeriod := flattenBinlogRetentionPeriod(&config.BinlogRetentionPeriod, diags)
	if diags.HasError() {
		return
	}
	data.BinlogRetentionPeriod = *binlogRetentionPeriod

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		diags.AddError("Error marshalling json", err.Error())
		return
	}

	hash := sha256.Sum256(jsonBytes)
	id := hex.EncodeToString(hash[:])

	data.ID = types.StringValue(id)
}

func flattenBinlogRetentionPeriod(binlogRetentionPeriod *linodego.MySQLDatabaseConfigInfoBinlogRetentionPeriod, diags *diag.Diagnostics) *basetypes.ListValue {
	result := map[string]attr.Value{
		"description":      types.StringValue(binlogRetentionPeriod.Description),
		"example":          types.Int64Value(int64(binlogRetentionPeriod.Example)),
		"maximum":          types.Int64Value(int64(binlogRetentionPeriod.Maximum)),
		"minimum":          types.Int64Value(int64(binlogRetentionPeriod.Minimum)),
		"requires_restart": types.BoolValue(binlogRetentionPeriod.RequiresRestart),
		"type":             types.StringValue(binlogRetentionPeriod.Type),
	}

	obj, d := types.ObjectValue(BinlogRetentionPeriodObjectType.AttrTypes, result)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	resultList := helper.GenericSliceToList(
		[]attr.Value{obj}, BinlogRetentionPeriodObjectType, helper.FwValueEchoConverter(), diags,
	)
	if diags.HasError() {
		return nil
	}

	return &resultList
}

func flattenMySQL(mysql *linodego.MySQLDatabaseConfigInfoMySQL, diags *diag.Diagnostics) *basetypes.ListValue {
	result := make(map[string]attr.Value)

	connectTimeout, d := flattenConnectTimeout(&mysql.ConnectTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	defaultTimeZone, d := flattenDefaultTimeZone(&mysql.DefaultTimeZone)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	groupConcatMaxLen, d := flattenGroupConcatMaxLen(&mysql.GroupConcatMaxLen)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	informationSchemaStatsExpiry, d := flattenInformationSchemaStatsExpiry(&mysql.InformationSchemaStatsExpiry)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBChangeBufferMaxSize, d := flattenInnoDBChangeBufferMaxSize(&mysql.InnoDBChangeBufferMaxSize)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBFlushNeighbors, d := flattenInnoDBFlushNeighbors(&mysql.InnoDBFlushNeighbors)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBFTMinTokenSize, d := flattenInnoDBFTMinTokenSize(&mysql.InnoDBFTMinTokenSize)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBFTServerStopwordTable, d := flattenInnoDBFTServerStopwordTable(&mysql.InnoDBFTServerStopwordTable)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBLockWaitTimeout, d := flattenInnoDBLockWaitTimeout(&mysql.InnoDBLockWaitTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBLogBufferSize, d := flattenInnoDBLogBufferSize(&mysql.InnoDBLogBufferSize)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBOnlineAlterLogMaxSize, d := flattenInnoDBOnlineAlterLogMaxSize(&mysql.InnoDBOnlineAlterLogMaxSize)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBReadIOThreads, d := flattenInnoDBReadIOThreads(&mysql.InnoDBReadIOThreads)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBRollbackOnTimeout, d := flattenInnoDBRollbackOnTimeout(&mysql.InnoDBRollbackOnTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBThreadConcurrency, d := flattenInnoDBThreadConcurrency(&mysql.InnoDBThreadConcurrency)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	innoDBWriteIOThreads, d := flattenInnoDBWriteIOThreads(&mysql.InnoDBWriteIOThreads)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	interactiveTimeout, d := flattenInteractiveTimeout(&mysql.InteractiveTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	internalTmpMemStorageEngine, d := flattenInternalTmpMemStorageEngine(&mysql.InternalTmpMemStorageEngine)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxAllowedPacket, d := flattenMaxAllowedPacket(&mysql.MaxAllowedPacket)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxHeapTableSize, d := flattenMaxHeapTableSize(&mysql.MaxHeapTableSize)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	netBufferLength, d := flattenNetBufferLength(&mysql.NetBufferLength)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	netReadTimeout, d := flattenNetReadTimeout(&mysql.NetReadTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	netWriteTimeout, d := flattenNetWriteTimeout(&mysql.NetWriteTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	sortBufferSize, d := flattenSortBufferSize(&mysql.SortBufferSize)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	sqlMode, d := flattenSQLMode(&mysql.SQLMode)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	sqlRequirePrimaryKey, d := flattenSQLRequirePrimaryKey(&mysql.SQLRequirePrimaryKey)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	tmpTableSize, d := flattenTmpTableSize(&mysql.TmpTableSize)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	waitTimeout, d := flattenWaitTimeout(&mysql.WaitTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	result["connect_timeout"] = connectTimeout
	result["default_time_zone"] = defaultTimeZone
	result["group_concat_max_len"] = groupConcatMaxLen
	result["information_schema_stats_expiry"] = informationSchemaStatsExpiry
	result["innodb_change_buffer_max_size"] = innoDBChangeBufferMaxSize
	result["innodb_flush_neighbors"] = innoDBFlushNeighbors
	result["innodb_ft_min_token_size"] = innoDBFTMinTokenSize
	result["innodb_ft_server_stopword_table"] = innoDBFTServerStopwordTable
	result["innodb_lock_wait_timeout"] = innoDBLockWaitTimeout
	result["innodb_log_buffer_size"] = innoDBLogBufferSize
	result["innodb_online_alter_log_max_size"] = innoDBOnlineAlterLogMaxSize
	result["innodb_read_io_threads"] = innoDBReadIOThreads
	result["innodb_rollback_on_timeout"] = innoDBRollbackOnTimeout
	result["innodb_thread_concurrency"] = innoDBThreadConcurrency
	result["innodb_write_io_threads"] = innoDBWriteIOThreads
	result["interactive_timeout"] = interactiveTimeout
	result["internal_tmp_mem_storage_engine"] = internalTmpMemStorageEngine
	result["max_allowed_packet"] = maxAllowedPacket
	result["max_heap_table_size"] = maxHeapTableSize
	result["net_buffer_length"] = netBufferLength
	result["net_read_timeout"] = netReadTimeout
	result["net_write_timeout"] = netWriteTimeout
	result["sort_buffer_size"] = sortBufferSize
	result["sql_mode"] = sqlMode
	result["sql_require_primary_key"] = sqlRequirePrimaryKey
	result["tmp_table_size"] = tmpTableSize
	result["wait_timeout"] = waitTimeout

	obj, d := types.ObjectValue(MySQLConfigMySQLObjectType.AttrTypes, result)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	resultList := helper.GenericSliceToList(
		[]attr.Value{obj}, MySQLConfigMySQLObjectType, helper.FwValueEchoConverter(), diags,
	)
	if diags.HasError() {
		return nil
	}

	return &resultList
}

func flattenConnectTimeout(val *linodego.ConnectTimeout) (
	*basetypes.ObjectValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	result["description"] = types.StringValue(val.Description)
	result["example"] = types.Int64Value(int64(val.Example))
	result["maximum"] = types.Int64Value(int64(val.Maximum))
	result["minimum"] = types.Int64Value(int64(val.Minimum))
	result["requires_restart"] = types.BoolValue(val.RequiresRestart)
	result["type"] = types.StringValue(val.Type)

	obj, d := types.ObjectValue(ConnectTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}

	return &obj, nil
}

func flattenDefaultTimeZone(val *linodego.DefaultTimeZone) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.StringValue(val.Example),
		"maxLength":        types.Int64Value(int64(val.MaxLength)),
		"minLength":        types.Int64Value(int64(val.MinLength)),
		"pattern":          types.StringValue(val.Pattern),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(DefaultTimeZoneObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenGroupConcatMaxLen(val *linodego.GroupConcatMaxLen) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Float64Value(val.Example),
		"maximum":          types.Float64Value(val.Maximum),
		"minimum":          types.Float64Value(val.Minimum),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(GroupConcatMaxLenObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInformationSchemaStatsExpiry(val *linodego.InformationSchemaStatsExpiry) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InformationSchemaStatsExpiryObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBChangeBufferMaxSize(val *linodego.InnoDBChangeBufferMaxSize) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBChangeBufferMaxSizeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBFlushNeighbors(val *linodego.InnoDBFlushNeighbors) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBFlushNeighborsObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBFTMinTokenSize(val *linodego.InnoDBFTMinTokenSize) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBFTMinTokenSizeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBFTServerStopwordTable(val *linodego.InnoDBFTServerStopwordTable) (*basetypes.ObjectValue, diag.Diagnostics) {
	typeVals := make([]attr.Value, len(val.Type))
	for i, s := range val.Type {
		typeVals[i] = types.StringValue(s)
	}
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.StringValue(val.Example),
		"maxLength":        types.Int64Value(int64(val.MaxLength)),
		"pattern":          types.StringValue(val.Pattern),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.ListValueMust(types.StringType, typeVals),
	}
	obj, d := types.ObjectValue(InnoDBFTServerStopwordTableObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBLockWaitTimeout(val *linodego.InnoDBLockWaitTimeout) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBLockWaitTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBLogBufferSize(val *linodego.InnoDBLogBufferSize) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBLogBufferSizeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBOnlineAlterLogMaxSize(val *linodego.InnoDBOnlineAlterLogMaxSize) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBOnlineAlterLogMaxSizeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBReadIOThreads(val *linodego.InnoDBReadIOThreads) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBReadIOThreadsObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBRollbackOnTimeout(val *linodego.InnoDBRollbackOnTimeout) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.BoolValue(val.Example),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBRollbackOnTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBThreadConcurrency(val *linodego.InnoDBThreadConcurrency) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBThreadConcurrencyObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInnoDBWriteIOThreads(val *linodego.InnoDBWriteIOThreads) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InnoDBWriteIOThreadsObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInteractiveTimeout(val *linodego.InteractiveTimeout) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InteractiveTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenInternalTmpMemStorageEngine(val *linodego.InternalTmpMemStorageEngine) (*basetypes.ObjectValue, diag.Diagnostics) {
	enumVals := make([]attr.Value, len(val.Enum))
	for i, v := range val.Enum {
		enumVals[i] = types.StringValue(v)
	}
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"enum":             types.ListValueMust(types.StringType, enumVals),
		"example":          types.StringValue(val.Example),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(InternalTmpMemStorageEngineObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxAllowedPacket(val *linodego.MaxAllowedPacket) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxAllowedPacketObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxHeapTableSize(val *linodego.MaxHeapTableSize) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxHeapTableSizeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenNetBufferLength(val *linodego.NetBufferLength) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(NetBufferLengthObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenNetReadTimeout(val *linodego.NetReadTimeout) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(NetReadTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenNetWriteTimeout(val *linodego.NetWriteTimeout) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(NetWriteTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenSortBufferSize(val *linodego.SortBufferSize) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(SortBufferSizeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenSQLMode(val *linodego.SQLMode) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.StringValue(val.Example),
		"maxLength":        types.Int64Value(int64(val.MaxLength)),
		"pattern":          types.StringValue(val.Pattern),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(SQLModeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenSQLRequirePrimaryKey(val *linodego.SQLRequirePrimaryKey) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.BoolValue(val.Example),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(SQLRequirePrimaryKeyObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenTmpTableSize(val *linodego.TmpTableSize) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(TmpTableSizeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenWaitTimeout(val *linodego.WaitTimeout) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(WaitTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}
