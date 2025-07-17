package databasepostgresqlconfig

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
	ID                      types.String `tfsdk:"id"`
	PG                      types.List   `tfsdk:"pg"`
	PGStatMonitorEnable     types.List   `tfsdk:"pg_stat_monitor_enable"`
	PGLookout               types.List   `tfsdk:"pglookout"`
	SharedBuffersPercentage types.List   `tfsdk:"shared_buffers_percentage"`
	WorkMem                 types.List   `tfsdk:"work_mem"`
}

func (data *DataSourceModel) ParsePostgreSQLConfig(
	config *linodego.PostgresDatabaseConfigInfo, diags *diag.Diagnostics,
) {
	pg := flattenPG(&config.PG, diags)
	if diags.HasError() {
		return
	}
	data.PG = *pg

	pgStatMonitorEnable := flattenPGStatMonitorEnable(&config.PGStatMonitorEnable, diags)
	if diags.HasError() {
		return
	}
	data.PGStatMonitorEnable = *pgStatMonitorEnable

	pgLookout := flattenPGLookout(&config.PGLookout, diags)
	if diags.HasError() {
		return
	}
	data.PGLookout = *pgLookout

	sharedBuffersPercentage := flattenSharedBuffersPercentage(&config.SharedBuffersPercentage, diags)
	if diags.HasError() {
		return
	}
	data.SharedBuffersPercentage = *sharedBuffersPercentage

	workMem := flattenWorkMem(&config.WorkMem, diags)
	if diags.HasError() {
		return
	}
	data.WorkMem = *workMem

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		diags.AddError("Error marshalling json", err.Error())
		return
	}

	hash := sha256.Sum256(jsonBytes)
	id := hex.EncodeToString(hash[:])

	data.ID = types.StringValue(id)
}

func flattenPGStatMonitorEnable(pgStatMonitorEnable *linodego.PostgresDatabaseConfigInfoPGStatMonitorEnable, diags *diag.Diagnostics) *basetypes.ListValue {
	result := map[string]attr.Value{
		"description":      types.StringValue(pgStatMonitorEnable.Description),
		"requires_restart": types.BoolValue(pgStatMonitorEnable.RequiresRestart),
		"type":             types.StringValue(pgStatMonitorEnable.Type),
	}

	obj, d := types.ObjectValue(PGStatMonitorEnableObjectType.AttrTypes, result)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	resultList := helper.GenericSliceToList(
		[]attr.Value{obj}, PGStatMonitorEnableObjectType, helper.FwValueEchoConverter(), diags,
	)
	if diags.HasError() {
		return nil
	}

	return &resultList
}

func flattenPGLookout(pgLookout *linodego.PostgresDatabaseConfigInfoPGLookout, diags *diag.Diagnostics) *basetypes.ListValue {
	maxFailoverReplicationTimeLag := flattenPGLookoutMaxFailoverReplicationTimeLag(&pgLookout.PGLookoutMaxFailoverReplicationTimeLag, diags)
	if maxFailoverReplicationTimeLag == nil {
		return nil
	}

	result := map[string]attr.Value{
		"max_failover_replication_time_lag": maxFailoverReplicationTimeLag,
	}

	obj, d := types.ObjectValue(PGLookoutObjectType.AttrTypes, result)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	resultList := helper.GenericSliceToList(
		[]attr.Value{obj}, PGLookoutObjectType, helper.FwValueEchoConverter(), diags,
	)
	if diags.HasError() {
		return nil
	}

	return &resultList
}

func flattenPGLookoutMaxFailoverReplicationTimeLag(val *linodego.PGLookoutMaxFailoverReplicationTimeLag, diags *diag.Diagnostics) attr.Value {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(val.Maximum),
		"minimum":          types.Int64Value(val.Minimum),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(PGLookoutMaxFailoverReplicationTimeLagObjectType.AttrTypes, result)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}
	return obj
}

func flattenSharedBuffersPercentage(
	sharedBuffersPercentage *linodego.PostgresDatabaseConfigInfoSharedBuffersPercentage,
	diags *diag.Diagnostics,
) *basetypes.ListValue {
	result := map[string]attr.Value{
		"description":      types.StringValue(sharedBuffersPercentage.Description),
		"example":          types.Float64Value(sharedBuffersPercentage.Example),
		"maximum":          types.Float64Value(sharedBuffersPercentage.Maximum),
		"minimum":          types.Float64Value(sharedBuffersPercentage.Minimum),
		"requires_restart": types.BoolValue(sharedBuffersPercentage.RequiresRestart),
		"type":             types.StringValue(sharedBuffersPercentage.Type),
	}

	obj, d := types.ObjectValue(SharedBuffersPercentageObjectType.AttrTypes, result)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	resultList := helper.GenericSliceToList(
		[]attr.Value{obj}, SharedBuffersPercentageObjectType, helper.FwValueEchoConverter(), diags,
	)
	if diags.HasError() {
		return nil
	}

	return &resultList
}

func flattenWorkMem(workMem *linodego.PostgresDatabaseConfigInfoWorkMem, diags *diag.Diagnostics) *basetypes.ListValue {
	result := map[string]attr.Value{
		"description":      types.StringValue(workMem.Description),
		"example":          types.Int64Value(int64(workMem.Example)),
		"maximum":          types.Int64Value(int64(workMem.Maximum)),
		"minimum":          types.Int64Value(int64(workMem.Minimum)),
		"requires_restart": types.BoolValue(workMem.RequiresRestart),
		"type":             types.StringValue(workMem.Type),
	}

	obj, d := types.ObjectValue(WorkMemObjectType.AttrTypes, result)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	resultList := helper.GenericSliceToList(
		[]attr.Value{obj}, WorkMemObjectType, helper.FwValueEchoConverter(), diags,
	)
	if diags.HasError() {
		return nil
	}

	return &resultList
}

func flattenPG(postgres *linodego.PostgresDatabaseConfigInfoPG, diags *diag.Diagnostics) *basetypes.ListValue {
	result := make(map[string]attr.Value)

	autovacuumAnalyzeScaleFactor, d := flattenAutovacuumAnalyzeScaleFactor(&postgres.AutovacuumAnalyzeScaleFactor)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	autovacuumAnalyzeThreshold, d := flattenAutovacuumAnalyzeThreshold(&postgres.AutovacuumAnalyzeThreshold)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	autovacuumMaxWorkers, d := flattenAutovacuumMaxWorkers(&postgres.AutovacuumMaxWorkers)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	autovacuumNaptime, d := flattenAutovacuumNaptime(&postgres.AutovacuumNaptime)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	autovacuumVacuumCostDelay, d := flattenAutovacuumVacuumCostDelay(&postgres.AutovacuumVacuumCostDelay)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	autovacuumVacuumCostLimit, d := flattenAutovacuumVacuumCostLimit(&postgres.AutovacuumVacuumCostLimit)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	autovacuumVacuumScaleFactor, d := flattenAutovacuumVacuumScaleFactor(&postgres.AutovacuumVacuumScaleFactor)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	autovacuumVacuumThreshold, d := flattenAutovacuumVacuumThreshold(&postgres.AutovacuumVacuumThreshold)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	bgWriterDelay, d := flattenBGWriterDelay(&postgres.BGWriterDelay)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	bgWriterFlushAfter, d := flattenBGWriterFlushAfter(&postgres.BGWriterFlushAfter)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	bgWriterLRUMaxPages, d := flattenBGWriterLRUMaxPages(&postgres.BGWriterLRUMaxPages)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	bgWriterLRUMultiplier, d := flattenBGWriterLRUMultiplier(&postgres.BGWriterLRUMultiplier)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	deadlockTimeout, d := flattenDeadlockTimeout(&postgres.DeadlockTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	defaultToastCompression, d := flattenDefaultToastCompression(&postgres.DefaultToastCompression)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	idleInTransactionSessionTimeout, d := flattenIdleInTransactionSessionTimeout(&postgres.IdleInTransactionSessionTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	jit, d := flattenJIT(&postgres.JIT)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxFilesPerProcess, d := flattenMaxFilesPerProcess(&postgres.MaxFilesPerProcess)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxLocksPerTransaction, d := flattenMaxLocksPerTransaction(&postgres.MaxLocksPerTransaction)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxLogicalReplicationWorkers, d := flattenMaxLogicalReplicationWorkers(&postgres.MaxLogicalReplicationWorkers)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxParallelWorkers, d := flattenMaxParallelWorkers(&postgres.MaxParallelWorkers)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxParallelWorkersPerGather, d := flattenMaxParallelWorkersPerGather(&postgres.MaxParallelWorkersPerGather)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxPredLocksPerTransaction, d := flattenMaxPredLocksPerTransaction(&postgres.MaxPredLocksPerTransaction)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxReplicationSlots, d := flattenMaxReplicationSlots(&postgres.MaxReplicationSlots)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxSlotWALKeepSize, d := flattenMaxSlotWALKeepSize(&postgres.MaxSlotWALKeepSize)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxStackDepth, d := flattenMaxStackDepth(&postgres.MaxStackDepth)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxStandbyArchiveDelay, d := flattenMaxStandbyArchiveDelay(&postgres.MaxStandbyArchiveDelay)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxStandbyStreamingDelay, d := flattenMaxStandbyStreamingDelay(&postgres.MaxStandbyStreamingDelay)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxWALSenders, d := flattenMaxWALSenders(&postgres.MaxWALSenders)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	maxWorkerProcesses, d := flattenMaxWorkerProcesses(&postgres.MaxWorkerProcesses)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	passwordEncryption, d := flattenPasswordEncryption(&postgres.PasswordEncryption)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	pgPartmanBGWInterval, d := flattenPGPartmanBGWInterval(&postgres.PGPartmanBGWInterval)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	pgPartmanBGWRole, d := flattenPGPartmanBGWRole(&postgres.PGPartmanBGWRole)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	pgStatMonitorPGSMEnableQueryPlan, d := flattenPGStatMonitorPGSMEnableQueryPlan(&postgres.PGStatMonitorPGSMEnableQueryPlan)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	pgStatMonitorPGSMMaxBuckets, d := flattenPGStatMonitorPGSMMaxBuckets(&postgres.PGStatMonitorPGSMMaxBuckets)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	pgStatStatementsTrack, d := flattenPGStatStatementsTrack(&postgres.PGStatStatementsTrack)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	tempFileLimit, d := flattenTempFileLimit(&postgres.TempFileLimit)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	timezone, d := flattenTimezone(&postgres.Timezone)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	trackActivityQuerySize, d := flattenTrackActivityQuerySize(&postgres.TrackActivityQuerySize)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	trackCommitTimestamp, d := flattenTrackCommitTimestamp(&postgres.TrackCommitTimestamp)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	trackFunctions, d := flattenTrackFunctions(&postgres.TrackFunctions)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	trackIOTiming, d := flattenTrackIOTiming(&postgres.TrackIOTiming)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	walSenderTimeout, d := flattenWALSenderTimeout(&postgres.WALSenderTimeout)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	walWriterDelay, d := flattenWALWriterDelay(&postgres.WALWriterDelay)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	result["autovacuum_analyze_scale_factor"] = autovacuumAnalyzeScaleFactor
	result["autovacuum_analyze_threshold"] = autovacuumAnalyzeThreshold
	result["autovacuum_max_workers"] = autovacuumMaxWorkers
	result["autovacuum_naptime"] = autovacuumNaptime
	result["autovacuum_vacuum_cost_delay"] = autovacuumVacuumCostDelay
	result["autovacuum_vacuum_cost_limit"] = autovacuumVacuumCostLimit
	result["autovacuum_vacuum_scale_factor"] = autovacuumVacuumScaleFactor
	result["autovacuum_vacuum_threshold"] = autovacuumVacuumThreshold
	result["bgwriter_delay"] = bgWriterDelay
	result["bgwriter_flush_after"] = bgWriterFlushAfter
	result["bgwriter_lru_maxpages"] = bgWriterLRUMaxPages
	result["bgwriter_lru_multiplier"] = bgWriterLRUMultiplier
	result["deadlock_timeout"] = deadlockTimeout
	result["default_toast_compression"] = defaultToastCompression
	result["idle_in_transaction_session_timeout"] = idleInTransactionSessionTimeout
	result["jit"] = jit
	result["max_files_per_process"] = maxFilesPerProcess
	result["max_locks_per_transaction"] = maxLocksPerTransaction
	result["max_logical_replication_workers"] = maxLogicalReplicationWorkers
	result["max_parallel_workers"] = maxParallelWorkers
	result["max_parallel_workers_per_gather"] = maxParallelWorkersPerGather
	result["max_pred_locks_per_transaction"] = maxPredLocksPerTransaction
	result["max_replication_slots"] = maxReplicationSlots
	result["max_slot_wal_keep_size"] = maxSlotWALKeepSize
	result["max_stack_depth"] = maxStackDepth
	result["max_standby_archive_delay"] = maxStandbyArchiveDelay
	result["max_standby_streaming_delay"] = maxStandbyStreamingDelay
	result["max_wal_senders"] = maxWALSenders
	result["max_worker_processes"] = maxWorkerProcesses
	result["password_encryption"] = passwordEncryption
	result["pg_partman_bgw.interval"] = pgPartmanBGWInterval
	result["pg_partman_bgw.role"] = pgPartmanBGWRole
	result["pg_stat_monitor.pgsm_enable_query_plan"] = pgStatMonitorPGSMEnableQueryPlan
	result["pg_stat_monitor.pgsm_max_buckets"] = pgStatMonitorPGSMMaxBuckets
	result["pg_stat_statements.track"] = pgStatStatementsTrack
	result["temp_file_limit"] = tempFileLimit
	result["timezone"] = timezone
	result["track_activity_query_size"] = trackActivityQuerySize
	result["track_commit_timestamp"] = trackCommitTimestamp
	result["track_functions"] = trackFunctions
	result["track_io_timing"] = trackIOTiming
	result["wal_sender_timeout"] = walSenderTimeout
	result["wal_writer_delay"] = walWriterDelay

	obj, d := types.ObjectValue(PostgreSQLConfigPGObjectType.AttrTypes, result)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	resultList := helper.GenericSliceToList(
		[]attr.Value{obj}, PostgreSQLConfigPGObjectType, helper.FwValueEchoConverter(), diags,
	)
	if diags.HasError() {
		return nil
	}

	return &resultList
}

func flattenAutovacuumAnalyzeScaleFactor(val *linodego.AutovacuumAnalyzeScaleFactor) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Float64Value(val.Maximum),
		"minimum":          types.Float64Value(val.Minimum),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(AutovacuumAnalyzeScaleFactorObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenAutovacuumAnalyzeThreshold(val *linodego.AutovacuumAnalyzeThreshold) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int32Value(val.Maximum),
		"minimum":          types.Int32Value(val.Minimum),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(AutovacuumAnalyzeThresholdObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenAutovacuumMaxWorkers(val *linodego.AutovacuumMaxWorkers) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(AutovacuumMaxWorkersObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenAutovacuumNaptime(val *linodego.AutovacuumNaptime) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(AutovacuumNaptimeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenAutovacuumVacuumCostDelay(val *linodego.AutovacuumVacuumCostDelay) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(AutovacuumVacuumCostDelayObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenAutovacuumVacuumCostLimit(val *linodego.AutovacuumVacuumCostLimit) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(AutovacuumVacuumCostLimitObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenAutovacuumVacuumScaleFactor(val *linodego.AutovacuumVacuumScaleFactor) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Float64Value(val.Maximum),
		"minimum":          types.Float64Value(val.Minimum),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(AutovacuumVacuumScaleFactorObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenAutovacuumVacuumThreshold(val *linodego.AutovacuumVacuumThreshold) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int32Value(val.Maximum),
		"minimum":          types.Int32Value(val.Minimum),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(AutovacuumVacuumThresholdObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenBGWriterDelay(val *linodego.BGWriterDelay) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(BGWriterDelayObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenBGWriterFlushAfter(val *linodego.BGWriterFlushAfter) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(BGWriterFlushAfterObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenBGWriterLRUMaxPages(val *linodego.BGWriterLRUMaxPages) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(BGWriterLRUMaxPagesObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenBGWriterLRUMultiplier(val *linodego.BGWriterLRUMultiplier) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Float64Value(val.Example),
		"maximum":          types.Float64Value(val.Maximum),
		"minimum":          types.Float64Value(val.Minimum),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(BGWriterLRUMultiplierObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenDeadlockTimeout(val *linodego.DeadlockTimeout) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(DeadlockTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenDefaultToastCompression(val *linodego.DefaultToastCompression) (*basetypes.ObjectValue, diag.Diagnostics) {
	enumVals := make([]attr.Value, len(val.Enum))
	for i, s := range val.Enum {
		enumVals[i] = types.StringValue(s)
	}
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"enum":             types.ListValueMust(types.StringType, enumVals),
		"example":          types.StringValue(val.Example),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(DefaultToastCompressionObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenIdleInTransactionSessionTimeout(val *linodego.IdleInTransactionSessionTimeout) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(IdleInTransactionSessionTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenJIT(val *linodego.JIT) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.BoolValue(val.Example),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(JITObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxFilesPerProcess(val *linodego.MaxFilesPerProcess) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxFilesPerProcessObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxLocksPerTransaction(val *linodego.MaxLocksPerTransaction) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxLocksPerTransactionObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxLogicalReplicationWorkers(val *linodego.MaxLogicalReplicationWorkers) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxLogicalReplicationWorkersObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxParallelWorkers(val *linodego.MaxParallelWorkers) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxParallelWorkersObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxParallelWorkersPerGather(val *linodego.MaxParallelWorkersPerGather) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxParallelWorkersPerGatherObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxPredLocksPerTransaction(val *linodego.MaxPredLocksPerTransaction) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxPredLocksPerTransactionObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxReplicationSlots(val *linodego.MaxReplicationSlots) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxReplicationSlotsObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxSlotWALKeepSize(val *linodego.MaxSlotWALKeepSize) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int32Value(int32(val.Maximum)),
		"minimum":          types.Int32Value(int32(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxSlotWALKeepSizeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxStackDepth(val *linodego.MaxStackDepth) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxStackDepthObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxStandbyArchiveDelay(val *linodego.MaxStandbyArchiveDelay) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxStandbyArchiveDelayObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxStandbyStreamingDelay(val *linodego.MaxStandbyStreamingDelay) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxStandbyStreamingDelayObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxWALSenders(val *linodego.MaxWALSenders) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxWALSendersObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenMaxWorkerProcesses(val *linodego.MaxWorkerProcesses) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(MaxWorkerProcessesObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenPasswordEncryption(val *linodego.PasswordEncryption) (*basetypes.ObjectValue, diag.Diagnostics) {
	enumVals := make([]attr.Value, len(val.Enum))
	for i, s := range val.Enum {
		enumVals[i] = types.StringValue(s)
	}
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"enum":             types.ListValueMust(types.StringType, enumVals),
		"example":          types.StringValue(val.Example),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(PasswordEncryptionObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenPGPartmanBGWInterval(val *linodego.PGPartmanBGWInterval) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(PGPartmanBGWIntervalObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenPGPartmanBGWRole(val *linodego.PGPartmanBGWRole) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.StringValue(val.Example),
		"maxLength":        types.Int64Value(int64(val.MaxLength)),
		"pattern":          types.StringValue(val.Pattern),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(PGPartmanBGWRoleObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenPGStatMonitorPGSMEnableQueryPlan(val *linodego.PGStatMonitorPGSMEnableQueryPlan) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.BoolValue(val.Example),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(PGStatMonitorPGSMEnableQueryPlanObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenPGStatMonitorPGSMMaxBuckets(val *linodego.PGStatMonitorPGSMMaxBuckets) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(PGStatMonitorPGSMMaxBucketsObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenPGStatStatementsTrack(val *linodego.PGStatStatementsTrack) (*basetypes.ObjectValue, diag.Diagnostics) {
	enumVals := make([]attr.Value, len(val.Enum))
	for i, s := range val.Enum {
		enumVals[i] = types.StringValue(s)
	}
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"enum":             types.ListValueMust(types.StringType, enumVals),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(PGStatStatementsTrackObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenTempFileLimit(val *linodego.TempFileLimit) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int32Value(val.Example),
		"maximum":          types.Int32Value(val.Maximum),
		"minimum":          types.Int32Value(val.Minimum),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(TempFileLimitObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenTimezone(val *linodego.Timezone) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.StringValue(val.Example),
		"maxLength":        types.Int64Value(int64(val.MaxLength)),
		"pattern":          types.StringValue(val.Pattern),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(TimezoneObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenTrackActivityQuerySize(val *linodego.TrackActivityQuerySize) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(TrackActivityQuerySizeObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenTrackCommitTimestamp(val *linodego.TrackCommitTimestamp) (*basetypes.ObjectValue, diag.Diagnostics) {
	enumVals := make([]attr.Value, len(val.Enum))
	for i, s := range val.Enum {
		enumVals[i] = types.StringValue(s)
	}
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"enum":             types.ListValueMust(types.StringType, enumVals),
		"example":          types.StringValue(val.Example),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(TrackCommitTimestampObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenTrackFunctions(val *linodego.TrackFunctions) (*basetypes.ObjectValue, diag.Diagnostics) {
	enumVals := make([]attr.Value, len(val.Enum))
	for i, s := range val.Enum {
		enumVals[i] = types.StringValue(s)
	}
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"enum":             types.ListValueMust(types.StringType, enumVals),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(TrackFunctionsObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenTrackIOTiming(val *linodego.TrackIOTiming) (*basetypes.ObjectValue, diag.Diagnostics) {
	enumVals := make([]attr.Value, len(val.Enum))
	for i, s := range val.Enum {
		enumVals[i] = types.StringValue(s)
	}
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"enum":             types.ListValueMust(types.StringType, enumVals),
		"example":          types.StringValue(val.Example),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(TrackIOTimingObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenWALSenderTimeout(val *linodego.WALSenderTimeout) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(WALSenderTimeoutObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}

func flattenWALWriterDelay(val *linodego.WALWriterDelay) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := map[string]attr.Value{
		"description":      types.StringValue(val.Description),
		"example":          types.Int64Value(int64(val.Example)),
		"maximum":          types.Int64Value(int64(val.Maximum)),
		"minimum":          types.Int64Value(int64(val.Minimum)),
		"requires_restart": types.BoolValue(val.RequiresRestart),
		"type":             types.StringValue(val.Type),
	}
	obj, d := types.ObjectValue(WALWriterDelayObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}
	return &obj, nil
}
