package nb

import (
	"context"
	"strconv"
	"time"

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

	nodebalancer, err := client.GetNodeBalancer(ctx, id)
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
