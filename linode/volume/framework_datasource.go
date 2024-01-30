package volume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_volume",
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
	client := r.Meta.Client

	var data VolumeDataSourceModel

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
	volumeID := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	volume, err := client.GetVolume(ctx, volumeID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get the volume.", err.Error())
	}
	data.ParseComputedAttributes(ctx, volume)
	data.ParseNonComputedAttributes(ctx, volume)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
