package linode

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/linode/linodego"
	"golang.org/x/crypto/sha3"
)

var (
	kernelList    []*linodego.LinodeKernel
	kernelListMap map[string]*linodego.LinodeKernel
	regionList    []*linodego.Region
	regionListMap map[string]*linodego.Region
	typeList      []*linodego.LinodeType
	typeListMap   map[string]*linodego.LinodeType
)

var (
	boolFalse = false
	boolTrue  = true
)

func resourceLinodeInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeInstanceCreate,
		Read:   resourceLinodeInstanceRead,
		Update: resourceLinodeInstanceUpdate,
		Delete: resourceLinodeInstanceDelete,
		Exists: resourceLinodeInstanceExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image": &schema.Schema{
				Type:        schema.TypeString,
				Description: "An Image ID to deploy the Disk from. Official Linode Images start with linode/, while your Images start with private/. See /images for more information on the Images available for you to use.",
				Required:    true,
				ForceNew:    true,
			},
			"backup_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "A Backup ID from another Linode's available backups. Your User must have read_write access to that Linode, the Backup must have a status of successful, and the Linode must be deployed to the same region as the Backup. See /linode/instances/{linodeId}/backups for a Linode's available backups. This field and the image field are mutually exclusive.",
				Optional:    true,
				ForceNew:    true,
			},
			"stackscript_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The StackScript to deploy to the newly created Linode. If provided, 'image' must also be provided, and must be an Image that is compatible with this StackScript.",
				Optional:    true,
				ForceNew:    true,
			},
			"stackscript_data": &schema.Schema{
				Type:        schema.TypeMap,
				Elem:          &schema.Schema{Type: schema.TypeString},

				Description: "An object containing responses to any User Defined Fields present in the StackScript being deployed to this Linode. Only accepted if 'stackscript_id' is given. The required values depend on the StackScript being deployed.",
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The Linode's label is for display purposes only. If no label is provided for a Linode, a default will be assigned",
				Optional:    true,
				Computed:    true,
			},
			"group": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display group of the Linode instance.",
				Optional:    true,
			},
			"region": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "This is the location where the Linode was deployed. This cannot be changed without opening a support ticket.",
				Required:     true,
				ForceNew:     true,
				InputDefault: "us-east",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The type of instance to be deployed, determining the price and size.",
				Optional:    true,
				Default:     "g6-standard-1",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The status of the instance, indicating the current readiness state.",
				Computed:    true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6": &schema.Schema{
				Type:        schema.TypeString,
				Description: "This Linode's IPv6 SLAAC addresses. This address is specific to a Linode, and may not be shared.",
				Computed:    true,
			},
			
			"ipv4": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "This Linode's IPv4 Addresses. Each Linode is assigned a single public IPv4 address upon creation, and may get a single private IPv4 address if needed. You may need to open a support ticket to get additional IPv4 addresses.",
				Computed:    true,
			},
		
			"private_ip": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If true, the created Linode will have private networking enabled, allowing use of the 192.168.0.0/17 network within the Linode's region.",
				Optional:    true,
			},
			"private_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorized_keys": &schema.Schema{
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "A list of SSH public keys to deploy for the root user on the newly created Linode. Only accepted if 'image' is provided.",
				Optional:      true,
				ForceNew:      true,
				StateFunc:     sshKeyState,
				PromoteSingle: true,
			},
			"root_pass": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The password that will be initialially assigned to the 'root' user account.",
				Sensitive:   true,
				Required:    true,
				ForceNew:    true,
				StateFunc:   rootPasswordState,
			},
			"helper_distro": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Controls the behavior of the Linode Config's Distribution Helper setting.",
				Optional:    true,
				Default:     true,
			},
			"helper_network": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Controls the behavior of the Linode Config's Network Helper setting, used to automatically configure additional IP addresses assigned to this instance.",
				Optional:    true,
				Default:     true,
			},
			"swap_size": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "When deploying from an Image, this field is optional, otherwise it is ignored. This is used to set the swap disk size for the newly-created Linode.",
				Optional:    true,
				Default:     512,
			},
			"backups_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If this field is set to true, the created Linode will automatically be enrolled in the Linode Backup service. This will incur an additional charge. The cost for the Backup service is dependent on the Type of Linode deployed.",
				Optional:    true,
				Default:     false,
			},
			"watchdog_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "The watchdog, named Lassie, is a Shutdown Watchdog that monitors your Linode and will reboot it if it powers off unexpectedly. It works by issuing a boot job when your Linode powers off without a shutdown job being responsible. To prevent a loop, Lassie will give up if there have been more than 5 boot jobs issued within 15 minutes.",
				Optional:    true,
				Default:     false,
			},
						

			"specs": &schema.Schema{
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"transfer": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"alerts": &schema.Schema{
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"network_in": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"network_out": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"transfer_quota": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"io": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"backups": &schema.Schema{
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"schedule": {
							Type: schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:     schema.TypeString,
										Required: true,
									},
									"window": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
							Required: true,
						},
					},
				},
			},
			"config": &schema.Schema{
				Computed:      true,
				PromoteSingle: true,
				Optional:      true,
				Type:          schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Required: true,
						},
						"helpers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"updatedb_disabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"distro": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"modules_dep": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"network": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"devtmpfs_automount": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"devices": {
							Type: schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sda": {
										Type: schema.TypeSet,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"volume_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										Required: true,
									},
									"sdb": {
										Type: schema.TypeSet,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"volume_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										Optional: true,
									},
									"sdc": {
										Type: schema.TypeSet,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"volume_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										Optional: true,
									},
									"sdd": {
										Type: schema.TypeSet,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"volume_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										Optional: true,
									},
									"sde": {
										Type: schema.TypeSet,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"volume_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										Optional: true,
									},
									"sdf": {
										Type: schema.TypeSet,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"volume_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										Optional: true,
									},
									"sdg": {
										Type: schema.TypeSet,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"volume_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										Optional: true,
									},
									"sdh": {
										Type: schema.TypeSet,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"volume_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										Optional: true,
									},
								},
							},
							Optional: true,
						},
						"kernel": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A Kernel ID to boot a Linode with. Default is based on image choice. (examples: linode/latest-64bit, linode/grub2, linode/direct-disk)",
						},
						"run_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines the state of your Linode after booting. Defaults to default.",
							Default:     "default",
						},
						"virt_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Controls the virtualization mode. Defaults to paravirt.",
							Default:     "paravirt",
						},
						"root_device": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The root device to boot. The corresponding disk must be attached.",
							Default:     "/dev/sda",
						},
						"comments": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional field for arbitrary User comments on this Config.",
						},

						"memory_limit": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Defaults to the total RAM of the Linode",
						},
					},
				},
			},
			"disk": &schema.Schema{
				Computed:      true,
				PromoteSingle: true,
				Optional:      true,
				Type:          schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"label": {
							Type:     schema.TypeString,
							Required: true,
						},
						"filesystem": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ext4",
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authorized_keys": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"stackscript_id": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The StackScript to deploy to the newly created Linode. If provided, 'image' must also be provided, and must be an Image that is compatible with this StackScript.",
							Optional:    true,
							ForceNew:    true,
						},
						"stackscript_data": &schema.Schema{
							Type:        schema.TypeMap,
							Description: "An object containing responses to any User Defined Fields present in the StackScript being deployed to this Linode. Only accepted if 'stackscript_id' is given. The required values depend on the StackScript being deployed.",
							Optional:    true,
							ForceNew:    true,
							Sensitive:   true,
						},
						"root_pass": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The password that will be initialially assigned to the 'root' user account.",
							Sensitive:   true,
							Required:    true,
							ForceNew:    true,
							StateFunc:   rootPasswordState,
						},
					},
				},
				
			},
		},
	}
}

func resourceLinodeInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Error parsing Linode instance ID %s as int: %s", d.Id(), err)
	}

	_, err = client.GetInstance(context.Background(), int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			return false, nil
		}

		return false, fmt.Errorf("Error getting Linode Instance %s: %s", d.Id(), err)
	}
	return true, nil
}

func resourceLinodeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode instance ID %s as int: %s", d.Id(), err)
	}

	instance, err := client.GetInstance(context.Background(), int(id))

	if err != nil {
		if lerr, ok := err.(linodego.Error); ok && lerr.Code == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error finding the specified Linode instance: %s", err)
	}

	instanceNetwork, err := client.GetInstanceIPAddresses(context.Background(), int(id))

	if err != nil {
		return fmt.Errorf("Error getting the IPs for Linode instance %s: %s", d.Id(), err)
	}


	d.Set("ipv4", instance.IPv4)
	public, private := instanceNetwork.IPv4.Public, instanceNetwork.IPv4.Private

	if len(public) > 0 {
		d.Set("ip_address", public[0].Address)

		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": public[0].Address,
		})
		// TODO(displague) to determine 'user', need to check disk.image
		// "linode/containerlinux" is "core", else "root"
		// might be better to make this a resource field and avoid lookups
	}

	if len(private) > 0 {
		d.Set("private_networking", true)
		d.Set("private_ip_address", private[0].Address)
	} else {
		d.Set("private_networking", false)
	}

	d.Set("label", instance.Label)
	d.Set("status", instance.Status)
	d.Set("type", instance.Type)
	d.Set("region", instance.Region)

	d.Set("group", instance.Group)

	instanceDisks, err := client.ListInstanceDisks(context.Background(), int(id), nil)

	if err != nil {
		return fmt.Errorf("Error getting the disks for the Linode instance %d: %s", id, err)
	}

	planStorageUtilized := 0
	swapSize := 0
	var disks []map[string[string]]

	for _, disk := range instanceDisks {
		// Determine if swap exists and the size.  If it does not exist, swap_size=0
		if disk.Filesystem == "swap" {
			swapSize += disk.Size
		}
		disks = append(disks, map[string[string]]{
			
			"size": disk.Size,
			"label": disk.Label,
			"filesystem": disk.Filesystem,
			"read_only": disk.ReadOnly,
			"image": disk.Image,
			"authorized_keys": disk.AuthorizedKeys,
			"stackscript_id": disk.StackScriptID,
		}
	}
	d.Set("disks", disks)
	d.Set("swap_size", swapSize)

	configs, err := client.ListInstanceConfigs(context.Background(), int(id), nil)
	if err != nil {
		return fmt.Errorf("Error getting the config for Linode instance %d (%s): %s", instance.ID, instance.Label, err)
	} else if len(configs) != 1 {
		return nil
	}
	config := configs[0]

	return nil
}

func resourceLinodeInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode Instance")
	}
	d.Partial(true)

	bootConfig := 0
	createOpts := linodego.InstanceCreateOptions{
		Region:         d.Get("region").(string),
		Type:           d.Get("type").(string),
		Label:          d.Get("label").(string),
		Group:          d.Get("group").(string),
		BackupsEnabled: d.Get("backups_enabled").(bool),
		PrivateIP:      d.Get("private_ip").(bool),
	}

	_, disksOk := d.GetOk("disk")
	_, configsOk := d.GetOk("config")

	// If we don't have disks and we don't have configs, use the single API call approach
	if !disksOk && !configsOk {
		for _, key := range(d.Get("authorized_keys").([]interface{})) {
			createOpts.AuthorizedKeys = append(createOpts.AuthorizedKeys, key.(string))
		}
		createOpts.RootPass = d.Get("root_pass").(string)
		createOpts.Image = d.Get("image").(string)
		createOpts.BackupID = d.Get("backup_id").(int)
		if swapSize := d.Get("swap_size").(int); swapSize > 0 {
			createOpts.SwapSize = &swapSize
		}

		createOpts.StackScriptID = d.Get("stackscript_id").(int)

		for name, value := range(d.Get("stackscript_data").(map[string]interface{})) {
			createOpts.StackScriptData[name] = value.(string)
		}
	} else {
		createOpts.Booted = &boolFalse // necessary to prepare disks and configs
	}

	instance, err := client.CreateInstance(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("Error creating a Linode Instance: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", instance.ID))

	// d.Set("backups_enabled", instance.BackupsEnabled)

	d.SetPartial("private_ip")
	d.SetPartial("authorized_keys")
	d.SetPartial("root_pass")
	d.SetPartial("kernel")
	d.SetPartial("image")
	d.SetPartial("backup_id")
	d.SetPartial("stackscript_id")
	d.SetPartial("stackscript_data")
	d.SetPartial("swap_size")

	d.Set("ipv4", instance.IPv4)
	d.Set("ipv6", instance.IPv6)

	for _, address := range instance.IPv4 {
		if private := privateIP(*address); private {
			d.Set("private_ip_address", address.String())
		} else {
			d.Set("ip_address", address.String())
		}
	}

	/*
		if d.Get("private_networking").(bool) {
			resp, err := client.AddInstanceIPAddress(context.Background(), instance.ID, false)
			if err != nil {
				return fmt.Errorf("Error adding a private ip address to Linode instance %d: %s", instance.ID, err)
			}
			d.Set("private_ip_address", resp.Address)
			d.SetPartial("private_ip_address")
		}
	*/

	if disksOk {
		swapSize := 0
		var swapDisk *linodego.InstanceDisk

		// Create the Swap Partition
		_, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeCreate, *instance.Created, int(d.Timeout("create").Seconds()))
		if swapSize = d.Get("swap_size").(int); swapSize > 0 {
			swapOpts := linodego.InstanceDiskCreateOptions{
				Label:      "linode" + strconv.Itoa(instance.ID) + "-swap",
				Filesystem: "swap",
				Size:       swapSize,
			}

			swapDisk, err = client.CreateInstanceDisk(context.Background(), instance.ID, swapOpts)

			if err != nil {
				return fmt.Errorf("Error creating Linode instance %d swap disk: %s", instance.ID, err)
			}

			_, err := client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionDiskCreate, swapDisk.Created, int(d.Timeout("create").Seconds()))
			if err != nil {
				return fmt.Errorf("Error waiting for Linode instance %d swap disk: %s", swapDisk.ID, err)
			}

		}
		d.SetPartial("swap_size")

		// Create the storage Partition

		storageSize := d.Get("storage").(int)

		if !ok || storageSize == 0 {
			storageSize = instance.Specs.Disk - d.Get("swap_size").(int)
		}

		diskOpts := linodego.InstanceDiskCreateOptions{
			Label:      "linode" + strconv.Itoa(instance.ID) + "-root",
			Filesystem: "ext4",
			Size:       storageSize,
		}

		if image, ok := d.GetOk("image"); ok {
			diskOpts.Image = image.(string)

			diskOpts.RootPass = d.Get("root_pass").(string)

			if sshKeys, ok := d.GetOk("ssh_key"); ok {
				if sshKeysArr, ok := sshKeys.([]interface{}); ok {
					diskOpts.AuthorizedKeys = make([]string, len(sshKeysArr))
					for k, v := range sshKeys.([]interface{}) {
						if val, ok := v.(string); ok {
							diskOpts.AuthorizedKeys[k] = val
						}
					}
				}
			}
		}

		storageDisk, err := client.CreateInstanceDisk(context.Background(), instance.ID, diskOpts)
		if err != nil {
			return fmt.Errorf("Error creating Linode instance %d root disk: %s", instance.ID, err)
		}

		_, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionDiskCreate, storageDisk.Created, int(d.Timeout("create").Seconds()))
		if err != nil {
			return fmt.Errorf("Error waiting for Linode instance %d root disk: %s", storageDisk.ID, err)
		}

		d.SetPartial("image")
		d.SetPartial("root_pass")
		d.SetPartial("ssh_key")

		if err != nil {
			return fmt.Errorf("Error creating Linode instance %d disk: %s", instance.ID, err)
		}
		d.SetPartial("storage")

		configDevices := &linodego.InstanceConfigDeviceMap{
			SDA: &linodego.InstanceConfigDevice{DiskID: storageDisk.ID},
		}

		if swapDisk.ID > 0 {
			configDevices.SDB = &linodego.InstanceConfigDevice{DiskID: swapDisk.ID}
		}

		if !configsOk {
			configOpts := linodego.InstanceConfigCreateOptions{
				Label:  fmt.Sprintf("linode%d-config", instance.ID),
				Kernel: d.Get("kernel").(string),
				// RootDevice: "/dev/sda",
				// RunLevel:   "default",
				// VirtMode:   "paravirt",
				Helpers: &linodego.InstanceConfigHelpers{
					Distro:  d.Get("helper_distro").(bool),
					Network: d.Get("helper_network").(bool),
				},
				Devices: *configDevices,
			}

			config, err := client.CreateInstanceConfig(context.Background(), instance.ID, configOpts)
			if err != nil {
				return fmt.Errorf("Error creating Linode instance %d config: %s", instance.ID, err)
			}
			bootConfig = config.ID

			d.SetPartial("helper_network")
			d.SetPartial("helper_distro")

		}
	}

	if createOpts.Booted != nil && *createOpts.Booted {
		booted, err := client.BootInstance(context.Background(), instance.ID, bootConfig)
		if !booted {
			return fmt.Errorf("Error booting Linode instance %d: %s", instance.ID, err)
		}
	}

	d.Partial(false)

	if _, err = client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceRunning, int(d.Timeout("create").Seconds())); err != nil {
		return fmt.Errorf("Timed-out waiting for Linode instance %d to boot: %s", instance.ID, err)
	}

	return resourceLinodeInstanceRead(d, meta)
}

func resourceLinodeInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	d.Partial(true)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Instance ID %s as int: %s", d.Id(), err)
	}

	instance, err := client.GetInstance(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Error fetching data about the current linode: %s", err)
	}

	if d.HasChange("label") {
		if instance, err = client.RenameInstance(context.Background(), instance.ID, d.Get("label").(string)); err != nil {
			return err
		}
		d.Set("label", instance.Label)
		d.SetPartial("label")
	}

	rebootInstance := false

	if d.HasChange("type") {
		err = changeLinodeSize(&client, instance, d)
		if err != nil {
			return err
		}
	}

	if d.HasChange("private_networking") {
		if !d.Get("private_networking").(bool) {
			return fmt.Errorf("Can't deactivate private networking for linode %s", d.Id())
		}

		resp, err := client.AddInstanceIPAddress(context.Background(), int(id), false)

		if err != nil {
			return fmt.Errorf("Error activating private networking on linode %s: %s", d.Id(), err)
		}
		d.SetPartial("private_networking")
		d.Set("private_ip_address", resp.Address)
		d.SetPartial("private_ip_address")
		rebootInstance = true
	}

	configs, err := client.ListInstanceConfigs(context.Background(), int(id), nil)
	if err != nil {
		return fmt.Errorf("Error fetching the config for linode %d: %s", id, err)
	}
	if len(configs) != 1 {
		return fmt.Errorf("Linode %d has an incorrect number of configs %d, this plugin can only handle 1", id, len(configs))
	}
	config := configs[0].GetUpdateOptions()
	updateConfig := false
	if d.HasChange("helper_distro") {
		updateConfig = true
		config.Helpers.Distro = d.Get("helper_distro").(bool)
	}
	if d.HasChange("helper_network") {
		updateConfig = true
		config.Helpers.Network = d.Get("helper_network").(bool)
	}
	if d.HasChange("kernel") {
		updateConfig = true
		config.Kernel = d.Get("kernel").(string)
	}

	if updateConfig {
		_, err := client.UpdateInstanceConfig(context.Background(), instance.ID, configs[0].ID, config)
		if err != nil {
			return fmt.Errorf("Error updating Linode %d config: %s", instance.ID, err)
		}
		d.SetPartial("helper_distro")
		d.SetPartial("helper_network")
		d.SetPartial("kernel")

		rebootInstance = true
	}

	if rebootInstance {
		_, err = client.RebootInstance(context.Background(), instance.ID, configs[0].ID)
		if err != nil {
			return fmt.Errorf("Error rebooting Linode instance %d: %s", instance.ID, err)
		}
		_, err = client.WaitForEventFinished(context.Background(), id, linodego.EntityLinode, linodego.ActionLinodeReboot, *instance.Created, int(d.Timeout("create").Seconds()))
		if err != nil {
			return fmt.Errorf("Error waiting for Linode instance %d to finish rebooting: %s", instance.ID, err)
		}
	}

	return nil // resourceLinodeInstanceRead(d, meta)
}

func resourceLinodeInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Instance ID %s as int", d.Id())
	}
	err = client.DeleteInstance(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Error deleting Linode instance %d: %s", id, err)
	}
	d.SetId("")
	return nil
}

// getTotalDiskSize returns the number of disks and their total size.
func getTotalDiskSize(client *linodego.Client, linodeID int) (totalDiskSize int, err error) {
	disks, err := client.ListInstanceDisks(context.Background(), linodeID, nil)
	if err != nil {
		return 0, err
	}

	for _, disk := range disks {
		totalDiskSize += disk.Size
	}

	return totalDiskSize, nil
}

// getBiggestDisk returns the ID and Size of the largest disk attached to the Linode
func getBiggestDisk(client *linodego.Client, linodeID int) (biggestDiskID int, biggestDiskSize int, err error) {
	diskFilter := "{\"+order_by\": \"size\", \"+order\": \"desc\"}"
	disks, err := client.ListInstanceDisks(context.Background(), linodeID, linodego.NewListOptions(1, diskFilter))
	if err != nil {
		return 0, 0, err
	}

	for _, disk := range disks {
		// Find Biggest Disk ID & Size
		if disk.Size > biggestDiskSize {
			biggestDiskID = disk.ID
			biggestDiskSize = disk.Size
		}
	}
	return biggestDiskID, biggestDiskSize, nil
}

// sshKeyState hashes a string passed in as an interface
func sshKeyState(val interface{}) string {
	return hashString(strings.Join(val.([]string), "\n"))
}

// rootPasswordState hashes a string passed in as an interface
func rootPasswordState(val interface{}) string {
	return hashString(val.(string))
}

// hashString hashes a string
func hashString(key string) string {
	hash := sha3.Sum512([]byte(key))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// changeLinodeSize resizes the current linode
func changeLinodeSize(client *linodego.Client, instance *linodego.Instance, d *schema.ResourceData) error {
	typeID, ok := d.Get("type").(string)
	if !ok {
		return fmt.Errorf("Unexpected value for type %v", d.Get("type"))
	}

	targetType, err := client.GetType(context.Background(), typeID)
	if err != nil {
		return fmt.Errorf("Error finding the instance type %s", typeID)
	}

	//biggestDiskID, biggestDiskSize, err := getBiggestDisk(client, instance.ID)

	//currentDiskSize, err := getTotalDiskSize(client, instance.ID)

	if ok, err := client.ResizeInstance(context.Background(), instance.ID, typeID); err != nil || !ok {
		return fmt.Errorf("Error resizing instance %d: %s", instance.ID, err)
	}

	event, err := client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeResize, *instance.Created, int(d.Timeout("update").Seconds()))
	if err != nil {
		return fmt.Errorf("Error waiting for instance %d to finish resizing: %s", instance.ID, err)
	}

	if d.Get("disk_expansion").(bool) && instance.Specs.Disk > targetType.Disk {
		// Determine the biggestDisk ID and Size
		biggestDiskID, biggestDiskSize, err := getBiggestDisk(client, instance.ID)
		if err != nil {
			return err
		}
		// Calculate new size, with other disks taken into consideration
		expandedDiskSize := biggestDiskSize + targetType.Disk - instance.Specs.Disk

		// Resize the Disk
		client.ResizeInstanceDisk(context.Background(), instance.ID, biggestDiskID, expandedDiskSize)

		// Wait for the Disk Resize Operation to Complete
		// waitForEventComplete(client, instance.ID, "linode_resize", waitMinutes)
		event, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionDiskResize, *event.Created, int(d.Timeout("update").Seconds()))
		if err != nil {
			return fmt.Errorf("Error waiting for resize of Disk %d for Linode %d: %s", biggestDiskID, instance.ID, err)
		}
	}

	// Return the new Linode size
	d.SetPartial("disk_expansion")
	d.SetPartial("type")
	return nil
}

// privateIP determines if an IP is for private use (RFC1918)
// https://stackoverflow.com/a/41273687
func privateIP(ip net.IP) bool {
	private := false
	_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
	_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
	_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
	private = private24BitBlock.Contains(ip) || private20BitBlock.Contains(ip) || private16BitBlock.Contains(ip)
	return private
}
