package kernels

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_kernels",
				Schema: &frameworkDatasourceSchema,
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
	tflog.Debug(ctx, "Read data.linode_kernels")

	var data KernelFilterModel

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
		ctx, d.Client, data.Filters, listKernels,
		data.Order, data.OrderBy)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	data.parseKernels(ctx, helper.AnySliceToTyped[linodego.LinodeKernel](result))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listKernels(
	ctx context.Context,
	client *linodego.Client,
	filter string,
) ([]any, error) {
	tflog.Trace(ctx, "client.ListKernels(...)", map[string]interface{}{
		"filter": filter,
	})

	kernels, err := client.ListKernels(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(kernels), nil
}
