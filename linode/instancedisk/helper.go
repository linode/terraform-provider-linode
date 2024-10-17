package instancedisk

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func AddDiskResource(ctx context.Context, diskID int, resp *resource.CreateResponse, plan ResourceModel) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(strconv.Itoa(diskID)))
	resp.State.SetAttribute(ctx, path.Root("linode_id"), plan.LinodeID)
}

func getLinodeIDAndDiskID(data ResourceModel, diags *diag.Diagnostics) (int, int) {
	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), diags)
	linodeID := helper.FrameworkSafeInt64ToInt(data.LinodeID.ValueInt64(), diags)
	return linodeID, id
}

// resizeDiskSync resizes the instance with the given ID to the given size,
// shutting down/booting the instance if necessary.
func resizeDiskSync(
	ctx context.Context,
	client *linodego.Client,
	meta *helper.FrameworkProviderMeta,
	linodeID,
	diskID,
	newSize,
	timeoutSeconds int,
) diag.Diagnostics {
	return runDiskOperation(
		ctx, client, meta, linodeID, timeoutSeconds,
		func() (resultDiag diag.Diagnostics) {
			disk, err := client.GetInstanceDisk(ctx, linodeID, diskID)
			if err != nil {
				resultDiag.AddError("Failed to get Linode instance disk", err.Error())
				return
			}

			tflog.Info(ctx, "Resizing instance disk", map[string]any{
				"old_size": disk.Size,
				"new_size": newSize,
			})

			p, err := client.NewEventPollerWithSecondary(
				ctx,
				linodeID,
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
			if err := client.ResizeInstanceDisk(ctx, linodeID, diskID, newSize); err != nil {
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
				linodeID,
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
			return
		},
	)
}

func runDiskOperation(
	ctx context.Context,
	client *linodego.Client,
	meta *helper.FrameworkProviderMeta,
	linodeID int,
	timeoutSeconds int,
	callOperation func() diag.Diagnostics,
) (resultDiag diag.Diagnostics) {
	originalStatus, err := helper.WaitForInstanceNonTransientStatus(
		ctx, client, linodeID, timeoutSeconds,
	)
	if err != nil {
		resultDiag.AddError("Failed to wait for instance to exit transient status", err.Error())
		return
	}

	if originalStatus == linodego.InstanceRunning {
		if meta.Config.SkipImplicitReboots.ValueBool() {
			resultDiag.AddError(
				"Linode instance shutdown is required to complete this operation",
				"Please consider setting 'skip_implicit_reboots' "+
					"to true in the Linode provider config.",
			)
			return
		}

		if err := helper.ShutDownInstanceSync(ctx, client, linodeID, timeoutSeconds); err != nil {
			resultDiag.AddError("Failed to shutdown Linode instance", err.Error())
			return
		}
	}

	// Run the operation
	resultDiag.Append(callOperation()...)
	if resultDiag.HasError() {
		return
	}

	// Boot the instance back up if necessary
	if originalStatus == linodego.InstanceRunning {
		// NOTE: A config ID of 0 will boot the instance into its previously booted config
		if err := helper.BootInstanceSync(ctx, client, linodeID, 0, timeoutSeconds); err != nil {
			resultDiag.AddError("Failed to boot Linode instance", err.Error())
			return
		}
	}

	return nil
}
