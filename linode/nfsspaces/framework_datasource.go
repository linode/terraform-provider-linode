package nfsspaces

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_nfs_spaces",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data NFSSpaceFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, d2 := filterConfig.GenerateID(data.Filters)
	if d2 != nil {
		resp.Diagnostics.Append(d2)
		return
	}
	data.ID = id

	result, d2 := filterConfig.GetAndFilter(
		ctx, d.Meta.Client, data.Filters, ListNFSSpaces,
		// There are no API filterable fields so we don't need to provide
		// order and order_by.
		types.StringNull(), types.StringNull())
	if d2 != nil {
		resp.Diagnostics.Append(d2)
		return
	}

	resp.Diagnostics.Append(data.FlattenNFSSpaces(ctx, helper.AnySliceToTyped[linodego.NFSSpace](result), false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func ListNFSSpaces(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	tflog.Trace(ctx, "client.ListNFSSpaces(...)", map[string]any{
		"filter": filter,
	})
	spaces, err := client.ListNFSSpaces(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(spaces), nil
}
