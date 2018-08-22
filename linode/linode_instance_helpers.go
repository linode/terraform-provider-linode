package linode

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/linode/linodego"
	"golang.org/x/crypto/sha3"
)

var (
	boolFalse = false
	boolTrue  = true
)

func flattenInstanceSpecs(instance linodego.Instance) []map[string]int {
	return []map[string]int{{
		"vcpus":    instance.Specs.VCPUs,
		"disk":     instance.Specs.Disk,
		"memory":   instance.Specs.Memory,
		"transfer": instance.Specs.Transfer,
	}}
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

func flattenInstanceDisks(instanceDisks []*linodego.InstanceDisk) (disks []map[string]interface{}, swapSize int) {
	for _, disk := range instanceDisks {
		// Determine if swap exists and the size.  If it does not exist, swap_size=0
		if disk.Filesystem == "swap" {
			swapSize += disk.Size
		}
		disks = append(disks, map[string]interface{}{
			"size":       disk.Size,
			"label":      disk.Label,
			"filesystem": string(disk.Filesystem),
			// TODO(displague) these can not be retrieved after the initial send
			// "read_only":       disk.ReadOnly,
			// "image":           disk.Image,
			// "authorized_keys": disk.AuthorizedKeys,
			// "stackscript_id":  disk.StackScriptID,
		})
	}
	return
}

func flattenInstanceConfigs(instanceConfigs []*linodego.InstanceConfig) (configs []map[string]interface{}) {
	for _, config := range instanceConfigs {

		devices := []map[string]interface{}{{
			"sda": flattenInstanceConfigDevice(config.Devices.SDA),
			"sdb": flattenInstanceConfigDevice(config.Devices.SDB),
			"sdc": flattenInstanceConfigDevice(config.Devices.SDC),
			"sdd": flattenInstanceConfigDevice(config.Devices.SDD),
			"sde": flattenInstanceConfigDevice(config.Devices.SDE),
			"sdf": flattenInstanceConfigDevice(config.Devices.SDF),
			"sdg": flattenInstanceConfigDevice(config.Devices.SDG),
			"sdh": flattenInstanceConfigDevice(config.Devices.SDH),
		}}

		// Determine if swap exists and the size.  If it does not exist, swap_size=0
		configs = append(configs, map[string]interface{}{
			"kernel":       config.Kernel,
			"run_level":    string(config.RunLevel),
			"virt_mode":    string(config.VirtMode),
			"root_device":  config.RootDevice,
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
			// panic: interface conversion: interface {} is map[string]map[string]int, not *schema.Set
			"devices": devices,

			// TODO(displague) these can not be retrieved after the initial send
			// "read_only":       disk.ReadOnly,
			// "image":           disk.Image,
			// "authorized_keys": disk.AuthorizedKeys,
			// "stackscript_id":  disk.StackScriptID,
		})
	}
	return
}

func flattenInstanceConfigDevice(dev *linodego.InstanceConfigDevice) []map[string]interface{} {
	if dev == nil {
		return []map[string]interface{}{{
			"disk_id":   0,
			"volume_id": 0,
		}}
	}

	return []map[string]interface{}{{
		"disk_id":   dev.DiskID,
		"volume_id": dev.VolumeID,
	}}
}

func expandInstanceconfigDevice(m map[string]interface{}) *linodego.InstanceConfigDevice {
	var dev *linodego.InstanceConfigDevice
	if m["disk_id"].(int) > 0 || m["volume_id"].(int) > 0 {
		dev = &linodego.InstanceConfigDevice{
			DiskID:   m["disk_id"].(int),
			VolumeID: m["volume_id"].(int),
		}
	}

	return dev
}

// getTotalDiskSize returns the number of disks and their total size.
func getTotalDiskSize(client *linodego.Client, linodeID int) (totalDiskSize int, err error) {
	disks, err := client.ListInstanceDisks(context.Background(), linodeID, nil)
	if err != nil {
		return 0, err
	}

	for _, disk := range disks {
		totalDiskSize += disk.Size
	}

	return totalDiskSize, nil
}

// getBiggestDisk returns the ID and Size of the largest disk attached to the Linode
func getBiggestDisk(client *linodego.Client, linodeID int) (biggestDiskID int, biggestDiskSize int, err error) {
	diskFilter := "{\"+order_by\": \"size\", \"+order\": \"desc\"}"
	disks, err := client.ListInstanceDisks(context.Background(), linodeID, linodego.NewListOptions(1, diskFilter))
	if err != nil {
		return 0, 0, err
	}

	for _, disk := range disks {
		// Find Biggest Disk ID & Size
		if disk.Size > biggestDiskSize {
			biggestDiskID = disk.ID
			biggestDiskSize = disk.Size
		}
	}
	return biggestDiskID, biggestDiskSize, nil
}

// sshKeyState hashes a string passed in as an interface
func sshKeyState(val interface{}) string {
	return hashString(strings.Join(val.([]string), "\n"))
}

// rootPasswordState hashes a string passed in as an interface
func rootPasswordState(val interface{}) string {
	return hashString(val.(string))
}

// hashString hashes a string
func hashString(key string) string {
	hash := sha3.Sum512([]byte(key))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// changeLinodeSize resizes the current linode
func changeLinodeSize(client *linodego.Client, instance *linodego.Instance, d *schema.ResourceData) error {
	typeID, ok := d.Get("type").(string)
	if !ok {
		return fmt.Errorf("Unexpected value for type %v", d.Get("type"))
	}

	targetType, err := client.GetType(context.Background(), typeID)
	if err != nil {
		return fmt.Errorf("Error finding the instance type %s", typeID)
	}

	//biggestDiskID, biggestDiskSize, err := getBiggestDisk(client, instance.ID)

	//currentDiskSize, err := getTotalDiskSize(client, instance.ID)

	if ok, err := client.ResizeInstance(context.Background(), instance.ID, typeID); err != nil || !ok {
		return fmt.Errorf("Error resizing instance %d: %s", instance.ID, err)
	}

	event, err := client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeResize, *instance.Created, int(d.Timeout(schema.TimeoutUpdate).Seconds()))
	if err != nil {
		return fmt.Errorf("Error waiting for instance %d to finish resizing: %s", instance.ID, err)
	}

	if d.Get("disk_expansion").(bool) && instance.Specs.Disk > targetType.Disk {
		// Determine the biggestDisk ID and Size
		biggestDiskID, biggestDiskSize, err := getBiggestDisk(client, instance.ID)
		if err != nil {
			return err
		}
		// Calculate new size, with other disks taken into consideration
		expandedDiskSize := biggestDiskSize + targetType.Disk - instance.Specs.Disk

		// Resize the Disk
		client.ResizeInstanceDisk(context.Background(), instance.ID, biggestDiskID, expandedDiskSize)

		// Wait for the Disk Resize Operation to Complete
		// waitForEventComplete(client, instance.ID, "linode_resize", waitMinutes)
		event, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionDiskResize, *event.Created, int(d.Timeout(schema.TimeoutUpdate).Seconds()))
		if err != nil {
			return fmt.Errorf("Error waiting for resize of Disk %d for Linode %d: %s", biggestDiskID, instance.ID, err)
		}
	}

	// Return the new Linode size
	d.SetPartial("disk_expansion")
	d.SetPartial("type")
	return nil
}

// privateIP determines if an IP is for private use (RFC1918)
// https://stackoverflow.com/a/41273687
func privateIP(ip net.IP) bool {
	private := false
	_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
	_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
	_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
	private = private24BitBlock.Contains(ip) || private20BitBlock.Contains(ip) || private16BitBlock.Contains(ip)
	return private
}

func labelHashcode(v interface{}) int {
	switch t := v.(type) {
	case linodego.InstanceConfig:
		return schema.HashString(t.Label)
	case linodego.InstanceDisk:
		return schema.HashString(t.Label)
	case map[string]interface{}:
		if label, ok := t["label"]; ok {
			return schema.HashString(label.(string))
		}
		panic(fmt.Sprintf("Error hashing label for unknown map: %#v", v))
	default:
		panic(fmt.Sprintf("Error hashing label for unknown interface: %#v", v))
	}
}
