package databasebackups

import (
	"context"

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
	resp.TypeName = "linode_database_backups"
}

func (r *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDataSourceSchema
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
			backups, err := client.ListMySQLDatabaseBackups(ctx, int(data.DatabaseID.ValueInt64()), &linodego.ListOptions{
				Filter: filter,
			})
			if err != nil {
				return nil, err
			}

			return helper.TypedSliceToAny(backups), nil
		}

		result, d := filterConfig.GetAndFilter(
			ctx, r.client, data.Filters, listMySQLBackups, data.Order, data.OrderBy)
		if d != nil {
			resp.Diagnostics.Append(d)
			return
		}

		data.parseMySQLBackups(helper.AnySliceToTyped[linodego.MySQLDatabaseBackup](result))
	} else if data.DatabaseType.ValueString() == "postgresql" {
		listPostgresSQLBackups := func(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
			backups, err := client.ListPostgresDatabaseBackups(ctx, int(data.DatabaseID.ValueInt64()), &linodego.ListOptions{
				Filter: filter,
			})
			if err != nil {
				return nil, err
			}

			return helper.TypedSliceToAny(backups), nil
		}

		result, d := filterConfig.GetAndFilter(
			ctx, r.client, data.Filters, listPostgresSQLBackups, data.Order, data.OrderBy)
		if d != nil {
			resp.Diagnostics.Append(d)
			return
		}

		data.parsePostgresSQLBackups(helper.AnySliceToTyped[linodego.PostgresDatabaseBackup](result))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
