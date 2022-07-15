package helper

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

// set bootConfig = 0 if using existing boot config
func RebootInstance(ctx context.Context, d *schema.ResourceData, entityID int,
	meta interface{}, bootConfig int,
) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client
	instance, err := client.GetInstance(ctx, int(entityID))
	if err != nil {
		return diag.Errorf("Error fetching data about the current linode: %s", err)
	}

	if instance.Status != linodego.InstanceRunning {
		log.Printf("[INFO] Instance [%d] is not running\n", instance.ID)
		return nil
	}

	log.Printf("[INFO] Instance [%d] will be rebooted\n", instance.ID)

	err = client.RebootInstance(ctx, instance.ID, bootConfig)

	if err != nil {
		return diag.Errorf("Error rebooting Instance [%d]: %s", instance.ID, err)
	}
	_, err = client.WaitForEventFinished(ctx, entityID, linodego.EntityLinode,
		linodego.ActionLinodeReboot, *instance.Created, GetDeadlineSeconds(ctx, d))
	if err != nil {
		return diag.Errorf("Error waiting for Instance [%d] to finish rebooting: %s", instance.ID, err)
	}
	if _, err = client.WaitForInstanceStatus(
		ctx, instance.ID, linodego.InstanceRunning, GetDeadlineSeconds(ctx, d),
	); err != nil {
		return diag.Errorf("Timed-out waiting for Linode instance [%d] to boot: %s", instance.ID, err)
	}
	return nil
}

// GetDeadlineSeconds gets the seconds remaining until deadline is met.
func GetDeadlineSeconds(ctx context.Context, d *schema.ResourceData) int {
	duration := d.Timeout(schema.TimeoutUpdate)
	if deadline, ok := ctx.Deadline(); ok {
		duration = time.Until(deadline)
	}
	return int(duration.Seconds())
}
