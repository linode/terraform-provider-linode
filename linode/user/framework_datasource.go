package user

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_user",
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
	client := d.Meta.Client

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
