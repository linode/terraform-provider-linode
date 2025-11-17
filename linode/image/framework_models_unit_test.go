//go:build unit

package image

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseImage(t *testing.T) {
	createdTime := &time.Time{}
	createdTimeFormatted := createdTime.Format(time.RFC3339)
	mockImage := linodego.Image{
		ID:           "linode/debian11",
		CreatedBy:    "linode",
		Capabilities: []string{},
		Label:        "Debian 11",
		Description:  "Example Image description.",
		Type:         "manual",
		Vendor:       "Debian",
		Status:       "available",
		Size:         2500,
		IsPublic:     true,
		IsShared:     false,
		Deprecated:   false,
		Created:      createdTime,
		Expiry:       nil,
		TotalSize:    2500,
		Tags:         []string{"test"},
		Regions: []linodego.ImageRegion{
			{
				Region: "us-east",
				Status: linodego.ImageRegionStatus("available"),
			},
			{
				Region: "us-west",
				Status: linodego.ImageRegionStatus("pending replication"),
			},
		},
		ImageSharing: linodego.ImageSharing{
			SharedWith: linodego.Pointer(linodego.ImageSharingSharedWith{
				ShareGroupCount:   1,
				ShareGroupListURL: "/images/private/1234/sharegroups",
			}),
			SharedBy: linodego.Pointer(linodego.ImageSharingSharedBy{
				ShareGroupID:    1,
				ShareGroupUUID:  "0ee8e1c1-b19b-4052-9487-e3b13faac111",
				ShareGroupLabel: "my-label",
				SourceImageID:   linodego.Pointer("private/1234"),
			}),
		},
	}

	var imageModel ImageModel
	imageModel.ParseImage(context.Background(), &mockImage)

	assert.Equal(t, types.StringValue("linode/debian11"), imageModel.ID)
	assert.Equal(t, types.StringValue("linode"), imageModel.CreatedBy)
	assert.Empty(t, imageModel.Capabilities)
	assert.Equal(t, types.StringValue("Debian 11"), imageModel.Label)
	assert.Equal(t, types.StringValue("Example Image description."), imageModel.Description)
	assert.Equal(t, types.StringValue("manual"), imageModel.Type)
	assert.Equal(t, types.StringValue("Debian"), imageModel.Vendor)
	assert.Equal(t, types.StringValue("available"), imageModel.Status)
	assert.Equal(t, types.Int64Value(2500), imageModel.Size)
	assert.Equal(t, types.BoolValue(true), imageModel.IsPublic)
	assert.Equal(t, types.BoolValue(false), imageModel.IsShared)
	assert.Equal(t, types.BoolValue(false), imageModel.Deprecated)
	assert.Equal(t, imageModel.Created, types.StringValue(createdTimeFormatted))
	assert.Empty(t, imageModel.Expiry)
	assert.Equal(t, types.Int64Value(2500), imageModel.TotalSize)
	assert.Equal(t, types.StringValue("us-east"), imageModel.Replications[0].Region)
	assert.Equal(t, types.StringValue("available"), imageModel.Replications[0].Status)
	assert.Equal(t, types.StringValue("us-west"), imageModel.Replications[1].Region)
	assert.Equal(t, types.StringValue("pending replication"), imageModel.Replications[1].Status)
	assert.Equal(t, types.Int64Value(1), imageModel.ImageSharing.SharedWith.ShareGroupCount)
	assert.Equal(t, types.StringValue("/images/private/1234/sharegroups"), imageModel.ImageSharing.SharedWith.ShareGroupListURL)
	assert.Equal(t, types.Int64Value(1), imageModel.ImageSharing.SharedBy.ShareGroupID)
	assert.Equal(t, types.StringValue("0ee8e1c1-b19b-4052-9487-e3b13faac111"), imageModel.ImageSharing.SharedBy.ShareGroupUUID)
	assert.Equal(t, types.StringValue("my-label"), imageModel.ImageSharing.SharedBy.ShareGroupLabel)
	assert.Equal(t, types.StringValue("private/1234"), imageModel.ImageSharing.SharedBy.SourceImageID)
	assert.Contains(t, imageModel.Tags.String(), "test")
}
