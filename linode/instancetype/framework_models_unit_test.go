//go:build unit

package instancetype

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseLinodeType(t *testing.T) {
	mockLinodeType := &linodego.LinodeType{
		ID:         "g6-standard-2",
		Disk:       81920,
		Class:      linodego.ClassStandard,
		Label:      "Linode 4GB",
		NetworkOut: 1000,
		Memory:     4096,
		Transfer:   4000,
		VCPUs:      2,
		GPUs:       0,
		Successor:  "",
		Price: &linodego.LinodePrice{
			Hourly:  0.03,
			Monthly: 20,
		},
		Addons: &linodego.LinodeAddons{
			Backups: &linodego.LinodeBackupsAddon{
				Price: &linodego.LinodePrice{
					Hourly:  0.008,
					Monthly: 5,
				},
			},
		},
	}

	data := &DataSourceModel{}

	diags := data.ParseLinodeType(context.Background(), mockLinodeType)
	assert.False(t, diags.HasError(), "Unexpected error")

	assert.Equal(t, types.StringValue("g6-standard-2"), data.ID)
	assert.Equal(t, types.Int64Value(81920), data.Disk)
	assert.Equal(t, types.StringValue("standard"), data.Class)
	assert.Equal(t, types.StringValue("Linode 4GB"), data.Label)
	assert.Equal(t, types.Int64Value(1000), data.NetworkOut)
	assert.Equal(t, types.Int64Value(4096), data.Memory)
	assert.Equal(t, types.Int64Value(4000), data.Transfer)
	assert.Equal(t, types.Int64Value(2), data.VCPUs)
}
