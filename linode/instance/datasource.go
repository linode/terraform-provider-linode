package instance

import (
	"context"
	"fmt"
	"strconv"

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

func flattenLinodeInstance(
	ctx context.Context, client *linodego.Client, instance *linodego.Instance) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	id := instance.ID

	instanceNetwork, err := client.GetInstanceIPAddresses(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get ips for linode instance %d: %s", id, err)
	}

	var ips []string
	for _, ip := range instance.IPv4 {
		ips = append(ips, ip.String())
	}

	result["ipv4"] = ips
	result["ipv6"] = instance.IPv6

	public, private := instanceNetwork.IPv4.Public, instanceNetwork.IPv4.Private

	if len(public) > 0 {
		result["ip_address"] = public[0].Address
	}

	if len(private) > 0 {
		result["private_ip_address"] = private[0].Address
	}

	result["label"] = instance.Label
	result["status"] = instance.Status
	result["type"] = instance.Type
	result["region"] = instance.Region
	result["watchdog_enabled"] = instance.WatchdogEnabled
	result["group"] = instance.Group
	result["tags"] = instance.Tags
	result["image"] = instance.Image

	result["backups"] = flattenInstanceBackups(*instance)
	result["specs"] = flattenInstanceSpecs(*instance)
	result["alerts"] = flattenInstanceAlerts(*instance)

	instanceDisks, err := client.ListInstanceDisks(ctx, int(id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the disks for the Linode instance %d: %s", id, err)
	}

	disks, swapSize := flattenInstanceDisks(instanceDisks)
	result["disk"] = disks
	result["swap_size"] = swapSize

	instanceConfigs, err := client.ListInstanceConfigs(ctx, int(id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the config for Linode instance %d (%s): %s", id, instance.Label, err)
	}

	diskLabelIDMap := make(map[int]string, len(instanceDisks))
	for _, disk := range instanceDisks {
		diskLabelIDMap[disk.ID] = disk.Label
	}

	configs := flattenInstanceConfigs(instanceConfigs, diskLabelIDMap)

	result["config"] = configs
	if len(instanceConfigs) == 1 {
		result["boot_config_label"] = instanceConfigs[0].Label
	}

	return result, nil
}

// instanceValueToFilterType converts the given value to the correct type depending on the filter name.
func instanceValueToFilterType(filterName, value string) (interface{}, error) {
	switch filterName {
	case "id":
		return strconv.Atoi(value)
	}

	return value, nil
}
