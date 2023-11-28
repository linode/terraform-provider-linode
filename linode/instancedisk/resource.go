package instancedisk

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
			StateContext: importResource,
		},
	}
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	tflog.Debug(ctx, "Import linode_instance_disk", map[string]any{
		"id": d.Id(),
	})

	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		// Validate that this is an ID by making sure it can be converted into an int
		_, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("invalid disk ID: %v", err)
		}

		instID, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, fmt.Errorf("invalid inst ID: %v", err)
		}

		d.SetId(s[1])
		d.Set("linode_id", instID)
	}

	err := readResource(ctx, d, meta)
	if err != nil {
		return nil, fmt.Errorf("unable to import %v as instance disk: %v", d.Id(), err)
	}

	results := make([]*schema.ResourceData, 0)
	results = append(results, d)

	return results, nil
}

func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Read linode_instance_disk")

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID := d.Get("linode_id").(int)

	client := meta.(*helper.ProviderMeta).Client

	disk, err := client.GetInstanceDisk(ctx, linodeID, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Instance Disk ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Instance Disk: %s", err)
	}

	d.Set("linode_id", linodeID)
	d.Set("created", disk.Created.Format(time.RFC3339))
	d.Set("filesystem", disk.Filesystem)
	d.Set("label", disk.Label)
	d.Set("size", disk.Size)
	d.Set("status", disk.Status)
	d.Set("updated", disk.Updated.Format(time.RFC3339))

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Create linode_instance_disk")

	linodeID := d.Get("linode_id").(int)

	client := meta.(*helper.ProviderMeta).Client

	createOpts := linodego.InstanceDiskCreateOptions{
		AuthorizedKeys:  helper.ExpandStringSet(d.Get("authorized_keys").(*schema.Set)),
		AuthorizedUsers: helper.ExpandStringSet(d.Get("authorized_users").(*schema.Set)),
		Filesystem:      d.Get("filesystem").(string),
		Image:           d.Get("image").(string),
		Label:           d.Get("label").(string),
		RootPass:        d.Get("root_pass").(string),
		Size:            d.Get("size").(int),
		StackscriptID:   d.Get("stackscript_id").(int),
	}

	if stackscriptData, ok := d.GetOk("stackscript_data"); ok {
		createOpts.StackscriptData = expandStackScriptData(stackscriptData)
	}

	createOpts.RootPass = d.Get("root_pass").(string)
	if createOpts.RootPass == "" {
		var err error
		createOpts.RootPass, err = helper.CreateRandomRootPassword()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	p, err := client.NewEventPoller(ctx, linodeID, linodego.EntityLinode, linodego.ActionDiskCreate)
	if err != nil {
		return diag.Errorf("failed to poll for events: %s", err)
	}

	disk, err := client.CreateInstanceDisk(ctx, linodeID, createOpts)
	if err != nil {
		return diag.Errorf("failed to create linode instance disk: %s", err)
	}

	ctx = tflog.SetField(ctx, "disk_id", disk.ID)
	tflog.Info(ctx, "Created Instance Disk; waiting for creation to finish", map[string]any{
		"body": createOpts,
	})

	d.SetId(strconv.Itoa(disk.ID))

	event, err := p.WaitForFinished(ctx, helper.GetDeadlineSeconds(ctx, d))
	if err != nil {
		return diag.Errorf("failed to wait for instance shutdown: %s", err)
	}

	if _, err := client.WaitForInstanceDiskStatus(
		ctx, linodeID, disk.ID, linodego.DiskReady, helper.GetDeadlineSeconds(ctx, d)); err != nil {
		return diag.Errorf("failed ot wait for disk ready: %s", err)
	}

	tflog.Debug(ctx, "Instance disk is ready", map[string]any{
		"event_id": event.ID,
	})

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Update linode_instance_disk")

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID := d.Get("linode_id").(int)

	client := meta.(*helper.ProviderMeta).Client

	if d.HasChange("size") {
		err = handleDiskResize(ctx, client, linodeID, id, d.Get("size").(int), helper.GetDeadlineSeconds(ctx, d))
		if err != nil {
			return diag.Errorf("failed to resize disk: %s", err)
		}
	}

	putRequest := linodego.InstanceDiskUpdateOptions{}
	shouldUpdate := false

	if d.HasChange("label") {
		putRequest.Label = d.Get("label").(string)
		shouldUpdate = true
	}

	if shouldUpdate {
		tflog.Debug(ctx, "Update Instance disk", map[string]any{
			"body": putRequest,
		})

		if _, err := client.UpdateInstanceDisk(ctx, linodeID, id, putRequest); err != nil {
			return diag.Errorf("failed to update instance disk: %s", err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Delete linode_instance_disk")

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID := d.Get("linode_id").(int)

	client := meta.(*helper.ProviderMeta).Client

	configID, err := helper.GetCurrentBootedConfig(ctx, &client, linodeID)
	if err != nil {
		return diag.Errorf("failed to get current booted config: %s", err)
	}

	isDiskInConfig := func() (bool, error) {
		if configID == 0 {
			return false, nil
		}

		cfg, err := client.GetInstanceConfig(ctx, linodeID, configID)
		if err != nil {
			return false, err
		}

		if cfg.Devices == nil {
			return false, nil
		}

		reflectMap := reflect.ValueOf(*cfg.Devices)

		for i := 0; i < reflectMap.NumField(); i++ {
			field := reflectMap.Field(i).Interface().(*linodego.InstanceConfigDevice)
			if field == nil {
				continue
			}

			if field.DiskID == id {
				return true, nil
			}
		}

		return false, nil
	}

	shouldShutdown := configID != 0
	diskInConfig, err := isDiskInConfig()
	if err != nil {
		return diag.Errorf("failed to check if disk is in use: %s", err)
	}

	// Shutdown instance if active
	if shouldShutdown {
		tflog.Info(ctx, "Shutting down instance for disk deletion")

		p, err := client.NewEventPoller(ctx, linodeID, linodego.EntityLinode, linodego.ActionLinodeShutdown)
		if err != nil {
			return diag.Errorf("failed to poll for events: %s", err)
		}

		if err := client.ShutdownInstance(ctx, linodeID); err != nil {
			return diag.Errorf("failed to shutdown instance: %s", err)
		}

		if _, err := p.WaitForFinished(ctx, helper.GetDeadlineSeconds(ctx, d)); err != nil {
			return diag.Errorf("failed to wait for instance shutdown: %s", err)
		}

		tflog.Debug(ctx, "Instance shutdown event finished")
	}

	tflog.Info(ctx, "Deleting instance disk")
	p, err := client.NewEventPollerWithSecondary(
		ctx,
		linodeID,
		linodego.EntityLinode,
		id,
		linodego.ActionDiskDelete,
	)
	if err != nil {
		return diag.Errorf("failed to initialize event poller: %s", err)
	}

	if err := client.DeleteInstanceDisk(ctx, linodeID, id); err != nil {
		return diag.Errorf("Error deleting Linode Instance Disk %d: %s", id, err)
	}

	if _, err := p.WaitForFinished(ctx, helper.GetDeadlineSeconds(ctx, d)); err != nil {
		return diag.Errorf(
			"Error waiting for Instance %d Disk %d to finish deleting: %s", linodeID, id, err)
	}

	d.SetId("")

	// Reboot the instance if necessary
	if shouldShutdown && !diskInConfig {
		tflog.Info(ctx, "Booting instance into previously booted config", map[string]any{
			"config_id": configID,
		})

		p, err := client.NewEventPoller(ctx, linodeID, linodego.EntityLinode, linodego.ActionLinodeBoot)
		if err != nil {
			return diag.Errorf("failed to poll for events: %s", err)
		}

		if err := client.BootInstance(ctx, linodeID, configID); err != nil {
			return diag.Errorf("failed to boot instance %d %d: %s", linodeID, configID, err)
		}

		if _, err := p.WaitForFinished(ctx, helper.GetDeadlineSeconds(ctx, d)); err != nil {
			return diag.Errorf("failed to wait for instance boot: %s", err)
		}

		tflog.Debug(ctx, "Instance boot event finished")
	}

	return nil
}

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"linode_id": d.Get("linode_id").(int),
		"id":        d.Id(),
	})
}
