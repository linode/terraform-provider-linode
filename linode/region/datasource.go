package region

import (
	"context"

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

	reqRegion := d.Get("id").(string)

	if reqRegion == "" {
		return diag.Errorf("Error region id is required")
	}

	region, err := client.GetRegion(ctx, reqRegion)
	if err != nil {
		return diag.Errorf("Error listing regions: %s", err)
	}

	if region != nil {
		d.SetId(region.ID)
		d.Set("country", region.Country)
		d.Set("label", region.Label)
		d.Set("capabilities", region.Capabilities)
		d.Set("status", region.Status)
		d.Set("resolvers", []map[string]interface{}{{
			"ipv4": region.Resolvers.IPv4,
			"ipv6": region.Resolvers.IPv6,
		}})
		return nil
	}

	return diag.Errorf("Linode Region %s was not found", reqRegion)
}

// func flattenRegionResolvers(data interface{}) map[string]interface{} {
// 	t := data.(linodego.RegionResolvers)

// 	result := make(map[string]interface{})

// 	result["ipv4"] = t.IPv4
// 	result["ipv6"] = t.IPv4

// 	return result
// }
