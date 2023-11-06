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
	"github.com/linode/terraform-provider-linode/linode/helper"
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

func flattenInterfaceIPv4(ipv4 linodego.VPCIPv4) []map[string]any {
	if ipv4.NAT1To1 == "" && ipv4.VPC == "" {
		return nil
	}
	return []map[string]any{
		{
			"vpc":     ipv4.VPC,
			"nat_1_1": ipv4.NAT1To1,
		},
	}
}

func flattenInterfaces(interfaces []linodego.InstanceConfigInterface) []map[string]any {
	result := make([]map[string]any, len(interfaces))

	for i, iface := range interfaces {
		result[i] = map[string]any{
			"purpose":      iface.Purpose,
			"ipam_address": iface.IPAMAddress,
			"label":        iface.Label,
			"id":           iface.ID,
			"vpc_id":       iface.VPCID,
			"subnet_id":    iface.SubnetID,
			"primary":      iface.Primary,
			"active":       iface.Active,
			"ipv4":         flattenInterfaceIPv4(iface.IPv4),
		}
	}

	return result
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
			tflog.Debug(ctx, "Waiting for instance to enter running status")
			if _, err := client.WaitForInstanceStatus(ctx, instance.ID, linodego.InstanceRunning, timeoutSeconds); err != nil {
				return fmt.Errorf("failed to wait for instance running: %s", err)
			}

			p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeReboot)
			if err != nil {
				return fmt.Errorf("failed to poll for events: %s", err)
			}

			tflog.Info(ctx, "Wrong config booted; rebooting into correct config", map[string]any{
				"current_config_id": currentConfig,
			})

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

func expandInterfaceIPv4(ipv4 any) *linodego.VPCIPv4 {
	IPv4 := ipv4.(map[string]any)
	vpcAddress := IPv4["vpc"].(string)
	nat1To1 := IPv4["nat_1_1"].(string)
	if vpcAddress == "" && nat1To1 == "" {
		return nil
	}
	return &linodego.VPCIPv4{
		VPC:     vpcAddress,
		NAT1To1: nat1To1,
	}
}

func expandInterfaces(ctx context.Context, ifaces []any) []linodego.InstanceConfigInterfaceCreateOptions {
	result := make([]linodego.InstanceConfigInterfaceCreateOptions, len(ifaces))

	for i, iface := range ifaces {
		ifaceMap := iface.(map[string]any)

		result[i] = linodego.InstanceConfigInterfaceCreateOptions{
			Purpose: linodego.ConfigInterfacePurpose(ifaceMap["purpose"].(string)),
			Primary: ifaceMap["primary"].(bool),
		}
		if ifaceMap["label"] != nil {
			result[i].Label = ifaceMap["label"].(string)
		}
		if ifaceMap["ipam_address"] != nil {
			result[i].IPAMAddress = ifaceMap["ipam_address"].(string)
		}
		if ifaceMap["subnet_id"] != nil {
			subnet_id := ifaceMap["subnet_id"].(int)
			if subnet_id != 0 {
				result[i].SubnetID = &subnet_id
			}
		}
		tflog.Info(ctx, fmt.Sprintf("%v", ifaceMap["subnet_id"]))
		tflog.Info(ctx, "_____________________________________")
		if ifaceMap["primary"] != nil {
			result[i].Primary = ifaceMap["primary"].(bool)
		}

		if ifaceMap["ipv4"] != nil {
			ipv4 := ifaceMap["ipv4"].([]any)
			if len(ipv4) > 0 {
				result[i].IPv4 = expandInterfaceIPv4(ipv4[0])
			}
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
