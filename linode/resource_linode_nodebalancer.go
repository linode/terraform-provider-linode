package linode

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/linodego"
)

func resourceLinodeNodeBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeNodeBalancerCreateContext,
		ReadContext:   resourceLinodeNodeBalancerReadContext,
		UpdateContext: resourceLinodeNodeBalancerUpdateContext,
		DeleteContext: resourceLinodeNodeBalancerDeleteContext,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The label of the Linode NodeBalancer.",
				Optional:    true,
			},
			"region": {
				Type:         schema.TypeString,
				Description:  "The region where this NodeBalancer will be deployed.",
				Required:     true,
				ForceNew:     true,
				InputDefault: "us-east",
			},
			"client_conn_throttle": {
				Type:         schema.TypeInt,
				Description:  "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
				ValidateFunc: validation.IntBetween(0, 20),
				Optional:     true,
				Default:      0,
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "This NodeBalancer's hostname, ending with .nodebalancer.linode.com",
				Computed:    true,
			},
			"ipv4": {
				Type:        schema.TypeString,
				Description: "The Public IPv4 Address of this NodeBalancer",
				Computed:    true,
			},
			"ipv6": {
				Type:        schema.TypeString,
				Description: "The Public IPv6 Address of this NodeBalancer",
				Computed:    true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transfer": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"in": {
							Type:        schema.TypeFloat,
							Description: "The total transfer, in MB, used by this NodeBalancer this month",
							Computed:    true,
						},
						"out": {
							Type:        schema.TypeFloat,
							Description: "The total inbound transfer, in MB, used for this NodeBalancer this month",
							Computed:    true,
						},
						"total": {
							Type:        schema.TypeFloat,
							Description: "The total outbound transfer, in MB, used for this NodeBalancer this month",
							Computed:    true,
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
			},
		},
	}
}

func resourceLinodeNodeBalancerReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %s as int: %s", d.Id(), err)
	}

	nodebalancer, err := client.GetNodeBalancer(context.Background(), int(id))

	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode NodeBalancer ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error finding the specified Linode NodeBalancer: %s", err)
	}

	d.Set("label", nodebalancer.Label)
	d.Set("hostname", nodebalancer.Hostname)
	d.Set("region", nodebalancer.Region)
	d.Set("ipv4", nodebalancer.IPv4)
	d.Set("ipv6", nodebalancer.IPv6)
	d.Set("tags", nodebalancer.Tags)
	d.Set("client_conn_throttle", nodebalancer.ClientConnThrottle)
	d.Set("created", nodebalancer.Created.Format(time.RFC3339))
	d.Set("updated", nodebalancer.Updated.Format(time.RFC3339))
	d.Set("transfer", flattenLinodeNodeBalancerTransfer(nodebalancer.Transfer))

	return nil
}

func resourceLinodeNodeBalancerCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(linodego.Client)
	if !ok {
		return diag.Errorf("Invalid Client when creating Linode NodeBalancer")
	}
	label := d.Get("label").(string)
	clientConnThrottle := d.Get("client_conn_throttle").(int)

	createOpts := linodego.NodeBalancerCreateOptions{
		Region:             d.Get("region").(string),
		Label:              &label,
		ClientConnThrottle: &clientConnThrottle,
	}

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		for _, tag := range tagsRaw.(*schema.Set).List() {
			createOpts.Tags = append(createOpts.Tags, tag.(string))
		}
	}

	nodebalancer, err := client.CreateNodeBalancer(context.Background(), createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode NodeBalancer: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", nodebalancer.ID))

	return resourceLinodeNodeBalancerReadContext(ctx, d, meta)
}

func resourceLinodeNodeBalancerUpdateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode NodeBalancer id %s as int: %s", d.Id(), err)
	}

	nodebalancer, err := client.GetNodeBalancer(context.Background(), int(id))
	if err != nil {
		return diag.Errorf("Error fetching data about the current NodeBalancer: %s", err)
	}

	if d.HasChange("label") || d.HasChange("client_conn_throttle") || d.HasChange("tags") {
		label := d.Get("label").(string)
		clientConnThrottle := d.Get("client_conn_throttle").(int)

		// @TODO nodebalancer.GetUpdateOptions, avoid clobbering client_conn_throttle
		updateOpts := linodego.NodeBalancerUpdateOptions{
			Label:              &label,
			ClientConnThrottle: &clientConnThrottle,
		}

		tags := []string{}
		for _, tag := range d.Get("tags").(*schema.Set).List() {
			tags = append(tags, tag.(string))
		}

		updateOpts.Tags = &tags

		if nodebalancer, err = client.UpdateNodeBalancer(context.Background(), nodebalancer.ID, updateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLinodeNodeBalancerReadContext(ctx, d, meta)
}

func resourceLinodeNodeBalancerDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode NodeBalancer id %s as int", d.Id())
	}
	err = client.DeleteNodeBalancer(context.Background(), int(id))
	if err != nil {
		return diag.Errorf("Error deleting Linode NodeBalancer %d: %s", id, err)
	}
	return nil
}

func flattenLinodeNodeBalancerTransfer(transfer linodego.NodeBalancerTransfer) []map[string]float64 {
	m := map[string]float64{"in": 0, "out": 0, "total": 0}
	if transfer.In != nil {
		m["in"] = *transfer.In
	}
	if transfer.Out != nil {
		m["out"] = *transfer.Out
	}
	if transfer.Total != nil {
		m["total"] = *transfer.Total
	}
	return []map[string]float64{m}
}
