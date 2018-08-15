package linode

import (
	"context"
	"fmt"

	"github.com/linode/linodego"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLinodeComputeIPv6Pool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeComputeIPv6PoolRead,

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

func dataSourceLinodeComputeIPv6PoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*linodego.Client)

	pools, err := client.ListIPv6Pools(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error listing pools: %s", err)
	}

	reqPool := d.Get("range").(string)

	for _, pool := range pools {
		if pool.Range == reqPool {
			d.SetId(pool.Range)
			d.Set("region", pool.Region)
			return nil
		}
	}

	d.SetId("")

	return fmt.Errorf("IPv6 Pool %s was not found", reqPool)
}
