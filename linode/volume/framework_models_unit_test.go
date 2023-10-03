//go:build unit

package volume

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestVolumeModelParsing(t *testing.T) {
	createdTime := time.Date(2023, time.August, 17, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, time.August, 18, 12, 0, 0, 0, time.UTC)
	linodeId := 12346

	data := &VolumeModel{}
	volumeData := &linodego.Volume{
		ID:             12345,
		Label:          "my-volume",
		Status:         "active",
		Region:         "us-east",
		Size:           30,
		LinodeID:       &linodeId,
		FilesystemPath: "/dev/disk/by-id/scsi-0Linode_Volume_my-volume",
		Tags:           []string{"example tag", "another example"},
		Created:        &createdTime,
		Updated:        &updatedTime,
	}

	ctx := context.Background()

	// Test ParseComputedAttributes
	diags := data.ParseComputedAttributes(ctx, volumeData)
	assert.Empty(t, diags)

	assert.Equal(t, types.Int64Value(12345), data.ID)
	assert.Equal(t, types.StringValue("active"), data.Status)
	assert.Equal(t, types.StringValue("us-east"), data.Region)
	assert.Equal(t, types.Int64Value(30), data.Size)
	assert.Equal(t, types.Int64Value(12346), data.LinodeID)
	assert.Equal(t, types.StringValue("/dev/disk/by-id/scsi-0Linode_Volume_my-volume"), data.FilesystemPath)

	assert.NotContains(t, data.Tags.String(), "example tag")
	assert.NotContains(t, data.Tags.String(), "another example")
	assert.NotNil(t, data.Created)
	assert.NotNil(t, data.Updated)

	// Test ParseNonComputedAttributes
	diags = data.ParseNonComputedAttributes(ctx, volumeData)
	assert.Empty(t, diags)
	assert.Contains(t, data.Tags.String(), "example tag")
	assert.Contains(t, data.Tags.String(), "another example")
	assert.Equal(t, types.StringValue("my-volume"), data.Label)
}
