package objectcluster

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

	reqObjectStorageCluster := d.Get("id").(string)

	if reqObjectStorageCluster == "" {
		return diag.Errorf("Error object storage cluster id is required")
	}

	objectStorageCluster, err := client.GetObjectStorageCluster(ctx, reqObjectStorageCluster)
	if err != nil {
		return diag.Errorf("Error listing object storage clusters: %s", err)
	}

	if objectStorageCluster != nil {
		d.SetId(objectStorageCluster.ID)
		d.Set("domain", objectStorageCluster.Domain)
		d.Set("status", objectStorageCluster.Status)
		d.Set("region", objectStorageCluster.Region)
		d.Set("static_site_domain", objectStorageCluster.StaticSiteDomain)

		return nil
	}

	return diag.Errorf("Linode object storage cluster %s was not found", reqObjectStorageCluster)
}
