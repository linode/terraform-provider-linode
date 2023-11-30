package helper

import (
	"context"
	"fmt"
	"strings"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func S3ConnectionV2(endpoint, accessKey, secretKey string) (*s3.Client, error) {
	awsSDKConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: "https://" + endpoint,
				}, nil
			})),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsSDKConfig), nil
}

// S3ConnectionFromDataV1 requires endpoint, access_key and secret_key in the data.
// if endpoint is empty a bucket and cluster are required
func S3ConnectionFromDataV2(ctx context.Context, d *schema.ResourceData, meta interface{}) (*s3.Client, error) {
	endpoint := d.Get("endpoint").(string)
	if endpoint == "" {
		var err error
		if endpoint, err = ComputeS3Endpoint(ctx, d, meta); err != nil {
			return nil, err
		}
	}

	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)

	return S3ConnectionV2(endpoint, accessKey, secretKey)
}

func ComputeS3Endpoint(ctx context.Context, d *schema.ResourceData, meta interface{}) (string, error) {
	cluster := d.Get("cluster").(string)
	bucket := d.Get("bucket").(string)

	b, err := meta.(*ProviderMeta).Client.GetObjectStorageBucket(ctx, cluster, bucket)
	if err != nil {
		return "", fmt.Errorf("failed to find the specified Linode ObjectStorageBucket: %s", err)
	}

	return ComputeS3EndpointFromBucket(*b), nil
}

func ComputeS3EndpointFromBucket(bucket linodego.ObjectStorageBucket) string {
	return strings.TrimPrefix(bucket.Hostname, fmt.Sprintf("%s.", bucket.Label))
}

func BuildObjectStorageObjectID(d *schema.ResourceData) string {
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	return fmt.Sprintf("%s/%s", bucket, key)
}

func DeleteAllObjects(bucket string, client *s3.Client) error {
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return err
		}

		objects := make([]s3types.ObjectIdentifier, len(page.Contents))
		for i, object := range page.Contents {
			objects[i] = s3types.ObjectIdentifier{
				Key: object.Key,
			}
		}

		_, err = client.DeleteObjects(
			context.Background(),
			&s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: &s3types.Delete{
					Objects: objects,
				},
			})
		if err != nil {
			return err
		}
	}
	return nil
}
