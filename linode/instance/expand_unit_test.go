//go:build unit

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
	tests := []struct {
		name string
		m    map[string]interface{}
		want *linodego.InstanceConfigDevice
	}{
		{
			name: "Valid DiskID",
			m: map[string]interface{}{
				"disk_id": 123,
			},
			want: &linodego.InstanceConfigDevice{
				DiskID: 123,
			},
		},
		{
			name: "Valid VolumeID",
			m: map[string]interface{}{
				"volume_id": 456,
			},
			want: &linodego.InstanceConfigDevice{
				VolumeID: 456,
			},
		},
		{
			name: "Invalid IDs",
			m: map[string]interface{}{
				"disk_id":   0,
				"volume_id": -1,
			},
			want: nil,
		},
		{
			name: "No IDs",
			m:    map[string]interface{}{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandInstanceConfigDevice(tt.m)
			if got != nil && tt.want == nil || got == nil && tt.want != nil {
				t.Errorf("expandInstanceConfigDevice() = %v, want %v", got, tt.want)
			} else if got != nil && tt.want != nil {
				if got.DiskID != tt.want.DiskID || got.VolumeID != tt.want.VolumeID {
					t.Errorf("expandInstanceConfigDevice() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestExpandConfigInterface(t *testing.T) {
	interfaceInput := map[string]interface{}{
		"label":        "eth0.100",
		"purpose":      "vlan",
		"ipam_address": "192.168.1.2/24",
	}

	interfaceResult := expandConfigInterface(interfaceInput)

	expectedLabel := "eth0.100"
	if interfaceResult.Label != expectedLabel {
		t.Errorf("Expected label %s, but got %s", expectedLabel, interfaceResult.Label)
	}

	expectedPurpose := linodego.InterfacePurposeVLAN
	if interfaceResult.Purpose != expectedPurpose {
		t.Errorf("Expected purpose %s, but got %s", expectedPurpose, interfaceResult.Purpose)
	}

	expectedIPAMAddress := "192.168.1.2/24"
	if interfaceResult.IPAMAddress != expectedIPAMAddress {
		t.Errorf("Expected IPAMAddress %s, but got %s", expectedIPAMAddress, interfaceResult.IPAMAddress)
	}
}
