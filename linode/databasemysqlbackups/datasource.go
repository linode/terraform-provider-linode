package databasemysqlbackups

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated.",
		Schema:             dataSourceSchema,
		ReadContext:        readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listBackups, flattenBackup)
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
) ([]any, error) {
	dbID := d.Get("database_id").(int)

	backups, err := client.ListMySQLDatabaseBackups(ctx, dbID, options)
	if err != nil {
		return nil, err
	}

	result := make([]any, len(backups))

	for i, v := range backups {
		result[i] = v
	}

	return result, nil
}

func flattenBackup(data any) map[string]any {
	backup := data.(linodego.MySQLDatabaseBackup)

	result := make(map[string]any)

	result["id"] = backup.ID
	result["label"] = backup.Label
	result["type"] = backup.Type
	result["created"] = backup.Created.Format(time.RFC3339)

	return result
}
