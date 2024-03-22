package obj

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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

// getObjKeysFromProvider gets obj_access_key and obj_secret_key from provider configuration.
// Return whether both of the keys exist.
func getObjKeysFromProvider(
	keys ObjectKeys,
	config *helper.Config,
) (ObjectKeys, bool) {
	keys.AccessKey = config.ObjAccessKey
	keys.SecretKey = config.ObjSecretKey

	return keys, checkObjKeysConfigured(keys)
}

// createTempKeys creates temporary Object Storage Keys to use.
// The temporary keys are scoped only to the target cluster and bucket with limited permissions.
// Keys only exist for the duration of the apply time.
func createTempKeys(
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

// checkObjKeysConfigured checks whether AccessKey and SecretKey both exist.
func checkObjKeysConfigured(keys ObjectKeys) bool {
	return keys.AccessKey != "" && keys.SecretKey != ""
}

// cleanUpTempKeys deleted the temporarily created object keys.
func cleanUpTempKeys(
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

// GetObjKeys gets object access_key and secret_key in the following order:
// 1) Whether the keys are specified in the resource configuration;
// 2) Whether the provider-level object keys exist;
// 3) Whether user opts-in temporary keys generation.
func GetObjKeys(
	ctx context.Context,
	d *schema.ResourceData,
	config *helper.Config,
	client linodego.Client,
	bucket, cluster, permission string,
) (ObjectKeys, diag.Diagnostics, func()) {
	var teardownTempKeysCleanUp func() = nil

	objKeys := ObjectKeys{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	if !checkObjKeysConfigured(objKeys) {
		// If object keys don't exist in the resource configuration, firstly look for the keys from provider configuration
		if providerKeys, ok := getObjKeysFromProvider(objKeys, config); ok {
			objKeys = providerKeys
		} else if config.ObjUseTempKeys {
			// Implicitly create temporary object storage keys
			keys, diag := createTempKeys(ctx, client, bucket, cluster, permission)
			if diag != nil {
				return objKeys, diag, nil
			}
			objKeys.AccessKey = keys.AccessKey
			objKeys.SecretKey = keys.SecretKey
			teardownTempKeysCleanUp = func() { cleanUpTempKeys(ctx, client, keys.ID) }
		} else {
			return objKeys, diag.Errorf(
				"access_key and secret_key are required.",
			), nil
		}
	}

	return objKeys, nil, teardownTempKeysCleanUp
}

func putObjectWithRetries(
	ctx context.Context,
	s3client *s3.Client,
	putInput *s3.PutObjectInput,
	retryDuration time.Duration,
) error {
	tflog.Debug(ctx, "Attempting to put object with retries")

	ticker := time.NewTicker(retryDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tflog.Debug(ctx, "putting the object", map[string]any{
				"PutObjectInput": putInput,
			})

			if _, err := s3client.PutObject(ctx, putInput); err != nil {
				tflog.Debug(ctx,
					fmt.Sprintf(
						"Failed to put Bucket (%v) Object (%v): %s. Retrying...",
						putInput.Bucket,
						putInput.Key,
						err.Error(),
					),
				)
				continue
			}

			return nil

		case <-ctx.Done():
			// The timeout for this context will implicitly be handled by Terraform
			return fmt.Errorf("failed to put the object: %s", ctx.Err())
		}
	}
}
