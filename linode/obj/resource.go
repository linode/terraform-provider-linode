package obj

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema: resourceSchema,

		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,

		CustomizeDiff: diffResource,
	}
}

func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "creating linode_object_storage_object")
	return putObject(ctx, d, meta)
}

func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "reading linode_object_storage_object")

	config := meta.(*helper.ProviderMeta).Config
	client := meta.(*helper.ProviderMeta).Client
	cluster := d.Get("cluster").(string)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	if !CheckObjKeysConfiged(d) {
		// If object keys don't exist in the plan, firstly look for the keys from provider configuration
		if !GetObjKeysFromProvider(d, config) && config.ObjUseTempKeys {
			// Implicitly create temporary object storage keys
			tempKeyId, diag := CreateTempKeys(ctx, d, client, bucket, cluster, "read_only")
			if diag != nil {
				return diag
			}

			defer CleanUpTempKeys(ctx, client, tempKeyId)
		}
		defer CleanUpKeysFromSchema(ctx, d)
	}

	if !CheckObjKeysConfiged(d) {
		return diag.Errorf("access_key and secret_key are required to read linode_object_storage_object")
	}

	s3client, err := helper.S3ConnectionFromData(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	headObjectInput := &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	tflog.Debug(ctx, "getting object header", map[string]any{"HeadObjectInput": headObjectInput})
	headOutput, err := s3client.HeadObject(
		ctx,
		headObjectInput,
	)
	if err != nil {
		if helper.IsObjNotFoundErr(err) {
			d.SetId("")
			tflog.Warn(ctx,
				"couldn't find the bucket or object, "+
					"removing the object from the TF state")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("cache_control", headOutput.CacheControl)
	d.Set("content_disposition", headOutput.ContentDisposition)
	d.Set("content_encoding", headOutput.ContentEncoding)
	d.Set("content_language", headOutput.ContentLanguage)
	d.Set("content_type", headOutput.ContentType)
	d.Set("etag", strings.Trim(helper.StringValue(headOutput.ETag), `"`))
	d.Set("website_redirect", headOutput.WebsiteRedirectLocation)
	d.Set("version_id", headOutput.VersionId)
	d.Set("metadata", flattenObjectMetadata(headOutput.Metadata))

	// Compute s3 endpoint when it's not configured by the user
	if _, ok := d.GetOk("endpoint"); !ok {
		tflog.Debug(ctx, "'endpoint' wasn't configured, computing it from cluster name")
		endpoint, err := helper.ComputeS3Endpoint(ctx, d, meta)
		if err != nil {
			return diag.Errorf("failed to compute object storage endpoint: %s", err)
		}
		tflog.Debug(ctx, fmt.Sprintf("computed endpoint: '%s'", endpoint))
		d.Set("endpoint", endpoint)
	}

	return nil
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "updating linode_object_storage_object")
	if d.HasChanges("cache_control", "content_base64", "content_disposition",
		"content_encoding", "content_language", "content_type", "content",
		"etag", "metadata", "source", "website_redirect") {
		tflog.Debug(ctx, "detected qualified change(s), calling 'putObject'")
		return putObject(ctx, d, meta)
	}

	cluster := d.Get("cluster").(string)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	acl := s3types.ObjectCannedACL(d.Get("acl").(string))

	if d.HasChange("acl") {
		config := meta.(*helper.ProviderMeta).Config
		client := meta.(*helper.ProviderMeta).Client

		if !CheckObjKeysConfiged(d) {
			// If object keys don't exist in the plan, firstly look for the keys from provider configuration
			if !GetObjKeysFromProvider(d, config) && config.ObjUseTempKeys {
				// Implicitly create temporary object storage keys
				tempKeyId, diag := CreateTempKeys(ctx, d, client, bucket, cluster, "read_write")
				if diag != nil {
					return diag
				}

				defer CleanUpTempKeys(ctx, client, tempKeyId)
			}
			defer CleanUpKeysFromSchema(ctx, d)
		}

		if !CheckObjKeysConfiged(d) {
			return diag.Errorf("access_key and secret_key are required to update linode_object_storage_object")
		}

		s3client, err := helper.S3ConnectionFromData(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		aclPutInput := &s3.PutObjectAclInput{
			Bucket: &bucket,
			Key:    &key,
			ACL:    acl,
		}
		tflog.Debug(
			ctx,
			"detected ACL change in TF files, updating it on the cloud",
			map[string]any{"PutObjectAclInput": aclPutInput},
		)

		_, err = s3client.PutObjectAcl(ctx, aclPutInput)
		if err != nil {
			return diag.Errorf("failed to put Bucket (%s) Object (%s) ACL: %s", bucket, key, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "deleting linode_object_storage_object")

	config := meta.(*helper.ProviderMeta).Config
	client := meta.(*helper.ProviderMeta).Client
	cluster := d.Get("cluster").(string)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	force := d.Get("force_destroy").(bool)

	if !CheckObjKeysConfiged(d) {
		// If object keys don't exist in the plan, firstly look for the keys from provider configuration
		if !GetObjKeysFromProvider(d, config) && config.ObjUseTempKeys {
			// Implicitly create temporary object storage keys
			tempKeyId, diag := CreateTempKeys(ctx, d, client, bucket, cluster, "read_write")
			if diag != nil {
				return diag
			}

			defer CleanUpTempKeys(ctx, client, tempKeyId)
		}
		defer CleanUpKeysFromSchema(ctx, d)
	}

	if !CheckObjKeysConfiged(d) {
		return diag.Errorf("access_key and secret_key are required to delete linode_object_storage_object")
	}

	s3client, err := helper.S3ConnectionFromData(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk("version_id"); ok {
		tflog.Debug(ctx, "versioning was enabled for this object, deleting all versions and delete markers")
		return diag.FromErr(
			helper.DeleteAllObjectVersionsAndDeleteMarkers(ctx, s3client, bucket, key, force, true),
		)
	}
	tflog.Debug(ctx, "versioning was disabled for this object, simply delete the object")
	return diag.FromErr(deleteObject(ctx, s3client, bucket, strings.TrimPrefix(key, "/"), "", force))
}

func diffResource(
	ctx context.Context, d *schema.ResourceDiff, meta any,
) error {
	if d.HasChange("etag") {
		tflog.Debug(ctx, "'etag' has been changed, computing new 'version_id'")
		d.SetNewComputed("version_id")
	}
	return nil
}

// putObject builds the object from spec and puts it in the
// specified bucket via the *schema.ResourceData, then it calls
// readResource.
func putObject(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tflog.Debug(ctx, "entered 'putObject' function")

	config := meta.(*helper.ProviderMeta).Config
	client := meta.(*helper.ProviderMeta).Client
	cluster := d.Get("cluster").(string)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	if !CheckObjKeysConfiged(d) {
		// If object keys don't exist in the plan, firstly look for the keys from provider configuration
		if !GetObjKeysFromProvider(d, config) && config.ObjUseTempKeys {
			// Implicitly create temporary object storage keys
			tempKeyId, diag := CreateTempKeys(ctx, d, client, bucket, cluster, "read_write")
			if diag != nil {
				return diag
			}

			defer CleanUpTempKeys(ctx, client, tempKeyId)
		}
		defer CleanUpKeysFromSchema(ctx, d)
	}

	if !CheckObjKeysConfiged(d) {
		return diag.Errorf("access_key and secret_key are required to create linode_object_storage_object")
	}

	s3client, err := helper.S3ConnectionFromData(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "getting object body from resource data")
	body, err := objectBodyFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	defer body.Close()

	nilOrValue := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	putInput := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   &body,

		CacheControl:            nilOrValue(d.Get("cache_control").(string)),
		ContentDisposition:      nilOrValue(d.Get("content_disposition").(string)),
		ContentEncoding:         nilOrValue(d.Get("content_encoding").(string)),
		ContentLanguage:         nilOrValue(d.Get("content_language").(string)),
		ContentType:             nilOrValue(d.Get("content_type").(string)),
		WebsiteRedirectLocation: nilOrValue(d.Get("website_redirect").(string)),
	}

	if acl := nilOrValue(d.Get("acl").(string)); acl != nil {
		putInput.ACL = s3types.ObjectCannedACL(*acl)
	}

	if metadata, ok := d.GetOk("metadata"); ok {
		putInput.Metadata = expandObjectMetadata(metadata.(map[string]any))
		tflog.Debug(ctx, fmt.Sprintf("got Metadata: %v", putInput.Metadata))
	}

	tflog.Debug(ctx, "putting the object", map[string]any{"PutObjectInput": putInput})
	if _, err := s3client.PutObject(ctx, putInput); err != nil {
		return diag.Errorf("failed to put Bucket (%s) Object (%s): %s", bucket, key, err)
	}

	d.SetId(helper.BuildObjectStorageObjectID(d))

	return readResource(ctx, d, meta)
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

	tflog.Debug(ctx, "DeleteObjectInput", map[string]any{"DeleteObjectInput": deleteObjectInput})
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

func objectBodyFromResourceData(d *schema.ResourceData) (body s3manager.ReaderSeekerCloser, err error) {
	if source, ok := d.GetOk("source"); ok {
		sourceFilePath := source.(string)

		file, err := os.Open(filepath.Clean(sourceFilePath))
		if err != nil {
			return s3manager.ReaderSeekerCloser{}, err
		}
		return *s3manager.ReadSeekCloser(file), err
	}

	var contentBytes []byte
	if encodedContent, ok := d.GetOk("content_base64"); ok {
		contentBytes, err = base64.StdEncoding.DecodeString(encodedContent.(string))
	} else {
		content := d.Get("content").(string)
		contentBytes = []byte(content)
	}

	body = *s3manager.ReadSeekCloser(bytes.NewReader(contentBytes))
	return
}

func expandObjectMetadata(metadata map[string]any) map[string]string {
	metadataMap := make(map[string]string, len(metadata))
	for key, value := range metadata {
		metadataMap[key] = value.(string)
	}
	return metadataMap
}

func flattenObjectMetadata(metadata map[string]string) map[string]string {
	metadataObject := make(map[string]string, len(metadata))
	for key, value := range metadata {
		key := strings.ToLower(key)
		metadataObject[key] = value
	}

	return metadataObject
}
