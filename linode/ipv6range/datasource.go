package ipv6range

import (
	"context"
	"strings"

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

	rangeStrSplit := strings.Split(d.Get("range").(string), "/")
	rangeStr := rangeStrSplit[0]

	rangeData, err := client.GetIPv6Range(ctx, rangeStr)
	if err != nil {
		return diag.Errorf("failed to get ipv6 range %s: %s", rangeStr, err)
	}

	d.Set("is_bgp", rangeData.IsBGP)
	d.Set("linodes", rangeData.Linodes)
	d.Set("prefix", rangeData.Prefix)
	d.Set("region", rangeData.Region)

	d.SetId(rangeData.Range)

	return nil
}
