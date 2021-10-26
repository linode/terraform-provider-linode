package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	PubKey string
	Type   string
	Image  string
	Group  string

	SwapSize int

	StackScriptName string
}

func Basic(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_basic", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func Updates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_updates", TemplateData{
			Label: label,
		})
}

func WatchdogDisabled(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_watchdog_disabled", TemplateData{
			Label: label,
		})
}

func WithType(t *testing.T, label, pubKey, typ string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_type", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Type:   typ,
		})
}

func WithSwapSize(t *testing.T, label, pubKey string, swapSize int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_swap_size", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			SwapSize: swapSize,
		})
}

func FullDisk(t *testing.T, label, pubKey, stackScriptName string, swapSize int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_full_disk", TemplateData{
			Label:           label,
			PubKey:          pubKey,
			SwapSize:        swapSize,
			StackScriptName: stackScriptName,
		})
}

func WithConfig(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_config", TemplateData{
			Label: label,
		})
}

func MultipleConfigs(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_multiple_configs", TemplateData{
			Label: label,
		})
}

func Interfaces(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_interfaces", TemplateData{
			Label: label,
		})
}

func InterfacesUpdate(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_interfaces_update", TemplateData{
			Label: label,
		})
}

func InterfacesUpdateEmpty(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_interfaces_update_empty", TemplateData{
			Label: label,
		})
}

func ConfigInterfaces(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces", TemplateData{
			Label: label,
		})
}

func ConfigInterfacesMultiple(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_multiple", TemplateData{
			Label: label,
		})
}

func ConfigInterfacesUpdate(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_update", TemplateData{
			Label: label,
		})
}

func ConfigInterfacesUpdateEmpty(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_update_empty", TemplateData{
			Label: label,
		})
}

func ConfigUpdates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_updates", TemplateData{
			Label: label,
		})
}

func ConfigsAllUpdated(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_configs_all_updated", TemplateData{
			Label: label,
		})
}

func RawDisk(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_raw_disk", TemplateData{
			Label: label,
		})
}

func RawDiskDeleted(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_raw_disk_deleted", TemplateData{
			Label: label,
		})
}

func Tag(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_tag", TemplateData{
			Label: label,
		})
}

func TagUpdate(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_tag_update", TemplateData{
			Label: label,
		})
}

func RawDiskExpanded(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_raw_disk_expanded", TemplateData{
			Label: label,
		})
}

func Disk(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DiskMultiple(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_multiple", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DiskMultipleReadOnly(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_multiple_readonly", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DiskConfig(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DiskConfigExpanded(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_expanded", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DiskConfigResized(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_resized", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DiskConfigResizedExpanded(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_resized_expanded", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DiskConfigReordered(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_reordered", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DiskConfigMultiple(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_multiple", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DiskBootImage(t *testing.T, label, image string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_boot_image", TemplateData{
			Label: label,
			Image: image,
		})
}

func VolumeConfig(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_volume_config", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func PrivateImage(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_private_image", TemplateData{
			Label: label,
		})
}

func NoImage(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_no_image", TemplateData{
			Label: label,
		})
}

func PrivateNetworking(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_private_networking", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func AuthorizedUsers(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_authorized_users", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func StackScript(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_stackscript", TemplateData{
			Label: label,
		})
}

func DiskStackScript(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_stackscript", TemplateData{
			Label:  label,
			PubKey: pubKey,
		})
}

func DataBasic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_basic", TemplateData{
			Label: label,
		})
}

func DataMultiple(t *testing.T, label, group string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_multiple", TemplateData{
			Label: label,
			Group: group,
		})
}
