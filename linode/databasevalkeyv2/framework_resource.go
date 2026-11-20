package databasevalkeyv2

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/databaseshared"
)

const (
	DefaultCreateTimeout = 60 * time.Minute
	DefaultUpdateTimeout = 60 * time.Minute
	DefaultDeleteTimeout = 5 * time.Minute
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(helper.BaseResourceConfig{
			Name:        "linode_database_valkey_v2",
			IDType:      types.StringType,
			Schema:      &frameworkResourceSchema,
			TimeoutOpts: &timeouts.Opts{Create: true, Update: true, Delete: true},
		}),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create "+r.Config.Name)
	var data ResourceModel
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

	privateNetwork := data.GetPrivateNetwork(ctx, &resp.Diagnostics)
	createOpts := linodego.ValkeyCreateOptions{
		Label:        data.Label.ValueString(),
		Region:       data.Region.ValueString(),
		Type:         data.Type.ValueString(),
		Engine:       data.EngineID.ValueString(),
		ClusterSize:  helper.FrameworkSafeInt64PointerToIntPointer(data.ClusterSize.ValueInt64Pointer(), &resp.Diagnostics),
		AllowList:    data.GetAllowList(ctx, &resp.Diagnostics),
		Fork:         data.GetFork(&resp.Diagnostics),
		EngineConfig: data.GetEngineConfig(&resp.Diagnostics),
	}
	if privateNetwork != nil {
		createOpts.PrivateNetwork = privateNetwork.ToLinodego(resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	createPoller, err := client.NewEventPollerWithoutEntity(linodego.EntityDatabase, linodego.ActionDatabaseCreate)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create event poller", err.Error())
		return
	}

	tflog.Debug(ctx, "client.CreateValkeyDatabase(...)", map[string]any{"options": createOpts})
	db, err := client.CreateValkeyDatabase(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Valkey database", err.Error())
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), strconv.Itoa(db.ID))
	ctx = tflog.SetField(ctx, "id", db.ID)
	createPoller.EntityID = db.ID
	if _, err := createPoller.WaitForFinished(ctx); err != nil {
		resp.Diagnostics.AddError("Failed to wait for Valkey database creation", err.Error())
		return
	}

	if err := client.WaitForDatabaseStatus(ctx, db.ID, linodego.DatabaseEngineTypeValkey, linodego.DatabaseStatusActive); err != nil {
		resp.Diagnostics.AddError("Failed to wait for Valkey database to become active", err.Error())
		return
	}

	updates := data.GetUpdates(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if updates != nil {
		updateOpts := linodego.ValkeyUpdateOptions{Updates: updates.ToLinodego(resp.Diagnostics)}
		if resp.Diagnostics.HasError() {
			return
		}
		if _, err := client.UpdateValkeyDatabase(ctx, db.ID, updateOpts); err != nil {
			resp.Diagnostics.AddError("Failed to update Valkey database maintenance window", err.Error())
			return
		}
		if err := client.WaitForDatabaseStatus(ctx, db.ID, linodego.DatabaseEngineTypeValkey, linodego.DatabaseStatusActive); err != nil {
			resp.Diagnostics.AddError("Failed to wait for Valkey database to become active", err.Error())
			return
		}
	}

	if err := databaseshared.ReconcileSuspensionSync(
		ctx,
		client,
		db.ID,
		linodego.DatabaseEngineTypeValkey,
		false,
		data.Suspended.ValueBool(),
		createTimeout,
	); err != nil {
		resp.Diagnostics.AddError("Failed to reconcile Valkey database suspension", err.Error())
		return
	}

	resp.Diagnostics.Append(data.Refresh(ctx, client, db.ID, true)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = types.StringValue(strconv.Itoa(db.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read "+r.Config.Name)
	var data ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	db, err := r.Meta.Client.GetValkeyDatabase(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning("Valkey database no longer exists", fmt.Sprintf("Removing Valkey database with ID %v from state", id))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to refresh the Valkey database", err.Error())
		return
	}

	resp.Diagnostics.Append(data.Refresh(ctx, r.Meta.Client, db.ID, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update "+r.Config.Name)
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
	client := r.Meta.Client
	if err := databaseshared.ReconcileSuspensionSync(
		ctx,
		client,
		id,
		linodego.DatabaseEngineTypeValkey,
		state.Suspended.ValueBool(),
		plan.Suspended.ValueBool(),
		updateTimeout,
	); err != nil {
		resp.Diagnostics.AddError("Failed to reconcile Valkey database suspension", err.Error())
		return
	}

	var updateOpts linodego.ValkeyUpdateOptions
	shouldUpdate := false
	shouldResize := false
	if !state.Label.Equal(plan.Label) {
		updateOpts.Label = plan.Label.ValueStringPointer()
		shouldUpdate = true
	}
	if !state.AllowList.Equal(plan.AllowList) {
		updateOpts.AllowList = plan.GetAllowList(ctx, &resp.Diagnostics)
		shouldUpdate = true
	}
	if !state.Type.Equal(plan.Type) {
		updateOpts.Type = plan.Type.ValueStringPointer()
		shouldResize = true
	}
	if !state.ClusterSize.Equal(plan.ClusterSize) {
		updateOpts.ClusterSize = helper.FrameworkSafeInt64PointerToIntPointer(plan.ClusterSize.ValueInt64Pointer(), &resp.Diagnostics)
		shouldResize = true
	}
	if !state.PrivateNetwork.Equal(plan.PrivateNetwork) {
		privateNetwork := plan.GetPrivateNetwork(ctx, &resp.Diagnostics)
		if privateNetwork == nil {
			updateOpts.PrivateNetwork = linodego.DoublePointerNull[linodego.DatabasePrivateNetwork]()
		} else {
			updateOpts.PrivateNetwork = linodego.Pointer(privateNetwork.ToLinodego(resp.Diagnostics))
		}
		shouldUpdate = true
	}
	if !state.Updates.Equal(plan.Updates) {
		updates := plan.GetUpdates(ctx, &resp.Diagnostics)
		if updates != nil {
			updateOpts.Updates = updates.ToLinodego(resp.Diagnostics)
		}
		shouldUpdate = true
	}

	configChanged := []bool{
		!helper.FrameworkValuesShallowEqual(state.EngineConfigBackupHour, plan.EngineConfigBackupHour),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigBackupMinute, plan.EngineConfigBackupMinute),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigFrequentSnapshots, plan.EngineConfigFrequentSnapshots),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyACLChannelsDefault, plan.EngineConfigValkeyACLChannelsDefault),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyActiveExpireEffort, plan.EngineConfigValkeyActiveExpireEffort),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyActiveDefrag, plan.EngineConfigValkeyActiveDefrag),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyLFUDecayTime, plan.EngineConfigValkeyLFUDecayTime),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyLFULogFactor, plan.EngineConfigValkeyLFULogFactor),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyMaxmemoryPolicy, plan.EngineConfigValkeyMaxmemoryPolicy),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyNumberOfDatabases, plan.EngineConfigValkeyNumberOfDatabases),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyPersistence, plan.EngineConfigValkeyPersistence),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyPubsubClientOutputBufferLimit, plan.EngineConfigValkeyPubsubClientOutputBufferLimit),
		!helper.FrameworkValuesShallowEqual(state.EngineConfigValkeyTimeout, plan.EngineConfigValkeyTimeout),
	}
	if slices.Contains(configChanged, true) {
		updateOpts.EngineConfig = plan.GetEngineConfig(&resp.Diagnostics)
		shouldUpdate = true
	}

	if !state.EngineID.Equal(plan.EngineID) {
		engine, version, err := databaseshared.ParseDatabaseEngineSlug(plan.EngineID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to parse Valkey engine ID", err.Error())
			return
		}
		if engine != state.Engine.ValueString() {
			resp.Diagnostics.AddError("Cannot change the engine component of engine_id", fmt.Sprintf("%s != %s", engine, state.Engine.ValueString()))
			return
		}
		updateOpts.Version = linodego.Pointer(version)
		shouldUpdate = true
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var updatePoller, resizePoller *linodego.EventPoller
	var err error
	if shouldUpdate {
		updatePoller, err = client.NewEventPoller(ctx, id, linodego.EntityDatabase, linodego.ActionDatabaseUpdate)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create Valkey update event poller", err.Error())
			return
		}
	}
	if shouldResize {
		resizePoller, err = client.NewEventPoller(ctx, id, linodego.EntityDatabase, linodego.ActionDatabaseResize)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create Valkey resize event poller", err.Error())
			return
		}
	}
	if shouldUpdate || shouldResize {
		if _, err := client.UpdateValkeyDatabase(ctx, id, updateOpts); err != nil {
			resp.Diagnostics.AddError("Failed to update Valkey database", err.Error())
			return
		}
		if updatePoller != nil {
			if _, err := updatePoller.WaitForFinished(ctx); err != nil {
				resp.Diagnostics.AddError("Failed to wait for Valkey update event", err.Error())
				return
			}
		}
		if resizePoller != nil {
			if _, err := resizePoller.WaitForFinished(ctx); err != nil {
				resp.Diagnostics.AddError("Failed to wait for Valkey resize event", err.Error())
				return
			}
		}
		if err := client.WaitForDatabaseStatus(ctx, id, linodego.DatabaseEngineTypeValkey, linodego.DatabaseStatusActive); err != nil {
			resp.Diagnostics.AddError("Failed to wait for Valkey database to become active", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(plan.Refresh(ctx, client, id, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.CopyFrom(&state.Model, true)
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.Meta.Client.DeleteValkeyDatabase(ctx, id); err != nil {
		if lerr, ok := err.(*linodego.Error); (ok && lerr.Code != 404) || !ok {
			resp.Diagnostics.AddError("Failed to delete Valkey database", err.Error())
		}
	}
}
