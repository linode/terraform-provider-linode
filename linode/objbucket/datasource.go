package objbucket

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      bucketDataSourceSchema,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	//	cluster, label, err := DecodeBucketID(d.Id())
	cluster := d.Get("cluster").(string)
	label := d.Get("label").(string)

	// if err != nil {
	// 	return diag.Errorf("failed to parse Linode ObjectStorageBucket id %s", d.Id())
	// }

	bucket, err := client.GetObjectStorageBucket(ctx, cluster, label)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Object Storage Bucket %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to find the specified Linode ObjectStorageBucket: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", bucket.Cluster, bucket.Label))
	d.Set("cluster", bucket.Cluster)
	d.Set("created", bucket.Created.Format(time.RFC3339))
	d.Set("hostname", bucket.Hostname)
	d.Set("label", bucket.Label)

	return nil
}
