package databaseengines

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

func readDataSource(ctx context.Context, d *schema.ResourceData,
	meta interface{},
) diag.Diagnostics {
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listEngines, flattenEngine)
	if err != nil {
		return nil
	}

	results, err = filterConfig.FilterLatestVersion(d, results)
	if err != nil {
		return diag.Errorf("failed to filter versions: %s", err)
	}

	d.Set("engines", results)

	return nil
}

func listEngines(
	ctx context.Context, d *schema.ResourceData, client *linodego.Client,
	options *linodego.ListOptions,
) ([]interface{}, error) {
	engines, err := client.ListDatabaseEngines(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(engines))

	for i, v := range engines {
		result[i] = v
	}

	return result, nil
}

func flattenEngine(data interface{}) map[string]interface{} {
	engine := data.(linodego.DatabaseEngine)

	result := make(map[string]interface{})

	result["id"] = engine.ID
	result["engine"] = engine.Engine
	result["version"] = engine.Version

	return result
}
