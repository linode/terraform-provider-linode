package objkey

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
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

func (rm *ResourceModel) parseConfiguredAttributes(key *linodego.ObjectStorageKey) {
	rm.Label = types.StringValue(key.Label)
	// No access is configured; we can return here
	if key.BucketAccess == nil {
		rm.BucketAccess = nil
		return
	}

	bucketAccess := make([]BucketAccessModelEntry, len(*key.BucketAccess))

	keyBucketAccess := *key.BucketAccess

	for i := range keyBucketAccess {
		var entry BucketAccessModelEntry

		entry.parseBucketAccess(&keyBucketAccess[i])

		bucketAccess[i] = entry
	}
}

func (rm *ResourceModel) parseComputedAttributes(key *linodego.ObjectStorageKey) {
	rm.ID = types.StringValue(strconv.Itoa(key.ID))
	rm.AccessKey = types.StringValue(key.AccessKey)
	rm.Limited = types.BoolValue(key.Limited)

	// We only want to populate this field if a key is returned,
	// else we should preserve the old value.
	if key.SecretKey != "[REDACTED]" {
		rm.SecretKey = types.StringValue(key.SecretKey)
	}
}

func (b *BucketAccessModelEntry) parseBucketAccess(access *linodego.ObjectStorageKeyBucketAccess) {
	b.BucketName = types.StringValue(access.BucketName)
	b.Cluster = types.StringValue(access.Cluster)
	b.Permissions = types.StringValue(access.Permissions)
}

func (b *BucketAccessModelEntry) toLinodeObject() linodego.ObjectStorageKeyBucketAccess {
	var result linodego.ObjectStorageKeyBucketAccess

	result.BucketName = b.BucketName.ValueString()
	result.Cluster = b.Cluster.ValueString()
	result.Permissions = b.Permissions.ValueString()

	return result
}
