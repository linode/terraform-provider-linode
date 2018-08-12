package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLinodeNodeBalancerNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeNodeBalancerNodeCreate,
		Read:   resourceLinodeNodeBalancerNodeRead,
		Update: resourceLinodeNodeBalancerNodeUpdate,
		Delete: resourceLinodeNodeBalancerNodeDelete,
		Exists: resourceLinodeNodeBalancerNodeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"nodebalancer_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The ID of the NodeBalancer to access.",
				Required:    true,
				ForceNew:    true,
			},
			"config_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The ID of the NodeBalancerConfig to access.",
				Required:    true,
				ForceNew:    true,
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The label for this node. This is for display purposes only.",
				Optional:    true,
			},
			"weight": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Used when picking a backend to serve a request and is not pinned to a single backend yet. Nodes with a higher weight will receive more traffic. (1-255)",
				Optional:    true,
				Computed:    true,
			},
			"mode": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The mode this NodeBalancer should use when sending traffic to this backend. If set to `accept` this backend is accepting traffic. If set to `reject` this backend will not receive traffic. If set to `drain` this backend will not receive new traffic, but connections already pinned to it will continue to be routed to it.",
				Optional:    true,
				Computed:    true,
			},
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The private IP Address where this backend can be reached. This must be a private IP address.",
				Required:    true,
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The current status of this node, based on the configured checks of its NodeBalancer Config. (unknown, UP, DOWN)",
				Computed:    true,
			},
		},
	}
}

func resourceLinodeNodeBalancerNodeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Error parsing Linode NodeBalancerNode ID %s as int: %s", d.Id(), err)
	}

	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return false, fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	configID, ok := d.Get("config_id").(int)
	if !ok {
		return false, fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("config_id"))
	}

	_, err = client.GetNodeBalancerNode(context.Background(), nodebalancerID, configID, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			d.SetId("")
			return false, nil
		}

		return false, fmt.Errorf("Error getting Linode NodeBalancerNode ID %s: %s", d.Id(), err)
	}
	return true, nil
}

func syncNodeBalancerNodeResourceData(d *schema.ResourceData, node *linodego.NodeBalancerNode) {
	d.Set("label", node.Label)
	d.Set("weight", node.Weight)
	d.Set("mode", node.Mode)
	d.Set("address", node.Address)
	d.Set("status", node.Status)
}

func resourceLinodeNodeBalancerNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode NodeBalancerNode ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	configID, ok := d.Get("config_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("config_id"))
	}

	node, err := client.GetNodeBalancerNode(context.Background(), nodebalancerID, configID, int(id))

	if err != nil {
		return fmt.Errorf("Error finding the specified Linode NodeBalancerNode: %s", err)
	}

	syncNodeBalancerNodeResourceData(d, node)

	return nil
}

func resourceLinodeNodeBalancerNodeCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode NodeBalancerNode")
	}

	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	configID, ok := d.Get("config_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("config_id"))
	}

	createOpts := linodego.NodeBalancerNodeCreateOptions{
		Address: d.Get("address").(string),
		Label:   d.Get("label").(string),
		Mode:    linodego.NodeMode(d.Get("mode").(string)),
		Weight:  d.Get("weight").(int),
	}
	node, err := client.CreateNodeBalancerNode(context.Background(), int(nodebalancerID), int(configID), createOpts)
	if err != nil {
		return fmt.Errorf("Error creating a Linode NodeBalancerNode: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", node.ID))
	d.Set("config_id", configID)
	d.Set("nodebalancer_id", nodebalancerID)

	syncNodeBalancerNodeResourceData(d, node)

	return nil
}

func resourceLinodeNodeBalancerNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode NodeBalancerConfig ID %v as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	configID, ok := d.Get("config_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("config_id"))
	}

	node, err := client.GetNodeBalancerNode(context.Background(), nodebalancerID, configID, int(id))
	if err != nil {
		return fmt.Errorf("Error fetching data about the current NodeBalancerNode: %s", err)
	}

	updateOpts := linodego.NodeBalancerNodeUpdateOptions{
		Address: d.Get("address").(string),
		Label:   d.Get("label").(string),
		Mode:    linodego.NodeMode(d.Get("mode").(string)),
		Weight:  d.Get("weight").(int),
	}

	if node, err = client.UpdateNodeBalancerNode(context.Background(), nodebalancerID, configID, int(id), updateOpts); err != nil {
		return err
	}
	syncNodeBalancerNodeResourceData(d, node)

	return nil
}

func resourceLinodeNodeBalancerNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode NodeBalancerConfig ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	configID, ok := d.Get("config_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("config_id"))
	}
	err = client.DeleteNodeBalancerNode(context.Background(), nodebalancerID, configID, int(id))
	if err != nil {
		return fmt.Errorf("Error deleting Linode NodeBalancerNode %d: %s", id, err)
	}
	return nil
}
