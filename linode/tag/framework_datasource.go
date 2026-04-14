package tag

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_tag",
				Schema: &frameworkDatasourceSchema,
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
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	label := data.Label.ValueString()
	ctx = tflog.SetField(ctx, "tag_label", label)

	objects, err := d.Meta.Client.ListTaggedObjects(ctx, label, nil)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Tag %q not found", label),
				err.Error(),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to list objects for Tag %q", label),
			err.Error(),
		)
		return
	}

	data.ID = types.StringValue(label)
	data.FlattenTaggedObjects(ctx, objects, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
