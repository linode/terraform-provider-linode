package linode

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLinodeObjectStorageObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeObjectStorageObjectCreate,
		Read:   resourceLinodeObjectStorageObjectRead,
		Update: resourceLinodeObjectStorageObjectUpdate,
		Delete: resourceLinodeObjectStorageObjectDelete,

		CustomizeDiff: resourceLinodeObjectStorageObjectCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Description: "The target bucket to put this object in.",
				Required:    true,
			},
			"cluster": {
				Type:        schema.TypeString,
				Description: "The target cluster that the bucket is in.",
				Required:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The name of the uploaded object.",
				Required:    true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The S3 secret key with access to the target bucket.",
				Required:    true,
			},
			"access_key": {
				Type:        schema.TypeString,
				Description: "The S3 access key with access to the target bucket.",
				Required:    true,
			},
			"content": {
				Type:         schema.TypeString,
				Description:  "The contents of the Object to upload.",
				Optional:     true,
				ExactlyOneOf: []string{"content", "content_base64", "source"},
			},
			"content_base64": {
				Type:        schema.TypeString,
				Description: "The base64 contents of the Object to upload.",
				Optional:    true,
			},
			"source": {
				Type:        schema.TypeString,
				Description: "The source file to upload.",
				Optional:    true,
			},
			"acl": {
				Type:        schema.TypeString,
				Description: "The ACL config given to this object.",
				Default:     s3.ObjectCannedACLPrivate,
				Optional:    true,
			},
			"cache_control": {
				Type:        schema.TypeString,
				Description: "This cache_control configuration of this object.",
				Optional:    true,
			},
			"content_disposition": {
				Type:        schema.TypeString,
				Description: "The content disposition configuration of this object.",
				Optional:    true,
			},
			"content_encoding": {
				Type:        schema.TypeString,
				Description: "The encoding of the content of this object.",
				Optional:    true,
			},
			"content_language": {
				Type:        schema.TypeString,
				Description: "The language metadata of this object.",
				Optional:    true,
			},
			"content_type": {
				Type:        schema.TypeString,
				Description: "The MIME type of the content.",
				Optional:    true,
				Computed:    true,
			},
			"etag": {
				Type:        schema.TypeString,
				Description: "The specific version of this object.",
				Optional:    true,
				Computed:    true,
			},
			"force_destroy": {
				Type:        schema.TypeBool,
				Description: "Whether the object should bypass deletion restrictions.",
				Optional:    true,
				Default:     false,
			},
			"metadata": {
				Type:        schema.TypeMap,
				Description: "The metadata of this object",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"version_id": {
				Type:        schema.TypeString,
				Description: "The version ID of this object.",
				Computed:    true,
			},
			"website_redirect": {
				Type:        schema.TypeString,
				Description: "The website redirect location of this object.",
				Optional:    true,
			},
		},
	}
}

func resourceLinodeObjectStorageObjectCreate(d *schema.ResourceData, meta interface{}) error {
	return putLinodeObjectStorageObject(d, meta)
}

func resourceLinodeObjectStorageObjectRead(d *schema.ResourceData, meta interface{}) error {
	client := s3ConnFromResourceData(d)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	headOutput, err := client.HeadObject(&s3.HeadObjectInput{
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
		return err
	}

	d.Set("cache_control", headOutput.CacheControl)
	d.Set("content_disposition", headOutput.ContentDisposition)
	d.Set("content_encoding", headOutput.ContentEncoding)
	d.Set("content_language", headOutput.ContentLanguage)
	d.Set("content_type", headOutput.ContentType)
	d.Set("etag", strings.Trim(aws.StringValue(headOutput.ETag), `"`))
	d.Set("website_redirect", headOutput.WebsiteRedirectLocation)
	d.Set("version_id", headOutput.VersionId)

	d.Set("metadata", flattenLinodeObjectStorageObjectMetadata(headOutput.Metadata))

	return nil
}

func resourceLinodeObjectStorageObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChanges("cache_control", "content_base64", "content_disposition",
		"content_encoding", "content_language", "content_type", "content",
		"etag", "metadata", "source", "website_redirect") {
		return putLinodeObjectStorageObject(d, meta)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	acl := d.Get("acl").(string)

	if d.HasChange("acl") {
		client := s3ConnFromResourceData(d)
		if _, err := client.PutObjectAcl(&s3.PutObjectAclInput{
			Bucket: &bucket,
			Key:    &key,
			ACL:    &acl,
		}); err != nil {
			return fmt.Errorf("failed to put Bucket (%s) Object (%s) ACL: %s", bucket, key, err)
		}
	}

	return resourceLinodeObjectStorageBucketRead(d, meta)
}

func resourceLinodeObjectStorageObjectDelete(d *schema.ResourceData, meta interface{}) (err error) {
	conn := s3ConnFromResourceData(d)

	if _, ok := d.GetOk("version_id"); ok {
		return deleteAllLinodeObjectStorageObjectVersions(d)
	}

	bucket := d.Get("bucket").(string)
	key := strings.TrimPrefix(d.Get("key").(string), "/")
	force := d.Get("force_destroy").(bool)
	return deleteLinodeObjectStorageObject(conn, bucket, key, "", force)
}

func resourceLinodeObjectStorageObjectCustomizeDiff(
	ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.HasChange("etag") {
		d.SetNewComputed("version_id")
	}
	return nil
}

// putLinodeObjectStorageObject builds the object from spec and puts it in the
// specified bucket via the *schema.ResourceData, then it calls
// resourceLinodeObjectStorageObjectRead.
func putLinodeObjectStorageObject(d *schema.ResourceData, meta interface{}) error {
	client := s3ConnFromResourceData(d)
	body, err := objectBodyFromResourceData(d)
	if err != nil {
		return err
	}
	defer body.Close()

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

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
		putInput.Metadata = expandLinodeObjectStorageObjectMetadata(metadata.(map[string]interface{}))
	}

	if _, err := client.PutObject(putInput); err != nil {
		return fmt.Errorf("failed to put Bucket (%s) Object (%s): %s", bucket, key, err)
	}

	d.SetId(buildObjectStorageObjectID(d))

	return resourceLinodeObjectStorageObjectRead(d, meta)
}

// deleteAllLinodeObjectStorageObjectVersions deletes all versions of a given
// object.
func deleteAllLinodeObjectStorageObjectVersions(d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	force := d.Get("force_destroy").(bool)

	conn := s3ConnFromResourceData(d)

	var versions []string
	listObjectVersionsInput := &s3.ListObjectVersionsInput{
		Bucket: &bucket,
	}

	if key != "" {
		listObjectVersionsInput.Prefix = &key
	}

	// accumulate all versions of the current object to be deleted
	if err := conn.ListObjectVersionsPages(
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
			return fmt.Errorf("failed to list Bucket (%s) Object (%s) versions: %s", bucket, key, err)
		}
	}

	// delete all version of the current object
	for _, version := range versions {
		if err := deleteLinodeObjectStorageObject(conn, bucket, key, version, force); err != nil {
			return err
		}
	}
	return nil
}

func deleteLinodeObjectStorageObject(client *s3.S3, bucket, key, version string, force bool) error {
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

		file, err := os.Open(sourceFilePath)
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

func expandLinodeObjectStorageObjectMetadata(metadata map[string]interface{}) map[string]*string {
	metadataMap := make(map[string]*string, len(metadata))
	for key, value := range metadata {
		metadataMap[key] = aws.String(value.(string))
	}
	return metadataMap
}

func flattenLinodeObjectStorageObjectMetadata(metadata map[string]*string) map[string]string {
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
