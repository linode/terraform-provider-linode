package images

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_images",
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
	tflog.Debug(ctx, "Read data.linode_images")

	var data ImageFilterModel

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
		ctx, d.Meta.Client, data.Filters, listImages,
		data.Order, data.OrderBy)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	if data.Latest.ValueBool() {
		result, diag = filterConfig.GetLatestCreated(result, "Created")
		if diag != nil {
			resp.Diagnostics.Append(diag)
			return
		}
	}

	data.parseImages(helper.AnySliceToTyped[linodego.Image](result))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listImages(
	ctx context.Context,
	client *linodego.Client,
	filter string,
) ([]any, error) {
	tflog.Trace(ctx, "client.ListImages", map[string]any{
		"filter": filter,
	})

	images, err := client.ListImages(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(images), nil
}
