package instance

import (
	"context"
	"fmt"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func flattenInstance(
	ctx context.Context, client *linodego.Client, instance *linodego.Instance,
) (map[string]any, error) {
	result := make(map[string]any)

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
	result["maintenance_policy"] = instance.MaintenancePolicy
	result["watchdog_enabled"] = instance.WatchdogEnabled
	result["group"] = instance.Group
	result["tags"] = instance.Tags
	result["capabilities"] = instance.Capabilities
	result["locks"] = instance.Locks
	result["image"] = instance.Image
	result["interface_generation"] = instance.InterfaceGeneration
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

func flattenInstanceBackups(instance linodego.Instance) []map[string]any {
	return []map[string]any{{
		"available": instance.Backups.Available,
		"enabled":   instance.Backups.Enabled,
		"schedule": []map[string]any{{
			"day":    instance.Backups.Schedule.Day,
			"window": instance.Backups.Schedule.Window,
		}},
	}}
}

func flattenInstanceDisks(instanceDisks []linodego.InstanceDisk) (disks []map[string]any, swapSize int) {
	for _, disk := range instanceDisks {
		// Determine if swap exists and the size.  If it does not exist, swap_size=0
		if disk.Filesystem == "swap" {
			swapSize += disk.Size
		}
		disks = append(disks, map[string]any{
			"id":         disk.ID,
			"size":       disk.Size,
			"label":      disk.Label,
			"filesystem": string(disk.Filesystem),
		})
	}
	return disks, swapSize
}

func flattenInstanceConfigDevice(
	dev *linodego.InstanceConfigDevice, diskLabelIDMap map[int]string,
) []map[string]any {
	if dev == nil || emptyInstanceConfigDevice(*dev) {
		return nil
	}

	if dev.DiskID > 0 {
		ret := map[string]any{
			"disk_id": dev.DiskID,
		}
		if label, found := diskLabelIDMap[dev.DiskID]; found {
			ret["disk_label"] = label
		}
		return []map[string]any{ret}
	}
	return []map[string]any{{
		"volume_id": dev.VolumeID,
	}}
}

func flattenInstanceConfigs(
	instanceConfigs []linodego.InstanceConfig, diskLabelIDMap map[int]string,
) (configs []map[string]any) {
	for _, config := range instanceConfigs {

		devices := []map[string]any{{
			"sda":  flattenInstanceConfigDevice(config.Devices.SDA, diskLabelIDMap),
			"sdb":  flattenInstanceConfigDevice(config.Devices.SDB, diskLabelIDMap),
			"sdc":  flattenInstanceConfigDevice(config.Devices.SDC, diskLabelIDMap),
			"sdd":  flattenInstanceConfigDevice(config.Devices.SDD, diskLabelIDMap),
			"sde":  flattenInstanceConfigDevice(config.Devices.SDE, diskLabelIDMap),
			"sdf":  flattenInstanceConfigDevice(config.Devices.SDF, diskLabelIDMap),
			"sdg":  flattenInstanceConfigDevice(config.Devices.SDG, diskLabelIDMap),
			"sdh":  flattenInstanceConfigDevice(config.Devices.SDH, diskLabelIDMap),
			"sdi":  flattenInstanceConfigDevice(config.Devices.SDI, diskLabelIDMap),
			"sdj":  flattenInstanceConfigDevice(config.Devices.SDJ, diskLabelIDMap),
			"sdk":  flattenInstanceConfigDevice(config.Devices.SDK, diskLabelIDMap),
			"sdl":  flattenInstanceConfigDevice(config.Devices.SDL, diskLabelIDMap),
			"sdm":  flattenInstanceConfigDevice(config.Devices.SDM, diskLabelIDMap),
			"sdn":  flattenInstanceConfigDevice(config.Devices.SDN, diskLabelIDMap),
			"sdo":  flattenInstanceConfigDevice(config.Devices.SDO, diskLabelIDMap),
			"sdp":  flattenInstanceConfigDevice(config.Devices.SDP, diskLabelIDMap),
			"sdq":  flattenInstanceConfigDevice(config.Devices.SDQ, diskLabelIDMap),
			"sdr":  flattenInstanceConfigDevice(config.Devices.SDR, diskLabelIDMap),
			"sds":  flattenInstanceConfigDevice(config.Devices.SDS, diskLabelIDMap),
			"sdt":  flattenInstanceConfigDevice(config.Devices.SDT, diskLabelIDMap),
			"sdu":  flattenInstanceConfigDevice(config.Devices.SDU, diskLabelIDMap),
			"sdv":  flattenInstanceConfigDevice(config.Devices.SDV, diskLabelIDMap),
			"sdw":  flattenInstanceConfigDevice(config.Devices.SDW, diskLabelIDMap),
			"sdx":  flattenInstanceConfigDevice(config.Devices.SDX, diskLabelIDMap),
			"sdy":  flattenInstanceConfigDevice(config.Devices.SDY, diskLabelIDMap),
			"sdz":  flattenInstanceConfigDevice(config.Devices.SDZ, diskLabelIDMap),
			"sdaa": flattenInstanceConfigDevice(config.Devices.SDAA, diskLabelIDMap),
			"sdab": flattenInstanceConfigDevice(config.Devices.SDAB, diskLabelIDMap),
			"sdac": flattenInstanceConfigDevice(config.Devices.SDAC, diskLabelIDMap),
			"sdad": flattenInstanceConfigDevice(config.Devices.SDAD, diskLabelIDMap),
			"sdae": flattenInstanceConfigDevice(config.Devices.SDAE, diskLabelIDMap),
			"sdaf": flattenInstanceConfigDevice(config.Devices.SDAF, diskLabelIDMap),
			"sdag": flattenInstanceConfigDevice(config.Devices.SDAG, diskLabelIDMap),
			"sdah": flattenInstanceConfigDevice(config.Devices.SDAH, diskLabelIDMap),
			"sdai": flattenInstanceConfigDevice(config.Devices.SDAI, diskLabelIDMap),
			"sdaj": flattenInstanceConfigDevice(config.Devices.SDAJ, diskLabelIDMap),
			"sdak": flattenInstanceConfigDevice(config.Devices.SDAK, diskLabelIDMap),
			"sdal": flattenInstanceConfigDevice(config.Devices.SDAL, diskLabelIDMap),
			"sdam": flattenInstanceConfigDevice(config.Devices.SDAM, diskLabelIDMap),
			"sdan": flattenInstanceConfigDevice(config.Devices.SDAN, diskLabelIDMap),
			"sdao": flattenInstanceConfigDevice(config.Devices.SDAO, diskLabelIDMap),
			"sdap": flattenInstanceConfigDevice(config.Devices.SDAP, diskLabelIDMap),
			"sdaq": flattenInstanceConfigDevice(config.Devices.SDAQ, diskLabelIDMap),
			"sdar": flattenInstanceConfigDevice(config.Devices.SDAR, diskLabelIDMap),
			"sdas": flattenInstanceConfigDevice(config.Devices.SDAS, diskLabelIDMap),
			"sdat": flattenInstanceConfigDevice(config.Devices.SDAT, diskLabelIDMap),
			"sdau": flattenInstanceConfigDevice(config.Devices.SDAU, diskLabelIDMap),
			"sdav": flattenInstanceConfigDevice(config.Devices.SDAV, diskLabelIDMap),
			"sdaw": flattenInstanceConfigDevice(config.Devices.SDAW, diskLabelIDMap),
			"sdax": flattenInstanceConfigDevice(config.Devices.SDAX, diskLabelIDMap),
			"sday": flattenInstanceConfigDevice(config.Devices.SDAY, diskLabelIDMap),
			"sdaz": flattenInstanceConfigDevice(config.Devices.SDAZ, diskLabelIDMap),
			"sdba": flattenInstanceConfigDevice(config.Devices.SDBA, diskLabelIDMap),
			"sdbb": flattenInstanceConfigDevice(config.Devices.SDBB, diskLabelIDMap),
			"sdbc": flattenInstanceConfigDevice(config.Devices.SDBC, diskLabelIDMap),
			"sdbd": flattenInstanceConfigDevice(config.Devices.SDBD, diskLabelIDMap),
			"sdbe": flattenInstanceConfigDevice(config.Devices.SDBE, diskLabelIDMap),
			"sdbf": flattenInstanceConfigDevice(config.Devices.SDBF, diskLabelIDMap),
			"sdbg": flattenInstanceConfigDevice(config.Devices.SDBG, diskLabelIDMap),
			"sdbh": flattenInstanceConfigDevice(config.Devices.SDBH, diskLabelIDMap),
			"sdbi": flattenInstanceConfigDevice(config.Devices.SDBI, diskLabelIDMap),
			"sdbj": flattenInstanceConfigDevice(config.Devices.SDBJ, diskLabelIDMap),
			"sdbk": flattenInstanceConfigDevice(config.Devices.SDBK, diskLabelIDMap),
			"sdbl": flattenInstanceConfigDevice(config.Devices.SDBL, diskLabelIDMap),
		}}

		interfaces := helper.FlattenInterfaces(config.Interfaces)

		// Determine if swap exists and the size.  If it does not exist, swap_size=0
		c := map[string]any{
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
	return configs
}

func flattenInstanceSpecs(instance linodego.Instance) []map[string]int {
	return []map[string]int{{
		"vcpus":               instance.Specs.VCPUs,
		"disk":                instance.Specs.Disk,
		"memory":              instance.Specs.Memory,
		"transfer":            instance.Specs.Transfer,
		"accelerated_devices": instance.Specs.AcceleratedDevices,
		"gpus":                instance.Specs.GPUs,
	}}
}

func flattenInstanceSimple(instance *linodego.Instance) (map[string]any, error) {
	result := make(map[string]any)

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
	result["maintenance_policy"] = instance.MaintenancePolicy
	result["watchdog_enabled"] = instance.WatchdogEnabled
	result["group"] = instance.Group
	result["tags"] = instance.Tags
	result["capabilities"] = instance.Capabilities
	result["image"] = instance.Image
	result["interface_generation"] = instance.InterfaceGeneration
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
