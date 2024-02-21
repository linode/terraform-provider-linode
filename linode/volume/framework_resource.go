package volume

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
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
				Name:   "linode_volume",
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

func cloneCheck(data *VolumeResourceModel, sourceVolume *linodego.Volume, diags *diag.Diagnostics) {
	if sourceVolume == nil {
		diags.AddError(
			"sourceVolume is nil",
			"sourceVolume is nil. This is a bug in Terraform Provider for Linode",
		)
		return
	}

	if !data.Region.IsNull() && !data.Region.IsUnknown() {
		if data.Region.ValueString() != sourceVolume.Region {
			diags.AddError(
				"Region Mismatched",
				fmt.Sprintf(
					"`region` of source volume differs from specified region: %s != %s",
					data.Region.ValueString(), sourceVolume.Region,
				),
			)
			return
		}
	}

	if !data.Size.IsNull() && !data.Size.IsUnknown() {
		if int64(sourceVolume.Size) > data.Size.ValueInt64() {
			diags.AddError(
				"Insufficient Size",
				fmt.Sprintf(
					"`size` must be greater than or equal to the size of the source volume: %d < %d",
					data.Size.ValueInt64(), sourceVolume.Size,
				),
			)
			return
		}
	}
}

func (r *Resource) CreateVolumeFromSource(
	ctx context.Context, data *VolumeResourceModel, diags *diag.Diagnostics, timeoutSeconds int,
) *linodego.Volume {
	tflog.Debug(ctx, "Create volume from source")

	client := r.Meta.Client
	sourceVolumeID := helper.FrameworkSafeInt64ToInt(data.SourceVolumeID.ValueInt64(), diags)
	if diags.HasError() {
		return nil
	}

	ctx = tflog.SetField(ctx, "source_volume_id", sourceVolumeID)
	tflog.Trace(ctx, "client.GetVolume(...)")

	sourceVolume, err := client.GetVolume(ctx, sourceVolumeID)
	if err != nil {
		diags.AddError("Failed to Get Source Volume", err.Error())
		return nil
	}

	cloneCheck(data, sourceVolume, diags)
	if diags.HasError() {
		return nil
	}

	tflog.Debug(ctx, "client.CloneVolume(...)")

	clonedVolume, err := client.CloneVolume(ctx, sourceVolumeID, data.Label.ValueString())
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Failed to Clone Volume %d", sourceVolumeID),
			err.Error(),
		)
		return clonedVolume
	}

	tflog.Trace(ctx, "client.WaitForVolumeStatus(...)")

	_, err = client.WaitForVolumeStatus(
		ctx, clonedVolume.ID, linodego.VolumeActive, timeoutSeconds,
	)
	if err != nil {
		diags.AddError(
			"Failed to Wait for Volume Cloned to be Ready",
			err.Error(),
		)
		return clonedVolume
	}

	var updateOpts linodego.VolumeUpdateOptions

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		diags.Append(data.Tags.ElementsAs(ctx, &updateOpts.Tags, false)...)
		if diags.HasError() {
			return clonedVolume
		}

		tflog.Debug(ctx, "Update cloned volume", map[string]interface{}{
			"options": updateOpts,
		})

		updatedVolume, err := client.UpdateVolume(ctx, clonedVolume.ID, updateOpts)

		if updatedVolume != nil {
			clonedVolume = updatedVolume
		}

		if err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to Update Cloned Volume %d", clonedVolume.ID),
				err.Error(),
			)
			return clonedVolume
		}
	}

	if !data.Size.IsNull() && !data.Size.IsUnknown() {
		newSize := helper.FrameworkSafeInt64ToInt(data.Size.ValueInt64(), diags)
		if diags.HasError() {
			return clonedVolume
		}
		if newSize != clonedVolume.Size {
			resizedVolume := HandleResize(ctx, client, clonedVolume.ID, newSize, timeoutSeconds, diags)
			if resizedVolume != nil {
				clonedVolume = resizedVolume
			}
			if diags.HasError() {
				return clonedVolume
			}
		}
	}

	if !data.LinodeID.IsNull() && !data.LinodeID.IsUnknown() {
		linodeID := helper.FrameworkSafeInt64ToInt(data.LinodeID.ValueInt64(), diags)
		if diags.HasError() {
			return clonedVolume
		}

		tflog.Debug(ctx, "client.AttachVolume(...)", map[string]interface{}{
			"linode_id": linodeID,
		})

		attachedVolume, err := client.AttachVolume(
			ctx, clonedVolume.ID, &linodego.VolumeAttachOptions{LinodeID: linodeID},
		)

		if attachedVolume != nil {
			clonedVolume = attachedVolume
		}

		if err != nil {
			diags.AddError(
				fmt.Sprintf(
					"Failed to Attach Cloned Volume (%d) to Linode (%d)",
					clonedVolume.ID, linodeID,
				),
				err.Error(),
			)
			return clonedVolume
		}
	}

	return clonedVolume
}

func (r *Resource) CreateVolume(
	ctx context.Context, data *VolumeResourceModel, diags *diag.Diagnostics, timeoutSeconds int,
) *linodego.Volume {
	tflog.Debug(ctx, "Create new volume")

	client := r.Meta.Client

	size := helper.FrameworkSafeInt64ToInt(data.Size.ValueInt64(), diags)

	createOpts := linodego.VolumeCreateOptions{
		Label:  data.Label.ValueString(),
		Region: data.Region.ValueString(),
		Size:   size,
	}

	diags.Append(data.Tags.ElementsAs(ctx, &createOpts.Tags, false)...)
	if diags.HasError() {
		return nil
	}

	if !data.LinodeID.IsNull() && !data.LinodeID.IsUnknown() {
		linodeID := helper.FrameworkSafeInt64ToInt(data.LinodeID.ValueInt64(), diags)
		if diags.HasError() {
			return nil
		}
		createOpts.LinodeID = linodeID
	}

	tflog.Debug(ctx, "client.CreateVolume(...)", map[string]interface{}{
		"options": createOpts,
	})

	volume, err := client.CreateVolume(ctx, createOpts)
	if err != nil {
		diags.AddError("Failed to Create a Volume", err.Error())
		return volume
	}

	if volume != nil {
		tflog.Trace(ctx, "client.WaitForVolumeStatus(...)")

		volume, err = client.WaitForVolumeStatus(
			ctx, volume.ID, linodego.VolumeActive, timeoutSeconds,
		)
		if err != nil {
			diags.AddError(
				"Failed to Wait for Created Volume be Active", err.Error(),
			)
		}
	}
	return volume
}

func (r *Resource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create linode_volume")

	var plan VolumeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, DefaultVolumeCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	timeoutSeconds, err := helper.SafeFloat64ToInt(createTimeout.Seconds())
	if err != nil {
		resp.Diagnostics.AddError("Failed to Convert Timeout to int", err.Error())
	}

	var volume *linodego.Volume

	if !plan.SourceVolumeID.IsNull() {
		volume = r.CreateVolumeFromSource(ctx, &plan, &resp.Diagnostics, timeoutSeconds)
	} else {
		volume = r.CreateVolume(ctx, &plan, &resp.Diagnostics, timeoutSeconds)
	}

	if volume != nil {
		// We should always set the created resource into state even if there is an error
		// to prevent untracked resources created on the cloud
		plan.FlattenVolume(volume, true)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	}
}

func (r *Resource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_volume")

	var state VolumeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	client := r.Meta.Client

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "volume_id", id)

	tflog.Trace(ctx, "client.GetVolume(...)")

	volume, err := client.GetVolume(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Volume Not Found",
				fmt.Sprintf("removing Volume ID %d from state because it no longer exists", id),
			)
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				"Failed to Get the Volume",
				fmt.Sprintf("Error finding the specified Linode Volume: %s", err.Error()),
			)
		}
		return
	}

	resp.Diagnostics.Append(state.FlattenVolume(volume, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func HandleResize(
	ctx context.Context,
	client *linodego.Client,
	volumeID, newSize, timeoutSeconds int,
	diags *diag.Diagnostics,
) *linodego.Volume {
	tflog.Debug(ctx, "Resize volume", map[string]interface{}{
		"volume_id": volumeID,
		"new_size":  newSize,
	})

	tflog.Trace(ctx, "client.ResizeVolume(...)")
	if err := client.ResizeVolume(ctx, volumeID, newSize); err != nil {
		diags.AddError(
			fmt.Sprintf("Failed to Resize Volume %d", volumeID),
			err.Error(),
		)
		return nil
	}

	tflog.Trace(ctx, "client.WaitForVolumeStatus(...)")
	volume, err := client.WaitForVolumeStatus(ctx, volumeID, linodego.VolumeActive, timeoutSeconds)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Failed to Wait for Volume %d Ready from Resizing", volumeID),
			err.Error(),
		)
	}

	return volume
}

func (r *Resource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_volume")

	var plan, state VolumeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := plan.Timeouts.Create(ctx, DefaultVolumeUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	client := r.Meta.Client

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	size := helper.FrameworkSafeInt64ToInt(plan.Size.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "volume_id", id)

	timeoutSeconds, err := helper.SafeFloat64ToInt(updateTimeout.Seconds())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Convert Timeout Seconds %v to int", updateTimeout.Seconds()),
			err.Error(),
		)
		return
	}

	if !state.Size.Equal(plan.Size) {
		volume := HandleResize(ctx, client, id, size, timeoutSeconds, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		id = volume.ID
		resp.Diagnostics.Append(plan.FlattenVolume(volume, true)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	updateOpts := linodego.VolumeUpdateOptions{
		Label: plan.Label.ValueString(),
	}
	doUpdate := false

	if !state.Tags.Equal(plan.Tags) {
		doUpdate = true
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &updateOpts.Tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !state.Label.Equal(plan.Label) {
		doUpdate = true
		updateOpts.Label = plan.Label.ValueString()
	}

	if doUpdate {
		tflog.Debug(ctx, "client.UpdateVolume(...)", map[string]interface{}{
			"options": updateOpts,
		})

		volume, err := client.UpdateVolume(ctx, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to Update Volume %d", id),
				err.Error(),
			)
			return
		}

		resp.Diagnostics.Append(plan.FlattenVolume(volume, true)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// if linode_id change, we should handle attach or detach
	if !state.LinodeID.Equal(plan.LinodeID) {
		if !state.LinodeID.IsNull() {
			// we need to detach if this volume attached to any linode

			volume := DetachVolumeAndWait(ctx, client, id, timeoutSeconds, &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}

			resp.Diagnostics.Append(plan.FlattenVolume(volume, true)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if !plan.LinodeID.IsNull() && !plan.LinodeID.IsUnknown() {
			// we need to attach this volume to another linode if linode_id presents in plan

			linodeID := helper.FrameworkSafeInt64ToInt(plan.LinodeID.ValueInt64(), &resp.Diagnostics)
			if resp.Diagnostics.HasError() {
				return
			}

			attachOptions := linodego.VolumeAttachOptions{
				LinodeID: linodeID,
				ConfigID: 0,
			}

			tflog.Debug(ctx, "client.AttachVolume(...)", map[string]interface{}{
				"options": attachOptions,
			})

			_, err := client.AttachVolume(ctx, id, &attachOptions)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Failed to Attach Volume %d to Linode %d", id, linodeID),
					err.Error(),
				)
				return
			}

			tflog.Trace(ctx, "client.WaitForVolumeLinodeID(...)")

			volume, err := client.WaitForVolumeLinodeID(ctx, id, &linodeID, timeoutSeconds)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Failed to Wait for Volume %d Attached to Linode %d", id, linodeID),
					err.Error(),
				)
				return
			}

			resp.Diagnostics.Append(plan.FlattenVolume(volume, true)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	plan.CopyFrom(state, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func DetachVolumeAndWait(
	ctx context.Context,
	client *linodego.Client,
	id, timeoutSeconds int,
	diags *diag.Diagnostics,
) *linodego.Volume {
	tflog.Debug(ctx, "Detach volume and wait...")

	tflog.Trace(ctx, "client.DetachVolume(...)")

	if err := client.DetachVolume(ctx, id); err != nil {
		diags.AddError(
			fmt.Sprintf("Failed to Detach Volume %d", id),
			err.Error(),
		)
		return nil
	}

	tflog.Trace(ctx, "client.WaitForVolumeLinodeID(...)")

	volume, err := client.WaitForVolumeLinodeID(ctx, id, nil, timeoutSeconds)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Failed to Wait for Volume %d Detached", id),
			err.Error(),
		)
		return nil
	}
	return volume
}

func (r *Resource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_volume")

	var state VolumeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Create(ctx, DefaultVolumeDeleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	client := r.Meta.Client

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "volume_id", id)

	timeoutSeconds, err := helper.SafeFloat64ToInt(deleteTimeout.Seconds())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Convert Timeout Seconds %v to int", deleteTimeout.Seconds()),
			err.Error(),
		)
		return
	}

	if !state.LinodeID.IsNull() {
		DetachVolumeAndWait(ctx, client, id, timeoutSeconds, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Trace(ctx, "client.DeleteVolume(...)")

	err = client.DeleteVolume(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Delete Volume %d", id),
			err.Error(),
		)
	}
}
