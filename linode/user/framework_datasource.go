package user

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
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
	resp.TypeName = "linode_user"
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

	user, err := client.GetUser(ctx, data.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Linode User with username %s was not found", data.Username.ValueString()), err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseUser(ctx, user)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if user.Restricted {
		grants, err := client.GetUserGrants(ctx, data.Username.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to get User Grants (%s): ", data.Username.ValueString()), err.Error(),
			)
			return
		}
		resp.Diagnostics.Append(data.ParseUserGrants(ctx, grants)...)
	} else {
		data.ParseNonUserGrants()
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
