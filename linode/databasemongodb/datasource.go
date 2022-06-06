package databasemongodb

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id := d.Get("database_id").(int)

	db, err := client.GetMongoDatabase(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get mysql database: %s", err)
	}

	cert, err := client.GetMongoDatabaseSSL(ctx, int(id))
	if err != nil {
		return diag.Errorf("failed to get cert for the specified mysql database: %s", err)
	}

	creds, err := client.GetMongoDatabaseCredentials(ctx, int(id))
	if err != nil {
		return diag.Errorf("failed to get credentials for mysql database: %s", err)
	}

	d.Set("engine_id", helper.CreateDatabaseEngineSlug(db.Engine, db.Version))
	d.Set("engine", db.Engine)
	d.Set("label", db.Label)
	d.Set("region", db.Region)
	d.Set("type", db.Type)
	d.Set("allow_list", db.AllowList)
	d.Set("cluster_size", db.ClusterSize)
	d.Set("compression_type", db.CompressionType)
	d.Set("encrypted", db.Encrypted)
	d.Set("storage_engine", db.StorageEngine)
	d.Set("ssl_connection", db.SSLConnection)
	d.Set("ca_cert", string(cert.CACertificate))
	d.Set("created", db.Created.Format(time.RFC3339))
	d.Set("host_primary", db.Hosts.Primary)
	d.Set("host_secondary", db.Hosts.Secondary)
	d.Set("peers", db.Peers)
	d.Set("replica_set", db.ReplicaSet)
	d.Set("port", db.Port)
	d.Set("root_password", creds.Password)
	d.Set("status", db.Status)
	d.Set("updated", db.Updated.Format(time.RFC3339))
	d.Set("root_username", creds.Username)
	d.Set("version", db.Version)
	d.Set("updates", []interface{}{helper.FlattenMaintenanceWindow(db.Updates)})

	d.SetId(strconv.Itoa(db.ID))

	return nil
}
