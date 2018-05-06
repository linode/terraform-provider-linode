package linode

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	golinode "github.com/chiefy/go-linode"
	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/crypto/sha3"
)

const (
	LINODE_BOOTING       = "booting"
	LINODE_RUNNING       = "running"
	LINODE_POWERED_OFF   = "offline"
	LINODE_SHUTTING_DOWN = "shutting_down"
	LINODE_REBOOTING     = "rebooting"
	LINODE_PROVISIONING  = "provisioning"
	LINODE_DELETING      = "deleting"
	LINODE_MIGRATING     = "migrating"
	LINODE_REBUILDING    = "rebuilding"
	LINODE_CLONING       = "cloning"
	LINODE_RESTORING     = "restoring"
)

const WAIT_TIMEOUT = 600

var (
	kernelList        *[]golinode.LinodeKernel
	regionList        *[]golinode.Region
	typeList          *[]golinode.LinodeType
	typeListMap       *map[string]golinode.LinodeType
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Description: "The image to deploy to the disk.",
				Required: true,
				ForceNew: true,
				InputDefault: "linode/debian9",
			},
			"kernel": &schema.Schema{
				Type:     schema.TypeString,
				Description: "The kernel used at boot by the Linode Config. (examples: linode/latest-64bit, linode/grub2, linode/direct-disk)",
				Required: true,
				Default: "linode/direct-disk",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Description: "The name (or label) of the Linode instance.",
				Optional: true,
			},
			"group": &schema.Schema{
				Removed: "See 'tags'",
			},
			"tags": &schema.Schema{
				Type:     []schema.TypeString,
				Description: "The tags to apply to the Linode instance.",
				Optional: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Description: "The region where this instance will be deployed.",
				Required: true,
				ForceNew: true,
				InputDefault "us-east-1a",
			},
			"size": &schema.Schema{
				Removed: "See 'type'",
			},
			"type": &schema.Schema{
				Type:     schema.String,
				Description: "The type of instance to be deployed, determining the price and size.",
				Required: true,
				Default: "g5-standard-1",
			},
			"status": &schema.Schema{
				Type:     schema.TypeInt,
				Description: "The status of the instance, indicating the current readiness state.",
				Computed: true,
			},
			"plan_storage": &schema.Schema{
				Removed: "See 'storage'",
			},
			"storage": &schema.Schema{
				Type:     schema.TypeInt,
				Description: "The total amount of disk space (MB) available to this Linode instance.",
				Computed: true,
			},
			"plan_storage_utilized": &schema.Schema{
				Removed: "See 'storage_utilized'",
			},
			"storage_utilized": &schema.Schema{
				Type:     schema.TypeInt,
				Description: "The total ",
				Computed: true,
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
				Type:      schema.TypeString,
				Description: "The public key to be used for accessing the root account via ssh.",
				Required:  true,
				ForceNew:  true,
				StateFunc: ssh_key_state,
			},
			"root_password": &schema.Schema{
				Type:      schema.TypeString,
				Description: "The password that will be initialially assigned to the 'root' user account.",
				Required:  true,
				ForceNew:  true,
				StateFunc: root_password_state,
			},
			"helper_distro": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Controls the behavior of the Linode Config's Distribution Helper setting.",
				Optional: true,
				Default:  true,
			},
			"manage_private_ip_automatically": &schema.Schema{
				Removed: "See 'helper_network'",
			},
			"helper_network": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Controls the behavior of the Linode Config's Network Helper setting, used to automatically configure additional IP addresses assigned to this instance.",
				Optional: true,
				Default:  true,
			},
			"disk_expansion": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Controls the Linode Terraform provider's behavior of resizing the disk to full size after resizing to a larger Linode type.",
				Optional: true,
				Default:  false,
			},
			"swap_size": &schema.Schema{
				Type:     schema.TypeInt,
				Description: "Storage (MB) to dedicate to swap disk (memory) space.",
				Optional: true,
				Default:  512,
			},
		},
	}
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

	instanceNetwork, err := client.GetInstanceIPAddress(id, nil)
	
	if err != nil {
		return fmt.Errorf("Failed to get the IPs for Linode instance %s because %s", d.Id(), err)
	}
	
	public, private := instanceNetwork.IPv4.Public, instanceNetwork.IPv4.Private

	if len(public) > 0 {
		d.Set("ip_address", public[0])

		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": public[0],
		})
	}

	if len(private) > 0 {
		d.Set("private_networking", true)
		d.Set("private_ip_address", private[0])
	} else {
		d.Set("private_networking", false)
	}


	d.Set("name", instance.Label)
	d.Set("status", instance.Status)
	d.Set("type", instance.Type)
	d.Set("region", instance.Region)

	d.Set("group", instance.Group)

	plan_storage := instance.Specs.Disk
	d.Set("plan_storage", plan_storage)

	instanceDisks, err := client.ListInstanceDisks(id, nil)

	if err != nil {
		return fmt.Errorf("Failed to get the disks for the linode because %s", err)
	}

	plan_storage_utilized:=0
	swap_size := 0

	for _, disk := range instanceDisks {
		plan_storage_utilized += disk.Size
		// Determine if swap exists and the size.  If it does not exist, swap_size=0
		if disk.Filesystem == "swap" {
			swap_size = disk.Size
			d.Set("swap_size", swap_size)
		}
	}

	d.Set("plan_storage_utilized", plan_storage_utilized)

	d.Set("disk_expansion", boolToString(d.Get("disk_expansion").(bool)))



	configs, err := client.ListInstanceConfigs(id)
	if err != nil {
		return fmt.Errorf("Failed to get the config for Linode instance %d (%s) because %s", instance.ID, instance.Label, err)
	} else if len(configs.LinodeConfigs) != 1 {
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

	d.Set("helper_distro", boolToString(config.Helpers.Distro.Bool))
	d.Set("helper_network", boolToString(config.Helpers.Network.Bool))
	d.Set("kernel", config.Kernel)

	return nil
}

func resourceLinodeLinodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*golinode.Client)
	d.Partial(true)

	regionId, err := getRegionID(client, d.Get("region").(string))
	if err != nil {
		return fmt.Errorf("Failed to locate region %s because %s", d.Get("region").(string), err)
	}

	typeId, err := getTypeId(client, d.Get("type").(int))
	if err != nil {
		return fmt.Errorf("Failed to find a Lindoe type %s because %s", d.Get("type"), err)
	}
	create, err := client.Linode.Create(regionId, typeId, 1)
	if err != nil {
		return fmt.Errorf("Failed to create a Linode instance in region %s of type %d because %s", d.Get("region"), d.Get("type"), err)
	}

	d.SetId(fmt.Sprintf("%d", create.LinodeId.LinodeId))
	d.SetPartial("region")
	d.SetPartial("type")

	// Create the Swap Partition
	if d.Get("swap_size").(int) > 0 {
		emptyArgs := make(map[string]string)
		_, err = client.Disk.Create(create.LinodeId.LinodeId, "swap", "swap", d.Get("swap_size").(int), emptyArgs)
		if err != nil {
			return fmt.Errorf("Failed to create a swap drive because %s", err)
		}
	}
	d.SetPartial("swap_size")

	// Load the basic data about the current linode
	linodes, err := client.Linode.List(create.LinodeId.LinodeId)
	if err != nil {
		return fmt.Errorf("Failed to load data about the newly created linode because %s", err)
	} else if len(linodes.Linodes) != 1 {
		return fmt.Errorf("An incorrect number of linodes (%d) was returned for id %s", len(linodes.Linodes), d.Id())
	}
	linode := linodes.Linodes[0]

	if err = changeLinodeSettings(client, linode, d); err != nil {
		return err
	}

	if d.Get("private_networking").(bool) {
		resp, err := client.Ip.AddPrivate(linode.LinodeId)
		if err != nil {
			return fmt.Errorf("Failed to add a private ip address to linode %d because %s", linode.LinodeId, err)
		}
		d.Set("private_ip_address", resp.IPAddress.IPAddress)
		d.SetPartial("private_ip_address")
	}
	d.SetPartial("private_networking")

	ssh_key := d.Get("ssh_key").(string)
	password := d.Get("root_password").(string)
	disk_size := (linode.TotalHD - d.Get("swap_size").(int))
	err = deployImage(client, linode, d.Get("image").(string), disk_size, ssh_key, password)
	if err != nil {
		return fmt.Errorf("Failed to create disk for image %s because %s", d.Get("image"), err)
	}

	d.SetPartial("root_password")
	d.SetPartial("ssh_key")

	diskResp, err := client.Disk.List(linode.LinodeId, -1)
	if err != nil {
		return fmt.Errorf("Failed to get the disks for the newly created linode because %s", err)
	}
	var rootDisk int
	var swapDisk int
	for i := range diskResp.Disks {
		if strings.EqualFold(diskResp.Disks[i].Type, "swap") {
			swapDisk = diskResp.Disks[i].DiskId
		} else {
			rootDisk = diskResp.Disks[i].DiskId
		}
	}

	kernelId, err := getKernelID(client, d.Get("kernel").(string))
	if err != nil {
		return fmt.Errorf("Failed to find kernel %s because %s", d.Get("kernel").(string), err)
	}

	confArgs := make(map[string]string)
	if d.Get("manage_private_ip_automatically").(bool) {
		confArgs["helper_network"] = "true"
	} else {
		confArgs["helper_network"] = "false"
	}
	if d.Get("helper_distro").(bool) {
		confArgs["helper_distro"] = "true"
	} else {
		confArgs["helper_distro"] = "false"
	}
	if d.Get("swap_size").(int) > 0 {
		confArgs["DiskList"] = fmt.Sprintf("%d,%d", rootDisk, swapDisk)
	} else {
		confArgs["DiskList"] = fmt.Sprintf("%d", rootDisk)
	}

	confArgs["RootDeviceNum"] = "1"
	c, err := client.Create.Create(linode.LinodeId, kernelId, d.Get("image").(string), confArgs)
	if err != nil {
		log.Printf("diskList: %s", confArgs["DiskList"])
		log.Println(confArgs["DiskList"])
		return fmt.Errorf("Failed to create config for linode %d because %s", linode.LinodeId, err)
	}
	confID := c.LinodeConfigId
	d.SetPartial("image")
	d.SetPartial("manage_private_ip_automatically")
	d.SetPartial("helper_distro")
	client.Linode.Boot(linode.LinodeId, confID.LinodeConfigId)

	d.Partial(false)
	err = waitForJobsToComplete(client, linode.LinodeId)
	if err != nil {
		return fmt.Errorf("Failed to wait for linode %d to boot because %s", linode.LinodeId, err)
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

	linode, err := client.GetInstance(int(id))
	if err != nil {
		return fmt.Errorf("Failed to fetch data about the current linode because %s", err)
	}

	if d.HasChange("name") || d.HasChange("group") {
		if err = changeLinodeSettings(client, linode, d); err != nil {
			return err
		}
	}

	if d.HasChange("type") {
		if err = changeLinodeType(client, linode, d); err != nil {
			return err
		}
		if err = waitForJobsToComplete(client, int(id)); err != nil {
			return fmt.Errorf("Failed while waiting for linode %s to finish resizing because %s", d.Id(), err)
		}
	}

	configResp, err := client.Config.List(int(id), -1)
	if err != nil {
		return fmt.Errorf("Failed to fetch the config for linode %d because %s", id, err)
	}
	if len(configResp.LinodeConfigs) != 1 {
		return fmt.Errorf("Linode %d has an incorrect number of configs %d, this plugin can only handle 1", id, len(configResp.LinodeConfigs))
	}
	config := configResp.LinodeConfigs[0]

	if err = changeLinodeConfig(client, config, d); err != nil {
		return fmt.Errorf("Failed to update Linode %d config because %s", id, err)
	}

	if d.HasChange("private_networking") {
		if !d.Get("private_networking").(bool) {
			return fmt.Errorf("Can't deactivate private networking for linode %s", d.Id())
		} else {
			_, err = client.Ip.AddPrivate(int(id))
			if err != nil {
				return fmt.Errorf("Failed to activate private networking on linode %s because %s", d.Id(), err)
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
	_, err = client.Linode.Delete(int(id), true)
	if err != nil {
		return fmt.Errorf("Failed to delete linode %d because %s", id, err)
	}
	return nil
}

// getDisks gets all of the disks that are attached to the instance. It only returns the names of those disks
func getDisks(client *golinode.Client, id int) ([]string, error) {
	resp, err := client.Disk.List(id, -1)
	if err != nil {
		return []string{}, err
	}
	if len(resp.Disks) != 2 {
		return []string{}, fmt.Errorf("Found %d disks attached to Linode instance %s. This plugin can only handle exactly 2.", len(resp.Disks), err)
	}
	disks := []string{}
	for i := range resp.Disks {
		disks = append(disks, resp.Disks[i].Label.String())
	}
	return disks, nil
}

// getImage Finds out what image was used to create the server.
func getImage(client *golinode.Client, id int) (string, error) {
	disks, err := getDisks(client, id)
	if err != nil {
		return "", err
	}

	// Assumes disk naming convention of Root(LINODEID)__Base(IMAGEID)
	grabId := regexp.MustCompile(`Base\(([0-9]+)\)`)

	for i := range disks {
		// Check if we match the pattern at all
		if grabId.MatchString(disks[i]) {
			// Print out the first group match
			return grabId.FindStringSubmatch(disks[i])[1], nil
		} else if strings.HasSuffix(disks[i], " Disk") {
			// Keep the old method for backward compatibility
			return disks[i][:(len(disks[i]) - 5)], nil
		}
	}
	return "", errors.New("Unable to find the image based on the disk names")
}

// getKernelName gets the name of the kernel from the id.
func getKernelName(client *golinode.Client, kernelId int) (string, error) {
	if kernelList == nil {
		if err := getKernelList(client); err != nil {
			return "", err
		}
	}
	k := *kernelList
	for i := range k {
		if k[i].KernelId == kernelId {
			if strings.HasPrefix(k[i].Label.String(), "Latest") {
				return latestKernelStrip.ReplaceAllString(k[i].Label.String(), ""), nil
			} else {
				return k[i].Label.String(), nil
			}
		}
	}
	return "", fmt.Errorf("Failed to find kernel id %d", kernelId)
}

// getKernelID gets the id of the kernel from the specified id.
func getKernelID(client *golinode.Client, kernelName string) (int, error) {
	if kernelList == nil {
		if err := getKernelList(client); err != nil {
			return -1, err
		}
	}
	k := *kernelList
	for i := range k {
		if strings.HasPrefix(kernelName, "Latest") {
			if strings.HasPrefix(k[i].Label.String(), kernelName) {
				return k[i].KernelId, nil
			}
		} else {
			if k[i].Label.String() == kernelName {
				return k[i].KernelId, nil
			}
		}
	}
	return -1, fmt.Errorf("Failed to find kernel %s", kernelName)
}

// getKernelList populates kernelList with all of the available kernels. kernelList is purely to reduce the number of
// api calls as the available kernels are unlikely to change within a single terraform run.
func getKernelList(client *golinode.Client) error {
	kernels, err := client.Avail.Kernels()
	if err != nil {
		return err
	}
	kernelList = &kernels.Kernels
	return nil
}

// getRegionName gets the region name from the region id
func getRegionName(client *golinode.Client, regionId int) (string, error) {
	if regionList == nil {
		if err := getRegionList(client); err != nil {
			return "", err
		}
	}

	r := *regionList
	for i := range r {
		if r[i].DataCenterId == regionId {
			return r[i].Location, nil
		}
	}
	return "", fmt.Errorf("Failed to find region id %d", regionId)
}

// getRegionID gets the region id from the name of the region
func getRegionID(client *golinode.Client, regionName string) (int, error) {
	if regionList == nil {
		if err := getRegionList(client); err != nil {
			return -1, err
		}
	}

	r := *regionList
	for i := range r {
		if r[i].Location == regionName {
			return r[i].DataCenterId, nil
		}
	}
	return -1, fmt.Errorf("Failed to find the region name %s", regionName)
}

// getRegionList populates regionList with the available regions. regionList is used to reduce the number of api
// requests required as it is unlikely that the available regions will change during a single terraform run.
func getRegionList(client *golinode.Client) error {
	resp, err := client.Avail.DataCenters()
	if err != nil {
		return err
	}
	regionList = &resp.DataCenters
	return nil
}

// getSizeId gets the id from the specified amount of ram
func getTypeId(client *golinode.Client, type int) (int, error) {
	if typeList == nil {
		if err := getTypeList(client); err != nil {
			return -1, err
		}
	}

	s := *typeList
	for i := range s {
		if s[i].RAM == type {
			return s[i].PlanId, nil
		}
	}
	return -1, fmt.Errorf("Unable to locate the plan with RAM %d", type)
}

// getType gets the amount of ram from the plan id
func getType(client *golinode.Client, type string) (LinodeType, error) {
	if typeList == nil {
		if err := getTypeList(client); err != nil {
			return nil, err
		}
	}

	if t, ok := typeListMap[type]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("Unabled to find Linode Type %s", type)
}

// getTypeList populates typeList and typeListMap. typeList is used to reduce
//  the number of api requests required as its unlikely that
// the plans will change during a single terraform run.
func getTypeList(client *golinode.Client) error {
	if typeList == nil {
		types, err := client.ListTypes(nil)
		
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

// root_password_state hashes a string passed in as an interface
func ssh_key_state(val interface{}) string {
	return hash_string(val.(string))
}

// root_password_state hashes a string passed in as an interface
func root_password_state(val interface{}) string {
	return hash_string(val.(string))
}

// hash_string hashes a string
func hash_string(key string) string {
	hash := sha3.Sum256([]byte(key))
	return base64.StdEncoding.EncodeToString(hash[:])
}

const (
	PREBUILT = iota
	CUSTOM_IMAGE
)

// findImage finds the specified image. It checks the prebuilt images first and then any custom images. It returns both
// the image type and the images id
func findImage(client *golinode.Client, imageName string) (imageType, imageId int, err error) {
	// Get Available Distributions
	distResp, err := client.Avail.Distributions()
	if err != nil {
		return -1, -1, err
	}
	prebuilt := distResp.Distributions
	for i := range prebuilt {
		if prebuilt[i].Label.String() == imageName {
			return PREBUILT, prebuilt[i].DistributionId, nil
		}
	}

	// Get Available Client Images
	custResp, err := client.Image.List()
	if err != nil {
		return -1, -1, err
	}
	customImages := custResp.Images
	for i := range customImages {
		if customImages[i].Label.String() == imageName {
			return CUSTOM_IMAGE, customImages[i].ImageId, nil
		}
		if strconv.Itoa(customImages[i].ImageId) == imageName {
			return CUSTOM_IMAGE, customImages[i].ImageId, nil
		}
	}

	return -1, -1, fmt.Errorf("Failed to find image %s", imageName)
}

// deployImage deploys the specified image
// DiskLabel has 50 characters maximum!!!
func deployImage(client *golinode.Client, linode golinode.LinodeInstance, imageName string, diskSize int, key, password string) error {
	imageType, imageId, err := findImage(client, imageName)
	if err != nil {
		return err
	}
	args := make(map[string]string)
	args["rootSSHKey"] = key
	args["rootPass"] = password
	diskLabel := fmt.Sprintf("Root(%d)__Base(%d)", linode.LinodeId, imageId)
	if imageType == PREBUILT {
		_, err = client.Disk.CreateFromDistribution(imageId, linode.LinodeId, diskLabel, diskSize, args)
		if err != nil {
			return err
		}
	} else if imageType == CUSTOM_IMAGE {
		_, err = client.Disk.CreateFromImage(imageId, linode.LinodeId, diskLabel, diskSize, args)
		if err != nil {
			return err
		}
	} else {
		panic("Invalid image type returned")
	}
	if 	err = waitForInstanceStatus(client, linode.LinodeId, 'offline', WAIT_TIMEOUT); err != nil {
		return fmt.Errorf("Image %d failed to thaw for linode %d because %s", imageId, linode.LinodeId, err)
	}
	return nil
}

// waitForInstanceStatus waits for the Linode instance to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func waitForInstanceStatus(client *golinode.Client, linodeId int, status string, timeoutSeconds int) error {
	start := time.Now()
	for {
		linode, err := client.GetInstance(linodeId)
		if err != nil {
			return err
		}
		complete := (linode.Status == status)

		if complete {
			return nil
		}

		time.Sleep(1 * time.Second)
		if time.Since(start) > timeoutSeconds * time.Second {
			return fmt.Errorf("Linode %d didn't reach '%s' status in %d seconds", linodeId, status, timeoutSeconds)
		}
	}
}

// changeLinodeSettings changes linode level settings. This is things like the name or the group
func changeLinodeSettings(client *golinode.Client, linode golinode.LinodeInstance, d *schema.ResourceData) error {
	updates := make(map[string]interface{})
	if d.Get("group").(string) != linode.LpmDisplayGroup {
		updates["lpm_displayGroup"] = d.Get("group")
	}

	if d.Get("name").(string) != linode.Label.String() {
		updates["Label"] = d.Get("name")
	}

	if len(updates) > 0 {
		_, err := client.Linode.Update(linode.LinodeId, updates)
		if err != nil {
			return fmt.Errorf("Failed to update the linode's group because %s", err)
		}
	}
	d.SetPartial("group")
	d.SetPartial("name")
	return nil
}

// changeLinodeSize resizes the current linode
func changeLinodeSize(client *golinode.Client, linode golinode.LinodeInstance, d *schema.ResourceData) error {
	var newPlanID int
	var waitMinutes int

	// Get the Linode Plan Size
	sizeId, err := getSizeId(client, d.Get("size").(int))
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
	err = waitForInstanceStatus(client, linode.LinodeId, 'offline')
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
func changeLinodeConfig(client *golinode.Client, config golinode.LinodeConfig, d *schema.ResourceData) error {
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
