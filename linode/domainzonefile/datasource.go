package domainzonefile

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      dataSourceSchema,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	domainID := d.Get("domain_id").(int)

	zf, err := client.GetDomainZoneFile(ctx, domainID)
	if err != nil {
		return diag.Errorf("Error fetching domain record: %v", err)
	}

	d.SetId(strconv.Itoa(domainID))
	d.Set("domain_id", domainID)
	d.Set("zone_file", zf.ZoneFile)

	return nil
}
