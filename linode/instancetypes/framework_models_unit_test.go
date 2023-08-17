package instancetypes

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseInstanceTypes(t *testing.T) {
	mockTypes := []linodego.LinodeType{
		{
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
		},
		{
			ID:         "g2-nanode-1",
			Disk:       81920,
			Class:      linodego.ClassNanode,
			Label:      "Linode 4GB",
			NetworkOut: 1000,
			Memory:     2048,
			Transfer:   4000,
			VCPUs:      2,
			GPUs:       0,
			Successor:  "",
			Price: &linodego.LinodePrice{
				Hourly:  0.01,
				Monthly: 10,
			},
			Addons: &linodego.LinodeAddons{
				Backups: &linodego.LinodeBackupsAddon{
					Price: &linodego.LinodePrice{
						Hourly:  0.008,
						Monthly: 5,
					},
				},
			},
		},
	}

	model := &InstanceTypeFilterModel{}

	diags := model.parseInstanceTypes(context.Background(), mockTypes)
	assert.False(t, diags.HasError(), "Unexpected error")

	assert.Len(t, model.Types, len(mockTypes), "Number of parsed types does not match")

	// Assertions for each mock instance type
	for i, mockType := range mockTypes {
		assert.Equal(t, model.Types[i].ID, types.StringValue(mockType.ID), "ID doesn't match")
		assert.Equal(t, model.Types[i].Disk, types.Int64Value(int64(mockType.Disk)), "Disk size doesn't match")
		assert.Equal(t, model.Types[i].Class, types.StringValue(string(mockType.Class)), "Class doesn't match")
		assert.Equal(t, model.Types[i].Label, types.StringValue(mockType.Label), "Label doesn't match")
		assert.Equal(t, model.Types[i].NetworkOut, types.Int64Value(int64(mockType.NetworkOut)), "NetworkOut doesn't match")
		assert.Equal(t, model.Types[i].Memory, types.Int64Value(int64(mockType.Memory)), "Memory size doesn't match")
		assert.Equal(t, model.Types[i].Transfer, types.Int64Value(int64(mockType.Transfer)), "Transfer size doesn't match")
		assert.Equal(t, model.Types[i].VCPUs, types.Int64Value(int64(mockType.VCPUs)), "VCPUs count doesn't match")

		// Assertions for Price
		assert.NotNil(t, model.Types[i].Price, "Price should not be nil")
		assert.Contains(t, model.Types[i].Price.String(), strconv.FormatFloat(float64(mockType.Price.Hourly), 'f', -1, 32), "Hourly price doesn't match")
		assert.Contains(t, model.Types[i].Price.String(), strconv.FormatFloat(float64(mockType.Price.Monthly), 'f', -1, 32), "Monthly price doesn't match")

		// Assertions for Addons
		assert.Contains(t, model.Types[i].Addons.String(), strconv.FormatFloat(float64(mockType.Addons.Backups.Price.Hourly), 'f', -1, 32), "Backups hourly price doesn't match")
		assert.Contains(t, model.Types[i].Addons.String(), strconv.FormatFloat(float64(mockType.Addons.Backups.Price.Monthly), 'f', -1, 32), "Backups monthly price doesn't match")
	}
}
