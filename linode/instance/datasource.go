package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var filterConfig = helper.FilterConfig{
	"group":          {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"id":             {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"image":          {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"label":          {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"region":         {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"lke_cluster_id": {APIFilterable: true, TypeFunc: helper.FilterTypeInt},

	// Tags must be filtered on the client
	"tags":             {TypeFunc: helper.FilterTypeString},
	"status":           {TypeFunc: helper.FilterTypeString},
	"type":             {TypeFunc: helper.FilterTypeString},
	"watchdog_enabled": {TypeFunc: helper.FilterTypeBool},
	"disk_encryption":  {TypeFunc: helper.FilterTypeString},
}

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Schema: instanceDataSourceSchema,
	}
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema: map[string]*schema.Schema{
			"filter":   filterConfig.FilterSchema(),
			"order_by": filterConfig.OrderBySchema(),
			"order":    filterConfig.OrderSchema(),
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
	tflog.Debug(ctx, "Read data.linode_instances")

	client := meta.(*helper.ProviderMeta).Client

	filterID, err := filterConfig.GetFilterID(d)
	if err != nil {
		return diag.Errorf("failed to generate filter id: %s", err)
	}

	filter, err := filterConfig.ConstructFilterString(d)
	if err != nil {
		return diag.Errorf("failed to construct filter: %s", err)
	}

	instances, err := client.ListInstances(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return diag.Errorf("failed to get instances: %s", err)
	}

	instanceIDMap := make(map[int]linodego.Instance, len(instances))

	// Create a list of filterable instance maps
	flattenedInstances := make([]interface{}, len(instances))
	for i, instance := range instances {
		instance := instance
		instanceIDMap[instance.ID] = instance

		instanceMap, err := flattenInstanceSimple(&instance)
		if err != nil {
			return diag.Errorf("failed to translate instance to filterable map: %s", err)
		}

		flattenedInstances[i] = instanceMap
	}

	instancesFiltered, err := filterConfig.FilterResults(d, flattenedInstances)
	if err != nil {
		return diag.Errorf("failed to filter returned instances: %s", err)
	}

	// Fully populate returned instances
	for i, instance := range instancesFiltered {
		instanceObject := instanceIDMap[instance["id"].(int)]

		instanceMap, err := flattenInstance(ctx, &client, &instanceObject)
		if err != nil {
			return diag.Errorf("failed to translate instance to map: %s", err)
		}

		instancesFiltered[i] = instanceMap
	}

	d.SetId(filterID)
	d.Set("instances", instancesFiltered)

	return nil
}
