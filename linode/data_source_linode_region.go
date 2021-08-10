package linode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func dataSourceLinodeRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLinodeRegionRead,

		Schema: map[string]*schema.Schema{
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The country where this Region resides.",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The unique ID of this Region.",
				Required:    true,
			},
		},
	}
}

func dataSourceLinodeRegionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return nil
	}

	return diag.Errorf("Linode Region %s was not found", reqRegion)
}
