package tmpl

import (
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label    string
	PubKey   string
	Type     string
	Image    string
	Group    string
	Tag      string
	Region   string
	RootPass string

	SwapSize int

	StackScriptName string

	Booted     bool
	ResizeDisk bool

	PlacementGroups []string
	AssignedGroup   string

	DiskEncryption *linodego.InstanceDiskEncryption
}

func Basic(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_basic", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func Updates(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_updates", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func WatchdogDisabled(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_watchdog_disabled", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func WithType(t testing.TB, label, pubKey, typ, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_type", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Type:     typ,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func WithSwapSize(t testing.TB, label, pubKey, region string, swapSize int, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_swap_size", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			SwapSize: swapSize,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func FullDisk(t testing.TB, label, pubKey, stackScriptName, region string, swapSize int, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_full_disk", TemplateData{
			Label:           label,
			PubKey:          pubKey,
			SwapSize:        swapSize,
			StackScriptName: stackScriptName,
			Image:           acceptance.TestImageLatest,
			Region:          region,
			RootPass:        rootPass,
		})
}

func WithConfig(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_config", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func MultipleConfigs(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_multiple_configs", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func Interfaces(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_interfaces", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func InterfacesUpdate(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_interfaces_update", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func InterfacesUpdateEmpty(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_interfaces_update_empty", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func ConfigInterfaces(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func ConfigInterfacesMultiple(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_multiple", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func ConfigInterfacesUpdate(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_update", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func ConfigInterfacesUpdateNoReboot(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_update_no_reboot", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func ConfigInterfacesUpdateEmpty(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_interfaces_update_empty", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func ConfigUpdates(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_updates", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func ConfigsAllUpdated(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_configs_all_updated", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func RawDisk(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_raw_disk", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func RawDiskDeleted(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_raw_disk_deleted", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func Tag(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_tag", TemplateData{
			Label:  label,
			Region: region,
			Image:  acceptance.TestImageLatest,
		})
}

func TagUpdate(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_tag_update", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func TagUpdateCaseChange(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_tag_update_case_change", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func TagVolume(t testing.TB, label, tag, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_tag_volume", TemplateData{
			Label:  label,
			Tag:    tag,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func RawDiskExpanded(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_raw_disk_expanded", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func Disk(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DiskMultiple(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_multiple", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DiskConfig(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DiskConfigExpanded(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_expanded", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DiskConfigResized(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_resized", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DiskConfigResizedExpanded(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_resized_expanded", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DiskConfigReordered(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_reordered", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DiskConfigMultiple(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_config_multiple", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DiskBootImage(t testing.TB, label, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_boot_image", TemplateData{
			Label:  label,
			Image:  image,
			Region: region,
		})
}

func VolumeConfig(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_volume_config", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func PrivateImage(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_private_image", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func NoImage(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_no_image", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func PrivateNetworking(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_private_networking", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func AuthorizedUsers(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_authorized_users", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func AuthorizedKeysEmpty(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_authorized_keys_empty", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func DiskAuthorizedKeysEmpty(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_authorized_keys_empty", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func StackScript(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_stackscript", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func DiskStackScript(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_stackscript", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func BootState(t testing.TB, label, region string, booted bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_boot_state", TemplateData{
			Label:  label,
			Image:  acceptance.TestImageLatest,
			Booted: booted,
			Region: region,
		})
}

func BootStateNoImage(t testing.TB, label, region string, booted bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_boot_state_noimage", TemplateData{
			Label:  label,
			Booted: booted,
			Region: region,
		})
}

func BootStateInterface(t testing.TB, label, region string, booted bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_boot_state_interface", TemplateData{
			Label:  label,
			Booted: booted,
			Image:  acceptance.TestImageLatest,
			Region: region,
		})
}

func BootStateConfig(t testing.TB, label, region string, booted bool, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_boot_state_config", TemplateData{
			Label:    label,
			Booted:   booted,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func TypeChangeDisk(t testing.TB, label, instanceType, region string, resizeDisk bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_type_change_disk", TemplateData{
			Label:      label,
			Type:       instanceType,
			Image:      acceptance.TestImageLatest,
			ResizeDisk: resizeDisk,
			Region:     region,
		})
}

func TypeChangeDiskExplicit(t testing.TB, label, instanceType, region string, resizeDisk bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_type_change_disk_explicit", TemplateData{
			Label:      label,
			Type:       instanceType,
			ResizeDisk: resizeDisk,
			Region:     region,
		})
}

func TypeChangeDiskNone(t testing.TB, label, instanceType, region string, resizeDisk bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_type_change_disk_none", TemplateData{
			Label:      label,
			Type:       instanceType,
			ResizeDisk: resizeDisk,
			Region:     region,
		})
}

func TypeChangeWarm(t testing.TB, label, instanceType, region string, resizeDisk bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_type_change_warm", TemplateData{
			Label:      label,
			Type:       instanceType,
			Image:      acceptance.TestImageLatest,
			ResizeDisk: resizeDisk,
			Region:     region,
		})
}

func IPv4Sharing(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ipv4_sharing", TemplateData{
			Label:  label,
			Region: region,
		})
}

func IPv4SharingEmpty(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ipv4_sharing_empty", TemplateData{
			Label:  label,
			Region: region,
		})
}

func IPv4SharingAllocation(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ipv4_sharing_allocation", TemplateData{
			Label:  label,
			Region: region,
		})
}

func IPv4SharingBadInput(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ipv4_sharing_bad_input", TemplateData{
			Label:  label,
			Region: region,
		})
}

func ManyLinodes(t testing.TB, label, pubKey, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_many_linodes", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			PubKey:   pubKey,
			Region:   region,
			RootPass: rootPass,
		})
}

func UserData(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_userdata", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DiskEncryption(
	t testing.TB,
	label,
	region,
	rootPass string,
	diskEncryption *linodego.InstanceDiskEncryption,
) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_encryption", TemplateData{
			Label:          label,
			Image:          acceptance.TestImageLatest,
			Region:         region,
			RootPass:       rootPass,
			DiskEncryption: diskEncryption,
		})
}

func DataBasic(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_basic", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DataWithBlockStorageEncryption(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_with_block_storage_encryption", TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		})
}

func DataWithPG(t testing.TB, label, region, assignedGroup string, groups []string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_with_pg", TemplateData{
			Label:           label,
			Region:          region,
			PlacementGroups: groups,
			AssignedGroup:   assignedGroup,
		})
}

func DataMultiple(t testing.TB, label, tag, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_multiple", TemplateData{
			Label:    label,
			Tag:      tag,
			Region:   region,
			Image:    acceptance.TestImageLatest,
			RootPass: rootPass,
		})
}

func DataMultipleOrder(t testing.TB, label, tag, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_multiple_order", TemplateData{
			Label:    label,
			Tag:      tag,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DataMultipleRegex(t testing.TB, label, tag, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_multiple_regex", TemplateData{
			Label:    label,
			Tag:      tag,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func DataClientFilter(t testing.TB, label, tag, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_data_clientfilter", TemplateData{
			Label:    label,
			Tag:      tag,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func FirewallOnCreation(t testing.TB, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_firewall_on_creation", TemplateData{
			Label:    label,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
}

func VPCInterface(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_vpc_interface", TemplateData{
			Label:  label,
			Region: region,
			Image:  acceptance.TestImageLatest,
		})
}

func VPCAndPublicInterfaces(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_vpc_public_interface", TemplateData{
			Label:  label,
			Region: region,
			Image:  acceptance.TestImageLatest,
		})
}

func PublicAndVPCInterfaces(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_public_vpc_interface", TemplateData{
			Label:  label,
			Region: region,
			Image:  acceptance.TestImageLatest,
		})
}

func PublicInterface(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_public_interface", TemplateData{
			Label:  label,
			Region: region,
			Image:  acceptance.TestImageLatest,
		})
}

func WithPG(t testing.TB, label, region, assignedGroup string, groups []string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_with_pg", TemplateData{
			Label:           label,
			Region:          region,
			PlacementGroups: groups,
			AssignedGroup:   assignedGroup,
		})
}

func WithReservedIP(t *testing.T, label, pubKey, region, rootPass string) string {
	generatedConfig := acceptance.ExecuteTemplate(t,
		"instance_with_reserved_ip", TemplateData{
			Label:    label,
			PubKey:   pubKey,
			Image:    acceptance.TestImageLatest,
			Region:   region,
			RootPass: rootPass,
		})
	return generatedConfig
}
