package databaseaccesscontrols

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"strconv"
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

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id for database firewall: %s", err)
	}

	db, err := getDBByID(ctx, client, id)
	if err != nil {
		return diag.Errorf("failed to find the specified mysql database: %s", err)
	}

	d.Set("database_id", id)
	d.Set("allow_list", db.AllowList)

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dbID := d.Get("database_id").(int)

	d.SetId(strconv.Itoa(dbID))

	return updateResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id for database firewall: %s", err)
	}

	if d.HasChange("allow_list") {
		if err := updateAllowList(ctx, d, client, id, helper.ExpandStringSet(d.Get("allow_list").(*schema.Set))); err != nil {
			return diag.Errorf("failed to update allow_list for database %d: %s", id, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id for database firewall: %s", err)
	}

	if err := updateAllowList(ctx, d, client, id, []string{}); err != nil {
		return diag.Errorf("failed to update allow_list for database %d: %s", id, err)
	}

	d.SetId("")

	return nil
}

func updateAllowList(ctx context.Context, d *schema.ResourceData,
	client linodego.Client, dbID int, allowList []string,
) error {
	db, err := getDBByID(ctx, client, dbID)
	if err != nil {
		return err
	}

	if err := updateDBAllowListByEngine(ctx, client, db.Engine, dbID, allowList); err != nil {
		return err
	}

	return helper.WaitForDatabaseUpdated(ctx, client, dbID, linodego.DatabaseEngineType(db.Engine),
		db.Created, int(d.Timeout(schema.TimeoutUpdate).Seconds()))
}

var updateDBAllowListEngineMap = map[string]func(context.Context, linodego.Client, string, int, []string) error{
	"mysql": func(ctx context.Context, client linodego.Client, engine string, id int, allowList []string) error {
		_, err := client.UpdateMySQLDatabase(ctx, id, linodego.MySQLUpdateOptions{
			AllowList: &allowList,
		})
		return err
	},
	"mongodb": func(ctx context.Context, client linodego.Client, engine string, id int, allowList []string) error {
		_, err := client.UpdateMongoDatabase(ctx, id, linodego.MongoUpdateOptions{
			AllowList: &allowList,
		})
		return err
	},
}

func updateDBAllowListByEngine(ctx context.Context, client linodego.Client,
	engine string, id int, allowList []string,
) error {
	// Future-proofing for more DB types
	f, ok := updateDBAllowListEngineMap[engine]
	if !ok {
		return fmt.Errorf("invalid database engine: %s", engine)
	}

	return f(ctx, client, engine, id, allowList)
}

func getDBByID(ctx context.Context, client linodego.Client, id int) (*linodego.Database, error) {
	instances, err := client.ListDatabases(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, db := range instances {
		if db.ID == id {
			return &db, nil
		}
	}

	return nil, fmt.Errorf("unable to find database with id %d", id)
}
