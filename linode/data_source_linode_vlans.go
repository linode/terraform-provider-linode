package linode

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"

	"context"
	"fmt"
	"time"
)

func dataSourceLinodeVLAN() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The unique label of this VLAN.",
				Computed:    true,
			},
			"linodes": {
				Type:        schema.TypeList,
				Description: "The Linodes currently attached to this VLAN.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The region this VLAN is located in.",
				Computed:    true,
			},
			"created": {
				Type:        schema.TypeString,
				Description: "When this VLAN was created.",
				Computed:    true,
			},
		},
	}
}

func dataSourceLinodeVLANs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeVLANsRead,
		Schema: map[string]*schema.Schema{
			"filter": filterSchema([]string{"label", "region"}),
			"vlans": {
				Type:        schema.TypeList,
				Description: "The returned list of VLANs.",
				Computed:    true,
				Elem:        dataSourceLinodeVLAN(),
			},
		},
	}
}

func dataSourceLinodeVLANsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client

	filter, err := constructFilterString(d, vlanValueToFilterType)
	if err != nil {
		return fmt.Errorf("failed to construct filter: %s", err)
	}

	vlans, err := client.ListVLANs(context.Background(), &linodego.ListOptions{
		Filter: filter,
	})

	if err != nil {
		return fmt.Errorf("failed to list linode vlans: %s", err)
	}

	vlansFlattened := make([]interface{}, len(vlans))
	for i, vlan := range vlans {
		vlansFlattened[i] = flattenLinodeVLAN(&vlan)
	}

	d.SetId(filter)
	d.Set("vlans", vlansFlattened)

	return nil
}

func vlanValueToFilterType(_, value string) (interface{}, error) {
	return value, nil
}

func flattenLinodeVLAN(vlan *linodego.VLAN) map[string]interface{} {
	result := make(map[string]interface{})

	result["label"] = vlan.Label
	result["linodes"] = vlan.Linodes
	result["region"] = vlan.Region

	if vlan.Created != nil {
		result["created"] = vlan.Created.Format(time.RFC3339)
	}

	return result
}
