package instance

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var filterableFields = []string{"group", "id", "image", "label", "region", "tags"}

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Schema: instanceDataSourceSchema,
	}
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema: map[string]*schema.Schema{
			"filter":   helper.FilterSchema(filterableFields),
			"order_by": helper.OrderBySchema(filterableFields),
			"order":    helper.OrderSchema(),
			"instances": {
				Type:        schema.TypeList,
				Description: "The returned list of Instances.",
				Computed:    true,
				Elem:        dataSourceInstance(),
			},
		},
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	filterID, err := helper.GetFilterID(d)
	if err != nil {
		return diag.Errorf("failed to generate filter id: %s", err)
	}

	filter, err := helper.ConstructFilterString(d, instanceValueToFilterType)
	if err != nil {
		return diag.Errorf("failed to construct filter: %s", err)
	}

	instances, err := client.ListInstances(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return diag.Errorf("failed to get instances: %s", err)
	}

	flattenedInstances := make([]interface{}, len(instances))
	for i, instance := range instances {
		instanceMap, err := flattenInstance(ctx, &client, &instance)
		if err != nil {
			return diag.Errorf("failed to translate instance to map: %s", err)
		}

		// Merge additional fields
		instanceMap["id"] = instance.ID

		flattenedInstances[i] = instanceMap
	}

	instancesFiltered, err := helper.FilterResults(d, flattenedInstances)
	if err != nil {
		return diag.Errorf("failed to filter returned instances: %s", err)
	}

	d.SetId(filterID)
	d.Set("instances", instancesFiltered)

	return nil
}

// instanceValueToFilterType converts the given value to the correct type depending on the filter name.
func instanceValueToFilterType(filterName, value string) (interface{}, error) {
	switch filterName {
	case "id":
		return strconv.Atoi(value)
	}

	return value, nil
}
