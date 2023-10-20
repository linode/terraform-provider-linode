package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func ComputeS3Endpoint(ctx context.Context, d *schema.ResourceData, meta interface{}) (string, error) {
	cluster := d.Get("cluster").(string)
	bucket := d.Get("bucket").(string)

	b, err := meta.(*ProviderMeta).Client.GetObjectStorageBucket(ctx, cluster, bucket)
	if err != nil {
		return "", fmt.Errorf("failed to find the specified Linode ObjectStorageBucket: %s", err)
	}

	return strings.TrimPrefix(b.Hostname, fmt.Sprintf("%s.", bucket)), nil
}

func BuildObjectStorageObjectID(d *schema.ResourceData) string {
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	return fmt.Sprintf("%s/%s", bucket, key)
}
