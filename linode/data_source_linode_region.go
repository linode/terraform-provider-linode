package linode

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/linode/linodego"
)

func dataSourceLinodeRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeRegionRead,

		Schema: map[string]*schema.Schema{
			"country": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLinodeRegionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*linodego.Client)

	regions, err := client.ListRegions(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error listing regions: %s", err)
	}

	reqRegion := d.Get("id").(string)

	for _, r := range regions {
		if r.ID == reqRegion {
			d.SetId(r.ID)
			d.Set("country", r.Country)
			return nil
		}
	}

	d.SetId("")

	return fmt.Errorf("Linode Region %s was not found", reqRegion)
}
