package objglobalquotas

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

type ObjectStorageGlobalQuotaFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Quotas  []ObjectStorageGlobalQuotaModel  `tfsdk:"quotas"`
}

type ObjectStorageGlobalQuotaModel struct {
	QuotaID        types.String `tfsdk:"quota_id"`
	QuotaName      types.String `tfsdk:"quota_name"`
	Description    types.String `tfsdk:"description"`
	QuotaLimit     types.Int64  `tfsdk:"quota_limit"`
	ResourceMetric types.String `tfsdk:"resource_metric"`
	QuotaType      types.String `tfsdk:"quota_type"`
	HasUsage       types.Bool   `tfsdk:"has_usage"`
}

func (model *ObjectStorageGlobalQuotaFilterModel) parseQuotas(
	quotas []linodego.ObjectStorageGlobalQuota,
) {
	quotaModels := make([]ObjectStorageGlobalQuotaModel, len(quotas))

	for i, quota := range quotas {
		var quotaModel ObjectStorageGlobalQuotaModel

		quotaModel.QuotaID = types.StringValue(quota.QuotaID)
		quotaModel.QuotaName = types.StringValue(quota.QuotaName)
		quotaModel.Description = types.StringValue(quota.Description)
		quotaModel.QuotaLimit = types.Int64Value(int64(quota.QuotaLimit))
		quotaModel.ResourceMetric = types.StringValue(quota.ResourceMetric)
		quotaModel.QuotaType = types.StringValue(quota.QuotaType)
		quotaModel.HasUsage = types.BoolValue(quota.HasUsage)

		quotaModels[i] = quotaModel
	}

	model.Quotas = quotaModels
}
