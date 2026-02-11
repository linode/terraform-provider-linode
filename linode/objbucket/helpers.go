package objbucket

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

// validateRegionIfPresent validates the region specified in the resource data.
// It checks if the region is valid for Object Storage and returns helpful error
// messages with suggested regions if validation fails.
func validateRegionIfPresent(
	ctx context.Context,
	d *schema.ResourceData,
	client *linodego.Client,
) diag.Diagnostics {
	region, ok := d.GetOk("region")
	if !ok {
		return nil
	}

	valid, suggestedRegions, err := validateRegion(ctx, region.(string), client)
	if err != nil {
		return diag.FromErr(err)
	}

	if !valid {
		errorMsg := fmt.Sprintf("Region '%s' is not valid for Object Storage.", region.(string))
		if len(suggestedRegions) > 0 {
			errorMsg += fmt.Sprintf(" Suggested regions: %s", strings.Join(suggestedRegions, ", "))
		}
		return diag.Errorf("%s", errorMsg)
	}

	return nil
}
