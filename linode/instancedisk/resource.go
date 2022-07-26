package instancedisk

import (
	"context"
	"log"
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
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	ids, err := helper.ParseMultiSegmentID(d.Id(), 2)
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID, id := ids[0], ids[1]

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

	disk, err := client.CreateInstanceDisk(ctx, linodeID, createOpts)
	if err != nil {
		return diag.Errorf("failed to create linode instance disk: %s", err)
	}

	d.SetId(helper.FormatMultiSegmentID(linodeID, disk.ID))

	// Wait for the resize event to complete
	_, err = client.WaitForEventFinished(ctx, linodeID, linodego.EntityLinode, linodego.ActionDiskCreate,
		*disk.Updated, helper.GetDeadlineSeconds(ctx, d))
	if err != nil {
		return diag.Errorf("failed to wait for disk creation: %s", err)
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	ids, err := helper.ParseMultiSegmentID(d.Id(), 2)
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID, id := ids[0], ids[1]

	disk, err := client.GetInstanceDisk(ctx, linodeID, id)
	if err != nil {
		return diag.Errorf("failed to get instance disk: %s", err)
	}

	if d.HasChange("size") {
		newSize := d.Get("size").(int)

		if err := client.ResizeInstanceDisk(ctx, linodeID, id, newSize); err != nil {
			return diag.Errorf("failed to resize disk: %s", err)
		}

		// Wait for the resize event to complete
		_, err := client.WaitForEventFinished(ctx, linodeID, linodego.EntityLinode, linodego.ActionDiskResize,
			*disk.Updated, helper.GetDeadlineSeconds(ctx, d))
		if err != nil {
			return diag.Errorf("failed to resize disk: %s", err)
		}

		// Check to see if the resize operation worked
		if updatedDisk, err := client.WaitForInstanceDiskStatus(ctx, linodeID, disk.ID, linodego.DiskReady,
			helper.GetDeadlineSeconds(ctx, d)); err != nil {
			return diag.Errorf("failed to wait for disk ready: %s", err)
		} else if updatedDisk.Size != newSize {
			return diag.Errorf(
				"failed to resize disk %d from %d to %d", disk.ID, disk.Size, newSize)
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
	ids, err := helper.ParseMultiSegmentID(d.Id(), 2)
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID, id := ids[0], ids[1]

	err = client.DeleteInstanceDisk(ctx, linodeID, id)
	if err != nil {
		return diag.Errorf("Error deleting Linode Instance Disk %d: %s", id, err)
	}

	d.SetId("")
	return nil
}
