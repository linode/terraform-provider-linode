package linode

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/linode/linodego"
)

const (
	LinodeLKECreateTimeout = 10 * time.Minute
	LinodeLKEUpdateTimeout = 20 * time.Minute
	LinodeLKEDeleteTimeout = 10 * time.Minute
)

func resourceLinodeLKE() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeLKECreate,
		Read:   resourceLinodeLKERead,
		Update: nil,
		Delete: resourceLinodeLKEDelete,
		Exists: resourceLinodeLKEExists,
		Importer: &schema.ResourceImporter{
			State: nil,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(LinodeLKECreateTimeout),
			Update: schema.DefaultTimeout(LinodeLKEUpdateTimeout),
			Delete: schema.DefaultTimeout(LinodeLKEDeleteTimeout),
		},
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
			},
			"version": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: "The desired Kubernetes version for this Kubernetes cluster in the format of <major>.<minor>, and the latest supported patch version will be deployed.",
			},
			"label": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: "Cluster label.",
			},
			"region": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: "This Kubernetes cluster's location.",
			},
			"node_pools": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeMap},
				Required:    true,
				ForceNew:    true,
				Description: "An array of Linode's instance type for your cluster",
			},
		},
	}
}

func resourceLinodeLKEExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Error parsing Linode LKE ID %s as int: %s", d.Id(), err)
	}

	_, err = client.GetLKECluster(context.Background(), id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			return false, nil
		}

		return false, fmt.Errorf("Error getting Linode LKE ID %s: %s", d.Id(), err)
	}
	return true, nil
}

func resourceLinodeLKERead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Linode LKE ID %s as int: %s", d.Id(), err)
	}

	lke, err := client.GetLKECluster(context.Background(), id)

	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing LKE ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error finding the specified Linode LKE Cluster: %s", err)
	}

	d.Set("label", lke.Label)
	d.Set("region", lke.Region)
	d.Set("version", lke.Version)
	d.Set("tags", lke.Tags)
	return nil
}

func resourceLinodeLKECreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode LKE")
	}
	d.Partial(true)

	createOpts := linodego.LKEClusterCreateOptions{
		Label:   d.Get("label").(string),
		Region:  d.Get("region").(string),
		Version: d.Get("version").(string),
	}

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		for _, tag := range tagsRaw.(*schema.Set).List() {
			createOpts.Tags = append(createOpts.Tags, tag.(string))
		}
	}

	// FIX: Iterate through the pool and create the LKEClusterPoolCreateOptions.
	// for _, nodePoolsRaw := range d.Get("node_pool").([]interface{}) {
	// 	//createOpts.NodePools = append(createOpts.NodePools, v.(linodego.LKEClusterPoolCreateOptions))
	// }

	clusterLKE, err := client.CreateLKECluster(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("error createing LKE cluster %s", err)
	}

	client.WaitForEventFinished(context.Background(), clusterLKE.ID, linodego.EntityLinode, linodego.ActionLKECreate, *clusterLKE.Created, int(d.Timeout(schema.TimeoutCreate).Seconds()))
	if err != nil {
		return fmt.Errorf("Error waiting for Instance to finish creating")
	}
	d.SetId(fmt.Sprintf("%d", clusterLKE.ID))

	return resourceLinodeLKERead(d, meta)
}

func resourceLinodeLKEDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Linode LKE ID %s as int: %s", d.Id(), err)
	}

	minDelete := time.Now().AddDate(0, 0, -1)

	err = client.DeleteLKECluster(context.Background(), id)

	if err != nil {
		return fmt.Errorf("Error deleting Linode LKE cluster %d: %s", id, err)
	}
	client.WaitForEventFinished(context.Background(), id, linodego.EntityLinode, linodego.ActionLKEDelete, minDelete, int(d.Timeout(schema.TimeoutDelete).Seconds()))

	d.SetId("")
	return nil
}
