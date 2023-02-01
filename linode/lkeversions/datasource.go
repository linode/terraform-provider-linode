package lkeversions

import (
	"context"
	"encoding/json"

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
	client := meta.(*helper.ProviderMeta).Client

	data, err := listLKEVersions(ctx, d, &client, nil)
	if err != nil {
		return diag.Errorf("Error getting lke versions: %s", err)
	}

	id, err := json.Marshal(data)
	if err != nil {
		return diag.Errorf("failed to marshal id: %s", err)
	}

	d.SetId(string(id))
	d.Set("versions", data)

	return nil
}

func listLKEVersions(
	ctx context.Context, d *schema.ResourceData, client *linodego.Client,
	options *linodego.ListOptions,
) ([]interface{}, error) {
	types, err := client.ListLKEVersions(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(types))

	for i, v := range types {
		result[i] = flattenLKEVersions(v)
	}

	return result, nil
}

func flattenLKEVersions(data interface{}) map[string]interface{} {
	t := data.(linodego.LKEVersion)

	result := make(map[string]interface{})

	result["id"] = t.ID

	return result
}
