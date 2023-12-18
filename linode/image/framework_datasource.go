package image

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_image",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := d.Meta.Client
	var data ImageModel

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

	data.ParseImage(image)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
