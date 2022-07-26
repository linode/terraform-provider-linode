package instancedisk

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func expandStackScriptData(data any) map[string]string {
	dataMap := data.(map[string]any)
	result := make(map[string]string, len(dataMap))

	for k, v := range dataMap {
		result[k] = v.(string)
	}

	return result
}

func handleDiskResize(ctx context.Context, client linodego.Client, instID, diskID, newSize, timeoutSeconds int) error {
	configID, err := helper.GetCurrentBootedConfig(ctx, &client, instID)
	if err != nil {
		return err
	}

	shouldShutdown := configID != 0

	if shouldShutdown {
		log.Printf("[INFO] Shutting down instance %d for disk %d resize", instID, diskID)

		if err := client.ShutdownInstance(ctx, instID); err != nil {
			return fmt.Errorf("failed to shutdown instance: %s", err)
		}

		if _, err := client.WaitForEventFinished(ctx, instID, linodego.EntityLinode,
			linodego.ActionLinodeShutdown, time.Now(), timeoutSeconds); err != nil {
			return fmt.Errorf("failed to wait for instance shutdown: %s", err)
		}
	}

	disk, err := client.GetInstanceDisk(ctx, instID, diskID)
	if err != nil {
		return fmt.Errorf("failed to get instance disk: %s", err)
	}

	if err := client.ResizeInstanceDisk(ctx, instID, diskID, newSize); err != nil {
		return fmt.Errorf("failed to resize disk: %s", err)
	}

	// Wait for the resize event to complete
	_, err = client.WaitForEventFinished(ctx, instID, linodego.EntityLinode, linodego.ActionDiskResize,
		time.Now(), timeoutSeconds)
	if err != nil {
		return fmt.Errorf("failed to resize disk: %s", err)
	}

	// Check to see if the resize operation worked
	if updatedDisk, err := client.WaitForInstanceDiskStatus(ctx, instID, disk.ID, linodego.DiskReady,
		timeoutSeconds); err != nil {
		return fmt.Errorf("failed to wait for disk ready: %s", err)
	} else if updatedDisk.Size != newSize {
		return fmt.Errorf(
			"failed to resize disk %d from %d to %d", disk.ID, disk.Size, newSize)
	}

	// Reboot the instance if necessary
	if shouldShutdown {
		log.Printf("[INFO] Booting instance %d to config %d", instID, configID)
		if err := client.BootInstance(ctx, instID, configID); err != nil {
			return fmt.Errorf("failed to boot instance %d %d: %s", instID, configID, err)
		}

		if _, err := client.WaitForEventFinished(ctx, instID, linodego.EntityLinode,
			linodego.ActionLinodeBoot, time.Now(), timeoutSeconds); err != nil {
			return fmt.Errorf("failed to wait for instance boot: %s", err)
		}
	}

	return nil
}
