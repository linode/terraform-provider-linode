package obj

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"bucket":     d.Get("bucket"),
		"cluster":    d.Get("cluster"),
		"object_key": d.Get("key"),
	})
}

// GetObjKeysFromProvider gets object access_key and secret_key from provider configuration.
// Return whether both of the keys exist.
func GetObjKeysFromProvider(
	d *schema.ResourceData,
	config *helper.Config,
) bool {
	d.Set("access_key", config.ObjAccessKey)
	d.Set("secret_key", config.ObjSecretKey)

	return CheckObjKeysConfiged(d)
}

// CreateTempKeys creates temporary Object Storage Keys to use.
// The temporary keys are scoped only to the target cluster and bucket with limited permissions.
// Keys only exist for the duration of the apply time.
func CreateTempKeys(
	ctx context.Context,
	d *schema.ResourceData,
	client linodego.Client,
	bucket, cluster, permissions string,
) (int, diag.Diagnostics) {
	tflog.Debug(ctx, "Create temporary object storage access keys implicitly.")

	createOpts := linodego.ObjectStorageKeyCreateOptions{
		Label: fmt.Sprintf("temp_%s_%v", bucket, time.Now().Unix()),
		BucketAccess: &[]linodego.ObjectStorageKeyBucketAccess{{
			BucketName:  bucket,
			Cluster:     cluster,
			Permissions: permissions,
		}},
	}

	tflog.Debug(ctx, "client.CreateObjectStorageKey(...)", map[string]interface{}{
		"options": createOpts,
	})

	key, err := client.CreateObjectStorageKey(ctx, createOpts)
	if err != nil {
		return 0, diag.FromErr(err)
	}

	d.Set("access_key", key.AccessKey)
	d.Set("secret_key", key.SecretKey)

	return key.ID, nil
}

// CheckObjKeysConfiged checks whether access_key and secret_key both exist in the schema.
func CheckObjKeysConfiged(d *schema.ResourceData) bool {
	_, accessKeyConfiged := d.GetOk("access_key")
	_, secretKeyConfiged := d.GetOk("secret_key")

	return accessKeyConfiged && secretKeyConfiged
}

// CleanUpTempKeys deleted the temporarily created object keys
func CleanUpTempKeys(
	ctx context.Context,
	client linodego.Client,
	keyId int,
) {
	tflog.Trace(ctx, "Clean up temporary keys: client.DeleteObjectStorageKey(...)", map[string]interface{}{
		"key_id": keyId,
	})

	if err := client.DeleteObjectStorageKey(ctx, keyId); err != nil {
		log.Printf("[WARN] Failed to clean up temporary object storage keys: %s\n", err)
	}
}

func CleanUpKeysFromSchema(d *schema.ResourceData) {
	d.Set("access_key", "")
	d.Set("secret_key", "")
}
