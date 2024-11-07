package networkingips

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
				Name:   "linode_networking_ips",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (data *DataSourceModel) parseIP(ip *linodego.InstanceIP) {
	data.Address = types.StringValue(ip.Address)
	data.Gateway = types.StringValue(ip.Gateway)
	data.SubnetMask = types.StringValue(ip.SubnetMask)
	data.Prefix = types.Int64Value(int64(ip.Prefix))
	data.Type = types.StringValue(string(ip.Type))
	data.Public = types.BoolValue(ip.Public)
	data.RDNS = types.StringValue(ip.RDNS)
	data.LinodeID = types.Int64Value(int64(ip.LinodeID))
	data.Region = types.StringValue(ip.Region)

	id, _ := json.Marshal(ip)
	data.Reserved = types.BoolValue(ip.Reserved)
	data.ID = types.StringValue(string(id))
}

type DataSourceModel struct {
	Address        types.String `tfsdk:"address"`
	Gateway        types.String `tfsdk:"gateway"`
	SubnetMask     types.String `tfsdk:"subnet_mask"`
	Prefix         types.Int64  `tfsdk:"prefix"`
	Type           types.String `tfsdk:"type"`
	Public         types.Bool   `tfsdk:"public"`
	RDNS           types.String `tfsdk:"rdns"`
	LinodeID       types.Int64  `tfsdk:"linode_id"`
	Region         types.String `tfsdk:"region"`
	ID             types.String `tfsdk:"id"`
	Reserved       types.Bool   `tfsdk:"reserved"`
	IPAddresses    types.List   `tfsdk:"ip_addresses"`
	FilterReserved types.Bool   `tfsdk:"filter_reserved"`
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_networking_ip")

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Address.IsNull() {

		// List all IP addresses with filter on reservation status

		filter, err := buildFilter(data)
		if err != nil {
			resp.Diagnostics.AddError("Unable to build filter", err.Error())
			return
		}

		tflog.Debug(ctx, "Generated filter", map[string]interface{}{
			"filter": filter,
		})

		opts := &linodego.ListOptions{Filter: filter}
		ips, err := d.Meta.Client.ListIPAddresses(ctx, opts)
		if err != nil {
			resp.Diagnostics.AddError("Unable to list IP Addresses", err.Error())
			return
		}

		ipList := make([]attr.Value, len(ips))
		for i, ip := range ips {
			ipObj := map[string]attr.Value{
				"address":     types.StringValue(ip.Address),
				"region":      types.StringValue(ip.Region),
				"gateway":     types.StringValue(ip.Gateway),
				"subnet_mask": types.StringValue(ip.SubnetMask),
				"prefix":      types.Int64Value(int64(ip.Prefix)),
				"type":        types.StringValue(string(ip.Type)),
				"public":      types.BoolValue(ip.Public),
				"rdns":        types.StringValue(ip.RDNS),
				"linode_id":   types.Int64Value(int64(ip.LinodeID)),
				"reserved":    types.BoolValue(ip.Reserved),
			}
			ipList[i] = types.ObjectValueMust(updatedIPObjectType.AttrTypes, ipObj)
		}

		var diags diag.Diagnostics
		data.IPAddresses, diags = types.ListValue(updatedIPObjectType, ipList)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildFilter(data DataSourceModel) (string, error) {
	filters := make(map[string]string)

	if !data.FilterReserved.IsNull() {
		filters["reserved"] = fmt.Sprintf("%t", data.FilterReserved.ValueBool())
	}

	if len(filters) == 0 {
		return "", nil
	}

	jsonFilter, err := json.Marshal(filters)
	if err != nil {
		return "", fmt.Errorf("error creating filter: %v", err)
	}

	return string(jsonFilter), nil
}
