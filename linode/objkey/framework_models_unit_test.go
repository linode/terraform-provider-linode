//go:build unit

package objkey

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFlattenObjectStorageKey(t *testing.T) {
	bucketAccessData := []linodego.ObjectStorageKeyBucketAccess{
		{
			Cluster:     "ap-south-1",
			BucketName:  "example-bucket",
			Permissions: "read_only",
		},
	}

	key := linodego.ObjectStorageKey{
		ID:           123,
		Label:        "my-key",
		AccessKey:    "KVAKUTGBA4WTR2NSJQ81",
		SecretKey:    "OiA6F5r0niLs3QA2stbyq7mY5VCV7KqOzcmitmHw",
		Limited:      true,
		BucketAccess: &bucketAccessData,
	}

	data := ResourceModel{}
	var diags diag.Diagnostics
	data.FlattenObjectStorageKey(context.Background(), &key, false, &diags)
	assert.False(t, diags.HasError(), "error flattening obj key")
	assert.True(t, types.StringValue("123").Equal(data.ID))
	assert.Equal(t, types.StringValue("my-key"), data.Label)
	assert.Equal(t, types.StringValue("KVAKUTGBA4WTR2NSJQ81"), data.AccessKey)
	assert.Equal(t, types.StringValue("OiA6F5r0niLs3QA2stbyq7mY5VCV7KqOzcmitmHw"), data.SecretKey)
	assert.Equal(t, types.BoolValue(true), data.Limited)

	assert.NotNil(t, data.BucketAccess)

	bucketAccessEntry := data.BucketAccess[0]
	assert.Equal(t, types.StringValue("ap-south-1"), bucketAccessEntry.Cluster)
	assert.Equal(t, types.StringValue("example-bucket"), bucketAccessEntry.BucketName)
	assert.Equal(t, types.StringValue("read_only"), bucketAccessEntry.Permissions)
}

func TestFlattenObjectStorageKeyPreserveKnown(t *testing.T) {
	key := linodego.ObjectStorageKey{
		ID:           123,
		AccessKey:    "KVAKUTGBA4WTR2NSJQ81",
		SecretKey:    "[REDACTED]",
		Limited:      true,
		BucketAccess: nil,
	}

	expectedID := types.StringValue("123")
	expectedSecretKey := types.StringValue("OiA6F5r0niLs3QA2stbyq7mY5VCV7KqOzcmitmHw")

	rm := ResourceModel{
		ID:        types.StringUnknown(),
		SecretKey: expectedSecretKey,
	}

	var diags diag.Diagnostics
	rm.FlattenObjectStorageKey(context.Background(), &key, true, &diags)
	assert.False(t, diags.HasError(), "error flattening obj key")

	assert.True(t, expectedID.Equal(rm.ID))
	assert.True(t, expectedSecretKey.Equal(rm.SecretKey))
}
