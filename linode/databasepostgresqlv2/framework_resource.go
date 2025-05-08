package databasepostgresqlv2

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

const (
	DefaultCreateTimeout = 60 * time.Minute
	DefaultUpdateTimeout = 60 * time.Minute
	DefaultDeleteTimeout = 5 * time.Minute
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_database_postgresql_v2",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
				TimeoutOpts: &timeouts.Opts{
					Update: true,
					Create: true,
					Delete: true,
				},
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, DefaultCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	createOpts := linodego.PostgresCreateOptions{
		Label:        data.Label.ValueString(),
		Region:       data.Region.ValueString(),
		Type:         data.Type.ValueString(),
		Engine:       data.EngineID.ValueString(),
		ClusterSize:  helper.FrameworkSafeInt64ToInt(data.ClusterSize.ValueInt64(), &resp.Diagnostics),
		Fork:         data.GetFork(resp.Diagnostics),
		AllowList:    data.GetAllowList(ctx, resp.Diagnostics),
		EngineConfig: data.GetEngineConfig(resp.Diagnostics),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createPoller, err := client.NewEventPollerWithoutEntity(linodego.EntityDatabase, linodego.ActionDatabaseCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create event poller",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "client.CreatePostgresDatabase(...)", map[string]any{
		"options": createOpts,
	})

	db, err := client.CreatePostgresDatabase(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create PostgreSQL database",
			err.Error(),
		)
		return
	}

	// We explicitly set the ID in state here to prevent leaking the resource
	// in the case of a polling failure
	resp.State.SetAttribute(ctx, path.Root("id"), strconv.Itoa(db.ID))

	ctx = tflog.SetField(ctx, "id", db.ID)

	createPoller.EntityID = db.ID

	tflog.Debug(ctx, "Waiting for database to finish provisioning")

	if _, err := createPoller.WaitForFinished(ctx, int(createTimeout.Seconds())); err != nil {
		resp.Diagnostics.AddError(
			"Failed to wait for PostgreSQL database to finish creating",
			err.Error(),
		)
	}

	// Sometimes the creation event finishes before the status becomes `active`
	tflog.Debug(ctx, "Waiting for database to enter active status", map[string]any{
		"options": createOpts,
	})

	if err = client.WaitForDatabaseStatus(
		ctx,
		db.ID,
		linodego.DatabaseEngineTypePostgres,
		linodego.DatabaseStatusActive,
		int(createTimeout.Seconds()),
	); err != nil {
		resp.Diagnostics.AddError("Failed to wait for PostgreSQL database active", err.Error())
		return
	}

	// The `updates` field can only be changed using PUT requests
	updates := data.GetUpdates(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if updates != nil {
		updateOpts := linodego.PostgresUpdateOptions{Updates: updates.ToLinodego(resp.Diagnostics)}
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "client.UpdatePostgresDatabase(...)", map[string]any{
			"options": updateOpts,
		})

		db, err = client.UpdatePostgresDatabase(
			ctx,
			db.ID,
			updateOpts,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to update PostgreSQL database",
				err.Error(),
			)
			return
		}
	}

	if err := helper.ReconcileDatabaseSuspensionSync(
		ctx,
		client,
		db.ID,
		linodego.DatabaseEngineTypePostgres,
		false,
		data.Suspended.ValueBool(),
		createTimeout,
	); err != nil {
		resp.Diagnostics.AddError(
			"Failed to reconcile database suspension",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.Refresh(ctx, client, db.ID, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	data.ID = types.StringValue(strconv.Itoa(db.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.GetPostgresDatabase(...)")

	db, err := client.GetPostgresDatabase(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Database no longer exists",
				fmt.Sprintf(
					"Removing PostgreSQL database with ID %v from state because it no longer exists",
					id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the Database",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.Refresh(ctx, client, db.ID, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	client := r.Meta.Client
	var plan, state ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.FrameworkSafeStringToInt(plan.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, DefaultUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	ctx = populateLogAttributes(ctx, state)

	// Suspend or resume the database if necessary
	// NOTE: Other fields cannot be updated if suspended = true
	if err := helper.ReconcileDatabaseSuspensionSync(
		ctx,
		client,
		id,
		linodego.DatabaseEngineTypePostgres,
		state.Suspended.ValueBool(),
		plan.Suspended.ValueBool(),
		updateTimeout,
	); err != nil {
		resp.Diagnostics.AddError(
			"Failed to reconcile database suspension",
			err.Error(),
		)
		return
	}

	var updateOpts linodego.PostgresUpdateOptions
	shouldUpdate := false
	shouldResize := false

	// `label` field updates
	if !state.Label.Equal(plan.Label) {
		shouldUpdate = true
		updateOpts.Label = plan.Label.ValueString()
	}

	// `allow_list` field updates
	if !state.AllowList.Equal(plan.AllowList) {
		shouldUpdate = true

		allowList := plan.GetAllowList(ctx, resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		updateOpts.AllowList = &allowList
	}

	// `type` field updates
	if !state.Type.Equal(plan.Type) {
		shouldResize = true
		updateOpts.Type = plan.Type.ValueString()
	}

	// `updates` field updates
	if !state.Updates.Equal(plan.Updates) {
		shouldUpdate = true

		updates := plan.GetUpdates(ctx, resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		updateOpts.Updates = updates.ToLinodego(resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	engineConfigFields := []bool{
		!cmp.Equal(state.EngineConfigPGAutovacuumAnalyzeScaleFactor, plan.EngineConfigPGAutovacuumAnalyzeScaleFactor),
		!cmp.Equal(state.EngineConfigPGAutovacuumAnalyzeThreshold, plan.EngineConfigPGAutovacuumAnalyzeThreshold),
		!cmp.Equal(state.EngineConfigPGAutovacuumMaxWorkers, plan.EngineConfigPGAutovacuumMaxWorkers),
		!cmp.Equal(state.EngineConfigPGAutovacuumNaptime, plan.EngineConfigPGAutovacuumNaptime),
		!cmp.Equal(state.EngineConfigPGAutovacuumVacuumCostDelay, plan.EngineConfigPGAutovacuumVacuumCostDelay),
		!cmp.Equal(state.EngineConfigPGAutovacuumVacuumCostLimit, plan.EngineConfigPGAutovacuumVacuumCostLimit),
		!cmp.Equal(state.EngineConfigPGAutovacuumVacuumScaleFactor, plan.EngineConfigPGAutovacuumVacuumScaleFactor),
		!cmp.Equal(state.EngineConfigPGAutovacuumVacuumThreshold, plan.EngineConfigPGAutovacuumVacuumThreshold),
		!cmp.Equal(state.EngineConfigPGBGWriterDelay, plan.EngineConfigPGBGWriterDelay),
		!cmp.Equal(state.EngineConfigPGBGWriterFlushAfter, plan.EngineConfigPGBGWriterFlushAfter),
		!cmp.Equal(state.EngineConfigPGBGWriterLRUMaxpages, plan.EngineConfigPGBGWriterLRUMaxpages),
		!cmp.Equal(state.EngineConfigPGBGWriterLRUMultiplier, plan.EngineConfigPGBGWriterLRUMultiplier),
		!cmp.Equal(state.EngineConfigPGDeadlockTimeout, plan.EngineConfigPGDeadlockTimeout),
		!cmp.Equal(state.EngineConfigPGDefaultToastCompression, plan.EngineConfigPGDefaultToastCompression),
		!cmp.Equal(state.EngineConfigPGIdleInTransactionSessionTimeout, plan.EngineConfigPGIdleInTransactionSessionTimeout),
		!cmp.Equal(state.EngineConfigPGJIT, plan.EngineConfigPGJIT),
		!cmp.Equal(state.EngineConfigPGMaxFilesPerProcess, plan.EngineConfigPGMaxFilesPerProcess),
		!cmp.Equal(state.EngineConfigPGMaxLocksPerTransaction, plan.EngineConfigPGMaxLocksPerTransaction),
		!cmp.Equal(state.EngineConfigPGMaxLogicalReplicationWorkers, plan.EngineConfigPGMaxLogicalReplicationWorkers),
		!cmp.Equal(state.EngineConfigPGMaxParallelWorkers, plan.EngineConfigPGMaxParallelWorkers),
		!cmp.Equal(state.EngineConfigPGMaxParallelWorkersPerGather, plan.EngineConfigPGMaxParallelWorkersPerGather),
		!cmp.Equal(state.EngineConfigPGMaxPredLocksPerTransaction, plan.EngineConfigPGMaxPredLocksPerTransaction),
		!cmp.Equal(state.EngineConfigPGMaxReplicationSlots, plan.EngineConfigPGMaxReplicationSlots),
		!cmp.Equal(state.EngineConfigPGMaxSlotWALKeepSize, plan.EngineConfigPGMaxSlotWALKeepSize),
		!cmp.Equal(state.EngineConfigPGMaxStackDepth, plan.EngineConfigPGMaxStackDepth),
		!cmp.Equal(state.EngineConfigPGMaxStandbyArchiveDelay, plan.EngineConfigPGMaxStandbyArchiveDelay),
		!cmp.Equal(state.EngineConfigPGMaxStandbyStreamingDelay, plan.EngineConfigPGMaxStandbyStreamingDelay),
		!cmp.Equal(state.EngineConfigPGMaxWALSenders, plan.EngineConfigPGMaxWALSenders),
		!cmp.Equal(state.EngineConfigPGMaxWorkerProcesses, plan.EngineConfigPGMaxWorkerProcesses),
		!cmp.Equal(state.EngineConfigPGPasswordEncryption, plan.EngineConfigPGPasswordEncryption),
		!cmp.Equal(state.EngineConfigPGPGPartmanBGWInterval, plan.EngineConfigPGPGPartmanBGWInterval),
		!cmp.Equal(state.EngineConfigPGPGPartmanBGWRole, plan.EngineConfigPGPGPartmanBGWRole),
		!cmp.Equal(state.EngineConfigPGPGStatMonitorPGSMEnableQueryPlan, plan.EngineConfigPGPGStatMonitorPGSMEnableQueryPlan),
		!cmp.Equal(state.EngineConfigPGPGStatMonitorPGSMMaxBuckets, plan.EngineConfigPGPGStatMonitorPGSMMaxBuckets),
		!cmp.Equal(state.EngineConfigPGPGStatStatementsTrack, plan.EngineConfigPGPGStatStatementsTrack),
		!cmp.Equal(state.EngineConfigPGTempFileLimit, plan.EngineConfigPGTempFileLimit),
		!cmp.Equal(state.EngineConfigPGTimezone, plan.EngineConfigPGTimezone),
		!cmp.Equal(state.EngineConfigPGTrackActivityQuerySize, plan.EngineConfigPGTrackActivityQuerySize),
		!cmp.Equal(state.EngineConfigPGTrackCommitTimestamp, plan.EngineConfigPGTrackCommitTimestamp),
		!cmp.Equal(state.EngineConfigPGTrackFunctions, plan.EngineConfigPGTrackFunctions),
		!cmp.Equal(state.EngineConfigPGTrackIOTiming, plan.EngineConfigPGTrackIOTiming),
		!cmp.Equal(state.EngineConfigPGWALSenderTimeout, plan.EngineConfigPGWALSenderTimeout),
		!cmp.Equal(state.EngineConfigPGWALWriterDelay, plan.EngineConfigPGWALWriterDelay),
		!cmp.Equal(state.EngineConfigPGStatMonitorEnable, plan.EngineConfigPGStatMonitorEnable),
		!cmp.Equal(state.EngineConfigPGLookoutMaxFailoverReplicationTimeLag, plan.EngineConfigPGLookoutMaxFailoverReplicationTimeLag),
		!cmp.Equal(state.EngineConfigSharedBuffersPercentage, plan.EngineConfigSharedBuffersPercentage),
		!cmp.Equal(state.EngineConfigWorkMem, plan.EngineConfigWorkMem),
	}
	for _, changed := range engineConfigFields {
		if changed {
			shouldUpdate = true

			engineConfig := plan.GetEngineConfig(resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}

			updateOpts.EngineConfig = engineConfig
			if resp.Diagnostics.HasError() {
				return
			}
			break // already handled all updates in one call to GetEngineConfig()
		}
	}

	// `engine_id` field updates
	if !state.EngineID.Equal(plan.EngineID) {
		engine, version, err := helper.ParseDatabaseEngineSlug(plan.EngineID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to parse database engine slug", err.Error())
			return
		}

		if engine != state.Engine.ValueString() {
			resp.Diagnostics.AddError(
				"Cannot update engine component of engine_id",
				fmt.Sprintf("%s != %s", engine, state.Engine.ValueString()),
			)
		}

		shouldUpdate = true
		updateOpts.Version = version
	}

	// `cluster_size` field updates
	if !state.ClusterSize.Equal(plan.ClusterSize) {
		shouldResize = true

		updateOpts.ClusterSize = helper.FrameworkSafeInt64ToInt(plan.ClusterSize.ValueInt64(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if shouldUpdate || shouldResize {
		var updatePoller, resizePoller *linodego.EventPoller
		var err error

		if shouldUpdate {
			updatePoller, err = client.NewEventPoller(ctx, id, linodego.EntityDatabase, linodego.ActionDatabaseUpdate)
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to create update EventPoller for database",
					err.Error(),
				)
				return
			}
		}

		if shouldResize {
			resizePoller, err = client.NewEventPoller(ctx, id, linodego.EntityDatabase, linodego.ActionDatabaseResize)
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to create resize EventPoller for database",
					err.Error(),
				)
				return
			}
		}

		tflog.Debug(ctx, "client.UpdatePostgresDatabase(...)", map[string]any{
			"options": updateOpts,
		})
		if _, err := client.UpdatePostgresDatabase(ctx, id, updateOpts); err != nil {
			resp.Diagnostics.AddError(
				"Failed to update database",
				err.Error(),
			)
			return
		}

		timeoutSeconds := helper.FrameworkSafeFloat64ToInt(updateTimeout.Seconds(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if updatePoller != nil {
			if _, err := updatePoller.WaitForFinished(ctx, timeoutSeconds); err != nil {
				resp.Diagnostics.AddError(
					"Failed to poll for database update event to finish",
					err.Error(),
				)
				return
			}
		}

		if resizePoller != nil {
			if _, err := resizePoller.WaitForFinished(ctx, timeoutSeconds); err != nil {
				resp.Diagnostics.AddError(
					"Failed to poll for database resize event to finish",
					err.Error(),
				)
				return
			}
		}
	}

	resp.Diagnostics.Append(plan.Refresh(ctx, client, id, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.CopyFrom(&state.Model, true)

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	client := r.Meta.Client
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := data.Timeouts.Delete(ctx, DefaultDeleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	ctx = populateLogAttributes(ctx, data)

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeletePostgresDatabase(...)")
	err := client.DeletePostgresDatabase(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); (ok && lerr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				"Failed to delete the database",
				err.Error(),
			)
		}
		return
	}
}

func populateLogAttributes(ctx context.Context, data ResourceModel) context.Context {
	return tflog.SetField(ctx, "id", data.ID)
}
