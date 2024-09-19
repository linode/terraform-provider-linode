package instancedisk

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func AddDiskResource(ctx context.Context, disk linodego.InstanceDisk, resp *resource.CreateResponse, plan ResourceModel) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(strconv.Itoa(disk.ID)))
	resp.State.SetAttribute(ctx, path.Root("linode_id"), plan.LinodeID)
}

func getLinodeIDAndDiskID(data ResourceModel, diags *diag.Diagnostics) (int, int) {
	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), diags)
	linodeID := helper.FrameworkSafeInt64ToInt(data.LinodeID.ValueInt64(), diags)
	return linodeID, id
}

func handleDiskResize(
	ctx context.Context,
	client *linodego.Client,
	meta *helper.FrameworkProviderMeta,
	instID,
	diskID,
	newSize,
	timeoutSeconds int,
) (resultDiag diag.Diagnostics) {
	originalStatus, err := helper.WaitForInstanceNonTransientStatus(
		ctx, client, instID, timeoutSeconds,
	)
	if err != nil {
		resultDiag.AddError("Failed to wait for instance to exit transient status", err.Error())
		return
	}

	if originalStatus == linodego.InstanceRunning {
		if meta.Config.SkipImplicitReboots.ValueBool() {
			resultDiag.AddError(
				"Linode instance shutdown is required to delete this disk.",
				"Please consider setting 'skip_implicit_reboots' "+
					"to true in the Linode provider config.",
			)
			return
		}

		if err := helper.ShutDownInstanceSync(ctx, client, instID, timeoutSeconds); err != nil {
			resultDiag.AddError("Failed to shutdown Linode instance", err.Error())
			return
		}
	}

	disk, err := client.GetInstanceDisk(ctx, instID, diskID)
	if err != nil {
		resultDiag.AddError("Failed to get Linode instance disk", err.Error())
		return
	}

	tflog.Info(ctx, "Resizing Instance disk", map[string]any{
		"old_size": disk.Size,
		"new_size": newSize,
	})

	p, err := client.NewEventPollerWithSecondary(
		ctx,
		instID,
		linodego.EntityLinode,
		diskID,
		linodego.ActionDiskResize)
	if err != nil {
		resultDiag.AddError("Failed to create event poller", err.Error())
		return
	}

	tflog.Debug(ctx, "client.ResizeInstanceDisk(...)", map[string]any{
		"new_size": newSize,
	})
	if err := client.ResizeInstanceDisk(ctx, instID, diskID, newSize); err != nil {
		resultDiag.AddError("Failed to resize Linode instance disk", err.Error())
		return
	}

	// Wait for the resize event to complete
	if _, err := p.WaitForFinished(ctx, timeoutSeconds); err != nil {
		resultDiag.AddError(
			"Failed to wait for Linode instance disk resize to complete",
			err.Error(),
		)
		return
	}

	// Check to see if the resize operation worked
	if updatedDisk, err := client.WaitForInstanceDiskStatus(
		ctx,
		instID,
		disk.ID,
		linodego.DiskReady,
		timeoutSeconds,
	); err != nil {
		resultDiag.AddError(
			"Failed to wait for Linode instance disk to be ready",
			err.Error(),
		)
		return
	} else if updatedDisk.Size != newSize {
		resultDiag.AddError(
			"Failed to resize Linode instance disk",
			fmt.Sprintf("Disk %d has size %d, expected %d", disk.ID, disk.Size, newSize),
		)
		return
	}

	tflog.Debug(ctx, "Resize operation complete")

	// Boot the instance back up if necessary
	if originalStatus == linodego.InstanceRunning {
		// NOTE: A config ID of 0 will boot the instance into its previously booted config
		if err := helper.BootInstanceSync(ctx, client, instID, 0, timeoutSeconds); err != nil {
			resultDiag.AddError("Failed to boot Linode instance", err.Error())
			return
		}
	}

	return nil
}
