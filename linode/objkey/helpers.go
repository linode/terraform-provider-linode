package objkey

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func getRegionFromCluster(cluster string) (string, error) {
	s := strings.Split(cluster, "-")
	if len(s) <= 2 {
		return "", fmt.Errorf("failed to parse cluster %q", cluster)
	}
	return strings.Join(s[:2], "-"), nil
}

func validateRegionsAgainstBucketAccesses(ctx context.Context, plan ResourceModel, diags *diag.Diagnostics) {
	// regions will be computed if not configured, so it's okay to be null or unknown.
	if plan.BucketAccess == nil || plan.Regions.IsNull() || plan.Regions.IsUnknown() {
		return
	}

	var regions []string
	var bucketRegions []string

	plan.Regions.ElementsAs(ctx, &regions, true)

	for _, ba := range plan.BucketAccess {
		var bucketRegion string
		var err error

		if ba.Region.IsNull() || ba.Region.IsUnknown() {
			bucketRegion, err = getRegionFromCluster(ba.Cluster.ValueString())
			if err != nil {
				diags.AddWarning("Failed to Parse Cluster", err.Error())
				continue
			}
		} else {
			bucketRegion = ba.Region.ValueString()
		}

		if !slices.Contains(bucketRegions, bucketRegion) {
			bucketRegions = append(bucketRegions, bucketRegion)
		}
	}

	if !helper.ValidateStringSubset(regions, bucketRegions) {
		diags.AddAttributeError(
			path.Root("regions"),
			"Incomplete Regions",
			"All regions of the buckets defined in `bucket_access` blocks are expected in the `regions` set attribute.\n"+
				fmt.Sprintf("Regions in the `regions` set attribute: %v\n", regions)+
				fmt.Sprintf("Regions in the `bucket_access` blocks: %v\n", bucketRegions),
		)
	}
}
