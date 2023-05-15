package instancenetworking

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
}

func (data *DataSourceModel) parseInstanceIPAddressResponse(
	ctx context.Context, ip *linodego.InstanceIPAddressResponse,
) diag.Diagnostics {
	ipv4, diags := flattenIPv4(ctx, ip.IPv4)
	if diags.HasError() {
		return diags
	}

	data.IPV4 = *ipv4

	ipv6, diags := flattenIPv6(ctx, ip.IPv6)
	if diags.HasError() {
		return diags
	}

	data.IPV6 = *ipv6

	id, err := json.Marshal(ip)
	if err != nil {
		diags.AddError("Error marshalling json: %s", err.Error())
		return diags
	}

	data.ID = types.StringValue(string(id))

	return nil
}

func (d *DataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client = meta.Client
}

type DataSourceModel struct {
	LinodeID types.Int64  `tfsdk:"linode_id"`
	IPV4     types.Object `tfsdk:"ipv4"`
	IPV6     types.Object `tfsdk:"ipv6"`
	ID       types.String `tfsdk:"id"`
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_instance_networking"
}

func (d *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDatasourceSchema
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := d.client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	netInfo, err := client.GetInstanceIPAddresses(ctx, int(data.LinodeID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Instance Networking Information: ", err.Error(),
		)
		return
	}

	data.parseInstanceIPAddressResponse(ctx, netInfo)

	resp.Diagnostics.Append(data.parseInstanceIPAddressResponse(ctx, netInfo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
