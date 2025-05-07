package objquotas

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type ObjectStorageQuotaFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Quotas  []ObjectStorageQuotaModel        `tfsdk:"quotas"`
}

type ObjectStorageQuotaModel struct {
	QuotaID        types.String `tfsdk:"quota_id"`
	QuotaName      types.String `tfsdk:"quota_name"`
	EndpointType   types.String `tfsdk:"endpoint_type"`
	S3Endpoint     types.String `tfsdk:"s3_endpoint"`
	Description    types.String `tfsdk:"description"`
	QuotaLimit     types.Int64  `tfsdk:"quota_limit"`
	ResourceMetric types.String `tfsdk:"resource_metric"`
}

func (model *ObjectStorageQuotaFilterModel) parseQuotas(
	quotas []linodego.ObjectStorageQuota,
) {
	quotaModels := make([]ObjectStorageQuotaModel, len(quotas))

	for i, quota := range quotas {
		var quotaModel ObjectStorageQuotaModel

		quotaModel.QuotaID = types.StringValue(quota.QuotaID)
		quotaModel.QuotaName = types.StringValue(quota.QuotaName)
		quotaModel.EndpointType = types.StringValue(quota.EndpointType)
		quotaModel.S3Endpoint = types.StringValue(quota.S3Endpoint)
		quotaModel.Description = types.StringValue(quota.Description)
		quotaModel.QuotaLimit = types.Int64Value(int64(quota.QuotaLimit))
		quotaModel.ResourceMetric = types.StringValue(quota.ResourceMetric)

		quotaModels[i] = quotaModel
	}

	model.Quotas = quotaModels
}
