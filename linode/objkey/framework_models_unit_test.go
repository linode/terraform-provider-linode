//go:build unit

package objkey

import (
	"testing"

	"github.com/linode/linodego"
)

func TestParseConfiguredAttributes(t *testing.T) {
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
	data.parseConfiguredAttributes(&key)
	//assert.Equal(t, types.Int64Value(123), data.ID)
	//assert.Equal(t, types.StringValue("my-key"), data.Label)
	//assert.Equal(t, types.StringValue("KVAKUTGBA4WTR2NSJQ81"), data.AccessKey)
	//assert.Equal(t, types.StringValue("OiA6F5r0niLs3QA2stbyq7mY5VCV7KqOzcmitmHw"), data.SecretKey)
	//assert.Equal(t, types.BoolValue(true), data.Limited)

	//assert.NotNil(t, data.BucketAccess)
	//
	//bucketAccessEntry := data.BucketAccess[0]
	//assert.Equal(t, types.StringValue("ap-south-1"), bucketAccessEntry.Cluster)
	//assert.Equal(t, types.StringValue("example-bucket"), bucketAccessEntry.BucketName)
	//assert.Equal(t, types.StringValue("read_only"), bucketAccessEntry.Permissions)
}

func TestParseComputedAttributes(t *testing.T) {
	key := linodego.ObjectStorageKey{
		ID:           123,
		AccessKey:    "KVAKUTGBA4WTR2NSJQ81",
		SecretKey:    "[REDACTED]",
		Limited:      true,
		BucketAccess: nil,
	}

	rm := ResourceModel{}
	rm.parseComputedAttributes(&key)

	//assert.Equal(t, types.Int64Value(123), rm.ID)
	//assert.Equal(t, types.StringValue("KVAKUTGBA4WTR2NSJQ81"), rm.AccessKey)
	//assert.Equal(t, rm.Limited, types.BoolValue(true))
	//
	//assert.Equal(t, types.StringValue(""), rm.SecretKey)
}
