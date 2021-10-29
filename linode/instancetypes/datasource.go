package instancetypes

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"

	"context"
	"strconv"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	filterID, err := helper.GetFilterID(d)
	if err != nil {
		return diag.Errorf("failed to generate filter id: %s", err)
	}

	filter, err := helper.ConstructFilterString(d, typeValueToFilterType)
	if err != nil {
		return diag.Errorf("failed to construct filter: %s", err)
	}

	types, err := client.ListTypes(ctx, &linodego.ListOptions{
		Filter: filter,
	})

	if err != nil {
		return diag.Errorf("failed to list linode types: %s", err)
	}

	typesFlattened := make([]interface{}, len(types))
	for i, t := range types {
		typesFlattened[i] = flattenType(&t)
	}

	typesFiltered, err := helper.FilterResults(d, typesFlattened)
	if err != nil {
		return diag.Errorf("failed to filter returned types: %s", err)
	}

	d.SetId(filterID)
	d.Set("types", typesFiltered)

	return nil
}

func flattenType(t *linodego.LinodeType) map[string]interface{} {
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

func typeValueToFilterType(filterName, value string) (interface{}, error) {
	switch filterName {
	case "disk", "gpus", "memory", "transfer", "vcpus":
		return strconv.Atoi(value)
	}

	return value, nil
}
