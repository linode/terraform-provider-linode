package regions

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
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listRegions, flattenRegions)
	if err != nil {
		return nil
	}

	d.Set("regions", results)

	return nil
}

func listRegions(
	ctx context.Context, d *schema.ResourceData, client *linodego.Client,
	options *linodego.ListOptions,
) ([]interface{}, error) {
	types, err := client.ListRegions(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(types))

	for i, v := range types {
		result[i] = v
	}

	return result, nil
}

func flattenRegions(data interface{}) map[string]interface{} {
	t := data.(linodego.Region)

	result := make(map[string]interface{})

	result["capabilities"] = t.Capabilities
	result["country"] = t.Country
	result["id"] = t.ID
	result["label"] = t.Label
	result["resolvers"] =
		[]map[string]interface{}{{
			"ipv4": t.Resolvers.IPv4,
			"ipv6": t.Resolvers.IPv6,
		}}
	result["status"] = t.Status

	return result
}
