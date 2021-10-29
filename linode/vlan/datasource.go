package vlan

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"

	"context"
	"time"
)

var filterableFields = []string{"label", "region"}

func dataSourceVLAN() *schema.Resource {
	return &schema.Resource{
		Schema: resourceSchema,
	}
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema: map[string]*schema.Schema{
			"order_by": helper.OrderBySchema(filterableFields),
			"order":    helper.OrderSchema(),
			"filter":   helper.FilterSchema(filterableFields),
			"vlans": {
				Type:        schema.TypeList,
				Description: "The returned list of VLANs.",
				Computed:    true,
				Elem:        dataSourceVLAN(),
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

	filter, err := helper.ConstructFilterString(d, vlanValueToFilterType)
	if err != nil {
		return diag.Errorf("failed to construct filter: %s", err)
	}

	vlans, err := client.ListVLANs(ctx, &linodego.ListOptions{
		Filter: filter,
	})

	if err != nil {
		return diag.Errorf("failed to list linode vlans: %s", err)
	}

	vlansFlattened := make([]interface{}, len(vlans))
	for i, vlan := range vlans {
		vlansFlattened[i] = flattenVLAN(&vlan)
	}

	vlansFiltered, err := helper.FilterResults(d, vlansFlattened)
	if err != nil {
		return diag.Errorf("failed to filter returned vlans: %s", err)
	}

	d.SetId(filterID)
	d.Set("vlans", vlansFiltered)

	return nil
}

func vlanValueToFilterType(_, value string) (interface{}, error) {
	return value, nil
}

func flattenVLAN(vlan *linodego.VLAN) map[string]interface{} {
	result := make(map[string]interface{})

	result["label"] = vlan.Label
	result["linodes"] = vlan.Linodes
	result["region"] = vlan.Region

	if vlan.Created != nil {
		result["created"] = vlan.Created.Format(time.RFC3339)
	}

	return result
}
