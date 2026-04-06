package instance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

			newDeviceMap, err := changeInstanceConfigDevice(*deviceMap, k, tDevice)
			if err != nil {
				return nil, fmt.Errorf("failed to change config device map: %w", err)
			}

			*deviceMap = newDeviceMap
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
		systemAlertsSet := v.(*schema.Set)
		systemAlerts := make([]int, 0)
		for _, alert := range systemAlertsSet.List() {
			systemAlerts = append(systemAlerts, alert.(int))
		}
		alertsACLPOpts.SystemAlerts = &systemAlerts
	}

	if v, ok := m["user_alerts"]; ok {
		userAlertsSet := v.(*schema.Set)
		userAlerts := make([]int, 0)
		for _, alert := range userAlertsSet.List() {
			userAlerts = append(userAlerts, alert.(int))
		}

		alertsACLPOpts.UserAlerts = &userAlerts
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
		systemAlertsSet := v.(*schema.Set)
		systemAlerts := make([]int, 0)
		for _, alert := range systemAlertsSet.List() {
			systemAlerts = append(systemAlerts, alert.(int))
		}

		alertsUpdateOpts.SystemAlerts = &systemAlerts
	}

	if v, ok := m["user_alerts"]; ok {
		userAlertsSet := v.(*schema.Set)
		userAlerts := make([]int, 0)
		for _, alert := range userAlertsSet.List() {
			userAlerts = append(userAlerts, alert.(int))
		}

		alertsUpdateOpts.UserAlerts = &userAlerts
	}

	return &alertsUpdateOpts
}
