package lkeversion

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

	version := d.Get("id").(string)

	versionInfo, err := client.GetLKEVersion(ctx, version)
	if err != nil {
		return diag.Errorf("Error getting lke version %s: %s", version, err)
	}

	d.Set("id", versionInfo.ID)
	d.SetId(versionInfo.ID)
	return nil
}
