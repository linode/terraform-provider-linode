package accountavailabilities

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
				Name:   "linode_account_availabilities",
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
	tflog.Debug(ctx, "Read data.linode_account_availabilities")
	var data AccountAvailabilityFilterModel

	client := r.Client

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
		ctx, client, data.Filters, listAvailabilities,
		types.StringNull(), types.StringNull())
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	if d := data.parseAvailabilities(
		ctx,
		helper.AnySliceToTyped[linodego.AccountAvailability](result),
	); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listAvailabilities(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	tflog.Trace(ctx, "client.ListAccountAvailabilities(...)", map[string]any{
		"filter": filter,
	})
	logins, err := client.ListAccountAvailabilities(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(logins), nil
}
