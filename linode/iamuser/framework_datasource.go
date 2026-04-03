package iamuser

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_iam_user",
				Schema: &frameworkSchema,
			},
		),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	client := d.Meta.Client
	var data IAMUserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "username", data.Username.ValueString())

	perms, err := client.GetUserRolePermissions(ctx, data.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error listing permissions: %s", data.Username.ValueString()),
			err.Error(),
		)
		return
	}

	if perms == nil {
		resp.Diagnostics.AddError(
			"Permissions not found.",
			fmt.Sprintf("Permissions for %s was not found", data.Username.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseIAMUserModel(ctx, perms)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
