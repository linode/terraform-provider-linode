package databases

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
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listDatabases, flattenDatabase)
	if err != nil {
		return nil
	}

	results = filterConfig.FilterLatest(d, results)

	d.Set("databases", results)

	return nil
}

func listDatabases(
	ctx context.Context, d *schema.ResourceData,
	client *linodego.Client, options *linodego.ListOptions) ([]interface{}, error) {
	// TODO: return a list of engines

	return nil, nil
}

func flattenDatabase(data interface{}) map[string]interface{} {
	// TODO: Flatten the engine info into a map
	return nil
}
