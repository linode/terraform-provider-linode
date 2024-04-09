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
	ctx context.Context, client *linodego.Client, instID, diskID, newSize, timeoutSeconds int,
) error {
	configID, err := helper.GetCurrentBootedConfig(ctx, client, instID)
	if err != nil {
		return err
	}

	shouldShutdown := configID != 0

	if shouldShutdown {
		tflog.Info(ctx, "Shutting down Instance for disk resize")

		p, err := client.NewEventPoller(ctx, instID, linodego.EntityLinode, linodego.ActionLinodeShutdown)
		if err != nil {
			return fmt.Errorf("failed to poll for events: %s", err)
		}

		tflog.Debug(ctx, "client.ShutdownInstance(...)")

		if err := client.ShutdownInstance(ctx, instID); err != nil {
			return fmt.Errorf("failed to shutdown instance: %s", err)
		}

		tflog.Debug(ctx, "Waiting for Instance shutdown operation to complete")

		if _, err := p.WaitForFinished(ctx, timeoutSeconds); err != nil {
			return fmt.Errorf("failed to wait for instance shutdown: %s", err)
		}

		tflog.Debug(ctx, "Instance finished shutting down")
	}

	disk, err := client.GetInstanceDisk(ctx, instID, diskID)
	if err != nil {
		return fmt.Errorf("failed to get instance disk: %s", err)
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
		return fmt.Errorf("failed to poll for events: %s", err)
	}

	tflog.Debug(ctx, "client.ResizeInstanceDisk(...)", map[string]any{
		"new_size": newSize,
	})
	if err := client.ResizeInstanceDisk(ctx, instID, diskID, newSize); err != nil {
		return fmt.Errorf("failed to resize disk: %s", err)
	}

	// Wait for the resize event to complete
	if _, err := p.WaitForFinished(ctx, timeoutSeconds); err != nil {
		return fmt.Errorf("failed to wait for disk resize: %s", err)
	}

	// Check to see if the resize operation worked
	if updatedDisk, err := client.WaitForInstanceDiskStatus(ctx, instID, disk.ID, linodego.DiskReady,
		timeoutSeconds); err != nil {
		return fmt.Errorf("failed to wait for disk ready: %s", err)
	} else if updatedDisk.Size != newSize {
		return fmt.Errorf(
			"failed to resize disk %d from %d to %d", disk.ID, disk.Size, newSize)
	}

	tflog.Debug(ctx, "Resize operation complete")

	// Reboot the instance if necessary
	if shouldShutdown {
		tflog.Info(ctx, "Rebooting instance to previously booted config")

		p, err := client.NewEventPoller(ctx, instID, linodego.EntityLinode, linodego.ActionLinodeBoot)
		if err != nil {
			return fmt.Errorf("failed to poll for events: %s", err)
		}

		tflog.Debug(ctx, "client.BootInstance(...)", map[string]any{
			"config_id": configID,
		})
		if err := client.BootInstance(ctx, instID, configID); err != nil {
			return fmt.Errorf("failed to boot instance %d %d: %s", instID, configID, err)
		}

		if _, err := p.WaitForFinished(ctx, timeoutSeconds); err != nil {
			return fmt.Errorf("failed to wait for instance boot: %s", err)
		}

		tflog.Debug(ctx, "Reboot event finished")
	}

	return nil
}
