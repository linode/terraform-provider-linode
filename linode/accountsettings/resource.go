package accountsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	settings, err := client.GetAccountSettings(ctx)
	if err != nil {
		return diag.Errorf("Error getting account settings: %s", err)
	}

	d.Set("backups_enabled", settings.BackupsEnabled)
	d.Set("longview_subscription", settings.LongviewSubscription)
	d.Set("managed", settings.Managed)
	d.Set("network_helper", settings.NetworkHelper)
	d.Set("object_storage", settings.ObjectStorage)

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	account, err := client.GetAccount(ctx)
	if err != nil {
		return diag.Errorf("Error getting account: %s", err)
	}

	d.SetId(account.Email)
	return updateResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	_, errSettings := client.GetAccountSettings(ctx)
	if errSettings != nil {
		return diag.Errorf("Error fetching the account settings: %s", errSettings)
	}

	accountUpdateOpts := linodego.AccountSettingsUpdateOptions{}
	longviewUpdateOpts := linodego.LongviewPlanUpdateOptions{}

	accountUpdate := false
	longviewUpdate := false

	if d.HasChange("backups_enabled") {
		backupsEnabled := d.Get("backups_enabled").(bool)
		accountUpdateOpts.BackupsEnabled = &backupsEnabled
		accountUpdate = true

	}

	if d.HasChange("network_helper") {
		networkHelper := d.Get("network_helper").(bool)
		accountUpdateOpts.NetworkHelper = &networkHelper
		accountUpdate = true
	}

	if d.HasChange("longview_subscription") {
		longviewUpdateOpts.LongviewSubscription = d.Get("longview_subscription").(string)
		longviewUpdate = true
	}

	if accountUpdate {
		_, updateErr := client.UpdateAccountSettings(ctx, accountUpdateOpts)
		if updateErr != nil {
			return diag.Errorf("Error updating the account settings: %s", updateErr)
		}
	}

	if longviewUpdate {
		_, updateErr := client.UpdateLongviewPlan(ctx, longviewUpdateOpts)
		if updateErr != nil {
			return diag.Errorf("Error updating the longview plan: %s", updateErr)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
