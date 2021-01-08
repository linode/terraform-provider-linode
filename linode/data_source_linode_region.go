package linode

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLinodeRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeRegionRead,

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

func dataSourceLinodeRegionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client

	reqRegion := d.Get("id").(string)

	if reqRegion == "" {
		return fmt.Errorf("Error region id is required")
	}

	region, err := client.GetRegion(context.Background(), reqRegion)
	if err != nil {
		return fmt.Errorf("Error listing regions: %s", err)
	}

	if region != nil {
		d.SetId(region.ID)
		d.Set("country", region.Country)
		return nil
	}

	return fmt.Errorf("Linode Region %s was not found", reqRegion)
}
