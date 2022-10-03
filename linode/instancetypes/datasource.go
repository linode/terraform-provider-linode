package instancetypes

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
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listTypes, flattenType)
	if err != nil {
		return nil
	}

	d.Set("types", results)

	return nil
}

func listTypes(
	ctx context.Context, d *schema.ResourceData, client *linodego.Client,
	options *linodego.ListOptions,
) ([]interface{}, error) {
	types, err := client.ListTypes(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(types))

	for i, v := range types {
		result[i] = v
	}

	return result, nil
}

func flattenType(data interface{}) map[string]interface{} {
	t := data.(linodego.LinodeType)

	result := make(map[string]interface{})

	result["id"] = t.ID
	result["label"] = t.Label
	result["disk"] = t.Disk
	result["class"] = t.Class
	result["network_out"] = t.NetworkOut
	result["memory"] = t.Memory
	result["transfer"] = t.Transfer
	result["vcpus"] = t.VCPUs

	result["price"] = []map[string]float32{
		{
			"hourly":  t.Price.Hourly,
			"monthly": t.Price.Monthly,
		},
	}

	result["addons"] = []map[string]interface{}{
		{
			"backups": []map[string]interface{}{
				{
					"price": []map[string]float32{
						{
							"hourly":  t.Addons.Backups.Price.Hourly,
							"monthly": t.Addons.Backups.Price.Monthly,
						},
					},
				},
			},
		},
	}

	return result
}
