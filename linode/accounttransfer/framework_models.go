package accounttransfer

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

type RegionTransferModel struct {
	ID       types.String `tfsdk:"id"`
	Billable types.Int64  `tfsdk:"billable"`
	Quota    types.Int64  `tfsdk:"quota"`
	Used     types.Int64  `tfsdk:"used"`
}

type DataSourceModel struct {
	ID              types.String          `tfsdk:"id"`
	Billable        types.Int64           `tfsdk:"billable"`
	Quota           types.Int64           `tfsdk:"quota"`
	Used            types.Int64           `tfsdk:"used"`
	RegionTransfers []RegionTransferModel `tfsdk:"region_transfers"`
}

func (data *DataSourceModel) parseAccountTransfer(transfer *linodego.AccountTransfer) {
	data.ID = types.StringValue("account_transfer")
	data.Billable = types.Int64Value(int64(transfer.Billable))
	data.Quota = types.Int64Value(int64(transfer.Quota))
	data.Used = types.Int64Value(int64(transfer.Used))

	regionTransfers := make([]RegionTransferModel, len(transfer.RegionTransfers))
	for i, regionTransfer := range transfer.RegionTransfers {
		regionTransfers[i] = RegionTransferModel{
			ID:       types.StringValue(regionTransfer.ID),
			Billable: types.Int64Value(int64(regionTransfer.Billable)),
			Quota:    types.Int64Value(int64(regionTransfer.Quota)),
			Used:     types.Int64Value(int64(regionTransfer.Used)),
		}
	}
	data.RegionTransfers = regionTransfers
}
