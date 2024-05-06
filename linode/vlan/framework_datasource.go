package vlan

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				Name:          "linode_vlans",
				Schema:        &frameworkDatasourceSchema,
				IsEarlyAccess: true,
			},
		),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_vlans")

	var data VLANsFilterModel

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

	results, diag := filterConfig.GetAndFilter(
		ctx, d.Meta.Client,
		data.Filters,
		listVLANs,
		data.Order,
		data.OrderBy,
	)

	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	data.parseVLANs(ctx, helper.AnySliceToTyped[linodego.VLAN](results))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listVLANs(
	ctx context.Context,
	client *linodego.Client,
	filter string,
) ([]any, error) {
	tflog.Trace(ctx, "client.ListVLANs(...)", map[string]any{
		"filter": filter,
	})

	vlans, err := client.ListVLANs(
		ctx,
		&linodego.ListOptions{Filter: filter},
	)
	if err != nil {
		return nil, err
	}
	return helper.TypedSliceToAny(vlans), nil
}
