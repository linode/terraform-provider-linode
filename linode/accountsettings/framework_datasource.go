package accountsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.BaseDataSource{
			TypeName:     "linode_account_settings",
			SchemaObject: frameworkDataSourceSchema,
		},
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
	client := r.GetMeta().Client

	var data AccountSettingsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := client.GetAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Linode Account",
			err.Error(),
		)
		return
	}

	settings, err := client.GetAccountSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Linode Account Settings",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseAccountSettings(
		ctx,
		account.Email,
		settings,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
