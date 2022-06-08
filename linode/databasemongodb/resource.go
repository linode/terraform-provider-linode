package databasemongodb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

const (
	createDBTimeout = 60 * time.Minute
	updateDBTimeout = 5 * time.Minute
	deleteDBTimeout = 5 * time.Minute
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(createDBTimeout),
			Update: schema.DefaultTimeout(updateDBTimeout),
			Delete: schema.DefaultTimeout(deleteDBTimeout),
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MongoDB database ID %s as int: %s", d.Id(), err)
	}

	db, err := client.GetMongoDatabase(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing MongoDB database ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("failed to find the specified mongodb database: %s", err)
	}

	cert, err := client.GetMongoDatabaseSSL(ctx, int(id))
	if err != nil {
		return diag.Errorf("failed to get cert for the specified mongodb database: %s", err)
	}

	creds, err := client.GetMongoDatabaseCredentials(ctx, int(id))
	if err != nil {
		return diag.Errorf("failed to get credentials for the specified mongodb database: %s", err)
	}

	d.Set("engine_id", helper.CreateDatabaseEngineSlug(db.Engine, db.Version))
	d.Set("engine", db.Engine)
	d.Set("label", db.Label)
	d.Set("region", db.Region)
	d.Set("type", db.Type)
	d.Set("allow_list", db.AllowList)
	d.Set("cluster_size", db.ClusterSize)
	d.Set("compression_type", db.CompressionType)
	d.Set("encrypted", db.Encrypted)
	d.Set("storage_engine", db.StorageEngine)
	d.Set("ssl_connection", db.SSLConnection)
	d.Set("ca_cert", string(cert.CACertificate))
	d.Set("created", db.Created.Format(time.RFC3339))
	d.Set("host_primary", db.Hosts.Primary)
	d.Set("host_secondary", db.Hosts.Secondary)
	d.Set("peers", db.Peers)
	d.Set("replica_set", db.ReplicaSet)
	d.Set("port", db.Port)
	d.Set("root_password", creds.Password)
	d.Set("status", db.Status)
	d.Set("updated", db.Updated.Format(time.RFC3339))
	d.Set("root_username", creds.Username)
	d.Set("version", db.Version)
	d.Set("updates", []interface{}{helper.FlattenMaintenanceWindow(db.Updates)})

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	db, err := client.CreateMongoDatabase(ctx, linodego.MongoCreateOptions{
		Label:           d.Get("label").(string),
		Region:          d.Get("region").(string),
		Type:            d.Get("type").(string),
		Engine:          d.Get("engine_id").(string),
		Encrypted:       d.Get("encrypted").(bool),
		ClusterSize:     d.Get("cluster_size").(int),
		CompressionType: linodego.MongoCompressionType(d.Get("compression_type").(string)),
		StorageEngine:   linodego.MongoStorageEngine(d.Get("storage_engine").(string)),
		SSLConnection:   d.Get("ssl_connection").(bool),
		AllowList:       helper.ExpandStringSet(d.Get("allow_list").(*schema.Set)),
	})
	if err != nil {
		return diag.Errorf("failed to create mongodb database: %s", err)
	}

	d.SetId(strconv.Itoa(db.ID))

	_, err = client.WaitForEventFinished(ctx, db.ID, linodego.EntityDatabase,
		linodego.ActionDatabaseCreate, *db.Created, int(d.Timeout(schema.TimeoutCreate).Seconds()))
	if err != nil {
		return diag.Errorf("failed to wait for mongodb database creation: %s", err)
	}

	updateList := d.Get("updates").([]interface{})

	if !d.GetRawConfig().GetAttr("updates").IsNull() && len(updateList) > 0 {
		updates, err := helper.ExpandMaintenanceWindow(updateList[0].(map[string]interface{}))
		if err != nil {
			return diag.Errorf("failed to read maintenance window config: %s", err)
		}

		updatedDB, err := client.UpdateMongoDatabase(ctx, db.ID, linodego.MongoUpdateOptions{
			Updates: &updates,
		})
		if err != nil {
			return diag.Errorf("failed to update mongodb database maintenance window: %s", err)
		}

		err = helper.WaitForDatabaseUpdated(ctx, client, db.ID,
			linodego.DatabaseEngineTypeMongo, updatedDB.Created, int(d.Timeout(schema.TimeoutUpdate).Seconds()))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MongoDB database ID %s as int: %s", d.Id(), err)
	}

	updateOpts := linodego.MongoUpdateOptions{
		Label: d.Get("label").(string),
	}

	shouldUpdate := false

	if d.HasChange("allow_list") {
		allowList := helper.ExpandStringSet(d.Get("allow_list").(*schema.Set))
		updateOpts.AllowList = &allowList
		shouldUpdate = true
	}

	if d.HasChange("updates") {
		var updates *linodego.DatabaseMaintenanceWindow

		updatesRaw := d.Get("updates")
		if updatesRaw != nil && len(updatesRaw.([]interface{})) > 0 {
			expanded, err := helper.ExpandMaintenanceWindow(updatesRaw.([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return diag.Errorf("failed to update maintenance window: %s", err)
			}

			updates = &expanded
		}

		updateOpts.Updates = updates
		shouldUpdate = true
	}

	if shouldUpdate {
		updatedDB, err := client.UpdateMongoDatabase(ctx, int(id), updateOpts)
		if err != nil {
			return diag.Errorf("failed to update mongodb database: %s", err)
		}

		err = helper.WaitForDatabaseUpdated(ctx, client, int(id),
			linodego.DatabaseEngineTypeMongo, updatedDB.Created, int(d.Timeout(schema.TimeoutUpdate).Seconds()))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MongoDB database ID %s as int: %s", d.Id(), err)
	}

	// We should retry on intermittent deletion errors
	return diag.FromErr(resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		err := client.DeleteMongoDatabase(ctx, int(id))
		if err != nil {
			if lerr, ok := err.(*linodego.Error); ok &&
				lerr.Code == 500 && strings.Contains(lerr.Message, "Unable to delete instance") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(fmt.Errorf("failed to delete mongo database %d: %s", id, err))
		}

		return nil
	}))
}
