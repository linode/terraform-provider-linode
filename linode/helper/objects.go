package helper

import (
	"context"
	"fmt"
	"strings"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	credentialsv2 "github.com/aws/aws-sdk-go-v2/credentials"
	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

const (
	LinodeObjectsEndpoint = "https://%s.linodeobjects.com"
)

// S3Connection create a client for s3 from an endpoint and keys
func S3Connection(endpoint, accessKey, secretKey string) (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		// This region is hardcoded strictly for preflight validation purposes.
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String(endpoint),
	})
	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}

func S3ConnectionV2(endpoint, accessKey, secretKey string) (*s3v2.Client, error) {
	awsSDKConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(
			credentialsv2.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}

	return s3v2.NewFromConfig(awsSDKConfig, func(o *s3v2.Options) {
		o.BaseEndpoint = awsv2.String("https://" + endpoint)
	}), nil
}

// S3ConnectionFromData requires endpoint, access_key and secret_key in the data.
// if endpoint is empty a bucket and cluster are required
func S3ConnectionFromData(ctx context.Context, d *schema.ResourceData, meta interface{}) (*s3.S3, error) {
	endpoint := d.Get("endpoint").(string)
	if endpoint == "" {
		var err error
		if endpoint, err = ComputeS3Endpoint(ctx, d, meta); err != nil {
			return nil, err
		}
	}

	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)

	return S3Connection(endpoint, accessKey, secretKey)
}

// S3ConnectionFromDataV1 requires endpoint, access_key and secret_key in the data.
// if endpoint is empty a bucket and cluster are required
func S3ConnectionFromDataV2(ctx context.Context, d *schema.ResourceData, meta interface{}) (*s3v2.Client, error) {
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

func DeleteAllObjects(bucket string, client *s3v2.Client) error {
	paginator := s3v2.NewListObjectsV2Paginator(client, &s3v2.ListObjectsV2Input{
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
			&s3v2.DeleteObjectsInput{
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
