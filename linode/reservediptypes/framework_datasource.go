package reservediptypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_reserved_ip_types",
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

	var data reservedIPTypeFilterModel

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
		ctx, r.Meta.Client, data.Filters, listReservedIPTypes, data.Order, data.OrderBy)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	diags := data.parseReservedIPTypes(ctx, helper.AnySliceToTyped[linodego.ReservedIPType](result))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listReservedIPTypes(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	tflog.Debug(ctx, "Listing Reserved IP types", map[string]any{
		"filter_header": filter,
	})

	types, err := client.ListReservedIPTypes(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(types), nil
}
