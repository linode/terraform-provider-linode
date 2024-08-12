package instance

import (
	"context"
	"fmt"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func flattenInstance(
	ctx context.Context, client *linodego.Client, instance *linodego.Instance,
) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	id := instance.ID

	instanceNetwork, err := client.GetInstanceIPAddresses(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ips for linode instance %d: %s", id, err)
	}

	var ips []string
	for _, ip := range instance.IPv4 {
		ips = append(ips, ip.String())
	}

	result["ipv4"] = ips
	result["ipv6"] = instance.IPv6

	public, private := instanceNetwork.IPv4.Public, instanceNetwork.IPv4.Private

	if len(public) > 0 {
		result["ip_address"] = public[0].Address
	}

	if len(private) > 0 {
		result["private_ip_address"] = private[0].Address
	}

	result["id"] = instance.ID
	result["label"] = instance.Label
	result["status"] = instance.Status
	result["type"] = instance.Type
	result["region"] = instance.Region
	result["watchdog_enabled"] = instance.WatchdogEnabled
	result["group"] = instance.Group
	result["tags"] = instance.Tags
	result["image"] = instance.Image
	result["host_uuid"] = instance.HostUUID
	result["has_user_data"] = instance.HasUserData
	result["disk_encryption"] = instance.DiskEncryption
	result["lke_cluster_id"] = instance.LKEClusterID

	result["backups"] = flattenInstanceBackups(*instance)
	result["specs"] = flattenInstanceSpecs(*instance)
	result["alerts"] = flattenInstanceAlerts(*instance)
	result["placement_group"] = flattenInstancePlacementGroup(*instance)

	instanceDisks, err := client.ListInstanceDisks(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the disks for the Linode instance %d: %s", id, err)
	}

	disks, swapSize := flattenInstanceDisks(instanceDisks)
	result["disk"] = disks
	result["swap_size"] = swapSize

	instanceConfigs, err := client.ListInstanceConfigs(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the config for Linode instance %d (%s): %s", id, instance.Label, err)
	}

	diskLabelIDMap := make(map[int]string, len(instanceDisks))
	for _, disk := range instanceDisks {
		diskLabelIDMap[disk.ID] = disk.Label
	}

	configs := flattenInstanceConfigs(instanceConfigs, diskLabelIDMap)

	result["config"] = configs
	if len(instanceConfigs) == 1 {
		result["boot_config_label"] = instanceConfigs[0].Label
	}

	return result, nil
}

func flattenInstanceAlerts(instance linodego.Instance) []map[string]int {
	return []map[string]int{{
		"cpu":            instance.Alerts.CPU,
		"io":             instance.Alerts.IO,
		"network_in":     instance.Alerts.NetworkIn,
		"network_out":    instance.Alerts.NetworkOut,
		"transfer_quota": instance.Alerts.TransferQuota,
	}}
}

func flattenInstanceBackups(instance linodego.Instance) []map[string]interface{} {
	return []map[string]interface{}{{
		"available": instance.Backups.Available,
		"enabled":   instance.Backups.Enabled,
		"schedule": []map[string]interface{}{{
			"day":    instance.Backups.Schedule.Day,
			"window": instance.Backups.Schedule.Window,
		}},
	}}
}

func flattenInstanceDisks(instanceDisks []linodego.InstanceDisk) (disks []map[string]interface{}, swapSize int) {
	for _, disk := range instanceDisks {
		// Determine if swap exists and the size.  If it does not exist, swap_size=0
		if disk.Filesystem == "swap" {
			swapSize += disk.Size
		}
		disks = append(disks, map[string]interface{}{
			"id":         disk.ID,
			"size":       disk.Size,
			"label":      disk.Label,
			"filesystem": string(disk.Filesystem),
		})
	}
	return
}

func flattenInstanceConfigDevice(
	dev *linodego.InstanceConfigDevice, diskLabelIDMap map[int]string,
) []map[string]interface{} {
	if dev == nil || emptyInstanceConfigDevice(*dev) {
		return nil
	}

	if dev.DiskID > 0 {
		ret := map[string]interface{}{
			"disk_id": dev.DiskID,
		}
		if label, found := diskLabelIDMap[dev.DiskID]; found {
			ret["disk_label"] = label
		}
		return []map[string]interface{}{ret}
	}
	return []map[string]interface{}{{
		"volume_id": dev.VolumeID,
	}}
}

func flattenInstanceConfigs(
	instanceConfigs []linodego.InstanceConfig, diskLabelIDMap map[int]string,
) (configs []map[string]interface{}) {
	for _, config := range instanceConfigs {

		devices := []map[string]interface{}{{
			"sda": flattenInstanceConfigDevice(config.Devices.SDA, diskLabelIDMap),
			"sdb": flattenInstanceConfigDevice(config.Devices.SDB, diskLabelIDMap),
			"sdc": flattenInstanceConfigDevice(config.Devices.SDC, diskLabelIDMap),
			"sdd": flattenInstanceConfigDevice(config.Devices.SDD, diskLabelIDMap),
			"sde": flattenInstanceConfigDevice(config.Devices.SDE, diskLabelIDMap),
			"sdf": flattenInstanceConfigDevice(config.Devices.SDF, diskLabelIDMap),
			"sdg": flattenInstanceConfigDevice(config.Devices.SDG, diskLabelIDMap),
			"sdh": flattenInstanceConfigDevice(config.Devices.SDH, diskLabelIDMap),
		}}

		interfaces := helper.FlattenInterfaces(config.Interfaces)

		// Determine if swap exists and the size.  If it does not exist, swap_size=0
		c := map[string]interface{}{
			"id":           config.ID,
			"root_device":  config.RootDevice,
			"kernel":       config.Kernel,
			"run_level":    string(config.RunLevel),
			"virt_mode":    string(config.VirtMode),
			"comments":     config.Comments,
			"memory_limit": config.MemoryLimit,
			"label":        config.Label,
			"helpers": []map[string]bool{{
				"updatedb_disabled":  config.Helpers.UpdateDBDisabled,
				"distro":             config.Helpers.Distro,
				"modules_dep":        config.Helpers.ModulesDep,
				"network":            config.Helpers.Network,
				"devtmpfs_automount": config.Helpers.DevTmpFsAutomount,
			}},
			"devices":   devices,
			"interface": interfaces,
		}

		configs = append(configs, c)
	}
	return
}

func flattenInstanceSpecs(instance linodego.Instance) []map[string]int {
	return []map[string]int{{
		"vcpus":    instance.Specs.VCPUs,
		"disk":     instance.Specs.Disk,
		"memory":   instance.Specs.Memory,
		"transfer": instance.Specs.Transfer,
	}}
}

func flattenInstanceSimple(instance *linodego.Instance) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var ips []string
	for _, ip := range instance.IPv4 {
		ips = append(ips, ip.String())
	}

	result["id"] = instance.ID
	result["ipv4"] = ips
	result["ipv6"] = instance.IPv6
	result["label"] = instance.Label
	result["status"] = instance.Status
	result["type"] = instance.Type
	result["region"] = instance.Region
	result["watchdog_enabled"] = instance.WatchdogEnabled
	result["group"] = instance.Group
	result["tags"] = instance.Tags
	result["image"] = instance.Image
	result["host_uuid"] = instance.HostUUID
	result["backups"] = flattenInstanceBackups(*instance)
	result["specs"] = flattenInstanceSpecs(*instance)
	result["alerts"] = flattenInstanceAlerts(*instance)

	return result, nil
}

func flattenInstancePlacementGroup(instance linodego.Instance) []map[string]any {
	if instance.PlacementGroup == nil {
		return nil
	}

	result := map[string]any{
		"id":                     instance.PlacementGroup.ID,
		"label":                  instance.PlacementGroup.Label,
		"placement_group_type":   instance.PlacementGroup.PlacementGroupType,
		"placement_group_policy": instance.PlacementGroup.PlacementGroupPolicy,
	}

	return []map[string]any{result}
}
