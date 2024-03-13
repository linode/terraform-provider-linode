package objwsconfig

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func GetS3WebsiteDomain(ctx context.Context, linodeClient *linodego.Client, cluster string) (string, error) {
	storageCluster, err := linodeClient.GetObjectStorageCluster(ctx, cluster)
	if err != nil {
		return "", err
	}
	return storageCluster.StaticSiteDomain, nil
}

func s3ConnectionFromData(ctx context.Context, linodeClient *linodego.Client, data ResourceModel) (*s3.Client, error) {
	b, err := linodeClient.GetObjectStorageBucket(ctx, data.Cluster.ValueString(), data.Bucket.ValueString())
	if err != nil {
		return nil, fmt.Errorf("failed to find the specified Linode ObjectStorageBucket: %s", err)
	}
	endpoint := helper.ComputeS3EndpointFromBucket(ctx, *b)
	return helper.S3Connection(ctx, endpoint, data.AccessKey.ValueString(), data.SecretKey.ValueString())
}

func putBucketWebsite(ctx context.Context, linodeClient *linodego.Client, data ResourceModel) error {
	s3client, err := s3ConnectionFromData(ctx, linodeClient, data)
	if err != nil {
		return err
	}

	websiteConfig := &s3types.WebsiteConfiguration{}
	if data.IndexDocument.ValueStringPointer() != nil {
		websiteConfig.IndexDocument = &s3types.IndexDocument{Suffix: aws.String(data.IndexDocument.ValueString())}
	}
	if data.ErrorDocument.ValueStringPointer() != nil {
		websiteConfig.ErrorDocument = &s3types.ErrorDocument{Key: aws.String(data.ErrorDocument.ValueString())}
	}

	_, err = s3client.PutBucketWebsite(ctx, &s3.PutBucketWebsiteInput{
		Bucket:               aws.String(data.Bucket.ValueString()),
		WebsiteConfiguration: websiteConfig,
	})
	return err
}

func deleteBucketWebsite(ctx context.Context, linodeClient *linodego.Client, data ResourceModel) error {
	s3client, err := s3ConnectionFromData(ctx, linodeClient, data)
	if err != nil {
		return err
	}
	_, err = s3client.DeleteBucketWebsite(ctx, &s3.DeleteBucketWebsiteInput{
		Bucket: aws.String(data.Bucket.ValueString()),
	})
	return err
}

func isNotFoundError(err error) bool {
	var re *awshttp.ResponseError
	if errors.As(err, &re) {
		return re.HTTPStatusCode() == 404
	}
	return false
}
