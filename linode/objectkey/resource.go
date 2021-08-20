package objectkey

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
	}
}

func createResource(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	createOpts := linodego.ObjectStorageKeyCreateOptions{
		Label: d.Get("label").(string),
	}

	if bucketAccess, bucketAccessOk := d.GetOk("bucket_access"); bucketAccessOk {
		createOpts.BucketAccess = expandKeyBucketAccess(bucketAccess.([]interface{}))
	}

	objectStorageKey, err := client.CreateObjectStorageKey(ctx, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Object Storage Key: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", objectStorageKey.ID))
	d.Set("label", objectStorageKey.Label)
	d.Set("access_key", objectStorageKey.AccessKey)

	// secret_key only available on creation
	d.Set("secret_key", objectStorageKey.SecretKey)

	d.Set("limited", objectStorageKey.Limited)

	bucketAccess := flattenKeyBucketAccess(objectStorageKey.BucketAccess)
	if bucketAccess != nil {
		d.Set("bucket_access", bucketAccess)
	}

	return readResource(ctx, d, meta)
}

func readResource(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Object Storage Key ID %s as int: %s", d.Id(), err)
	}

	objectStorageKey, err := client.GetObjectStorageKey(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Object Storage Key %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Object Storage Key: %s", err)
	}

	d.Set("label", objectStorageKey.Label)
	d.Set("access_key", objectStorageKey.AccessKey)
	d.Set("limited", objectStorageKey.Limited)

	bucketAccess := flattenKeyBucketAccess(objectStorageKey.BucketAccess)
	if bucketAccess != nil {
		d.Set("bucket_access", bucketAccess)
	}
	return nil
}

func updateResource(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Object Storage Key id %s as int: %s", d.Id(), err)
	}

	if d.HasChange("label") {
		objectStorageKey, err := client.GetObjectStorageKey(ctx, int(id))

		updateOpts := linodego.ObjectStorageKeyUpdateOptions{
			Label: d.Get("label").(string),
		}

		if err != nil {
			return diag.Errorf("Error fetching data about the current Linode Object Storage Key: %s", err)
		}

		if objectStorageKey, err = client.UpdateObjectStorageKey(ctx, int(id), updateOpts); err != nil {
			return diag.FromErr(err)
		}
		d.Set("label", objectStorageKey.Label)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Object Storage Key id %s as int", d.Id())
	}
	err = client.DeleteObjectStorageKey(ctx, int(id))
	if err != nil {
		return diag.Errorf("Error deleting Linode Object Storage Key %d: %s", id, err)
	}
	return nil
}

func flattenKeyBucketAccess(
	bucketAccesses *[]linodego.ObjectStorageKeyBucketAccess) *[]map[string]interface{} {
	if bucketAccesses == nil {
		return nil
	}
	specs := make([]map[string]interface{}, len(*bucketAccesses))

	for i, bucketAccess := range *bucketAccesses {
		specs[i] = map[string]interface{}{
			"bucket_name": bucketAccess.BucketName,
			"cluster":     bucketAccess.Cluster,
			"permissions": bucketAccess.Permissions,
		}
	}
	return &specs
}

func expandKeyBucketAccess(
	bucketAccessSpecs []interface{}) *[]linodego.ObjectStorageKeyBucketAccess {
	bucketAccesses := make([]linodego.ObjectStorageKeyBucketAccess, len(bucketAccessSpecs))
	for i, bucketAccessSpec := range bucketAccessSpecs {
		bucketAccessSpec := bucketAccessSpec.(map[string]interface{})
		bucketAccess := linodego.ObjectStorageKeyBucketAccess{
			BucketName:  bucketAccessSpec["bucket_name"].(string),
			Cluster:     bucketAccessSpec["cluster"].(string),
			Permissions: bucketAccessSpec["permissions"].(string),
		}
		bucketAccesses[i] = bucketAccess
	}
	return &bucketAccesses
}
