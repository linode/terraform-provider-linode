package linode

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"transfer": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"in": &schema.Schema{
							Type:        schema.TypeFloat,
							Description: "The total transfer, in MB, used by this NodeBalancer this month",
							Optional:    true,
							Computed:    true,
						},
						"out": &schema.Schema{
							Type:        schema.TypeFloat,
							Description: "The total inbound transfer, in MB, used for this NodeBalancer this month",
							Optional:    true,
							Computed:    true,
						},
						"total": &schema.Schema{
							Type:        schema.TypeFloat,
							Description: "The total outbound transfer, in MB, used for this NodeBalancer this month",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceLinodeNodeBalancerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Error parsing Linode NodeBalancer ID %s as int: %s", d.Id(), err)
	}

	_, err = client.GetNodeBalancer(context.Background(), int(id))
	if err != nil {
		if _, ok := err.(*linodego.Error); ok {
			d.SetId("")
			return false, nil
		}

		return false, fmt.Errorf("Error getting sdd marty Linode NodeBalancer ID %s: %s", d.Id(), err)
	}
	return true, nil
}

func syncNodeBalancerData(d *schema.ResourceData, nodebalancer *linodego.NodeBalancer) {
	d.Set("label", nodebalancer.Label)
	d.Set("hostname", nodebalancer.Hostname)
	d.Set("region", nodebalancer.Region)
	d.Set("ipv4", nodebalancer.IPv4)
	d.Set("ipv6", nodebalancer.IPv6)
	d.Set("client_conn_throttle", nodebalancer.ClientConnThrottle)
	d.Set("created", nodebalancer.Created.Format(time.RFC3339))
	d.Set("updated", nodebalancer.Updated.Format(time.RFC3339))
	transfer := map[string]interface{}{
		"in":    floatString(nodebalancer.Transfer.In),
		"out":   floatString(nodebalancer.Transfer.Out),
		"total": floatString(nodebalancer.Transfer.Total),
	}
	if err := d.Set("transfer", transfer); err != nil {
		panic(err)
	}
}

// floatString returns nil or the string representation of the supplied *float64
//   this is needed because ResourceData.Set will not accept *float64 and expects a string
func floatString(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%g", *f)
}

func resourceLinodeNodeBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %s as int: %s", d.Id(), err)
	}

	nodebalancer, err := client.GetNodeBalancer(context.Background(), int(id))

	if err != nil {
		return fmt.Errorf("Error finding the specified Linode NodeBalancer: %s", err)
	}

	syncNodeBalancerData(d, nodebalancer)

	return nil
}

func resourceLinodeNodeBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode NodeBalancer")
	}
	label := d.Get("label").(string)
	clientConnThrottle := d.Get("client_conn_throttle").(int)

	createOpts := linodego.NodeBalancerCreateOptions{
		Region:             d.Get("region").(string),
		Label:              &label,
		ClientConnThrottle: &clientConnThrottle,
	}
	nodebalancer, err := client.CreateNodeBalancer(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("Error creating a Linode NodeBalancer: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", nodebalancer.ID))

	syncNodeBalancerData(d, nodebalancer)

	return nil
}

func resourceLinodeNodeBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode NodeBalancer id %s as int: %s", d.Id(), err)
	}

	nodebalancer, err := client.GetNodeBalancer(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Error fetching data about the current NodeBalancer: %s", err)
	}

	if d.HasChange("label") || d.HasChange("client_conn_throttle") {
		label := d.Get("label").(string)
		clientConnThrottle := d.Get("client_conn_throttle").(int)
		// @TODO nodebalancer.GetUpdateOptions, avoid clobbering client_conn_throttle
		updateOpts := linodego.NodeBalancerUpdateOptions{
			Label:              &label,
			ClientConnThrottle: &clientConnThrottle,
		}
		if nodebalancer, err = client.UpdateNodeBalancer(context.Background(), nodebalancer.ID, updateOpts); err != nil {
			return err
		}
		syncNodeBalancerData(d, nodebalancer)
	}

	return nil
}

func resourceLinodeNodeBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode NodeBalancer id %s as int", d.Id())
	}
	err = client.DeleteNodeBalancer(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Error deleting Linode NodeBalancer %d: %s", id, err)
	}
	return nil
}
