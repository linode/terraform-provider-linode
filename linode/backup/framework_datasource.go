package backup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type DataSource struct {
	client *linodego.Client
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
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

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_instance_backups"
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
	var data DataSourceModel
	client := d.client

	var linodeId types.Int64
	diag := req.Config.GetAttribute(ctx, path.Root("linode_id"), &linodeId)

	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	backups, err := client.GetInstanceBackups(ctx, int(linodeId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get backups",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseBackups(ctx, backups, linodeId)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
