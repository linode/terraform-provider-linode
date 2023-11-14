//go:build unit

package instance

import (
	"reflect"
	"testing"
	"time"

	"github.com/linode/linodego"

	"github.com/stretchr/testify/assert"
)

// Test helper functions
func isEqual(a, b []map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		for key := range a[i] {
			if a[i][key] != b[i][key] {
				return false
			}
		}
	}

	return true
}

func areMapsEqual(a, b map[string]interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// Unit tests for functions in flatten.go
func TestFlattenInstanceAlerts(t *testing.T) {
	instance := linodego.Instance{
		ID:      123,
		Created: &time.Time{},
		Updated: &time.Time{},
		Region:  "us-east",
		Alerts: &linodego.InstanceAlert{
			CPU:           180,
			IO:            10000,
			NetworkIn:     10,
			NetworkOut:    10,
			TransferQuota: 80,
		},
		Backups: &linodego.InstanceBackup{
			Available: true,
			Enabled:   true,
		},
		Image:       "linode/debian10",
		Group:       "Linode-Group",
		IPv6:        "c001:d00d::1337/128",
		Label:       "linode123",
		Type:        "g6-standard-1",
		Status:      linodego.InstanceStatus("running"),
		HasUserData: false,
		Hypervisor:  "kvm",
		HostUUID:    "3a3ddd59d9a78bb8de041391075df44de62bfec8",
		Specs: &linodego.InstanceSpec{
			Disk:     81920,
			GPUs:     0,
			Memory:   4096,
			Transfer: 4000,
			VCPUs:    2,
		},
		WatchdogEnabled: true,
		Tags:            []string{"example tag", "another example"},
	}
	alerts := flattenInstanceAlerts(instance)

	expectedAlerts := []map[string]int{
		{
			"cpu":            180,
			"io":             10000,
			"network_in":     10,
			"network_out":    10,
			"transfer_quota": 80,
		},
	}

	assert.Equal(t, expectedAlerts, alerts)
}

func TestFlattenInstanceBackups(t *testing.T) {
	instance := linodego.Instance{
		ID:      123,
		Created: &time.Time{},
		Updated: &time.Time{},
		Region:  "us-east",
		Backups: &linodego.InstanceBackup{
			Available: true,
			Enabled:   true,
			Schedule: struct {
				Day    string `json:"day,omitempty"`
				Window string `json:"window,omitempty"`
			}{
				Day:    "Saturday",
				Window: "W22",
			},
		},
	}

	backups := flattenInstanceBackups(instance)

	expectedBackups := []map[string]interface{}{
		{
			"available": true,
			"enabled":   true,
			"schedule": []map[string]interface{}{
				{
					"day":    "Saturday",
					"window": "W22",
				},
			},
		},
	}

	assert.Equal(t, expectedBackups, backups)
}

func TestFlattenInstanceDisks(t *testing.T) {
	instanceDisks := []linodego.InstanceDisk{
		{
			ID:         25674,
			Label:      "Debian 9 Disk",
			Status:     "ready",
			Size:       48640,
			Filesystem: "ext4",
			Created:    nil,
			Updated:    nil,
		},
	}

	disks, swapSize := flattenInstanceDisks(instanceDisks)

	expectedDisks := []map[string]interface{}{
		{
			"id":         25674,
			"size":       48640,
			"label":      "Debian 9 Disk",
			"filesystem": "ext4",
		},
	}

	assert.Equal(t, expectedDisks, disks, "Flattened disks do not match expected")
	assert.Equal(t, 0, swapSize, "Swap size does not match expected")
}

func TestFlattenInstanceConfigDevice(t *testing.T) {
	diskLabelIDMap := map[int]string{
		1: "DiskLabel1",
		2: "DiskLabel2",
	}

	deviceWithDisk := &linodego.InstanceConfigDevice{
		DiskID: 1,
	}
	resultDisk := flattenInstanceConfigDevice(deviceWithDisk, diskLabelIDMap)
	expectedResultDisk := []map[string]interface{}{
		{
			"disk_id":    1,
			"disk_label": "DiskLabel1",
		},
	}
	if !isEqual(resultDisk, expectedResultDisk) {
		t.Errorf("Expected %v, but got %v", expectedResultDisk, resultDisk)
	}

	deviceWithVolume := &linodego.InstanceConfigDevice{
		VolumeID: 3,
	}
	resultVolume := flattenInstanceConfigDevice(deviceWithVolume, diskLabelIDMap)
	expectedResultVolume := []map[string]interface{}{
		{
			"volume_id": 3,
		},
	}
	if !isEqual(resultVolume, expectedResultVolume) {
		t.Errorf("Expected %v, but got %v", expectedResultVolume, resultVolume)
	}

	emptyDevice := &linodego.InstanceConfigDevice{}
	resultEmpty := flattenInstanceConfigDevice(emptyDevice, diskLabelIDMap)
	if resultEmpty != nil {
		t.Errorf("Expected nil for an empty device, but got %v", resultEmpty)
	}
}

func TestFlattenInstanceConfigs(t *testing.T) {
	diskLabelIDMap := map[int]string{
		124458: "disk_label",
	}

	instanceConfigs := []linodego.InstanceConfig{
		{
			ID:       1,
			Label:    "config1",
			Comments: "test config",
			Devices: &linodego.InstanceConfigDeviceMap{
				SDA: &linodego.InstanceConfigDevice{
					DiskID:   124458,
					VolumeID: 1,
				},
				SDB: &linodego.InstanceConfigDevice{
					DiskID:   124458,
					VolumeID: 2,
				},
				SDC: &linodego.InstanceConfigDevice{
					DiskID:   124458,
					VolumeID: 3,
				},
				SDD: &linodego.InstanceConfigDevice{
					DiskID:   124458,
					VolumeID: 4,
				},
				SDE: &linodego.InstanceConfigDevice{
					DiskID:   124458,
					VolumeID: 5,
				},
				SDF: &linodego.InstanceConfigDevice{
					DiskID:   124458,
					VolumeID: 6,
				},
				SDG: &linodego.InstanceConfigDevice{
					DiskID:   124458,
					VolumeID: 7,
				},
				SDH: &linodego.InstanceConfigDevice{
					DiskID:   124458,
					VolumeID: 8,
				},
			},
			Helpers:     &linodego.InstanceConfigHelpers{},
			Interfaces:  []linodego.InstanceConfigInterface{},
			MemoryLimit: 2048,
			Kernel:      "linode/latest-64bit",
			RootDevice:  "/dev/sda",
			RunLevel:    "default",
			VirtMode:    "paravirt",
		},
	}

	expectedConfigs := []map[string]interface{}{
		{
			"id":           1,
			"root_device":  "/dev/sda",
			"kernel":       "linode/latest-64bit",
			"run_level":    "default",
			"virt_mode":    "paravirt",
			"comments":     "test config",
			"memory_limit": 2048,
			"label":        "config1",
			"helpers": []map[string]bool{
				{
					"updatedb_disabled":  false,
					"distro":             false,
					"modules_dep":        false,
					"network":            false,
					"devtmpfs_automount": false,
				},
			},
			"devices": []map[string]interface{}{
				{
					"sda": []map[string]interface{}{
						{
							"disk_id":    124458,
							"disk_label": "disk_label",
						},
					},
					"sdb": []map[string]interface{}{
						{
							"disk_id":    124458,
							"disk_label": "disk_label",
						},
					},
					"sdc": []map[string]interface{}{
						{
							"disk_id":    124458,
							"disk_label": "disk_label",
						},
					},
					"sdd": []map[string]interface{}{
						{
							"disk_id":    124458,
							"disk_label": "disk_label",
						},
					},
					"sde": []map[string]interface{}{
						{
							"disk_id":    124458,
							"disk_label": "disk_label",
						},
					},
					"sdf": []map[string]interface{}{
						{
							"disk_id":    124458,
							"disk_label": "disk_label",
						},
					},
					"sdg": []map[string]interface{}{
						{
							"disk_id":    124458,
							"disk_label": "disk_label",
						},
					},
					"sdh": []map[string]interface{}{
						{
							"disk_id":    124458,
							"disk_label": "disk_label",
						},
					},
				},
			},
			"interface": []interface{}{},
		},
	}

	actualConfigs := flattenInstanceConfigs(instanceConfigs, diskLabelIDMap)

	if len(actualConfigs) != len(expectedConfigs) {
		t.Fatalf("Expected %d configs, but got %d", len(expectedConfigs), len(actualConfigs))
	}

	for i := 0; i < len(expectedConfigs); i++ {
		if !areMapsEqual(actualConfigs[i], expectedConfigs[i]) {
			t.Errorf("Config %d: Mismatch between expected and actual config", i)
		}
	}
}

func TestFlattenInstanceSpecs(t *testing.T) {
	instance := linodego.Instance{
		ID:     1,
		Region: "us-east",
		Image:  "linode/debian10",
		Label:  "test-instance",
		Type:   "g6-standard-1",
		Status: linodego.InstanceRunning,
		Specs: &linodego.InstanceSpec{
			VCPUs:    2,
			Disk:     50,
			Memory:   4096,
			Transfer: 2000,
		},
	}

	result := flattenInstanceSpecs(instance)

	expected := []map[string]int{
		{
			"vcpus":    2,
			"disk":     50,
			"memory":   4096,
			"transfer": 2000,
		},
	}

	if len(result) != len(expected) {
		t.Errorf("Result slice length does not match the expected slice length")
	}

	for i := 0; i < len(result); i++ {
		for key, expectedValue := range expected[i] {
			if resultValue, ok := result[i][key]; !ok || resultValue != expectedValue {
				t.Errorf("Mismatch for key %s: Expected %d, but got %d", key, expectedValue, resultValue)
			}
		}
	}
}
