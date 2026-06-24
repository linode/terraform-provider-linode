package objkey

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func validateRegionsAgainstBucketAccesses(ctx context.Context, plan ResourceModel, diags *diag.Diagnostics) {
	// regions will be computed if not configured, so it's okay to be null or unknown.
	if plan.BucketAccess == nil || plan.Regions.IsNull() || plan.Regions.IsUnknown() {
		return
	}

	var regions []string
	var bucketRegions []string

	plan.Regions.ElementsAs(ctx, &regions, true)

	for _, ba := range plan.BucketAccess {
		if ba.Region.IsNull() || ba.Region.IsUnknown() {
			continue
		}

		bucketRegion := ba.Region.ValueString()

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
