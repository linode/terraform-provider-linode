package instance2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
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

	// Optional fields
	"alerts": {
		Description: "Configuration options for alert triggers on this Linode.",
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: alertsSchema,
		},
	},
	"backups_enabled": {
		Type: schema.TypeBool,
		Description: "If this field is set to true, the created Linode will automatically be enrolled in the " +
			"Linode Backup service.",
		Optional: true,
		Default:  false,
	},
	"backup_window": {
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: backupScheduleSchema,
		},
	},
	"label": {
		Type:        schema.TypeString,
		Description: "The Linode’s label is for display purposes only. If no label is provided for a Linode, a default will be assigned.",
		Optional:    true,
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

	// Computed fields
	"created": {
		Type:        schema.TypeString,
		Description: "When this Linode was created.",
		Computed:    true,
	},
	"hypervisor": {
		Type:        schema.TypeString,
		Description: "The virtualization software powering this Linode.",
		Computed:    true,
	},
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

var alertsSchema = map[string]*schema.Schema{
	"cpu": {
		Type:     schema.TypeInt,
		Computed: true,
		Optional: true,
		Description: "The percentage of CPU usage required to trigger an alert. If the average CPU usage " +
			"over two hours exceeds this value, we'll send you an alert. If this is set to 0, the alert is disabled.",
	},
	"network_in": {
		Type:     schema.TypeInt,
		Computed: true,
		Optional: true,
		Description: "The amount of incoming traffic, in Mbit/s, required to trigger an alert. If the average " +
			"incoming traffic over two hours exceeds this value, we'll send you an alert. " +
			"If this is set to 0 (zero), the alert is disabled.",
	},
	"network_out": {
		Type:     schema.TypeInt,
		Computed: true,
		Optional: true,
		Description: "The amount of outbound traffic, in Mbit/s, required to trigger an alert. " +
			"If the average outbound traffic over two hours exceeds this value, we'll send you an alert. " +
			"If this is set to 0 (zero), the alert is disabled.",
	},
	"transfer_quota": {
		Type:     schema.TypeInt,
		Computed: true,
		Optional: true,
		Description: "The percentage of network transfer that may be used before an alert is triggered. " +
			"When this value is exceeded, we'll alert you. " +
			"If this is set to 0 (zero), the alert is disabled.",
	},
	"io": {
		Type:     schema.TypeInt,
		Computed: true,
		Optional: true,
		Description: "The amount of disk IO operation per second required to trigger an alert. " +
			"If the average disk IO over two hours exceeds this value, we'll send you an alert. " +
			"If set to 0, this alert is disabled.",
	},
}

var backupScheduleSchema = map[string]*schema.Schema{
	"day": {
		Type: schema.TypeString,
		Description: "The day ('Sunday'-'Saturday') of the week that your Linode's weekly Backup is " +
			"taken. If not set manually, a day will be chosen for you. Backups are taken every day, " +
			"but backups taken on this day are preferred when selecting backups to retain for a " +
			"longer period.  If not set manually, then when backups are initially enabled, this " +
			"may come back as 'Scheduling' until the day is automatically selected.",
		Required: true,
	},
	"window": {
		Type: schema.TypeString,
		Description: "The window ('W0'-'W22') in which your backups will be taken, in UTC. A " +
			"backups window is a two-hour span of time in which the backup may occur. For example, " +
			"'W10' indicates that your backups should be taken between 10:00 and 12:00. If you do " +
			"not choose a backup window, one will be selected for you automatically.  If not set " +
			"manually, when backups are initially enabled this may come back as Scheduling until " +
			"the window is automatically selected.",
		Required: true,
	},
}
