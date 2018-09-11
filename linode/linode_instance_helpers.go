package linode

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

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

func flattenInstanceDisks(instanceDisks []linodego.InstanceDisk) (disks []map[string]interface{}, swapSize int) {
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

func flattenInstanceConfigs(instanceConfigs []linodego.InstanceConfig) (configs []map[string]interface{}) {
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

// TODO(displague) do we need a disk_label map?
func expandInstanceConfigDeviceMap(m map[string]interface{}, diskIDLabelMap map[string]int) (deviceMap *linodego.InstanceConfigDeviceMap, err error) {
	if len(m) == 0 {
		return nil, nil
	}
	deviceMap = &linodego.InstanceConfigDeviceMap{}
	for k, rdev := range m {
		devSlots := rdev.([]interface{})
		for _, rrdev := range devSlots {
			dev := rrdev.(map[string]interface{})
			if k == "sda" {
				deviceMap.SDA = &linodego.InstanceConfigDevice{}
				if err := assignConfigDevice(deviceMap.SDA, dev, diskIDLabelMap); err != nil {
					return nil, err
				}
			}
			if k == "sdb" {
				deviceMap.SDB = &linodego.InstanceConfigDevice{}
				if err := assignConfigDevice(deviceMap.SDB, dev, diskIDLabelMap); err != nil {
					return nil, err
				}
			}
			if k == "sdc" {
				deviceMap.SDC = &linodego.InstanceConfigDevice{}
				if err := assignConfigDevice(deviceMap.SDC, dev, diskIDLabelMap); err != nil {
					return nil, err
				}
			}
			if k == "sdd" {
				deviceMap.SDD = &linodego.InstanceConfigDevice{}
				if err := assignConfigDevice(deviceMap.SDD, dev, diskIDLabelMap); err != nil {
					return nil, err
				}
			}
			if k == "sde" {
				deviceMap.SDE = &linodego.InstanceConfigDevice{}
				if err := assignConfigDevice(deviceMap.SDE, dev, diskIDLabelMap); err != nil {
					return nil, err
				}
			}
			if k == "sdf" {
				deviceMap.SDF = &linodego.InstanceConfigDevice{}

				if err := assignConfigDevice(deviceMap.SDF, dev, diskIDLabelMap); err != nil {
					return nil, err
				}
			}
			if k == "sdg" {
				deviceMap.SDG = &linodego.InstanceConfigDevice{}
				if err := assignConfigDevice(deviceMap.SDG, dev, diskIDLabelMap); err != nil {
					return nil, err
				}
			}
			if k == "sdh" {
				deviceMap.SDH = &linodego.InstanceConfigDevice{}
				if err := assignConfigDevice(deviceMap.SDH, dev, diskIDLabelMap); err != nil {
					return nil, err
				}
			}
		}
	}
	return deviceMap, nil
}

type volumeDetacher func(context.Context, int, string) error

func makeVolumeDetacher(client linodego.Client, d *schema.ResourceData) volumeDetacher {
	return func(ctx context.Context, volumeID int, reason string) error {
		log.Printf("[INFO] Detaching Linode Volume %d %s", volumeID, reason)
		if err := client.DetachVolume(ctx, volumeID); err != nil {
			return err
		}

		log.Printf("[INFO] Waiting for Linode Volume %d to detach ...", volumeID)
		if _, err := client.WaitForVolumeLinodeID(ctx, volumeID, nil, int(d.Timeout("update").Seconds())); err != nil {
			return err
		}
		return nil
	}
}

func expandInstanceConfigDevice(m map[string]interface{}) *linodego.InstanceConfigDevice {
	var dev *linodego.InstanceConfigDevice
	// be careful of `disk_label string` in m
	if diskID, ok := m["disk_id"]; ok && diskID.(int) > 0 {
		dev = &linodego.InstanceConfigDevice{
			DiskID: diskID.(int),
		}
	} else if volumeID, ok := m["volume_id"]; ok && volumeID.(int) > 0 {
		dev = &linodego.InstanceConfigDevice{
			VolumeID: m["volume_id"].(int),
		}
	}

	return dev
}

func createDiskFromSet(client linodego.Client, instance linodego.Instance, v interface{}, d *schema.ResourceData) (*linodego.InstanceDisk, error) {
	disk, ok := v.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("Error converting disk from Terraform Set to golang map")
	}

	diskOpts := linodego.InstanceDiskCreateOptions{
		Label:      disk["label"].(string),
		Filesystem: disk["filesystem"].(string),
		Size:       disk["size"].(int),
	}

	if image, ok := disk["image"]; ok {
		diskOpts.Image = image.(string)

		if rootPass, ok := disk["root_pass"]; ok {
			diskOpts.RootPass = rootPass.(string)
		}

		if authorizedKeys, ok := disk["authorized_keys"]; ok {
			for _, sshKey := range authorizedKeys.([]interface{}) {
				diskOpts.AuthorizedKeys = append(diskOpts.AuthorizedKeys, sshKey.(string))
			}
		}

		if stackscriptID, ok := disk["stackscript_id"]; ok {
			diskOpts.StackscriptID = stackscriptID.(int)
		}

		if stackscriptData, ok := disk["stackscript_data"]; ok {
			for name, value := range stackscriptData.(map[string]interface{}) {
				diskOpts.StackscriptData[name] = value.(string)
			}
		}

		/*
			if sshKeys, ok := d.GetOk("authorized_keys"); ok {
				if sshKeysArr, ok := sshKeys.([]interface{}); ok {
					diskOpts.AuthorizedKeys = make([]string, len(sshKeysArr))
					for k, v := range sshKeys.([]interface{}) {
						if val, ok := v.(string); ok {
							diskOpts.AuthorizedKeys[k] = val
						}
					}
				}
			}
		*/
	}

	instanceDisk, err := client.CreateInstanceDisk(context.Background(), instance.ID, diskOpts)

	if err != nil {
		return nil, fmt.Errorf("Error creating Linode instance %d disk: %s", instance.ID, err)
	}

	_, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionDiskCreate, instanceDisk.Created, int(d.Timeout(schema.TimeoutCreate).Seconds()))
	if err != nil {
		return nil, fmt.Errorf("Error waiting for Linode instance %d disk: %s", instanceDisk.ID, err)
	}

	return instanceDisk, err
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

// changeInstanceType resizes the Linode Instance
func changeInstanceType(client *linodego.Client, instance *linodego.Instance, targetType string, d *schema.ResourceData) error {
	// Instance must be either offline or running (with no extra activity) to resize.
	if instance.Status == linodego.InstanceOffline || instance.Status == linodego.InstanceShuttingDown {
		if _, err := client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceOffline, int(d.Timeout(schema.TimeoutUpdate).Seconds())); err != nil {
			return fmt.Errorf("Error waiting for instance %d to go offline: %s", instance.ID, err)
		}
	} else {
		if _, err := client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceRunning, int(d.Timeout(schema.TimeoutUpdate).Seconds())); err != nil {
			return fmt.Errorf("Error waiting for instance %d readiness: %s", instance.ID, err)
		}
	}

	if err := client.ResizeInstance(context.Background(), instance.ID, targetType); err != nil {
		return fmt.Errorf("Error resizing instance %d: %s", instance.ID, err)
	}

	_, err := client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeResize, *instance.Created, int(d.Timeout(schema.TimeoutUpdate).Seconds()))
	if err != nil {
		return fmt.Errorf("Error waiting for instance %d to finish resizing: %s", instance.ID, err)
	}

	return nil
}

func changeInstanceDiskSize(client *linodego.Client, instance *linodego.Instance, disk *linodego.InstanceDisk, targetSize int, d *schema.ResourceData) error {
	if instance.Specs.Disk > targetSize {
		client.ResizeInstanceDisk(context.Background(), instance.ID, disk.ID, targetSize)

		// Wait for the Disk Resize Operation to Complete
		// waitForEventComplete(client, instance.ID, "linode_resize", waitMinutes)
		_, err := client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionDiskResize, disk.Updated, int(d.Timeout(schema.TimeoutUpdate).Seconds()))
		if err != nil {
			return fmt.Errorf("Error waiting for resize of Instance %d Disk %d: %s", instance.ID, disk.ID, err)
		}
	} else {
		return fmt.Errorf("Error resizing Disk %d: size exceeds disk size for Instance %d", disk.ID, instance.ID)
	}
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
	case string:
		return schema.HashString(v)
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

func assignConfigDevice(device *linodego.InstanceConfigDevice, dev map[string]interface{}, diskIDLabelMap map[string]int) error {
	if label, ok := dev["disk_label"].(string); ok && len(label) > 0 {
		if dev["disk_id"], ok = diskIDLabelMap[label]; !ok {
			return fmt.Errorf("Error mapping disk label %s to ID", dev["disk_label"])
		}
	}
	expanded := expandInstanceConfigDevice(dev)
	*device = *expanded
	return nil
}

// detachConfigVolumes detaches any volumes associated with an InstanceConfig.Devices struct
func detachConfigVolumes(dmap linodego.InstanceConfigDeviceMap, detacher volumeDetacher) error {
	// Preallocate our slice of config devices
	drives := []*linodego.InstanceConfigDevice{
		dmap.SDA, dmap.SDB, dmap.SDC, dmap.SDD, dmap.SDE, dmap.SDF, dmap.SDG, dmap.SDH,
	}

	// Make a buffered error channel for our goroutines to send error values back on
	errCh := make(chan error, len(drives))

	// Make a sync.WaitGroup so our devices can signal they're finished
	var wg sync.WaitGroup
	wg.Add(len(drives))

	// For each drive, spawn a goroutine to detach the volume, send an error on the err channel
	// if one exists, and signal the worker process is done
	for _, d := range drives {
		go func(dev *linodego.InstanceConfigDevice) {
			defer wg.Done()

			if dev != nil && dev.VolumeID > 0 {
				err := detacher(context.Background(), dev.VolumeID, "for config attachment")
				if err != nil {
					errCh <- err
				}
			}
		}(d)
	}

	// Wait until all processes are finished and close the error channel so we can range over it
	wg.Wait()
	close(errCh)

	// Build the error from the errors in the channel and return the combined error if any exist
	var errStr string
	for err := range errCh {
		if len(errStr) == 0 {
			errStr += ", "
		}

		errStr += err.Error()
	}

	if len(errStr) > 0 {
		return fmt.Errorf("Error detaching volumes: %s", errStr)
	}

	return nil
}
