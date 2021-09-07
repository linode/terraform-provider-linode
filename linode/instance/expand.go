package instance

import (
	"github.com/linode/linodego"
)

// expandInstanceConfigDeviceMap converts a terraform linode_instance config.*.devices map to a InstanceConfigDeviceMap
// for the Linode API.
func expandInstanceConfigDeviceMap(
	m map[string]interface{}, diskIDLabelMap map[string]int) (deviceMap *linodego.InstanceConfigDeviceMap, err error) {
	if len(m) == 0 {
		return nil, nil
	}
	deviceMap = &linodego.InstanceConfigDeviceMap{}
	for k, rdev := range m {
		devSlots := rdev.([]interface{})
		for _, rrdev := range devSlots {
			dev := rrdev.(map[string]interface{})
			tDevice := new(linodego.InstanceConfigDevice)
			if err := assignConfigDevice(tDevice, dev, diskIDLabelMap); err != nil {
				return nil, err
			}

			*deviceMap = changeInstanceConfigDevice(*deviceMap, k, tDevice)
		}
	}
	return deviceMap, nil
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

func expandConfigInterface(i map[string]interface{}) linodego.InstanceConfigInterface {
	result := linodego.InstanceConfigInterface{}

	result.Label = i["label"].(string)
	result.Purpose = linodego.ConfigInterfacePurpose(i["purpose"].(string))
	result.IPAMAddress = i["ipam_address"].(string)

	return result
}

func expandStringList(list []interface{}) []string {
	slice := make([]string, 0, len(list))
	for _, s := range list {
		if val, ok := s.(string); ok && val != "" {
			slice = append(slice, val)
		}
	}
	return slice
}

func expandIntList(list []interface{}) []int {
	slice := make([]int, 0, len(list))
	for _, n := range list {
		if val, ok := n.(int); ok {
			slice = append(slice, val)
		}
	}
	return slice
}
