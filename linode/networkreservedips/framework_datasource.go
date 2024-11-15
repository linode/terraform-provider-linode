package networkreservedips

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_reserved_ips",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data.linode_reserved_ips")

	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := d.Meta.Client

	// List all reserved IPs
	ips, err := client.ListReservedIPAddresses(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list Reserved IP Addresses",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Retrieved Reserved IP Addresses", map[string]interface{}{
		"ips": ips,
	})

	// Filter IPs by region if specified
	var filteredIPs []linodego.InstanceIP
	var regionFilter types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("region"), &regionFilter)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !regionFilter.IsNull() {
		for _, ip := range ips {
			if ip.Region == regionFilter.ValueString() {
				filteredIPs = append(filteredIPs, ip)
			}
		}
	} else {
		filteredIPs = ips // No filtering if region is not specified
	}

	reservedIPs := make([]ReservedIPObject, len(filteredIPs))
	for i, ip := range filteredIPs {
		var vpcNAT1To1List types.List
		if ip.VPCNAT1To1 == nil {
			vpcNAT1To1List = types.ListNull(instancenetworking.VPCNAT1To1Type)
		} else {
			vpcNAT1To1, _ := instancenetworking.FlattenIPVPCNAT1To1(ip.VPCNAT1To1)
			vpcNAT1To1List, _ = types.ListValue(
				instancenetworking.VPCNAT1To1Type,
				[]attr.Value{vpcNAT1To1},
			)
		}

		reservedIPs[i] = ReservedIPObject{
			ID:           types.StringValue(ip.Address),
			Address:      types.StringValue(ip.Address),
			Region:       types.StringValue(ip.Region),
			Gateway:      types.StringValue(ip.Gateway),
			SubnetMask:   types.StringValue(ip.SubnetMask),
			Prefix:       types.Int64Value(int64(ip.Prefix)),
			Type:         types.StringValue(string(ip.Type)),
			Public:       types.BoolValue(ip.Public),
			RDNS:         types.StringValue(ip.RDNS),
			LinodeID:     types.Int64Value(int64(ip.LinodeID)),
			Reserved:     types.BoolValue(ip.Reserved),
			IPVPCNAT1To1: vpcNAT1To1List,
		}
	}

	reservedIPsValue, diags := types.ListValueFrom(ctx, reservedIPObjectType, reservedIPs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ReservedIPs = reservedIPsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
