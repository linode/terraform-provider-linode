package objkey

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type RegionDetail struct {
	ID         types.String `tfsdk:"id"`
	S3Endpoint types.String `tfsdk:"s3_endpoint"`
}

type BucketAccessModelEntry struct {
	BucketName  types.String `tfsdk:"bucket_name"`
	Cluster     types.String `tfsdk:"cluster"`
	Permissions types.String `tfsdk:"permissions"`
	Region      types.String `tfsdk:"region"`
}

// ResourceModel describes the Terraform resource rm model to match the
// resource schema.
type ResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Label         types.String `tfsdk:"label"`
	AccessKey     types.String `tfsdk:"access_key"`
	SecretKey     types.String `tfsdk:"secret_key"`
	Limited       types.Bool   `tfsdk:"limited"`
	Regions       types.Set    `tfsdk:"regions"`
	RegionDetails types.Set    `tfsdk:"regions_details"`

	BucketAccess []BucketAccessModelEntry `tfsdk:"bucket_access"`
}

func (plan ResourceModel) GetUpdateOptions(
	ctx context.Context,
	state ResourceModel,
) (opts linodego.ObjectStorageKeyUpdateOptions, shouldUpdate bool) {
	if !state.Label.Equal(plan.Label) {
		opts.Label = plan.Label.ValueString()
		shouldUpdate = true
	}

	if !state.Regions.Equal(plan.Regions) {
		plan.Regions.ElementsAs(ctx, &opts.Regions, false)
		shouldUpdate = true
	}
	return
}

func getObjectStorageKeyRegionIDs(regions []linodego.ObjectStorageKeyRegion) []string {
	regionIDs := make([]string, len(regions))
	for i, r := range regions {
		regionIDs[i] = r.ID
	}
	return regionIDs
}

func getRegionDetails(regions []linodego.ObjectStorageKeyRegion) []RegionDetail {
	regionDetails := make([]RegionDetail, len(regions))
	for i, rd := range regions {
		regionDetails[i] = FlattenRegionDetail(rd)
	}
	return regionDetails
}

func (rm *ResourceModel) FlattenObjectStorageKey(
	ctx context.Context,
	key *linodego.ObjectStorageKey,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	rm.Label = helper.KeepOrUpdateString(rm.Label, key.Label, preserveKnown)

	rm.ID = helper.KeepOrUpdateString(rm.ID, strconv.Itoa(key.ID), preserveKnown)
	rm.AccessKey = helper.KeepOrUpdateString(rm.AccessKey, key.AccessKey, preserveKnown)
	rm.Limited = helper.KeepOrUpdateBool(rm.Limited, key.Limited, preserveKnown)

	// We only want to populate this field if a key is returned,
	// else we should preserve the old value.
	if key.SecretKey != "[REDACTED]" {
		rm.SecretKey = helper.KeepOrUpdateString(rm.SecretKey, key.SecretKey, preserveKnown)
	}

	newRegions := getObjectStorageKeyRegionIDs(key.Regions)
	rm.Regions = helper.KeepOrUpdateStringSet(rm.Regions, newRegions, preserveKnown, diags)

	rm.BucketAccess = FlattenBucketAccessEntries(key.BucketAccess, rm.BucketAccess, preserveKnown)

	regionDetailsSet, newDiags := types.SetValueFrom(ctx, RegionDetailType, getRegionDetails(key.Regions))
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	rm.RegionDetails = helper.KeepOrUpdateValue(rm.RegionDetails, regionDetailsSet, preserveKnown)
}

func FlattenRegionDetail(region linodego.ObjectStorageKeyRegion) (regionDetail RegionDetail) {
	regionDetail.ID = types.StringValue(region.ID)
	regionDetail.S3Endpoint = types.StringValue(region.S3Endpoint)
	return
}

func FlattenBucketAccessEntries(
	entriesPtr *[]linodego.ObjectStorageKeyBucketAccess,
	knownEntries []BucketAccessModelEntry,
	preserveKnown bool,
) (resultEntries []BucketAccessModelEntry) {
	if entriesPtr == nil {
		if preserveKnown {
			return knownEntries
		} else {
			return nil
		}
	}

	entries := *entriesPtr
	if preserveKnown && knownEntries == nil {
		return make([]BucketAccessModelEntry, 0)
	}

	if !preserveKnown {
		resultEntries = make([]BucketAccessModelEntry, len(entries))
	} else {
		resultEntries = knownEntries
	}

	for i := range resultEntries {
		if i > len(entries) {
			break
		}

		resultEntries[i].FlattenBucketAccess(&entries[i], preserveKnown)
	}

	return
}

func (rm *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	rm.ID = helper.KeepOrUpdateValue(rm.ID, other.ID, preserveKnown)
	rm.Label = helper.KeepOrUpdateValue(rm.Label, other.Label, preserveKnown)
	rm.AccessKey = helper.KeepOrUpdateValue(rm.AccessKey, other.AccessKey, preserveKnown)
	rm.SecretKey = helper.KeepOrUpdateValue(rm.SecretKey, other.SecretKey, preserveKnown)
	rm.Limited = helper.KeepOrUpdateValue(rm.Limited, other.Limited, preserveKnown)
	rm.Regions = helper.KeepOrUpdateValue(rm.Regions, other.Regions, preserveKnown)
	rm.RegionDetails = helper.KeepOrUpdateValue(rm.RegionDetails, other.RegionDetails, preserveKnown)

	if !preserveKnown {
		rm.BucketAccess = other.BucketAccess
	}
}

func (b *BucketAccessModelEntry) FlattenBucketAccess(access *linodego.ObjectStorageKeyBucketAccess, preserveKnown bool) {
	b.BucketName = helper.KeepOrUpdateString(b.BucketName, access.BucketName, preserveKnown)
	b.Cluster = helper.KeepOrUpdateString(b.Cluster, access.Cluster, preserveKnown)
	b.Region = helper.KeepOrUpdateString(b.Region, access.Region, preserveKnown)
	b.Permissions = helper.KeepOrUpdateString(b.Permissions, access.Permissions, preserveKnown)
}

func (b *BucketAccessModelEntry) toLinodeObject() linodego.ObjectStorageKeyBucketAccess {
	var result linodego.ObjectStorageKeyBucketAccess

	result.BucketName = b.BucketName.ValueString()
	result.Cluster = b.Cluster.ValueString()
	result.Region = b.Region.ValueString()

	result.Permissions = b.Permissions.ValueString()

	return result
}
