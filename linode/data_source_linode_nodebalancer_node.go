package linode

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLinodeNodeBalancerNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceLinodeNodeBalancerNodeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "The ID of the NodeBalancer node.",
				Required:    true,
			},
			"nodebalancer_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the NodeBalancer to access.",
				Required:    true,
			},
			"config_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the NodeBalancerConfig to access.",
				Required:    true,
			},
			"label": {
				Type:        schema.TypeString,
				Description: "The label for this node. This is for display purposes only.",
				Computed:    true,
			},
			"weight": {
				Type: schema.TypeInt,
				Description: "Used when picking a backend to serve a request and is not pinned to a single backend yet. " +
					"Nodes with a higher weight will receive more traffic. (1-255)",
				Computed: true,
			},
			"mode": {
				Type: schema.TypeString,
				Description: "The mode this NodeBalancer should use when sending traffic to this backend. " +
					"If set to `accept` this backend is accepting traffic. If set to `reject` this backend will not " +
					"receive traffic. If set to `drain` this backend will not receive new traffic, but connections " +
					"already pinned to it will continue to be routed to it. If set to `backup` this backend will only " +
					"accept traffic if all other nodes are down.",
				Computed: true,
			},
			"address": {
				Type: schema.TypeString,
				Description: "The private IP Address and port (IP:PORT) where this backend can be reached. " +
					"This must be a private IP address.",
				Computed: true,
			},
			"status": {
				Type: schema.TypeString,
				Description: "The current status of this node, based on the configured checks of its NodeBalancer Config. " +
					"(unknown, UP, DOWN)",
				Computed: true,
			},
		},
	}
}

func datasourceLinodeNodeBalancerNodeRead(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	id := d.Get("id").(int)
	nodebalancerID := d.Get("nodebalancer_id").(int)
	configID := d.Get("config_id").(int)

	node, err := client.GetNodeBalancerNode(context.Background(), nodebalancerID, configID, id)
	if err != nil {
		return diag.Errorf("failed to get nodebalancer node %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(node.ID))
	d.Set("nodebalancer_id", nodebalancerID)
	d.Set("config_id", configID)
	d.Set("label", node.Label)
	d.Set("weight", node.Weight)
	d.Set("mode", node.Mode)
	d.Set("address", node.Address)
	d.Set("status", node.Status)

	return nil
}
