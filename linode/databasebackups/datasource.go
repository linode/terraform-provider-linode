package databasebackups

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listBackups, FlattenBackup)
	if err != nil {
		return nil
	}

	results = filterConfig.FilterLatest(d, results)

	d.Set("backups", results)

	return nil
}

func listBackups(
	ctx context.Context, d *schema.ResourceData, client *linodego.Client,
	options *linodego.ListOptions,
) ([]interface{}, error) {
	dbID := d.Get("database_id").(int)
	dbType := d.Get("database_type").(string)

	switch dbType {
	case "mysql":
		return helper.ListResultToInterface(
			client.ListMySQLDatabaseBackups(ctx, dbID, nil))
	case "mongodb":
		return helper.ListResultToInterface(
			client.ListMongoDatabaseBackups(ctx, dbID, nil))
	case "postgresql":
		return helper.ListResultToInterface(
			client.ListPostgresDatabaseBackups(ctx, dbID, nil))
	}

	return nil, fmt.Errorf("invalid database type: %s", dbType)
}

func flattenMySQLBackup(backup linodego.MySQLDatabaseBackup) map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = backup.ID
	result["label"] = backup.Label
	result["type"] = backup.Type

	if backup.Created != nil {
		result["created"] = backup.Created.Format(time.RFC3339)
	}

	return result
}

func flattenMongoBackup(backup linodego.MongoDatabaseBackup) map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = backup.ID
	result["label"] = backup.Label
	result["type"] = backup.Type

	if backup.Created != nil {
		result["created"] = backup.Created.Format(time.RFC3339)
	}

	return result
}

func flattenPostgresBackup(backup linodego.PostgresDatabaseBackup) map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = backup.ID
	result["label"] = backup.Label
	result["type"] = backup.Type

	if backup.Created != nil {
		result["created"] = backup.Created.Format(time.RFC3339)
	}

	return result
}

func FlattenBackup(data interface{}) map[string]interface{} {
	switch data.(type) {
	case linodego.MySQLDatabaseBackup:
		return flattenMySQLBackup(data.(linodego.MySQLDatabaseBackup))
	case linodego.MongoDatabaseBackup:
		return flattenMongoBackup(data.(linodego.MongoDatabaseBackup))
	case linodego.PostgresDatabaseBackup:
		return flattenPostgresBackup(data.(linodego.PostgresDatabaseBackup))
	}

	log.Printf("[WARN] Could not flatten backup due to invalid type")

	return nil
}
