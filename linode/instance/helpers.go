package instance

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"golang.org/x/crypto/sha3"
)

var (
	boolFalse = false
	boolTrue  = true
)

type diskSpec map[string]interface{}

// getDeadlineSeconds gets the seconds remaining until deadline is met.
func getDeadlineSeconds(ctx context.Context, d *schema.ResourceData) int {
	duration := d.Timeout(schema.TimeoutUpdate)
	if deadline, ok := ctx.Deadline(); ok {
		duration = time.Until(deadline)
	}
	return int(duration.Seconds())
}

func createInstanceConfigsFromSet(
	ctx context.Context,
	client linodego.Client,
	instanceID int,
	cset []interface{},
	diskIDLabelMap map[string]int,
	detacher volumeDetacher,
) (map[int]linodego.InstanceConfig, error) {
	configIDMap := make(map[int]linodego.InstanceConfig, len(cset))
	for _, v := range cset {
		config := v.(map[string]interface{})

		configOpts := linodego.InstanceConfigCreateOptions{}
		configOpts.Kernel = config["kernel"].(string)
		configOpts.Label = config["label"].(string)
		configOpts.Comments = config["comments"].(string)
		configOpts.VirtMode = config["virt_mode"].(string)
		configOpts.MemoryLimit = config["memory_limit"].(int)
		configOpts.RunLevel = config["run_level"].(string)

		if helpers, ok := config["helpers"].([]interface{}); ok {
			for _, helper := range helpers {
				if helperMap, helperMapOk := helper.(map[string]interface{}); helperMapOk {
					configOpts.Helpers = &linodego.InstanceConfigHelpers{}
					if updateDBDisabled, ok := helperMap["updatedb_disabled"]; ok {
						configOpts.Helpers.UpdateDBDisabled = updateDBDisabled.(bool)
					}
					if distro, ok := helperMap["distro"]; ok {
						configOpts.Helpers.Distro = distro.(bool)
					}
					if modulesDep, ok := helperMap["modules_dep"]; ok {
						configOpts.Helpers.ModulesDep = modulesDep.(bool)
					}
					if network, ok := helperMap["network"]; ok {
						configOpts.Helpers.Network = network.(bool)
					}
					if devTmpFsAutomount, ok := helperMap["devtmpfs_automount"]; ok {
						configOpts.Helpers.DevTmpFsAutomount = devTmpFsAutomount.(bool)
					}
				}
			}
		}

		if interfaces, ok := config["interface"]; ok {
			interfaces := interfaces.([]interface{})
			configOpts.Interfaces = make([]linodego.InstanceConfigInterfaceCreateOptions, len(interfaces))

			for i, ni := range interfaces {
				configOpts.Interfaces[i] = helper.ExpandConfigInterface(ni.(map[string]interface{}))
			}
		}

		rootDevice := config["root_device"].(string)
		if rootDevice != "" {
			configOpts.RootDevice = &rootDevice
		}

		// configOpts.InitRD = config["initrd"].(string)
		// TODO(displague) need a disk_label to initrd lookup?
		devices, ok := config["devices"].([]interface{})
		if !ok {
			return configIDMap, fmt.Errorf("Error converting config devices")
		}
		for _, device := range devices {
			deviceMap := device.(map[string]interface{})
			confDevices, err := expandInstanceConfigDeviceMap(deviceMap, diskIDLabelMap)
			if err != nil {
				return configIDMap, err
			}
			configOpts.Devices = *confDevices
		}

		if err := detachConfigVolumes(ctx, configOpts.Devices, detacher); err != nil {
			return configIDMap, err
		}

		tflog.Info(ctx, "Creating instance config", map[string]any{
			"body": configOpts,
		})

		instanceConfig, err := client.CreateInstanceConfig(ctx, instanceID, configOpts)
		if err != nil {
			return configIDMap, fmt.Errorf("Error creating Instance Config: %s", err)
		}
		configIDMap[instanceConfig.ID] = *instanceConfig
	}
	return configIDMap, nil
}

func updateInstanceConfigs(
	ctx context.Context,
	client linodego.Client,
	d *schema.ResourceData,
	instance linodego.Instance,
	tfConfigsOld, tfConfigsNew interface{},
	diskIDLabelMap map[string]int,
	bootConfigLabel string,
) (bool, map[string]int, []*linodego.InstanceConfig, error) {
	var updatedConfigMap map[string]int
	var rebootInstance bool
	var updatedConfigs []*linodego.InstanceConfig

	configs, err := client.ListInstanceConfigs(ctx, instance.ID, nil)
	if err != nil {
		return rebootInstance, updatedConfigMap, updatedConfigs, fmt.Errorf(
			"Error fetching the config for Instance %d: %s", instance.ID, err)
	}

	configMap := make(map[string]linodego.InstanceConfig, len(configs))
	for _, config := range configs {
		if _, duplicate := configMap[config.Label]; duplicate {
			return rebootInstance, updatedConfigMap, updatedConfigs, fmt.Errorf(
				"Error indexing Instance %d Configs: Label '%s' is assigned to multiple configs", instance.ID, config.Label)
		}
		configMap[config.Label] = config
	}

	oldConfigLabels := make([]string, 0)

	var oldBootInterfaces, newBootInterfaces []string
	for _, tfConfigOld := range tfConfigsOld.([]interface{}) {
		if oldConfig, ok := tfConfigOld.(map[string]interface{}); ok {
			oldLabel := oldConfig["label"].(string)
			oldConfigLabels = append(oldConfigLabels, oldLabel)
			if oldLabel == bootConfigLabel {
				for _, iface := range oldConfig["interface"].([]interface{}) {
					oldBootInterfaces = append(oldBootInterfaces, iface.(map[string]interface{})["ipam_address"].(string))
				}
			}
		}
	}
	tfConfigs := tfConfigsNew.([]interface{})
	updatedConfigs = make([]*linodego.InstanceConfig, 0)
	updatedConfigMap = make(map[string]int, len(tfConfigs))
	for _, tfConfig := range tfConfigs {
		tfc, _ := tfConfig.(map[string]interface{})
		label, _ := tfc["label"].(string)
		rootDevice, _ := tfc["root_device"].(string)
		if existingConfig, existing := configMap[label]; existing {
			configUpdateOpts := existingConfig.GetUpdateOptions()
			configUpdateOpts.Kernel = tfc["kernel"].(string)
			configUpdateOpts.RunLevel = tfc["run_level"].(string)
			configUpdateOpts.VirtMode = tfc["virt_mode"].(string)
			configUpdateOpts.RootDevice = rootDevice
			configUpdateOpts.Comments = tfc["comments"].(string)
			configUpdateOpts.MemoryLimit = tfc["memory_limit"].(int)

			tfcHelpersRaw, helpersFound := tfc["helpers"]
			if tfcHelpers, ok := tfcHelpersRaw.([]interface{}); helpersFound && ok {
				helpersMap := tfcHelpers[0].(map[string]interface{})
				configUpdateOpts.Helpers = &linodego.InstanceConfigHelpers{
					UpdateDBDisabled:  helpersMap["updatedb_disabled"].(bool),
					Distro:            helpersMap["distro"].(bool),
					ModulesDep:        helpersMap["modules_dep"].(bool),
					Network:           helpersMap["network"].(bool),
					DevTmpFsAutomount: helpersMap["devtmpfs_automount"].(bool),
				}

			}

			configUpdateOpts.Interfaces = make([]linodego.InstanceConfigInterfaceCreateOptions, 0)

			if interfaces, ok := tfc["interface"]; ok {
				interfaces := interfaces.([]interface{})

				configUpdateOpts.Interfaces = make([]linodego.InstanceConfigInterfaceCreateOptions, len(interfaces))

				for i, ni := range interfaces {
					mappedInterface := ni.(map[string]interface{})
					configUpdateOpts.Interfaces[i] = helper.ExpandConfigInterface(mappedInterface)
					if label == bootConfigLabel {
						newBootInterfaces = append(newBootInterfaces, mappedInterface["ipam_address"].(string))
					}
				}

			}

			if reflect.DeepEqual(configUpdateOpts.Interfaces, existingConfig.GetUpdateOptions().Interfaces) {
				configUpdateOpts.Interfaces = nil
			}

			tfcDevicesRaw, devicesFound := tfc["devices"]
			if tfcDevices, ok := tfcDevicesRaw.([]interface{}); devicesFound && ok {
				devices := tfcDevices[0].(map[string]interface{})

				configUpdateOpts.Devices, err = expandInstanceConfigDeviceMap(devices, diskIDLabelMap)
				if err != nil {
					return rebootInstance, updatedConfigMap, updatedConfigs, err
				}
				if configUpdateOpts.Devices != nil && emptyConfigDeviceMap(*configUpdateOpts.Devices) {
					configUpdateOpts.Devices = nil
				}
			} else {
				configUpdateOpts.Devices = nil
			}

			if configUpdateOpts.Devices != nil {
				detacher := makeVolumeDetacherIgnoreAttached(client, d)

				if detachErr := detachConfigVolumes(ctx, *configUpdateOpts.Devices, detacher); detachErr != nil {
					return rebootInstance, updatedConfigMap, updatedConfigs, detachErr
				}
			}

			tflog.Debug(ctx, "Updating instance config", map[string]any{
				"body":      configUpdateOpts,
				"config_id": existingConfig.ID,
			})
			tflog.Trace(ctx, "Config update options", map[string]any{"value": configUpdateOpts.Interfaces})
			updatedConfig, err := client.UpdateInstanceConfig(ctx, instance.ID, existingConfig.ID, configUpdateOpts)
			if err != nil {
				return rebootInstance, updatedConfigMap, updatedConfigs, fmt.Errorf(
					"Error updating Instance %d Config %d: %s", instance.ID, existingConfig.ID, err)
			}

			updatedConfigMap[updatedConfig.Label] = updatedConfig.ID
		} else {
			detacher := makeVolumeDetacher(client, d)

			configIDMap, err := createInstanceConfigsFromSet(
				ctx, client, instance.ID, []interface{}{tfc}, diskIDLabelMap, detacher)
			if err != nil {
				return rebootInstance, updatedConfigMap, updatedConfigs, err
			}
			for _, config := range configIDMap {
				cfg := config

				updatedConfigMap[config.Label] = config.ID
				updatedConfigs = append(updatedConfigs, &cfg)
			}
		}
	}

	if !reflect.DeepEqual(oldBootInterfaces, newBootInterfaces) {
		rebootInstance = true
	}

	updatedConfigMap, err = deleteInstanceConfigs(ctx, client, instance.ID, oldConfigLabels, updatedConfigMap, configMap)
	if err != nil {
		return rebootInstance, updatedConfigMap, updatedConfigs, err
	}

	return rebootInstance, updatedConfigMap, updatedConfigs, nil
}

func deleteInstanceConfigs(
	ctx context.Context,
	client linodego.Client,
	instanceID int,
	oldConfigLabels []string,
	newConfigLabels map[string]int,
	configMap map[string]linodego.InstanceConfig,
) (map[string]int, error) {
	for _, oldLabel := range oldConfigLabels {
		if _, found := newConfigLabels[oldLabel]; !found {
			if listedConfig, found := configMap[oldLabel]; found {
				tflog.Info(ctx, "Deleting instance config", map[string]any{
					"config_id": listedConfig.ID,
				})
				if err := client.DeleteInstanceConfig(ctx, instanceID, listedConfig.ID); err != nil {
					return newConfigLabels, err
				}
				delete(newConfigLabels, oldLabel)
			}
		}
	}
	return newConfigLabels, nil
}

// changeInstanceConfigDevice returns a copy of a config device map with the specified disk slot changed to the
// provided device.
func changeInstanceConfigDevice(
	deviceMap linodego.InstanceConfigDeviceMap,
	namedSlot string,
	device *linodego.InstanceConfigDevice,
) linodego.InstanceConfigDeviceMap {
	tDevice := device
	if tDevice != nil && emptyInstanceConfigDevice(*tDevice) {
		tDevice = nil
	}
	switch namedSlot {
	case "sda":
		deviceMap.SDA = tDevice
	case "sdb":
		deviceMap.SDB = tDevice
	case "sdc":
		deviceMap.SDC = tDevice
	case "sdd":
		deviceMap.SDD = tDevice
	case "sde":
		deviceMap.SDE = tDevice
	case "sdf":
		deviceMap.SDF = tDevice
	case "sdg":
		deviceMap.SDG = tDevice
	case "sdh":
		deviceMap.SDH = tDevice
	}

	return deviceMap
}

// emptyInstanceConfigDevice returns true only when neither the disk or volume have been assigned to a config device.
func emptyInstanceConfigDevice(dev linodego.InstanceConfigDevice) bool {
	return (dev.DiskID == 0 && dev.VolumeID == 0)
}

// emptyConfigDeviceMap returns true only when none of the disks in a config device map have been assigned.
func emptyConfigDeviceMap(dmap linodego.InstanceConfigDeviceMap) bool {
	drives := []*linodego.InstanceConfigDevice{
		dmap.SDA, dmap.SDB, dmap.SDC, dmap.SDD, dmap.SDE, dmap.SDF, dmap.SDG, dmap.SDH,
	}
	empty := true
	for _, drive := range drives {
		if drive != nil && !emptyInstanceConfigDevice(*drive) {
			empty = false
			break
		}
	}
	return empty
}

type volumeDetacher func(context.Context, int, string) error

func makeVolumeDetacher(client linodego.Client, d *schema.ResourceData) volumeDetacher {
	return func(ctx context.Context, volumeID int, reason string) error {
		ctx = helper.SetLogFieldBulk(ctx, map[string]any{
			"volume_id": volumeID,
			"reason":    reason,
		})

		tflog.Info(ctx, "Detaching Linode Volume")
		if err := client.DetachVolume(ctx, volumeID); err != nil {
			return err
		}

		tflog.Debug(ctx, "Waiting for Linode Volume to detach ...")
		if _, err := client.WaitForVolumeLinodeID(ctx, volumeID, nil, getDeadlineSeconds(ctx, d)); err != nil {
			return err
		}

		tflog.Debug(ctx, "Volume finished detaching")

		return nil
	}
}

func makeVolumeDetacherIgnoreAttached(client linodego.Client, d *schema.ResourceData) volumeDetacher {
	return func(ctx context.Context, volumeID int, reason string) error {
		ctx = helper.SetLogFieldBulk(ctx, map[string]any{
			"volume_id": volumeID,
			"reason":    reason,
		})

		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return err
		}

		vol, err := client.GetVolume(ctx, volumeID)
		if err != nil {
			return err
		}

		if vol.LinodeID != nil && *vol.LinodeID == id {
			tflog.Info(ctx, "Volume is already attached to Linode, ignoring ...")
			return nil
		}

		return makeVolumeDetacher(client, d)(ctx, volumeID, reason)
	}
}

func createInstanceDisk(
	ctx context.Context,
	client linodego.Client,
	instance linodego.Instance,
	disk diskSpec,
	d *schema.ResourceData,
) (*linodego.InstanceDisk, error) {
	diskOpts := linodego.InstanceDiskCreateOptions{
		Label:      disk["label"].(string),
		Filesystem: disk["filesystem"].(string),
		Size:       disk["size"].(int),
	}

	if image, ok := disk["image"]; ok {
		diskOpts.Image = image.(string)

		if rootPass, ok := disk["root_pass"]; ok && rootPass != "" {
			diskOpts.RootPass = rootPass.(string)
		} else {
			var err error
			diskOpts.RootPass, err = helper.CreateRandomRootPassword()
			if err != nil {
				return nil, err
			}
		}

		if authorizedKeys, ok := disk["authorized_keys"]; ok {
			for _, sshKey := range authorizedKeys.([]interface{}) {
				if sshKey == nil {
					return nil, fmt.Errorf(
						"invalid input for disk authorized_keys: keys cannot be empty or null")
				}

				diskOpts.AuthorizedKeys = append(diskOpts.AuthorizedKeys, sshKey.(string))
			}
		}

		if authorizedUsers, ok := disk["authorized_users"]; ok {
			for _, user := range authorizedUsers.([]interface{}) {
				if user == nil {
					return nil, fmt.Errorf(
						"invalid input for disk authorized_users: users cannot be empty or null")
				}

				diskOpts.AuthorizedUsers = append(diskOpts.AuthorizedUsers, user.(string))
			}
		}

		if stackscriptID, ok := disk["stackscript_id"]; ok {
			diskOpts.StackscriptID = stackscriptID.(int)
		}

		if stackscriptDataRaw, ok := disk["stackscript_data"]; ok {
			stackscriptData, ok := stackscriptDataRaw.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Error parsing stackscript_data: expected map[string]interface{}")
			}
			diskOpts.StackscriptData = make(map[string]string, len(stackscriptData))
			for name, value := range stackscriptData {
				diskOpts.StackscriptData[name] = value.(string)
			}
		}
	}

	tflog.Info(ctx, "Creating new instance disk", map[string]any{
		"body": diskOpts,
	})

	p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionDiskCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize event poller: %s", err)
	}

	instanceDisk, err := client.CreateInstanceDisk(ctx, instance.ID, diskOpts)
	if err != nil {
		return nil, fmt.Errorf("Error creating Linode instance %d disk: %s", instance.ID, err)
	}

	_, err = p.WaitForFinished(ctx, getDeadlineSeconds(ctx, d))
	if err != nil {
		return nil, fmt.Errorf("Error waiting for Linode instance %d disk: %s", instanceDisk.ID, err)
	}

	tflog.Debug(ctx, "Instance disk creation complete", map[string]any{
		"disk_id": instanceDisk.ID,
	})

	return instanceDisk, err
}

// getInstanceDisks returns a map of disks for a given instance that is indexed by label.
func getInstanceDisks(
	ctx context.Context, client linodego.Client, instanceID int,
) (map[string]linodego.InstanceDisk, error) {
	disks, err := client.ListInstanceDisks(ctx, instanceID, nil)
	if err != nil {
		return nil, fmt.Errorf("Error fetching the disks for Instance %d: %s", instanceID, err)
	}

	diskMap := make(map[string]linodego.InstanceDisk)
	for _, disk := range disks {
		if _, duplicate := diskMap[disk.Label]; duplicate {
			return nil, fmt.Errorf("Error indexing Instance %d Disks: Label '%s' is assigned to multiple disks",
				instanceID, disk.Label)
		}
		diskMap[disk.Label] = disk
	}
	return diskMap, nil
}

// getInstanceDiskLabelIDMap returns a map of an instances disk labels to their corresponding IDs.
func getInstanceDiskLabelIDMap(
	ctx context.Context, client linodego.Client, d *schema.ResourceData, instanceID int,
) (map[string]int, error) {
	diskSpecs := d.Get("disk").([]interface{})
	disks, err := getInstanceDisks(ctx, client, instanceID)
	if err != nil {
		return nil, err
	}

	labelIDMap := make(map[string]int)
	for _, spec := range diskSpecs {
		label := spec.(map[string]interface{})["label"].(string)
		disk, ok := disks[label]
		if !ok {
			return nil, fmt.Errorf(`could not map disk label "%s" to an ID; not found`, label)
		}
		labelIDMap[label] = disk.ID
	}
	return labelIDMap, nil
}

// getInstanceDiskSpecChange returns a map of disk specs indexed by label.
func getInstanceDiskSpecChange(d *schema.ResourceData) (oldDiskSpecs, newDiskSpecs map[string]diskSpec) {
	old, new := d.GetChange("disk")
	oldDisk := old.([]interface{})
	newDisk := new.([]interface{})
	oldDiskSpecs = make(map[string]diskSpec)
	newDiskSpecs = make(map[string]diskSpec)

	for _, spec := range oldDisk {
		spec := spec.(map[string]interface{})
		oldDiskSpecs[spec["label"].(string)] = spec
	}
	for _, spec := range newDisk {
		spec := spec.(map[string]interface{})
		newDiskSpecs[spec["label"].(string)] = spec
	}
	return
}

// getInstanceDiskSpecDiffs sorts the disk specs by added, removed, and existing.
func getInstanceDiskSpecDiffs(
	oldDiskSpecs, newDiskSpecs map[string]diskSpec,
) (added, removed, existing map[string]diskSpec) {
	added = make(map[string]diskSpec)
	removed = make(map[string]diskSpec)
	existing = make(map[string]diskSpec)

	placed := make(map[string]struct{})
	for label, spec := range newDiskSpecs {
		_, exists := oldDiskSpecs[label]
		if exists {
			existing[label] = spec
		} else {
			added[label] = spec
		}
		placed[label] = struct{}{}
	}

	for label, spec := range oldDiskSpecs {
		if _, ok := placed[label]; !ok {
			removed[label] = spec
		}
	}
	return
}

// updateInstanceDisks ensures the disk specification matches the instance disk state. This means creating,
// updating, and deleting disks as needed.
//
// This function will also warn when there are disks attached to an instance which are not managed by
// terraform.
func updateInstanceDisks(
	ctx context.Context, client linodego.Client, d *schema.ResourceData, instance linodego.Instance,
) (bool, error) {
	oldDisk, newDisk := getInstanceDiskSpecChange(d)
	added, removed, existing := getInstanceDiskSpecDiffs(oldDisk, newDisk)
	disks, err := getInstanceDisks(ctx, client, instance.ID)
	if err != nil {
		return false, err
	}
	hasChanges := len(added)+len(removed) > 0

	// ensure all old disks which have not been removed are present on the instance
	for label := range oldDisk {
		_, wasRemoved := removed[label]
		if _, ok := disks[label]; !ok && !wasRemoved {
			return hasChanges, fmt.Errorf(`disk "%s" was not found on instance %d`, label, instance.ID)
		}
	}
	// keep track of all disks visited for accounting
	visited := make(map[string]struct{})

	// remove disks staged for removal
	for label := range removed {
		disk, ok := disks[label]
		if !ok {
			// It's ok if a removed disk is not found
			continue
		}

		tflog.Info(ctx, "Deleting unused disk", map[string]any{
			"disk_id": disk.ID,
		})

		p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionDiskDelete)
		if err != nil {
			return false, fmt.Errorf("failed to initialize event poller: %s", err)
		}

		if err := client.DeleteInstanceDisk(ctx, instance.ID, disk.ID); err != nil {
			return hasChanges, err
		}

		_, err = p.WaitForFinished(ctx, getDeadlineSeconds(ctx, d))
		if err != nil {
			return hasChanges, fmt.Errorf(
				"error waiting for Instance %d Disk %d to finish deleting: %s", instance.ID, disk.ID, err)
		}

		tflog.Debug(ctx, "Unused disk finished deleting", map[string]any{
			"disk_id": disk.ID,
		})

		visited[label] = struct{}{}
	}

	// ensure state is consistent with existing disks specs
	for label, spec := range existing {
		existingDisk := disks[label]
		// The only non-destructive change supported is resize.
		// Label renames are not supported because this TF provider relies on the label as an identifier.
		if spec["size"].(int) != existingDisk.Size {
			if err := changeInstanceDiskSize(ctx, &client, instance, existingDisk, spec["size"].(int), d); err != nil {
				return hasChanges, err
			}
			hasChanges = true
		}
		if spec["filesystem"].(string) != string(existingDisk.Filesystem) {
			return hasChanges, fmt.Errorf("failed to update disk %d; filesystems can not be changed", existingDisk.ID)
		}
		visited[label] = struct{}{}
	}

	// create disks staged for creation
	for _, spec := range added {
		if _, err := createInstanceDisk(ctx, client, instance, spec, d); err != nil {
			return hasChanges, err
		}
	}

	// ensure all disks visited
	for label, disk := range disks {
		if _, ok := visited[label]; !ok {
			// warn if disk found that is not in terraform state
			fmt.Printf("[WARN] found disk %s (%d) on instance %d not found in state", label, disk.ID, instance.ID)
		}
	}

	return hasChanges, nil
}

// sshKeyState hashes a string passed in as an interface.
func sshKeyState(val interface{}) string {
	return hashString(strings.Join(val.([]string), "\n"))
}

// rootPasswordState hashes a string passed in as an interface.
func rootPasswordState(val interface{}) string {
	return hashString(val.(string))
}

// hashString hashes a string.
func hashString(key string) string {
	hash := sha3.Sum512([]byte(key))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// ensureInstanceOffline ensures that a given instance is offline.
func ensureInstanceOffline(
	ctx context.Context, client *linodego.Client, instanceID, timeout int,
) (instance *linodego.Instance, err error) {
	if instance, err = client.GetInstance(ctx, instanceID); err != nil {
		return
	}

	if instance.Status == linodego.InstanceOffline {
		return
	} else if instance.Status != linodego.InstanceShuttingDown {
		tflog.Info(ctx, "Shutting down instance to reach offline status")
		err = client.ShutdownInstance(ctx, instanceID)
	}

	if err != nil {
		return
	}

	inst, err := client.WaitForInstanceStatus(ctx, instanceID, linodego.InstanceOffline, timeout)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "Instance has reached offline state")

	return inst, nil
}

// changeInstanceType resizes the Linode Instance.
func changeInstanceType(
	ctx context.Context,
	client *linodego.Client,
	instanceID int,
	targetType string,
	migrationType linodego.InstanceMigrationType,
	diskResize bool,
	d *schema.ResourceData,
) (*linodego.Instance, error) {
	instance, err := ensureInstanceOffline(ctx, client, instanceID, getDeadlineSeconds(ctx, d))
	if err != nil {
		return nil, err
	}

	var primaryDisk *linodego.InstanceDisk

	// We only need to make this request if we plan on resizing the disk
	if diskResize {
		primaryDisk, err = getPrimaryImplicitDisk(ctx, d, client, instanceID)
		if err != nil {
			return nil, err
		}
	}

	resizeOpts := linodego.InstanceResizeOptions{
		AllowAutoDiskResize: &diskResize,
		Type:                targetType,
		MigrationType:       migrationType,
	}

	tflog.Debug(ctx, "Resizing instance", map[string]any{
		"body":           resizeOpts,
		"target_type":    targetType,
		"migration_type": migrationType,
	})

	p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeResize)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize event poller %d: %s", instance.ID, err)
	}

	if err := client.ResizeInstance(ctx, instance.ID, resizeOpts); err != nil {
		return nil, fmt.Errorf("Error resizing Instance %d: %s", instance.ID, err)
	}

	_, err = p.WaitForFinished(ctx, getDeadlineSeconds(ctx, d))
	if err != nil {
		return nil, fmt.Errorf("Error waiting for instance %d to finish resizing: %s", instance.ID, err)
	}

	tflog.Debug(ctx, "Instance resize completed")

	// Wait for instance status to go back to idle, offline state
	if instance, err = client.WaitForInstanceStatus(
		ctx, instance.ID, linodego.InstanceOffline, getDeadlineSeconds(ctx, d),
	); err != nil {
		return nil, fmt.Errorf("Error waiting for Instance %d to enter offline state: %s", instance.ID, err)
	}

	// Sometimes the API falls behind on updating the disk info, let's make sure it has updated
	if diskResize && primaryDisk != nil {
		if _, err := waitForInstanceDiskSizeChange(ctx, client, instance.ID, primaryDisk.ID, primaryDisk.Size,
			getDeadlineSeconds(ctx, d)); err != nil {
			return nil, fmt.Errorf("failed to wait for disk resize: %s", err)
		}
	}

	return instance, nil
}

// returns the amount of disk space used by the new plan and old plan.
func getDiskSizeChange(oldDisk interface{}, newDisk interface{}) (int, int) {
	tfDisksOldInterface := oldDisk.([]interface{})
	tfDisksNewInterface := newDisk.([]interface{})

	oldDiskSize := 0
	newDiskSize := 0
	// Get total amount of disk usage before and after
	for _, disk := range tfDisksOldInterface {
		oldDiskSize += disk.(map[string]interface{})["size"].(int)
	}

	for _, disk := range tfDisksNewInterface {
		newDiskSize += disk.(map[string]interface{})["size"].(int)
	}
	return oldDiskSize, newDiskSize
}

func changeInstanceDiskSize(
	ctx context.Context,
	client *linodego.Client,
	instance linodego.Instance,
	disk linodego.InstanceDisk,
	targetSize int,
	d *schema.ResourceData,
) error {
	if instance.Specs.Disk < targetSize {
		return fmt.Errorf("Error resizing disk %d: size exceeds disk size for Instance %d", disk.ID, instance.ID)
	}

	ctx = tflog.SetField(ctx, "disk_id", disk.ID)

	switch instance.Status {
	case linodego.InstanceShuttingDown:
		tflog.Debug(ctx, "Waiting for instance to reach offline status for disk resize")

		if _, err := client.WaitForInstanceStatus(
			ctx, instance.ID, linodego.InstanceOffline, getDeadlineSeconds(ctx, d),
		); err != nil {
			return fmt.Errorf("Error waiting for Instance %d to go offline: %s", instance.ID, err)
		}
	case linodego.InstanceOffline:
	default:
		tflog.Debug(ctx, "Shutting down instance for disk resize")

		if err := client.ShutdownInstance(ctx, instance.ID); err != nil {
			return err
		}
	}

	// Wait for instance to go offline. Resize the disk once Linode is shut down.
	if _, err := client.WaitForInstanceStatus(
		ctx, instance.ID, linodego.InstanceOffline, getDeadlineSeconds(ctx, d),
	); err != nil {
		return fmt.Errorf("Error waiting for Instance %d to go offline: %s", instance.ID, err)
	}

	tflog.Info(ctx, "Instance has reached offline status, resizing disk")

	p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionDiskResize)
	if err != nil {
		return fmt.Errorf("failed to initialize event poller: %s", err)
	}

	if err := client.ResizeInstanceDisk(ctx, instance.ID, disk.ID, targetSize); err != nil {
		return fmt.Errorf("Error resizing disk %d for Instance %d: %s", disk.ID, instance.ID, err)
	}

	// Wait for the disk resize operation to complete.
	_, err = p.WaitForFinished(ctx, getDeadlineSeconds(ctx, d))
	if err != nil {
		return fmt.Errorf("Error waiting for resize of Instance %d Disk %d: %s", instance.ID, disk.ID, err)
	}

	// Check to see if the resize operation worked
	if updatedDisk, err := client.WaitForInstanceDiskStatus(ctx, instance.ID, disk.ID, linodego.DiskReady,
		getDeadlineSeconds(ctx, d)); err != nil {
		return fmt.Errorf("Error waiting disk %d on instance %d to be ready: %s", disk.ID, instance.ID, err)
	} else if updatedDisk.Size != targetSize {
		return fmt.Errorf(
			"Error resizing disk %d on instance %d from %d to %d", disk.ID, instance.ID, disk.Size, targetSize)
	}

	tflog.Debug(ctx, "Disk has finished resizing")

	return nil
}

// privateIP determines if an IP is for private use (RFC1918)
// https://stackoverflow.com/a/41273687
func privateIP(ip net.IP) bool {
	_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
	_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
	_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
	private := private24BitBlock.Contains(ip) || private20BitBlock.Contains(ip) || private16BitBlock.Contains(ip)
	return private
}

func assignConfigDevice(
	device *linodego.InstanceConfigDevice, dev map[string]interface{}, diskIDLabelMap map[string]int,
) error {
	if label, ok := dev["disk_label"].(string); ok && len(label) > 0 {
		if dev["disk_id"], ok = diskIDLabelMap[label]; !ok {
			return fmt.Errorf(`Error mapping disk label "%s" to ID`, dev["disk_label"])
		}
	}
	expanded := expandInstanceConfigDevice(dev)
	if expanded != nil {
		*device = *expanded
	}
	return nil
}

// getInstanceDefaultDisks gets the boot and swap disk for an instance which has implicit, default disks.
func getInstanceDefaultDisks(
	ctx context.Context,
	instanceID int,
	client *linodego.Client,
) (bootDisk, swapDisk *linodego.InstanceDisk, err error) {
	var disks []linodego.InstanceDisk
	disks, err = client.ListInstanceDisks(ctx, instanceID, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting instance managed disks: %s", err)
	}
	bootDisk = findDiskByFS(disks, linodego.FilesystemExt4)
	swapDisk = findDiskByFS(disks, linodego.FilesystemSwap)
	return
}

// getInstanceTypeChange checks to see if the linode itself was resized.
func getInstanceTypeChange(
	ctx context.Context,
	d *schema.ResourceData,
	client *linodego.Client,
) (oldSpec, newSpec *linodego.LinodeType, err error) {
	old, new := d.GetChange("type")
	oldType, newType := old.(string), new.(string)

	oldSpec, err = client.GetType(ctx, oldType)
	if err != nil {
		return
	}
	newSpec, err = client.GetType(ctx, newType)
	if err != nil {
		return
	}
	return
}

// applyInstanceDiskSpec checks to see if the staged disk changes can be supported by the instance specification's
// capacity. If there is sufficient space, it attempts to update the disks.
//
// returns bool describing whether change has occurred.
func applyInstanceDiskSpec(
	ctx context.Context,
	d *schema.ResourceData,
	client *linodego.Client,
	instance *linodego.Instance,
	typ *linodego.LinodeType,
) (bool, error) {
	if err := assertDiskConfigFitsInstanceType(d, typ); err != nil {
		return false, err
	}

	return updateInstanceDisks(ctx, *client, d, *instance)
}

// assertDiskConfigFitsInstanceType asserts that the cumulative disk space used by a given disk config fits a given
// linode type spec for disk capacity.
func assertDiskConfigFitsInstanceType(d *schema.ResourceData, typ *linodego.LinodeType) error {
	oldDisks, newDisks := d.GetChange("disk")

	_, newDiskSize := getDiskSizeChange(oldDisks, newDisks)

	if typ.Disk < newDiskSize {
		return fmt.Errorf(
			"linode type %s has insufficient disk capacity for the config. Have %d; want %d",
			typ.Label, typ.Disk, newDiskSize)
	}
	return nil
}

// applyInstanceTypeChange checks to see if the staged disk changes can be supported by the new instance
// specification. If there is sufficient space, it attempts to update the instance type.
func applyInstanceTypeChange(
	ctx context.Context,
	d *schema.ResourceData,
	client *linodego.Client,
	instance *linodego.Instance,
	typ *linodego.LinodeType,
) (*linodego.Instance, error) {
	resizeDisk := d.Get("resize_disk").(bool)
	migrationType := linodego.InstanceMigrationType(
		d.Get("migration_type").(string),
	)

	if resizeDisk {
		// Verify that there are implicit disks defined
		if d.GetRawConfig().GetAttr("image").IsNull() && d.GetRawConfig().GetAttr("disk").LengthInt() > 0 {
			return nil, fmt.Errorf("resize_disk requires that no explicit disks are defined")
		}

		if err := validateImplicitDisks(ctx, client, instance.ID); err != nil {
			return nil, err
		}

		usedSpace, err := getDiskSizeSum(ctx, d, client, instance.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate total disk size: %s", err)
		}

		if typ.Disk < usedSpace {
			return nil, fmt.Errorf("failed to resize instance and disk: %s", downsizeFailedMessage)
		}
	}

	if err := assertDiskConfigFitsInstanceType(d, typ); err != nil {
		return nil, err
	}

	return changeInstanceType(ctx, client, instance.ID, typ.ID, migrationType, resizeDisk, d)
}

// applyInstanceMigration synchronously migrates a Linode to a new region.
func applyInstanceMigration(
	ctx context.Context,
	d *schema.ResourceData,
	client *linodego.Client,
	instance *linodego.Instance,
	targetRegion string,
) (*linodego.Instance, error) {
	migrationType := linodego.InstanceMigrationType(
		d.Get("migration_type").(string),
	)

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"target_region":  targetRegion,
		"migration_type": migrationType,
	})

	tflog.Debug(ctx, "Migrating instance to new region")

	p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeMigrateDatacenter)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize event poller %d: %s", instance.ID, err)
	}

	if err := client.MigrateInstance(ctx, instance.ID, linodego.InstanceMigrateOptions{
		Region: targetRegion,
		Type:   migrationType,
	}); err != nil {
		return nil, fmt.Errorf("failed to migrate instance %d to region %s: %w", instance.ID, targetRegion, err)
	}

	_, err = p.WaitForFinished(ctx, getDeadlineSeconds(ctx, d))
	if err != nil {
		return nil, fmt.Errorf("failed to wait for instance %d to finish migration: %w", instance.ID, err)
	}

	tflog.Debug(ctx, "Instance migration has finished")

	result, err := client.GetInstance(ctx, instance.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh instance %d: %w", instance.ID, err)
	}

	return result, nil
}

// detachConfigVolumes detaches any volumes associated with an InstanceConfig.Devices struct.
func detachConfigVolumes(
	ctx context.Context, dmap linodego.InstanceConfigDeviceMap, detacher volumeDetacher,
) error {
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
				err := detacher(ctx, dev.VolumeID, "for config attachment")
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

func isInstanceBooted(instance *linodego.Instance) bool {
	// For diffing purposes, transition states need to be treated as
	// booted == true. This is because these statuses will eventually
	// result in a powered on Linode.
	return instance.Status == linodego.InstanceRunning ||
		instance.Status == linodego.InstanceRebooting ||
		instance.Status == linodego.InstanceBooting
}

func validateBooted(ctx context.Context, d *schema.ResourceData) error {
	booted := d.Get("booted").(bool)
	bootedNull := d.GetRawConfig().GetAttr("booted").IsNull()

	_, imageOk := d.GetOk("image")
	_, disksOk := d.GetOk("disk")
	_, configsOk := d.GetOk("config")

	if !bootedNull && booted && !imageOk && !(disksOk && configsOk) {
		return fmt.Errorf("booted requires an image or disk/config be defined")
	}

	return nil
}

func handleBootedUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}, instanceID, configID int,
) error {
	client := meta.(*helper.ProviderMeta).Client

	deadlineSeconds := getDeadlineSeconds(ctx, d)

	booted := d.Get("booted")
	bootedNull := d.GetRawConfig().GetAttr("booted").IsNull()

	if bootedNull {
		return nil
	}

	inst, err := client.GetInstance(ctx, instanceID)
	if err != nil {
		return err
	}

	instStatus, err := waitForRunningOrOfflineState(ctx, inst.Status, &client, instanceID)
	if err != nil {
		return err
	}

	// Boot or shutdown the instance if necessary
	if instStatus != linodego.InstanceRunning && booted.(bool) {
		if err := bootInstanceSync(ctx, client, instanceID, configID, deadlineSeconds); err != nil {
			return err
		}
	}

	if instStatus != linodego.InstanceOffline && !booted.(bool) {
		if err := shutDownInstanceSync(ctx, client, instanceID, deadlineSeconds); err != nil {
			return err
		}
	}

	return nil
}

func waitForRunningOrOfflineState(
	ctx context.Context, status linodego.InstanceStatus, client *linodego.Client, instanceID int,
) (instStatus linodego.InstanceStatus, err error) {
	// Ensure the Linode reaches a running or offline state
	// if in a transition state (shutting down, booting, rebooting)
	// TODO: clean up this logic
	switch status {

	// These cases can be ignored
	case linodego.InstanceRunning:
	case linodego.InstanceOffline:

	case linodego.InstanceShuttingDown:
		tflog.Info(ctx, "Awaiting instance shutdown before continuing")

		_, err = client.WaitForInstanceStatus(ctx, instanceID, linodego.InstanceOffline, 120)
		if err != nil {
			return "", fmt.Errorf("failed to wait for instance offline: %s", err)
		}

		instStatus = linodego.InstanceOffline

	case linodego.InstanceBooting, linodego.InstanceRebooting:
		tflog.Info(ctx, "Awaiting instance boot before continuing")

		_, err = client.WaitForInstanceStatus(ctx, instanceID, linodego.InstanceRunning, 120)
		if err != nil {
			return "", fmt.Errorf("failed to wait for instance running: %s", err)
		}

		instStatus = linodego.InstanceRunning

	default:
		return "", fmt.Errorf("instance is in unhandled state %s", status)
	}
	return instStatus, nil
}

func shutDownInstanceSync(ctx context.Context, client linodego.Client, instanceID, deadlineSeconds int) error {
	tflog.Info(ctx, "Shutting down instance")

	p, err := client.NewEventPoller(ctx, instanceID, linodego.EntityLinode, linodego.ActionLinodeShutdown)
	if err != nil {
		return fmt.Errorf("failed to initialize event poller: %s", err)
	}

	if err := client.ShutdownInstance(ctx, instanceID); err != nil {
		return fmt.Errorf("failed to shutdown instance: %s", err)
	}

	if _, err := p.WaitForFinished(ctx, deadlineSeconds); err != nil {
		return fmt.Errorf("failed to wait for instance shutdown: %s", err)
	}

	tflog.Debug(ctx, "Instance has finished shutting down")

	return nil
}

func bootInstanceSync(ctx context.Context, client linodego.Client, instanceID, configID, deadlineSeconds int) error {
	tflog.Info(ctx, "Booting instance", map[string]any{
		"config_id": configID,
	})

	p, err := client.NewEventPoller(ctx, instanceID, linodego.EntityLinode, linodego.ActionLinodeBoot)
	if err != nil {
		return fmt.Errorf("failed to initialize event poller: %s", err)
	}

	if err := client.BootInstance(ctx, instanceID, configID); err != nil {
		return fmt.Errorf("failed to boot instance: %s", err)
	}

	if _, err := p.WaitForFinished(ctx, deadlineSeconds); err != nil {
		return fmt.Errorf("failed to wait for instance boot: %s", err)
	}

	tflog.Debug(ctx, "Instance has finished booting")

	return nil
}

func getDiskSizeSum(ctx context.Context, d *schema.ResourceData,
	client *linodego.Client, instanceID int,
) (int, error) {
	disks, err := client.ListInstanceDisks(ctx, instanceID, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get instance disks: %s", err)
	}

	sum := 0
	for _, disk := range disks {
		sum += disk.Size
	}

	return sum, nil
}

func getFirstDiskWithFilesystem(disks []linodego.InstanceDisk,
	filesystems []linodego.DiskFilesystem,
) *linodego.InstanceDisk {
	for _, disk := range disks {
		for _, filesystem := range filesystems {
			if disk.Filesystem == filesystem {
				return &disk
			}
		}
	}

	return nil
}

func validateImplicitDisks(ctx context.Context,
	client *linodego.Client, instanceID int,
) error {
	disks, err := client.ListInstanceDisks(ctx, instanceID, nil)
	if err != nil {
		return fmt.Errorf("failed to get instance disks: %s", err)
	}

	// No disks are an acceptable case
	if len(disks) < 1 {
		return nil
	}

	if getFirstDiskWithFilesystem(disks,
		[]linodego.DiskFilesystem{linodego.FilesystemExt4, linodego.FilesystemExt3}) == nil || len(disks) > 2 {
		return fmt.Errorf("invalid implicit disk configuration: %s", invalidImplicitDiskConfigMessage)
	}

	return nil
}

func getPrimaryImplicitDisk(ctx context.Context, d *schema.ResourceData,
	client *linodego.Client, instanceID int,
) (*linodego.InstanceDisk, error) {
	disks, err := client.ListInstanceDisks(ctx, instanceID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance disks: %s", err)
	}

	if len(disks) > 2 {
		return nil, fmt.Errorf("invalid implicit disk configuration: %s", invalidImplicitDiskConfigMessage)
	}

	targetDisk := getFirstDiskWithFilesystem(disks,
		[]linodego.DiskFilesystem{linodego.FilesystemExt4, linodego.FilesystemExt3})

	return targetDisk, nil
}

func waitForInstanceDiskSizeChange(ctx context.Context, client *linodego.Client, instanceID,
	diskID, oldSize, timeoutSeconds int,
) (*linodego.InstanceDisk, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			disk, err := client.GetInstanceDisk(ctx, instanceID, diskID)
			if err != nil {
				return nil, err
			}

			if disk.Size != oldSize {
				return disk, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to wait for instance %d disk %d size change from %d: %s",
				instanceID, diskID, oldSize, ctx.Err())
		}
	}
}

func instanceIPSliceToString(ips []*linodego.InstanceIP) []string {
	result := make([]string, len(ips))

	for i, ip := range ips {
		result[i] = ip.Address
	}

	return result
}
