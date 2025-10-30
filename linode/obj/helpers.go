package obj

import (
	"context"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"io"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkv2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type ObjectKeys struct {
	AccessKey string
	SecretKey string
}

func getS3ClientFromModel(
	ctx context.Context,
	client *linodego.Client,
	config *helper.FrameworkProviderModel,
	data ResourceModel,
	permission string,
	endpointType *linodego.ObjectStorageEndpointType,
	diags *diag.Diagnostics,
) (*s3.Client, func()) {
	keys, teardownKeys := data.GetObjectStorageKeys(ctx, client, config, permission, endpointType, diags)
	if diags.HasError() {
		return nil, teardownKeys
	}

	s3client := helper.FwS3Connection(
		ctx,
		data.Endpoint.ValueString(),
		keys.AccessKey,
		keys.SecretKey,
		diags,
	)
	if diags.HasError() {
		return nil, teardownKeys
	}

	return s3client, teardownKeys
}

// getObjKeysFromProvider gets obj_access_key and obj_secret_key from provider configuration.
// Return whether both of the keys exist.
func getObjKeysFromProvider(
	keys ObjectKeys,
	config *helper.Config,
) (ObjectKeys, bool) {
	keys.AccessKey = config.ObjAccessKey
	keys.SecretKey = config.ObjSecretKey

	return keys, keys.Ok()
}

func isCluster(regionOrCluster string) bool {
	pattern := `^[a-z]{2}-[a-z]+-[0-9]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(regionOrCluster)
}

// fwCreateTempKeys creates temporary Object Storage Keys to use.
// The temporary keys are scoped only to the target cluster and bucket with limited permissions.
// Keys only exist for the duration of the apply time.
func fwCreateTempKeys(
	ctx context.Context,
	client *linodego.Client,
	bucketLabel, regionOrCluster, permissions string,
	endpointType *linodego.ObjectStorageEndpointType,
	diags *diag.Diagnostics,
) *linodego.ObjectStorageKey {
	tflog.Debug(ctx, "Create temporary object storage access keys implicitly.")

	tempBucketAccess := linodego.ObjectStorageKeyBucketAccess{
		BucketName:  bucketLabel,
		Permissions: permissions,
	}

	if isCluster(regionOrCluster) {
		tflog.Warn(ctx, "Cluster is deprecated for Linode Object Storage service, please consider switch to using region.")
		tempBucketAccess.Cluster = regionOrCluster
	} else {
		tflog.Info(ctx, fmt.Sprintf("%q Is Region", regionOrCluster))
		tempBucketAccess.Region = regionOrCluster
	}

	createOpts := linodego.ObjectStorageKeyCreateOptions{
		Label:        fmt.Sprintf("temp_%s_%v", bucketLabel, time.Now().Unix()),
		BucketAccess: &[]linodego.ObjectStorageKeyBucketAccess{tempBucketAccess},
	}

	tflog.Debug(ctx, "client.CreateObjectStorageKey(...)", map[string]interface{}{
		"options": createOpts,
	})

	keys, err := client.CreateObjectStorageKey(ctx, createOpts)
	if err != nil {
		diags.AddError("Failed to Create Object Storage Key", err.Error())
		return nil
	}

	if endpointType == nil {
		et, err := getBucketEndpointType(ctx, client, regionOrCluster, bucketLabel)
		if err != nil {
			diags.AddWarning(
				"Can't determine the type of the object storage endpoint. If the it's an E2/E3 OBJ clusters, "+
					"it may lead to an issue that temporary limited key is used before becoming effective",
				err.Error(),
			)
		} else {
			endpointType = &et
		}
	}

	// OBJ limited key for OBJ gen2 takes at most 30s to refresh the cache can becomes effective
	if endpointType != nil && *endpointType != linodego.ObjectStorageEndpointE0 && *endpointType != linodego.ObjectStorageEndpointE1 {
		time.Sleep(30 * time.Second)
	}

	return keys
}

func getBucketEndpointType(
	ctx context.Context, client *linodego.Client, cluster, label string,
) (linodego.ObjectStorageEndpointType, error) {
	bucket, err := client.GetObjectStorageBucket(ctx, cluster, label)
	if err != nil {
		return "", err
	}

	return bucket.EndpointType, nil
}

// createTempKeys creates temporary Object Storage Keys to use.
// The temporary keys are scoped only to the target cluster and bucket with limited permissions.
// Keys only exist for the duration of the apply time.
func createTempKeys(
	ctx context.Context,
	client *linodego.Client,
	bucketLabel, regionOrCluster, permissions string,
	endpointType *linodego.ObjectStorageEndpointType,
) (*linodego.ObjectStorageKey, sdkv2diag.Diagnostics) {
	tflog.Debug(ctx, "Create temporary object storage access keys implicitly.")

	tempBucketAccess := linodego.ObjectStorageKeyBucketAccess{
		BucketName:  bucketLabel,
		Permissions: permissions,
	}

	if isCluster(regionOrCluster) {
		tflog.Warn(ctx, "Cluster is deprecated for Linode Object Storage service, please consider switch to using region.")
		tempBucketAccess.Cluster = regionOrCluster
	} else {
		tempBucketAccess.Region = regionOrCluster
	}

	// Bucket key labels are a maximum of 50 characters - if the bucket name is
	// too long, then truncate it.
	// We use 16 characters for `temp__{timestamp}`, so the maximum length of a
	// full bucket name is 34.
	if len(bucketLabel) > 34 {
		bucketLabel = bucketLabel[:34]
	}
	createOpts := linodego.ObjectStorageKeyCreateOptions{
		Label:        fmt.Sprintf("temp_%s_%v", bucketLabel, time.Now().Unix()),
		BucketAccess: &[]linodego.ObjectStorageKeyBucketAccess{tempBucketAccess},
	}

	tflog.Debug(ctx, "client.CreateObjectStorageKey(...)", map[string]interface{}{
		"options": createOpts,
	})

	keys, err := client.CreateObjectStorageKey(ctx, createOpts)
	if err != nil {
		return nil, sdkv2diag.FromErr(err)
	}
	if endpointType == nil {
		et, err := getBucketEndpointType(ctx, client, regionOrCluster, bucketLabel)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Can't determine the type of the object storage endpoint: %s", err.Error()))
		} else {
			endpointType = &et
		}
	}

	// OBJ limited key for OBJ gen2 takes at most 30s to refresh the cache can becomes effective
	if endpointType != nil && *endpointType != linodego.ObjectStorageEndpointE0 && *endpointType != linodego.ObjectStorageEndpointE1 {
		time.Sleep(30 * time.Second)
	}
	// panic(fmt.Sprintf("%v", endpointType))

	return keys, nil
}

// checkObjKeysConfigured checks whether AccessKey and SecretKey both exist.
func (keys ObjectKeys) Ok() bool {
	return keys.AccessKey != "" && keys.SecretKey != ""
}

// cleanUpTempKeys deleted the temporarily created object keys.
func cleanUpTempKeys(
	ctx context.Context,
	client *linodego.Client,
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
	bucket, regionOrCluster, permission string,
	endpointType *linodego.ObjectStorageEndpointType,
) (ObjectKeys, sdkv2diag.Diagnostics, func()) {
	var teardownTempKeysCleanUp func() = nil

	objKeys := ObjectKeys{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	if !objKeys.Ok() {
		// If object keys don't exist in the resource configuration, firstly look for the keys from provider configuration
		if providerKeys, ok := getObjKeysFromProvider(objKeys, config); ok {
			objKeys = providerKeys
		} else if config.ObjUseTempKeys {
			// Implicitly create temporary object storage keys
			keys, diag := createTempKeys(ctx, &client, bucket, regionOrCluster, permission, endpointType)
			if diag != nil {
				return objKeys, diag, nil
			}
			objKeys.AccessKey = keys.AccessKey
			objKeys.SecretKey = keys.SecretKey
			teardownTempKeysCleanUp = func() { cleanUpTempKeys(ctx, &client, keys.ID) }
		} else {
			return objKeys, sdkv2diag.Errorf(
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
	diags *diag.Diagnostics,
) {
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
						"Failed to put Bucket (%v) Object (%v) with input %v: %s. Retrying...",
						aws.ToString(putInput.Bucket),
						aws.ToString(putInput.Key),
						putInput,
						err.Error(),
					),
				)
				continue
			}

			return

		case <-ctx.Done():
			// The timeout for this context will implicitly be handled by Terraform
			diags.AddError("Failed to Put the Object", ctx.Err().Error())
			return
		}
	}
}

func getQuotesTrimmedETag(
	obj s3.HeadObjectOutput,
) *string {
	if obj.ETag != nil {
		result := strings.Trim(*obj.ETag, `"`)
		return &result
	}
	return nil
}

func deleteObject(ctx context.Context, client *s3.Client, bucket, key, version string, force bool) error {
	tflog.Debug(ctx, "deleting the object key")
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket:                    &bucket,
		Key:                       &key,
		BypassGovernanceRetention: aws.Bool(force),
	}
	if version != "" {
		deleteObjectInput.VersionId = &version
	}

	tflog.Debug(ctx, "client.DeleteObject(...)", map[string]any{"options": deleteObjectInput})
	_, err := client.DeleteObject(ctx, deleteObjectInput)
	if err != nil {
		msg := fmt.Sprintf("failed to delete object version (%s): %s", version, err)
		tflog.Error(ctx, msg)
		if !helper.IsObjNotFoundErr(err) {
			return fmt.Errorf("%s: %w", msg, err)
		}
	}
	return nil
}

func flattenObjectMetadata(metadata map[string]string) map[string]attr.Value {
	metadataObject := make(map[string]attr.Value, len(metadata))
	for key, value := range metadata {
		key := strings.ToLower(key)
		metadataObject[key] = types.StringValue(value)
	}

	return metadataObject
}

func deleteBucketNotFound(diags diag.Diagnostics) diag.Diagnostics {
	return slices.DeleteFunc(diags, func(d diag.Diagnostic) bool {
		return strings.Contains(d.Detail(), "Bucket not found")
	})
}

func AddObjectResource(
	ctx context.Context,
	resp *resource.CreateResponse,
	plan ResourceModel,
) {
	plan.GenerateObjectStorageObjectID(true, true)
	resp.State.SetAttribute(ctx, path.Root("bucket"), plan.Bucket)
	resp.State.SetAttribute(ctx, path.Root("key"), plan.Key)
	resp.State.SetAttribute(ctx, path.Root("cluster"), plan.Cluster)
	resp.State.SetAttribute(ctx, path.Root("region"), plan.Region)
}

func fwPutObject(
	ctx context.Context,
	data ResourceModel,
	s3client *s3.Client,
	diags *diag.Diagnostics,
) {
	tflog.Debug(ctx, "getting object body from resource data")

	body := data.getObjectBody(diags)
	if diags.HasError() {
		return
	}
	defer body.Close()

	sumHandler := crc32.NewIEEE()
	if _, err := io.Copy(sumHandler, body); err != nil {
		diags.AddError(
			"Failed to calculate object body CRC32 sum",
			err.Error(),
		)
		return
	}

	encodedSum := base64.StdEncoding.EncodeToString(sumHandler.Sum(nil))

	// Seek the beginning of the body for uploading
	if _, err := body.Seek(0, 0); err != nil {
		diags.AddError(
			"Failed to seek beginning of object body",
			err.Error(),
		)
		return
	}

	putInput := &s3.PutObjectInput{
		Body:   body,
		Bucket: data.Bucket.ValueStringPointer(),
		Key:    data.Key.ValueStringPointer(),

		ChecksumCRC32: &encodedSum,

		ACL:                     s3types.ObjectCannedACL(data.ACL.ValueString()),
		CacheControl:            data.CacheControl.ValueStringPointer(),
		ContentDisposition:      data.ContentDisposition.ValueStringPointer(),
		ContentEncoding:         data.ContentEncoding.ValueStringPointer(),
		ContentLanguage:         data.ContentLanguage.ValueStringPointer(),
		ContentType:             data.ContentType.ValueStringPointer(),
		WebsiteRedirectLocation: data.WebsiteRedirect.ValueStringPointer(),
	}

	if len(data.Metadata.Elements()) > 0 {
		data.Metadata.ElementsAs(ctx, &putInput.Metadata, false)
		tflog.Debug(ctx, fmt.Sprintf("got Metadata: %v", putInput.Metadata))
	}

	putObjectWithRetries(ctx, s3client, putInput, time.Second*5, diags)
	if diags.HasError() {
		return
	}
}
