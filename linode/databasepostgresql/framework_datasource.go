package databasepostgresql

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
	resp.TypeName = "linode_database_postgresql"
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

	var id int

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		id = int(data.ID.ValueInt64())
	} else {
		id = int(data.DatabaseID.ValueInt64())
	}

	if id == 0 {
		resp.Diagnostics.AddError(
			"ID not provided properly.", "",
		)
		return
	}

	db, err := client.GetPostgresDatabase(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get PostgreSQL database with ID %d: ", id), err.Error(),
		)
		return
	}

	cert, err := client.GetPostgresDatabaseSSL(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get cert for the specified PostgreSQL database with ID %d: ", id), err.Error(),
		)
		return
	}

	cred, err := client.GetPostgresDatabaseCredentials(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get credentials for PostgreSQL database with ID %d: ", id), err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parsePostgresDatabase(ctx, db)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.parsePostgresDatabaseSSL(cert)
	data.parsePostgresDatabaseCredentials(cred)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
