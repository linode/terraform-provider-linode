package objquota

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

type DataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	QuotaID        types.String `tfsdk:"quota_id"`
	QuotaName      types.String `tfsdk:"quota_name"`
	EndpointType   types.String `tfsdk:"endpoint_type"`
	S3Endpoint     types.String `tfsdk:"s3_endpoint"`
	Description    types.String `tfsdk:"description"`
	QuotaLimit     types.Int64  `tfsdk:"quota_limit"`
	ResourceMetric types.String `tfsdk:"resource_metric"`
	QuotaType      types.String `tfsdk:"quota_type"`
	HasUsage       types.Bool   `tfsdk:"has_usage"`
	QuotaUsage     types.Object `tfsdk:"quota_usage"`
}

type QuotaUsageModel struct {
	QuotaLimit types.Int64 `tfsdk:"quota_limit"`
	Usage      types.Int64 `tfsdk:"usage"`
}

func (data *DataSourceModel) parseObjectStorageQuota(
	ctx context.Context,
	quota *linodego.ObjectStorageQuota,
	usage *linodego.ObjectStorageQuotaUsage,
) diag.Diagnostics {
	data.ID = types.StringValue(quota.QuotaID)
	data.QuotaID = types.StringValue(quota.QuotaID)
	data.QuotaName = types.StringValue(quota.QuotaName)
	data.EndpointType = types.StringValue(quota.EndpointType)
	data.S3Endpoint = types.StringValue(quota.S3Endpoint)
	data.Description = types.StringValue(quota.Description)
	data.QuotaLimit = types.Int64Value(int64(quota.QuotaLimit))
	data.ResourceMetric = types.StringValue(quota.ResourceMetric)
	data.QuotaType = types.StringValue(quota.QuotaType)
	data.HasUsage = types.BoolValue(quota.HasUsage)

	if !quota.HasUsage || usage == nil {
		data.QuotaUsage = types.ObjectNull(quotaUsageAttributes)
		return nil
	}

	quotaUsage, diag := types.ObjectValueFrom(
		ctx,
		quotaUsageAttributes,
		QuotaUsageModel{
			QuotaLimit: types.Int64Value(int64(usage.QuotaLimit)),
			Usage:      helper.IntPointerValueWithDefault(usage.Usage),
		},
	)
	if diag != nil {
		return diag
	}

	data.QuotaUsage = quotaUsage

	return nil
}
