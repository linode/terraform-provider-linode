package instanceip

import (
	"context"
	"fmt"

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

	linodeID := d.Get("linode_id").(int)

	netInfo, err := client.GetInstanceIPAddresses(ctx, linodeID)
	if err != nil {
		return diag.Errorf("Error listing network info: %s", err)
	}

	d.Set("ipv4", flattenIPv4(netInfo.IPv4))
	d.Set("ipv6", flattenIPv6(netInfo.IPv6))

	d.SetId(fmt.Sprintf("%d", linodeID))
	return nil
}
