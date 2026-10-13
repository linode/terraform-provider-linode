//go:build unit

package nfsfilesystem

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestFlattenNFSFilesystem(t *testing.T) {
	created := time.Date(2026, time.May, 1, 14, 0, 0, 0, time.UTC)
	updated := time.Date(2026, time.May, 10, 10, 0, 0, 0, time.UTC)
	collected := time.Date(2026, time.May, 21, 10, 0, 0, 0, time.UTC)
	fqdn := "prod-ml-space-3039.nfs.us-rno.linode.com:/training-data-10932"
	snapshotUsage := int64(10737418240)
	usedCapacity := int64(107374182400)
	usedFiles := int64(250000)

	fs := &linodego.NFSFilesystem{
		ID:                 67890,
		SpaceID:            12345,
		Label:              "training-data",
		Region:             "us-rno",
		ProtocolVersions:   []linodego.NFSProtocolVersion{linodego.NFSProtocolVersionV4},
		Status:             linodego.NFSFilesystemStatusActive,
		MountTargetIPs:     []string{"[2400:8900::10]:/training-data-10932"},
		MountTargetFQDN:    &fqdn,
		MaxCapacityBytes:   1099511627776,
		SnapshotUsageBytes: &snapshotUsage,
		Created:            &created,
		Updated:            &updated,
		Tags:               []string{"training", "ml"},
		Stats: linodego.NFSFilesystemStats{
			UsedCapacityBytes: &usedCapacity,
			UsedFileCount:     &usedFiles,
			CollectedAt:       &collected,
		},
	}

	var model BaseModel
	diags := model.Flatten(context.Background(), fs, false)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, types.StringValue("67890"), model.ID)
	assert.Equal(t, types.Int64Value(12345), model.SpaceID)
	assert.Equal(t, types.StringValue("training-data"), model.Label)
	assert.Equal(t, types.StringValue("us-rno"), model.Region)
	assert.Equal(t, types.StringValue("active"), model.Status)
	assert.Equal(t, types.StringValue(fqdn), model.MountTargetFQDN)
	assert.Equal(t, types.Int64Value(1099511627776), model.MaxCapacityBytes)
	assert.Equal(t, types.Int64Value(snapshotUsage), model.SnapshotUsageBytes)
	assert.False(t, model.ProtocolVersions.IsNull())
	assert.False(t, model.MountTargetIPs.IsNull())
	assert.False(t, model.Tags.IsNull())
	assert.False(t, model.Stats.IsNull())
}

func TestFlattenNFSFilesystemPreserveKnown(t *testing.T) {
	fs := &linodego.NFSFilesystem{
		ID:    67890,
		Label: "training-data",
	}

	model := BaseModel{
		ID:    types.StringValue("11111"),
		Label: types.StringUnknown(),
	}

	diags := model.Flatten(context.Background(), fs, true)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, types.StringValue("11111"), model.ID)
	assert.Equal(t, types.StringValue("training-data"), model.Label)
}
