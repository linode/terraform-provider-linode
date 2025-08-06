package vpcips

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_vpc_ips",
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

	var data Model

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

	var result []any

	// Check if VPC ID is provided
	vpcID := helper.FrameworkSafeInt64ToInt(data.VPCID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var listFunc func(ctx context.Context, client *linodego.Client, filter string, vpcID int) ([]any, error)
	var listAllFunc func(context.Context, *linodego.Client, string) ([]any, error)

	// Ensure we're listing calling the correct endpoint for the
	// user-specified IP type.
	if data.IPv6.ValueBool() {
		listFunc = listVPCIPv6s
		listAllFunc = listAllVPCIPv6s
	} else {
		listFunc = listVPCIPv4s
		listAllFunc = listAllVPCIPv4s
	}

	if !data.VPCID.IsNull() {
		tflog.Debug(ctx, "Filtering IPs for specific VPC", map[string]interface{}{
			"vpc_id": vpcID,
		})

		result, d = filterConfig.GetAndFilter(
			ctx, r.Meta.Client, data.Filters,
			func(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
				return listFunc(ctx, client, filter, vpcID)
			},
			types.StringNull(), types.StringNull(),
		)
	} else {
		tflog.Debug(ctx, "Filtering all IPs in the account")

		result, d = filterConfig.GetAndFilter(
			ctx, r.Meta.Client, data.Filters,
			listAllFunc,
			types.StringNull(), types.StringNull(),
		)
	}

	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	data.FlattenVPCIPs(ctx, helper.AnySliceToTyped[linodego.VPCIP](result), false)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listAllVPCIPv4s(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	tflog.Trace(ctx, "client.ListAllVPCIPAddresses(...)", map[string]any{
		"filter": filter,
	})
	vpcIps, err := client.ListAllVPCIPAddresses(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(vpcIps), nil
}

func listVPCIPv4s(ctx context.Context, client *linodego.Client, filter string, vpcID int) ([]any, error) {
	tflog.Trace(ctx, "client.ListVPCIPAddresses(...)", map[string]any{
		"filter": filter,
	})
	vpcIps, err := client.ListVPCIPAddresses(ctx, vpcID, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(vpcIps), nil
}

func listAllVPCIPv6s(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	tflog.Trace(ctx, "client.ListAllVPCIPv6Addresses(...)", map[string]any{
		"filter": filter,
	})
	vpcIps, err := client.ListAllVPCIPv6Addresses(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(vpcIps), nil
}

func listVPCIPv6s(ctx context.Context, client *linodego.Client, filter string, vpcID int) ([]any, error) {
	tflog.Trace(ctx, "client.ListVPCIPv6Addresses(...)", map[string]any{
		"filter": filter,
	})
	vpcIps, err := client.ListVPCIPv6Addresses(ctx, vpcID, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(vpcIps), nil
}
