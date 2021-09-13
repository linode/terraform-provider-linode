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

	types, err := client.ListTypes(ctx, nil)
	if err != nil {
		return diag.Errorf("Error listing ranges: %s", err)
	}

	reqType := d.Get("id").(string)

	for _, r := range types {
		if r.ID == reqType {
			d.SetId(r.ID)
			d.Set("label", r.Label)
			d.Set("disk", r.Disk)
			d.Set("memory", r.Memory)
			d.Set("vcpus", r.VCPUs)
			d.Set("network_out", r.NetworkOut)
			d.Set("transfer", r.Transfer)
			d.Set("class", r.Class)

			d.Set("price", []map[string]interface{}{{
				"hourly":  r.Price.Hourly,
				"monthly": r.Price.Monthly,
			}})

			d.Set("addons", []map[string]interface{}{{
				"backups": []map[string]interface{}{{
					"price": []map[string]interface{}{{
						"hourly":  r.Addons.Backups.Price.Hourly,
						"monthly": r.Addons.Backups.Price.Monthly,
					}},
				}},
			}})
			return nil
		}
	}

	d.SetId("")

	return diag.Errorf("Instance Type %s was not found", reqType)
}
