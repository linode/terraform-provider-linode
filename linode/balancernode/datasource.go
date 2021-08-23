package balancernode

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id := d.Get("id").(int)
	nodebalancerID := d.Get("nodebalancer_id").(int)
	configID := d.Get("config_id").(int)

	node, err := client.GetNodeBalancerNode(ctx, nodebalancerID, configID, id)
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
