//go:build unit

package backup

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseBackups(t *testing.T) {
	ctx := context.Background()

	// Create mock data for InstanceSnapshot
	mockSnapshot := &linodego.InstanceSnapshot{
		ID:        1,
		Label:     "Linode Snapshot Label",
		Status:    "successful",
		Type:      "snapshot",
		Created:   nil,
		Updated:   nil,
		Finished:  nil,
		Configs:   []string{"config1", "config2"},
		Disks:     []*linodego.InstanceSnapshotDisk{}, // You can populate this with mock disk data
		Available: true,
	}

	mockBackupSnapshotResponse := &linodego.InstanceBackupSnapshotResponse{
		Current:    mockSnapshot,
		InProgress: mockSnapshot,
	}

	mockBackups := &linodego.InstanceBackupsResponse{
		Automatic: []*linodego.InstanceSnapshot{mockSnapshot},
		Snapshot:  mockBackupSnapshotResponse,
	}

	linodeId := int64(123)

	data := &DataSourceModel{}

	diags := data.parseBackups(ctx, mockBackups, types.Int64Value(linodeId))

	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	assert.Equal(t, types.Int64Value(linodeId), data.ID)
}
