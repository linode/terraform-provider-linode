package accountsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
<<<<<<< HEAD
	"github.com/linode/linodego"
=======
>>>>>>> 953248a85dd3a5b6c5dde77146267067f81c9ec3
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
<<<<<<< HEAD
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
}

func (r *DataSource) Configure(
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

	r.client = meta.Client
}

func (r *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_account_settings"
}

func (r *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDataSourceSchema
=======
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_account_settings",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
>>>>>>> 953248a85dd3a5b6c5dde77146267067f81c9ec3
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
<<<<<<< HEAD
	client := r.client
=======
	client := r.Meta.Client
>>>>>>> 953248a85dd3a5b6c5dde77146267067f81c9ec3

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
