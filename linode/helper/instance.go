package helper

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

const (
	RootPassMinimumCharacters = 11
	RootPassMaximumCharacters = 128
)

var bootEvents = []linodego.EventAction{linodego.ActionLinodeBoot, linodego.ActionLinodeReboot}

// set bootConfig = 0 if using existing boot config
func RebootInstance(ctx context.Context, d *schema.ResourceData, entityID int,
	meta interface{}, bootConfig int,
) diag.Diagnostics {
	ctx = SetLogFieldBulk(ctx, map[string]any{
		"linode_id": entityID,
		"config_id": bootConfig,
	})

	client := meta.(*ProviderMeta).Client
	instance, err := client.GetInstance(ctx, int(entityID))
	if err != nil {
		return diag.Errorf("Error fetching data about the current linode: %s", err)
	}

	if instance.Status != linodego.InstanceRunning {
		tflog.Info(ctx, "Instance is not running")
		return nil
	}

	tflog.Info(ctx, "Rebooting instance")

	p, err := client.NewEventPoller(ctx, entityID, linodego.EntityLinode, linodego.ActionLinodeReboot)
	if err != nil {
		return diag.Errorf("failed to initialize event poller: %s", err)
	}

	err = client.RebootInstance(ctx, instance.ID, bootConfig)

	if err != nil {
		return diag.Errorf("Error rebooting Instance [%d]: %s", instance.ID, err)
	}
	_, err = p.WaitForFinished(ctx, GetDeadlineSeconds(ctx, d))
	if err != nil {
		return diag.Errorf("Error waiting for Instance [%d] to finish rebooting: %s", instance.ID, err)
	}

	if _, err = client.WaitForInstanceStatus(
		ctx, instance.ID, linodego.InstanceRunning, GetDeadlineSeconds(ctx, d),
	); err != nil {
		return diag.Errorf("Timed-out waiting for Linode instance [%d] to boot: %s", instance.ID, err)
	}

	tflog.Debug(ctx, "Instance has finished rebooting")

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

// IsInstanceInBootedState checks whether an instance is in a booted or booting state
func IsInstanceInBootedState(status linodego.InstanceStatus) bool {
	// For diffing purposes, transition states need to be treated as
	// booted == true. This is because these statuses will eventually
	// result in a powered on Linode.
	return status == linodego.InstanceRunning ||
		status == linodego.InstanceRebooting ||
		status == linodego.InstanceBooting
}

// GetCurrentBootedConfig gets the config a linode instance is current booted to
func GetCurrentBootedConfig(ctx context.Context, client *linodego.Client, instID int) (int, error) {
	inst, err := client.GetInstance(ctx, instID)
	if err != nil {
		return 0, err
	}

	// Valid exit condition where no config is booted
	if !IsInstanceInBootedState(inst.Status) {
		return 0, nil
	}

	filter := map[string]any{
		"entity.id":   instID,
		"entity.type": linodego.EntityLinode,
		"+or":         []map[string]any{},
		"+order_by":   "created",
		"+order":      "desc",
	}

	for _, v := range bootEvents {
		filter["+or"] = append(filter["+or"].([]map[string]any), map[string]any{"action": v})
	}

	filterBytes, err := json.Marshal(filter)
	if err != nil {
		return 0, err
	}

	events, err := client.ListEvents(ctx, &linodego.ListOptions{
		Filter: string(filterBytes),
	})
	if err != nil {
		return 0, err
	}

	if len(events) < 1 {
		// This is a valid exit case
		return 0, nil
	}

	return int(events[0].SecondaryEntity.ID.(float64)), nil
}

func CreateRandomRootPassword() (string, error) {
	rawRootPass := make([]byte, 50)
	_, err := rand.Read(rawRootPass)
	if err != nil {
		return "", fmt.Errorf("Failed to generate random password")
	}
	rootPass := base64.StdEncoding.EncodeToString(rawRootPass)
	return rootPass, nil
}
