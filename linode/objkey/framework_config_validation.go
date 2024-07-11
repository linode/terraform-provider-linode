package objkey

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func (r Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	validateRegionsAgainstBucketAccessEntries(ctx, config, &resp.Diagnostics)
}

func validateRegionsAgainstBucketAccessEntries(ctx context.Context, config ResourceModel, diags *diag.Diagnostics) {
	// regions will be computed if not configured, so it's okay to be null.
	if config.BucketAccess == nil || config.Regions.IsNull() {
		return
	}

	var regions []string
	var bucketRegions []string

	config.Regions.ElementsAs(ctx, &regions, true)

	for _, ba := range config.BucketAccess {
		bucketRegions = append(bucketRegions, ba.Region.ValueString())
	}

	if !helper.ValidateStringSubset(regions, bucketRegions) {
		diags.AddAttributeError(
			path.Root("regions"),
			"Incomplete Regions",
			"All regions of the buckets defined in `bucket_access` blocks are expected in the `regions` attribute.",
		)
	}
}
