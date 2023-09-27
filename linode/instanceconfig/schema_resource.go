package instanceconfig

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const deviceDescription = "Device can be either a Disk or Volume identified by disk_id or " +
	"volume_id. Only one type per slot allowed."

var resourceSchema = map[string]*schema.Schema{
	"linode_id": {
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "The ID of the Linode to create this configuration profile under.",
	},
	"label": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Config's label for display purposes only.",
	},

	"device": {
		Type:          schema.TypeSet,
		Elem:          &schema.Resource{Schema: deviceV2Schema},
		Optional:      true,
		Computed:      true,
		ConflictsWith: []string{"devices"},
		Description:   "Blocks for device disks in a Linode's configuration profile.",
	},

	"devices": {
		Type:          schema.TypeList,
		Elem:          &schema.Resource{Schema: devicesSchema},
		Optional:      true,
		Computed:      true,
		MaxItems:      1,
		ConflictsWith: []string{"device"},
		Deprecated:    "Devices attribute is deprecated in favor of `device`.",
		Description:   "A dictionary of device disks to use as a device map in a Linode's configuration profile.",
	},

	"booted": {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
		Description: "If true, the Linode will be booted to running state. " +
			"If false, the Linode will be shutdown. If undefined, no action will be taken.",
	},
	"comments": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Optional field for arbitrary User comments on this Config.",
	},
	"helpers": {
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: helpersSchema},
		Optional:    true,
		Computed:    true,
		Description: "Helpers enabled when booting to this Linode Config.",
	},
	"interface": {
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: interfaceSchema},
		Optional:    true,
		Description: "An array of Network Interfaces to add to this Linode's Configuration Profile.",
	},
	"kernel": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "linode/latest-64bit",
		Description: "A Kernel ID to boot a Linode with. Defaults to “linode/latest-64bit”.",
	},
	"memory_limit": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "The memory limit of the Linode.",
	},
	"root_device": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "/dev/sda",
		Description: "The root device to boot. " +
			"If no value or an invalid value is provided, root device will default to /dev/sda. " +
			"If the device specified at the root device location is not mounted, " +
			"the Linode will not boot until a device is mounted.",
	},
	"run_level": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "default",
		Description: "Defines the state of your Linode after booting.",
		ValidateDiagFunc: validation.ToDiagFunc(
			validation.StringInSlice([]string{"default", "single", "binbash"}, true),
		),
	},
	"virt_mode": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "paravirt",
		Description: "Controls the virtualization mode.",
		ValidateDiagFunc: validation.ToDiagFunc(
			validation.StringInSlice([]string{"paravirt", "fullvirt"}, true),
		),
	},
}

var devicesSchema = map[string]*schema.Schema{
	"sda": {
		Type:        schema.TypeList,
		Description: deviceDescription,
		MaxItems:    1,
		Optional:    true,
		Elem:        &schema.Resource{Schema: deviceSchema},
	},
	"sdb": {
		Type:        schema.TypeList,
		Description: deviceDescription,
		MaxItems:    1,
		Optional:    true,
		Elem:        &schema.Resource{Schema: deviceSchema},
	},
	"sdc": {
		Type:        schema.TypeList,
		Description: deviceDescription,
		MaxItems:    1,
		Optional:    true,
		Elem:        &schema.Resource{Schema: deviceSchema},
	},
	"sdd": {
		Type:        schema.TypeList,
		Description: deviceDescription,
		MaxItems:    1,
		Optional:    true,
		Elem:        &schema.Resource{Schema: deviceSchema},
	},
	"sde": {
		Type:        schema.TypeList,
		Description: deviceDescription,
		MaxItems:    1,
		Optional:    true,
		Elem:        &schema.Resource{Schema: deviceSchema},
	},
	"sdf": {
		Type:        schema.TypeList,
		Description: deviceDescription,
		MaxItems:    1,
		Optional:    true,
		Elem:        &schema.Resource{Schema: deviceSchema},
	},
	"sdg": {
		Type:        schema.TypeList,
		Description: deviceDescription,
		MaxItems:    1,
		Optional:    true,
		Elem:        &schema.Resource{Schema: deviceSchema},
	},
	"sdh": {
		Type:        schema.TypeList,
		Description: deviceDescription,
		MaxItems:    1,
		Optional:    true,
		Elem:        &schema.Resource{Schema: deviceSchema},
	},
}

var deviceV2Schema = map[string]*schema.Schema{
	"device_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Disk ID to map to this disk slot",
	},
	"disk_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The Disk ID to map to this disk slot",
	},
	"volume_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The Block Storage volume ID to map to this disk slot",
	},
}

var deviceSchema = map[string]*schema.Schema{
	"disk_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The Disk ID to map to this disk slot",
	},
	"volume_id": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The Block Storage volume ID to map to this disk slot",
	},
}

var helpersSchema = map[string]*schema.Schema{
	"devtmpfs_automount": {
		Type:        schema.TypeBool,
		Description: "Populates the /dev directory early during boot without udev.",
		Optional:    true,
		Default:     true,
	},
	"distro": {
		Type:        schema.TypeBool,
		Description: "Helps maintain correct inittab/upstart console device.",
		Optional:    true,
		Default:     true,
	},
	"modules_dep": {
		Type:        schema.TypeBool,
		Description: "Creates a modules dependency file for the Kernel you run.",
		Optional:    true,
		Default:     true,
	},
	"network": {
		Type:        schema.TypeBool,
		Description: "Automatically configures static networking.",
		Optional:    true,
		Default:     true,
	},
	"updatedb_disabled": {
		Type:        schema.TypeBool,
		Description: "Disables updatedb cron job to avoid disk thrashing.",
		Optional:    true,
		Default:     true,
	},
}

const (
	onlyAllowedForVPCMsg  = "This attribute is only allowed for VPC interfaces."
	onlyAllowedForVLANMsg = "This attribute is only allowed for VLAN interfaces."
	requiredForVPCMsg     = "This attribute is required for VPC interfaces."
	requiredForVLANMsg    = "This attribute is required for VPC interfaces."
)

var interfaceSchema = map[string]*schema.Schema{
	"purpose": {
		Type:        schema.TypeString,
		Description: "The type of interface.",
		Required:    true,
		ValidateDiagFunc: validation.ToDiagFunc(
			validation.StringInSlice([]string{"public", "vlan", "vpc"}, true),
		),
	},
	"ipam_address": {
		Type: schema.TypeString,
		Description: "This Network Interface's private IP address in " +
			"Classless Inter-Domain Routing (CIDR) notation." +
			onlyAllowedForVLANMsg,
		Optional: true,
	},
	"label": {
		Type: schema.TypeString,
		Description: "The name of the VALN. " + requiredForVLANMsg +
			" " + onlyAllowedForVLANMsg,
		Optional: true,
	},
	"id": {
		Type:        schema.TypeInt,
		Description: "The ID of the interface.",
		Computed:    true,
	},
	"subnet_id": {
		Type: schema.TypeInt,
		Description: "The ID of the subnet which the VPC interface is connected to." +
			requiredForVPCMsg + onlyAllowedForVPCMsg,
		Optional: true,
	},
	"vpc_id": {
		Type: schema.TypeInt,
		Description: "The ID of VPC of the subnet which the VPC " +
			"interface is connected to.",
		Computed: true,
	},
	"primary": {
		Type: schema.TypeBool,
		Description: "Whether the interface is the primary interface that should " +
			"have the default route for this Linode.",
		Optional: true,
		Default:  false,
	},
	"ipv4": {
		Type: schema.TypeList,
		Description: "The IPv4 configuration of the VPC interface." +
			onlyAllowedForVPCMsg,
		Computed: true,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"vpc": {
					Type:        schema.TypeString,
					Description: "The IP from the VPC subnet to use for this interface.",
					Computed:    true,
					Optional:    true,
				},
				"nat_1_1": {
					Type: schema.TypeString,
					Description: "The public IP that will be used for the " +
						"one-to-one NAT purpose.",
					Computed: true,
					Optional: true,
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						if new == "any" && old != "" {
							return true
						}
						return old == new
					},
				},
			},
		},
	},
}
