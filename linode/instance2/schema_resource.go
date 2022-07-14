package instance2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: "The Linode’s label is for display purposes only. If no label is provided for a Linode, a default will be assigned.",
		Optional:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Description: "The region where the Linode will be located.",
		Required:    true,
		ForceNew:    true,
	},
	"type": {
		Type:        schema.TypeString,
		Description: "The Linode Type of the Linode to be created.",
		Required:    true,
		ForceNew:    true,
	},
	"tags": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "An array of tags applied to this object.",
		Optional:    true,
	},
	"private_ip": {
		Type:        schema.TypeBool,
		Description: "If true, the created Linode will have private networking enabled and assigned a private IPv4 address.",
		Optional:    true,
		Default:     false,
	},

	// IP related fields
	"ipv4_public": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "A set of public IPv4 addresses assigned to this Linode.",
		Computed:    true,
	},
	"ipv4_private": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "A set of private IPv4 addresses assigned to this Linode.",
		Computed:    true,
	},
	"ipv4_shared": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "A set of IPv4 addresses shared to this Linode.",
		Computed:    true,
	},
	"ipv4_reserved": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "A set of IPv4 addresses reserved for this Linode.",
		Computed:    true,
	},
	"ipv6_slaac": {
		Type:        schema.TypeString,
		Description: "The SLAAC IPv6 address assigned to this Linode.",
		Computed:    true,
	},
	"ipv6_link_local": {
		Type:        schema.TypeString,
		Description: "The Link Local IPv6 address assigned to this Linode.",
		Computed:    true,
	},
	"ipv6_global": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "A set of IPv6 ranges assigned to this Linode.",
		Computed:    true,
	},

	"hypervisor": {
		Type:        schema.TypeString,
		Description: "The virtualization software powering this Linode.",
		Computed:    true,
	},
	"specs": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: specSchema,
		},
		Description: "Information about the resources available to this Linode.",
		Computed:    true,
	},
	"status": {
		Type:        schema.TypeString,
		Description: "A brief description of this Linode’s current state.",
		Computed:    true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "When this Linode was created.",
		Computed:    true,
	},
	"updated": {
		Type:        schema.TypeString,
		Description: "When this Linode was last updated.",
		Computed:    true,
	},
	"watchdog_enabled": {
		Type:        schema.TypeBool,
		Description: "The watchdog, named Lassie, is a Shutdown Watchdog that monitors your Linode and will reboot it if it powers off unexpectedly.",
		Computed:    true,
	},
}

var specSchema = map[string]*schema.Schema{
	"disk": {
		Type:     schema.TypeInt,
		Computed: true,
		Description: "The amount of storage space, in GB. this Linode has access to. A typical Linode " +
			"will divide this space between a primary disk with an image deployed to it, and a swap disk, " +
			"usually 512 MB. This is the default configuration created when deploying a Linode with an image " +
			"without specifying disks.",
	},
	"memory": {
		Type:     schema.TypeInt,
		Computed: true,
		Description: "The amount of RAM, in MB, this Linode has access to. Typically a Linode will choose " +
			"to boot with all of its available RAM, but this can be configured in a Config profile.",
	},
	"vcpus": {
		Type:     schema.TypeInt,
		Computed: true,
		Description: "The number of vcpus this Linode has access to. Typically a Linode will choose to boot " +
			"with all of its available vcpus, but this can be configured in a Config Profile.",
	},
	"transfer": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The amount of network transfer this Linode is allotted each month.",
	},
}
