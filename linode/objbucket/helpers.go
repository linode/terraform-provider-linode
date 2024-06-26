package objbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"bucket":  d.Get("label"),
		"cluster": d.Get("cluster"),
	})
}

func getRegionOrCluster(d *schema.ResourceData) (regionOrCluster string) {
	if region, ok := d.GetOk("region"); ok {
		regionOrCluster = region.(string)
	} else {
		regionOrCluster = d.Get("cluster").(string)
	}
	return
}
