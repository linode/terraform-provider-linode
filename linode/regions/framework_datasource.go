package regions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_regions",
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
	tflog.Debug(ctx, "Read data.linode_regions")

	var data RegionFilterModel

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
		ctx, r.Meta.Client, data.Filters, listRegions,
		// There are no API filterable fields so we don't need to provide
		// order and order_by.
		types.StringNull(), types.StringNull())
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	data.parseRegions(helper.AnySliceToTyped[linodego.Region](result))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listRegions(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	tflog.Trace(ctx, "client.ListRegions(...)", map[string]interface{}{
		"filter": filter,
	})

	regions, err := client.ListRegions(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(regions), nil
}
