package linode

import (
	"fmt"
	"strconv"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLinodeNodeBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeNodeBalancerCreate,
		Read:   resourceLinodeNodeBalancerRead,
		Update: resourceLinodeNodeBalancerUpdate,
		Delete: resourceLinodeNodeBalancerDelete,
		Exists: resourceLinodeNodeBalancerExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The label of the Linode NodeBalancer.",
				Optional:    true,
			},
			"region": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The region where this NodeBalancer will be deployed.",
				Required:     true,
				ForceNew:     true,
				InputDefault: "us-east",
			},
			"client_conn_throttle": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
				Optional:    true,
				Default:     0,
			},
			"hostname": &schema.Schema{
				Type:        schema.TypeString,
				Description: "This NodeBalancer's hostname, ending with .nodebalancer.linode.com",
				Computed:    true,
			},
			"ipv4": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The Public IPv4 Address of this NodeBalancer",
				Computed:    true,
			},
			"ipv6": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The Public IPv6 Address of this NodeBalancer",
				Computed:    true,
			},
		},
	}
}

func resourceLinodeNodeBalancerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Failed to parse Linode NodeBalancer ID %s as int because %s", d.Id(), err)
	}

	_, err = client.GetNodeBalancer(int(id))
	if err != nil {
		return false, fmt.Errorf("Failed to get Linode NodeBalancer ID %s because %s", d.Id(), err)
	}
	return true, nil
}

func syncResourceData(d *schema.ResourceData, nodebalancer *linodego.NodeBalancer) {
	d.Set("name", nodebalancer.Label)
	d.Set("hostname", nodebalancer.Hostname)
	d.Set("region", nodebalancer.Region)
	d.Set("ipv4", nodebalancer.IPv4)
	d.Set("ipv6", nodebalancer.IPv6)
	d.Set("client_conn_throttle", nodebalancer.ClientConnThrottle)
}

func resourceLinodeNodeBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode NodeBalancer ID %s as int because %s", d.Id(), err)
	}

	nodebalancer, err := client.GetNodeBalancer(int(id))

	if err != nil {
		return fmt.Errorf("Failed to find the specified Linode NodeBalancer because %s", err)
	}

	syncResourceData(d, nodebalancer)

	return nil
}

func resourceLinodeNodeBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode NodeBalancer")
	}
	name := d.Get("name").(string)
	clientConnThrottle := d.Get("client_conn_throttle").(int)

	createOpts := linodego.NodeBalancerCreateOptions{
		Region:             d.Get("region").(string),
		Label:              &name,
		ClientConnThrottle: &clientConnThrottle,
	}
	nodebalancer, err := client.CreateNodeBalancer(&createOpts)
	if err != nil {
		return fmt.Errorf("Failed to create a Linode NodeBalancer in because %s", err)
	}
	d.SetId(fmt.Sprintf("%d", nodebalancer.ID))

	syncResourceData(d, nodebalancer)

	return nil
}

func resourceLinodeNodeBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode NodeBalancer id %s as an int because %s", d.Id(), err)
	}

	nodebalancer, err := client.GetNodeBalancer(int(id))
	if err != nil {
		return fmt.Errorf("Failed to fetch data about the current NodeBalancer because %s", err)
	}

	if d.HasChange("name") || d.HasChange("client_conn_throttle") {
		name := d.Get("name").(string)
		clientConnThrottle := d.Get("client_conn_throttle").(int)
		// @TODO nodebalancer.GetUpdateOptions, avoid clobbering client_conn_throttle
		updateOpts := linodego.NodeBalancerUpdateOptions{
			Label:              &name,
			ClientConnThrottle: &clientConnThrottle,
		}
		if nodebalancer, err = client.UpdateNodeBalancer(nodebalancer.ID, updateOpts); err != nil {
			return err
		}
		syncResourceData(d, nodebalancer)
	}

	return nil
}

func resourceLinodeNodeBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode NodeBalancer id %s as int", d.Id())
	}
	err = client.DeleteNodeBalancer(int(id))
	if err != nil {
		return fmt.Errorf("Failed to delete Linode NodeBalancer %d because %s", id, err)
	}
	return nil
}
