package helper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func S3Connection(ctx context.Context, endpoint, accessKey, secretKey string) (*s3.Client, error) {
	tflog.Debug(ctx, "creating object storage client")
	awsSDKConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: "https://" + endpoint}, nil
				},
			),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		tflog.Error(ctx, "failed to create object storage client")
		return nil, err
	}

	return s3.NewFromConfig(awsSDKConfig), nil
}

// S3ConnectionFromDataV1 requires endpoint, access_key and secret_key in the data.
// if endpoint is empty a bucket and cluster are required
func S3ConnectionFromData(ctx context.Context, d *schema.ResourceData, meta interface{}) (*s3.Client, error) {
	tflog.Debug(ctx, "creating object storage client from resource data")
	endpoint := d.Get("endpoint").(string)
	if endpoint == "" {
		var err error
		if endpoint, err = ComputeS3Endpoint(ctx, d, meta); err != nil {
			return nil, err
		}
	}

	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)

	return S3Connection(ctx, endpoint, accessKey, secretKey)
}

func ComputeS3Endpoint(ctx context.Context, d *schema.ResourceData, meta interface{}) (string, error) {
	tflog.Debug(ctx, "getting object storage bucket from resource data")
	cluster := d.Get("cluster").(string)
	bucket := d.Get("bucket").(string)

	b, err := meta.(*ProviderMeta).Client.GetObjectStorageBucket(ctx, cluster, bucket)
	if err != nil {
		return "", fmt.Errorf("failed to find the specified Linode ObjectStorageBucket: %s", err)
	}

	return ComputeS3EndpointFromBucket(ctx, *b), nil
}

func ComputeS3EndpointFromBucket(ctx context.Context, bucket linodego.ObjectStorageBucket) string {
	tflog.Debug(ctx, "computing object storage endpoint from bucket instance")
	return strings.TrimPrefix(bucket.Hostname, fmt.Sprintf("%s.", bucket.Label))
}

func BuildObjectStorageObjectID(d *schema.ResourceData) string {
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	return fmt.Sprintf("%s/%s", bucket, key)
}

func IsObjNotFoundErr(err error) bool {
	tflog.Debug(
		context.Background(),
		fmt.Sprintf("received an error: %s, checking whether it's a object not found error", err),
	)
	var apiErr smithy.APIError
	// Error code is 'Forbidden' when the bucket has been removed
	return errors.As(err, &apiErr) && (apiErr.ErrorCode() == "NotFound" || apiErr.ErrorCode() == "Forbidden")
}

// Purge all objects, wiping out all versions and delete markers for versioned objects.
func PurgeAllObjects(
	ctx context.Context,
	bucket string,
	s3client *s3.Client,
	bypassRetention,
	ignoreNotFound bool,
) error {
	tflog.Debug(ctx, fmt.Sprintf("purge all objects in bucket: %s", bucket))

	tflog.Debug(ctx, fmt.Sprintf("getting versioning config of bucket: %s", bucket))
	versioning, err := s3client.GetBucketVersioning(context.Background(), &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	if versioning.Status == s3types.BucketVersioningStatusEnabled {
		tflog.Debug(ctx, fmt.Sprintf("bucket '%s' is a versioned bucket", bucket))
		err = DeleteAllObjectVersionsAndDeleteMarkers(
			context.Background(),
			s3client,
			bucket,
			"",
			bypassRetention,
			ignoreNotFound,
		)
	} else {
		tflog.Debug(ctx, fmt.Sprintf("bucket '%s' isn't a versioned bucket", bucket))
		err = DeleteAllObjects(ctx, s3client, bucket, bypassRetention)
	}
	return err
}

// Send delete requests for every objects.
// Versioned objects will get a deletion marker instead of being fully purged.
func DeleteAllObjects(
	ctx context.Context,
	s3client *s3.Client,
	bucketName string,
	bypassRetention bool,
) error {
	tflog.Debug(ctx, fmt.Sprintf("deleting all objects in bucket '%s'", bucketName))
	objPaginator := s3.NewListObjectsV2Paginator(s3client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	var objectsToDelete []s3types.ObjectIdentifier
	for objPaginator.HasMorePages() {
		tflog.Debug(ctx, fmt.Sprintf("getting next page of the list of objects'%s'", bucketName))
		page, err := objPaginator.NextPage(context.Background())
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			tflog.Debug(ctx, fmt.Sprintf("adding key to delete list: %v", obj.Key))
			objectsToDelete = append(objectsToDelete, s3types.ObjectIdentifier{
				Key: obj.Key,
			})
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("deleting all keys in the list: %v", objectsToDelete))
	_, err := s3client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
		Bucket:                    aws.String(bucketName),
		Delete:                    &s3types.Delete{Objects: objectsToDelete},
		BypassGovernanceRetention: &bypassRetention,
	})

	return err
}

// deleteAllObjectVersions deletes all versions of a given object
func DeleteAllObjectVersionsAndDeleteMarkers(ctx context.Context, client *s3.Client, bucket, key string, bypassRetention, ignoreNotFound bool) error {
	tflog.Debug(ctx, fmt.Sprintf("deleting all versions and deletion marker in bucket '%s'", bucket))
	paginator := s3.NewListObjectVersionsPaginator(client, &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	})

	var objectsToDelete []s3types.ObjectIdentifier
	for paginator.HasMorePages() {
		tflog.Debug(
			ctx,
			fmt.Sprintf("getting next page of the list of versions and delete markers in bucket '%s'", bucket),
		)
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !IsObjNotFoundErr(err) || !ignoreNotFound {
				return err
			}
			tflog.Warn(ctx, fmt.Sprintf("bucket or object does not exist: %v", err))
		}

		for _, version := range page.Versions {
			tflog.Debug(ctx, fmt.Sprintf("adding version '%v' of object key '%v' into delete list", version.VersionId, version.Key))
			objectsToDelete = append(
				objectsToDelete,
				s3types.ObjectIdentifier{
					Key:       version.Key,
					VersionId: version.VersionId,
				},
			)
		}
		for _, marker := range page.DeleteMarkers {
			tflog.Debug(ctx, fmt.Sprintf("adding delete marker '%v' of object key '%v' into delete list", marker.VersionId, marker.Key))
			objectsToDelete = append(
				objectsToDelete,
				s3types.ObjectIdentifier{
					Key:       marker.Key,
					VersionId: marker.VersionId,
				},
			)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("delete all versions and delete markers in the list: %v", objectsToDelete))
	_, err := client.DeleteObjects(
		context.Background(),
		&s3.DeleteObjectsInput{
			Bucket:                    aws.String(bucket),
			Delete:                    &s3types.Delete{Objects: objectsToDelete},
			BypassGovernanceRetention: &bypassRetention,
		},
	)
	if err != nil {
		if !IsObjNotFoundErr(err) || !ignoreNotFound {
			return err
		}
		tflog.Warn(ctx, fmt.Sprintf("bucket or object does not exist: %v", err))
	}
	return nil
}
