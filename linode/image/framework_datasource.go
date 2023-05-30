package image

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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

type DataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
	Created     types.String `tfsdk:"created"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Deprecated  types.Bool   `tfsdk:"deprecated"`
	IsPublic    types.Bool   `tfsdk:"is_public"`
	Size        types.Int64  `tfsdk:"size"`
	Status      types.String `tfsdk:"status"`
	Type        types.String `tfsdk:"type"`
	Expiry      types.String `tfsdk:"expiry"`
	Vendor      types.String `tfsdk:"vendor"`
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
	resp.TypeName = "linode_image"
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

	if data.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Image ID is required.",
			"Image ID can't be empty.",
		)
		return
	}

	image, err := client.GetImage(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error listing images: %s", data.ID.ValueString()),
			err.Error(),
		)
		return
	}

	if image == nil {
		resp.Diagnostics.AddError(
			"Image not found.",
			fmt.Sprintf("Image %s was not found", data.ID.ValueString()),
		)
		return
	}

	data.parseImage(ctx, image)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (data *DataSourceModel) parseImage(
	ctx context.Context,
	image *linodego.Image,
) {
	data.ID = types.StringValue(image.ID)
	data.Label = types.StringValue(image.Label)
	data.Description = types.StringValue(image.Description)
	if image.Created != nil {
		data.Created = types.StringValue(image.Created.Format(time.RFC3339))
	}
	if image.Expiry != nil {
		data.Expiry = types.StringValue(image.Expiry.Format(time.RFC3339))
	}
	data.CreatedBy = types.StringValue(image.CreatedBy)
	data.Deprecated = types.BoolValue(image.Deprecated)
	data.IsPublic = types.BoolValue(image.IsPublic)
	data.Size = types.Int64Value(int64(image.Size))
	data.Status = types.StringValue(string(image.Status))
	data.Type = types.StringValue(string(image.Type))
	data.Vendor = types.StringValue(string(image.Vendor))
}
