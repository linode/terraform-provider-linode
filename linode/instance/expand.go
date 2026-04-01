package instance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

// expandIntSet converts a value that is either a *schema.Set or []any of ints
// into a []int slice.
func expandIntSet(v any) []int {
	var items []any
	switch val := v.(type) {
	case *schema.Set:
		items = val.List()
	case []any:
		items = val
	default:
		return nil
	}
	result := make([]int, 0, len(items))
	for _, raw := range items {
		result = append(result, raw.(int))
	}
	return result
}

// expandInstanceConfigDeviceMap converts a terraform linode_instance config.*.devices map to a InstanceConfigDeviceMap
// for the Linode API.
func expandInstanceConfigDeviceMap(
	m map[string]any, diskIDLabelMap map[string]int,
) (deviceMap *linodego.InstanceConfigDeviceMap, err error) {
	if len(m) == 0 {
		return nil, nil
	}
	deviceMap = &linodego.InstanceConfigDeviceMap{}
	for k, rdev := range m {
		devSlots := rdev.([]any)
		for _, rrdev := range devSlots {
			dev := rrdev.(map[string]any)
			tDevice := new(linodego.InstanceConfigDevice)
			if err := assignConfigDevice(tDevice, dev, diskIDLabelMap); err != nil {
				return nil, err
			}

			*deviceMap = changeInstanceConfigDevice(*deviceMap, k, tDevice)
		}
	}
	return deviceMap, nil
}

func expandInstanceConfigDevice(m map[string]any) *linodego.InstanceConfigDevice {
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

func expandInstanceACLPAlertsOpts(m map[string]any) *linodego.InstanceACLPAlertsOptions {
	var alertsACLPOpts linodego.InstanceACLPAlertsOptions

	if v, ok := m["system_alerts"]; ok {
		l := expandIntSet(v)
		alertsACLPOpts.SystemAlerts = l
	}

	if v, ok := m["user_alerts"]; ok {
		l := expandIntSet(v)
		alertsACLPOpts.UserAlerts = l
	}

	return &alertsACLPOpts
}

func expandInstanceAlertsUpdateOpts(m map[string]any) *linodego.InstanceAlert {
	var alertsUpdateOpts linodego.InstanceAlert

	// TODO(displague) only set specified alerts
	if v, ok := m["cpu"]; ok {
		alertsUpdateOpts.CPU = v.(int)
	}
	if v, ok := m["io"]; ok {
		alertsUpdateOpts.IO = v.(int)
	}
	if v, ok := m["network_in"]; ok {
		alertsUpdateOpts.NetworkIn = v.(int)
	}
	if v, ok := m["network_out"]; ok {
		alertsUpdateOpts.NetworkOut = v.(int)
	}
	if v, ok := m["transfer_quota"]; ok {
		alertsUpdateOpts.TransferQuota = v.(int)
	}

	if v, ok := m["system_alerts"]; ok {
		l := expandIntSet(v)
		alertsUpdateOpts.SystemAlerts = l
	}

	if v, ok := m["user_alerts"]; ok {
		l := expandIntSet(v)
		alertsUpdateOpts.UserAlerts = l
	}

	return &alertsUpdateOpts
}
