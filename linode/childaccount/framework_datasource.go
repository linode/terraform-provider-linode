package childaccount

import (
	"context"

	"github.com/linode/terraform-provider-linode/v2/linode/account"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_child_account",
				Schema: dataSourceSchema(),
			},
		),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_child_account")

	var data account.DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	euuid := data.EUUID.ValueString()

	ctx = tflog.SetField(ctx, "euuid", euuid)

	tflog.Trace(ctx, "client.GetChildAccount(...)")
	childAccount, err := d.Meta.Client.GetChildAccount(ctx, euuid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Child Account",
			err.Error(),
		)
		return
	}

	data.ParseAccount(childAccount)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
