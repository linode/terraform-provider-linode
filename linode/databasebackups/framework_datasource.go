package databasebackups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_database_backups",
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
	var data DatabaseBackupFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.DatabaseID

	if data.DatabaseType.ValueString() == "mysql" {
		listMySQLBackups := func(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
			databaseID, err := helper.SafeInt64ToInt(
				data.DatabaseID.ValueInt64(),
			)
			if err != nil {
				return nil, err
			}

			backups, err := client.ListMySQLDatabaseBackups(ctx, databaseID, &linodego.ListOptions{
				Filter: filter,
			})
			if err != nil {
				return nil, err
			}

			return helper.TypedSliceToAny(backups), nil
		}

		result, d := filterConfig.GetAndFilter(
			ctx, r.Meta.Client, data.Filters, listMySQLBackups, data.Order, data.OrderBy)
		if d != nil {
			resp.Diagnostics.Append(d)
			return
		}

		data.parseMySQLBackups(helper.AnySliceToTyped[linodego.MySQLDatabaseBackup](result))
	} else if data.DatabaseType.ValueString() == "postgresql" {
		listPostgresSQLBackups := func(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
			databaseID, err := helper.SafeInt64ToInt(data.DatabaseID.ValueInt64())
			if err != nil {
				return nil, err
			}

			backups, err := client.ListPostgresDatabaseBackups(
				ctx,
				databaseID,
				&linodego.ListOptions{
					Filter: filter,
				},
			)
			if err != nil {
				return nil, err
			}

			return helper.TypedSliceToAny(backups), nil
		}

		result, d := filterConfig.GetAndFilter(
			ctx, r.Meta.Client, data.Filters, listPostgresSQLBackups, data.Order, data.OrderBy)
		if d != nil {
			resp.Diagnostics.Append(d)
			return
		}

		data.parsePostgresSQLBackups(helper.AnySliceToTyped[linodego.PostgresDatabaseBackup](result))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
