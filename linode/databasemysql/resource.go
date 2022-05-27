package databasemysql

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
		return diag.Errorf("Error parsing Linode MySQL database ID %s as int: %s", d.Id(), err)
	}

	db, err := client.GetMySQLDatabase(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing MySQL database ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("failed to find the specified mysql database: %s", err)
	}

	cert, err := client.GetMySQLDatabaseSSL(ctx, int(id))
	if err != nil {
		return diag.Errorf("failed to get cert for the specified mysql database: %s", err)
	}

	creds, err := client.GetMySQLDatabaseCredentials(ctx, int(id))
	if err != nil {
		return diag.Errorf("failed to get credentials for the specified mysql database: %s", err)
	}

	d.Set("engine_id", createEngineSlug(db.Engine, db.Version))
	d.Set("engine", db.Engine)
	d.Set("label", db.Label)
	d.Set("region", db.Region)
	d.Set("type", db.Type)
	d.Set("allow_list", db.AllowList)
	d.Set("cluster_size", db.ClusterSize)
	d.Set("encrypted", db.Encrypted)
	d.Set("replication_type", db.ReplicationType)
	d.Set("ssl_connection", db.SSLConnection)
	d.Set("ca_cert", string(cert.CACertificate))
	d.Set("created", db.Created.Format(time.RFC3339))
	d.Set("host_primary", db.Hosts.Primary)
	d.Set("host_secondary", db.Hosts.Secondary)
	d.Set("root_password", creds.Password)
	d.Set("status", db.Status)
	d.Set("updated", db.Updated.Format(time.RFC3339))
	d.Set("root_username", creds.Username)
	d.Set("version", db.Version)
	d.Set("updates", []interface{}{FlattenMaintenanceWindow(db.Updates)})

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	db, err := client.CreateMySQLDatabase(ctx, linodego.MySQLCreateOptions{
		Label:           d.Get("label").(string),
		Region:          d.Get("region").(string),
		Type:            d.Get("type").(string),
		Engine:          d.Get("engine_id").(string),
		Encrypted:       d.Get("encrypted").(bool),
		ClusterSize:     d.Get("cluster_size").(int),
		ReplicationType: d.Get("replication_type").(string),
		SSLConnection:   d.Get("ssl_connection").(bool),
		AllowList:       helper.ExpandStringSet(d.Get("allow_list").(*schema.Set)),
	})
	if err != nil {
		return diag.Errorf("failed to create mysql database: %s", err)
	}

	d.SetId(strconv.Itoa(db.ID))

	_, err = client.WaitForEventFinished(ctx, db.ID, linodego.EntityDatabase,
		linodego.ActionDatabaseCreate, *db.Created, int(d.Timeout(schema.TimeoutCreate).Seconds()))
	if err != nil {
		return diag.Errorf("failed to wait for mysql database creation: %s", err)
	}

	updateList := d.Get("updates").([]interface{})

	if !d.GetRawConfig().GetAttr("updates").IsNull() && len(updateList) > 0 {
		updates, err := ExpandMaintenanceWindow(updateList[0].(map[string]interface{}))
		if err != nil {
			return diag.Errorf("failed to read maintenance window config: %s", err)
		}

		_, err = client.UpdateMySQLDatabase(ctx, db.ID, linodego.MySQLUpdateOptions{
			Updates: &updates,
		})
		if err != nil {
			return diag.Errorf("failed to update mysql database maintenance window: %s", err)
		}
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MySQL database ID %s as int: %s", d.Id(), err)
	}

	updateOpts := linodego.MySQLUpdateOptions{
		Label: d.Get("label").(string),
	}

	if d.HasChange("allow_list") {
		updateOpts.AllowList = helper.ExpandStringSet(d.Get("allow_list").(*schema.Set))
	}

	if d.HasChange("updates") {
		var updates *linodego.MySQLDatabaseMaintenanceWindow

		updatesRaw := d.Get("updates")
		if updatesRaw != nil && len(updatesRaw.([]interface{})) > 0 {
			expanded, err := ExpandMaintenanceWindow(updatesRaw.([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return diag.Errorf("failed to update maintenance window: %s", err)
			}

			updates = &expanded
		}

		updateOpts.Updates = updates
	}

	_, err = client.UpdateMySQLDatabase(ctx, int(id), updateOpts)
	if err != nil {
		return diag.Errorf("failed to update mysql database: %s", err)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MySQL database ID %s as int: %s", d.Id(), err)
	}

	// We should retry on intermittent deletion errors
	return diag.FromErr(resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		err := client.DeleteMySQLDatabase(ctx, int(id))
		if err != nil {
			if lerr, ok := err.(*linodego.Error); ok &&
				lerr.Code == 500 && strings.Contains(lerr.Message, "Unable to delete instance") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(fmt.Errorf("failed to delete mysql database %d: %s", id, err))
		}

		return nil
	}))
}

func createEngineSlug(engine, version string) string {
	return fmt.Sprintf("%s/%s", engine, version)
}

func FlattenMaintenanceWindow(window linodego.MySQLDatabaseMaintenanceWindow) map[string]interface{} {
	result := make(map[string]interface{})

	result["day_of_week"] = helper.FlattenDayOfWeek(window.DayOfWeek)
	result["duration"] = window.Duration
	result["frequency"] = string(window.Frequency)
	result["hour_of_day"] = window.HourOfDay

	// Nullable
	if window.WeekOfMonth != nil {
		result["week_of_month"] = window.WeekOfMonth
	}

	return result
}

func ExpandMaintenanceWindow(window map[string]interface{}) (linodego.MySQLDatabaseMaintenanceWindow, error) {
	result := linodego.MySQLDatabaseMaintenanceWindow{
		Duration:    window["duration"].(int),
		Frequency:   linodego.DatabaseMaintenanceFrequency(window["frequency"].(string)),
		HourOfDay:   window["hour_of_day"].(int),
		WeekOfMonth: nil,
	}

	dayOfWeek, err := helper.ExpandDayOfWeek(window["day_of_week"].(string))
	if err != nil {
		return result, err
	}
	result.DayOfWeek = dayOfWeek

	if val, ok := window["week_of_month"]; ok && val.(int) > 0 {
		valInt := val.(int)
		result.WeekOfMonth = &valInt
	}

	return result, nil
}
