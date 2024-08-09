package objkey

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func (r Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	validateRegionsAgainstBucketAccessEntries(ctx, req.Config, &resp.Diagnostics)
}

func validateRegionsAgainstBucketAccessEntries(ctx context.Context, config tfsdk.Config, diags *diag.Diagnostics) {
	var bucketAccesses types.Set
	var regions []string
	var bucketRegions []string

	diags.Append(config.GetAttribute(ctx, path.Root("bucket_access"), &bucketAccesses)...)
	diags.Append(config.GetAttribute(ctx, path.Root("regions"), &regions)...)

	if diags.HasError() {
		return
	}

	if bucketAccesses.IsNull() || bucketAccesses.IsUnknown() || regions == nil {
		return
	}

	for _, ba := range bucketAccesses.Elements() {
		if baObj, ok := ba.(types.Object); ok {
			if region, ok := baObj.Attributes()["region"]; ok {
				if regionStringValue, ok := region.(types.String); ok {
					bucketRegions = append(bucketRegions, regionStringValue.ValueString())
				}
			}
		}
	}

	if !helper.ValidateStringSubset(regions, bucketRegions) {
		diags.AddAttributeError(
			path.Root("regions"),
			"Incomplete Regions",
			"All regions of the buckets defined in `bucket_access` blocks are expected in the `regions` attribute.",
		)
	}
}
