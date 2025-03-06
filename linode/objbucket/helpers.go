package objbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
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
	return
}
