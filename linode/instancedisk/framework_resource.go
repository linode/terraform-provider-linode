package instancedisk

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/instance"
)

const (
	DefaultVolumeCreateTimeout = 15 * time.Minute
	DefaultVolumeUpdateTimeout = 20 * time.Minute
	DefaultVolumeDeleteTimeout = 10 * time.Minute
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_instance_disk",
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

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	helper.SetLogFieldBulk(ctx, map[string]any{"linode_id": plan.LinodeID})
	createTimeout, diags := plan.Timeouts.Create(ctx, DefaultVolumeCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	client := r.Meta.Client

	timeoutSeconds := helper.FrameworkSafeFloat64ToInt(createTimeout.Seconds(), &resp.Diagnostics)
	linodeID := helper.FrameworkSafeInt64ToInt(plan.LinodeID.ValueInt64(), &resp.Diagnostics)
	diskSize := helper.FrameworkSafeInt64ToInt(plan.Size.ValueInt64(), &resp.Diagnostics)
	stackScriptID := helper.FrameworkSafeInt64ToInt(
		plan.StackScriptID.ValueInt64(), &resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.InstanceDiskCreateOptions{
		Filesystem:    plan.Filesystem.ValueString(),
		Image:         plan.Image.ValueString(),
		Label:         plan.Label.ValueString(),
		Size:          diskSize,
		StackscriptID: stackScriptID,
	}

	resp.Diagnostics.Append(plan.AuthorizedKeys.ElementsAs(ctx, &createOpts.AuthorizedKeys, false)...)
	resp.Diagnostics.Append(plan.AuthorizedUsers.ElementsAs(ctx, &createOpts.AuthorizedUsers, false)...)
	resp.Diagnostics.Append(plan.StackScriptData.ElementsAs(ctx, &createOpts.StackscriptData, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.RootPass.IsNull() {
		createOpts.RootPass = helper.FrameworkCreateRandomRootPassword(&resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		createOpts.RootPass = plan.RootPass.ValueString()
	}

	p, err := client.NewEventPoller(ctx, linodeID, linodego.EntityLinode, linodego.ActionDiskCreate)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Poll for Events", err.Error())
		return
	}

	tflog.Debug(ctx, "client.CreateInstanceDisk(...)")
	disk, err := client.CreateInstanceDisk(ctx, linodeID, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Create Disk on Linode Instance %d", linodeID),
			err.Error(),
		)
		return
	}

	// Add resource to TF states earlier to prevent
	// dangling resources (resources created but not managed by TF)
	AddDiskResource(ctx, *disk, resp, plan)

	ctx = tflog.SetField(ctx, "disk_id", disk.ID)

	event, err := p.WaitForFinished(ctx, timeoutSeconds)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Wait for the Instance Shutdown Event (%d)", event.ID),
			err.Error(),
		)
	}

	if _, err := client.WaitForInstanceDiskStatus(
		ctx, linodeID, disk.ID, linodego.DiskReady, timeoutSeconds,
	); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Wait for Disk (%d) to be Ready", disk.ID), err.Error(),
		)
	}

	// get latest status of the disk
	tflog.Trace(ctx, "client.GetInstanceDisk(...)")
	disk, err = client.GetInstanceDisk(ctx, linodeID, disk.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Get Disk %d of Linode Instance %d", disk.ID, linodeID),
			err.Error(),
		)
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(strconv.Itoa(disk.ID))

	plan.FlattenDisk(disk, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	linodeID, id := getLinodeIDAndDiskID(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	tflog.Trace(ctx, "client.GetInstanceDisk(...)")
	disk, err := client.GetInstanceDisk(ctx, linodeID, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Disk Not Found",
				fmt.Sprintf(
					"%s\nRemoving Disk (%d) of Linode instance (%d) from state because it no longer exists",
					err.Error(), id, linodeID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Find the Linode Instance Disk (%d)", id),
			err.Error(),
		)
		return
	}

	state.FlattenDisk(disk, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	ctx = populateLogAttributes(ctx, state)

	updateTimeout, diags := plan.Timeouts.Update(ctx, DefaultVolumeUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	linodeID, id := getLinodeIDAndDiskID(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	client := r.Meta.Client
	timeoutSeconds := helper.FrameworkSafeFloat64ToInt(updateTimeout.Seconds(), &resp.Diagnostics)
	size := helper.FrameworkSafeInt64ToInt(plan.Size.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Size.Equal(plan.Size) {
		if err := handleDiskResize(
			ctx, client, linodeID, id, size, timeoutSeconds,
		); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Resize Disk %d", id), err.Error(),
			)
			return
		}
	}

	updateOpts := linodego.InstanceDiskUpdateOptions{}
	shouldUpdate := false

	if !state.Label.Equal(plan.Label) {
		updateOpts.Label = plan.Label.ValueString()
		shouldUpdate = true
	}

	if shouldUpdate {
		tflog.Debug(ctx, "client.UpdateInstanceDisk(...)", map[string]any{
			"options": updateOpts,
		})
		if _, err := client.UpdateInstanceDisk(ctx, linodeID, id, updateOpts); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Update Disk %d", id), err.Error(),
			)
			return
		}
		tflog.Debug(ctx, "Disk update event finished")
	}

	// get latest status of the disk
	tflog.Trace(ctx, "client.GetInstanceDisk(...)")
	disk, err := client.GetInstanceDisk(ctx, linodeID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Find the Disk (%d)", id),
			err.Error(),
		)
	}

	plan.FlattenDisk(disk, true)

	plan.CopyFrom(state, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Meta.Client

	linodeID, id := getLinodeIDAndDiskID(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, DefaultVolumeDeleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutSeconds := helper.FrameworkSafeFloat64ToInt(
		deleteTimeout.Seconds(), &resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	configID, err := helper.GetCurrentBootedConfig(ctx, client, linodeID)
	if err != nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Failed to Get the Current Booted Config of Linode %d", linodeID),
			fmt.Sprintf(
				"Will attempt to delete disk without without rebooting the instance. Error: %s",
				err.Error(),
			),
		)
	}

	shouldShutdown := configID != 0
	diskInConfig, err := diskInConfig(ctx, client, id, linodeID, configID)
	if err != nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf(
				"Failed to Check If Disk (%d) is Used in the Booted Config (%d) of Linode Instance (%d)",
				id, configID, linodeID,
			),
			fmt.Sprintf(
				"%s\nWill attempt to delete disk without without rebooting the instance.",
				err.Error(),
			),
		)
	}

	// Shutdown instance if active
	if shouldShutdown {
		if r.Meta.Config.SkipInstanceReadyPoll.ValueBool() {
			resp.Diagnostics.AddError(
				"Linode Instance Shutdown is Required for this Disk Deletion",
				"Please consider set please consider setting 'skip_implicit_reboots' "+
					"to true in the Linode provider config.",
			)
			return
		}
		if err := instance.SafeShutdownInstance(
			ctx, client, linodeID, timeoutSeconds,
		); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Shutdown Linode Instance %d", linodeID),
				err.Error(),
			)
		}
	}

	tflog.Info(ctx, "Deleting instance disk")
	p, err := client.NewEventPollerWithSecondary(
		ctx,
		linodeID,
		linodego.EntityLinode,
		id,
		linodego.ActionDiskDelete,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Initialize Event Poller", err.Error())
		return
	}

	tflog.Debug(ctx, "client.DeleteInstanceDisk(...)")
	if err := client.DeleteInstanceDisk(ctx, linodeID, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Delete Linode Instance Disk %d", id), err.Error(),
		)
		return
	}

	if _, err := p.WaitForFinished(ctx, timeoutSeconds); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Wait for Linode Instance Disk Deletion Finished", err.Error(),
		)
		return
	}

	// Reboot the instance if necessary
	if shouldShutdown && !diskInConfig {
		if err := instance.BootInstanceSync(
			ctx, client, linodeID, configID, timeoutSeconds,
		); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Boot Instance %d", linodeID), err.Error(),
			)
		}
	}
}

func diskInConfig(
	ctx context.Context, client *linodego.Client, diskID, linodeID, configID int,
) (bool, error) {
	if configID == 0 {
		return false, nil
	}

	tflog.Trace(ctx, "client.GetInstanceConfig(...)")
	cfg, err := client.GetInstanceConfig(ctx, linodeID, configID)
	if err != nil {
		return false, err
	}

	if cfg.Devices == nil {
		return false, nil
	}

	reflectMap := reflect.ValueOf(*cfg.Devices)

	for i := 0; i < reflectMap.NumField(); i++ {
		field := reflectMap.Field(i).Interface().(*linodego.InstanceConfigDevice)
		if field == nil {
			continue
		}

		if field.DiskID == diskID {
			return true, nil
		}
	}

	return false, nil
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import "+r.Config.Name)
	helper.ImportStateWithMultipleIDs(
		ctx, req, resp,
		[]helper.ImportableID{
			{
				Name:          "linode_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "id",
				TypeConverter: helper.IDTypeConverterString,
			},
		},
	)
}

func populateLogAttributes(ctx context.Context, model ResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"linode_id": model.LinodeID,
		"disk_id":   model.ID,
	})
}
