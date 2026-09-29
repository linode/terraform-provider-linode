//go:build unit

package accounttransfer

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestParseAccountTransfer(t *testing.T) {
	t.Parallel()

	transfer := &linodego.AccountTransfer{
		Billable: 0,
		Quota:    9141,
		Used:     2,
		RegionTransfers: []linodego.AccountTransferRegion{
			{
				ID:       "us-east",
				Billable: 0,
				Quota:    5010,
				Used:     1,
			},
			{
				ID:       "us-west",
				Billable: 1,
				Quota:    1000,
				Used:     3,
			},
		},
	}

	var data DataSourceModel
	data.parseAccountTransfer(transfer)

	assert.Equal(t, types.StringValue("account_transfer"), data.ID)
	assert.Equal(t, types.Int64Value(0), data.Billable)
	assert.Equal(t, types.Int64Value(9141), data.Quota)
	assert.Equal(t, types.Int64Value(2), data.Used)
	assert.Len(t, data.RegionTransfers, 2)

	assert.Equal(t, types.StringValue("us-east"), data.RegionTransfers[0].ID)
	assert.Equal(t, types.Int64Value(0), data.RegionTransfers[0].Billable)
	assert.Equal(t, types.Int64Value(5010), data.RegionTransfers[0].Quota)
	assert.Equal(t, types.Int64Value(1), data.RegionTransfers[0].Used)

	assert.Equal(t, types.StringValue("us-west"), data.RegionTransfers[1].ID)
	assert.Equal(t, types.Int64Value(1), data.RegionTransfers[1].Billable)
	assert.Equal(t, types.Int64Value(1000), data.RegionTransfers[1].Quota)
	assert.Equal(t, types.Int64Value(3), data.RegionTransfers[1].Used)
}
