package objquotas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_object_storage_quotas",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+r.Config.Name)

	var data ObjectStorageQuotaFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, d := filterConfig.GenerateID(data.Filters)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}
	data.ID = id

	result, d := filterConfig.GetAndFilter(
		ctx,
		r.Meta.Client,
		data.Filters,
		listObjectStorageQuotas,
		types.StringNull(),
		types.StringNull(),
	)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	data.parseQuotas(helper.AnySliceToTyped[linodego.ObjectStorageQuota](result))
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listObjectStorageQuotas(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	quotas, err := client.ListObjectStorageQuotas(ctx, nil)
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(quotas), nil
}
