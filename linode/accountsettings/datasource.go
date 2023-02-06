package accountsettings

import (
	"context"
	"encoding/json"

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

	data, err := client.GetAccountSettings(ctx)
	if err != nil {
		return diag.Errorf("Error getting account settings: %s", err)
	}

	id, err := json.Marshal(data)
	if err != nil {
		return diag.Errorf("failed to marshal id: %s", err)
	}

	d.SetId(string(id))
	d.Set("backups_enabled", data.BackupsEnabled)
	d.Set("longview_subscription", data.LongviewSubscription)
	d.Set("managed", data.Managed)
	d.Set("network_helper", data.NetworkHelper)
	d.Set("object_storage", data.ObjectStorage)

	return nil
}
