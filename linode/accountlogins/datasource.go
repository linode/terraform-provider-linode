package accountlogins

import (
	"context"
	"time"

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
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listLogins, flattenLogins)
	if err != nil {
		return nil
	}

	d.Set("logins", results)

	return nil
}

func listLogins(
	ctx context.Context, d *schema.ResourceData, client *linodego.Client,
	options *linodego.ListOptions,
) ([]interface{}, error) {
	types, err := client.ListLogins(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(types))

	for i, v := range types {
		result[i] = v
	}

	return result, nil
}

func flattenLogins(data interface{}) map[string]interface{} {
	t := data.(linodego.Login)

	result := make(map[string]interface{})

	result["id"] = t.ID
	result["datetime"] = t.Datetime.Format(time.RFC3339)
	result["ip"] = t.IP
	result["restricted"] = t.Restricted
	result["username"] = t.Username
	result["status"] = t.Status

	return result
}
