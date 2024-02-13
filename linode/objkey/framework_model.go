package objkey

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type BucketAccessModelEntry struct {
	BucketName  types.String `tfsdk:"bucket_name"`
	Cluster     types.String `tfsdk:"cluster"`
	Permissions types.String `tfsdk:"permissions"`
}

// ResourceModel describes the Terraform resource rm model to match the
// resource schema.
type ResourceModel struct {
	ID           types.String             `tfsdk:"id"`
	Label        types.String             `tfsdk:"label"`
	AccessKey    types.String             `tfsdk:"access_key"`
	SecretKey    types.String             `tfsdk:"secret_key"`
	Limited      types.Bool               `tfsdk:"limited"`
	BucketAccess []BucketAccessModelEntry `tfsdk:"bucket_access"`
}

func (rm *ResourceModel) FlattenObjectStorageKey(key *linodego.ObjectStorageKey, preserveKnown bool) {
	rm.Label = helper.KeepOrUpdateString(rm.Label, key.Label, preserveKnown)

	rm.ID = helper.KeepOrUpdateString(rm.ID, strconv.Itoa(key.ID), preserveKnown)
	rm.AccessKey = helper.KeepOrUpdateString(rm.AccessKey, key.AccessKey, preserveKnown)
	rm.Limited = helper.KeepOrUpdateBool(rm.Limited, key.Limited, preserveKnown)

	// We only want to populate this field if a key is returned,
	// else we should preserve the old value.
	if key.SecretKey != "[REDACTED]" {
		rm.SecretKey = helper.KeepOrUpdateString(rm.SecretKey, key.SecretKey, preserveKnown)
	}

	// rm.BucketAccess should only be changed when known values are not preserved
	if !preserveKnown {
		// No access is configured; we can return here
		if key.BucketAccess == nil {
			rm.BucketAccess = nil
			return
		}

		keyBucketAccess := *key.BucketAccess
		bucketAccess := make([]BucketAccessModelEntry, len(keyBucketAccess))

		for i := range bucketAccess {
			var entry BucketAccessModelEntry
			entry.FlattenBucketAccess(&keyBucketAccess[i], preserveKnown)
			bucketAccess[i] = entry
		}

		rm.BucketAccess = bucketAccess
	}
}

func (rm *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	rm.ID = helper.KeepOrUpdateValue(rm.ID, other.ID, preserveKnown)
	rm.Label = helper.KeepOrUpdateValue(rm.Label, other.Label, preserveKnown)
	rm.AccessKey = helper.KeepOrUpdateValue(rm.AccessKey, other.AccessKey, preserveKnown)
	rm.SecretKey = helper.KeepOrUpdateValue(rm.SecretKey, other.SecretKey, preserveKnown)
	rm.Limited = helper.KeepOrUpdateValue(rm.Limited, other.Limited, preserveKnown)
	if !preserveKnown {
		rm.BucketAccess = other.BucketAccess
	}
}

func (b *BucketAccessModelEntry) FlattenBucketAccess(access *linodego.ObjectStorageKeyBucketAccess, preserveKnown bool) {
	b.BucketName = helper.KeepOrUpdateString(b.BucketName, access.BucketName, preserveKnown)
	b.Cluster = helper.KeepOrUpdateString(b.Cluster, access.Cluster, preserveKnown)
	b.Permissions = helper.KeepOrUpdateString(b.Permissions, access.Permissions, preserveKnown)
}

func (b *BucketAccessModelEntry) toLinodeObject() linodego.ObjectStorageKeyBucketAccess {
	var result linodego.ObjectStorageKeyBucketAccess

	result.BucketName = b.BucketName.ValueString()
	result.Cluster = b.Cluster.ValueString()
	result.Permissions = b.Permissions.ValueString()

	return result
}
