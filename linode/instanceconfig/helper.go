package instanceconfig

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func flattenDeviceMap(deviceMap linodego.InstanceConfigDeviceMap) []map[string]any {
	result := make(map[string]any)

	reflectMap := reflect.ValueOf(deviceMap)

	for i := 0; i < reflectMap.NumField(); i++ {
		field := reflectMap.Field(i).Interface().(*linodego.InstanceConfigDevice)
		if field == nil {
			continue
		}

		fieldName := strings.ToLower(reflectMap.Type().Field(i).Name)

		result[fieldName] = []map[string]any{
			{
				"disk_id":   field.DiskID,
				"volume_id": field.VolumeID,
			},
		}
	}

	return []map[string]any{result}
}

func flattenHelpers(helpers linodego.InstanceConfigHelpers) []map[string]any {
	result := make(map[string]any)

	result["devtmpfs_automount"] = helpers.DevTmpFsAutomount
	result["distro"] = helpers.Distro
	result["modules_dep"] = helpers.ModulesDep
	result["network"] = helpers.Network
	result["updatedb_disabled"] = helpers.UpdateDBDisabled

	return []map[string]any{result}
}

func flattenInterfaces(interfaces []linodego.InstanceConfigInterface) []map[string]any {
	result := make([]map[string]any, len(interfaces))

	for i, iface := range interfaces {
		// Workaround for "222" responses for null IPAM
		// addresses from the API.
		// TODO: Remove this when issue is resolved.
		if iface.IPAMAddress == "222" {
			iface.IPAMAddress = ""
		}

		result[i] = map[string]any{
			"purpose":      iface.Purpose,
			"ipam_address": iface.IPAMAddress,
			"label":        iface.Label,
		}
	}

	return result
}

func expandDeviceMap(deviceMap any) *linodego.InstanceConfigDeviceMap {
	var result linodego.InstanceConfigDeviceMap

	deviceMapSlice := deviceMap.([]any)

	if len(deviceMapSlice) < 1 {
		return nil
	}

	devices := deviceMapSlice[0].(map[string]any)

	for k, v := range devices {
		currentDeviceSlice := v.([]any)
		if len(currentDeviceSlice) < 1 {
			continue
		}

		currentDevice := currentDeviceSlice[0].(map[string]any)

		device := linodego.InstanceConfigDevice{}

		if diskID, ok := currentDevice["disk_id"]; ok {
			device.DiskID = diskID.(int)
		}

		if volumeID, ok := currentDevice["volume_id"]; ok {
			device.VolumeID = volumeID.(int)
		}

		// Get the corresponding struct field and set it to the correct device
		field := reflect.Indirect(reflect.ValueOf(&result)).FieldByName(strings.ToUpper(k))
		field.Set(reflect.ValueOf(&device))
	}

	return &result
}

func expandHelpers(helpersRaw any) *linodego.InstanceConfigHelpers {
	helpersSlice := helpersRaw.([]any)

	if len(helpersSlice) < 1 {
		return nil
	}

	helpers := helpersSlice[0].(map[string]any)

	return &linodego.InstanceConfigHelpers{
		UpdateDBDisabled:  helpers["updatedb_disabled"].(bool),
		Distro:            helpers["distro"].(bool),
		ModulesDep:        helpers["modules_dep"].(bool),
		Network:           helpers["network"].(bool),
		DevTmpFsAutomount: helpers["devtmpfs_automount"].(bool),
	}
}

func applyBootStatus(ctx context.Context, client *linodego.Client, instance *linodego.Instance, configID int,
	timeoutSeconds int, booted bool,
) error {
	isBooted := helper.IsInstanceInBootedState(instance.Status)
	currentConfig, err := helper.GetCurrentBootedConfig(ctx, client, instance.ID)
	if err != nil {
		return fmt.Errorf("failed to get current booted config id: %s", err)
	}

	bootedTrue := func() error {
		// Instance is already in desired state
		if currentConfig == configID && isBooted {
			return nil
		}

		// Instance is booted into the wrong config
		if isBooted && currentConfig != configID {
			if _, err := client.WaitForInstanceStatus(ctx, instance.ID, linodego.InstanceRunning, timeoutSeconds); err != nil {
				return fmt.Errorf("failed to wait for instance running: %s", err)
			}

			p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeReboot)
			if err != nil {
				return fmt.Errorf("failed to poll for events: %s", err)
			}

			if err := client.RebootInstance(ctx, instance.ID, configID); err != nil {
				return fmt.Errorf("failed to reboot instance %d: %s", instance.ID, err)
			}

			if _, err := p.WaitForFinished(ctx, timeoutSeconds); err != nil {
				return fmt.Errorf("failed to wait for instance reboot: %s", err)
			}

			return nil
		}

		// Boot the instance
		if !isBooted {
			p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeBoot)
			if err != nil {
				return fmt.Errorf("failed to poll for events: %s", err)
			}

			if err := client.BootInstance(ctx, instance.ID, configID); err != nil {
				return fmt.Errorf("failed to boot instance %d %d: %s", instance.ID, configID, err)
			}

			if _, err := p.WaitForFinished(ctx, timeoutSeconds); err != nil {
				return fmt.Errorf("failed to wait for instance boot: %s", err)
			}
		}

		return nil
	}

	bootedFalse := func() error {
		// Instance is already in desired state
		if !isBooted || currentConfig != configID {
			return nil
		}

		if _, err := client.WaitForInstanceStatus(ctx, instance.ID, linodego.InstanceRunning, timeoutSeconds); err != nil {
			return fmt.Errorf("failed to wait for instance running: %s", err)
		}

		p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeShutdown)
		if err != nil {
			return fmt.Errorf("failed to poll for events: %s", err)
		}

		if err := client.ShutdownInstance(ctx, instance.ID); err != nil {
			return fmt.Errorf("failed to shutdown instance: %s", err)
		}

		if _, err := p.WaitForFinished(ctx, timeoutSeconds); err != nil {
			return fmt.Errorf("failed to wait for instance shutdown: %s", err)
		}

		return nil
	}

	if booted {
		err = bootedTrue()
	} else {
		err = bootedFalse()
	}

	return err
}

func expandInterfaces(ifaces []any) []linodego.InstanceConfigInterface {
	result := make([]linodego.InstanceConfigInterface, len(ifaces))

	for i, iface := range ifaces {
		ifaceMap := iface.(map[string]any)

		result[i] = linodego.InstanceConfigInterface{
			IPAMAddress: ifaceMap["ipam_address"].(string),
			Label:       ifaceMap["label"].(string),
			Purpose:     linodego.ConfigInterfacePurpose(ifaceMap["purpose"].(string)),
		}
	}

	return result
}

func isConfigBooted(ctx context.Context, client *linodego.Client,
	instance *linodego.Instance, configID int,
) (bool, error) {
	currentConfig, err := helper.GetCurrentBootedConfig(ctx, client, instance.ID)
	if err != nil {
		return false, fmt.Errorf("failed to get current booted config id: %s", err)
	}

	return helper.IsInstanceInBootedState(instance.Status) && currentConfig == configID, nil
}
