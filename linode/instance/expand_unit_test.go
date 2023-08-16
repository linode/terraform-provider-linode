package instance

import (
	"github.com/linode/linodego"
	"testing"
)

func TestExpandInstanceConfigDeviceMap(t *testing.T) {
	deviceMapInput := map[string]interface{}{
		"sda": []interface{}{
			map[string]interface{}{
				"disk_id":   124458,
				"volume_id": nil,
			},
		},
		"sdb": []interface{}{
			map[string]interface{}{
				"disk_id":   124459,
				"volume_id": nil,
			},
		},
	}

	diskIDLabelMap := map[string]int{
		"example_label_sda": 124458,
		"example_label_sdb": 124459,
	}

	deviceMap, err := expandInstanceConfigDeviceMap(deviceMapInput, diskIDLabelMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if deviceMap == nil {
		t.Error("Expected deviceMap to not be nil")
	}

	// Assert the DiskID for SDA and SDB
	expectedDiskID_SDA := 124458
	if deviceMap.SDA.DiskID != expectedDiskID_SDA {
		t.Errorf("Expected DiskID %d for SDA, but got %d", expectedDiskID_SDA, deviceMap.SDA.DiskID)
	}

	expectedDiskID_SDB := 124459
	if deviceMap.SDB.DiskID != expectedDiskID_SDB {
		t.Errorf("Expected DiskID %d for SDB, but got %d", expectedDiskID_SDB, deviceMap.SDB.DiskID)
	}

	// Assert that other device slots are nil
	if deviceMap.SDC != nil || deviceMap.SDD != nil || deviceMap.SDE != nil ||
		deviceMap.SDF != nil || deviceMap.SDG != nil || deviceMap.SDH != nil {
		t.Errorf("Expected all other device slots to be nil")
	}
}

func TestExpandInstanceConfigDevice(t *testing.T) {
	// Test case 1: Disk ID provided
	deviceInput1 := map[string]interface{}{
		"disk_id":   123,
		"volume_id": nil,
	}
	device1 := expandInstanceConfigDevice(deviceInput1)
	if device1 == nil {
		t.Error("Expected device1 to not be nil")
	}
	if device1.DiskID != 123 {
		t.Errorf("Expected DiskID %d, but got %d", 123, device1.DiskID)
	}

	// Test case 2: Volume ID provided
	deviceInput2 := map[string]interface{}{
		"disk_id":   nil,
		"volume_id": 456,
	}
	device2 := expandInstanceConfigDevice(deviceInput2)
	if device2 == nil {
		t.Error("Expected device2 to not be nil")
	}
	if device2.VolumeID != 456 {
		t.Errorf("Expected VolumeID %d, but got %d", 456, device2.VolumeID)
	}

	// Test case 3: No valid ID provided
	deviceInput3 := map[string]interface{}{
		"disk_id":   nil,
		"volume_id": nil,
	}
	device3 := expandInstanceConfigDevice(deviceInput3)
	if device3 != nil {
		t.Error("Expected device3 to be nil")
	}
}

func TestExpandConfigInterface(t *testing.T) {
	// Create example input data
	interfaceInput := map[string]interface{}{
		"label":        "eth0",
		"purpose":      "public",
		"ipam_address": "192.168.1.2",
	}

	// Call the function being tested
	interfaceResult := expandConfigInterface(interfaceInput)

	// Perform assertions
	expectedLabel := "eth0"
	if interfaceResult.Label != expectedLabel {
		t.Errorf("Expected label %s, but got %s", expectedLabel, interfaceResult.Label)
	}

	expectedPurpose := linodego.ConfigInterfacePurpose("public")
	if interfaceResult.Purpose != expectedPurpose {
		t.Errorf("Expected purpose %s, but got %s", expectedPurpose, interfaceResult.Purpose)
	}

	expectedIPAMAddress := "192.168.1.2"
	if interfaceResult.IPAMAddress != expectedIPAMAddress {
		t.Errorf("Expected IPAMAddress %s, but got %s", expectedIPAMAddress, interfaceResult.IPAMAddress)
	}
}
