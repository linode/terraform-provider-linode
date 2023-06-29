package ipv6range

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_ipv6_range",
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
	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rangeStrSplit := strings.Split(data.Range.ValueString(), "/")
	rangeStr := rangeStrSplit[0]

	rangeData, err := client.GetIPv6Range(ctx, rangeStr)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get ipv6 range %s :", rangeStr), err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseIPv6Range(ctx, rangeData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
