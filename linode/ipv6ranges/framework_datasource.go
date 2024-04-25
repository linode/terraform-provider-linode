package ipv6ranges

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_ipv6_ranges",
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
	tflog.Debug(ctx, "Read data.linode_ipv6_ranges")

	client := r.Client

	var data IPv6RangeFilterModel

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
		ctx, client, data.Filters, listRanges,
		types.StringNull(), types.StringNull())
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	resp.Diagnostics.Append(data.parseRanges(helper.AnySliceToTyped[linodego.IPv6Range](result))...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listRanges(
	ctx context.Context,
	client *linodego.Client,
	filter string,
) ([]any, error) {
	ctx = tflog.SetField(ctx, "filter", filter)
	tflog.Trace(ctx, "client.ListIPv6Ranges(...)")

	ranges, err := client.ListIPv6Ranges(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(ranges), nil
}
