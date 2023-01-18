package instancetype

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id := d.Get("id").(string)

	typeInfo, err := client.GetType(ctx, id)
	if err != nil {
		return diag.Errorf("Error getting type %s: %s", id, err)
	}

	d.SetId(typeInfo.ID)
	d.Set("label", typeInfo.Label)
	d.Set("disk", typeInfo.Disk)
	d.Set("memory", typeInfo.Memory)
	d.Set("vcpus", typeInfo.VCPUs)
	d.Set("network_out", typeInfo.NetworkOut)
	d.Set("transfer", typeInfo.Transfer)
	d.Set("class", typeInfo.Class)

	d.Set("price", []map[string]interface{}{{
		"hourly":  typeInfo.Price.Hourly,
		"monthly": typeInfo.Price.Monthly,
	}})

	d.Set("addons", []map[string]interface{}{{
		"backups": []map[string]interface{}{{
			"price": []map[string]interface{}{{
				"hourly":  typeInfo.Addons.Backups.Price.Hourly,
				"monthly": typeInfo.Addons.Backups.Price.Monthly,
			}},
		}},
	}})

	return nil
}
