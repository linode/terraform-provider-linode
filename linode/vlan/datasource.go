package vlan

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"

	"context"
	"time"
)

func dataSourceVLAN() *schema.Resource {
	return &schema.Resource{
		Schema: resourceSchema,
	}
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema: map[string]*schema.Schema{
			"order_by": helper.OrderBySchema(filterConfig),
			"order":    helper.OrderSchema(),
			"filter":   helper.FilterSchema(filterConfig),
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
	results, err := helper.FilterResource(ctx, d, meta, filterConfig, listVLANs, flattenVLAN)
	if err != nil {
		return nil
	}

	d.Set("vlans", results)

	return nil
}

func listVLANs(
	ctx context.Context, client *linodego.Client, options *linodego.ListOptions) ([]interface{}, error) {
	vlans, err := client.ListVLANs(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(vlans))

	for i, v := range vlans {
		result[i] = v
	}

	return result, nil
}

func flattenVLAN(data interface{}) map[string]interface{} {
	vlan := data.(linodego.VLAN)

	result := make(map[string]interface{})

	result["label"] = vlan.Label
	result["linodes"] = vlan.Linodes
	result["region"] = vlan.Region

	if vlan.Created != nil {
		result["created"] = vlan.Created.Format(time.RFC3339)
	}

	return result
}
