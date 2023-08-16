package images

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseImages(t *testing.T) {
	createdTime := &time.Time{}
	createdTimeFormatted := createdTime.Format(time.RFC3339)
	images := []linodego.Image{
		{
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
			Deprecated:   false,
			Created:      createdTime,
			Expiry:       nil,
		},
		{
			ID:           "linode/debian10",
			CreatedBy:    "linode",
			Capabilities: []string{},
			Label:        "Debian 10",
			Description:  "Example Image description 2.",
			Type:         "manual",
			Vendor:       "Debian",
			Status:       "available",
			Size:         2500,
			IsPublic:     true,
			Deprecated:   false,
			Created:      createdTime,
			Expiry:       nil,
		},
	}

	data := ImageFilterModel{}

	data.parseImages(images)

	assert.Len(t, data.Images, len(images))

	// Image 1 Assertions
	assert.Equal(t, types.StringValue("linode/debian11"), data.Images[0].ID)
	assert.Equal(t, types.StringValue("linode"), data.Images[0].CreatedBy)
	assert.Empty(t, data.Images[0].Capabilities)
	assert.Equal(t, types.StringValue("Debian 11"), data.Images[0].Label)
	assert.Equal(t, types.StringValue("Example Image description."), data.Images[0].Description)
	assert.Equal(t, types.StringValue("manual"), data.Images[0].Type)
	assert.Equal(t, types.StringValue("Debian"), data.Images[0].Vendor)
	assert.Equal(t, types.StringValue("available"), data.Images[0].Status)
	assert.Equal(t, types.Int64Value(2500), data.Images[0].Size)
	assert.Equal(t, types.BoolValue(true), data.Images[0].IsPublic)
	assert.Equal(t, types.BoolValue(false), data.Images[0].Deprecated)
	assert.Equal(t, data.Images[0].Created, types.StringValue(createdTimeFormatted))
	assert.Empty(t, data.Images[0].Expiry)

	// Image 2 Assertions
	assert.Equal(t, types.StringValue("linode/debian10"), data.Images[1].ID)
	assert.Equal(t, types.StringValue("linode"), data.Images[1].CreatedBy)
	assert.Empty(t, data.Images[0].Capabilities)
	assert.Equal(t, types.StringValue("Debian 10"), data.Images[1].Label)
	assert.Equal(t, types.StringValue("Example Image description 2."), data.Images[1].Description)
	assert.Equal(t, types.StringValue("manual"), data.Images[1].Type)
	assert.Equal(t, types.StringValue("Debian"), data.Images[1].Vendor)
	assert.Equal(t, types.StringValue("available"), data.Images[1].Status)
	assert.Equal(t, types.Int64Value(2500), data.Images[1].Size)
	assert.Equal(t, types.BoolValue(true), data.Images[1].IsPublic)
	assert.Equal(t, types.BoolValue(false), data.Images[1].Deprecated)
	assert.Equal(t, data.Images[1].Created, types.StringValue(createdTimeFormatted))
	assert.Empty(t, data.Images[1].Expiry)
}
