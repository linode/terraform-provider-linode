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
	Tag    string

	SwapSize int

	StackScriptName string

	Booted bool
}

func Basic(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_basic", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func Updates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_updates", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func WatchdogDisabled(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_watchdog_disabled", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func WithType(t *testing.T, label, pubKey, typ string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_type", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Type:   typ,
			Image:  acceptance.TestImageLatest,
		})
}

func WithSwapSize(t *testing.T, label, pubKey string, swapSize int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_swap_size", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			SwapSize: swapSize,
			Image:    acceptance.TestImageLatest,
		})
}

func FullDisk(t *testing.T, label, pubKey, stackScriptName string, swapSize int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_full_disk", TemplateData{
			Label:           label,
			PubKey:          pubKey,
			SwapSize:        swapSize,
			StackScriptName: stackScriptName,
			Image:           acceptance.TestImageLatest,
		})
}

func WithConfig(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_config", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func MultipleConfigs(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_multiple_configs", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func Interfaces(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_interfaces", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func InterfacesUpdate(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_interfaces_update", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func InterfacesUpdateEmpty(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_interfaces_update_empty", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func ConfigInterfaces(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func ConfigInterfacesMultiple(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_multiple", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func ConfigInterfacesUpdate(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_update", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func ConfigInterfacesUpdateEmpty(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_update_empty", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func ConfigUpdates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_updates", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func ConfigsAllUpdated(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_configs_all_updated", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func RawDisk(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_raw_disk", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func RawDiskDeleted(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_raw_disk_deleted", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func Tag(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_tag", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func TagUpdate(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_tag_update", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func TagVolume(t *testing.T, label, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_tag_volume", TemplateData{
			Label: label,
			Tag:   tag,
			Image: acceptance.TestImageLatest,
		})
}

func RawDiskExpanded(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_raw_disk_expanded", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func Disk(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func DiskMultiple(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_multiple", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func DiskConfig(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func DiskConfigExpanded(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_expanded", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func DiskConfigResized(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_resized", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func DiskConfigResizedExpanded(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_resized_expanded", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func DiskConfigReordered(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_reordered", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func DiskConfigMultiple(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_multiple", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
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
			Image:  acceptance.TestImageLatest,
		})
}

func PrivateImage(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_private_image", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func NoImage(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_no_image", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func PrivateNetworking(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_private_networking", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func AuthorizedUsers(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_authorized_users", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func StackScript(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_stackscript", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func DiskStackScript(t *testing.T, label, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_stackscript", TemplateData{
			Label:  label,
			PubKey: pubKey,
			Image:  acceptance.TestImageLatest,
		})
}

func BootState(t *testing.T, label string, booted bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_boot_state", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Booted: booted,
		})
}

func BootStateNoImage(t *testing.T, label string, booted bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_boot_state_noimage", TemplateData{
			Label:  label,
			Booted: booted,
		})
}

func DataBasic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_basic", TemplateData{
			Label: label,
			Image: acceptance.TestImageLatest,
		})
}

func DataMultiple(t *testing.T, label, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_multiple", TemplateData{
			Label: label,
			Tag:   tag,
			Image: acceptance.TestImageLatest,
		})
}

func DataMultipleOrder(t *testing.T, label, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_multiple_order", TemplateData{
			Label: label,
			Tag:   tag,
			Image: acceptance.TestImageLatest,
		})
}

func DataMultipleRegex(t *testing.T, label, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_multiple_regex", TemplateData{
			Label: label,
			Tag:   tag,
			Image: acceptance.TestImageLatest,
		})
}

func DataClientFilter(t *testing.T, label, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_clientfilter", TemplateData{
			Label: label,
			Tag:   tag,
			Image: acceptance.TestImageLatest,
		})
}
