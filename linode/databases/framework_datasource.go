package databases

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_databases",
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
	var data DatabaseFilterModel

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
		ctx, r.Meta.Client, data.Filters, listDatabases, data.Order, data.OrderBy)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	data.parseDatabases(helper.AnySliceToTyped[linodego.Database](result))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listDatabases(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	databases, err := client.ListDatabases(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(databases), nil
}
