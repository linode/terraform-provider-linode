package linode

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLinodeNodeBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceLinodeNodeBalancerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "The unique ID of the Linode NodeBalancer.",
				Required:    true,
			},
			"label": {
				Type:        schema.TypeString,
				Description: "The label of the Linode NodeBalancer.",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The region where this NodeBalancer will be deployed.",
				Computed:    true,
			},
			"client_conn_throttle": {
				Type:        schema.TypeInt,
				Description: "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
				Computed:    true,
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
				Type:        schema.TypeString,
				Description: "When this NodeBalancer was created.",
				Computed:    true,
			},
			"updated": {
				Type:        schema.TypeString,
				Description: "When this NodeBalancer was last updated.",
				Computed:    true,
			},
			"transfer": {
				Type:        schema.TypeList,
				Description: "Information about the amount of transfer this NodeBalancer has had so far this month.",
				Computed:    true,
				Elem:        resourceLinodeNodeBalancerTransfer(),
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
			},
		},
	}
}

func datasourceLinodeNodeBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client
	id := d.Get("id").(int)

	nodebalancer, err := client.GetNodeBalancer(context.Background(), id)
	if err != nil {
		return diag.Errorf("failed to get nodebalancer %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(nodebalancer.ID))
	d.Set("label", nodebalancer.Label)
	d.Set("hostname", nodebalancer.Hostname)
	d.Set("ipv4", nodebalancer.IPv4)
	d.Set("ipv6", nodebalancer.IPv6)
	d.Set("tags", nodebalancer.Tags)
	d.Set("client_conn_throttle", nodebalancer.ClientConnThrottle)
	d.Set("region", nodebalancer.Region)
	d.Set("created", nodebalancer.Created.Format(time.RFC3339))
	d.Set("updated", nodebalancer.Updated.Format(time.RFC3339))
	d.Set("transfer", []map[string]interface{}{{
		"in":    nodebalancer.Transfer.In,
		"out":   nodebalancer.Transfer.Out,
		"total": nodebalancer.Transfer.Total,
	}})

	return nil
}
