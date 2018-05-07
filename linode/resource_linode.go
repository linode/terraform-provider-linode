package linode

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"time"

	golinode "github.com/chiefy/go-linode"
	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/crypto/sha3"
)

// Instance.Status values
const (
	InstanceBooting      = "booting"
	InstanceRunning      = "running"
	InstanceOffline      = "offline"
	InstanceShuttingDown = "shutting_down"
	InstanceRebooting    = "rebooting"
	InstanceProvisioning = "provisioning"
	InstanceDeleting     = "deleting"
	InstanceMigrating    = "migrating"
	InstanceRebuilding   = "rebuilding"
	InstanceCloning      = "cloning"
	InstanceRestoring    = "restoring"
)

// WaitTimeout is the default number of seconds to wait for Linode instance status changes
const WaitTimeout = 600

var (
	kernelList        []*golinode.LinodeKernel
	kernelListMap     map[string]*golinode.LinodeKernel
	regionList        []*golinode.Region
	regionListMap     map[string]*golinode.LinodeRegion
	typeList          []*golinode.LinodeType
	typeListMap       map[string]*golinode.LinodeType
	latestKernelStrip *regexp.Regexp
)

func init() {
	latestKernelStrip = regexp.MustCompile("\\s*\\(.*\\)\\s*")
}

func resourceLinodeLinode() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeLinodeCreate,
		Read:   resourceLinodeLinodeRead,
		Update: resourceLinodeLinodeUpdate,
		Delete: resourceLinodeLinodeDelete,
		Exists: resourceLinodeLinodeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The image to deploy to the disk.",
				Required:     true,
				ForceNew:     true,
				InputDefault: "linode/debian9",
			},
			"kernel": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The kernel used at boot by the Linode Config. (examples: linode/latest-64bit, linode/grub2, linode/direct-disk)",
				Required:    true,
				Default:     "linode/direct-disk",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name (or label) of the Linode instance.",
				Optional:    true,
			},
			"group": &schema.Schema{
				Removed: "See 'tags'",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The tags to apply to the Linode instance.",
				Optional:    true,
			},
			"region": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The region where this instance will be deployed.",
				Required:     true,
				ForceNew:     true,
				InputDefault: "us-east-1a",
			},
			"size": &schema.Schema{
				Removed: "See 'type'",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The type of instance to be deployed, determining the price and size.",
				Required:    true,
				Default:     "g5-standard-1",
			},
			"status": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The status of the instance, indicating the current readiness state.",
				Computed:    true,
			},
			"plan_storage": &schema.Schema{
				Removed: "See 'storage'",
			},
			"storage": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The total amount of disk space (MB) available to this Linode instance.",
				Computed:    true,
			},
			"plan_storage_utilized": &schema.Schema{
				Removed: "See 'storage_utilized'",
			},
			"storage_utilized": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The total ",
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
				Required:      true,
				ForceNew:      true,
				StateFunc:     ssh_key_state,
				PromoteSingle: true,
			},
			"root_password": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The password that will be initialially assigned to the 'root' user account.",
				Required:    true,
				ForceNew:    true,
				StateFunc:   root_password_state,
			},
			"helper_distro": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Controls the behavior of the Linode Config's Distribution Helper setting.",
				Optional:    true,
				Default:     true,
			},
			"manage_private_ip_automatically": &schema.Schema{
				Removed: "See 'helper_network'",
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
				Description: "Storage (MB) to dedicate to swap disk (memory) space.",
				Optional:    true,
				Default:     512,
			},
		},
	}
}

func resourceLinodeLinodeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*golinode.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Failed to parse Linode instance ID %s as int because %s", d.Id(), err)
	}

	instance, err := client.GetInstance(int(id))
	if err != nil {
		return false, fmt.Errorf("Failed to parse Linode instance ID %s as int because %s", d.Id(), err)
	}
	return true, nil
}

func resourceLinodeLinodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*golinode.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode instance ID %s as int because %s", d.Id(), err)
	}

	instance, err := client.GetInstance(int(id))

	if err != nil {
		return fmt.Errorf("Failed to find the specified Linode instance because %s", err)
	}

	instanceNetwork, err := client.GetInstanceIPAddresses(int(id))

	if err != nil {
		return fmt.Errorf("Failed to get the IPs for Linode instance %s because %s", d.Id(), err)
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

	d.Set("name", instance.Label)
	d.Set("status", instance.Status)
	d.Set("type", instance.Type)
	d.Set("region", instance.Region)

	d.Set("group", instance.Group)

	planStorage := instance.Specs.Disk
	d.Set("plan_storage", planStorage)

	instanceDisks, err := client.ListInstanceDisks(int(id), nil)

	if err != nil {
		return fmt.Errorf("Failed to get the disks for the Linode instance %d because %s", id, err)
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

	d.Set("disk_expansion", boolToString(d.Get("disk_expansion").(bool)))

	configs, err := client.ListInstanceConfigs(int(id), nil)
	if err != nil {
		return fmt.Errorf("Failed to get the config for Linode instance %d (%s) because %s", instance.ID, instance.Label, err)
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
		return fmt.Errorf("Failed to get the image because %s", err)
	}
	d.Set("image", image)
	d.SetPartial("image")
	**/

	d.Set("helper_distro", boolToString(config.Helpers.Distro))
	d.Set("helper_network", boolToString(config.Helpers.Network))
	d.Set("kernel", config.Kernel)

	return nil
}

func resourceLinodeLinodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*golinode.Client)
	d.Partial(true)

	region, err := getRegion(client, d.Get("region").(string))
	if err != nil {
		return fmt.Errorf("Failed to locate region %s because %s", d.Get("region").(string), err)
	}

	linodetype, err := getType(client, d.Get("type").(string))
	if err != nil {
		return fmt.Errorf("Failed to find a Linode type %s because %s", d.Get("type"), err)
	}
	createOpts := golinode.InstanceCreateOptions{
		Region: region.ID,
		Type:   linodetype.ID,
		Label:  d.Get("name").(string),
		Group:  d.Get("group").(string),
	}
	instance, err := client.CreateInstance(&createOpts)
	if err != nil {
		return fmt.Errorf("Failed to create a Linode instance in region %s of type %d because %s", d.Get("region"), d.Get("type"), err)
	}
	d.SetId(fmt.Sprintf("%d", instance.ID))
	d.Set("name", instance.Label)

	d.SetPartial("region")
	d.SetPartial("type")
	d.SetPartial("name")
	d.SetPartial("group")

	swapSize := 0
	var swapDisk *golinode.InstanceDisk

	// Create the Swap Partition
	if swapSize = d.Get("swap_size").(int); swapSize > 0 {
		swapOpts := golinode.InstanceDiskCreateOptions{
			Filesystem: "swap",
			Size:       swapSize,
		}

		swapDisk, err = client.CreateInstanceDisk(instance.ID, swapOpts)

		if err != nil {
			return fmt.Errorf("Failed to create Linode instance %d swap disk because %s", instance.ID, err)
		}

	}
	d.SetPartial("swap_size")

	// Create the storage Partition
	var storageSize int

	if storageSize = d.Get("storage_size").(int); storageSize != 0 {
		storageSize = instance.Specs.Disk - d.Get("swap_size").(int)
	}

	diskOpts := golinode.InstanceDiskCreateOptions{
		Filesystem: "ext4",
		Size:       storageSize,
	}

	if image, ok := d.GetOk("image"); ok {
		diskOpts.Image = image.(string)

		diskOpts.RootPass = d.Get("root_pass").(string)

		if sshKeys, ok := d.GetOk("ssh_authorized_keys"); ok {
			diskOpts.AuthorizedKeys = sshKeys.([]string)
		}

	}

	storageDisk, err := client.CreateInstanceDisk(instance.ID, diskOpts)
	d.SetPartial("image")
	d.SetPartial("root_pass")
	d.SetPartial("ssh_authorized_keys")

	if err != nil {
		return fmt.Errorf("Failed to create Linode instance %d disk because %s", instance.ID, err)
	}
	d.SetPartial("storage_size")

	if d.Get("private_networking").(bool) {
		resp, err := client.AddInstanceIPAddress(instance.ID, false)
		if err != nil {
			return fmt.Errorf("Failed to add a private ip address to Linode instance %d because %s", instance.ID, err)
		}
		d.Set("private_ip_address", resp.Address)
		d.SetPartial("private_ip_address")
	}

	configDevices := &golinode.InstanceConfigDeviceMap{
		SDA: &golinode.InstanceConfigDevice{DiskID: storageDisk.ID},
	}

	if swapDisk.ID > 0 {
		configDevices.SDB = &golinode.InstanceConfigDevice{DiskID: swapDisk.ID}
	}

	configOpts := golinode.InstanceConfigCreateOptions{
		Kernel: d.Get("kernel").(string),
		Helpers: &golinode.InstanceConfigHelpers{
			Distro:  d.Get("helper_distro").(bool),
			Network: d.Get("helper_network").(bool),
		},
		Devices:    configDevices,
		RootDevice: "sda",
	}

	config, err := client.CreateInstanceConfig(instance.ID, configOpts)
	d.SetPartial("helper_network")
	d.SetPartial("helper_distro")

	client.BootInstance(instance.ID, config.ID)

	d.Partial(false)
	if err = waitForInstanceStatus(client, instance.ID, InstanceRunning, WaitTimeout); err != nil {
		return fmt.Errorf("Timed-out waiting for Linode instance %d to boot because %s", instance.ID, err)
	}

	return resourceLinodeLinodeRead(d, meta)
}

func resourceLinodeLinodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*golinode.Client)
	d.Partial(true)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse linode id %s as an int because %s", d.Id(), err)
	}

	instance, err := client.GetInstance(int(id))
	if err != nil {
		return fmt.Errorf("Failed to fetch data about the current linode because %s", err)
	}

	if d.HasChange("name") {
		if instance, err = client.RenameInstance(instance.ID, d.Get("name").(string)); err != nil {
			return err
		}
	}

	var ok bool
	if d.HasChange("type") {
		if ok, err = client.ResizeInstance(instance.ID, d.Get("type").(string)); err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("Failed resizing linode %d because %s", instance.ID, err)
		}
		if err = waitForInstanceStatus(client, instance.ID, InstanceOffline, WaitTimeout); err != nil {
			return fmt.Errorf("Failed while waiting for linode %d to finish resizing because %s", instance.ID, err)
		}
	}

	if d.HasChange("private_networking") {
		if !d.Get("private_networking").(bool) {
			return fmt.Errorf("Can't deactivate private networking for linode %s", d.Id())
		}

		resp, err := client.AddInstanceIPAddress(int(id), false)

		if err != nil {
			return fmt.Errorf("Failed to activate private networking on linode %s because %s", d.Id(), err)
		}
		d.SetPartial("private_networking")
		d.Set("private_ip_address", resp.Address)
		d.SetPartial("private_ip_address")
	}

	configs, err := client.ListInstanceConfigs(int(id), nil)
	if err != nil {
		return fmt.Errorf("Failed to fetch the config for linode %d because %s", id, err)
	}
	if len(configs) != 1 {
		return fmt.Errorf("Linode %d has an incorrect number of configs %d, this plugin can only handle 1", id, len(configResp.LinodeConfigs))
	}
	config := configs[0]
	config.(InstanceConfigUpdateOptions)
	if err = update(client, config, d); err != nil {
		return fmt.Errorf("Failed to update Linode %d config because %s", id, err)
	}

	if d.Get("helper_network").(bool) {
		_, err = client.RebootInstance(int(id), 0)
		if err != nil {
			return fmt.Errorf("Failed to reboot linode %s because %s", d.Id(), err)
		}
		err = waitForJobsToComplete(client, int(id))
		if err != nil {
			return fmt.Errorf("Failed while waiting for linode %s to finish rebooting because %s", d.Id(), err)
		}
	}

	return resourceLinodeLinodeRead(d, meta)
}

func resourceLinodeLinodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*golinode.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse linode id %d as int", d.Id())
	}
	err = client.DeleteInstance(int(id))
	if err != nil {
		return fmt.Errorf("Failed to delete Linode instance %d because %s", id, err)
	}
	return nil
}

// getKernel gets the kernel from the id of the kernel
func getKernel(client *golinode.Client, kernelID string) (golinode.LinodeKernel, error) {
	if KernelList == nil {
		if err := getKernelList(client); err != nil {
			return nil, err
		}
	}

	if t, ok := KernelListMap[kernelID]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("Unabled to find Linode Kernel %s", kernelID)
}

// getKernelList populates kernelList with the available kernels. kernelList is used to reduce the number of api
// requests required as it is unlikely that the available kernels will change during a single terraform run.
func getKernelList(client *golinode.Client) error {
	var err error
	if kernelList == nil {
		if kernelList, err = client.ListKernels(nil); err != nil {
			return err
		}

		for t := range kernelList {
			kernelListMap[t.ID] = t
		}
	}

	return nil
}

// getRegion gets the region from the id of the region
func getRegion(client *golinode.Client, regionID string) (golinode.Region, error) {
	if RegionList == nil {
		if err := getRegionList(client); err != nil {
			return nil, err
		}
	}

	if t, ok := RegionListMap[regionID]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("Unabled to find Linode Region %s", regionID)
}

// getRegionList populates regionList with the available regions. regionList is used to reduce the number of api
// requests required as it is unlikely that the available regions will change during a single terraform run.
func getRegionList(client *golinode.Client) error {
	if regionList == nil {
		if regionList, err = client.ListRegions(nil); err != nil {
			return err
		}

		for t := range regionList {
			regionListMap[t.ID] = t
		}
	}

	return nil
}

// getType gets the amount of ram from the plan id
func getType(client *golinode.Client, typeID string) (golinode.LinodeType, error) {
	if typeList == nil {
		if err := getTypeList(client); err != nil {
			return nil, err
		}
	}

	if t, ok := typeListMap[typeID]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("Unabled to find Linode Type %s", typeID)
}

// getTypeList populates typeList and typeListMap. typeList is used to reduce
//  the number of api requests required as its unlikely that
// the plans will change during a single terraform run.
func getTypeList(client *golinode.Client) error {
	if typeList == nil {
		typeList, err = client.ListTypes(nil)

		if err != nil {
			return err
		}

		for t := range types {
			typeListMap[t.ID] = t
		}
	}

	return nil
}

// getTotalDiskSize returns the number of disks and their total size.
func getTotalDiskSize(client *golinode.Client, linodeID int) (int, error) {
	var totalDiskSize int
	diskList, err := client.Disk.List(linodeID, -1)
	if err != nil {
		return -1, err
	}

	totalDiskSize = 0
	disks := diskList.Disks
	for i := range disks {
		// Calculate Total Disk Size
		totalDiskSize = totalDiskSize + disks[i].Size
	}

	return totalDiskSize, nil
}

// getBiggestDisk returns the ID and Size of the largest disk attached to the Linode
func getBiggestDisk(client *golinode.Client, linodeID int) (biggestDiskID int, biggestDiskSize int, err error) {
	// Retrieve the Linode's list of disks
	diskList, err := client.Disk.List(linodeID, -1)
	if err != nil {
		return -1, -1, err
	}

	biggestDiskID = 0
	biggestDiskSize = 0
	disks := diskList.Disks
	for i := range disks {
		// Find Biggest Disk ID & Size
		if disks[i].Size > biggestDiskSize {
			biggestDiskID = disks[i].DiskId
			biggestDiskSize = disks[i].Size
		}
	}
	return biggestDiskID, biggestDiskSize, nil
}

// getIps gets the ips assigned to the linode
func getIps(client *golinode.Client, linodeId int) (publicIp string, privateIp string, err error) {
	resp, err := client.Ip.List(linodeId, -1)
	if err != nil {
		return "", "", err
	}
	ips := resp.FullIPAddresses
	for i := range ips {
		if ips[i].IsPublic == 1 {
			publicIp = ips[i].IPAddress
		} else {
			privateIp = ips[i].IPAddress
		}
	}

	return publicIp, privateIp, nil
}

// ssh_key_state hashes a string passed in as an interface
func ssh_key_state(val interface{}) string {
	return hash_string(strings.join(val.([]string), "\n"))
}

// root_password_state hashes a string passed in as an interface
func root_password_state(val interface{}) string {
	return hash_string(val.(string))
}

// hash_string hashes a string
func hash_string(key string) string {
	hash := sha3.Sum512([]byte(key))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// waitForInstanceStatus waits for the Linode instance to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func waitForInstanceStatus(client *golinode.Client, linodeID int, status string, timeoutSeconds int) error {
	start := time.Now()
	for {
		linode, err := client.GetInstance(linodeID)
		if err != nil {
			return err
		}
		complete := (linode.Status == status)

		if complete {
			return nil
		}

		time.Sleep(1 * time.Second)
		if time.Since(start) > timeoutSeconds*time.Second {
			return fmt.Errorf("Linode %d didn't reach '%s' status in %d seconds", linodeId, status, timeoutSeconds)
		}
	}
}

// changeLinodeSize resizes the current linode
func changeLinodeSize(client *golinode.Client, linode golinode.Instance, d *schema.ResourceData) error {
	var newPlanID int
	var waitMinutes int

	// Get the Linode Plan Size
	sizeID, err := getSizeId(client, d.Get("size").(int))
	if err != nil {
		return fmt.Errorf("Failed to find a Plan with %d RAM because %s", d.Get("size"), err)
	}
	newPlanID = sizeId

	// Check if we can safely resize, with Disk Size considered
	currentDiskSize, err := getTotalDiskSize(client, linode.LinodeId)
	newDiskSize, err := getPlanDiskSize(client, newPlanID)
	if currentDiskSize > newDiskSize {
		return fmt.Errorf("Cannot resize linode %d because currentDisk(%d GB) is bigger than newDisk(%d GB)", linode.LinodeId, currentDiskSize, newDiskSize)
	}

	// Resize the Linode
	client.Linode.Resize(linode.LinodeId, newPlanID)
	// Linode says 1-3 minutes per gigabyte for Resize time... Let's be safe with 3
	waitMinutes = ((linode.TotalHD / 1024) * 3)
	// Wait for the Resize Operation to Complete
	err = waitForInstanceStatus(client, linode.LinodeId, "offline", WAIT_TIMEOUT)
	if err != nil {
		return fmt.Errorf("Failed to wait for linode %d resize because %s", linode.LinodeId, err)
	}

	if d.Get("disk_expansion").(bool) {
		// Determine the biggestDisk ID and Size
		biggestDiskID, biggestDiskSize, err := getBiggestDisk(client, linode.LinodeId)
		if err != nil {
			return err
		}
		// Calculate new size, with other disks taken into consideration
		expandedDiskSize := (newDiskSize - (currentDiskSize - biggestDiskSize))

		// Resize the Disk
		client.Disk.Resize(linode.LinodeId, biggestDiskID, expandedDiskSize)
		// Wait for the Disk Resize Operation to Complete
		err = waitForJobsToCompleteWaitMinutes(client, linode.LinodeId, waitMinutes)
		if err != nil {
			return fmt.Errorf("Failed to wait for resize of Disk %d for Linode %d because %s", biggestDiskID, linode.LinodeId, err)
		}
	}

	// Boot up the resized Linode
	client.Linode.Boot(linode.LinodeId, 0)

	// Return the new Linode size
	d.SetPartial("size")
	return nil
}

// changeLinodeConfig changes Config level settings. This is things like the various helpers
func changeLinodeConfig(client *golinode.Client, config golinode.InstanceConfig, d *schema.ResourceData) error {
	updates := make(map[string]string)
	if d.HasChange("helper_distro") {
		updates["helper_distro"] = boolToString(d.Get("helper_distro").(bool))
	}
	if d.HasChange("manage_private_ip_automatically") {
		updates["helper_network"] = boolToString(d.Get("manage_private_ip_automatically").(bool))
	}

	if len(updates) > 0 {
		_, err := client.Config.Update(config.ConfigId, config.LinodeId, config.KernelId, updates)
		if err != nil {
			return fmt.Errorf("Failed to update the linode's config because %s", err)
		}
	}
	d.SetPartial("helper_distro")
	d.SetPartial("manage_private_ip_automatically")
	return nil
}

// Converts a bool to a string
func boolToString(val bool) string {
	if val {
		return "true"
	}
	return "false"
}
