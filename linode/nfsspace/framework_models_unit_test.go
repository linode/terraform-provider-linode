//go:build unit

package nfsspace

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestFlattenNFSSpace(t *testing.T) {
	description := "test description"
	created := time.Date(2026, time.May, 1, 12, 0, 0, 0, time.UTC)
	updated := time.Date(2026, time.May, 10, 9, 30, 0, 0, time.UTC)

	space := &linodego.NFSSpace{
		ID:                12345,
		Label:             "prod-ml-space",
		Description:       &description,
		Status:            linodego.NFSSpaceStatusActive,
		MaxCapacityBytes:  10995116277760,
		UsedCapacityBytes: 5497558138880,
		Created:           &created,
		Updated:           &updated,
		Tags:              []string{"prod", "ml"},
	}

	var model BaseModel
	diags := model.Flatten(context.Background(), space, false)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, types.StringValue("12345"), model.ID)
	assert.Equal(t, types.StringValue("prod-ml-space"), model.Label)
	assert.Equal(t, types.StringValue("test description"), model.Description)
	assert.Equal(t, types.StringValue("active"), model.Status)
	assert.Equal(t, types.Int64Value(10995116277760), model.MaxCapacityBytes)
	assert.Equal(t, types.Int64Value(5497558138880), model.UsedCapacityBytes)
	assert.False(t, model.Tags.IsNull())
}

func TestFlattenNFSSpacePreserveKnown(t *testing.T) {
	space := &linodego.NFSSpace{
		ID:    12345,
		Label: "prod-ml-space",
	}

	model := BaseModel{
		ID:    types.StringValue("99999"),
		Label: types.StringUnknown(),
	}

	diags := model.Flatten(context.Background(), space, true)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	// Known values should be preserved.
	assert.Equal(t, types.StringValue("99999"), model.ID)
	// Unknown values should be updated.
	assert.Equal(t, types.StringValue("prod-ml-space"), model.Label)
}
