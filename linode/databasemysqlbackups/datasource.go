package databasemysqlbackups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listBackups, flattenBackup)
	if err != nil {
		return nil
	}

	results = filterConfig.FilterLatest(d, results)

	d.Set("backups", results)

	return nil
}

func listBackups(
	ctx context.Context, client *linodego.Client, options *linodego.ListOptions) ([]interface{}, error) {
	// TODO: return a list of backups

	return nil, nil
}

func flattenBackup(data interface{}) map[string]interface{} {
	// TODO: Flatten the backup info a map
	return nil
}
