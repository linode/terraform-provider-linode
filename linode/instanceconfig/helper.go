package instanceconfig

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func getDeviceMapFields(deviceMap linodego.InstanceConfigDeviceMap) [][2]any {
	result := make([][2]any, 0)

	reflectMap := reflect.ValueOf(deviceMap)

	for i := 0; i < reflectMap.NumField(); i++ {
		field := reflectMap.Field(i).Interface().(*linodego.InstanceConfigDevice)
		if field == nil {
			continue
		}

		fieldName := strings.ToLower(reflectMap.Type().Field(i).Name)

		result = append(result, [2]any{fieldName, field})
	}

	return result
}

func flattenDeviceMapToBlock(deviceMap linodego.InstanceConfigDeviceMap) []map[string]any {
	result := make([]map[string]any, 0)

	for _, pair := range getDeviceMapFields(deviceMap) {
		fieldName := pair[0].(string)
		field := pair[1].(*linodego.InstanceConfigDevice)

		result = append(
			result,
			map[string]any{
				"device_name": fieldName,
				"disk_id":     field.DiskID,
				"volume_id":   field.VolumeID,
			},
		)
	}

	return result
}

func flattenDeviceMapToNamedBlock(deviceMap linodego.InstanceConfigDeviceMap) []map[string]any {
	result := make(map[string]any)

	for _, pair := range getDeviceMapFields(deviceMap) {
		fieldName := pair[0].(string)
		field := pair[1].(*linodego.InstanceConfigDevice)

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

func createDevice(deviceMap map[string]any) linodego.InstanceConfigDevice {
	device := linodego.InstanceConfigDevice{}

	if diskID, ok := deviceMap["disk_id"]; ok {
		device.DiskID = diskID.(int)
	}

	if volumeID, ok := deviceMap["volume_id"]; ok {
		device.VolumeID = volumeID.(int)
	}

	return device
}

func expandDevicesBlock(devicesBlock any) *linodego.InstanceConfigDeviceMap {
	var result linodego.InstanceConfigDeviceMap

	devices := devicesBlock.(*schema.Set).List()

	if len(devices) <= 0 {
		return nil
	}

	seenDevices := make(map[string]bool)

	for _, rawDevice := range devices {
		device := rawDevice.(map[string]any)
		linodeGoDevice := createDevice(device)

		if deviceName, ok := device["device_name"]; ok {
			if seenDevices[deviceName.(string)] {
				log.Printf("[WARN] device %v was defined more than once", deviceName)
			} else {
				seenDevices[deviceName.(string)] = true
			}

			field := reflect.Indirect(
				reflect.ValueOf(&result),
			).FieldByName(
				strings.ToUpper(deviceName.(string)),
			)

			field.Set(reflect.ValueOf(&linodeGoDevice))
		}
	}

	return &result
}

func expandDevicesNamedBlock(devicesNamedBlock any) *linodego.InstanceConfigDeviceMap {
	var result linodego.InstanceConfigDeviceMap

	deviceMapSlice := devicesNamedBlock.([]any)

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
		linodeGoDevice := createDevice(currentDevice)

		// Get the corresponding struct field and set it to the correct device
		field := reflect.Indirect(reflect.ValueOf(&result)).FieldByName(strings.ToUpper(k))
		field.Set(reflect.ValueOf(&linodeGoDevice))
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

func applyBootStatus(ctx context.Context, client *linodego.Client, linodeID int, configID int,
	timeoutSeconds int, booted bool, reboot bool,
) error {
	instance, err := client.GetInstance(ctx, linodeID)
	if err != nil {
		return fmt.Errorf("Error finding the specified Linode Instance: %s", err)
	}

	isBooted := helper.IsInstanceInBootedState(instance.Status)
	currentConfig, err := helper.GetCurrentBootedConfig(ctx, client, instance.ID)
	if err != nil {
		return fmt.Errorf("failed to get current booted config id: %s", err)
	}

	bootedTrue := func() error {
		// Instance is already in desired state
		if currentConfig == configID && isBooted && !reboot {
			return nil
		}

		// Instance is booted into the wrong config or the booted config requires reboot
		if isBooted && (currentConfig != configID || reboot) {
			tflog.Debug(ctx, "Waiting for instance to enter running status")
			if _, err := client.WaitForInstanceStatus(ctx, instance.ID, linodego.InstanceRunning, timeoutSeconds); err != nil {
				return fmt.Errorf("failed to wait for instance running: %s", err)
			}

			p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeReboot)
			if err != nil {
				return fmt.Errorf("failed to poll for events: %s", err)
			}

			if currentConfig != configID {
				tflog.Info(ctx, "Wrong config booted; rebooting into correct config", map[string]any{
					"current_config_id": currentConfig,
				})
			}
			if reboot {
				tflog.Info(ctx, "Current config updated; rebooting to adopt changes", map[string]any{
					"current_config_id": currentConfig,
				})
			}

			if err := client.RebootInstance(ctx, instance.ID, configID); err != nil {
				return fmt.Errorf("failed to reboot instance %d: %s", instance.ID, err)
			}

			tflog.Debug(ctx, "Instance reboot triggered, waiting for event finished")

			event, err := p.WaitForFinished(ctx, timeoutSeconds)
			if err != nil {
				return fmt.Errorf("failed to wait for instance reboot: %s", err)
			}

			tflog.Debug(ctx, "Instance reboot finished", map[string]any{
				"event_id": event.ID,
			})

			return nil
		}

		// Boot the instance
		if !isBooted {
			tflog.Info(ctx, "Instance is not booted; booting into config")

			p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeBoot)
			if err != nil {
				return fmt.Errorf("failed to poll for events: %s", err)
			}

			if err := client.BootInstance(ctx, instance.ID, configID); err != nil {
				return fmt.Errorf("failed to boot instance %d %d: %s", instance.ID, configID, err)
			}

			tflog.Debug(ctx, "Instance boot triggered, waiting for event finished")

			event, err := p.WaitForFinished(ctx, timeoutSeconds)
			if err != nil {
				return fmt.Errorf("failed to wait for instance boot: %s", err)
			}

			tflog.Debug(ctx, "Instance boot finished", map[string]any{
				"event_id": event.ID,
			})
		}

		return nil
	}

	bootedFalse := func() error {
		// Instance is already in desired state
		if !isBooted || currentConfig != configID {
			return nil
		}

		tflog.Info(ctx, "Handling instance shutdown")

		if _, err := client.WaitForInstanceStatus(ctx, instance.ID, linodego.InstanceRunning, timeoutSeconds); err != nil {
			return fmt.Errorf("failed to wait for instance running: %s", err)
		}

		p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeShutdown)
		if err != nil {
			return fmt.Errorf("failed to poll for events: %s", err)
		}

		tflog.Debug(ctx, "Instance has entered running state, shutting down")

		if err := client.ShutdownInstance(ctx, instance.ID); err != nil {
			return fmt.Errorf("failed to shutdown instance: %s", err)
		}

		event, err := p.WaitForFinished(ctx, timeoutSeconds)
		if err != nil {
			return fmt.Errorf("failed to wait for instance shutdown: %s", err)
		}

		tflog.Debug(ctx, "Instance shutdown complete", map[string]any{
			"event_id": event.ID,
		})

		return nil
	}

	if booted {
		err = bootedTrue()
	} else {
		err = bootedFalse()
	}

	return err
}

func isConfigBooted(
	ctx context.Context,
	client *linodego.Client,
	instance *linodego.Instance,
	configID int,
	previousValue bool,
) (bool, error) {
	if !helper.IsInstanceInBootedState(instance.Status) {
		return false, nil
	}

	currentConfig, err := helper.GetCurrentBootedConfig(ctx, client, instance.ID)
	if err != nil {
		return false, fmt.Errorf("failed to get current booted config id: %s", err)
	}

	// If the last booted config couldn't be resolved,
	// we assume the event either expired (older than 90 days)
	// or the instance was transferred across accounts.
	//
	// In this case, we want to preserve the existing booted value from state.
	if currentConfig == 0 {
		tflog.Info(
			ctx,
			"booted: The instance is online but a boot event could not be resolved. "+
				"The previous value from state will be used.",
			map[string]any{
				"assumed_value": previousValue,
			},
		)
		return previousValue, nil
	}

	return currentConfig == configID, nil
}
