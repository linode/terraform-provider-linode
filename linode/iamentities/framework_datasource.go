package iamentities

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
				Name:   "linode_iam_entities",
				Schema: &frameworkSchema,
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

	client := d.Meta.Client
	var data ListEntitiesModel
	diags := resp.Diagnostics

	diags.Append(req.Config.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	id, diag := filterConfig.GenerateID(data.Filters)
	if diag != nil {
		diags.Append(diag)
		return
	}
	data.ID = id

	result, diag := filterConfig.GetAndFilter(
		ctx, client, data.Filters, listEntities, data.Order, data.OrderBy,
	)
	if diag != nil {
		diags.Append(diag)
		return
	}

	data.parseEntities(helper.AnySliceToTyped[linodego.LinodeEntity](result))
	diags.Append(resp.State.Set(ctx, &data)...)
}

func listEntities(
	ctx context.Context, client *linodego.Client, filter string,
) ([]any, error) {
	entities, err := client.ListEntities(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}
	return helper.TypedSliceToAny(entities), nil
}
