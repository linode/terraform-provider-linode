package objbucket

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"bucket":  d.Get("label"),
		"cluster": d.Get("cluster"),
	})
}

func getS3Endpoint(ctx context.Context, bucket linodego.ObjectStorageBucket) (endpoint string) {
	if bucket.S3Endpoint != "" {
		endpoint = bucket.S3Endpoint
	} else {
		endpoint = helper.ComputeS3EndpointFromBucket(ctx, bucket)
	}
	return endpoint
}

func validateRegion(ctx context.Context, region string, client *linodego.Client) (valid bool, suggestedRegions []string, err error) {
	endpoints, err := client.ListObjectStorageEndpoints(ctx, nil)
	if err != nil {
		return false, nil, err
	}

	for _, endpoint := range endpoints {
		if endpoint.Region == region {
			return true, nil, nil
		} else if endpoint.S3Endpoint != nil && strings.Contains(*endpoint.S3Endpoint, region) {
			suggestedRegions = append(suggestedRegions, endpoint.Region)
		}
	}

	return false, suggestedRegions, nil
}
