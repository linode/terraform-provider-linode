package obj

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

// BucketAccessor provides the minimal information needed to resolve Object Storage
// access keys for an operation.
//
// Implementations typically come from either a resource plan/state model or a
// data model that can supply:
//   - explicit access keys on the resource, or
//   - a bucket label + region/cluster so temporary keys can be created.
type BucketAccessor interface {
	ObjectStorageKeys() ObjectKeys
	BucketLabel() string
	RegionOrCluster(context.Context, *diag.Diagnostics) string
}

// GetObjectStorageKeys resolves Object Storage access keys used by OBJ resources.
//
// Resolution order:
//  1. Keys specified on the resource itself.
//  2. Keys specified in provider configuration (obj_access_key/obj_secret_key).
//  3. If enabled, temporary keys created via the Linode API (obj_use_temp_keys).
//
// When temporary keys are created, a non-nil teardown function is returned to
// delete them after the operation.
func GetObjectStorageKeys(
	ctx context.Context,
	data BucketAccessor,
	client *linodego.Client,
	config *helper.FrameworkProviderModel,
	permissions string,
	endpointType *linodego.ObjectStorageEndpointType,
	diags *diag.Diagnostics,
) (*ObjectKeys, func()) {
	result := data.ObjectStorageKeys()
	if result.Ok() {
		return &result, nil
	}

	result.AccessKey = config.ObjAccessKey.ValueString()
	result.SecretKey = config.ObjSecretKey.ValueString()
	if result.Ok() {
		return &result, nil
	}

	if config.ObjUseTempKeys.ValueBool() {
		clusterOrRegion := data.RegionOrCluster(ctx, diags)
		if diags.HasError() {
			return nil, nil
		}

		objKey := fwCreateTempKeys(ctx, client, data.BucketLabel(), clusterOrRegion, permissions, endpointType, diags)
		if diags.HasError() {
			return nil, nil
		}

		result.AccessKey = objKey.AccessKey
		result.SecretKey = objKey.SecretKey

		teardownTempKeysCleanUp := func() {
			cleanUpTempKeys(ctx, client, objKey.ID)
		}

		return &result, teardownTempKeysCleanUp
	}

	diags.AddError(
		"Keys Not Found",
		"`access_key` and `secret_key` are Required but not Configured",
	)

	return nil, nil
}
