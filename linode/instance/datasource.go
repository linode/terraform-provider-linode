package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Schema: instanceDataSourceSchema,
	}
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema: map[string]*schema.Schema{
			"filter": filterSchema([]string{"group", "id", "image", "label", "region", "tags"}),
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

	filter, err := constructFilterString(d, instanceValueToFilterType)
	if err != nil {
		return diag.Errorf("failed to construct filter: %s", err)
	}

	instances, err := client.ListInstances(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return diag.Errorf("failed to get instances: %s", err)
	}

	flattenedInstances := make([]map[string]interface{}, len(instances))
	for i, instance := range instances {
		instanceMap, err := flattenLinodeInstance(ctx, &client, &instance)
		if err != nil {
			return diag.Errorf("failed to translate instance to map: %s", err)
		}

		flattenedInstances[i] = instanceMap
	}

	d.SetId(fmt.Sprintf(filter))
	d.Set("instances", flattenedInstances)

	return nil
}
