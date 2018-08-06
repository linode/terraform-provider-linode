package linode

import (
	"context"
	"fmt"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLinodeComputeIPv6Range() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeComputeIPv6RangeRead,

		Schema: map[string]*schema.Schema{
			"range": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLinodeComputeIPv6RangeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*linodego.Client)

	ranges, err := client.ListIPv6Ranges(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error listing ranges: %s", err)
	}

	reqRange := d.Get("range").(string)

	for _, r := range ranges {
		if r.Range == reqRange {
			d.SetId(r.Range)
			d.Set("region", r.Region)
			return nil
		}
	}

	d.SetId("")

	return fmt.Errorf("IPv6 Range %s was not found", reqRange)
}
