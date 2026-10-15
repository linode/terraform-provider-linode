//go:build unit

package instance

import (
	"fmt"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertNoNullValues walks the given value and fails if any null or unknown
// value is found.
//
// The former SDKv2 data source never wrote a null into state: absent scalars
// became zero values and absent collections became empty. Reproducing that
// matters because expressions such as length(...instances[0].disk) fail on a
// null collection.
func assertNoNullValues(t *testing.T, path string, value attr.Value) {
	t.Helper()

	require.Falsef(t, value.IsNull(), "%s is null", path)
	require.Falsef(t, value.IsUnknown(), "%s is unknown", path)

	switch typed := value.(type) {
	case types.Object:
		for name, child := range typed.Attributes() {
			assertNoNullValues(t, path+"."+name, child)
		}
	case types.List:
		for i, child := range typed.Elements() {
			assertNoNullValues(t, fmt.Sprintf("%s[%d]", path, i), child)
		}
	case types.Set:
		for i, child := range typed.Elements() {
			assertNoNullValues(t, fmt.Sprintf("%s[%d]", path, i), child)
		}
	}
}

func objectAttr(t *testing.T, object types.Object, name string) attr.Value {
	t.Helper()

	value, ok := object.Attributes()[name]
	require.Truef(t, ok, "attribute %q is missing", name)

	return value
}

// listElem returns the object at the given index of a list-typed value.
func listElem(t *testing.T, value attr.Value, index int) types.Object {
	t.Helper()

	list, ok := value.(types.List)
	require.True(t, ok, "value is not a list")
	require.Greater(t, len(list.Elements()), index)

	object, ok := list.Elements()[index].(types.Object)
	require.True(t, ok, "list element is not an object")

	return object
}

func listLen(t *testing.T, value attr.Value) int {
	t.Helper()

	list, ok := value.(types.List)
	require.True(t, ok, "value is not a list")

	return len(list.Elements())
}

func testInstance() *linodego.Instance {
	return &linodego.Instance{
		ID:                  123,
		Label:               "test-instance",
		Status:              linodego.InstanceRunning,
		Type:                "g6-nanode-1",
		Region:              "us-east",
		MaintenancePolicy:   "linode/migrate",
		WatchdogEnabled:     true,
		Tags:                []string{"foo", "bar"},
		Capabilities:        []string{"SMTP Enabled"},
		Locks:               []linodego.LockType{linodego.LockTypeCannotDelete},
		Image:               "linode/ubuntu22.04",
		InterfaceGeneration: linodego.GenerationLinode,
		HostUUID:            "abc123",
		HasUserData:         true,
		DiskEncryption:      linodego.InstanceDiskEncryptionEnabled,
		LKEClusterID:        321,
		IPv4:                []net.IP{net.ParseIP("192.0.2.1")},
		IPv6:                "2600:3c03::/64",
		Specs:               &linodego.InstanceSpec{Disk: 25600, Memory: 1024, VCPUs: 1, Transfer: 1000},
		Alerts:              &linodego.InstanceAlert{CPU: 90, IO: 10000, NetworkIn: 10},
		Backups:             &linodego.InstanceBackup{Available: true, Enabled: true},
		PlacementGroup: &linodego.InstancePlacementGroup{
			ID:                   1,
			Label:                "pg",
			PlacementGroupType:   linodego.PlacementGroupTypeAntiAffinityLocal,
			PlacementGroupPolicy: linodego.PlacementGroupPolicyStrict,
		},
	}
}

func testDisks() []linodego.InstanceDisk {
	return []linodego.InstanceDisk{
		{ID: 1000, Label: "boot", Size: 25088, Filesystem: linodego.FilesystemExt4},
		{ID: 1001, Label: "swap", Size: 512, Filesystem: linodego.FilesystemSwap},
	}
}

func testConfigs() []linodego.InstanceConfig {
	isPublic := true
	vpcID := 456
	subnetID := 789
	nat := "192.0.2.10"

	return []linodego.InstanceConfig{
		{
			ID:          2000,
			Label:       "My Config",
			Comments:    "a comment",
			Kernel:      "linode/latest-64bit",
			RootDevice:  "/dev/sda",
			RunLevel:    "default",
			VirtMode:    "paravirt",
			MemoryLimit: 2048,
			Helpers:     &linodego.InstanceConfigHelpers{Distro: true, Network: true},
			Devices: &linodego.InstanceConfigDeviceMap{
				SDA: &linodego.InstanceConfigDevice{DiskID: 1000},
				SDB: &linodego.InstanceConfigDevice{VolumeID: 3000},
				// An assigned disk that is not present in diskLabelIDMap.
				SDC: &linodego.InstanceConfigDevice{DiskID: 9999},
			},
			Interfaces: []linodego.InstanceConfigInterface{
				{
					ID:          10,
					Purpose:     linodego.InterfacePurposeVLAN,
					Label:       "my-vlan",
					IPAMAddress: "10.0.0.1/24",
				},
				{
					ID:       11,
					Purpose:  linodego.InterfacePurposeVPC,
					Primary:  true,
					Active:   true,
					VPCID:    &vpcID,
					SubnetID: &subnetID,
					IPRanges: []string{"192.0.2.0/24"},
					IPv4:     &linodego.VPCIPv4{VPC: "10.0.0.2", NAT1To1: &nat},
					IPv6: &linodego.InstanceConfigInterfaceIPv6{
						SLAAC:    []linodego.InstanceConfigInterfaceIPv6SLAAC{{Range: "2600::/64", Address: "2600::1"}},
						Ranges:   []linodego.InstanceConfigInterfaceIPv6Range{{Range: "2600:1::/64"}},
						IsPublic: &isPublic,
					},
				},
			},
		},
	}
}

func testNetwork() *linodego.InstanceIPAddressResponse {
	return &linodego.InstanceIPAddressResponse{
		IPv4: &linodego.InstanceIPv4Response{
			Public:  []linodego.InstanceIP{{Address: "192.0.2.1"}},
			Private: []linodego.InstanceIP{{Address: "192.168.128.1"}},
		},
	}
}

func buildTestObject(
	t *testing.T,
	instance *linodego.Instance,
	network *linodego.InstanceIPAddressResponse,
	disks []linodego.InstanceDisk,
	configs []linodego.InstanceConfig,
) types.Object {
	t.Helper()

	var diags diag.Diagnostics

	result := buildInstanceObject(instance, network, disks, configs, &diags)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags.Errors())

	return result
}

func TestBuildInstanceObject_populated(t *testing.T) {
	object := buildTestObject(t, testInstance(), testNetwork(), testDisks(), testConfigs())

	assertNoNullValues(t, "instance", object)

	assert.Equal(t, types.Int64Value(123), objectAttr(t, object, "id"))
	assert.Equal(t, types.StringValue("running"), objectAttr(t, object, "status"))
	assert.Equal(t, types.StringValue("enabled"), objectAttr(t, object, "disk_encryption"))
	assert.Equal(t, types.Int64Value(321), objectAttr(t, object, "lke_cluster_id"))
	assert.Equal(t, types.StringValue("192.0.2.1"), objectAttr(t, object, "ip_address"))
	assert.Equal(t, types.StringValue("192.168.128.1"), objectAttr(t, object, "private_ip_address"))

	// swap_size is the total size of the swap disks.
	assert.Equal(t, types.Int64Value(512), objectAttr(t, object, "swap_size"))

	// boot_config_label is only populated for a single-config instance.
	assert.Equal(t, types.StringValue("My Config"), objectAttr(t, object, "boot_config_label"))

	// locks are exposed as plain strings.
	locks, ok := objectAttr(t, object, "locks").(types.Set)
	require.True(t, ok)
	assert.Equal(t, []attr.Value{types.StringValue("cannot_delete")}, locks.Elements())

	// The single-element wrappers SDKv2 produced are preserved.
	assert.Equal(t, 1, listLen(t, objectAttr(t, object, "specs")))
	assert.Equal(t, 1, listLen(t, objectAttr(t, object, "alerts")))
	assert.Equal(t, 1, listLen(t, objectAttr(t, object, "backups")))
	assert.Equal(t, 1, listLen(t, objectAttr(t, object, "placement_group")))
	assert.Equal(t, 2, listLen(t, objectAttr(t, object, "disk")))
	assert.Equal(t, 1, listLen(t, objectAttr(t, object, "config")))

	backups := listElem(t, objectAttr(t, object, "backups"), 0)
	assert.Equal(t, 1, listLen(t, objectAttr(t, backups, "schedule")))

	config := listElem(t, objectAttr(t, object, "config"), 0)
	assert.Equal(t, types.Int64Value(2000), objectAttr(t, config, "id"))
	assert.Equal(t, 1, listLen(t, objectAttr(t, config, "helpers")))
	assert.Equal(t, 1, listLen(t, objectAttr(t, config, "devices")))
	assert.Equal(t, 2, listLen(t, objectAttr(t, config, "interface")))
}

func TestBuildInstanceObject_configDevices(t *testing.T) {
	object := buildTestObject(t, testInstance(), testNetwork(), testDisks(), testConfigs())

	config := listElem(t, objectAttr(t, object, "config"), 0)
	devices := listElem(t, objectAttr(t, config, "devices"), 0)

	// Every valid device slot is present.
	keys := helper.GetConfigDeviceKeys()
	assert.Len(t, devices.Attributes(), len(keys))

	for _, key := range keys {
		value := objectAttr(t, devices, key)
		require.Falsef(t, value.IsNull(), "device slot %q is null", key)
	}

	// A disk device carries disk_id and disk_label; volume_id stays at zero.
	sda := listElem(t, objectAttr(t, devices, "sda"), 0)
	assert.Equal(t, types.Int64Value(1000), objectAttr(t, sda, "disk_id"))
	assert.Equal(t, types.StringValue("boot"), objectAttr(t, sda, "disk_label"))
	assert.Equal(t, types.Int64Value(0), objectAttr(t, sda, "volume_id"))

	// A volume device carries volume_id; disk fields stay at their zero values.
	sdb := listElem(t, objectAttr(t, devices, "sdb"), 0)
	assert.Equal(t, types.Int64Value(3000), objectAttr(t, sdb, "volume_id"))
	assert.Equal(t, types.Int64Value(0), objectAttr(t, sdb, "disk_id"))
	assert.Equal(t, types.StringValue(""), objectAttr(t, sdb, "disk_label"))

	// A disk device whose label cannot be resolved keeps an empty disk_label.
	sdc := listElem(t, objectAttr(t, devices, "sdc"), 0)
	assert.Equal(t, types.Int64Value(9999), objectAttr(t, sdc, "disk_id"))
	assert.Equal(t, types.StringValue(""), objectAttr(t, sdc, "disk_label"))

	// An unassigned slot is an empty list, not a null one.
	assert.Equal(t, 0, listLen(t, objectAttr(t, devices, "sdd")))
}

func TestBuildInstanceObject_interfaces(t *testing.T) {
	object := buildTestObject(t, testInstance(), testNetwork(), testDisks(), testConfigs())

	config := listElem(t, objectAttr(t, object, "config"), 0)
	interfaces := objectAttr(t, config, "interface")

	// A VLAN interface has no VPC configuration, but none of its fields is null.
	vlan := listElem(t, interfaces, 0)
	assert.Equal(t, types.StringValue("vlan"), objectAttr(t, vlan, "purpose"))
	assert.Equal(t, types.Int64Value(0), objectAttr(t, vlan, "vpc_id"))
	assert.Equal(t, types.Int64Value(0), objectAttr(t, vlan, "subnet_id"))
	assert.Equal(t, 0, listLen(t, objectAttr(t, vlan, "ipv4")))
	assert.Equal(t, 0, listLen(t, objectAttr(t, vlan, "ipv6")))
	assert.Equal(t, 0, listLen(t, objectAttr(t, vlan, "ip_ranges")))

	vpc := listElem(t, interfaces, 1)
	assert.Equal(t, types.Int64Value(456), objectAttr(t, vpc, "vpc_id"))
	assert.Equal(t, types.Int64Value(789), objectAttr(t, vpc, "subnet_id"))
	assert.Equal(t, 1, listLen(t, objectAttr(t, vpc, "ip_ranges")))

	ipv4 := listElem(t, objectAttr(t, vpc, "ipv4"), 0)
	assert.Equal(t, types.StringValue("10.0.0.2"), objectAttr(t, ipv4, "vpc"))
	assert.Equal(t, types.StringValue("192.0.2.10"), objectAttr(t, ipv4, "nat_1_1"))

	ipv6 := listElem(t, objectAttr(t, vpc, "ipv6"), 0)
	assert.Equal(t, types.BoolValue(true), objectAttr(t, ipv6, "is_public"))

	slaac := listElem(t, objectAttr(t, ipv6, "slaac"), 0)
	assert.Equal(t, types.StringValue("2600::/64"), objectAttr(t, slaac, "range"))
	assert.Equal(t, types.StringValue("2600::/64"), objectAttr(t, slaac, "assigned_range"))
	assert.Equal(t, types.StringValue("2600::1"), objectAttr(t, slaac, "address"))

	ipv6Range := listElem(t, objectAttr(t, ipv6, "range"), 0)
	assert.Equal(t, types.StringValue("2600:1::/64"), objectAttr(t, ipv6Range, "range"))
	assert.Equal(t, types.StringValue("2600:1::/64"), objectAttr(t, ipv6Range, "assigned_range"))
}

// TestBuildInstanceObject_empty covers an Instance with every optional pointer
// unset and no disks, configs or addresses. The former SDKv2 implementation
// dereferenced these pointers unguarded; state must still come back fully
// populated with zero values.
func TestBuildInstanceObject_empty(t *testing.T) {
	object := buildTestObject(t, &linodego.Instance{ID: 1}, nil, nil, nil)

	assertNoNullValues(t, "instance", object)

	assert.Equal(t, types.StringValue(""), objectAttr(t, object, "ip_address"))
	assert.Equal(t, types.StringValue(""), objectAttr(t, object, "private_ip_address"))
	assert.Equal(t, types.StringValue(""), objectAttr(t, object, "boot_config_label"))
	assert.Equal(t, types.Int64Value(0), objectAttr(t, object, "swap_size"))

	// SDKv2 always emitted exactly one specs/alerts/backups element.
	assert.Equal(t, 1, listLen(t, objectAttr(t, object, "specs")))
	assert.Equal(t, 1, listLen(t, objectAttr(t, object, "alerts")))
	assert.Equal(t, 1, listLen(t, objectAttr(t, object, "backups")))

	// ...and an empty list when the instance has no placement group.
	assert.Equal(t, 0, listLen(t, objectAttr(t, object, "placement_group")))
	assert.Equal(t, 0, listLen(t, objectAttr(t, object, "disk")))
	assert.Equal(t, 0, listLen(t, objectAttr(t, object, "config")))
}

// TestBuildInstanceObject_nilConfigMembers covers a Config whose Devices and
// Helpers pointers are unset.
func TestBuildInstanceObject_nilConfigMembers(t *testing.T) {
	configs := []linodego.InstanceConfig{{ID: 1, Label: "bare"}}

	object := buildTestObject(t, &linodego.Instance{ID: 1}, nil, nil, configs)

	assertNoNullValues(t, "instance", object)

	config := listElem(t, objectAttr(t, object, "config"), 0)
	assert.Equal(t, 1, listLen(t, objectAttr(t, config, "helpers")))
	assert.Equal(t, 1, listLen(t, objectAttr(t, config, "devices")))
	assert.Equal(t, 0, listLen(t, objectAttr(t, config, "interface")))

	devices := listElem(t, objectAttr(t, config, "devices"), 0)
	for _, key := range helper.GetConfigDeviceKeys() {
		assert.Equalf(t, 0, listLen(t, objectAttr(t, devices, key)), "device slot %q", key)
	}
}

// TestInstanceObjectTypeMatchesSchema guards against the model and the schema
// drifting apart.
func TestInstanceObjectTypeMatchesSchema(t *testing.T) {
	object := buildTestObject(t, testInstance(), testNetwork(), testDisks(), testConfigs())

	assert.Equal(t, instanceObjectType, object.Type(t.Context()))
	assert.Len(t, object.Attributes(), len(instanceAttributes))
}
