package monitorlogsstreamsquotas

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

type LogStreamQuotaListModel struct {
	ID     types.String          `tfsdk:"id"`
	Quotas []LogStreamQuotaModel `tfsdk:"quotas"`
}

type LogStreamQuotaModel struct {
	QuotaID     types.String `tfsdk:"quota_id"`
	QuotaName   types.String `tfsdk:"quota_name"`
	Description types.String `tfsdk:"description"`
	QuotaLimit  types.Int64  `tfsdk:"quota_limit"`
	QuotaType   types.String `tfsdk:"quota_type"`
}

func (model *LogStreamQuotaListModel) parseQuotas(quotas []linodego.LogStreamQuota) {
	quotaModels := make([]LogStreamQuotaModel, len(quotas))

	for i, quota := range quotas {
		quotaModels[i] = LogStreamQuotaModel{
			QuotaID:     types.StringValue(quota.QuotaID),
			QuotaName:   types.StringValue(quota.QuotaName),
			Description: types.StringValue(quota.Description),
			QuotaLimit:  types.Int64Value(int64(quota.QuotaLimit)),
			QuotaType:   types.StringValue(quota.QuotaType),
		}
	}

	model.Quotas = quotaModels
}
