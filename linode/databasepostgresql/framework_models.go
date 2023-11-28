package databasepostgresql

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSourceModel struct {
	DatabaseID            types.Int64  `tfsdk:"database_id"`
	EngineID              types.String `tfsdk:"engine_id"`
	Label                 types.String `tfsdk:"label"`
	Region                types.String `tfsdk:"region"`
	Type                  types.String `tfsdk:"type"`
	AllowList             types.Set    `tfsdk:"allow_list"`
	ClusterSize           types.Int64  `tfsdk:"cluster_size"`
	Encrypted             types.Bool   `tfsdk:"encrypted"`
	ReplicationType       types.String `tfsdk:"replication_type"`
	ReplicationCommitType types.String `tfsdk:"replication_commit_type"`
	SSLConnection         types.Bool   `tfsdk:"ssl_connection"`
	CACert                types.String `tfsdk:"ca_cert"`
	Created               types.String `tfsdk:"created"`
	Engine                types.String `tfsdk:"engine"`
	HostPrimary           types.String `tfsdk:"host_primary"`
	HostSecondary         types.String `tfsdk:"host_secondary"`
	Port                  types.Int64  `tfsdk:"port"`
	RootPassword          types.String `tfsdk:"root_password"`
	Status                types.String `tfsdk:"status"`
	Updated               types.String `tfsdk:"updated"`
	Updates               types.List   `tfsdk:"updates"`
	RootUsername          types.String `tfsdk:"root_username"`
	Version               types.String `tfsdk:"version"`
	ID                    types.Int64  `tfsdk:"id"`
}

func (data *DataSourceModel) parsePostgresDatabase(
	ctx context.Context, db *linodego.PostgresDatabase,
) diag.Diagnostics {
	data.DatabaseID = types.Int64Value(int64(db.ID))
	data.Status = types.StringValue(string(db.Status))
	data.Label = types.StringValue(string(db.Label))
	data.HostPrimary = types.StringValue(string(db.Hosts.Primary))
	data.HostSecondary = types.StringValue(string(db.Hosts.Secondary))
	data.Region = types.StringValue(string(db.Region))
	data.Type = types.StringValue(string(db.Type))
	data.Engine = types.StringValue(string(db.Engine))
	data.Port = types.Int64Value(int64(db.Port))
	data.Version = types.StringValue(string(db.Version))
	data.ClusterSize = types.Int64Value(int64(db.ClusterSize))
	data.ReplicationType = types.StringValue(string(db.ReplicationType))
	data.ReplicationCommitType = types.StringValue(string(db.ReplicationCommitType))
	data.SSLConnection = types.BoolValue(db.SSLConnection)
	data.Encrypted = types.BoolValue(db.Encrypted)

	allowList, diags := types.SetValueFrom(ctx, types.StringType, db.AllowList)
	if diags.HasError() {
		return diags
	}
	data.AllowList = allowList

	data.Created = types.StringValue(db.Created.Format(time.RFC3339))
	data.Updated = types.StringValue(db.Updated.Format(time.RFC3339))

	data.EngineID = types.StringValue(helper.CreateDatabaseEngineSlug(db.Engine, db.Version))

	updates, diags := helper.FlattenDatabaseMaintenanceWindow(ctx, db.Updates)
	if diags.HasError() {
		return diags
	}

	data.Updates = *updates

	data.ID = types.Int64Value(int64(db.ID))

	return nil
}

func (data *DataSourceModel) parsePostgresDatabaseSSL(db *linodego.PostgresDatabaseSSL) {
	data.CACert = types.StringValue(string(db.CACertificate))
}

func (data *DataSourceModel) parsePostgresDatabaseCredentials(db *linodego.PostgresDatabaseCredential) {
	data.RootUsername = types.StringValue(db.Username)
	data.RootPassword = types.StringValue(db.Password)
}
