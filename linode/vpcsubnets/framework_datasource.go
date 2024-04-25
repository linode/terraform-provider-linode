package vpcsubnets

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_vpc_subnets",
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
	tflog.Debug(ctx, "Read data.linode_vpc_subnets")
	var data VPCSubnetFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "vpc_id", data.VPCId)

	id, d := filterConfig.GenerateID(data.Filters)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}
	data.ID = id

	result, d := filterConfig.GetAndFilter(
		ctx, r.Client, data.Filters, data.ListVPCSubnets,
		// There are no API filterable fields so we don't need to provide
		// order and order_by.
		types.StringNull(), types.StringNull())
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resp.Diagnostics.Append(data.FlattenSubnets(ctx, helper.AnySliceToTyped[linodego.VPCSubnet](result), false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (data *VPCSubnetFilterModel) ListVPCSubnets(
	ctx context.Context,
	client *linodego.Client,
	filter string,
) ([]any, error) {
	var diags diag.Diagnostics
	vpcId := helper.FrameworkSafeInt64ToInt(data.VPCId.ValueInt64(), &diags)
	if diags.HasError() {
		for _, err := range diags.Errors() {
			return nil, errors.New(err.Detail())
		}
	}

	tflog.Trace(ctx, "client.ListVPCSubnets(...)", map[string]any{
		"filter": filter,
	})
	vpcs, err := client.ListVPCSubnets(ctx, vpcId, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(vpcs), nil
}
