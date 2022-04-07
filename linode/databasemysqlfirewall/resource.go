package databasemysqlfirewall

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MySQL database ID %s as int: %s", d.Id(), err)
	}

	db, err := client.GetMySQLDatabase(ctx, int(id))
	if err != nil {
		return diag.Errorf("failed to find the specified mysql database: %s", err)
	}

	d.Set("database_id", db.ID)
	d.Set("allow_list", db.AllowList)

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	dbID := d.Get("database_id").(int)

	d.SetId(strconv.Itoa(dbID))

	if err := updateAllowList(ctx, d, client, dbID,
		helper.ExpandStringSet(d.Get("allow_list").(*schema.Set))); err != nil {
		return diag.Errorf("failed to update allow_list for database %d: %s", dbID, err)
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MySQL database ID %s as int: %s", d.Id(), err)
	}

	if d.HasChange("allow_list") {
		if err := updateAllowList(ctx, d, client, int(id),
			helper.ExpandStringSet(d.Get("allow_list").(*schema.Set))); err != nil {
			return diag.Errorf("failed to update allow_list for database %d: %s", id, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MySQL database ID %s as int: %s", d.Id(), err)
	}

	if err := updateAllowList(ctx, d, client, int(id),
		[]string{}); err != nil {
		return diag.Errorf("failed to update allow_list for database %d: %s", id, err)
	}

	d.SetId("")

	return nil
}

func updateAllowList(ctx context.Context, d *schema.ResourceData,
	client linodego.Client, dbID int, allowList []string) error {
	db, err := client.GetMySQLDatabase(ctx, dbID)
	if err != nil {
		return err
	}

	return resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		updateOpts := linodego.MySQLUpdateOptions{
			Label:     db.Label,
			AllowList: allowList,
		}

		if _, err := client.UpdateMySQLDatabase(ctx, dbID, updateOpts); err != nil {
			if lerr, ok := err.(*linodego.Error); ok &&
				lerr.Code == 500 && strings.Contains(lerr.Message, "Unable to update allow_list on database") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(
				fmt.Errorf("failed to update mysql database allow_list %d: %s", db.ID, err))
		}

		return nil
	})

	return nil
}
