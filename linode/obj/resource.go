package obj

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/smithy-go"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
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

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return putObject(ctx, d, meta)
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	s3client, err := helper.S3ConnectionFromDataV2(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	headOutput, err := s3client.HeadObject(
		ctx,
		&s3v2.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NoSuchKey" || apiErr.ErrorCode() == "NoSuchKey" {
				d.SetId("")
				tflog.Warn(ctx, fmt.Sprintf("couldn't find Bucket (%s) or Object (%s)", bucket, key))
				return nil
			}
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
		endpoint, err := helper.ComputeS3Endpoint(ctx, d, meta)
		if err != nil {
			return diag.Errorf("failed to compute object storage endpoint: %s", err)
		}
		d.Set("endpoint", endpoint)
	}

	return nil
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChanges("cache_control", "content_base64", "content_disposition",
		"content_encoding", "content_language", "content_type", "content",
		"etag", "metadata", "source", "website_redirect") {
		return putObject(ctx, d, meta)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	acl := s3types.ObjectCannedACL(d.Get("acl").(string))

	if d.HasChange("acl") {
		s3client, err := helper.S3ConnectionFromDataV2(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = s3client.PutObjectAcl(
			ctx, &s3v2.PutObjectAclInput{
				Bucket: &bucket,
				Key:    &key,
				ACL:    acl,
			},
		)
		if err != nil {
			return diag.Errorf("failed to put Bucket (%s) Object (%s) ACL: %s", bucket, key, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	s3client, err := helper.S3ConnectionFromDataV2(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	force := d.Get("force_destroy").(bool)

	if _, ok := d.GetOk("version_id"); ok {
		return deleteAllObjectVersions(ctx, s3client, bucket, key, force)
	}

	return diag.FromErr(deleteObject(ctx, s3client, bucket, strings.TrimPrefix(key, "/"), "", force))
}

func diffResource(
	ctx context.Context, d *schema.ResourceDiff, meta interface{},
) error {
	if d.HasChange("etag") {
		d.SetNewComputed("version_id")
	}
	return nil
}

// putObject builds the object from spec and puts it in the
// specified bucket via the *schema.ResourceData, then it calls
// readResource.
func putObject(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	s3client, err := helper.S3ConnectionFromDataV2(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

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

	putInputV2 := &s3v2.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   &body,

		ACL:                     s3types.ObjectCannedACL(d.Get("acl").(string)),
		CacheControl:            nilOrValue(d.Get("cache_control").(string)),
		ContentDisposition:      nilOrValue(d.Get("content_disposition").(string)),
		ContentEncoding:         nilOrValue(d.Get("content_encoding").(string)),
		ContentLanguage:         nilOrValue(d.Get("content_language").(string)),
		ContentType:             nilOrValue(d.Get("content_type").(string)),
		WebsiteRedirectLocation: nilOrValue(d.Get("website_redirect").(string)),
	}

	if metadata, ok := d.GetOk("metadata"); ok {
		putInputV2.Metadata = expandObjectMetadata(metadata.(map[string]interface{}))
	}

	if _, err := s3client.PutObject(ctx, putInputV2); err != nil {
		return diag.Errorf("failed to put Bucket (%s) Object (%s): %s", bucket, key, err)
	}

	d.SetId(helper.BuildObjectStorageObjectID(d))

	return readResource(ctx, d, meta)
}

// deleteAllObjectVersions deletes all versions of a given object
func deleteAllObjectVersions(ctx context.Context, client *s3v2.Client, bucket, key string, force bool) diag.Diagnostics {
	paginator := s3v2.NewListObjectVersionsPaginator(client, &s3v2.ListObjectVersionsInput{
		Bucket: awsv2.String(bucket),
		Prefix: awsv2.String(key),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		shouldIgnore := func(code string) bool {
			return (code == "NoSuchBucket" ||
				code == "NoSuchKey" ||
				code == "NoSuchVersion")
		}

		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && shouldIgnore(apiErr.ErrorCode()) {
				tflog.Error(
					ctx,
					fmt.Sprintf("bucket or object does not exist: %v", err),
				)
			} else {
				return diag.FromErr(err)
			}
		}

		for _, version := range page.Versions {
			_, err := client.DeleteObject(
				context.Background(),
				&s3v2.DeleteObjectInput{
					Bucket:    awsv2.String(bucket),
					Key:       version.Key,
					VersionId: version.VersionId,
				})
			if err != nil {
				var apiErr smithy.APIError
				if errors.As(err, &apiErr) && shouldIgnore(apiErr.ErrorCode()) {
					tflog.Error(
						ctx,
						fmt.Sprintf("bucket or object does not exist: %v", err),
					)
				} else {
					return diag.FromErr(err)
				}
			}
		}
	}
	return nil
}

func deleteObject(ctx context.Context, client *s3v2.Client, bucket, key, version string, force bool) error {
	deleteObjectInput := &s3v2.DeleteObjectInput{
		Bucket:                    &bucket,
		Key:                       &key,
		BypassGovernanceRetention: awsv2.Bool(force),
	}
	if version != "" {
		deleteObjectInput.VersionId = &version
	}

	_, err := client.DeleteObject(ctx, deleteObjectInput)
	if err != nil {
		msg := fmt.Sprintf(
			"failed to delete Bucket (%s) Object (%s) Version (%s): %s",
			bucket,
			key,
			version,
			err,
		)

		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() != "NoSuchBucket" && apiErr.ErrorCode() != "NoSuchKey" {
				return fmt.Errorf("%s: %w", msg, err)
			}
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

func expandObjectMetadata(metadata map[string]interface{}) map[string]string {
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
