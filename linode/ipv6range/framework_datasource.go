package ipv6range

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			"linode_ipv6_range",
			frameworkDatasourceSchema,
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (data *DataSourceModel) parseIPv6Range(
	ctx context.Context, ipv6Range *linodego.IPv6Range,
) diag.Diagnostics {
	data.Range = types.StringValue(ipv6Range.Range)
	data.IsBGP = types.BoolValue(ipv6Range.IsBGP)

	linodes, diag := types.SetValueFrom(ctx, types.Int64Type, ipv6Range.Linodes)
	if diag.HasError() {
		return diag
	}
	data.Linodes = linodes

	data.Prefix = types.Int64Value(int64(ipv6Range.Prefix))
	data.Region = types.StringValue(ipv6Range.Region)

	id, _ := json.Marshal(ipv6Range)

	data.ID = types.StringValue(string(id))

	return nil
}

type DataSourceModel struct {
	Range   types.String `tfsdk:"range"`
	IsBGP   types.Bool   `tfsdk:"is_bgp"`
	Linodes types.Set    `tfsdk:"linodes"`
	Prefix  types.Int64  `tfsdk:"prefix"`
	Region  types.String `tfsdk:"region"`
	ID      types.String `tfsdk:"id"`
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
