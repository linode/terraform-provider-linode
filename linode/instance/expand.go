package instance

import (
	"github.com/linode/linodego"
)

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

func expandInstanceACLPAlertsOpts(m map[string]interface{}) *linodego.InstanceACLPAlertsOptions {
	var alertsACLPOpts linodego.InstanceACLPAlertsOptions

	if v, ok := m["system_alerts"]; ok {
		l := v.([]interface{})
		systemAlerts := make([]int, 0, len(l))
		for _, raw := range l {
			systemAlerts = append(systemAlerts, raw.(int))
		}
		alertsACLPOpts.SystemAlerts = systemAlerts
	}

	if v, ok := m["user_alerts"]; ok {
		l := v.([]interface{})
		userAlerts := make([]int, 0, len(l))
		for _, raw := range l {
			userAlerts = append(userAlerts, raw.(int))
		}
		alertsACLPOpts.UserAlerts = userAlerts
	}

	return &alertsACLPOpts
}

func expandInstanceAlertsUpdateOpts(m map[string]interface{}) *linodego.InstanceAlert {
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
		l := v.([]interface{})
		systemAlerts := make([]int, 0, len(l))
		for _, raw := range l {
			systemAlerts = append(systemAlerts, raw.(int))
		}
		alertsUpdateOpts.SystemAlerts = systemAlerts
	}

	if v, ok := m["user_alerts"]; ok {
		l := v.([]interface{})
		userAlerts := make([]int, 0, len(l))
		for _, raw := range l {
			userAlerts = append(userAlerts, raw.(int))
		}
		alertsUpdateOpts.UserAlerts = userAlerts
	}

	return &alertsUpdateOpts
}
