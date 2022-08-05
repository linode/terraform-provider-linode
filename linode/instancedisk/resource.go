package instancedisk

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

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
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
	}
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID := d.Get("linode_id").(int)

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
	client := meta.(*helper.ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)

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

	p, err := client.NewEventPoller(ctx, linodeID, linodego.EntityLinode, linodego.ActionDiskCreate)
	if err != nil {
		return diag.Errorf("failed to poll for events: %s", err)
	}

	disk, err := client.CreateInstanceDisk(ctx, linodeID, createOpts)
	if err != nil {
		return diag.Errorf("failed to create linode instance disk: %s", err)
	}

	d.SetId(strconv.Itoa(disk.ID))

	if _, err := p.WaitForFinished(ctx, helper.GetDeadlineSeconds(ctx, d)); err != nil {
		return diag.Errorf("failed to wait for instance shutdown: %s", err)
	}

	if _, err := client.WaitForInstanceDiskStatus(
		ctx, linodeID, disk.ID, linodego.DiskReady, helper.GetDeadlineSeconds(ctx, d)); err != nil {
		return diag.Errorf("failed ot wait for disk ready: %s", err)
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID := d.Get("linode_id").(int)

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
		if _, err := client.UpdateInstanceDisk(ctx, linodeID, id, putRequest); err != nil {
			return diag.Errorf("failed to update instance disk: %s", err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID := d.Get("linode_id").(int)

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
		log.Printf("[INFO] Shutting down instance %d for disk %d deletion", linodeID, id)

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
	}

	err = client.DeleteInstanceDisk(ctx, linodeID, id)
	if err != nil {
		return diag.Errorf("Error deleting Linode Instance Disk %d: %s", id, err)
	}

	d.SetId("")

	// Reboot the instance if necessary
	if shouldShutdown && !diskInConfig {
		log.Printf("[INFO] Booting instance %d to config %d", linodeID, configID)

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
	}

	return nil
}
