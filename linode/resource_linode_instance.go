package linode

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:         schema.TypeString,
				Description:  "The image to deploy to the disk.",
				Optional:     true,
				ForceNew:     true,
				InputDefault: "linode/debian9",
			},
			"kernel": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The kernel used at boot by the Linode Config. (examples: linode/latest-64bit, linode/grub2, linode/direct-disk)",
				Optional:     true,
				InputDefault: "linode/grub2",
				Computed:     true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true, // @TODO seems a bug that Optional is required when Removed is set
				Removed:  "See 'label'",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The label of the Linode instance.",
				Optional:    true,
				Computed:    true,
			},
			"group": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display group of the Linode instance.",
				Optional:    true,
			},
			"tags": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The tags to apply to the Linode instance.",
				Optional:    true,
			},
			"region": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The region where this instance will be deployed.",
				Required:     true,
				ForceNew:     true,
				InputDefault: "us-east",
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true, // @TODO seems a bug that Optional is required when Removed is set
				Removed:  "See 'type'",
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
			"plan_storage": &schema.Schema{
				Type: schema.TypeInt,
				// Optional: true, // @TODO seems a bug that Optional is required when Removed is set
				Computed: true,
				Removed:  "See 'storage'",
			},
			"storage": &schema.Schema{
				Type: schema.TypeInt,
				// Optional:    true, // @TODO seems a bug that Optional is required when Removed is set
				Computed:    true,
				Description: "The total amount of local disk space (MB) available to this Linode instance.",
			},
			"plan_storage_utilized": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				// Optional: true, // @TODO seems a bug that Optional is required when Removed is set
				Removed: "See 'storage_utilized'",
			},
			"storage_utilized": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The total amount of local disk space (MB) utilized by this Linode instance.",
				Computed:    true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_networking": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"private_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_key": &schema.Schema{
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "The public keys to be used for accessing the root account via ssh.",
				Optional:      true,
				ForceNew:      true,
				StateFunc:     sshKeyState,
				PromoteSingle: true,
			},
			"root_password": &schema.Schema{
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
			"manage_private_ip_automatically": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true, // @TODO seems a bug that Optional is required when Removed is set
				Removed:  "See 'helper_network'",
			},
			"helper_network": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Controls the behavior of the Linode Config's Network Helper setting, used to automatically configure additional IP addresses assigned to this instance.",
				Optional:    true,
				Default:     true,
			},
			"disk_expansion": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Controls the Linode Terraform provider's behavior of resizing the disk to full size after resizing to a larger Linode type.",
				Optional:    true,
				Default:     false,
			},
			"swap_size": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Storage (MB) to dedicate to local swap disk (memory) space.",
				Optional:    true,
				Default:     512,
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
			d.SetId("")
			return false, nil
		}

		return false, fmt.Errorf("Error parsing Linode instance ID %s as int: %s", d.Id(), err)
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

	public, private := instanceNetwork.IPv4.Public, instanceNetwork.IPv4.Private

	if len(public) > 0 {
		d.Set("ip_address", public[0].Address)

		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": public[0].Address,
		})
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

	planStorage := instance.Specs.Disk
	d.Set("plan_storage", planStorage)
	d.Set("storage", planStorage)

	instanceDisks, err := client.ListInstanceDisks(context.Background(), int(id), nil)

	if err != nil {
		return fmt.Errorf("Error getting the disks for the Linode instance %d: %s", id, err)
	}

	planStorageUtilized := 0
	swapSize := 0

	for _, disk := range instanceDisks {
		planStorageUtilized += disk.Size
		// Determine if swap exists and the size.  If it does not exist, swap_size=0
		if disk.Filesystem == "swap" {
			swapSize = disk.Size
			d.Set("swap_size", swapSize)
		}
	}

	d.Set("plan_storage_utilized", planStorageUtilized)
	d.Set("storage_utilized", planStorageUtilized)

	//diskExpansion := d.Get("disk_expansion").(bool)
	//d.Set("disk_expansion", diskExpansion)

	configs, err := client.ListInstanceConfigs(context.Background(), int(id), nil)
	if err != nil {
		return fmt.Errorf("Error getting the config for Linode instance %d (%s): %s", instance.ID, instance.Label, err)
	} else if len(configs) != 1 {
		return nil
	}
	config := configs[0]

	/**
	// This doesn't really tell us much.  This will flunk if an ImageName is used to deploy, since getImage will return
	// an imageID.  Trying to derive the imageName from an imageID could be bad if the image happens to be deleted, which would
	// likely occur in an environment where base image's are lifecycled.
	image, err := getImage(client, int(id))
	if err != nil {
		return fmt.Errorf("Error getting the image: %s", err)
	}
	d.Set("image", image)
	d.SetPartial("image")
	**/

	d.Set("helper_distro", config.Helpers.Distro)
	d.Set("helper_network", config.Helpers.Network)
	d.Set("kernel", config.Kernel)

	return nil
}

func resourceLinodeInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	waitSeconds := 180
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode Instance")
	}
	d.Partial(true)

	/**
	// we used to translate these, now we expect the linode api ids
	region, err := getRegion(&client, d.Get("region").(string))
	if err != nil {
		return fmt.Errorf("Error locating region %s: %s", d.Get("region").(string), err)
	}

	linodetype, err := getType(&client, d.Get("type").(string))
	if err != nil {
		return fmt.Errorf("Error finding a Linode type %s: %s", d.Get("type"), err)
	}
	**/

	createOpts := linodego.InstanceCreateOptions{
		// Region: region.ID,
		// Type:   linodetype.ID,
		Region: d.Get("region").(string),
		Type:   d.Get("type").(string),
		Label:  d.Get("label").(string),
		Group:  d.Get("group").(string),
	}
	instance, err := client.CreateInstance(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("Error creating a Linode instance in region %s of type %s: %s", d.Get("region"), d.Get("type"), err)
	}
	d.SetId(fmt.Sprintf("%d", instance.ID))
	d.Set("label", instance.Label)

	d.SetPartial("region")
	d.SetPartial("type")
	d.SetPartial("label")
	d.SetPartial("group")

	swapSize := 0
	var swapDisk *linodego.InstanceDisk

	// Create the Swap Partition
	_, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeCreate, *instance.Created, waitSeconds)
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

		_, err := client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionDiskCreate, swapDisk.Created, waitSeconds)
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

		diskOpts.RootPass = d.Get("root_password").(string)

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

	_, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionDiskCreate, storageDisk.Created, waitSeconds)
	if err != nil {
		return fmt.Errorf("Error waiting for Linode instance %d root disk: %s", storageDisk.ID, err)
	}

	d.SetPartial("image")
	d.SetPartial("root_password")
	d.SetPartial("ssh_key")

	if err != nil {
		return fmt.Errorf("Error creating Linode instance %d disk: %s", instance.ID, err)
	}
	d.SetPartial("storage")

	if d.Get("private_networking").(bool) {
		resp, err := client.AddInstanceIPAddress(context.Background(), instance.ID, false)
		if err != nil {
			return fmt.Errorf("Error adding a private ip address to Linode instance %d: %s", instance.ID, err)
		}
		d.Set("private_ip_address", resp.Address)
		d.SetPartial("private_ip_address")
	}

	configDevices := &linodego.InstanceConfigDeviceMap{
		SDA: &linodego.InstanceConfigDevice{DiskID: storageDisk.ID},
	}

	if swapDisk.ID > 0 {
		configDevices.SDB = &linodego.InstanceConfigDevice{DiskID: swapDisk.ID}
	}

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

	d.SetPartial("helper_network")
	d.SetPartial("helper_distro")

	booted, err := client.BootInstance(context.Background(), instance.ID, config.ID)
	if !booted {
		return fmt.Errorf("Error booting Linode instance %d: %s", instance.ID, err)
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

// getKernel gets the kernel from the id of the kernel
func getKernel(client *linodego.Client, kernelID string) (*linodego.LinodeKernel, error) {
	if kernelList == nil {
		if err := getKernelList(client); err != nil {
			return nil, err
		}
	}

	if t, ok := kernelListMap[kernelID]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("Unable to find Linode Kernel %s", kernelID)
}

// getKernelList populates kernelList with the available kernels. kernelList is used to reduce the number of api
// requests required as it is unlikely that the available kernels will change during a single terraform run.
func getKernelList(client *linodego.Client) error {
	var err error
	if kernelList == nil {
		if kernelList, err = client.ListKernels(context.Background(), nil); err != nil {
			return err
		}

		kernelListMap = make(map[string]*linodego.LinodeKernel)
		for _, t := range kernelList {
			kernelListMap[t.ID] = t
		}
	}

	return nil
}

// getRegion gets the region from the id of the region
func getRegion(client *linodego.Client, regionID string) (*linodego.Region, error) {
	if regionList == nil {
		if err := getRegionList(client); err != nil {
			return nil, err
		}
	}

	if t, ok := regionListMap[regionID]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("Unable to find Linode Region %s", regionID)
}

// getRegionList populates regionList with the available regions. regionList is used to reduce the number of api
// requests required as it is unlikely that the available regions will change during a single terraform run.
func getRegionList(client *linodego.Client) error {
	if regionList == nil {
		var err error
		if regionList, err = client.ListRegions(context.Background(), nil); err != nil {
			return err
		}

		regionListMap = make(map[string]*linodego.Region)
		for _, t := range regionList {
			regionListMap[t.ID] = t
		}
	}

	return nil
}

// getType gets the amount of ram from the plan id
func getType(client *linodego.Client, typeID string) (*linodego.LinodeType, error) {
	if typeList == nil {
		if err := getTypeList(client); err != nil {
			return nil, err
		}
	}

	if t, ok := typeListMap[typeID]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("Unable to find Linode Type %s", typeID)
}

// getTypeList populates typeList and typeListMap. typeList is used to reduce
//  the number of api requests required as its unlikely that
// the plans will change during a single terraform run.
func getTypeList(client *linodego.Client) error {
	if typeList == nil {
		var err error
		typeList, err = client.ListTypes(context.Background(), nil)

		if err != nil {
			return err
		}
		typeListMap = make(map[string]*linodego.LinodeType)
		for _, t := range typeList {
			typeListMap[t.ID] = t
		}
	}

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

	targetType, err := getType(client, typeID)
	if err != nil {
		return fmt.Errorf("Error finding the instance type %s", typeID)
	}

	//biggestDiskID, biggestDiskSize, err := getBiggestDisk(client, instance.ID)

	//currentDiskSize, err := getTotalDiskSize(client, instance.ID)

	if ok, err := client.ResizeInstance(context.Background(), instance.ID, typeID); err != nil || !ok {
		return fmt.Errorf("Error resizing instance %d: %s", instance.ID, err)
	}

	// Linode says 1-3 minutes per gigabyte for Resize time... Let's be safe with 3
	// This delay should be expected for both the host migration of the Linode,
	// and the filesystem expansion
	waitSeconds := ((instance.Specs.Disk / 1024) * 180)

	event, err := client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeResize, *instance.Created, waitSeconds)
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
		event, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionDiskResize, *event.Created, waitSeconds)
		if err != nil {
			return fmt.Errorf("Error waiting for resize of Disk %d for Linode %d: %s", biggestDiskID, instance.ID, err)
		}
	}

	// Return the new Linode size
	d.SetPartial("disk_expansion")
	d.SetPartial("type")
	return nil
}
