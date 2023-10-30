package obj

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
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
	s3client, err := helper.S3ConnectionFromData(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	headOutput, err := s3client.HeadObject(&s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		// If the object is not found, mark as destroyed so it can be recreated.
		if awsErr, ok := err.(awserr.RequestFailure); ok && awsErr.StatusCode() == http.StatusNotFound {
			d.SetId("")
			fmt.Printf(`[WARN] could not find Bucket (%s) Object (%s)`, bucket, key)
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("cache_control", headOutput.CacheControl)
	d.Set("content_disposition", headOutput.ContentDisposition)
	d.Set("content_encoding", headOutput.ContentEncoding)
	d.Set("content_language", headOutput.ContentLanguage)
	d.Set("content_type", headOutput.ContentType)
	d.Set("etag", strings.Trim(aws.StringValue(headOutput.ETag), `"`))
	d.Set("website_redirect", headOutput.WebsiteRedirectLocation)
	d.Set("version_id", headOutput.VersionId)
	d.Set("metadata", flattenObjectMetadata(headOutput.Metadata))
	d.Set("endpoint", s3client.Config.Endpoint)

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
	acl := d.Get("acl").(string)

	if d.HasChange("acl") {
		s3client, err := helper.S3ConnectionFromData(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		if _, err := s3client.PutObjectAcl(&s3.PutObjectAclInput{
			Bucket: &bucket,
			Key:    &key,
			ACL:    &acl,
		}); err != nil {
			return diag.Errorf("failed to put Bucket (%s) Object (%s) ACL: %s", bucket, key, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	s3client, err := helper.S3ConnectionFromData(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	force := d.Get("force_destroy").(bool)

	if _, ok := d.GetOk("version_id"); ok {
		return deleteAllObjectVersions(s3client, bucket, key, force)
	}

	return diag.FromErr(deleteObject(s3client, bucket, strings.TrimPrefix(key, "/"), "", force))
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
	s3client, err := helper.S3ConnectionFromData(ctx, d, meta)
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

	putInput := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   body,

		ACL:                     nilOrValue(d.Get("acl").(string)),
		CacheControl:            nilOrValue(d.Get("cache_control").(string)),
		ContentDisposition:      nilOrValue(d.Get("content_disposition").(string)),
		ContentEncoding:         nilOrValue(d.Get("content_encoding").(string)),
		ContentLanguage:         nilOrValue(d.Get("content_language").(string)),
		ContentType:             nilOrValue(d.Get("content_type").(string)),
		WebsiteRedirectLocation: nilOrValue(d.Get("website_redirect").(string)),
	}

	if metadata, ok := d.GetOk("metadata"); ok {
		putInput.Metadata = expandObjectMetadata(metadata.(map[string]interface{}))
	}

	if _, err := s3client.PutObject(putInput); err != nil {
		return diag.Errorf("failed to put Bucket (%s) Object (%s): %s", bucket, key, err)
	}

	d.SetId(helper.BuildObjectStorageObjectID(d))

	return readResource(ctx, d, meta)
}

// deleteAllObjectVersions deletes all versions of a given object
func deleteAllObjectVersions(client *s3.S3, bucket, key string, force bool) diag.Diagnostics {
	var versions []string
	listObjectVersionsInput := &s3.ListObjectVersionsInput{
		Bucket: &bucket,
	}

	if key != "" {
		listObjectVersionsInput.Prefix = &key
	}

	// accumulate all versions of the current object to be deleted
	if err := client.ListObjectVersionsPages(
		listObjectVersionsInput, func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
			if page == nil {
				return !lastPage
			}

			for _, objectVersion := range page.Versions {
				if objectKey := aws.StringValue(objectVersion.Key); objectKey == key {
					versions = append(versions, aws.StringValue(objectVersion.VersionId))
				}
			}

			return !lastPage
		}); err != nil {
		if err, ok := err.(awserr.Error); !(ok && err.Code() != s3.ErrCodeNoSuchBucket) {
			return diag.Errorf("failed to list Bucket (%s) Object (%s) versions: %s", bucket, key, err)
		}
	}

	// delete all version of the current object
	for _, version := range versions {
		if err := deleteObject(client, bucket, key, version, force); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func deleteObject(client *s3.S3, bucket, key, version string, force bool) error {
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket:                    &bucket,
		Key:                       &key,
		BypassGovernanceRetention: aws.Bool(force),
	}
	if version != "" {
		deleteObjectInput.VersionId = &version
	}

	_, err := client.DeleteObject(deleteObjectInput)
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf("failed to delete Bucket (%s) Object (%s) Version (%s): %s", bucket, key, version, err)

	//nolint:lll
	if awsErr, ok := err.(awserr.Error); ok && (awsErr.Code() == s3.ErrCodeNoSuchBucket || awsErr.Code() == s3.ErrCodeNoSuchKey) {
		return nil
	} else if ok {
		return awserr.New(awsErr.Code(), msg, awsErr)
	}

	return errors.New(msg)
}

func objectBodyFromResourceData(d *schema.ResourceData) (body aws.ReaderSeekerCloser, err error) {
	if source, ok := d.GetOk("source"); ok {
		sourceFilePath := source.(string)

		file, err := os.Open(filepath.Clean(sourceFilePath))
		if err != nil {
			return aws.ReaderSeekerCloser{}, err
		}
		return aws.ReadSeekCloser(file), err
	}

	var contentBytes []byte
	if encodedContent, ok := d.GetOk("content_base64"); ok {
		contentBytes, err = base64.StdEncoding.DecodeString(encodedContent.(string))
	} else {
		content := d.Get("content").(string)
		contentBytes = []byte(content)
	}

	body = aws.ReadSeekCloser(bytes.NewReader(contentBytes))
	return
}

func expandObjectMetadata(metadata map[string]interface{}) map[string]*string {
	metadataMap := make(map[string]*string, len(metadata))
	for key, value := range metadata {
		metadataMap[key] = aws.String(value.(string))
	}
	return metadataMap
}

func flattenObjectMetadata(metadata map[string]*string) map[string]string {
	metadataObject := make(map[string]string, len(metadata))
	for key, value := range metadata {
		if value == nil {
			continue
		}

		// AWS Go SDK capitalizes metadata, this is a workaround. https://github.com/aws/aws-sdk-go/issues/445
		key := strings.ToLower(key)
		metadataObject[key] = *value
	}

	return metadataObject
}
