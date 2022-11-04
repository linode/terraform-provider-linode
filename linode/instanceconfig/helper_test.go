package instanceconfig

import (
	"testing"
)

func TestExpandDeviceMap(t *testing.T) {
	inputValue := make([]any, 1)
	inputValue[0] = map[string]any{
		"sda": []any{
			map[string]any{"disk_id": 12345},
		},
		"sdb": []any{
			map[string]any{"volume_id": 54321},
		},
	}

	result := expandDeviceMap(inputValue)

	if result.SDA.DiskID != 12345 {
		t.Fatal("disk id != 12345")
	}

	if result.SDB.VolumeID != 54321 {
		t.Fatal("volume id != 54321")
	}
}
