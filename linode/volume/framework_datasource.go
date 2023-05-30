package volume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
}

func (r *DataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	r.client = meta.Client
}

func (r *DataSource) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_volume"
}

func (r *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDataSourceSchema
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := r.client

	var data VolumeModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if data.ID.IsNull() || data.ID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("id"),
			"ID Required",
			"ID is required but wasn't configured.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	volume, err := client.GetVolume(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get the volume.", err.Error())
	}
	data.parseComputedAttributes(ctx, volume)
	data.parseNonComputedAttributes(ctx, volume)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
