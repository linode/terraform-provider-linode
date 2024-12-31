package databaseaccesscontrols

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
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

	dbID, dbType, err := parseID(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse database id: %s", err)
	}

	allowList, err := getDBAllowListByEngine(ctx, client, dbType, dbID)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] allow_list %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("failed to get allow list for database %d: %s", dbID, err)
	}

	d.Set("database_id", dbID)
	d.Set("database_type", dbType)
	d.Set("allow_list", allowList)

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dbID := d.Get("database_id").(int)
	dbType := d.Get("database_type").(string)

	d.SetId(formatID(dbID, dbType))

	return updateResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	dbID, dbType, err := parseID(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse database id: %s", err)
	}

	if d.HasChange("allow_list") {
		allowList := helper.ExpandStringSet(d.Get("allow_list").(*schema.Set))

		if err := updateDBAllowListByEngine(ctx, client, d, dbType, dbID, allowList); err != nil {
			return diag.Errorf("failed to update allow_list for database %d: %s", dbID, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	dbID, dbType, err := parseID(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse database id: %s", err)
	}

	if err := updateDBAllowListByEngine(ctx, client, d, dbType, dbID, []string{}); err != nil {
		return diag.Errorf("failed to update allow_list for database %d: %s", dbID, err)
	}

	d.SetId("")

	return nil
}

func updateDBAllowListByEngine(
	ctx context.Context,
	client linodego.Client,
	d *schema.ResourceData,
	engine string,
	id int,
	allowList []string,
) error {
	updatePoller, err := client.NewEventPoller(ctx, id, linodego.EntityDatabase, linodego.ActionDatabaseUpdate)
	if err != nil {
		return fmt.Errorf("failed to create update EventPoller: %w", err)
	}

	switch engine {
	case "mysql":
		if _, err := client.UpdateMySQLDatabase(ctx, id, linodego.MySQLUpdateOptions{
			AllowList: &allowList,
		}); err != nil {
			return err
		}

	case "postgresql":
		if _, err := client.UpdatePostgresDatabase(ctx, id, linodego.PostgresUpdateOptions{
			AllowList: &allowList,
		}); err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid database engine: %s", engine)
	}

	timeoutSeconds, err := helper.SafeFloat64ToInt(d.Timeout(schema.TimeoutUpdate).Seconds())
	if err != nil {
		return err
	}

	if _, err := updatePoller.WaitForFinished(ctx, timeoutSeconds); err != nil {
		return fmt.Errorf("failed to wait for update event completion: %w", err)
	}

	return nil
}

func getDBAllowListByEngine(ctx context.Context, client linodego.Client, engine string, id int) ([]string, error) {
	switch engine {
	case "mysql":
		db, err := client.GetMySQLDatabase(ctx, id)
		if err != nil {
			return nil, err
		}

		return db.AllowList, nil
	case "postgresql":
		db, err := client.GetPostgresDatabase(ctx, id)
		if err != nil {
			return nil, err
		}

		return db.AllowList, nil
	}

	return nil, fmt.Errorf("invalid database type: %s", engine)
}

func formatID(dbID int, dbType string) string {
	return fmt.Sprintf("%d:%s", dbID, dbType)
}

func parseID(id string) (int, string, error) {
	split := strings.Split(id, ":")
	if len(split) != 2 {
		return 0, "", fmt.Errorf("invalid number of segments")
	}

	dbID, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, "", err
	}

	return dbID, split[1], nil
}
