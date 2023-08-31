//go:build unit

package instanceconfig

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExpandDevicesNamedBlock(t *testing.T) {
	inputValue := make([]any, 1)
	inputValue[0] = map[string]any{
		"sda": []any{
			map[string]any{"disk_id": 12345},
		},
		"sdb": []any{
			map[string]any{"volume_id": 54321},
		},
	}

	result := expandDevicesNamedBlock(inputValue)

	if result.SDA.DiskID != 12345 {
		t.Fatal("disk id != 12345")
	}

	if result.SDB.VolumeID != 54321 {
		t.Fatal("volume id != 54321")
	}
}

func TestExpandDevicesBlock(t *testing.T) {
	inputValue := make([]any, 2)
	inputValue[0] = map[string]any{
		"device_name": "sda",
		"disk_id":     12345,
	}
	inputValue[1] = map[string]any{
		"device_name": "sdb",
		"volume_id":   54321,
	}

	// Hack to create a *schema.Set hashing on device_name
	setValue := schema.NewSet(func(i any) int {
		return schema.HashString(i.(map[string]any)["device_name"])
	}, inputValue)

	result := expandDevicesBlock(setValue)

	if result.SDA.DiskID != 12345 {
		t.Fatal("disk id != 12345")
	}

	if result.SDB.VolumeID != 54321 {
		t.Fatal("volume id != 54321")
	}
}
