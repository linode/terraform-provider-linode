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
	originalStatus, err := helper.WaitForInstanceNonTransientStatus(
		ctx, client, instID, timeoutSeconds,
	)
	if err != nil {
		return fmt.Errorf("failed to wait for instance to exit transient status: %w", err)
	}

	if originalStatus == linodego.InstanceRunning {
		if err := helper.ShutDownInstanceSync(ctx, client, instID, timeoutSeconds); err != nil {
			return fmt.Errorf("failed to shut down instance: %w", err)
		}
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
	if updatedDisk, err := client.WaitForInstanceDiskStatus(
		ctx,
		instID,
		disk.ID,
		linodego.DiskReady,
		timeoutSeconds,
	); err != nil {
		return fmt.Errorf("failed to wait for disk ready: %s", err)
	} else if updatedDisk.Size != newSize {
		return fmt.Errorf(
			"failed to resize disk %d from %d to %d", disk.ID, disk.Size, newSize,
		)
	}

	tflog.Debug(ctx, "Resize operation complete")

	// Boot the instance back up if necessary
	if originalStatus == linodego.InstanceRunning {
		// NOTE: A config ID of 0 will boot the instance into its previously booted config
		if err := helper.BootInstanceSync(ctx, client, instID, 0, timeoutSeconds); err != nil {
			return fmt.Errorf("failed to boot instance: %w", err)
		}
	}

	return nil
}
