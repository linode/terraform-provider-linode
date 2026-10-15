package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

// InstanceFilterModel describes the Terraform data model of the
// linode_instances data source.
type InstanceFilterModel struct {
	ID        types.String                     `tfsdk:"id"`
	Filters   frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order     types.String                     `tfsdk:"order"`
	OrderBy   types.String                     `tfsdk:"order_by"`
	Instances types.List                       `tfsdk:"instances"`
}

// The object types below are derived from the schema so the models and the
// schema can never drift apart.
var (
	instanceObjectType = listElemObjectType(instancesAttribute)

	instanceSpecsObjectType  = listElemObjectType(instanceSpecsAttribute)
	instanceAlertsObjectType = listElemObjectType(instanceAlertsAttribute)

	instanceBackupsObjectType  = listElemObjectType(instanceBackupsAttribute)
	instanceScheduleObjectType = listElemObjectType(childAttr(instanceBackupsAttribute, "schedule"))

	instanceConfigObjectType        = listElemObjectType(instanceConfigAttribute)
	instanceConfigHelpersObjectType = listElemObjectType(childAttr(instanceConfigAttribute, "helpers"))
	instanceConfigDevicesObjectType = listElemObjectType(childAttr(instanceConfigAttribute, "devices"))
	instanceConfigDeviceObjectType  = listElemObjectType(
		childAttr(childListAttr(instanceConfigAttribute, "devices"), "sda"),
	)

	instanceInterfaceObjectType     = listElemObjectType(instanceConfigInterfaceAttribute)
	instanceInterfaceIPv4ObjectType = listElemObjectType(childAttr(instanceConfigInterfaceAttribute, "ipv4"))
	instanceInterfaceIPv6ObjectType = listElemObjectType(childAttr(instanceConfigInterfaceAttribute, "ipv6"))
	instanceIPv6SLAACObjectType     = listElemObjectType(
		childAttr(childListAttr(instanceConfigInterfaceAttribute, "ipv6"), "slaac"),
	)
	instanceIPv6RangeObjectType = listElemObjectType(
		childAttr(childListAttr(instanceConfigInterfaceAttribute, "ipv6"), "range"),
	)

	instanceDiskObjectType           = listElemObjectType(instanceDiskAttribute)
	instancePlacementGroupObjectType = listElemObjectType(instancePlacementGroupAttribute)
)

// listElemObjectType returns the object type of the elements of the given
// list-nested attribute.
func listElemObjectType(a schema.Attribute) types.ObjectType {
	return a.GetType().(types.ListType).ElemType.(types.ObjectType)
}

// childAttr returns the named child of the given list-nested attribute.
func childAttr(a schema.ListNestedAttribute, name string) schema.Attribute {
	return a.NestedObject.Attributes[name]
}

// childListAttr is childAttr for children that are themselves list-nested.
func childListAttr(a schema.ListNestedAttribute, name string) schema.ListNestedAttribute {
	return a.NestedObject.Attributes[name].(schema.ListNestedAttribute)
}

// FlattenInstances populates this model with the given instances, making the
// additional API calls needed to fully describe each one.
func (model *InstanceFilterModel) FlattenInstances(
	ctx context.Context,
	client *linodego.Client,
	instances []linodego.Instance,
	diags *diag.Diagnostics,
) {
	elems := make([]attr.Value, len(instances))

	for i := range instances {
		elems[i] = flattenInstanceObject(ctx, client, &instances[i], diags)
		if diags.HasError() {
			return
		}
	}

	model.Instances = listValue(instanceObjectType, elems, diags)
}

// flattenInstanceObject fetches the extra details of a single Instance and
// converts it into the object described by instanceAttributes.
func flattenInstanceObject(
	ctx context.Context,
	client *linodego.Client,
	instance *linodego.Instance,
	diags *diag.Diagnostics,
) types.Object {
	null := types.ObjectNull(instanceObjectType.AttrTypes)
	id := instance.ID

	instanceNetwork, err := client.GetInstanceIPAddresses(ctx, id)
	if err != nil {
		diags.AddError("Failed to get IPs for Linode instance", err.Error())
		return null
	}

	instanceDisks, err := client.ListInstanceDisks(ctx, id, nil)
	if err != nil {
		diags.AddError("Failed to get the disks for the Linode instance", err.Error())
		return null
	}

	instanceConfigs, err := client.ListInstanceConfigs(ctx, id, nil)
	if err != nil {
		diags.AddError("Failed to get the configs for the Linode instance", err.Error())
		return null
	}

	return buildInstanceObject(instance, instanceNetwork, instanceDisks, instanceConfigs, diags)
}

// buildInstanceObject converts an Instance and its already-resolved disks,
// configs and IP addresses into the object described by instanceAttributes.
//
// Every field is populated with a known value, including the zero values that
// the former SDKv2 implementation wrote for absent data.
func buildInstanceObject(
	instance *linodego.Instance,
	instanceNetwork *linodego.InstanceIPAddressResponse,
	instanceDisks []linodego.InstanceDisk,
	instanceConfigs []linodego.InstanceConfig,
	diags *diag.Diagnostics,
) types.Object {
	var ipAddress, privateIPAddress string

	if instanceNetwork != nil && instanceNetwork.IPv4 != nil {
		if public := instanceNetwork.IPv4.Public; len(public) > 0 {
			ipAddress = public[0].Address
		}

		if private := instanceNetwork.IPv4.Private; len(private) > 0 {
			privateIPAddress = private[0].Address
		}
	}

	disks, swapSize := flattenInstanceDisksFramework(instanceDisks, diags)

	diskLabelIDMap := make(map[int]string, len(instanceDisks))
	for _, disk := range instanceDisks {
		diskLabelIDMap[disk.ID] = disk.Label
	}

	// boot_config_label is only populated when the instance has exactly one
	// config, matching the behavior of the former SDKv2 implementation.
	var bootConfigLabel string
	if len(instanceConfigs) == 1 {
		bootConfigLabel = instanceConfigs[0].Label
	}

	values := map[string]attr.Value{
		"id":                   types.Int64Value(int64(instance.ID)),
		"label":                types.StringValue(instance.Label),
		"status":               types.StringValue(string(instance.Status)),
		"type":                 types.StringValue(instance.Type),
		"region":               types.StringValue(instance.Region),
		"maintenance_policy":   types.StringValue(instance.MaintenancePolicy),
		"watchdog_enabled":     types.BoolValue(instance.WatchdogEnabled),
		"tags":                 stringSetValue(instance.Tags, diags),
		"capabilities":         stringSetValue(instance.Capabilities, diags),
		"locks":                stringSetValue(flattenInstanceLocks(instance.Locks), diags),
		"image":                types.StringValue(instance.Image),
		"interface_generation": types.StringValue(string(instance.InterfaceGeneration)),
		"host_uuid":            types.StringValue(instance.HostUUID),
		"has_user_data":        types.BoolValue(instance.HasUserData),
		"disk_encryption":      types.StringValue(string(instance.DiskEncryption)),
		"lke_cluster_id":       types.Int64Value(int64(instance.LKEClusterID)),
		"ipv4":                 stringSetValue(flattenInstanceIPv4(instance.IPv4), diags),
		"ipv6":                 types.StringValue(instance.IPv6),
		"ip_address":           types.StringValue(ipAddress),
		"private_ip_address":   types.StringValue(privateIPAddress),
		"swap_size":            types.Int64Value(int64(swapSize)),
		"boot_config_label":    types.StringValue(bootConfigLabel),
		"backups":              flattenInstanceBackupsFramework(instance.Backups, diags),
		"specs":                flattenInstanceSpecsFramework(instance.Specs, diags),
		"alerts":               flattenInstanceAlertsFramework(instance.Alerts, diags),
		"placement_group":      flattenInstancePlacementGroupFramework(instance.PlacementGroup, diags),
		"disk":                 disks,
		"config":               flattenInstanceConfigsFramework(instanceConfigs, diskLabelIDMap, diags),
	}

	return objectValue(instanceObjectType, values, diags)
}

func flattenInstanceLocks(locks []linodego.LockType) []string {
	return helper.MapSlice(locks, func(lock linodego.LockType) string {
		return string(lock)
	})
}

func flattenInstanceSpecsFramework(specs *linodego.InstanceSpec, diags *diag.Diagnostics) types.List {
	if specs == nil {
		specs = &linodego.InstanceSpec{}
	}

	return singletonList(instanceSpecsObjectType, objectValue(instanceSpecsObjectType, map[string]attr.Value{
		"vcpus":               types.Int64Value(int64(specs.VCPUs)),
		"disk":                types.Int64Value(int64(specs.Disk)),
		"memory":              types.Int64Value(int64(specs.Memory)),
		"transfer":            types.Int64Value(int64(specs.Transfer)),
		"accelerated_devices": types.Int64Value(int64(specs.AcceleratedDevices)),
		"gpus":                types.Int64Value(int64(specs.GPUs)),
	}, diags), diags)
}

func flattenInstanceAlertsFramework(alerts *linodego.InstanceAlert, diags *diag.Diagnostics) types.List {
	if alerts == nil {
		alerts = &linodego.InstanceAlert{}
	}

	return singletonList(instanceAlertsObjectType, objectValue(instanceAlertsObjectType, map[string]attr.Value{
		"cpu":            types.Int64Value(int64(alerts.CPU)),
		"io":             types.Int64Value(int64(alerts.IO)),
		"network_in":     types.Int64Value(int64(alerts.NetworkIn)),
		"network_out":    types.Int64Value(int64(alerts.NetworkOut)),
		"transfer_quota": types.Int64Value(int64(alerts.TransferQuota)),
	}, diags), diags)
}

func flattenInstanceBackupsFramework(backups *linodego.InstanceBackup, diags *diag.Diagnostics) types.List {
	if backups == nil {
		backups = &linodego.InstanceBackup{}
	}

	schedule := singletonList(
		instanceScheduleObjectType,
		objectValue(instanceScheduleObjectType, map[string]attr.Value{
			"day":    types.StringValue(backups.Schedule.Day),
			"window": types.StringValue(backups.Schedule.Window),
		}, diags),
		diags,
	)

	return singletonList(instanceBackupsObjectType, objectValue(instanceBackupsObjectType, map[string]attr.Value{
		"available": types.BoolValue(backups.Available),
		"enabled":   types.BoolValue(backups.Enabled),
		"schedule":  schedule,
	}, diags), diags)
}

func flattenInstancePlacementGroupFramework(
	pg *linodego.InstancePlacementGroup, diags *diag.Diagnostics,
) types.List {
	if pg == nil {
		return listValue(instancePlacementGroupObjectType, nil, diags)
	}

	return singletonList(
		instancePlacementGroupObjectType,
		objectValue(instancePlacementGroupObjectType, map[string]attr.Value{
			"id":                     types.Int64Value(int64(pg.ID)),
			"label":                  types.StringValue(pg.Label),
			"placement_group_type":   types.StringValue(string(pg.PlacementGroupType)),
			"placement_group_policy": types.StringValue(string(pg.PlacementGroupPolicy)),
		}, diags),
		diags,
	)
}

// flattenInstanceDisksFramework returns the disks of an Instance alongside the
// total size of its swap disks.
func flattenInstanceDisksFramework(
	instanceDisks []linodego.InstanceDisk, diags *diag.Diagnostics,
) (types.List, int) {
	var swapSize int

	elems := make([]attr.Value, len(instanceDisks))

	for i, disk := range instanceDisks {
		// Determine if swap exists and the size. If it does not exist, swap_size=0
		if disk.Filesystem == "swap" {
			swapSize += disk.Size
		}

		elems[i] = objectValue(instanceDiskObjectType, map[string]attr.Value{
			"id":         types.Int64Value(int64(disk.ID)),
			"size":       types.Int64Value(int64(disk.Size)),
			"label":      types.StringValue(disk.Label),
			"filesystem": types.StringValue(string(disk.Filesystem)),
		}, diags)
	}

	return listValue(instanceDiskObjectType, elems, diags), swapSize
}

func flattenInstanceConfigsFramework(
	instanceConfigs []linodego.InstanceConfig,
	diskLabelIDMap map[int]string,
	diags *diag.Diagnostics,
) types.List {
	elems := make([]attr.Value, len(instanceConfigs))

	for i, config := range instanceConfigs {
		var helpers linodego.InstanceConfigHelpers
		if config.Helpers != nil {
			helpers = *config.Helpers
		}

		helpersValue := singletonList(
			instanceConfigHelpersObjectType,
			objectValue(instanceConfigHelpersObjectType, map[string]attr.Value{
				"updatedb_disabled":  types.BoolValue(helpers.UpdateDBDisabled),
				"distro":             types.BoolValue(helpers.Distro),
				"modules_dep":        types.BoolValue(helpers.ModulesDep),
				"network":            types.BoolValue(helpers.Network),
				"devtmpfs_automount": types.BoolValue(helpers.DevTmpFsAutomount),
			}, diags),
			diags,
		)

		elems[i] = objectValue(instanceConfigObjectType, map[string]attr.Value{
			"id":           types.Int64Value(int64(config.ID)),
			"label":        types.StringValue(config.Label),
			"comments":     types.StringValue(config.Comments),
			"kernel":       types.StringValue(config.Kernel),
			"root_device":  types.StringValue(config.RootDevice),
			"run_level":    types.StringValue(config.RunLevel),
			"virt_mode":    types.StringValue(config.VirtMode),
			"memory_limit": types.Int64Value(int64(config.MemoryLimit)),
			"helpers":      helpersValue,
			"devices":      flattenInstanceConfigDevicesFramework(config.Devices, diskLabelIDMap, diags),
			"interface":    flattenInstanceInterfacesFramework(config.Interfaces, diags),
		}, diags)
	}

	return listValue(instanceConfigObjectType, elems, diags)
}

// flattenInstanceConfigDevicesFramework converts a config's device map into the
// single-element list of device slots described by the schema.
func flattenInstanceConfigDevicesFramework(
	deviceMap *linodego.InstanceConfigDeviceMap,
	diskLabelIDMap map[int]string,
	diags *diag.Diagnostics,
) types.List {
	var devices []*linodego.InstanceConfigDevice
	if deviceMap != nil {
		devices = configDeviceSlice(*deviceMap)
	}

	keys := helper.GetConfigDeviceKeys()
	values := make(map[string]attr.Value, len(keys))

	for i, key := range keys {
		var device *linodego.InstanceConfigDevice
		if i < len(devices) {
			device = devices[i]
		}

		values[key] = flattenInstanceConfigDeviceFramework(device, diskLabelIDMap, diags)
	}

	return singletonList(
		instanceConfigDevicesObjectType,
		objectValue(instanceConfigDevicesObjectType, values, diags),
		diags,
	)
}

// flattenInstanceConfigDeviceFramework converts a single device slot. An
// unassigned slot flattens to an empty list; the fields that do not apply to
// the assigned device type keep their zero values, as they did under SDKv2.
func flattenInstanceConfigDeviceFramework(
	device *linodego.InstanceConfigDevice,
	diskLabelIDMap map[int]string,
	diags *diag.Diagnostics,
) types.List {
	if device == nil || emptyInstanceConfigDevice(*device) {
		return listValue(instanceConfigDeviceObjectType, nil, diags)
	}

	values := map[string]attr.Value{
		"disk_id":    types.Int64Value(0),
		"disk_label": types.StringValue(""),
		"volume_id":  types.Int64Value(0),
	}

	if device.DiskID > 0 {
		values["disk_id"] = types.Int64Value(int64(device.DiskID))

		if label, found := diskLabelIDMap[device.DiskID]; found {
			values["disk_label"] = types.StringValue(label)
		}
	} else {
		values["volume_id"] = types.Int64Value(int64(device.VolumeID))
	}

	return singletonList(
		instanceConfigDeviceObjectType,
		objectValue(instanceConfigDeviceObjectType, values, diags),
		diags,
	)
}

func flattenInstanceInterfacesFramework(
	interfaces []linodego.InstanceConfigInterface, diags *diag.Diagnostics,
) types.List {
	elems := make([]attr.Value, len(interfaces))

	for i, iface := range interfaces {
		elems[i] = objectValue(instanceInterfaceObjectType, map[string]attr.Value{
			"id":           types.Int64Value(int64(iface.ID)),
			"purpose":      types.StringValue(string(iface.Purpose)),
			"ipam_address": types.StringValue(iface.IPAMAddress),
			"label":        types.StringValue(iface.Label),
			"primary":      types.BoolValue(iface.Primary),
			"active":       types.BoolValue(iface.Active),
			"vpc_id":       int64PointerValue(iface.VPCID),
			"subnet_id":    int64PointerValue(iface.SubnetID),
			"ip_ranges":    stringListValue(iface.IPRanges, diags),
			"ipv4":         flattenInterfaceIPv4Framework(iface.IPv4, diags),
			"ipv6":         flattenInterfaceIPv6Framework(iface.IPv6, diags),
		}, diags)
	}

	return listValue(instanceInterfaceObjectType, elems, diags)
}

func flattenInterfaceIPv4Framework(ipv4 *linodego.VPCIPv4, diags *diag.Diagnostics) types.List {
	if ipv4 == nil {
		return listValue(instanceInterfaceIPv4ObjectType, nil, diags)
	}

	var nat1To1 string
	if ipv4.NAT1To1 != nil {
		nat1To1 = *ipv4.NAT1To1
	}

	return singletonList(
		instanceInterfaceIPv4ObjectType,
		objectValue(instanceInterfaceIPv4ObjectType, map[string]attr.Value{
			"vpc":     types.StringValue(ipv4.VPC),
			"nat_1_1": types.StringValue(nat1To1),
		}, diags),
		diags,
	)
}

func flattenInterfaceIPv6Framework(
	ipv6 *linodego.InstanceConfigInterfaceIPv6, diags *diag.Diagnostics,
) types.List {
	if ipv6 == nil {
		return listValue(instanceInterfaceIPv6ObjectType, nil, diags)
	}

	slaac := make([]attr.Value, len(ipv6.SLAAC))
	for i, entry := range ipv6.SLAAC {
		slaac[i] = objectValue(instanceIPv6SLAACObjectType, map[string]attr.Value{
			"range":          types.StringValue(entry.Range),
			"assigned_range": types.StringValue(entry.Range),
			"address":        types.StringValue(entry.Address),
		}, diags)
	}

	ranges := make([]attr.Value, len(ipv6.Ranges))
	for i, entry := range ipv6.Ranges {
		ranges[i] = objectValue(instanceIPv6RangeObjectType, map[string]attr.Value{
			"range":          types.StringValue(entry.Range),
			"assigned_range": types.StringValue(entry.Range),
		}, diags)
	}

	var isPublic bool
	if ipv6.IsPublic != nil {
		isPublic = *ipv6.IsPublic
	}

	return singletonList(
		instanceInterfaceIPv6ObjectType,
		objectValue(instanceInterfaceIPv6ObjectType, map[string]attr.Value{
			"slaac":     listValue(instanceIPv6SLAACObjectType, slaac, diags),
			"range":     listValue(instanceIPv6RangeObjectType, ranges, diags),
			"is_public": types.BoolValue(isPublic),
		}, diags),
		diags,
	)
}

// The helpers below keep the flatten functions above free of repeated
// diagnostics plumbing.

func objectValue(
	objectType types.ObjectType, values map[string]attr.Value, diags *diag.Diagnostics,
) types.Object {
	result, d := types.ObjectValue(objectType.AttrTypes, values)
	diags.Append(d...)

	return result
}

func listValue(elemType attr.Type, elems []attr.Value, diags *diag.Diagnostics) types.List {
	if elems == nil {
		elems = []attr.Value{}
	}

	result, d := types.ListValue(elemType, elems)
	diags.Append(d...)

	return result
}

func singletonList(elemType attr.Type, elem attr.Value, diags *diag.Diagnostics) types.List {
	return listValue(elemType, []attr.Value{elem}, diags)
}

func stringSetValue(values []string, diags *diag.Diagnostics) types.Set {
	result, d := types.SetValue(types.StringType, helper.StringSliceToFrameworkValueSlice(values))
	diags.Append(d...)

	return result
}

func stringListValue(values []string, diags *diag.Diagnostics) types.List {
	result, d := types.ListValue(types.StringType, helper.StringSliceToFrameworkValueSlice(values))
	diags.Append(d...)

	return result
}

func int64PointerValue(value *int) types.Int64 {
	if value == nil {
		return types.Int64Value(0)
	}

	return types.Int64Value(int64(*value))
}
