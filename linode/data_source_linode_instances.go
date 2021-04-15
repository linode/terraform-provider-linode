package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func dataSourceLinodeInstancesInstances() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image": {
				Type: schema.TypeString,
				Description: "An Image ID to deploy the Disk from. Official Linode Images start with linode/, while " +
					"your Images start with private/. See /images for more information on the Images available for you to use.",
				Computed: true,
			},
			"label": {
				Type: schema.TypeString,
				Description: "The Linode's label is for display purposes only. If no label is provided for a Linode, " +
					"a default will be assigned",
				Computed: true,
			},
			"group": {
				Type:        schema.TypeString,
				Description: "The display group of the Linode instance.",
				Computed:    true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"boot_config_label": {
				Type:        schema.TypeString,
				Description: "The Label of the Instance Config that should be used to boot the Linode instance.",
				Computed:    true,
			},
			"region": {
				Type: schema.TypeString,
				Description: "This is the location where the Linode was deployed. This cannot be changed without " +
					"opening a support ticket.",
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of instance to be deployed, determining the price and size.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The status of the instance, indicating the current readiness state.",
				Computed:    true,
			},
			"ip_address": {
				Type: schema.TypeString,
				Description: "This Linode's Public IPv4 Address. If there are multiple public IPv4 addresses on this " +
					"Instance, an arbitrary address will be used for this field.",
				Computed: true,
			},
			"ipv6": {
				Type:        schema.TypeString,
				Description: "This Linode's IPv6 SLAAC addresses. This address is specific to a Linode, and may not be shared.",
				Computed:    true,
			},

			"ipv4": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				Description: "This Linode's IPv4 Addresses. Each Linode is assigned a single public IPv4 address upon " +
					"creation, and may get a single private IPv4 address if needed. You may need to open a support " +
					"ticket to get additional IPv4 addresses.",
				Computed: true,
			},
			"private_ip_address": {
				Type: schema.TypeString,
				Description: "This Linode's Private IPv4 Address.  The regional private IP address range is " +
					"192.168.128/17 address shared by all Linode Instances in a region.",
				Computed: true,
			},
			"swap_size": {
				Type: schema.TypeInt,
				Description: "When deploying from an Image, this field is optional with a Linode API default of " +
					"512mb, otherwise it is ignored. This is used to set the swap disk size for the newly-created Linode.",
				Computed: true,
			},
			"watchdog_enabled": {
				Type: schema.TypeBool,
				Description: "The watchdog, named Lassie, is a Shutdown Watchdog that monitors your Linode and will " +
					"reboot it if it powers off unexpectedly. It works by issuing a boot job when your Linode powers " +
					"off without a shutdown job being responsible. To prevent a loop, Lassie will give up if there have " +
					"been more than 5 boot jobs issued within 15 minutes.",
				Computed: true,
			},
			"specs": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The amount of storage space, in GB. this Linode has access to. A typical " +
								"Linode will divide this space between a primary disk with an image deployed to it, " +
								"and a swap disk, usually 512 MB. This is the default configuration created when " +
								"deploying a Linode with an image without specifying disks.",
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The amount of RAM, in MB, this Linode has access to. Typically a Linode will " +
								"choose to boot with all of its available RAM, but this can be configured in a Config profile.",
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The number of vcpus this Linode has access to. Typically a Linode will " +
								"choose to boot with all of its available vcpus, but this can be configured in a Config Profile.",
						},
						"transfer": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of network transfer this Linode is allotted each month.",
						},
					},
				},
			},

			"alerts": {
				Computed: true,
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The percentage of CPU usage required to trigger an alert. If the average " +
								"CPU usage over two hours exceeds this value, we'll send you an alert. If this is set " +
								"to 0, the alert is disabled.",
						},
						"network_in": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The amount of incoming traffic, in Mbit/s, required to trigger an alert. " +
								"If the average incoming traffic over two hours exceeds this value, we'll send you an " +
								"alert. If this is set to 0 (zero), the alert is disabled.",
						},
						"network_out": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The amount of outbound traffic, in Mbit/s, required to trigger an alert. " +
								"If the average outbound traffic over two hours exceeds this value, we'll send you an alert. " +
								"If this is set to 0 (zero), the alert is disabled.",
						},
						"transfer_quota": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The percentage of network transfer that may be used before an alert is triggered. " +
								"When this value is exceeded, we'll alert you. If this is set to 0 (zero), the alert is disabled.",
						},
						"io": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "The amount of disk IO operation per second required to trigger an alert. " +
								"If the average disk IO over two hours exceeds this value, we'll send you an alert. " +
								"If set to 0, this alert is disabled.",
						},
					},
				},
			},
			"backups": {
				Type:        schema.TypeList,
				Description: "Information about this Linode's backups status.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If this Linode has the Backup service enabled.",
						},
						"schedule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type: schema.TypeString,
										Description: "The day ('Sunday'-'Saturday') of the week that your Linode's weekly " +
											"Backup is taken. If not set manually, a day will be chosen for you. Backups are " +
											"taken every day, but backups taken on this day are preferred when selecting backups " +
											"to retain for a longer period.  If not set manually, then when backups are initially " +
											"enabled, this may come back as 'Scheduling' until the day is automatically selected.",
										Computed: true,
									},
									"window": {
										Type: schema.TypeString,
										Description: "The window ('W0'-'W22') in which your backups will be taken, in UTC. " +
											"A backups window is a two-hour span of time in which the backup may occur. " +
											"For example, 'W10' indicates that your backups should be taken between 10:00 " +
											"and 12:00. If you do not choose a backup window, one will be selected for you " +
											"automatically.  If not set manually, when backups are initially enabled this " +
											"may come back as Scheduling until the window is automatically selected.",
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"config": {
				Description: "Configuration profiles define the VM settings and boot behavior of the Linode Instance.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:        schema.TypeString,
							Description: "The Config's label for display purposes.  Also used by `boot_config_label`.",
							Computed:    true,
						},
						"helpers": {
							Type:        schema.TypeList,
							Description: "Helpers enabled when booting to this Linode Config.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"updatedb_disabled": {
										Type:        schema.TypeBool,
										Description: "Disables updatedb cron job to avoid disk thrashing.",
										Computed:    true,
									},
									"distro": {
										Type:        schema.TypeBool,
										Description: "Controls the behavior of the Linode Config's Distribution Helper setting.",
										Computed:    true,
									},
									"modules_dep": {
										Type:        schema.TypeBool,
										Description: "Creates a modules dependency file for the Kernel you run.",
										Computed:    true,
									},
									"network": {
										Type: schema.TypeBool,
										Description: "Controls the behavior of the Linode Config's Network Helper setting, used to " +
											"automatically configure additional IP addresses assigned to this instance.",
										Computed: true,
									},
									"devtmpfs_automount": {
										Type:        schema.TypeBool,
										Description: "Populates the /dev directory early during boot without udev. Defaults to false.",
										Computed:    true,
									},
								},
							},
						},
						"devices": {
							Type: schema.TypeList,
							Description: "Device sda-sdh can be either a Disk or Volume identified by disk_label or " +
								"volume_id. Only one type per slot allowed.",
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sda": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceLinodeInstanceDeviceDisk(),
									},
									"sdb": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceLinodeInstanceDeviceDisk(),
									},
									"sdc": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceLinodeInstanceDeviceDisk(),
									},
									"sdd": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceLinodeInstanceDeviceDisk(),
									},
									"sde": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceLinodeInstanceDeviceDisk(),
									},
									"sdf": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceLinodeInstanceDeviceDisk(),
									},
									"sdg": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceLinodeInstanceDeviceDisk(),
									},
									"sdh": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     resourceLinodeInstanceDeviceDisk(),
									},
								},
							},
						},
						"kernel": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "A Kernel ID to boot a Linode with. Default is based on image choice. " +
								"(examples: linode/latest-64bit, linode/grub2, linode/direct-disk)",
						},
						"run_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the state of your Linode after booting. Defaults to default.",
						},
						"virt_mode": {
							Type:        schema.TypeString,
							Description: "Controls the virtualization mode. Defaults to paravirt.",
							Computed:    true,
						},
						"root_device": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The root device to boot. The corresponding disk must be attached.",
						},
						"comments": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional field for arbitrary User comments on this Config.",
						},

						"memory_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Defaults to the total RAM of the Linode",
						},
					},
				},
			},
			"disk": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:        schema.TypeString,
							Description: "The disks label, which acts as an identifier in Terraform.",
							Computed:    true,
						},
						"size": {
							Type:        schema.TypeInt,
							Description: "The size of the Disk in MB.",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeInt,
							Description: "The ID of the Disk (for use in Linode Image resources and Linode Instance Config Devices)",
							Computed:    true,
						},
						"filesystem": {
							Type:        schema.TypeString,
							Description: "The Disk filesystem can be one of: raw, swap, ext3, ext4, initrd (max 32mb)",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLinodeInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeInstancesRead,
		Schema: map[string]*schema.Schema{
			"filter": filterSchema([]string{"group", "id", "image", "label", "region", "tags"}),
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceLinodeInstancesInstances(),
			},
		},
	}
}

func dataSourceLinodeInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client

	filter, err := constructFilterString(d, instanceValueToFilterType)
	if err != nil {
		return fmt.Errorf("failed to construct filter: %s", err)
	}

	instances, err := client.ListInstances(context.Background(), &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return fmt.Errorf("failed to get instances: %s", err)
	}

	flattenedInstances := make([]map[string]interface{}, len(instances))
	for i, instance := range instances {
		instanceMap, err := flattenLinodeInstance(&client, &instance)
		if err != nil {
			return fmt.Errorf("failed to translate instance to map: %s", err)
		}

		flattenedInstances[i] = instanceMap
	}

	d.SetId(fmt.Sprintf(filter))
	d.Set("instances", flattenedInstances)

	return nil
}

func flattenLinodeInstance(client *linodego.Client, instance *linodego.Instance) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	id := instance.ID

	instanceNetwork, err := client.GetInstanceIPAddresses(context.Background(), int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get ips for linode instance %d: %s", id, err)
	}

	var ips []string
	for _, ip := range instance.IPv4 {
		ips = append(ips, ip.String())
	}

	result["ipv4"] = ips
	result["ipv6"] = instance.IPv6

	public, private := instanceNetwork.IPv4.Public, instanceNetwork.IPv4.Private

	if len(public) > 0 {
		result["ip_address"] = public[0].Address
	}

	if len(private) > 0 {
		result["private_ip_address"] = private[0].Address
	}

	result["label"] = instance.Label
	result["status"] = instance.Status
	result["type"] = instance.Type
	result["region"] = instance.Region
	result["watchdog_enabled"] = instance.WatchdogEnabled
	result["group"] = instance.Group
	result["tags"] = instance.Tags
	result["image"] = instance.Image

	result["backups"] = flattenInstanceBackups(*instance)
	result["specs"] = flattenInstanceSpecs(*instance)
	result["alerts"] = flattenInstanceAlerts(*instance)

	instanceDisks, err := client.ListInstanceDisks(context.Background(), int(id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the disks for the Linode instance %d: %s", id, err)
	}

	disks, swapSize := flattenInstanceDisks(instanceDisks)
	result["disk"] = disks
	result["swap_size"] = swapSize

	instanceConfigs, err := client.ListInstanceConfigs(context.Background(), int(id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the config for Linode instance %d (%s): %s", id, instance.Label, err)
	}

	diskLabelIDMap := make(map[int]string, len(instanceDisks))
	for _, disk := range instanceDisks {
		diskLabelIDMap[disk.ID] = disk.Label
	}

	configs := flattenInstanceConfigs(instanceConfigs, diskLabelIDMap)

	result["config"] = configs
	if len(instanceConfigs) == 1 {
		result["boot_config_label"] = instanceConfigs[0].Label
	}

	return result, nil
}

// instanceValueToFilterType converts the given value to the correct type depending on the filter name.
func instanceValueToFilterType(filterName, value string) (interface{}, error) {
	switch filterName {
	case "id":
		return strconv.Atoi(value)
	}

	return value, nil
}
