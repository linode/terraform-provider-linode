package instanceconfig

import (
	"fmt"
	"testing"
)

func TestExpandDeviceMap(t *testing.T) {
	inputValue := []map[string]interface{}{
		{
			"sda": map[string]interface{}{
				"disk_id": 12345,
			},
			"sdb": map[string]interface{}{
				"volume_id": 54321,
			},
		},
	}

	result, err := expandDeviceMap(inputValue)
	if err != nil {
		t.Fatal(err)
	}

	if result.SDA.DiskID != 12345 {
		t.Fatal("disk id != 12345")
	}

	if result.SDB.VolumeID != 54321 {
		t.Fatal("volume id != 54321")
	}

	fmt.Println(result)
}
