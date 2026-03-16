package prefixlists

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_prefix_lists",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data PrefixListFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diag := filterConfig.GenerateID(data.Filters)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	data.ID = id

	result, diag := filterConfig.GetAndFilter(
		ctx, d.Meta.Client, data.Filters, listPrefixLists,
		types.StringNull(), types.StringNull())
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	data.parsePrefixLists(
		ctx,
		helper.AnySliceToTyped[linodego.PrefixList](result),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listPrefixLists(
	ctx context.Context,
	client *linodego.Client,
	filter string,
) ([]any, error) {
	tflog.Trace(ctx, "client.ListPrefixLists(...)", map[string]any{
		"filter": filter,
	})
	lists, err := client.ListPrefixLists(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(lists), nil
}
