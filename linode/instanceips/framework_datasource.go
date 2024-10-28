package instanceips

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_instance_ips",
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
	var data InstanceIPDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Meta.Client
	linodeID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Linode ID", err.Error())
		return
	}

	ips, err := client.GetInstanceIPAddresses(ctx, linodeID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get instance IP addresses",
			fmt.Sprintf("Error getting IP addresses for Linode %d: %s", linodeID, err),
		)
		return
	}

	// Populate IPv4 addresses
	ipv4Value, diags := types.ObjectValue(
		map[string]attr.Type{
			"public":   types.ListType{ElemType: types.StringType},
			"private":  types.ListType{ElemType: types.StringType},
			"shared":   types.ListType{ElemType: types.StringType},
			"reserved": types.ListType{ElemType: types.StringType},
			"vpc":      types.ListType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"public":   types.ListValueMust(types.StringType, getIPAddresses(ips.IPv4.Public)),
			"private":  types.ListValueMust(types.StringType, getIPAddresses(ips.IPv4.Private)),
			"shared":   types.ListValueMust(types.StringType, getIPAddresses(ips.IPv4.Shared)),
			"reserved": types.ListValueMust(types.StringType, getIPAddresses(ips.IPv4.Reserved)),
			"vpc":      types.ListValueMust(types.StringType, getVPCIPAddresses(ips.IPv4.VPC)),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.IPv4 = ipv4Value

	// Populate IPv6 addresses
	ipv6Value, diags := types.ObjectValue(
		map[string]attr.Type{
			"link_local": types.StringType,
			"slaac":      types.StringType,
			"global":     types.ListType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"link_local": types.StringValue(ips.IPv6.LinkLocal.Address),
			"slaac":      types.StringValue(ips.IPv6.SLAAC.Address),
			"global":     types.ListValueMust(types.StringType, getIPv6Ranges(ips.IPv6.Global)),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.IPv6 = ipv6Value

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getIPAddresses(ips []*linodego.InstanceIP) []attr.Value {
	addresses := make([]attr.Value, len(ips))
	for i, ip := range ips {
		addresses[i] = types.StringValue(ip.Address)
	}
	return addresses
}

func getVPCIPAddresses(ips []*linodego.VPCIP) []attr.Value {
	addresses := make([]attr.Value, len(ips))
	for i, ip := range ips {
		if ip.Address != nil {
			addresses[i] = types.StringValue(*ip.Address)
		} else {
			addresses[i] = types.StringNull()
		}
	}
	return addresses
}

func getIPv6Ranges(ranges []linodego.IPv6Range) []attr.Value {
	result := make([]attr.Value, len(ranges))
	for i, r := range ranges {
		result[i] = types.StringValue(r.Range)
	}
	return result
}
