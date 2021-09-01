package instance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func expandStringList(list []interface{}) []string {
	slice := make([]string, 0, len(list))
	for _, s := range list {
		if val, ok := s.(string); ok && val != "" {
			slice = append(slice, val)
		}
	}
	return slice
}

func expandStringSet(set *schema.Set) []string {
	return expandStringList(set.List())
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

func expandIntSet(set *schema.Set) []int {
	return expandIntList(set.List())
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

func expandLinodeConfigInterface(i map[string]interface{}) linodego.InstanceConfigInterface {
	result := linodego.InstanceConfigInterface{}

	result.Label = i["label"].(string)
	result.Purpose = linodego.ConfigInterfacePurpose(i["purpose"].(string))
	result.IPAMAddress = i["ipam_address"].(string)

	return result
}
