package instance

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":             {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
	"image":          {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"label":          {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"region":         {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"lke_cluster_id": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},

	// Tags must be filtered on the client
	"tags":                 {TypeFunc: frameworkfilter.FilterTypeString},
	"status":               {TypeFunc: frameworkfilter.FilterTypeString},
	"type":                 {TypeFunc: frameworkfilter.FilterTypeString},
	"watchdog_enabled":     {TypeFunc: frameworkfilter.FilterTypeBool},
	"disk_encryption":      {TypeFunc: frameworkfilter.FilterTypeString},
	"interface_generation": {TypeFunc: frameworkfilter.FilterTypeString},
}

var instanceSpecsAttribute = schema.ListNestedAttribute{
	Computed: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"disk": schema.Int64Attribute{
				Computed: true,
				Description: "The amount of storage space, in GB. this Linode has access to. A typical " +
					"Linode will divide this space between a primary disk with an image deployed to it, " +
					"and a swap disk, usually 512 MB. This is the default configuration created when " +
					"deploying a Linode with an image without specifying disks.",
			},
			"memory": schema.Int64Attribute{
				Computed: true,
				Description: "The amount of RAM, in MB, this Linode has access to. Typically a Linode will " +
					"choose to boot with all of its available RAM, but this can be configured in a Config profile.",
			},
			"vcpus": schema.Int64Attribute{
				Computed: true,
				Description: "The number of vcpus this Linode has access to. Typically a Linode will " +
					"choose to boot with all of its available vcpus, but this can be configured in a Config Profile.",
			},
			"transfer": schema.Int64Attribute{
				Computed:    true,
				Description: "The amount of network transfer this Linode is allotted each month.",
			},
			"accelerated_devices": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of VPUs this Linode has access to.",
			},
			"gpus": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of GPUs this Linode has access to.",
			},
		},
	},
}

var instanceAlertsAttribute = schema.ListNestedAttribute{
	Computed: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"cpu": schema.Int64Attribute{
				Computed: true,
				Description: "The percentage of CPU usage required to trigger an alert. If the average " +
					"CPU usage over two hours exceeds this value, we'll Device can be either a Disk or Volume " +
					"identified by disk_id or volume_id. Only one type per slot allowed.send you an alert. If " +
					"this is set to 0, the alert is disabled.",
			},
			"network_in": schema.Int64Attribute{
				Computed: true,
				Description: "The amount of incoming traffic, in Mbit/s, required to trigger an alert. " +
					"If the average incoming traffic over two hours exceeds this value, we'll send you an " +
					"alert. If this is set to 0 (zero), the alert is disabled.",
			},
			"network_out": schema.Int64Attribute{
				Computed: true,
				Description: "The amount of outbound traffic, in Mbit/s, required to trigger an alert. " +
					"If the average outbound traffic over two hours exceeds this value, we'll send you an alert. " +
					"If this is set to 0 (zero), the alert is disabled.",
			},
			"transfer_quota": schema.Int64Attribute{
				Computed: true,
				Description: "The percentage of network transfer that may be used before an alert is triggered. " +
					"When this value is exceeded, we'll alert you. If this is set to 0 (zero), the alert is disabled.",
			},
			"io": schema.Int64Attribute{
				Computed: true,
				Description: "The amount of disk IO operation per second required to trigger an alert. " +
					"If the average disk IO over two hours exceeds this value, we'll send you an alert. " +
					"If set to 0, this alert is disabled.",
			},
		},
	},
}

var instanceBackupsAttribute = schema.ListNestedAttribute{
	Description: "Information about this Linode's backups status.",
	Computed:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"available": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this Backup is available for restoration.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "If this Linode has the Backup service enabled.",
			},
			"schedule": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"day": schema.StringAttribute{
							Description: "The day ('Sunday'-'Saturday') of the week that your Linode's weekly " +
								"Backup is taken. If not set manually, a day will be chosen for you. Backups are " +
								"taken every day, but backups taken on this day are preferred when selecting backups " +
								"to retain for a longer period.  If not set manually, then when backups are initially " +
								"enabled, this may come back as 'Scheduling' until the day is automatically selected.",
							Computed: true,
						},
						"window": schema.StringAttribute{
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
}

var instanceConfigInterfaceAttribute = schema.ListNestedAttribute{
	Description: "An array of Network Interfaces for this Linode’s Configuration Profile.",
	Computed:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"purpose": schema.StringAttribute{
				Description: "The type of interface.",
				Computed:    true,
			},
			"ipam_address": schema.StringAttribute{
				Description: "This Network Interface's private IP address in " +
					"Classless Inter-Domain Routing (CIDR) notation." +
					onlyAllowedForVLANMsg,
				Computed: true,
			},
			"label": schema.StringAttribute{
				Description: "The name of the VALN. " + requiredForVLANMsg +
					" " + onlyAllowedForVLANMsg,
				Computed: true,
			},
			"id": schema.Int64Attribute{
				Description: "The ID of the interface.",
				Computed:    true,
			},
			"subnet_id": schema.Int64Attribute{
				Description: "The ID of the subnet which the VPC interface is connected to." +
					requiredForVPCMsg + onlyAllowedForVPCMsg,
				Computed: true,
			},
			"vpc_id": schema.Int64Attribute{
				Description: "The ID of VPC of the subnet which the VPC " +
					"interface is connected to.",
				Computed: true,
			},
			"primary": schema.BoolAttribute{
				Description: "Whether the interface is the primary interface that should " +
					"have the default route for this Linode.",
				Computed: true,
			},
			"active": schema.BoolAttribute{
				Description: "Whether this interface is currently booted and active.",
				Computed:    true,
			},
			"ip_ranges": schema.ListAttribute{
				Description: "List of VPC IPs or IP ranges inside the VPC subnet.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"ipv4": schema.ListNestedAttribute{
				Description: "The IPv4 configuration of the VPC interface." +
					onlyAllowedForVPCMsg,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc": schema.StringAttribute{
							Description: "The IP from the VPC subnet to use for this interface.",
							Computed:    true,
						},
						"nat_1_1": schema.StringAttribute{
							Description: "The public IP that will be used for the " +
								"one-to-one NAT purpose.",
							Computed: true,
						},
					},
				},
			},
			"ipv6": schema.ListNestedAttribute{
				Description: "The IPv6 configuration of the VPC interface. " +
					onlyAllowedForVPCMsg,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"slaac": schema.ListNestedAttribute{
							Description: "An array of SLAAC prefixes to use for this interface.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"range": schema.StringAttribute{
										Description: "A SLAAC prefix to add to this interface, " +
											"or `auto` for a new IPv6 prefix to be automatically allocated.",
										Computed: true,
									},
									"assigned_range": schema.StringAttribute{
										Description: "The value of `range` computed by the API. " +
											"This is necessary when needing to access the range " +
											"implicitly allocated using `auto`.",
										Computed: true,
									},
									"address": schema.StringAttribute{
										Description: "The SLAAC address chosen for this interface.",
										Computed:    true,
									},
								},
							},
						},
						"range": schema.ListNestedAttribute{
							Description: "An array of SLAAC prefixes to use for this interface.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"range": schema.StringAttribute{
										Description: "A prefix to add to this interface, " +
											"or `auto` for a new IPv6 prefix to be automatically allocated.",
										Computed: true,
									},
									"assigned_range": schema.StringAttribute{
										Description: "The value of `range` computed by the API. " +
											"This is necessary when needing to access the range " +
											"implicitly allocated using `auto`.",
										Computed: true,
									},
								},
							},
						},
						"is_public": schema.BoolAttribute{
							Description: "If true, connections from the interface to IPv6 addresses outside the VPC, " +
								"and connections from IPv6 addresses outside the VPC to the interface will be permitted.",
							Computed: true,
						},
					},
				},
			},
		},
	},
}

// instanceConfigDeviceAttributes mirrors the SDKv2 `resourceDeviceDisk` schema
// used for each device slot of a Config.
var instanceConfigDeviceAttributes = map[string]schema.Attribute{
	"disk_label": schema.StringAttribute{
		Computed:    true,
		Description: "The `label` of the `disk` to map to this `device` slot.",
	},
	"disk_id": schema.Int64Attribute{
		Computed:    true,
		Description: "The Disk ID to map to this disk slot",
	},
	"volume_id": schema.Int64Attribute{
		Computed:    true,
		Description: "The Block Storage volume ID to map to this disk slot",
	},
}

// dataSourceFrameworkDevicesAttributes builds an attribute for every valid
// config device slot (sda-sdbl).
func dataSourceFrameworkDevicesAttributes() map[string]schema.Attribute {
	result := make(map[string]schema.Attribute, 64)

	for _, key := range helper.GetConfigDeviceKeys() {
		result[key] = schema.ListNestedAttribute{
			Description: deviceDescription,
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: instanceConfigDeviceAttributes,
			},
		}
	}

	return result
}

var instanceConfigAttribute = schema.ListNestedAttribute{
	Description: "Configuration profiles define the VM settings and boot behavior of the Linode Instance.",
	Computed:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique ID of this Config.",
				Computed:    true,
			},
			"label": schema.StringAttribute{
				Description: "The Config's label for display purposes.  Also used by `boot_config_label`.",
				Computed:    true,
			},
			"helpers": schema.ListNestedAttribute{
				Description: "Helpers enabled when booting to this Linode Config.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"updatedb_disabled": schema.BoolAttribute{
							Description: "Disables updatedb cron job to avoid disk thrashing.",
							Computed:    true,
						},
						"distro": schema.BoolAttribute{
							Description: "Controls the behavior of the Linode Config's Distribution Helper setting.",
							Computed:    true,
						},
						"modules_dep": schema.BoolAttribute{
							Description: "Creates a modules dependency file for the Kernel you run.",
							Computed:    true,
						},
						"network": schema.BoolAttribute{
							Description: "Controls the behavior of the Linode Config's Network Helper setting, used to " +
								"automatically configure additional IP addresses assigned to this instance.",
							Computed: true,
						},
						"devtmpfs_automount": schema.BoolAttribute{
							Description: "Populates the /dev directory early during boot without udev. Defaults to false.",
							Computed:    true,
						},
					},
				},
			},
			"devices": schema.ListNestedAttribute{
				Description: "Device sda-sdbl can be either a Disk or Volume identified by disk_label or " +
					"volume_id. Only one type per slot allowed.",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: dataSourceFrameworkDevicesAttributes(),
				},
			},
			"interface": instanceConfigInterfaceAttribute,
			"kernel": schema.StringAttribute{
				Computed: true,
				Description: "A Kernel ID to boot a Linode with. Default is based on image choice. " +
					"(examples: linode/latest-64bit, linode/grub2, linode/direct-disk)",
			},
			"run_level": schema.StringAttribute{
				Computed:    true,
				Description: "Defines the state of your Linode after booting. Defaults to default.",
			},
			"virt_mode": schema.StringAttribute{
				Description: "Controls the virtualization mode. Defaults to paravirt.",
				Computed:    true,
			},
			"root_device": schema.StringAttribute{
				Computed:    true,
				Description: "The root device to boot. The corresponding disk must be attached.",
			},
			"comments": schema.StringAttribute{
				Computed:    true,
				Description: "Optional field for arbitrary User comments on this Config.",
			},

			"memory_limit": schema.Int64Attribute{
				Computed:    true,
				Description: "Defaults to the total RAM of the Linode",
			},
		},
	},
}

var instanceDiskAttribute = schema.ListNestedAttribute{
	Computed:    true,
	Description: "Disks associated with this Linode.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"label": schema.StringAttribute{
				Description: "The disks label, which acts as an identifier in Terraform.",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "The size of the Disk in MB.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "The ID of the Disk (for use in Linode Image resources and Linode Instance Config Devices)",
				Computed:    true,
			},
			"filesystem": schema.StringAttribute{
				Description: "The Disk filesystem can be one of: raw, swap, ext3, ext4, initrd (max 32mb)",
				Computed:    true,
			},
		},
	},
}

var instancePlacementGroupAttribute = schema.ListNestedAttribute{
	Computed: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The placement group's ID. You need to provide it for all operations impacting it.",
				Computed:    true,
			},
			"label": schema.StringAttribute{
				Description: "The unique name set for the placement group.",
				Computed:    true,
			},
			"placement_group_type": schema.StringAttribute{
				Description: "How compute instances are distributed in your placement group. " +
					"anti-affinity:local places compute instances in separate fault domains, but still in the same region.",
				Computed: true,
			},
			"placement_group_policy": schema.StringAttribute{
				Description: "How the API enforces your placement_group_type. Set to strict, your group is strict. You can't " +
					"add more compute instances to your placement group if your preferred container lacks capacity or is" +
					" unavailable. Set to flexible, your group is flexible. You can add more compute instances to it even if " +
					"they violate the placement_group_type. If you violate the placement_group_type your placement group becomes " +
					"non-compliant and you need to wait for our assistance.",
				Computed: true,
			},
		},
	},
}

// instanceAttributes describes a single Instance returned by this data source.
var instanceAttributes = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the Linode instance.",
		Computed:    true,
	},
	"image": schema.StringAttribute{
		Description: "An Image ID to deploy the Disk from. Official Linode Images start with linode/, while " +
			"your Images start with private/. See /images for more information on the Images available for you to use.",
		Computed: true,
	},
	"label": schema.StringAttribute{
		Description: "The Linode's label is for display purposes only. If no label is provided for a Linode, " +
			"a default will be assigned",
		Computed: true,
	},
	"tags": schema.SetAttribute{
		Description: "The tags assigned to this Instance.",
		ElementType: types.StringType,
		Computed:    true,
	},
	"capabilities": schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "A list of capabilities of this Linode instance.",
	},
	"locks": schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "A list of locks applied to this Linode.",
	},
	"boot_config_label": schema.StringAttribute{
		Description: "The Label of the Instance Config that should be used to boot the Linode instance.",
		Computed:    true,
	},
	"region": schema.StringAttribute{
		Description: "This is the location where the Linode was deployed. This cannot be changed without " +
			"opening a support ticket.",
		Computed: true,
	},
	"maintenance_policy": schema.StringAttribute{
		Description: "This is the maintenance type of the Linode instance.",
		Computed:    true,
	},
	"type": schema.StringAttribute{
		Description: "The type of instance to be deployed, determining the price and size.",
		Computed:    true,
	},
	"status": schema.StringAttribute{
		Description: "The status of the instance, indicating the current readiness state.",
		Computed:    true,
	},
	"interface_generation": schema.StringAttribute{
		Description: "The interface type for the Linode. ",
		Computed:    true,
	},
	"ip_address": schema.StringAttribute{
		Description: "This Linode's Public IPv4 Address. If there are multiple public IPv4 addresses on this " +
			"Instance, an arbitrary address will be used for this field.",
		Computed: true,
	},
	"ipv6": schema.StringAttribute{
		Description: "This Linode's IPv6 SLAAC addresses. This address is specific to a Linode, and may not be shared.",
		Computed:    true,
	},

	"ipv4": schema.SetAttribute{
		ElementType: types.StringType,
		Description: "This Linode's IPv4 Addresses. Each Linode is assigned a single public IPv4 address upon " +
			"creation, and may get a single private IPv4 address if needed. You may need to open a support " +
			"ticket to get additional IPv4 addresses.",
		Computed: true,
	},
	"private_ip_address": schema.StringAttribute{
		Description: "This Linode's Private IPv4 Address.  The regional private IP address range is " +
			"192.168.128/17 address shared by all Linode Instances in a region.",
		Computed: true,
	},
	"swap_size": schema.Int64Attribute{
		Description: "When deploying from an Image, this field is optional with a Linode API default of " +
			"512mb, otherwise it is ignored. This is used to set the swap disk size for the newly-created Linode.",
		Computed: true,
	},
	"watchdog_enabled": schema.BoolAttribute{
		Description: "The watchdog, named Lassie, is a Shutdown Watchdog that monitors your Linode and will " +
			"reboot it if it powers off unexpectedly. It works by issuing a boot job when your Linode powers " +
			"off without a shutdown job being responsible. To prevent a loop, Lassie will give up if there have " +
			"been more than 5 boot jobs issued within 15 minutes.",
		Computed: true,
	},
	"host_uuid": schema.StringAttribute{
		Description: "The Linode’s host machine, as a UUID.",
		Computed:    true,
	},
	"has_user_data": schema.BoolAttribute{
		Description: "Whether this Instance was created with user-data.",
		Computed:    true,
	},
	"disk_encryption": schema.StringAttribute{
		Description: "The disk encryption policy for this Instance.",
		Computed:    true,
	},
	"lke_cluster_id": schema.Int64Attribute{
		Description: "If applicable, the ID of the LKE cluster this Instance is a node of.",
		Computed:    true,
	},
	"specs":           instanceSpecsAttribute,
	"alerts":          instanceAlertsAttribute,
	"backups":         instanceBackupsAttribute,
	"config":          instanceConfigAttribute,
	"disk":            instanceDiskAttribute,
	"placement_group": instancePlacementGroupAttribute,
}

var instancesAttribute = schema.ListNestedAttribute{
	Description: "The returned list of Instances.",
	Computed:    true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: instanceAttributes,
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":     filterConfig.OrderSchema(),
		"order_by":  filterConfig.OrderBySchema(),
		"instances": instancesAttribute,
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
