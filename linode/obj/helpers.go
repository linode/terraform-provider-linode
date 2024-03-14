package obj

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type ObjectKeys struct {
	AccessKey string
	SecretKey string
}

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"bucket":     d.Get("bucket"),
		"cluster":    d.Get("cluster"),
		"object_key": d.Get("key"),
	})
}

// GetObjKeysFromProvider gets obj_access_key and obj_secret_key from provider configuration.
// Return whether both of the keys exist.
func GetObjKeysFromProvider(
	keys ObjectKeys,
	config *helper.Config,
) (ObjectKeys, bool) {
	keys.AccessKey = config.ObjAccessKey
	keys.SecretKey = config.ObjSecretKey

	return keys, CheckObjKeysConfiged(keys)
}

// CreateTempKeys creates temporary Object Storage Keys to use.
// The temporary keys are scoped only to the target cluster and bucket with limited permissions.
// Keys only exist for the duration of the apply time.
func CreateTempKeys(
	ctx context.Context,
	client linodego.Client,
	bucket, cluster, permissions string,
) (*linodego.ObjectStorageKey, diag.Diagnostics) {
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

	keys, err := client.CreateObjectStorageKey(ctx, createOpts)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return keys, nil
}

// CheckObjKeysConfiged checks whether AccessKey and SecretKey both exist.
func CheckObjKeysConfiged(keys ObjectKeys) bool {
	return keys.AccessKey != "" && keys.SecretKey != ""
}

// CleanUpTempKeys deleted the temporarily created object keys.
func CleanUpTempKeys(
	ctx context.Context,
	client linodego.Client,
	keyId int,
) {
	tflog.Trace(ctx, "Clean up temporary keys: client.DeleteObjectStorageKey(...)", map[string]interface{}{
		"key_id": keyId,
	})

	if err := client.DeleteObjectStorageKey(ctx, keyId); err != nil {
		tflog.Warn(ctx, "Failed to clean up temporary object storage keys", map[string]interface{}{
			"details": err,
		})
	}
}
