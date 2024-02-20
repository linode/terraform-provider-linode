package instancesharedips

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type ResourceModel struct {
	ID        types.String `tfsdk:"id"`
	LinodeID  types.Int64  `tfsdk:"linode_id"`
	Addresses types.Set    `tfsdk:"addresses"`
}

func (data *ResourceModel) FlattenSharedIPs(
	linodeID int, sharedIPs []string, preserveKnown bool, diags *diag.Diagnostics,
) {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(linodeID), preserveKnown)
	data.LinodeID = helper.KeepOrUpdateInt64(data.LinodeID, int64(linodeID), preserveKnown)
	data.Addresses = helper.KeepOrUpdateStringSet(data.Addresses, sharedIPs, preserveKnown, diags)
}

func (data *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.LinodeID = helper.KeepOrUpdateValue(data.LinodeID, other.LinodeID, preserveKnown)
	data.Addresses = helper.KeepOrUpdateValue(data.Addresses, other.Addresses, preserveKnown)
}
